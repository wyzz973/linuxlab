package verify

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sd3/linuxlab/internal/challenge"
)

func TestScriptVerifier_Pass(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "check.sh")
	os.WriteFile(script, []byte("#!/bin/bash\nexit 0\n"), 0o755)

	rule := challenge.VerifyRule{
		Type: "script",
		Path: script,
	}

	v := &ScriptVerifier{}
	result := v.Verify(rule)
	if !result.Passed {
		t.Errorf("expected pass, got fail: %s", result.Message)
	}
}

func TestScriptVerifier_Fail(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "check.sh")
	os.WriteFile(script, []byte("#!/bin/bash\necho 'wrong answer'\nexit 1\n"), 0o755)

	rule := challenge.VerifyRule{
		Type: "script",
		Path: script,
	}

	v := &ScriptVerifier{}
	result := v.Verify(rule)
	if result.Passed {
		t.Error("expected fail, got pass")
	}
	if result.Message == "" {
		t.Error("expected non-empty message with script output")
	}
}

func TestScriptVerifier_NotFound(t *testing.T) {
	rule := challenge.VerifyRule{
		Type: "script",
		Path: "/nonexistent/check.sh",
	}

	v := &ScriptVerifier{}
	result := v.Verify(rule)
	if result.Passed {
		t.Error("expected fail for missing script")
	}
}

func TestScriptVerifier_NotExecutable(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "check.sh")
	os.WriteFile(script, []byte("#!/bin/bash\nexit 0\n"), 0o644) // no exec perm

	rule := challenge.VerifyRule{
		Type: "script",
		Path: script,
	}

	v := &ScriptVerifier{}
	result := v.Verify(rule)
	// Should still work — we run via "bash <script>"
	if !result.Passed {
		t.Errorf("expected pass even without exec perm: %s", result.Message)
	}
}

func TestScriptVerifier_CapturesOutput(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "check.sh")
	os.WriteFile(script, []byte("#!/bin/bash\necho 'test output'\nexit 1\n"), 0o755)

	rule := challenge.VerifyRule{
		Type: "script",
		Path: script,
	}

	v := &ScriptVerifier{}
	result := v.Verify(rule)
	if result.Passed {
		t.Error("expected fail")
	}
}
