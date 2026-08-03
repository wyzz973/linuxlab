package verify

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sd3/linuxlab/internal/challenge"
)

func TestRunAll_AllPass(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "test.txt")
	os.WriteFile(p, []byte("hello"), 0644)

	rules := []challenge.VerifyRule{
		{Type: "file_exists", Path: p},
		{Type: "file_content", Path: p, Expect: "hello"},
	}

	results := RunAll(rules)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for i, r := range results {
		if !r.Passed {
			t.Fatalf("rule %d failed: %s", i, r.Message)
		}
	}
}

func TestRunAll_OneFails(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "test.txt")
	os.WriteFile(p, []byte("wrong"), 0644)

	rules := []challenge.VerifyRule{
		{Type: "file_exists", Path: p},
		{Type: "file_content", Path: p, Expect: "hello"},
	}

	results := RunAll(rules)
	if !results[0].Passed {
		t.Fatal("file_exists should pass")
	}
	if results[1].Passed {
		t.Fatal("file_content should fail")
	}
}

func TestRunAll_AllPassed_Helper(t *testing.T) {
	passing := []Result{{Passed: true}, {Passed: true}}
	if !AllPassed(passing) {
		t.Fatal("AllPassed should return true for all passing results")
	}

	mixed := []Result{{Passed: true}, {Passed: false}}
	if AllPassed(mixed) {
		t.Fatal("AllPassed should return false when one fails")
	}
}

func TestRunAll_EmptyRulesFail(t *testing.T) {
	results := RunAll(nil)
	if len(results) != 1 {
		t.Fatalf("expected 1 result for empty rules, got %d", len(results))
	}
	if results[0].Passed {
		t.Fatal("empty rule set should not pass")
	}
	if !strings.Contains(results[0].Message, "缺少验证规则") {
		t.Fatalf("unexpected message: %s", results[0].Message)
	}
	if AllPassed(results) {
		t.Fatal("AllPassed should be false for empty rule set")
	}
}

func TestRunAllWithOptions_ScriptWorkDir(t *testing.T) {
	chDir := t.TempDir()
	workDir := t.TempDir()

	script := filepath.Join(chDir, "check.sh")
	os.WriteFile(script, []byte("#!/bin/bash\ngrep -q hello challenge.txt\n"), 0o755)
	os.WriteFile(filepath.Join(workDir, "challenge.txt"), []byte("hello\n"), 0o644)

	rules := []challenge.VerifyRule{{Type: "script", Path: script}}

	results := RunAllWithOptions(rules, Options{WorkDir: workDir})
	if !results[0].Passed {
		t.Fatalf("expected pass with WorkDir set: %s", results[0].Message)
	}

	// Without WorkDir the relative path in the script resolves against the
	// process cwd, where challenge.txt does not exist.
	results = RunAll(rules)
	if results[0].Passed {
		t.Fatal("expected fail without WorkDir")
	}
}

func TestRunAll_UnknownType(t *testing.T) {
	rules := []challenge.VerifyRule{
		{Type: "unknown_type"},
	}

	results := RunAll(rules)
	if results[0].Passed {
		t.Fatal("unknown type should fail")
	}
	if results[0].Message == "" {
		t.Fatal("expected error message for unknown type")
	}
}
