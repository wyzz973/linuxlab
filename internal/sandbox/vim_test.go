package sandbox

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sd3/linuxlab/internal/challenge"
)

func TestVimRunner_PrepareFiles(t *testing.T) {
	dir := t.TempDir()
	runner := &VimRunner{WorkDir: dir}

	setupFiles := []challenge.SetupFile{
		{Path: "test.txt", Content: "hello world"},
		{Path: "subdir/nested.txt", Content: "nested content"},
	}

	paths, err := runner.PrepareFiles(setupFiles)
	if err != nil {
		t.Fatalf("PrepareFiles failed: %v", err)
	}
	if len(paths) != 2 {
		t.Fatalf("got %d paths, want 2", len(paths))
	}

	content, _ := os.ReadFile(paths[0])
	if string(content) != "hello world" {
		t.Errorf("content = %q", string(content))
	}

	content2, _ := os.ReadFile(paths[1])
	if string(content2) != "nested content" {
		t.Errorf("content = %q", string(content2))
	}
}

func TestVimRunner_VimArgs(t *testing.T) {
	dir := t.TempDir()
	runner := &VimRunner{WorkDir: dir}
	filePath := filepath.Join(dir, "test.txt")

	args := runner.VimArgs(filePath)
	if args[0] != "vim" {
		t.Errorf("args[0] = %q, want vim", args[0])
	}
	if args[len(args)-1] != filePath {
		t.Errorf("last arg = %q, want %q", args[len(args)-1], filePath)
	}
}

func TestVimRunner_CheckResult(t *testing.T) {
	dir := t.TempDir()
	runner := &VimRunner{WorkDir: dir}

	filePath := filepath.Join(dir, "test.txt")
	os.WriteFile(filePath, []byte("bar is great\n"), 0o644)

	rules := []challenge.VerifyRule{
		{Type: "file_content", Path: filePath, Expect: "bar is great"},
	}

	passed := runner.CheckResult(rules)
	if !passed {
		t.Error("expected pass")
	}
}
