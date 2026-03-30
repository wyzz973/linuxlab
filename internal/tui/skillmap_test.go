package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestSkillMapModel_View_ShowsCategories(t *testing.T) {
	store := tmpStore(t)
	store.RecordAttempt("ls-basic", "linux-basics", "navigation", true, 0)
	store.RecordAttempt("vim-move", "vim", "movement", false, 1)

	m := NewSkillMapModel(store).(SkillMapModel)
	view := m.View()

	if !strings.Contains(view, "Linux 基础命令") && !strings.Contains(view, "linux-basics") {
		t.Error("expected view to contain linux-basics category")
	}
	if !strings.Contains(view, "Vim 操作") && !strings.Contains(view, "vim") {
		t.Error("expected view to contain vim category")
	}
}

func TestSkillMapModel_Toggle_Expand(t *testing.T) {
	store := tmpStore(t)
	store.RecordAttempt("ls-basic", "linux-basics", "navigation", true, 0)
	store.RecordAttempt("vim-move", "vim", "movement", false, 1)

	m := NewSkillMapModel(store).(SkillMapModel)

	// Initially not expanded
	if m.expanded[0] {
		t.Error("expected category 0 to not be expanded initially")
	}

	// Toggle expand
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(SkillMapModel)
	if !m.expanded[0] {
		t.Error("expected category 0 to be expanded after Enter")
	}

	// Toggle collapse
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(SkillMapModel)
	if m.expanded[0] {
		t.Error("expected category 0 to be collapsed after second Enter")
	}
}

func TestSkillMapModel_Escape_GoesBack(t *testing.T) {
	store := tmpStore(t)
	m := NewSkillMapModel(store).(SkillMapModel)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected a command")
	}
	msg := cmd()
	if _, ok := msg.(GoBackMsg); !ok {
		t.Fatalf("expected GoBackMsg, got %T", msg)
	}
}
