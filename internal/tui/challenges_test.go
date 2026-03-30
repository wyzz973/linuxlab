package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sd3/linuxlab/internal/challenge"
)

func sampleChallenges() []*challenge.Challenge {
	return []*challenge.Challenge{
		{ID: "ls-basic", Title: "ls 基础", Category: "linux-basics", Difficulty: 1},
		{ID: "cp-files", Title: "cp 文件", Category: "linux-basics", Difficulty: 2},
		{ID: "mv-rename", Title: "mv 重命名", Category: "linux-basics", Difficulty: 2},
	}
}

func TestChallengesModel_Init(t *testing.T) {
	store := tmpStore(t)
	chs := sampleChallenges()
	m := NewChallengesModel("linux-basics", chs, store).(ChallengesModel)
	if len(m.challenges) != 3 {
		t.Errorf("expected 3 challenges, got %d", len(m.challenges))
	}
	if m.category != "linux-basics" {
		t.Errorf("expected category linux-basics, got %q", m.category)
	}
}

func TestChallengesModel_Enter_SelectsChallenge(t *testing.T) {
	store := tmpStore(t)
	chs := sampleChallenges()
	m := NewChallengesModel("linux-basics", chs, store).(ChallengesModel)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected a command")
	}
	msg := cmd()
	sel, ok := msg.(ChallengeSelectedMsg)
	if !ok {
		t.Fatalf("expected ChallengeSelectedMsg, got %T", msg)
	}
	if sel.Challenge.ID != "ls-basic" {
		t.Errorf("expected ls-basic, got %q", sel.Challenge.ID)
	}
}

func TestChallengesModel_View_ShowsStatus(t *testing.T) {
	store := tmpStore(t)
	chs := sampleChallenges()
	m := NewChallengesModel("linux-basics", chs, store).(ChallengesModel)
	view := m.View()
	if view == "" {
		t.Error("expected non-empty view")
	}
}
