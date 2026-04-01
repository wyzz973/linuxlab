package verify

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/sd3/linuxlab/internal/challenge"
)

// ScriptVerifier runs a shell script and checks its exit code (0 = pass).
type ScriptVerifier struct{}

func (v *ScriptVerifier) Verify(rule challenge.VerifyRule) Result {
	cmd := exec.Command("bash", rule.Path)
	output, err := cmd.CombinedOutput()

	msg := strings.TrimSpace(string(output))

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return Result{
				Passed:  false,
				Message: fmt.Sprintf("脚本检测未通过 (exit %d): %s", exitErr.ExitCode(), msg),
			}
		}
		return Result{Passed: false, Message: fmt.Sprintf("脚本执行失败: %v", err)}
	}

	if msg == "" {
		msg = "脚本检测通过"
	}
	return Result{Passed: true, Message: msg}
}
