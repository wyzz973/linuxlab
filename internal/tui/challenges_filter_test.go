package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sd3/linuxlab/internal/challenge"
)

// filterFixture returns a challenge list with distinctive titles and
// subcategories for exercising the '/' filter mode.
func filterFixture() []*challenge.Challenge {
	return []*challenge.Challenge{
		{ID: "chmod-num", Title: "使用数字模式修改权限 - chmod", Category: "linux-basics", Subcategory: "permissions", Difficulty: 2},
		{ID: "chmod-sym", Title: "使用符号模式修改权限 - chmod u+x", Category: "linux-basics", Subcategory: "permissions", Difficulty: 2},
		{ID: "awk-basic", Title: "字段提取 - awk 基础", Category: "linux-basics", Subcategory: "text-processing", Difficulty: 2},
		{ID: "apt-basics", Title: "APT 包管理基础", Category: "linux-basics", Subcategory: "package-management", Difficulty: 1},
		{ID: "cat-view", Title: "查看和合并文件内容 - cat", Category: "linux-basics", Subcategory: "file-operations", Difficulty: 1},
	}
}

func filterModel(t *testing.T) ChallengesModel {
	t.Helper()
	m := NewChallengesModel("linux-basics", filterFixture(), tmpStore(t)).(ChallengesModel)
	m.width = 80
	m.height = 24
	return m
}

// stepCh feeds one key into the model and returns the updated ChallengesModel.
func stepCh(t *testing.T, m ChallengesModel, msg tea.Msg) (ChallengesModel, tea.Cmd) {
	t.Helper()
	updated, cmd := m.Update(msg)
	cm, ok := updated.(ChallengesModel)
	if !ok {
		t.Fatalf("Update returned %T, want ChallengesModel", updated)
	}
	return cm, cmd
}

func typeString(t *testing.T, m ChallengesModel, s string) ChallengesModel {
	t.Helper()
	for _, r := range s {
		m, _ = stepCh(t, m, runeKey(r))
	}
	return m
}

func TestChallengesFilter_SlashEntersFilterMode(t *testing.T) {
	m := filterModel(t)

	m, _ = stepCh(t, m, runeKey('/'))
	if !m.filtering {
		t.Fatal("expected / to enter filter mode")
	}
	if m.query != "" {
		t.Fatalf("query = %q, want empty on entry", m.query)
	}
	if len(m.filtered) != len(m.challenges) {
		t.Fatalf("filtered = %d items, want the full list (%d)", len(m.filtered), len(m.challenges))
	}
}

func TestChallengesFilter_TypeToFilterNarrowsAndResetsCursor(t *testing.T) {
	m := filterModel(t)
	m, _ = stepCh(t, m, tea.KeyMsg{Type: tea.KeyDown}) // cursor -> 1
	m, _ = stepCh(t, m, runeKey('/'))

	m = typeString(t, m, "chmod")
	if len(m.filtered) != 2 {
		t.Fatalf("filtered = %d, want 2 chmod matches", len(m.filtered))
	}
	if m.cursor != 0 || m.offset != 0 {
		t.Fatalf("cursor/offset = %d/%d, want 0/0 after query change", m.cursor, m.offset)
	}
	if got := m.itemAt(0).ID; got != "chmod-num" {
		t.Fatalf("first match = %q, want chmod-num", got)
	}
}

func TestChallengesFilter_MatchesSubcategory(t *testing.T) {
	m := filterModel(t)
	m, _ = stepCh(t, m, runeKey('/'))

	m = typeString(t, m, "package")
	if len(m.filtered) != 1 {
		t.Fatalf("filtered = %d, want 1 subcategory match", len(m.filtered))
	}
	if got := m.itemAt(0).ID; got != "apt-basics" {
		t.Fatalf("match = %q, want apt-basics", got)
	}
}

func TestChallengesFilter_CaseInsensitive(t *testing.T) {
	m := filterModel(t)
	m, _ = stepCh(t, m, runeKey('/'))

	m = typeString(t, m, "apt")
	if len(m.filtered) != 1 || m.itemAt(0).ID != "apt-basics" {
		t.Fatalf("lower-case query should match the upper-case title, got %d matches", len(m.filtered))
	}
}

func TestChallengesFilter_RuneNavKeysGoIntoQuery(t *testing.T) {
	// Search-first semantics: while filtering, j/k/g/q are query characters,
	// not navigation or back keys.
	m := filterModel(t)
	m, _ = stepCh(t, m, runeKey('/'))

	var cmd tea.Cmd
	for _, r := range "jkgq" {
		m, cmd = stepCh(t, m, runeKey(r))
		if cmd != nil {
			t.Fatalf("rune %q must not produce a command while filtering", r)
		}
	}
	if m.query != "jkgq" {
		t.Fatalf("query = %q, want %q", m.query, "jkgq")
	}
}

func TestChallengesFilter_SpaceGoesIntoQuery(t *testing.T) {
	m := filterModel(t)
	m, _ = stepCh(t, m, runeKey('/'))
	m = typeString(t, m, "awk")
	m, _ = stepCh(t, m, tea.KeyMsg{Type: tea.KeySpace})
	m = typeString(t, m, "基础")

	if m.query != "awk 基础" {
		t.Fatalf("query = %q, want %q", m.query, "awk 基础")
	}
	if len(m.filtered) != 1 || m.itemAt(0).ID != "awk-basic" {
		t.Fatalf("expected the CJK query to match awk-basic, got %d matches", len(m.filtered))
	}
}

func TestChallengesFilter_ArrowsNavigateWithinMatches(t *testing.T) {
	m := filterModel(t)
	m, _ = stepCh(t, m, runeKey('/'))
	m = typeString(t, m, "chmod")

	m, _ = stepCh(t, m, tea.KeyMsg{Type: tea.KeyDown})
	if m.cursor != 1 {
		t.Fatalf("cursor = %d, want 1 after Down", m.cursor)
	}
	m, _ = stepCh(t, m, tea.KeyMsg{Type: tea.KeyUp})
	if m.cursor != 0 {
		t.Fatalf("cursor = %d, want 0 after Up", m.cursor)
	}
}

func TestChallengesFilter_EnterOpensHighlightedAndExitsFilter(t *testing.T) {
	m := filterModel(t)
	m, _ = stepCh(t, m, runeKey('/'))
	m = typeString(t, m, "chmod")
	m, _ = stepCh(t, m, tea.KeyMsg{Type: tea.KeyDown}) // second match

	m, cmd := stepCh(t, m, tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected Enter to select the highlighted challenge")
	}
	sel, ok := cmd().(ChallengeSelectedMsg)
	if !ok {
		t.Fatalf("expected ChallengeSelectedMsg, got %T", cmd())
	}
	if sel.Challenge.ID != "chmod-sym" {
		t.Fatalf("selected = %q, want chmod-sym", sel.Challenge.ID)
	}
	if m.filtering {
		t.Fatal("Enter must leave filter mode so going back lands on the full list")
	}
	// The cursor must map back to the selected challenge's full-list index.
	if m.cursor != 1 {
		t.Fatalf("cursor = %d, want full-list index 1 after commit", m.cursor)
	}
}

func TestChallengesFilter_EnterWithNoMatchDoesNothing(t *testing.T) {
	m := filterModel(t)
	m, _ = stepCh(t, m, runeKey('/'))
	m = typeString(t, m, "zzz-no-match")
	if len(m.filtered) != 0 {
		t.Fatalf("filtered = %d, want 0", len(m.filtered))
	}

	m, cmd := stepCh(t, m, tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Fatal("Enter with no match must not select anything")
	}
	if !m.filtering {
		t.Fatal("Enter with no match must stay in filter mode")
	}
}

func TestChallengesFilter_EscRestoresFullList(t *testing.T) {
	m := filterModel(t)
	m, _ = stepCh(t, m, runeKey('/'))
	m = typeString(t, m, "chmod")
	m, _ = stepCh(t, m, tea.KeyMsg{Type: tea.KeyDown}) // highlight chmod-sym

	m, cmd := stepCh(t, m, tea.KeyMsg{Type: tea.KeyEsc})
	if cmd != nil {
		t.Fatal("Esc while filtering must exit the filter, not leave the screen")
	}
	if m.filtering {
		t.Fatal("expected Esc to exit filter mode")
	}
	if m.query != "" || m.filtered != nil {
		t.Fatalf("filter state not cleared: query=%q filtered=%v", m.query, m.filtered)
	}
	if got := m.listLen(); got != len(filterFixture()) {
		t.Fatalf("list length = %d, want the full list (%d)", got, len(filterFixture()))
	}
	// The cursor follows the challenge that was highlighted in the filter.
	if m.cursor != 1 {
		t.Fatalf("cursor = %d, want full-list index 1 (chmod-sym)", m.cursor)
	}

	// A second Esc now leaves the screen.
	_, cmd = stepCh(t, m, tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected Esc after exiting the filter to go back")
	}
	if _, ok := cmd().(GoBackMsg); !ok {
		t.Fatalf("expected GoBackMsg, got %T", cmd())
	}
}

func TestChallengesFilter_BackspaceRecomputesMatches(t *testing.T) {
	m := filterModel(t)
	m, _ = stepCh(t, m, runeKey('/'))
	m = typeString(t, m, "chmod u")
	if len(m.filtered) != 1 {
		t.Fatalf("filtered = %d, want 1", len(m.filtered))
	}

	m, _ = stepCh(t, m, tea.KeyMsg{Type: tea.KeyBackspace})
	m, _ = stepCh(t, m, tea.KeyMsg{Type: tea.KeyBackspace})
	if m.query != "chmod" {
		t.Fatalf("query = %q, want %q", m.query, "chmod")
	}
	if len(m.filtered) != 2 {
		t.Fatalf("filtered = %d, want 2 after backspace", len(m.filtered))
	}
}

func TestChallengesFilter_CtrlUClearsQuery(t *testing.T) {
	m := filterModel(t)
	m, _ = stepCh(t, m, runeKey('/'))
	m = typeString(t, m, "chmod")

	m, _ = stepCh(t, m, tea.KeyMsg{Type: tea.KeyCtrlU})
	if m.query != "" {
		t.Fatalf("query = %q, want empty after Ctrl+U", m.query)
	}
	if len(m.filtered) != len(m.challenges) {
		t.Fatalf("filtered = %d, want the full list restored", len(m.filtered))
	}
	if !m.filtering {
		t.Fatal("Ctrl+U must stay in filter mode")
	}
}

func TestChallengesFilter_ViewShowsQueryAndCount(t *testing.T) {
	m := filterModel(t)
	m, _ = stepCh(t, m, runeKey('/'))
	m = typeString(t, m, "chmod")

	view := m.View()
	for _, want := range []string{"chmod", "（2/5）", " 2/5 "} {
		if !strings.Contains(view, want) {
			t.Fatalf("filter view missing %q\n%s", want, view)
		}
	}
	// Filtered-out rows must not render.
	if strings.Contains(view, "APT 包管理基础") {
		t.Fatalf("filter view leaked a non-matching row\n%s", view)
	}
}

func TestChallengesFilter_NoMatchShowsEmptyState(t *testing.T) {
	m := filterModel(t)
	m, _ = stepCh(t, m, runeKey('/'))
	m = typeString(t, m, "xyz")

	view := m.View()
	for _, want := range []string{"没有匹配", "xyz", "退出过滤"} {
		if !strings.Contains(view, want) {
			t.Fatalf("no-match view missing %q\n%s", want, view)
		}
	}
}

func TestChallengesFilter_FocusChallengeDropsFilter(t *testing.T) {
	// Coming back from the detail screen always lands on the full list.
	m := filterModel(t)
	m, _ = stepCh(t, m, runeKey('/'))
	m = typeString(t, m, "chmod")

	m = m.focusChallenge("awk-basic")
	if m.filtering {
		t.Fatal("focusChallenge must drop an active filter")
	}
	if m.cursor != 2 {
		t.Fatalf("cursor = %d, want 2 (awk-basic)", m.cursor)
	}
}

func TestChallengesView_NonFilterShowsWindowLabelAndCursor(t *testing.T) {
	m := filterModel(t)
	view := m.View()
	if !strings.Contains(view, " 1-5/5 ") {
		t.Fatalf("expected window label 1-5/5 in the box border\n%s", view)
	}
	if !strings.Contains(view, IconCursor) {
		t.Fatalf("expected the cursor glyph on the selected row\n%s", view)
	}
	// The progress summary replaces the filter input line outside filter mode.
	if !strings.Contains(view, "完成进度") {
		t.Fatalf("expected the progress summary line\n%s", view)
	}
}

func TestChallengesFilter_NonFilterPagingUnchanged(t *testing.T) {
	// Guard: adding filter mode must not disturb the locked paging semantics
	// (mirrors ux_test's exact numbers on a 12-item list at height 14).
	store := tmpStore(t)
	m := NewChallengesModel("linux-basics", manyChallenges(12), store).(ChallengesModel)
	m.height = 14

	m, _ = stepCh(t, m, tea.KeyMsg{Type: tea.KeyPgDown})
	if m.cursor != 5 || m.offset != 1 {
		t.Fatalf("PgDown cursor/offset = %d/%d, want 5/1", m.cursor, m.offset)
	}
	m, _ = stepCh(t, m, tea.KeyMsg{Type: tea.KeyEnd})
	if m.cursor != 11 || m.offset != 7 {
		t.Fatalf("End cursor/offset = %d/%d, want 11/7", m.cursor, m.offset)
	}
	m, _ = stepCh(t, m, runeKey('g'))
	if m.cursor != 0 || m.offset != 0 {
		t.Fatalf("g cursor/offset = %d/%d, want 0/0", m.cursor, m.offset)
	}
}

func TestChallengesFilter_FilteredListScrolls(t *testing.T) {
	// 12 items all matching "挑战": the filtered window must page exactly like
	// the plain list.
	store := tmpStore(t)
	m := NewChallengesModel("linux-basics", manyChallenges(12), store).(ChallengesModel)
	m.height = 14 // 5 visible rows

	m, _ = stepCh(t, m, runeKey('/'))
	m = typeString(t, m, "挑战")
	if len(m.filtered) != 12 {
		t.Fatalf("filtered = %d, want 12", len(m.filtered))
	}

	m, _ = stepCh(t, m, tea.KeyMsg{Type: tea.KeyPgDown})
	if m.cursor != 5 || m.offset != 1 {
		t.Fatalf("filtered PgDown cursor/offset = %d/%d, want 5/1", m.cursor, m.offset)
	}

	view := m.View()
	if !strings.Contains(view, "还有") {
		t.Fatalf("scrolled filter view should show an overflow indicator\n%s", view)
	}
}
