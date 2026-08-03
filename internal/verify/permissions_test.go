package verify

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sd3/linuxlab/internal/challenge"
)

func TestPermissionsVerifier_Pass(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.sh")
	os.WriteFile(path, []byte("#!/bin/bash"), 0o755)

	rule := challenge.VerifyRule{
		Type:   "permissions",
		Path:   path,
		Expect: "755",
	}

	v := &PermissionsVerifier{}
	result := v.Verify(rule)
	if !result.Passed {
		t.Errorf("expected pass, got fail: %s", result.Message)
	}
}

func TestPermissionsVerifier_Fail(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.sh")
	os.WriteFile(path, []byte("#!/bin/bash"), 0o644)

	rule := challenge.VerifyRule{
		Type:   "permissions",
		Path:   path,
		Expect: "755",
	}

	v := &PermissionsVerifier{}
	result := v.Verify(rule)
	if result.Passed {
		t.Error("expected fail, got pass")
	}
}

func TestPermissionsVerifier_NoFile(t *testing.T) {
	rule := challenge.VerifyRule{
		Type:   "permissions",
		Path:   "/nonexistent",
		Expect: "644",
	}

	v := &PermissionsVerifier{}
	result := v.Verify(rule)
	if result.Passed {
		t.Error("expected fail for nonexistent file")
	}
}

func TestPermissionsVerifier_LeadingZeroExpect(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("data"), 0o644)
	os.Chmod(path, 0o644)

	for _, expect := range []string{"644", "0644"} {
		rule := challenge.VerifyRule{
			Type:   "permissions",
			Path:   path,
			Expect: expect,
		}
		v := &PermissionsVerifier{}
		result := v.Verify(rule)
		if !result.Passed {
			t.Errorf("expect %q: expected pass, got fail: %s", expect, result.Message)
		}
	}
}

func TestPermissionsVerifier_SetuidBit(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "suid.sh")
	os.WriteFile(path, []byte("#!/bin/bash"), 0o755)
	if err := os.Chmod(path, 0o755|os.ModeSetuid); err != nil {
		t.Skipf("cannot set setuid bit: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode()&os.ModeSetuid == 0 {
		t.Skip("filesystem does not preserve setuid bit")
	}

	v := &PermissionsVerifier{}
	result := v.Verify(challenge.VerifyRule{Type: "permissions", Path: path, Expect: "4755"})
	if !result.Passed {
		t.Errorf("expected pass for 4755, got fail: %s", result.Message)
	}

	result = v.Verify(challenge.VerifyRule{Type: "permissions", Path: path, Expect: "755"})
	if result.Passed {
		t.Error("expected fail: setuid bit present but expect is 755")
	}
}

func TestPermissionsVerifier_StickyBitDir(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "shared")
	os.MkdirAll(subdir, 0o777)
	if err := os.Chmod(subdir, 0o777|os.ModeSticky); err != nil {
		t.Skipf("cannot set sticky bit: %v", err)
	}
	info, err := os.Stat(subdir)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode()&os.ModeSticky == 0 {
		t.Skip("filesystem does not preserve sticky bit")
	}

	v := &PermissionsVerifier{}
	result := v.Verify(challenge.VerifyRule{Type: "permissions", Path: subdir, Expect: "1777"})
	if !result.Passed {
		t.Errorf("expected pass for 1777, got fail: %s", result.Message)
	}
}

func TestPermissionsVerifier_InvalidExpect(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("data"), 0o644)

	v := &PermissionsVerifier{}
	result := v.Verify(challenge.VerifyRule{Type: "permissions", Path: path, Expect: "rwxr-xr-x"})
	if result.Passed {
		t.Error("expected fail for non-octal expect value")
	}
}

func TestPermissionsVerifier_DirectoryPerms(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "mydir")
	os.MkdirAll(subdir, 0o750)

	rule := challenge.VerifyRule{
		Type:   "permissions",
		Path:   subdir,
		Expect: "750",
	}

	v := &PermissionsVerifier{}
	result := v.Verify(rule)
	if !result.Passed {
		t.Errorf("expected pass for dir, got fail: %s", result.Message)
	}
}
