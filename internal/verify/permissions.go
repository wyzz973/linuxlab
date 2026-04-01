package verify

import (
	"fmt"
	"os"
	"strings"

	"github.com/sd3/linuxlab/internal/challenge"
)

// PermissionsVerifier checks file permissions against expected octal string.
type PermissionsVerifier struct{}

func (v *PermissionsVerifier) Verify(rule challenge.VerifyRule) Result {
	info, err := os.Stat(rule.Path)
	if err != nil {
		return Result{Passed: false, Message: fmt.Sprintf("无法读取文件: %v", err)}
	}

	actual := fmt.Sprintf("%o", info.Mode().Perm())
	expected := strings.TrimSpace(rule.Expect)

	if actual != expected {
		return Result{
			Passed:  false,
			Message: fmt.Sprintf("权限不匹配\n期望: %s\n实际: %s", expected, actual),
		}
	}
	return Result{Passed: true, Message: fmt.Sprintf("权限匹配: %s", expected)}
}
