package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sd3/linuxlab/internal/progress"
)

// skillStore builds a store with progress in two categories / three
// subcategories so the skill map has both bars and fold targets.
func skillStore(t *testing.T) *progress.Store {
	t.Helper()
	store := tmpStore(t)
	store.RecordAttempt("ls-basic", "linux-basics", "navigation", true, 0)
	store.RecordAttempt("chmod-basic", "linux-basics", "permissions", false, 0)
	store.RecordAttempt("vim-move", "vim", "movement", false, 1)
	return store
}

func sizedSkillMap(t *testing.T, store *progress.Store, w, h int) tea.Model {
	t.Helper()
	m := NewSkillMapModel(store)
	m, _ = m.Update(tea.WindowSizeMsg{Width: w, Height: h})
	return m
}

func TestSkillMapModel_OverallLineUsesProgressSummary(t *testing.T) {
	m := sizedSkillMap(t, skillStore(t), 80, 24)

	view := m.View()
	// ProgressSummary format: "总体掌握  █░...  33%  (1/3)".
	for _, want := range []string{"总体掌握", BarFull, BarEmpty, "(1/3)"} {
		if !strings.Contains(view, want) {
			t.Fatalf("overall progress line missing %q\n%s", want, view)
		}
	}
}

func TestSkillMapModel_FoldArrowsAndSubcategories(t *testing.T) {
	m := sizedSkillMap(t, skillStore(t), 80, 24)

	view := m.View()
	if !strings.Contains(view, ArrowFold) {
		t.Fatalf("collapsed categories should show %q\n%s", ArrowFold, view)
	}
	if strings.Contains(view, "navigation") {
		t.Fatalf("subcategories must be hidden while collapsed\n%s", view)
	}

	// Expand the first category (linux-basics): arrow flips to ▼ and the
	// subcategory ProgressSummary lines appear with their (passed/total).
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	view = m.View()
	if !strings.Contains(view, ArrowUnfold) {
		t.Fatalf("expanded category should show %q\n%s", ArrowUnfold, view)
	}
	for _, want := range []string{"navigation", "permissions", "(1/1)", "(0/1)"} {
		if !strings.Contains(view, want) {
			t.Fatalf("expanded view missing %q\n%s", want, view)
		}
	}
}

func TestSkillMapModel_ToggleCopiesExpandedMap(t *testing.T) {
	orig := sizedSkillMap(t, skillStore(t), 80, 24).(SkillMapModel)

	updated, _ := orig.Update(tea.KeyMsg{Type: tea.KeySpace})
	next := updated.(SkillMapModel)

	if !next.expanded[0] {
		t.Fatal("expected category 0 expanded on the updated model")
	}
	if orig.expanded[0] {
		t.Fatal("toggle must copy the expanded map, not mutate the original")
	}
}

func TestSkillMapModel_CursorNavigationWraps(t *testing.T) {
	m := sizedSkillMap(t, skillStore(t), 80, 24).(SkillMapModel)
	catCount := len(m.skillMap().Categories)
	if catCount < 2 {
		t.Fatalf("fixture should produce >= 2 categories, got %d", catCount)
	}

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(SkillMapModel)
	if m.cursor != catCount-1 {
		t.Fatalf("Up from top should wrap to %d, got %d", catCount-1, m.cursor)
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(SkillMapModel)
	if m.cursor != 0 {
		t.Fatalf("Down from bottom should wrap to 0, got %d", m.cursor)
	}
}

func TestSkillMapModel_RenderCacheKeyedOnWidthAndVersion(t *testing.T) {
	store := skillStore(t)
	m := sizedSkillMap(t, store, 80, 24).(SkillMapModel)

	m.View()
	if got := m.cache.builds; got != 1 {
		t.Fatalf("first View should build the cache once, builds = %d", got)
	}

	// Cursor movement, fold toggles and repeated renders reuse the cache.
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(SkillMapModel)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(SkillMapModel)
	m.View()
	m.View()
	if got := m.cache.builds; got != 1 {
		t.Fatalf("cursor/toggle/re-render must not rebuild the cache, builds = %d", got)
	}

	// A height-only resize keeps the (width, version) key: still cached.
	updated, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 30})
	m = updated.(SkillMapModel)
	m.View()
	if got := m.cache.builds; got != 1 {
		t.Fatalf("height-only resize must not rebuild the cache, builds = %d", got)
	}

	// A width change invalidates the cache.
	updated, _ = m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m = updated.(SkillMapModel)
	m.View()
	if got := m.cache.builds; got != 2 {
		t.Fatalf("width change should rebuild the cache once, builds = %d", got)
	}

	// New progress bumps the data version: rebuild + fresh numbers.
	store.RecordAttempt("grep-basic", "linux-basics", "text", true, 0)
	view := m.View()
	if got := m.cache.builds; got != 3 {
		t.Fatalf("progress change should rebuild the cache, builds = %d", got)
	}
	if !strings.Contains(view, "(2/4)") {
		t.Fatalf("rebuilt view should reflect the new totals (2/4)\n%s", view)
	}
}

func TestSkillMapModel_EmptyStateOffersNextAction(t *testing.T) {
	m := sizedSkillMap(t, tmpStore(t), 80, 24)

	view := m.View()
	for _, want := range []string{"还没有练习记录", "完成一个挑战", "Esc 返回"} {
		if !strings.Contains(view, want) {
			t.Fatalf("empty state missing %q\n%s", want, view)
		}
	}
}

func TestSkillMapModel_ViewWidthInvariantAndFrameStable(t *testing.T) {
	for _, size := range []struct{ w, h int }{{60, 20}, {80, 24}, {100, 30}, {140, 40}} {
		m := sizedSkillMap(t, skillStore(t), size.w, size.h)
		// Expand the first category so subcategory rows join the invariant check.
		m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})

		view := m.View()
		for i, line := range strings.Split(view, "\n") {
			if got := lipgloss.Width(line); got > size.w {
				t.Fatalf("%dx%d line %d width = %d > %d: %q", size.w, size.h, i, got, size.w, line)
			}
		}
		if again := m.View(); again != view {
			t.Fatalf("same state must render identically at %dx%d", size.w, size.h)
		}
	}
}
