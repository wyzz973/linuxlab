// +build ignore

// validate_challenges.go — batch-validates all challenges in a Docker container.
// Usage: go run scripts/validate_challenges.go [category]
// Example: go run scripts/validate_challenges.go linux-basics
//          go run scripts/validate_challenges.go          (all categories)

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sd3/linuxlab/internal/challenge"
	"github.com/sd3/linuxlab/internal/sandbox"
)

func main() {
	filterCat := ""
	if len(os.Args) > 1 {
		filterCat = os.Args[1]
	}

	byCategory, err := challenge.LoadAllByCategory("challenges")
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载题库失败: %v\n", err)
		os.Exit(1)
	}

	if !sandbox.DockerAvailable() {
		fmt.Fprintf(os.Stderr, "Docker 不可用，无法测试\n")
		os.Exit(1)
	}

	totalPass := 0
	totalFail := 0
	totalSkip := 0
	var failures []string

	for cat, challenges := range byCategory {
		if filterCat != "" && cat != filterCat {
			continue
		}

		fmt.Printf("\n══════ %s (%d 题) ══════\n", cat, len(challenges))

		for _, ch := range challenges {
			// Skip container challenges (need compose or special setup)
			if ch.Category == "containers" {
				fmt.Printf("  ⊘ %-40s SKIP (container challenge)\n", ch.Title)
				totalSkip++
				continue
			}

			result := validateChallenge(ch)
			switch result.status {
			case "PASS":
				fmt.Printf("  ✓ %-40s PASS\n", ch.Title)
				totalPass++
			case "FAIL":
				fmt.Printf("  ✗ %-40s FAIL: %s\n", ch.Title, result.message)
				totalFail++
				failures = append(failures, fmt.Sprintf("[%s] %s: %s", ch.Category, ch.Title, result.message))
			case "SKIP":
				fmt.Printf("  ⊘ %-40s SKIP: %s\n", ch.Title, result.message)
				totalSkip++
			}
		}
	}

	fmt.Printf("\n══════ 结果汇总 ══════\n")
	fmt.Printf("  通过: %d\n", totalPass)
	fmt.Printf("  失败: %d\n", totalFail)
	fmt.Printf("  跳过: %d\n", totalSkip)
	fmt.Printf("  总计: %d\n", totalPass+totalFail+totalSkip)

	if len(failures) > 0 {
		fmt.Printf("\n══════ 失败详情 ══════\n")
		for _, f := range failures {
			fmt.Printf("  ✗ %s\n", f)
		}
		os.Exit(1)
	}
}

type validateResult struct {
	status  string // PASS, FAIL, SKIP
	message string
}

func validateChallenge(ch *challenge.Challenge) validateResult {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// Create a fresh container for each challenge
	sb, err := sandbox.NewDockerSandbox(ctx, "ubuntu:22.04")
	if err != nil {
		return validateResult{"FAIL", fmt.Sprintf("创建沙盒失败: %v", err)}
	}
	defer sb.Destroy(ctx)

	// Run init.sh if exists (ignore exit code for process-management challenges)
	initPath := filepath.Join(ch.Dir, "init.sh")
	if data, err := os.ReadFile(initPath); err == nil {
		sb.Exec(ctx, string(data))
	}

	// Write setup_files into the container (for Vim challenges)
	for _, sf := range ch.SetupFiles {
		dir := filepath.Dir(sf.Path)
		if dir != "." && dir != "/" {
			sb.Exec(ctx, "mkdir -p '"+dir+"'")
		}
		escaped := strings.ReplaceAll(sf.Content, "'", "'\"'\"'")
		sb.Exec(ctx, "cat > '"+sf.Path+"' << 'SETUP_EOF'\n"+escaped+"\nSETUP_EOF")
	}

	// Run solution.sh
	solutionPath := filepath.Join(ch.Dir, "solution.sh")
	solData, err := os.ReadFile(solutionPath)
	if err != nil {
		return validateResult{"SKIP", "无 solution.sh"}
	}

	// For shell-scripting challenges, copy solution.sh to /home/learner/solution.sh
	// because check.sh typically runs `bash /home/learner/solution.sh`
	if ch.Category == "shell-scripting" {
		sb.Exec(ctx, "mkdir -p /home/learner")
		escaped := strings.ReplaceAll(string(solData), "'", "'\"'\"'")
		sb.Exec(ctx, "cat > /home/learner/solution.sh << 'SOL_EOF'\n"+escaped+"\nSOL_EOF")
		sb.Exec(ctx, "chmod +x /home/learner/solution.sh")
	}

	// Run solution.sh
	sb.Exec(ctx, string(solData)) // ignore exit code — some solutions have intentional non-zero parts

	// Run verification inside the container
	for i, rule := range ch.Verify {
		var passed bool
		var msg string

		switch rule.Type {
		case "file_content":
			out, _, _ := sb.Exec(ctx, "cat '"+rule.Path+"' 2>/dev/null")
			actual := strings.TrimSpace(out)
			expected := strings.TrimSpace(rule.Expect)
			passed = actual == expected
			if !passed {
				msg = fmt.Sprintf("verify[%d] file_content 不匹配\n    期望: %q\n    实际: %q", i, truncate(expected, 80), truncate(actual, 80))
			}

		case "file_exists":
			_, code, _ := sb.Exec(ctx, "test -e '"+rule.Path+"'")
			passed = code == 0
			if !passed {
				msg = fmt.Sprintf("verify[%d] 路径不存在: %s", i, rule.Path)
			}

		case "command_output":
			out, _, execErr := sb.Exec(ctx, rule.Command)
			if execErr != nil {
				return validateResult{"FAIL", fmt.Sprintf("verify[%d] 命令错误: %v", i, execErr)}
			}
			actual := strings.TrimSpace(out)
			expected := strings.TrimSpace(rule.Expect)
			passed = actual == expected
			if !passed {
				msg = fmt.Sprintf("verify[%d] command_output 不匹配\n    命令: %s\n    期望: %q\n    实际: %q", i, rule.Command, truncate(expected, 80), truncate(actual, 80))
			}

		case "exit_code":
			_, code, _ := sb.Exec(ctx, rule.Command)
			expected := strings.TrimSpace(rule.Expect)
			passed = fmt.Sprintf("%d", code) == expected
			if !passed {
				msg = fmt.Sprintf("verify[%d] exit_code 不匹配: 期望 %s, 实际 %d", i, expected, code)
			}

		case "permissions":
			out, _, _ := sb.Exec(ctx, "stat -c '%a' '"+rule.Path+"' 2>/dev/null")
			actual := strings.TrimSpace(out)
			expected := strings.TrimSpace(rule.Expect)
			passed = actual == expected
			if !passed {
				msg = fmt.Sprintf("verify[%d] permissions 不匹配: 期望 %s, 实际 %s", i, expected, actual)
			}

		case "script":
			scriptPath := rule.Path
			if !filepath.IsAbs(scriptPath) {
				scriptPath = filepath.Join(ch.Dir, scriptPath)
			}
			if scriptData, readErr := os.ReadFile(scriptPath); readErr == nil {
				_, code, _ := sb.Exec(ctx, string(scriptData))
				passed = code == 0
				if !passed {
					msg = fmt.Sprintf("verify[%d] check.sh 退出码: %d", i, code)
				}
			} else {
				return validateResult{"FAIL", fmt.Sprintf("verify[%d] 无法读取脚本: %v", i, readErr)}
			}

		default:
			return validateResult{"SKIP", fmt.Sprintf("未知验证类型: %s", rule.Type)}
		}

		if !passed {
			return validateResult{"FAIL", msg}
		}
	}

	return validateResult{"PASS", ""}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
