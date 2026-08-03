package tui

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- ScrollList -------------------------------------------------------------

func TestScrollList_PagingMatchesLockedSemantics(t *testing.T) {
	// Mirrors ux_test's TestChallengesModel_PageAndEdgeNavigation: 12 items,
	// 5 visible.
	l := ScrollList{}

	if !l.Move(tea.KeyMsg{Type: tea.KeyPgDown}, 12, 5, true) {
		t.Fatal("PgDown should move")
	}
	if l.Cursor != 5 || l.Offset != 1 {
		t.Fatalf("PgDown cursor/offset = %d/%d, want 5/1", l.Cursor, l.Offset)
	}

	l.Move(tea.KeyMsg{Type: tea.KeyEnd}, 12, 5, true)
	if l.Cursor != 11 || l.Offset != 7 {
		t.Fatalf("End cursor/offset = %d/%d, want 11/7", l.Cursor, l.Offset)
	}

	l.Move(runeKey('g'), 12, 5, true)
	if l.Cursor != 0 || l.Offset != 0 {
		t.Fatalf("g cursor/offset = %d/%d, want 0/0", l.Cursor, l.Offset)
	}

	// Up from the top wraps to the bottom (locked wrap-around rule).
	l.Move(tea.KeyMsg{Type: tea.KeyUp}, 12, 5, true)
	if l.Cursor != 11 {
		t.Fatalf("wrap-around cursor = %d, want 11", l.Cursor)
	}
}

func TestScrollList_ClampAfterShrink(t *testing.T) {
	l := ScrollList{Cursor: 11, Offset: 7}
	l.Clamp(12, 3) // viewport shrank from 5 to 3 rows
	if l.Cursor < l.Offset || l.Cursor >= l.Offset+3 {
		t.Fatalf("cursor %d not visible in window [%d,%d)", l.Cursor, l.Offset, l.Offset+3)
	}

	// Item count shrank below the cursor.
	l = ScrollList{Cursor: 11, Offset: 7}
	l.Clamp(4, 5)
	if l.Cursor != 3 {
		t.Fatalf("cursor = %d, want 3 after count shrank to 4", l.Cursor)
	}

	l = ScrollList{Cursor: 2, Offset: 1}
	l.Clamp(0, 5)
	if l.Cursor != 0 || l.Offset != 0 {
		t.Fatalf("empty list should reset to 0/0, got %d/%d", l.Cursor, l.Offset)
	}
}

func TestScrollList_Window(t *testing.T) {
	l := ScrollList{Cursor: 5, Offset: 3}
	start, end := l.Window(12, 5)
	if start != 3 || end != 8 {
		t.Fatalf("window = [%d,%d), want [3,8)", start, end)
	}

	start, end = (&ScrollList{}).Window(2, 5)
	if start != 0 || end != 2 {
		t.Fatalf("short list window = [%d,%d), want [0,2)", start, end)
	}
}

func TestScrollList_RenderRowsIndicatorsAndWidth(t *testing.T) {
	l := ScrollList{Cursor: 5, Offset: 3}
	var selectedRow int = -1
	lines := l.RenderRows(12, 5, 30, func(i int, selected bool) string {
		if selected {
			selectedRow = i
		}
		return fmt.Sprintf("第 %02d 行，这是一个足够长的中文条目内容", i+1)
	})

	if selectedRow != 5 {
		t.Fatalf("selected row = %d, want 5", selectedRow)
	}
	// 5 rows + top indicator + bottom indicator.
	if len(lines) != 7 {
		t.Fatalf("lines = %d, want 7\n%s", len(lines), strings.Join(lines, "\n"))
	}
	if !strings.Contains(lines[0], "还有 3 项") || !strings.Contains(lines[0], ArrowUp) {
		t.Fatalf("top indicator missing: %q", lines[0])
	}
	if !strings.Contains(lines[len(lines)-1], "还有 4 项") {
		t.Fatalf("bottom indicator missing: %q", lines[len(lines)-1])
	}
	for i, line := range lines {
		if w := lipgloss.Width(line); w > 30 {
			t.Fatalf("line %d width = %d, want <= 30", i, w)
		}
	}

	// No indicators when everything is visible.
	lines = (&ScrollList{}).RenderRows(3, 5, 30, func(i int, selected bool) string {
		return fmt.Sprintf("row %d", i)
	})
	if len(lines) != 3 {
		t.Fatalf("full-fit lines = %d, want 3", len(lines))
	}

	if got := (&ScrollList{}).RenderRows(0, 5, 30, nil); got != nil {
		t.Fatalf("empty list should render nil, got %q", got)
	}
}

// --- Header -----------------------------------------------------------------

func TestRenderHeader_WidthAndContent(t *testing.T) {
	crumbs := []string{"LinuxLab", "练习", "Linux 基础命令"}
	right := "Docker ✓ · 12/278"

	for _, width := range []int{140, 80, 60, 30, 12, 1} {
		line := renderHeader(crumbs, right, width)
		if strings.Contains(line, "\n") {
			t.Fatalf("header must be a single line at width %d", width)
		}
		if w := lipgloss.Width(line); w > width {
			t.Fatalf("header width = %d > %d: %q", w, width, line)
		}
	}

	line := renderHeader(crumbs, right, 80)
	for _, want := range []string{"LinuxLab", "练习", "Linux 基础命令", "12/278"} {
		if !strings.Contains(line, want) {
			t.Fatalf("header missing %q: %q", want, line)
		}
	}

	// Narrow: the breadcrumb shrinks from the right, the right slot survives
	// as long as it fits.
	line = renderHeader(crumbs, "12/278", 24)
	if !strings.Contains(line, "12/278") {
		t.Fatalf("narrow header should keep the right slot: %q", line)
	}

	if got := renderHeader(crumbs, right, 0); got != "" {
		t.Fatalf("zero width header = %q, want empty", got)
	}
}

// --- Footer -----------------------------------------------------------------

func TestRenderFooter_CollapsePriority(t *testing.T) {
	// Wide: everything fits, notice included.
	line := renderFooter("已保存进度", challengesKeys, 120)
	for _, want := range []string{"已保存进度", "预览", "选择", "确认", "返回"} {
		if !strings.Contains(line, want) {
			t.Fatalf("wide footer missing %q: %q", want, line)
		}
	}
	if w := lipgloss.Width(line); w > 120 {
		t.Fatalf("footer width = %d > 120", w)
	}

	// Medium: the notice (left slot) is dropped first.
	line = renderFooter("已保存进度", challengesKeys, 46)
	if strings.Contains(line, "已保存进度") {
		t.Fatalf("notice should be dropped before key hints: %q", line)
	}
	if w := lipgloss.Width(line); w > 46 {
		t.Fatalf("footer width = %d > 46", w)
	}

	// Narrow: secondary hints are dropped one by one; the core Back binding
	// (last in ShortHelp) survives the longest.
	line = renderFooter("", challengesKeys, 16)
	if !strings.Contains(line, "返回") {
		t.Fatalf("core back key should survive collapse: %q", line)
	}
	if strings.Contains(line, "预览") {
		t.Fatalf("secondary hint should be dropped at width 16: %q", line)
	}
	if w := lipgloss.Width(line); w > 16 {
		t.Fatalf("footer width = %d > 16", w)
	}

	// Absurdly narrow: still within budget (plain truncation fallback).
	for _, width := range []int{8, 3, 1} {
		line = renderFooter("通知", challengesKeys, width)
		if w := lipgloss.Width(line); w > width {
			t.Fatalf("footer width = %d > %d: %q", w, width, line)
		}
	}
}

func TestCollapseShortHelp_SkipsDisabledBindings(t *testing.T) {
	km := newResultKeyMap()
	km.Next.SetEnabled(false)
	line := collapseShortHelp(km.ShortHelp(), 120)
	if strings.Contains(line, "下一题") {
		t.Fatalf("disabled binding should be hidden: %q", line)
	}
	if !strings.Contains(line, "再试一次") || !strings.Contains(line, "返回") {
		t.Fatalf("enabled bindings missing: %q", line)
	}
}

// --- EmptyState / ProgressSummary ------------------------------------------

func TestEmptyState(t *testing.T) {
	out := EmptyState("还没有练习记录", "先完成一个挑战吧", "Esc 返回")
	for _, want := range []string{"还没有练习记录", "先完成一个挑战吧", "Esc 返回"} {
		if !strings.Contains(out, want) {
			t.Fatalf("empty state missing %q:\n%s", want, out)
		}
	}
	if out := EmptyState("仅标题", "", ""); !strings.Contains(out, "仅标题") || strings.Contains(out, "（") {
		t.Fatalf("title-only empty state malformed:\n%s", out)
	}
}

func TestProgressSummary(t *testing.T) {
	out := ProgressSummary("完成进度", 3, 6, 10)
	for _, want := range []string{"完成进度", "(3/6)", "50%", BarFull, BarEmpty} {
		if !strings.Contains(out, want) {
			t.Fatalf("progress summary missing %q: %q", want, out)
		}
	}
	// Zero total must not divide by zero.
	out = ProgressSummary("", 0, 0, 10)
	if !strings.Contains(out, "(0/0)") {
		t.Fatalf("zero-total summary = %q", out)
	}
}
