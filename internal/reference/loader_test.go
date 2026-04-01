package reference

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadReferences_Valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "commands.yaml")
	data := `commands:
  - name: ls
    brief: "列出目录内容"
    examples:
      - desc: "列出所有文件"
        cmd: "ls -la {{目录}}"
      - desc: "按大小排序"
        cmd: "ls -lS"
    related_challenges:
      - ls-basic
  - name: grep
    brief: "搜索文本"
    examples:
      - desc: "搜索模式"
        cmd: "grep {{模式}} {{文件}}"
`
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	refs, err := LoadReferences(path)
	if err != nil {
		t.Fatalf("LoadReferences failed: %v", err)
	}

	if len(refs.Commands) != 2 {
		t.Fatalf("expected 2 commands, got %d", len(refs.Commands))
	}

	ls := refs.Commands[0]
	if ls.Name != "ls" {
		t.Errorf("expected name 'ls', got %q", ls.Name)
	}
	if ls.Brief != "列出目录内容" {
		t.Errorf("expected brief '列出目录内容', got %q", ls.Brief)
	}
	if len(ls.Examples) != 2 {
		t.Errorf("expected 2 examples, got %d", len(ls.Examples))
	}
	if ls.Examples[0].Desc != "列出所有文件" {
		t.Errorf("expected example desc '列出所有文件', got %q", ls.Examples[0].Desc)
	}
	if ls.Examples[0].Cmd != "ls -la {{目录}}" {
		t.Errorf("expected example cmd 'ls -la {{目录}}', got %q", ls.Examples[0].Cmd)
	}
	if len(ls.RelatedChallenges) != 1 || ls.RelatedChallenges[0] != "ls-basic" {
		t.Errorf("expected related_challenges [ls-basic], got %v", ls.RelatedChallenges)
	}
}

func TestLoadReferences_NonexistentFile(t *testing.T) {
	_, err := LoadReferences("/nonexistent/path/commands.yaml")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestLoadReferences_EmptyCommands(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "commands.yaml")
	data := `commands: []
`
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	refs, err := LoadReferences(path)
	if err != nil {
		t.Fatalf("LoadReferences failed: %v", err)
	}

	if len(refs.Commands) != 0 {
		t.Errorf("expected 0 commands, got %d", len(refs.Commands))
	}
}
