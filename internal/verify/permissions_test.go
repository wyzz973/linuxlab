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
