package runner

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sd3/linuxlab/internal/challenge"
)

func TestEventShapeHasStableJSONFields(t *testing.T) {
	passed := true
	event := Event{Type: "result", Passed: &passed, HintsUsed: 1}
	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("marshal event: %v", err)
	}
	if !json.Valid(data) {
		t.Fatalf("event is not valid json: %s", data)
	}
	if event.Type != "result" {
		t.Fatalf("unexpected event type: %s", event.Type)
	}
}

// installFakeVim puts a fake "vim" executable on PATH that runs the given
// shell body with the target file as $1.
func installFakeVim(t *testing.T, body string) {
	t.Helper()
	binDir := t.TempDir()
	script := "#!/bin/bash\n" + body + "\n"
	if err := os.WriteFile(filepath.Join(binDir, "vim"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func vimChallenge(t *testing.T, checkScript string) *challenge.Challenge {
	t.Helper()
	chDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(chDir, "check.sh"), []byte(checkScript), 0o755); err != nil {
		t.Fatal(err)
	}
	return &challenge.Challenge{
		ID:       "vim-test",
		Title:    "测试题",
		Category: "vim",
		Dir:      chDir,
		SetupFiles: []challenge.SetupFile{
			{Path: "challenge.txt", Content: "初始内容\n"},
		},
		Verify: []challenge.VerifyRule{
			{Type: "script", Path: "check.sh"},
		},
	}
}

func runOpts(ch *challenge.Challenge) Options {
	return Options{
		Challenge: ch,
		Stdin:     strings.NewReader(""),
		Stdout:    io.Discard,
		Stderr:    io.Discard,
	}
}

// runHostShell must execute scripts on the host in the challenge directory —
// compose init/check scripts rely on the host docker CLI and host paths.
func TestRunHostShellUsesChallengeDir(t *testing.T) {
	out, code, err := runHostShell(context.Background(), "/", "pwd")
	if err != nil {
		t.Fatalf("runHostShell: %v", err)
	}
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if strings.TrimSpace(out) != "/" {
		t.Errorf("pwd = %q, want %q", strings.TrimSpace(out), "/")
	}

	if _, code, _ := runHostShell(context.Background(), "/", "exit 7"); code != 7 {
		t.Errorf("exit code = %d, want 7 (ExitError must carry the status)", code)
	}
}

// Compose challenges are verified on the host (docker compose resources live
// there), so verifyInSandbox must drive its execFn with the challenge dir.
func TestVerifyInSandboxScriptRuleUsesExecFn(t *testing.T) {
	chDir := t.TempDir()
	checkPath := filepath.Join(chDir, "check.sh")
	if err := os.WriteFile(checkPath, []byte("#!/bin/sh\ntest -f host-marker.txt\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(chDir, "host-marker.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	// A fake execFn that runs the script on the host inside chDir — the same
	// wiring runSandbox uses for ComposeSandbox.
	execFn := func(ctx context.Context, command string) (string, int, error) {
		return runHostShell(ctx, chDir, command)
	}

	ch := &challenge.Challenge{Dir: chDir, Verify: []challenge.VerifyRule{{Type: "script", Path: "check.sh"}}}
	results := verifyInSandbox(context.Background(), ch, execFn)
	if !results[0].Passed {
		t.Fatalf("expected host-context check to pass, got: %+v", results)
	}
}

// A script rule that fails must surface its output, not vanish.
func TestVerifyInSandboxScriptRuleFailureMessage(t *testing.T) {
	chDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(chDir, "check.sh"), []byte("#!/bin/sh\necho 未完成\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	ch := &challenge.Challenge{Dir: chDir, Verify: []challenge.VerifyRule{{Type: "script", Path: "check.sh"}}}
	results := verifyInSandbox(context.Background(), ch, func(ctx context.Context, command string) (string, int, error) {
		return runHostShell(ctx, chDir, command)
	})
	if results[0].Passed {
		t.Fatal("expected failure")
	}
	if !strings.Contains(results[0].Message, "未完成") {
		t.Fatalf("failure message must include script output, got: %q", results[0].Message)
	}
}

// Regression for id=22: a script rule with a relative path (check.sh) must
// resolve against the challenge dir, and the script must run inside the vim
// work dir so its relative paths (cat challenge.txt) see the edited file.
func TestRunVim_RelativeScriptPathVerifiesEditedFile(t *testing.T) {
	tmpRoot := t.TempDir()
	t.Setenv("TMPDIR", tmpRoot)
	installFakeVim(t, `echo "已修改" >> "$1"`)

	ch := vimChallenge(t, "#!/bin/bash\ngrep -q '已修改' challenge.txt\n")

	result, err := RunInteractive(context.Background(), runOpts(ch))
	if err != nil {
		t.Fatalf("RunInteractive: %v", err)
	}
	if !result.Passed {
		t.Fatalf("expected pass, got: %+v", result.Results)
	}

	// Regression for id=29: the per-run work dir must be cleaned up.
	entries, err := os.ReadDir(tmpRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("vim work dir not cleaned up: %v", entries)
	}
}

// Regression for id=0: a non-zero editor exit code (e.g. :cq) must not abort
// the run — verification still executes.
func TestRunVim_IgnoresEditorExitError(t *testing.T) {
	installFakeVim(t, `echo "已修改" >> "$1"`+"\nexit 3")

	ch := vimChallenge(t, "#!/bin/bash\ngrep -q '已修改' challenge.txt\n")

	result, err := RunInteractive(context.Background(), runOpts(ch))
	if err != nil {
		t.Fatalf("editor exit code should not be fatal, got: %v", err)
	}
	if !result.Passed {
		t.Fatalf("expected pass, got: %+v", result.Results)
	}
}

// A missing editor binary is a startup failure and must surface as an error.
func TestRunVim_MissingEditorIsFatal(t *testing.T) {
	emptyDir := t.TempDir()
	t.Setenv("PATH", emptyDir)

	ch := vimChallenge(t, "#!/bin/bash\nexit 0\n")

	if _, err := RunInteractive(context.Background(), runOpts(ch)); err == nil {
		t.Fatal("expected error when vim cannot be started")
	}
}

// Regression for id=25: an empty verify rule set must fail instead of
// auto-passing.
func TestRunVim_EmptyVerifyRulesFail(t *testing.T) {
	installFakeVim(t, "exit 0")

	ch := vimChallenge(t, "#!/bin/bash\nexit 0\n")
	ch.Verify = nil

	result, err := RunInteractive(context.Background(), runOpts(ch))
	if err != nil {
		t.Fatalf("RunInteractive: %v", err)
	}
	if result.Passed {
		t.Fatal("challenge without verify rules must not pass")
	}
	if len(result.Results) == 0 || !strings.Contains(result.Results[0].Message, "缺少验证规则") {
		t.Fatalf("unexpected results: %+v", result.Results)
	}
}

// Regression for id=11: a file rule that matches no setup file must fail
// explicitly instead of silently keeping the relative path.
func TestRunVim_UnmatchedFileRuleFails(t *testing.T) {
	installFakeVim(t, "exit 0")

	ch := vimChallenge(t, "#!/bin/bash\nexit 0\n")
	ch.Verify = []challenge.VerifyRule{
		{Type: "file_content", Path: "nonexistent.txt", Expect: "x"},
	}

	result, err := RunInteractive(context.Background(), runOpts(ch))
	if err != nil {
		t.Fatalf("RunInteractive: %v", err)
	}
	if result.Passed {
		t.Fatal("unmatched file rule must not pass")
	}
	if len(result.Results) == 0 || !strings.Contains(result.Results[0].Message, "nonexistent.txt") {
		t.Fatalf("unexpected results: %+v", result.Results)
	}
}

// Regression for id=11: rule paths map to setup files by exact base-name
// match, not by suffix, so "notes.txt" must not bind to "draft_notes.txt".
func TestResolveVimRules(t *testing.T) {
	chDir := t.TempDir()
	workDir := t.TempDir()
	ch := &challenge.Challenge{
		Dir: chDir,
		Verify: []challenge.VerifyRule{
			{Type: "script", Path: "check.sh"},
			{Type: "file_content", Path: "notes.txt", Expect: "x"},
		},
	}
	paths := []string{
		filepath.Join(workDir, "draft_notes.txt"),
		filepath.Join(workDir, "notes.txt"),
	}

	rules, unmatched := resolveVimRules(ch, workDir, paths)
	if unmatched != "" {
		t.Fatalf("unexpected unmatched rule path: %s", unmatched)
	}
	if rules[0].Path != filepath.Join(chDir, "check.sh") {
		t.Errorf("script rule path = %q, want %q", rules[0].Path, filepath.Join(chDir, "check.sh"))
	}
	if rules[1].Path != filepath.Join(workDir, "notes.txt") {
		t.Errorf("file rule path = %q, want %q (suffix collision with draft_notes.txt)", rules[1].Path, filepath.Join(workDir, "notes.txt"))
	}

	// Original rules must stay untouched.
	if ch.Verify[0].Path != "check.sh" || ch.Verify[1].Path != "notes.txt" {
		t.Errorf("resolveVimRules mutated the challenge rules: %+v", ch.Verify)
	}

	// A file rule that matches nothing reports the offending path.
	ch.Verify = []challenge.VerifyRule{{Type: "file_exists", Path: "missing.txt"}}
	if _, unmatched := resolveVimRules(ch, workDir, paths); unmatched != "missing.txt" {
		t.Errorf("unmatched = %q, want missing.txt", unmatched)
	}
}
