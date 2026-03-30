package tui

import (
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sd3/linuxlab/internal/challenge"
	"github.com/sd3/linuxlab/internal/progress"
)

func tmpStore(t *testing.T) *progress.Store {
	t.Helper()
	dir := t.TempDir()
	s, err := progress.NewStore(filepath.Join(dir, "progress.json"))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func sampleCategories() map[string][]*challenge.Challenge {
	return map[string][]*challenge.Challenge{
		"linux-basics": {
			{ID: "ls-basic", Title: "ls 基础", Category: "linux-basics", Difficulty: 1},
			{ID: "cp-files", Title: "cp 文件", Category: "linux-basics", Difficulty: 2},
		},
		"vim": {
			{ID: "vim-move", Title: "Vim 移动", Category: "vim", Difficulty: 1},
		},
	}
}

func TestModulesModel_Init(t *testing.T) {
	cats := sampleCategories()
	store := tmpStore(t)
	m := NewModulesModel(cats, store).(ModulesModel)
	if len(m.modules) != len(cats) {
		t.Errorf("expected %d modules, got %d", len(cats), len(m.modules))
	}
}

func TestModulesModel_Enter_SelectsModule(t *testing.T) {
	cats := sampleCategories()
	store := tmpStore(t)
	m := NewModulesModel(cats, store).(ModulesModel)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected a command")
	}
	msg := cmd()
	sel, ok := msg.(ModuleSelectedMsg)
	if !ok {
		t.Fatalf("expected ModuleSelectedMsg, got %T", msg)
	}
	// The selected category should be one of the modules
	found := false
	for _, mod := range m.modules {
		if mod.category == sel.Category {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("selected category %q not in modules", sel.Category)
	}
}

func TestModulesModel_Escape_GoesBack(t *testing.T) {
	cats := sampleCategories()
	store := tmpStore(t)
	m := NewModulesModel(cats, store).(ModulesModel)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected a command")
	}
	msg := cmd()
	if _, ok := msg.(GoBackMsg); !ok {
		t.Fatalf("expected GoBackMsg, got %T", msg)
	}
}

// Ensure tmpStore doesn't leave files around (uses t.TempDir).
func TestModulesModel_ProgressDisplay(t *testing.T) {
	cats := sampleCategories()
	dir := t.TempDir()
	s, _ := progress.NewStore(filepath.Join(dir, "progress.json"))
	s.RecordAttempt("ls-basic", "linux-basics", "navigation", true, 0)
	_ = s.Save()
	// reload
	s2, _ := progress.NewStore(filepath.Join(dir, "progress.json"))
	m := NewModulesModel(cats, s2).(ModulesModel)
	view := m.View()
	if view == "" {
		t.Error("expected non-empty view")
	}
	// Cleanup
	os.RemoveAll(dir)
}
