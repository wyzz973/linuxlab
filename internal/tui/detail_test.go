package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sd3/linuxlab/internal/challenge"
)

func sampleChallenge() *challenge.Challenge {
	return &challenge.Challenge{
		ID:          "ls-basic",
		Title:       "ls 基础",
		Category:    "linux-basics",
		Subcategory: "navigation",
		Difficulty:  2,
		Tags:        []string{"ls", "list"},
		Description: "学习使用 ls 命令列出目录内容。",
		Hints: []challenge.Hint{
			{Level: 1, Text: "试试 ls -la"},
			{Level: 2, Text: "使用 ls -la /path"},
			{Level: 3, Text: "答案是 ls -la /home"},
		},
	}
}

func TestDetailModel_Init(t *testing.T) {
	ch := sampleChallenge()
	m := NewDetailModel(ch).(DetailModel)
	if m.hintLevel != 0 {
		t.Errorf("expected hintLevel=0, got %d", m.hintLevel)
	}
}

func TestDetailModel_Enter_LaunchesChallenge(t *testing.T) {
	ch := sampleChallenge()
	m := NewDetailModel(ch).(DetailModel)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected a command")
	}
	msg := cmd()
	launch, ok := msg.(LaunchChallengeMsg)
	if !ok {
		t.Fatalf("expected LaunchChallengeMsg, got %T", msg)
	}
	if launch.Challenge.ID != "ls-basic" {
		t.Errorf("expected ls-basic, got %q", launch.Challenge.ID)
	}
	if launch.HintsUsed != 0 {
		t.Errorf("expected HintsUsed=0, got %d", launch.HintsUsed)
	}
}

func TestDetailModel_Hint_Progressive(t *testing.T) {
	ch := sampleChallenge()
	m := NewDetailModel(ch).(DetailModel)

	// First hint
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	m = updated.(DetailModel)
	if m.hintLevel != 1 {
		t.Errorf("expected hintLevel=1, got %d", m.hintLevel)
	}

	// Second hint
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	m = updated.(DetailModel)
	if m.hintLevel != 2 {
		t.Errorf("expected hintLevel=2, got %d", m.hintLevel)
	}

	// Third hint
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	m = updated.(DetailModel)
	if m.hintLevel != 3 {
		t.Errorf("expected hintLevel=3, got %d", m.hintLevel)
	}

	// Should not exceed max
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	m = updated.(DetailModel)
	if m.hintLevel != 3 {
		t.Errorf("expected hintLevel=3 (capped), got %d", m.hintLevel)
	}
}

func TestDetailModel_HintsUsed(t *testing.T) {
	ch := sampleChallenge()
	m := NewDetailModel(ch).(DetailModel)

	// Reveal 2 hints
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	m = updated.(DetailModel)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	m = updated.(DetailModel)

	// Launch should report 2 hints used
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	msg := cmd()
	launch := msg.(LaunchChallengeMsg)
	if launch.HintsUsed != 2 {
		t.Errorf("expected HintsUsed=2, got %d", launch.HintsUsed)
	}
}
