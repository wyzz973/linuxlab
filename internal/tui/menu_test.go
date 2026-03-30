package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func newMenu() MenuModel {
	return NewMenuModel().(MenuModel)
}

func TestMenuModel_Init(t *testing.T) {
	m := newMenu()
	if m.cursor != 0 {
		t.Errorf("expected cursor=0, got %d", m.cursor)
	}
	if len(m.items) != 4 {
		t.Errorf("expected 4 items, got %d", len(m.items))
	}
}

func TestMenuModel_MoveDown(t *testing.T) {
	m := newMenu()
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	mm := updated.(MenuModel)
	if mm.cursor != 1 {
		t.Errorf("expected cursor=1, got %d", mm.cursor)
	}
}

func TestMenuModel_MoveUp_Wraps(t *testing.T) {
	m := newMenu()
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	mm := updated.(MenuModel)
	if mm.cursor != 3 {
		t.Errorf("expected cursor=3, got %d", mm.cursor)
	}
}

func TestMenuModel_MoveDown_Wraps(t *testing.T) {
	m := newMenu()
	m.cursor = 3
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	mm := updated.(MenuModel)
	if mm.cursor != 0 {
		t.Errorf("expected cursor=0, got %d", mm.cursor)
	}
}

func TestMenuModel_JK_Navigation(t *testing.T) {
	m := newMenu()
	// j = down
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	mm := updated.(MenuModel)
	if mm.cursor != 1 {
		t.Errorf("j: expected cursor=1, got %d", mm.cursor)
	}
	// k = up
	updated, _ = mm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	mm = updated.(MenuModel)
	if mm.cursor != 0 {
		t.Errorf("k: expected cursor=0, got %d", mm.cursor)
	}
}

func TestMenuModel_Enter_ReturnsChoice(t *testing.T) {
	m := newMenu()
	m.cursor = 1
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected a command, got nil")
	}
	msg := cmd()
	choice, ok := msg.(MenuChoiceMsg)
	if !ok {
		t.Fatalf("expected MenuChoiceMsg, got %T", msg)
	}
	if choice.Choice != "skillmap" {
		t.Errorf("expected 'skillmap', got %q", choice.Choice)
	}
}
