package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// menu_view_test.go locks the stage-2 menu rendering (spec §4.1): number
// badges, right-hand status words, and the full-row BgFocus highlight of the
// selected entry.

func sizedMenu(w, h int) MenuModel {
	m := NewMenuModel().(MenuModel)
	m.width = w
	m.height = h
	return m
}

func TestMenuView_SelectedRowIsPaddedToFullInnerWidth(t *testing.T) {
	m := sizedMenu(100, 28)
	innerW := maxInt(boxWidth(m.width)-6, 0)
	labelW, descW := menuColumnWidths(m.items)

	sel := m.renderItem(0, innerW, labelW, descW)
	if got := displayWidth(sel); got != innerW {
		t.Fatalf("selected row width = %d, want the full inner width %d\n%q", got, innerW, sel)
	}
	if !strings.Contains(sel, IconCursor) {
		t.Fatalf("selected row missing cursor %q: %q", IconCursor, sel)
	}

	plain := m.renderItem(1, innerW, labelW, descW)
	if strings.Contains(plain, IconCursor) {
		t.Fatalf("unselected row must not carry the cursor: %q", plain)
	}
	if got := displayWidth(plain); got >= innerW {
		t.Fatalf("unselected row width = %d, should stay below inner width %d", got, innerW)
	}
}

func TestMenuView_ShowsBadgesAndStatusWords(t *testing.T) {
	m := sizedMenu(100, 28)
	view := m.View()
	for _, want := range []string{"[1]", "[2]", "[3]", "[4]", "继续训练", "查漏补缺", "速查手册"} {
		if !strings.Contains(view, want) {
			t.Fatalf("menu view missing %q\n%s", want, view)
		}
	}
}

func TestMenuView_CursorFollowsSelection(t *testing.T) {
	m := sizedMenu(100, 28)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	mm := updated.(MenuModel)

	for _, line := range strings.Split(mm.View(), "\n") {
		if strings.Contains(line, "能力图谱") && !strings.Contains(line, IconCursor) {
			t.Fatalf("selected entry row missing cursor: %q", line)
		}
		if strings.Contains(line, "开始练习") && strings.Contains(line, IconCursor) {
			t.Fatalf("previously selected row still shows cursor: %q", line)
		}
	}
}

func TestMenuView_StatsLineOnlyWithStats(t *testing.T) {
	withStats := NewMenuModelWithStats(278, 5).(MenuModel)
	withStats.width, withStats.height = 100, 28
	if view := withStats.View(); !strings.Contains(view, "278 道挑战 · 5 个模块") {
		t.Fatalf("menu view missing stats line\n%s", view)
	}

	noStats := sizedMenu(100, 28)
	if view := noStats.View(); strings.Contains(view, "道挑战") {
		t.Fatalf("menu without stats must not render a stats line\n%s", view)
	}
}

func TestMenuView_NarrowWidthLinesFit(t *testing.T) {
	m := sizedMenu(60, 20)
	for i, line := range strings.Split(m.View(), "\n") {
		if got := lipgloss.Width(line); got > 60 {
			t.Fatalf("line %d width = %d > 60: %q", i, got, line)
		}
	}
}
