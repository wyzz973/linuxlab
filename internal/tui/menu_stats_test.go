package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sd3/linuxlab/internal/challenge"
)

// The menu's right-hand status words are derived from live counters; the zero
// value keeps the static fallback wording so a data-less menu still reads well.
func TestMenuStats_StatusWordsFromCounters(t *testing.T) {
	cases := []struct {
		name  string
		stats MenuStats
		want  []string
		avoid []string
	}{
		{
			name: "零值回落到静态词",
			want: []string{"继续训练", "查漏补缺", "速查手册"},
		},
		{
			name:  "有题未开始",
			stats: MenuStats{TotalChallenges: 10},
			want:  []string{"未开始"},
			avoid: []string{"继续训练"},
		},
		{
			name:  "部分完成显示比例",
			stats: MenuStats{TotalChallenges: 10, PassedChallenges: 3},
			want:  []string{"3/10"},
		},
		{
			name:  "全部完成",
			stats: MenuStats{TotalChallenges: 10, PassedChallenges: 10},
			want:  []string{"已完成"},
		},
		{
			name:  "薄弱项与速查条数",
			stats: MenuStats{WeakCount: 3, ReferenceCount: 79},
			want:  []string{"3 项薄弱", "79 条"},
			avoid: []string{"查漏补缺", "速查手册"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewMenuModelWithData(tc.stats)
			sized, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 28})
			view := sized.(MenuModel).View()
			for _, want := range tc.want {
				if !strings.Contains(view, want) {
					t.Fatalf("menu view missing %q\n%s", want, view)
				}
			}
			for _, avoid := range tc.avoid {
				if strings.Contains(view, avoid) {
					t.Fatalf("menu view should not contain %q\n%s", avoid, view)
				}
			}
		})
	}
}

// statsCategories mirrors real data closely enough for the weak-spot
// recommender: every challenge carries the subcategory progress is keyed by.
func statsCategories() map[string][]*challenge.Challenge {
	return map[string][]*challenge.Challenge{
		"linux-basics": {
			{ID: "ls-basic", Title: "ls 基础", Category: "linux-basics", Subcategory: "导航", Difficulty: 1},
			{ID: "cp-files", Title: "cp 文件", Category: "linux-basics", Subcategory: "导航", Difficulty: 2},
		},
		"vim": {
			{ID: "vim-move", Title: "Vim 移动", Category: "vim", Subcategory: "移动", Difficulty: 1},
		},
	}
}

func TestComputeMenuStats_ReadsProgressAndReferences(t *testing.T) {
	cats := statsCategories()
	store := tmpStore(t)
	store.RecordAttempt("ls-basic", "linux-basics", "导航", true, 0)
	store.RecordAttempt("cp-files", "linux-basics", "导航", false, 0)

	stats := computeMenuStats(cats, store, sampleRefs())

	total := 0
	for _, chs := range cats {
		total += len(chs)
	}
	if stats.TotalChallenges != total {
		t.Fatalf("TotalChallenges = %d, want %d", stats.TotalChallenges, total)
	}
	if stats.TotalModules != len(cats) {
		t.Fatalf("TotalModules = %d, want %d", stats.TotalModules, len(cats))
	}
	if stats.PassedChallenges != 1 {
		t.Fatalf("PassedChallenges = %d, want 1", stats.PassedChallenges)
	}
	if stats.ReferenceCount != len(sampleRefs().Commands) {
		t.Fatalf("ReferenceCount = %d, want %d", stats.ReferenceCount, len(sampleRefs().Commands))
	}
	if stats.WeakCount == 0 {
		t.Fatal("WeakCount should be non-zero once progress exists")
	}
	if stats.WeakCount > menuWeakLimit {
		t.Fatalf("WeakCount = %d, exceeds limit %d", stats.WeakCount, menuWeakLimit)
	}
}

func TestComputeMenuStats_NilStoreAndRefs(t *testing.T) {
	stats := computeMenuStats(statsCategories(), nil, nil)
	if stats.PassedChallenges != 0 || stats.WeakCount != 0 || stats.ReferenceCount != 0 {
		t.Fatalf("nil store/refs should leave counters zero: %+v", stats)
	}
	if stats.TotalChallenges == 0 {
		t.Fatal("challenge total should still be counted")
	}
}

// Completing a challenge while away from the menu must be reflected when the
// user navigates back — the cached menu model is refreshed, not rebuilt.
func TestAppModel_MenuStatsRefreshOnReturn(t *testing.T) {
	cats := map[string][]*challenge.Challenge{
		"linux-basics": {
			{ID: "c1", Title: "题一", Category: "linux-basics", Subcategory: "s", Difficulty: 1},
			{ID: "c2", Title: "题二", Category: "linux-basics", Subcategory: "s", Difficulty: 1},
		},
	}
	store := tmpStore(t)

	var m tea.Model = NewAppModel(cats, store, nil)
	feed := func(msg tea.Msg) { m, _ = m.Update(msg) }
	feed(tea.WindowSizeMsg{Width: 100, Height: 28})
	feed(tea.KeyMsg{Type: tea.KeyEsc}) // dismiss the opening cover

	if view := m.(AppModel).View(); !strings.Contains(view, "未开始") {
		t.Fatalf("fresh menu should report 未开始\n%s", view)
	}

	// Move the cursor, go into practice, complete a challenge, come back.
	feed(tea.KeyMsg{Type: tea.KeyDown})
	cursorBefore := m.(AppModel).menu.(MenuModel).cursor
	feed(MenuChoiceMsg{Choice: "practice"})
	store.RecordAttempt("c1", "linux-basics", "s", true, 0)
	feed(GoBackMsg{})

	app := m.(AppModel)
	if app.screen != screenMenu {
		t.Fatalf("screen = %v, want menu", app.screen)
	}
	menu := app.menu.(MenuModel)
	if menu.cursor != cursorBefore {
		t.Fatalf("cursor = %d, want %d (position must survive the refresh)", menu.cursor, cursorBefore)
	}
	if view := app.View(); !strings.Contains(view, "1/2") {
		t.Fatalf("menu should report 1/2 after the pass\n%s", view)
	}
}
