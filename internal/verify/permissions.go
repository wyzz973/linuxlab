package verify

import (
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"

	"github.com/sd3/linuxlab/internal/challenge"
)

// PermissionsVerifier checks file permissions against an expected octal
// string. Expected values are compared numerically, so "644" and "0644" are
// equivalent, and setuid/setgid/sticky bits ("4755", "1777") are supported.
type PermissionsVerifier struct{}

func (v *PermissionsVerifier) Verify(rule challenge.VerifyRule) Result {
	info, err := os.Stat(rule.Path)
	if err != nil {
		return Result{Passed: false, Message: fmt.Sprintf("无法读取文件: %v", err)}
	}

	expected := strings.TrimSpace(rule.Expect)
	expectedBits, parseErr := strconv.ParseUint(expected, 8, 32)
	if parseErr != nil {
		return Result{Passed: false, Message: fmt.Sprintf("期望权限值无效（应为八进制，如 644 或 4755）: %q", expected)}
	}

	actualBits := modeBits(info.Mode())
	if actualBits != uint32(expectedBits)&0o7777 {
		return Result{
			Passed:  false,
			Message: fmt.Sprintf("权限不匹配\n期望: %s\n实际: %s", expected, strconv.FormatUint(uint64(actualBits), 8)),
		}
	}
	return Result{Passed: true, Message: fmt.Sprintf("权限匹配: %s", expected)}
}

// modeBits maps a fs.FileMode to the traditional low 12 permission bits,
// including setuid (04000), setgid (02000) and sticky (01000), matching the
// output of `stat -c '%a'`.
func modeBits(mode fs.FileMode) uint32 {
	bits := uint32(mode.Perm())
	if mode&fs.ModeSetuid != 0 {
		bits |= 0o4000
	}
	if mode&fs.ModeSetgid != 0 {
		bits |= 0o2000
	}
	if mode&fs.ModeSticky != 0 {
		bits |= 0o1000
	}
	return bits
}
