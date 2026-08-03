package tui

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sd3/linuxlab/internal/challenge"
	"github.com/sd3/linuxlab/internal/verify"
)

func runeKey(r rune) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
}

func manyChallenges(count int) []*challenge.Challenge {
	chs := make([]*challenge.Challenge, count)
	for i := 0; i < count; i++ {
		chs[i] = &challenge.Challenge{
			ID:         fmt.Sprintf("challenge-%02d", i+1),
			Title:      fmt.Sprintf("挑战 %02d", i+1),
			Category:   "linux-basics",
			Difficulty: (i % 5) + 1,
		}
	}
	return chs
}

func TestChallengesModel_PageAndEdgeNavigation(t *testing.T) {
	store := tmpStore(t)
	m := NewChallengesModel("linux-basics", manyChallenges(12), store).(ChallengesModel)
	m.height = 14

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	m = updated.(ChallengesModel)
	if m.cursor != 5 {
		t.Fatalf("PgDown cursor = %d, want 5", m.cursor)
	}
	if m.offset != 1 {
		t.Fatalf("PgDown offset = %d, want 1", m.offset)
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnd})
	m = updated.(ChallengesModel)
	if m.cursor != 11 {
		t.Fatalf("End cursor = %d, want 11", m.cursor)
	}
	if m.offset != 7 {
		t.Fatalf("End offset = %d, want 7", m.offset)
	}

	updated, _ = m.Update(runeKey('g'))
	m = updated.(ChallengesModel)
	if m.cursor != 0 || m.offset != 0 {
		t.Fatalf("g should jump to top, got cursor=%d offset=%d", m.cursor, m.offset)
	}
}

func TestChallengesModel_QGoesBack(t *testing.T) {
	store := tmpStore(t)
	m := NewChallengesModel("linux-basics", sampleChallenges(), store).(ChallengesModel)

	_, cmd := m.Update(runeKey('q'))
	if cmd == nil {
		t.Fatal("expected q to return a GoBackMsg")
	}
	if _, ok := cmd().(GoBackMsg); !ok {
		t.Fatalf("expected GoBackMsg from q")
	}
}

func TestDetailModel_LongContentScrolls(t *testing.T) {
	ch := sampleChallenge()
	ch.Description = strings.Repeat("这一行用于制造较长的任务说明。\n", 24)
	m := NewDetailModel(ch).(DetailModel)
	m.width = 80
	m.height = 14

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	m = updated.(DetailModel)
	if m.scroll == 0 {
		t.Fatal("expected PgDown to scroll long detail content")
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyHome})
	m = updated.(DetailModel)
	if m.scroll != 0 {
		t.Fatalf("Home scroll = %d, want 0", m.scroll)
	}
}

func TestDetailModel_QGoesBack(t *testing.T) {
	m := NewDetailModel(sampleChallenge()).(DetailModel)

	_, cmd := m.Update(runeKey('q'))
	if cmd == nil {
		t.Fatal("expected q to return a GoBackMsg")
	}
	if _, ok := cmd().(GoBackMsg); !ok {
		t.Fatalf("expected GoBackMsg from q")
	}
}

func TestReferenceModel_SearchResultsScrollWithCursor(t *testing.T) {
	refs := sampleRefs()
	for i := 0; i < 12; i++ {
		refs.Commands = append(refs.Commands, refs.Commands[0])
		refs.Commands[len(refs.Commands)-1].Name = fmt.Sprintf("ls%d", i)
	}
	m := NewReferenceModel(refs).(ReferenceModel)
	m.height = 14

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	m = updated.(ReferenceModel)
	if m.cursor != 5 {
		t.Fatalf("PgDown cursor = %d, want 5", m.cursor)
	}
	if m.offset != 1 {
		t.Fatalf("PgDown offset = %d, want 1", m.offset)
	}

	view := m.View()
	if !strings.Contains(view, "ls2") {
		t.Fatalf("expected scrolled reference list to include visible cursor window")
	}
}

func TestReferenceModel_EscapeClearsNonEmptyQueryFirst(t *testing.T) {
	m := NewReferenceModel(sampleRefs())
	m, _ = m.Update(runeKey('l'))
	m, _ = m.Update(runeKey('s'))

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd != nil {
		t.Fatal("Esc with a search query should clear the query, not leave the screen")
	}
	rm := updated.(ReferenceModel)
	if rm.query != "" {
		t.Fatalf("query = %q, want empty", rm.query)
	}
	if len(rm.filtered) != len(rm.commands) {
		t.Fatalf("filtered len = %d, want %d", len(rm.filtered), len(rm.commands))
	}
}

func TestAppModel_ResultRetryLaunchesCurrentChallenge(t *testing.T) {
	store := tmpStore(t)
	cats := sampleCategories()
	m := NewAppModel(cats, store, nil).(AppModel)
	m.screen = screenResult
	m.currentCat = "linux-basics"
	m.currentChallenge = cats["linux-basics"][0]
	m.lastResult = &ChallengeResultMsg{Passed: false, Results: []verify.Result{{Passed: false, Message: "未通过"}}}

	_, cmd := m.Update(runeKey('r'))
	if cmd == nil {
		t.Fatal("expected r to launch the current challenge")
	}
	msg := cmd()
	launch, ok := msg.(LaunchChallengeMsg)
	if !ok {
		t.Fatalf("expected LaunchChallengeMsg, got %T", msg)
	}
	if launch.Challenge.ID != "ls-basic" {
		t.Fatalf("retry challenge = %q, want ls-basic", launch.Challenge.ID)
	}
}

func TestAppModel_ResultNextOpensNextChallenge(t *testing.T) {
	store := tmpStore(t)
	cats := sampleCategories()
	m := NewAppModel(cats, store, nil).(AppModel)
	m.width = 80
	m.height = 24
	m.screen = screenResult
	m.currentCat = "linux-basics"
	m.currentChallenge = cats["linux-basics"][0]
	m.lastResult = &ChallengeResultMsg{Passed: true}

	updated, cmd := m.Update(runeKey('n'))
	if cmd != nil {
		t.Fatal("expected n to navigate without launching immediately")
	}
	app := updated.(AppModel)
	if app.screen != screenDetail {
		t.Fatalf("screen = %d, want screenDetail", app.screen)
	}
	if app.currentChallenge.ID != "cp-files" {
		t.Fatalf("next challenge = %q, want cp-files", app.currentChallenge.ID)
	}
}

func TestAppModel_ChallengeSelectionStoresCategoryContext(t *testing.T) {
	store := tmpStore(t)
	cats := sampleCategories()
	m := NewAppModel(cats, store, nil)
	ch := cats["vim"][0]

	updated, _ := m.Update(ChallengeSelectedMsg{Challenge: ch})
	app := updated.(AppModel)
	if app.currentCat != "vim" {
		t.Fatalf("currentCat = %q, want vim", app.currentCat)
	}
}

func TestBoxWidth_DoesNotExceedSmallTerminal(t *testing.T) {
	if got := boxWidth(40); got > 40 {
		t.Fatalf("boxWidth(40) = %d, want <= 40", got)
	}
}
