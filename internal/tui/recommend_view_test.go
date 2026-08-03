package tui

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sd3/linuxlab/internal/challenge"
)

// recFixture builds n recommendation entries with CJK titles.
func recFixture(n int) []*challenge.Challenge {
	cats := []string{"linux-basics", "vim", "shell-scripting"}
	chs := make([]*challenge.Challenge, n)
	for i := 0; i < n; i++ {
		chs[i] = &challenge.Challenge{
			ID:         fmt.Sprintf("rec-%02d", i+1),
			Title:      fmt.Sprintf("推荐 %02d", i+1),
			Category:   cats[i%len(cats)],
			Difficulty: (i % 5) + 1,
		}
	}
	return chs
}

func recModel(n, w, h int) RecommendModel {
	m := NewRecommendModel(recFixture(n)).(RecommendModel)
	m.width = w
	m.height = h
	return m
}

// stepRec feeds one message into the model and returns the updated RecommendModel.
func stepRec(t *testing.T, m RecommendModel, msg tea.Msg) (RecommendModel, tea.Cmd) {
	t.Helper()
	updated, cmd := m.Update(msg)
	rm, ok := updated.(RecommendModel)
	if !ok {
		t.Fatalf("Update returned %T, want RecommendModel", updated)
	}
	return rm, cmd
}

// lineIndexContaining returns the index of the first view line containing sub.
func lineIndexContaining(t *testing.T, view, sub string) int {
	t.Helper()
	for i, line := range strings.Split(view, "\n") {
		if strings.Contains(line, sub) {
			return i
		}
	}
	t.Fatalf("view has no line containing %q\n%s", sub, view)
	return -1
}

func TestRecommendView_EmptyStateOffersNextAction(t *testing.T) {
	m := recModel(0, 80, 24)
	view := m.View()
	for _, want := range []string{"还没有练习记录", "先完成一个挑战", "Esc 返回"} {
		if !strings.Contains(view, want) {
			t.Fatalf("empty-state view missing %q\n%s", want, view)
		}
	}
}

func TestRecommendView_ShowsIntroRowsAndCount(t *testing.T) {
	m := recModel(5, 80, 24)
	view := m.View()

	wants := []string{
		"根据你的练习记录，优先补强以下题目：",
		" 5 题 ",      // count label in the box border
		"Linux 基础命令", // blue category column
		StarOn,       // difficulty gauge
	}
	for i := 1; i <= 5; i++ {
		wants = append(wants, fmt.Sprintf("推荐 %02d", i))
	}
	for _, want := range wants {
		if !strings.Contains(view, want) {
			t.Fatalf("recommend view missing %q\n%s", want, view)
		}
	}
}

func TestRecommendView_CursorMarksSelectedRow(t *testing.T) {
	m := recModel(5, 80, 24)
	m, _ = stepRec(t, m, tea.KeyMsg{Type: tea.KeyDown})

	view := m.View()
	lines := strings.Split(view, "\n")
	selected := lineIndexContaining(t, view, "推荐 02")
	if !strings.Contains(lines[selected], IconCursor) {
		t.Fatalf("selected row must carry the cursor glyph\n%s", lines[selected])
	}
	first := lineIndexContaining(t, view, "推荐 01")
	if strings.Contains(lines[first], IconCursor) {
		t.Fatalf("non-selected row must not carry the cursor glyph\n%s", lines[first])
	}
}

func TestRecommendView_SpacedRowsWhenListFits(t *testing.T) {
	// 5 rows at 80x24 (12-row budget) use the airy layout: one blank line
	// between rows, matching the design mock.
	m := recModel(5, 80, 24)
	view := m.View()
	first := lineIndexContaining(t, view, "推荐 01")
	second := lineIndexContaining(t, view, "推荐 02")
	if second-first != 2 {
		t.Fatalf("spaced layout: rows %d and %d should be 2 lines apart\n%s", first, second, view)
	}
}

func TestRecommendView_CompactWhenListOverflows(t *testing.T) {
	// 20 rows at 80x24 cannot be spaced: rows are contiguous and the overflow
	// indicator appears.
	m := recModel(20, 80, 24)
	view := m.View()
	first := lineIndexContaining(t, view, "推荐 01")
	second := lineIndexContaining(t, view, "推荐 02")
	if second-first != 1 {
		t.Fatalf("compact layout: rows %d and %d should be adjacent\n%s", first, second, view)
	}
	if !strings.Contains(view, ArrowDown+" 还有 8 项") {
		t.Fatalf("expected overflow indicator with the remaining count\n%s", view)
	}
}

func TestRecommend_PagingMirrorsLockedListSemantics(t *testing.T) {
	// Same exact numbers as the challenges list (ux_test): 12 items at
	// height 14 -> 5 visible rows.
	m := recModel(12, 80, 14)

	m, _ = stepRec(t, m, tea.KeyMsg{Type: tea.KeyPgDown})
	if m.cursor != 5 || m.offset != 1 {
		t.Fatalf("PgDown cursor/offset = %d/%d, want 5/1", m.cursor, m.offset)
	}
	m, _ = stepRec(t, m, tea.KeyMsg{Type: tea.KeyEnd})
	if m.cursor != 11 || m.offset != 7 {
		t.Fatalf("End cursor/offset = %d/%d, want 11/7", m.cursor, m.offset)
	}
	m, _ = stepRec(t, m, runeKey('g'))
	if m.cursor != 0 || m.offset != 0 {
		t.Fatalf("g cursor/offset = %d/%d, want 0/0", m.cursor, m.offset)
	}
}

func TestRecommend_EnterSelectsHighlighted(t *testing.T) {
	m := recModel(3, 80, 24)
	m, _ = stepRec(t, m, tea.KeyMsg{Type: tea.KeyDown})

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected Enter to select the highlighted recommendation")
	}
	sel, ok := cmd().(ChallengeSelectedMsg)
	if !ok {
		t.Fatalf("expected ChallengeSelectedMsg, got %T", cmd())
	}
	if sel.Challenge.ID != "rec-02" {
		t.Fatalf("selected = %q, want rec-02", sel.Challenge.ID)
	}
}

func TestRecommend_ResizeReclampsOffset(t *testing.T) {
	m := recModel(12, 80, 20) // 8 visible rows
	m, _ = stepRec(t, m, tea.KeyMsg{Type: tea.KeyEnd})
	if m.cursor != 11 {
		t.Fatalf("cursor = %d, want 11", m.cursor)
	}

	m, _ = stepRec(t, m, tea.WindowSizeMsg{Width: 80, Height: 14}) // 5 visible rows
	maxVis := m.maxVisible()
	if m.cursor < m.offset || m.cursor >= m.offset+maxVis {
		t.Fatalf("cursor %d not visible in window [%d,%d)", m.cursor, m.offset, m.offset+maxVis)
	}
}

func TestRecommendSpaced_Rule(t *testing.T) {
	cases := []struct {
		rows, visible int
		want          bool
	}{
		{0, 12, false},
		{1, 1, true},
		{5, 12, true},  // 9 spaced lines fit a 12-row budget
		{7, 12, false}, // 13 spaced lines do not
		{5, 8, false},
		{3, 5, true},
	}
	for _, c := range cases {
		if got := recommendSpaced(c.rows, c.visible); got != c.want {
			t.Fatalf("recommendSpaced(%d, %d) = %v, want %v", c.rows, c.visible, got, c.want)
		}
	}
}

// --- floating preview panel ---------------------------------------------------

func openRecPreview(t *testing.T, m RecommendModel) RecommendModel {
	t.Helper()
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")})
	if cmd == nil {
		t.Fatal("opening the preview should start the animation")
	}
	m = updated.(RecommendModel)
	for i := 0; i < panelSteps; i++ {
		next, _ := m.Update(panelTickMsg{})
		m = next.(RecommendModel)
	}
	return m
}

func recPreviewModel(t *testing.T, n int) RecommendModel {
	t.Helper()
	store := tmpStore(t)
	recs := recFixture(n)
	store.RecordAttempt(recs[0].ID, recs[0].Category, recs[0].Subcategory, false, 0)
	m := NewRecommendModelWithStore(recs, store)
	sized, _ := m.Update(tea.WindowSizeMsg{Width: 110, Height: 30})
	return sized.(RecommendModel)
}

// The panel floats over the list; the list keeps its own layout.
func TestRecommendPreview_DoesNotRelayoutTheList(t *testing.T) {
	m := recPreviewModel(t, 5)
	before := strings.Split(m.View(), "\n")
	after := strings.Split(openRecPreview(t, m).View(), "\n")

	if len(before) != len(after) {
		t.Fatalf("preview changed the list height: %d -> %d", len(before), len(after))
	}
	for i := range before {
		if a, b := displayWidth(before[i]), displayWidth(after[i]); a != b {
			t.Fatalf("line %d width changed %d -> %d", i, a, b)
		}
	}
}

func TestRecommendPreview_ShowsSelectedChallengeState(t *testing.T) {
	m := recPreviewModel(t, 5)
	view := openRecPreview(t, m).View()
	for _, want := range []string{"预览", m.challenges[0].Title, "状态", "未通过", "难度"} {
		if !strings.Contains(view, want) {
			t.Fatalf("panel missing %q\n%s", want, view)
		}
	}
}

// Without a store the panel omits the state row rather than guessing it.
func TestRecommendPreview_OmitsStateWithoutStore(t *testing.T) {
	m := NewRecommendModel(recFixture(5))
	sized, _ := m.Update(tea.WindowSizeMsg{Width: 110, Height: 30})
	view := openRecPreview(t, sized.(RecommendModel)).View()
	if !strings.Contains(view, "预览") {
		t.Fatalf("panel should still open without a store\n%s", view)
	}
	if strings.Contains(view, "状态") {
		t.Fatalf("state row should be omitted without a store\n%s", view)
	}
}

func TestRecommendPreview_EscClosesPanelFirst(t *testing.T) {
	m := openRecPreview(t, recPreviewModel(t, 5))

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = updated.(RecommendModel)
	if cmd != nil {
		if _, ok := cmd().(GoBackMsg); ok {
			t.Fatal("Esc closed the panel and left the screen in one press")
		}
	}
	for i := 0; i < panelSteps; i++ {
		next, _ := m.Update(panelTickMsg{})
		m = next.(RecommendModel)
	}
	if strings.Contains(m.View(), "预览") {
		t.Fatalf("the panel should be gone after closing\n%s", m.View())
	}
	if _, cmd = m.Update(tea.KeyMsg{Type: tea.KeyEsc}); cmd == nil {
		t.Fatal("Esc with the panel closed should leave the screen")
	}
}

func TestRecommendPreview_HonorsFrameInvariants(t *testing.T) {
	for _, sz := range [][2]int{{60, 18}, {80, 24}, {110, 30}, {160, 44}} {
		store := tmpStore(t)
		m := NewRecommendModelWithStore(recFixture(5), store)
		sized, _ := m.Update(tea.WindowSizeMsg{Width: sz[0], Height: sz[1]})
		rm := sized.(RecommendModel)
		updated, _ := rm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")})
		rm = updated.(RecommendModel)

		for frame := 0; frame <= panelSteps; frame++ {
			for i, line := range strings.Split(rm.View(), "\n") {
				if got := displayWidth(line); got > sz[0] {
					t.Fatalf("%dx%d frame %d: line %d width %d > %d: %q",
						sz[0], sz[1], frame, i, got, sz[0], line)
				}
			}
			next, _ := rm.Update(panelTickMsg{})
			rm = next.(RecommendModel)
		}
	}
}
