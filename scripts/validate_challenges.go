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
	"os/exec"
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
			// containers 类别在宿主上操作 docker（LocalSandbox/ComposeSandbox
			// 语义，与 runner 一致）：init/solution/check 都在宿主的题目目录
			// 执行，因为沙盒容器里没有 docker CLI。
			if ch.Category == "containers" {
				result := validateHostChallenge(ch)
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

// execFunc runs a command string and returns output, exit code, and error.
// The docker-backed and host-backed validation paths share it.
type execFunc func(ctx context.Context, command string) (string, int, error)

// challengeTimeout bounds one challenge's init → solution → verify sequence.
// 120s was too tight for fresh containers whose init runs apt-get (mirror
// latency under parallel validation routinely exceeded it, killing the
// challenge with exit -1); 300s keeps slow-apt challenges green without
// letting genuinely blocking commands run forever.
const challengeTimeout = 300 * time.Second

// validationImage pre-warms every tool the challenge init scripts install
// (see Dockerfile.sandbox), so apt inside validation containers is a no-op
// and runs are fast and stable regardless of mirror latency. Build it with:
//
//	docker build -t linuxlab/sandbox:22.04 -f Dockerfile.sandbox .
const validationImage = "linuxlab/sandbox:22.04"

func validateChallenge(ch *challenge.Challenge) validateResult {
	ctx, cancel := context.WithTimeout(context.Background(), challengeTimeout)
	defer cancel()

	// Create a fresh container for each challenge
	sb, err := sandbox.NewDockerSandbox(ctx, validationImage)
	if err != nil {
		return validateResult{"FAIL", fmt.Sprintf("创建沙盒失败: %v", err)}
	}
	// Destroy with a fresh context: the challenge ctx may already be expired
	// (120s timeout), which would make cleanup fail and leak the container.
	defer func() {
		dctx, dcancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer dcancel()
		if err := sb.Destroy(dctx); err != nil {
			fmt.Fprintf(os.Stderr, "  清理容器失败 [%s]: %v\n", ch.ID, err)
		}
	}()

	// Ensure /home/learner exists before init.sh runs
	sb.Exec(ctx, "mkdir -p /home/learner")

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
		// Quoted heredoc ('SETUP_EOF') prevents shell expansion, so no escaping needed.
		// Trim trailing newline to avoid an extra blank line before the delimiter.
		content := strings.TrimRight(sf.Content, "\n")
		sb.Exec(ctx, "cat > '"+sf.Path+"' << 'SETUP_EOF'\n"+content+"\nSETUP_EOF")
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
		// Use quoted heredoc delimiter so content is taken literally — no escaping needed
		sb.Exec(ctx, "cat > /home/learner/solution.sh << 'SOL_EOF'\n"+string(solData)+"\nSOL_EOF")
		sb.Exec(ctx, "chmod +x /home/learner/solution.sh")
	}

	// Run solution.sh
	sb.Exec(ctx, string(solData)) // ignore exit code — some solutions have intentional non-zero parts

	return verifyRules(ch, ctx, sb.Exec)
}

// validateHostChallenge validates a containers-category challenge on the host,
// matching the runner's LocalSandbox/ComposeSandbox semantics: init.sh,
// solution.sh and the verify rules all execute in the challenge directory,
// because the sandbox container has neither the docker CLI nor the daemon
// socket (the exercises themselves are host docker operations).
func validateHostChallenge(ch *challenge.Challenge) validateResult {
	ctx, cancel := context.WithTimeout(context.Background(), challengeTimeout)
	defer cancel()

	execFn := func(ctx context.Context, command string) (string, int, error) {
		cmd := exec.CommandContext(ctx, "/bin/sh", "-c", command)
		cmd.Dir = ch.Dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				return string(out), exitErr.ExitCode(), nil
			}
			return string(out), -1, err
		}
		return string(out), 0, nil
	}

	before := runningContainerIDs()

	// 宿主脚本以文件方式执行（bash <path>，工作目录 = 题目目录），而不是
	// 把内容传给 sh -c：后者会让 $0 变成 shell 路径（如 /bin/sh），导致
	// 脚本里 `cd "$(dirname "$0")"` 解析到错误目录。
	// 路径必须绝对化：ch.Dir 是相对路径，相对路径会被 bash 再相对 cmd.Dir
	// 解析，导致双重拼接找不到文件。
	runScript := func(path string) (string, int) {
		abs, err := filepath.Abs(path)
		if err != nil {
			return "", -1
		}
		cmd := exec.CommandContext(ctx, "bash", abs)
		cmd.Dir = ch.Dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				return string(out), exitErr.ExitCode()
			}
			return string(out), -1
		}
		return string(out), 0
	}

	if initPath := filepath.Join(ch.Dir, "init.sh"); fileExists(initPath) {
		runScript(initPath)
	}

	solutionPath := filepath.Join(ch.Dir, "solution.sh")
	if !fileExists(solutionPath) {
		return validateResult{"SKIP", "无 solution.sh"}
	}
	runScript(solutionPath) // ignore exit code — some solutions have intentional non-zero parts

	// 清理本挑战启动的资源（compose 项目 down + 新增容器删除），
	// 避免批量验证在宿主上留下运行中的容器。
	defer cleanupHostContainers(before, ch)

	return verifyRules(ch, ctx, execFn)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func runningContainerIDs() []string {
	out, err := exec.Command("docker", "ps", "-aq").Output()
	if err != nil {
		return nil
	}
	return strings.Fields(string(out))
}

func cleanupHostContainers(before []string, ch *challenge.Challenge) {
	// compose 挑战整体下架（同时移除其容器）
	if ch.ComposeFile != "" {
		cmd := exec.Command("docker", "compose", "-f", filepath.Join(ch.Dir, ch.ComposeFile), "down")
		cmd.Dir = ch.Dir
		_ = cmd.Run()
	}
	for _, id := range runningContainerIDs() {
		found := false
		for _, old := range before {
			if id == old {
				found = true
				break
			}
		}
		if !found {
			_ = exec.Command("docker", "rm", "-f", id).Run()
		}
	}
}

func verifyRules(ch *challenge.Challenge, ctx context.Context, execFn execFunc) validateResult {
	// A challenge without any verify rule must not pass silently.
	if len(ch.Verify) == 0 {
		return validateResult{"FAIL", "无验证规则: challenge.yaml 没有任何 verify 规则"}
	}
	for i, rule := range ch.Verify {
		var passed bool
		var msg string

		switch rule.Type {
		case "file_content":
			out, _, _ := execFn(ctx, "cat '"+rule.Path+"' 2>/dev/null")
			actual := strings.TrimSpace(out)
			expected := strings.TrimSpace(rule.Expect)
			passed = actual == expected
			if !passed {
				msg = fmt.Sprintf("verify[%d] file_content 不匹配\n    期望: %q\n    实际: %q", i, truncate(expected, 80), truncate(actual, 80))
			}

		case "file_exists":
			_, code, _ := execFn(ctx, "test -e '"+rule.Path+"'")
			passed = code == 0
			if !passed {
				msg = fmt.Sprintf("verify[%d] 路径不存在: %s", i, rule.Path)
			}

		case "command_output":
			out, _, execErr := execFn(ctx, rule.Command)
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
			_, code, _ := execFn(ctx, rule.Command)
			expected := strings.TrimSpace(rule.Expect)
			passed = fmt.Sprintf("%d", code) == expected
			if !passed {
				msg = fmt.Sprintf("verify[%d] exit_code 不匹配: 期望 %s, 实际 %d", i, expected, code)
			}

		case "permissions":
			out, _, _ := execFn(ctx, "stat -c '%a' '"+rule.Path+"' 2>/dev/null")
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
				out, code, _ := execFn(ctx, string(scriptData))
				passed = code == 0
				if !passed {
					msg = fmt.Sprintf("verify[%d] check.sh 退出码: %d\n    %s", i, code, truncate(strings.TrimSpace(out), 300))
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
