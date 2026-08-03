package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestMenuModel_ViewShowsConsoleHeaderAndNumberedEntries(t *testing.T) {
	m := NewMenuModelWithStats(42, 5).(MenuModel)
	m.width = 100
	m.height = 28

	view := m.View()
	for _, want := range []string{"LinuxLab 训练控制台", "真实 Shell / Vim / 运维场景训练", "1", "42 道挑战"} {
		if !strings.Contains(view, want) {
			t.Fatalf("menu view missing %q\n%s", want, view)
		}
	}
}

func TestMenuModel_NumberKeySelectsEntry(t *testing.T) {
	m := NewMenuModel()

	_, cmd := m.Update(runeKey('4'))
	if cmd == nil {
		t.Fatal("expected number key to select a menu entry")
	}
	choice := cmd().(MenuChoiceMsg)
	if choice.Choice != "reference" {
		t.Fatalf("choice = %q, want reference", choice.Choice)
	}
}

func TestModulesModel_ViewShowsProgressSummary(t *testing.T) {
	cats := sampleCategories()
	store := tmpStore(t)
	store.RecordAttempt("ls-basic", "linux-basics", "navigation", true, 0)

	m := NewModulesModel(cats, store).(ModulesModel)
	m.width = 100
	m.height = 28

	view := m.View()
	for _, want := range []string{"整体进度", "继续训练", "1"} {
		if !strings.Contains(view, want) {
			t.Fatalf("modules view missing %q\n%s", want, view)
		}
	}
}

func TestModulesModel_NumberKeySelectsModule(t *testing.T) {
	cats := sampleCategories()
	store := tmpStore(t)
	m := NewModulesModel(cats, store).(ModulesModel)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	if cmd == nil {
		t.Fatal("expected number key to select a module")
	}
	if _, ok := cmd().(ModuleSelectedMsg); !ok {
		t.Fatalf("expected ModuleSelectedMsg")
	}
}

func TestChallengesModel_ViewShowsProgressSummary(t *testing.T) {
	store := tmpStore(t)
	store.RecordAttempt("ls-basic", "linux-basics", "navigation", true, 0)
	m := NewChallengesModel("linux-basics", sampleChallenges(), store).(ChallengesModel)
	m.width = 100
	m.height = 28

	view := m.View()
	for _, want := range []string{"完成进度", "已通过", "未完成"} {
		if !strings.Contains(view, want) {
			t.Fatalf("challenges view missing %q\n%s", want, view)
		}
	}
}
