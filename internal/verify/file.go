package verify

import (
	"fmt"
	"os"
	"strings"

	"github.com/sd3/linuxlab/internal/challenge"
)

// FileContentVerifier reads a file and compares its trimmed content to the expected value.
type FileContentVerifier struct{}

func (v *FileContentVerifier) Verify(rule challenge.VerifyRule) Result {
	data, err := os.ReadFile(rule.Path)
	if err != nil {
		return Result{Passed: false, Message: fmt.Sprintf("无法读取文件 %s: %v", rule.Path, err)}
	}

	actual := strings.TrimSpace(string(data))
	expected := strings.TrimSpace(rule.Expect)

	if actual != expected {
		return Result{
			Passed:  false,
			Message: fmt.Sprintf("文件内容不匹配\n期望: %q\n实际: %q", expected, actual),
		}
	}
	return Result{Passed: true, Message: fmt.Sprintf("文件内容匹配: %s", rule.Path)}
}

// FileExistsVerifier checks whether a file or directory exists at the given path.
type FileExistsVerifier struct{}

func (v *FileExistsVerifier) Verify(rule challenge.VerifyRule) Result {
	_, err := os.Stat(rule.Path)
	if err != nil {
		return Result{Passed: false, Message: fmt.Sprintf("路径不存在: %s", rule.Path)}
	}
	return Result{Passed: true, Message: fmt.Sprintf("路径存在: %s", rule.Path)}
}
