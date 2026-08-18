package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/sd3/linuxlab/internal/challenge"
	"github.com/sd3/linuxlab/internal/progress"
	"github.com/sd3/linuxlab/internal/reference"
	"github.com/sd3/linuxlab/internal/runner"
)

// exitVerifyFailed is returned when the challenge ran but verification did not
// pass, distinct from exit code 1 which means an execution error.
const exitVerifyFailed = 2

func runChallenge(ctx context.Context, env Env) int {
	if len(env.Args) < 2 || env.Args[1] != "run" {
		return writeError(env.Stderr, "用法: linuxlab challenge run <challenge-id> [--json] [--hints N]")
	}
	if len(env.Args) < 3 {
		return writeError(env.Stderr, "缺少题目 ID")
	}

	challengeID := env.Args[2]
	jsonMode := hasFlag(env.Args[3:], "--json")
	ch, err := findChallenge(env.ChallengesDir, challengeID)
	if err != nil {
		return writeError(env.Stderr, err.Error())
	}

	// 提示使用次数与 Go TUI 的渐进式提示语义一致：已在详情页解锁的
	// 提示数透传到执行层，记入进度并影响得分（ScoreWithHints）。
	hints := hintsFromArgs(env.Args[3:])
	if hints > len(ch.Hints) {
		hints = len(ch.Hints)
	}

	refs := loadRefs(env.RefsPath)
	emit := func(event runner.Event) {
		if jsonMode {
			_ = json.NewEncoder(env.Stdout).Encode(event)
			return
		}
		renderEvent(env.Stdout, event)
	}

	runStdout := env.Stdout
	if jsonMode {
		runStdout = env.Stderr
	}

	result, err := runner.RunInteractive(ctx, runner.Options{
		Challenge: ch,
		Refs:      refs,
		HintsUsed: hints,
		Emit:      emit,
		Stdin:     env.Stdin,
		Stdout:    runStdout,
		Stderr:    env.Stderr,
	})
	if err != nil {
		if jsonMode {
			_ = json.NewEncoder(env.Stdout).Encode(runner.Event{Type: "error", Message: err.Error()})
		}
		return writeError(env.Stderr, fmt.Sprintf("挑战执行失败: %v", err))
	}

	recordProgress(env.ProgressPath, ch, result)
	if !result.Passed {
		return exitVerifyFailed
	}
	return 0
}

// renderEvent prints a runner event in plain-text mode. Result events get a
// dedicated rendering so the user sees pass/fail status and every check message.
func renderEvent(w io.Writer, event runner.Event) {
	if event.Type == "result" {
		if event.Passed != nil && *event.Passed {
			fmt.Fprintln(w, "检测结果: 通过")
		} else {
			fmt.Fprintln(w, "检测结果: 未通过")
		}
		for i, r := range event.Results {
			icon := "✓"
			if !r.Passed {
				icon = "✗"
			}
			fmt.Fprintf(w, "%s 检查 %d: %s\n", icon, i+1, r.Message)
		}
		if event.HintsUsed > 0 {
			fmt.Fprintf(w, "使用提示: %d\n", event.HintsUsed)
		}
		return
	}
	if event.Message != "" {
		fmt.Fprintln(w, event.Message)
	}
}

func findChallenge(challengesDir, id string) (*challenge.Challenge, error) {
	challenges, err := challenge.LoadAll(challengesDir)
	if err != nil {
		return nil, fmt.Errorf("加载题库失败: %w", err)
	}
	for _, ch := range challenges {
		if ch.ID == id {
			return ch, nil
		}
	}
	return nil, fmt.Errorf("未找到题目: %s", id)
}

// hintsFromArgs extracts the --hints N flag value. Malformed or missing
// values yield 0; the caller clamps the result to the challenge's hint count.
func hintsFromArgs(args []string) int {
	for i, arg := range args {
		if arg == "--hints" && i+1 < len(args) {
			if n, err := strconv.Atoi(args[i+1]); err == nil && n >= 0 {
				return n
			}
		}
	}
	return 0
}

func loadRefs(path string) *reference.ReferenceData {
	if path == "" {
		return nil
	}
	refs, err := reference.LoadReferences(path)
	if err != nil {
		return nil
	}
	return refs
}

func recordProgress(progressPath string, ch *challenge.Challenge, result runner.Result) {
	if progressPath == "" {
		return
	}
	store, err := progress.NewStore(progressPath)
	if err != nil {
		return
	}
	store.RecordAttempt(ch.ID, ch.Category, ch.Subcategory, result.Passed, result.HintsUsed)
	_ = store.Save()
}
