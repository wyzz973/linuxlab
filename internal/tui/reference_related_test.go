package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sd3/linuxlab/internal/progress"
)

// openLsDetail opens the detail view of the "ls" command, whose fixture lists
// "ls-basic" as a related challenge.
func openLsDetail(t *testing.T, m tea.Model) string {
	t.Helper()
	m, _ = m.Update(tea.WindowSizeMsg{Width: 90, Height: 26})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // first row is "ls"
	rm, ok := m.(ReferenceModel)
	if !ok {
		t.Fatalf("unexpected model type %T", m)
	}
	if !rm.showDetail {
		t.Fatal("Enter should open the command detail")
	}
	return rm.View()
}

// lsDetailWith renders the "ls" detail view against a fresh model wired to the
// given store, so each status variant is read live from progress.
func lsDetailWith(t *testing.T, store *progress.Store) string {
	t.Helper()
	return openLsDetail(t, NewReferenceModelWithChallenges(sampleRefs(), statsCategories(), store))
}

func TestReferenceDetail_RelatedChallengeShowsTitleAndStatus(t *testing.T) {
	store := tmpStore(t)

	// Not attempted yet: todo icon plus the real title, no raw ID.
	view := lsDetailWith(t, store)
	for _, want := range []string{"相关挑战", "ls 基础", IconTodo} {
		if !strings.Contains(view, want) {
			t.Fatalf("detail view missing %q\n%s", want, view)
		}
	}
	if strings.Contains(view, "ls-basic") {
		t.Fatalf("raw challenge ID should be replaced by the title\n%s", view)
	}

	store.RecordAttempt("ls-basic", "linux-basics", "导航", true, 0)
	if view := lsDetailWith(t, store); !strings.Contains(view, IconPass) {
		t.Fatalf("passed challenge should show %q\n%s", IconPass, view)
	}

	// A separate store: RecordAttempt never downgrades a pass, so the failed
	// state has to start from a clean slate.
	failed := tmpStore(t)
	failed.RecordAttempt("ls-basic", "linux-basics", "导航", false, 0)
	if view := lsDetailWith(t, failed); !strings.Contains(view, IconFail) {
		t.Fatalf("failed challenge should show %q\n%s", IconFail, view)
	}
}

// Without a challenge lookup (or with an ID missing from it) the list degrades
// to the bare ID rather than dropping the entry.
func TestReferenceDetail_RelatedChallengeFallsBackToID(t *testing.T) {
	if view := openLsDetail(t, NewReferenceModel(sampleRefs())); !strings.Contains(view, "ls-basic") {
		t.Fatalf("detail without lookup should show the raw ID\n%s", view)
	}

	empty := NewReferenceModelWithChallenges(sampleRefs(), nil, tmpStore(t))
	if view := openLsDetail(t, empty); !strings.Contains(view, "ls-basic") {
		t.Fatalf("unknown ID should fall back to the raw ID\n%s", view)
	}
}

func TestReferenceDetail_RelatedChallengeRespectsBoxWidth(t *testing.T) {
	cats := statsCategories()
	cats["linux-basics"][0].Title = strings.Repeat("超长的中文挑战标题", 12)
	m := NewReferenceModelWithChallenges(sampleRefs(), cats, tmpStore(t))

	for _, width := range []int{60, 80, 120} {
		sized, _ := m.Update(tea.WindowSizeMsg{Width: width, Height: 26})
		opened, _ := sized.Update(tea.KeyMsg{Type: tea.KeyEnter})
		for _, line := range strings.Split(opened.(ReferenceModel).View(), "\n") {
			if got := displayWidth(line); got > width {
				t.Fatalf("width %d: line overflows (%d): %q", width, got, line)
			}
		}
	}
}
