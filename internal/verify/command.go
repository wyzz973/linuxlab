package verify

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/sd3/linuxlab/internal/challenge"
)

// CommandOutputVerifier runs a command and compares its trimmed stdout to the expected value.
type CommandOutputVerifier struct{}

func (v *CommandOutputVerifier) Verify(rule challenge.VerifyRule) Result {
	cmd := exec.Command("sh", "-c", rule.Command)
	out, err := cmd.Output()
	if err != nil {
		return Result{Passed: false, Message: fmt.Sprintf("命令执行失败: %v", err)}
	}

	actual := strings.TrimSpace(string(out))
	expected := strings.TrimSpace(rule.Expect)

	if actual != expected {
		return Result{
			Passed:  false,
			Message: fmt.Sprintf("命令输出不匹配\n期望: %q\n实际: %q", expected, actual),
		}
	}
	return Result{Passed: true, Message: fmt.Sprintf("命令输出匹配: %s", rule.Command)}
}

// ExitCodeVerifier runs a command and compares its exit code to the expected value.
type ExitCodeVerifier struct{}

func (v *ExitCodeVerifier) Verify(rule challenge.VerifyRule) Result {
	expectedCode, err := strconv.Atoi(strings.TrimSpace(rule.Expect))
	if err != nil {
		return Result{Passed: false, Message: fmt.Sprintf("无效的期望退出码: %s", rule.Expect)}
	}

	cmd := exec.Command("sh", "-c", rule.Command)
	err = cmd.Run()

	actualCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			actualCode = exitErr.ExitCode()
		} else {
			return Result{Passed: false, Message: fmt.Sprintf("命令执行失败: %v", err)}
		}
	}

	if actualCode != expectedCode {
		return Result{
			Passed:  false,
			Message: fmt.Sprintf("退出码不匹配\n期望: %d\n实际: %d", expectedCode, actualCode),
		}
	}
	return Result{Passed: true, Message: fmt.Sprintf("退出码匹配: %d", actualCode)}
}
