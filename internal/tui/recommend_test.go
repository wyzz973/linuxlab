package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sd3/linuxlab/internal/challenge"
)

func TestRecommendModel_Init(t *testing.T) {
	recs := []*challenge.Challenge{
		{ID: "ch1", Title: "题目1", Difficulty: 1, Category: "linux-basics", Subcategory: "permissions"},
		{ID: "ch2", Title: "题目2", Difficulty: 2, Category: "vim", Subcategory: "editing"},
	}

	m := NewRecommendModel(recs).(RecommendModel)
	if len(m.challenges) != 2 {
		t.Errorf("challenges = %d, want 2", len(m.challenges))
	}
}

func TestRecommendModel_Enter_SelectsChallenge(t *testing.T) {
	recs := []*challenge.Challenge{
		{ID: "ch1", Title: "题目1", Difficulty: 1},
	}

	m := NewRecommendModel(recs)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected command")
	}
	msg := cmd()
	sel, ok := msg.(ChallengeSelectedMsg)
	if !ok {
		t.Fatalf("expected ChallengeSelectedMsg, got %T", msg)
	}
	if sel.Challenge.ID != "ch1" {
		t.Errorf("selected = %q", sel.Challenge.ID)
	}
}

func TestRecommendModel_Empty(t *testing.T) {
	m := NewRecommendModel(nil)
	view := m.View()
	if view == "" {
		t.Error("view is empty")
	}
}

func TestRecommendModel_Escape(t *testing.T) {
	m := NewRecommendModel(nil)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected command on Esc")
	}
	msg := cmd()
	if _, ok := msg.(GoBackMsg); !ok {
		t.Fatalf("expected GoBackMsg, got %T", msg)
	}
}
