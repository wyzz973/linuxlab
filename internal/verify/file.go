package verify

import (
	"fmt"
	"os"

	"github.com/sd3/linuxlab/internal/challenge"
)

// FileContentVerifier reads a file and compares its trimmed content to the expected value.
type FileContentVerifier struct{}

func (v *FileContentVerifier) Verify(rule challenge.VerifyRule) Result {
	data, err := os.ReadFile(rule.Path)
	if err != nil {
		return Result{Passed: false, Message: fmt.Sprintf("无法读取文件 %s: %v", rule.Path, err)}
	}
	return compareStrings(string(data), rule.Expect, "文件内容")
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
