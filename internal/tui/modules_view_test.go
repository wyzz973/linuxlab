package tui

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sd3/linuxlab/internal/challenge"
)

// modules_view_test.go locks the stage-2 modules rendering (spec §4.2): rows
// through ScrollList + padRight, the overall bar through ProgressSummary,
// action words 未开始/继续训练/已完成, and the full-row BgFocus highlight.

func sizedModules(t *testing.T, cats map[string][]*challenge.Challenge, w, h int) ModulesModel {
	t.Helper()
	m := NewModulesModel(cats, tmpStore(t)).(ModulesModel)
	m.width = w
	m.height = h
	return m
}

// manyModuleCats builds n single-challenge categories with stable sorted keys.
func manyModuleCats(n int) map[string][]*challenge.Challenge {
	cats := make(map[string][]*challenge.Challenge, n)
	for i := 0; i < n; i++ {
		key := fmt.Sprintf("cat-%02d", i+1)
		cats[key] = []*challenge.Challenge{
			{ID: key + "-1", Title: "挑战 " + key, Category: key, Difficulty: 1},
		}
	}
	return cats
}

func TestModulesView_OverallProgressUsesProgressSummary(t *testing.T) {
	cats := sampleCategories() // linux-basics: 2, vim: 1
	store := tmpStore(t)
	store.RecordAttempt("ls-basic", "linux-basics", "navigation", true, 0)
	m := NewModulesModel(cats, store).(ModulesModel)
	m.width, m.height = 100, 30

	want := ProgressSummary("整体进度", 1, 3, modulesOverallBarWidth)
	if view := m.View(); !strings.Contains(view, want) {
		t.Fatalf("modules view missing the ProgressSummary line %q\n%s", want, view)
	}
}

func TestModulesView_ActionWords(t *testing.T) {
	cats := map[string][]*challenge.Challenge{
		"a-basics": {
			{ID: "a1", Title: "题 A1", Category: "a-basics", Difficulty: 1},
			{ID: "a2", Title: "题 A2", Category: "a-basics", Difficulty: 1},
		},
		"b-vim": {{ID: "b1", Title: "题 B1", Category: "b-vim", Difficulty: 1}},
		"c-ops": {{ID: "c1", Title: "题 C1", Category: "c-ops", Difficulty: 1}},
	}
	store := tmpStore(t)
	store.RecordAttempt("a1", "a-basics", "s", true, 0) // partial -> 继续训练
	store.RecordAttempt("b1", "b-vim", "s", true, 0)    // complete -> 已完成
	m := NewModulesModel(cats, store).(ModulesModel)
	m.width, m.height = 100, 30

	view := m.View()
	for _, want := range []string{"继续训练", "已完成", "未开始"} {
		if !strings.Contains(view, want) {
			t.Fatalf("modules view missing action word %q\n%s", want, view)
		}
	}
}

func TestModulesView_SelectedRowIsPaddedToFullInnerWidth(t *testing.T) {
	m := sizedModules(t, sampleCategories(), 100, 30)
	g := m.rowGeom()

	sel := m.renderModuleRow(0, true, g)
	if got := displayWidth(sel); got != g.innerW {
		t.Fatalf("selected row width = %d, want the full inner width %d\n%q", got, g.innerW, sel)
	}
	if !strings.Contains(sel, IconCursor) {
		t.Fatalf("selected row missing cursor %q: %q", IconCursor, sel)
	}

	plain := m.renderModuleRow(1, false, g)
	if strings.Contains(plain, IconCursor) {
		t.Fatalf("unselected row must not carry the cursor: %q", plain)
	}
	if got := displayWidth(plain); got >= g.innerW {
		t.Fatalf("unselected row width = %d, should stay below inner width %d", got, g.innerW)
	}
}

func TestModulesView_RowsAlignAndCarryBadges(t *testing.T) {
	m := sizedModules(t, sampleCategories(), 100, 30)
	view := m.View()
	for _, want := range []string{"[1] ", "[2] ", "个模块"} {
		if !strings.Contains(view, want) {
			t.Fatalf("modules view missing %q\n%s", want, view)
		}
	}

	// Every module row's bar must start at the same display column.
	var barCols []int
	for _, line := range strings.Split(view, "\n") {
		if !strings.Contains(line, "[1] ") && !strings.Contains(line, "[2] ") {
			continue
		}
		idx := strings.IndexAny(line, BarFull+BarEmpty)
		if idx < 0 {
			t.Fatalf("module row missing progress bar: %q", line)
		}
		barCols = append(barCols, lipgloss.Width(line[:idx]))
	}
	if len(barCols) != 2 {
		t.Fatalf("expected 2 module rows, got %d\n%s", len(barCols), view)
	}
	if barCols[0] != barCols[1] {
		t.Fatalf("progress bars misaligned: columns %v\n%s", barCols, view)
	}
}

func TestModulesView_ScrollListWindowAndIndicators(t *testing.T) {
	m := sizedModules(t, manyModuleCats(12), 80, 24) // 8 visible rows (gapped)

	view := m.View()
	if !strings.Contains(view, ArrowDown+" 还有") {
		t.Fatalf("scrolled-out tail must show the %s indicator\n%s", ArrowDown, view)
	}
	if strings.Contains(view, ArrowUp+" 还有") {
		t.Fatalf("nothing above the window yet, %s indicator unexpected\n%s", ArrowUp, view)
	}

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnd})
	mm := updated.(ModulesModel)
	if mm.cursor != 11 {
		t.Fatalf("End: cursor = %d, want 11", mm.cursor)
	}
	view = mm.View()
	if !strings.Contains(view, ArrowUp+" 还有") {
		t.Fatalf("after End the %s indicator must appear\n%s", ArrowUp, view)
	}
	if strings.Contains(view, ArrowDown+" 还有") {
		t.Fatalf("window is at the bottom, %s indicator unexpected\n%s", ArrowDown, view)
	}
}

func TestModulesView_ResizeReclampsOffset(t *testing.T) {
	m := sizedModules(t, manyModuleCats(12), 80, 24)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnd})
	mm := updated.(ModulesModel)
	if mm.offset == 0 {
		t.Fatal("End on an overflowing list should scroll the window")
	}

	// Growing the terminal makes everything fit: the offset must reset so no
	// rows are stranded above the window.
	updated, _ = mm.Update(tea.WindowSizeMsg{Width: 80, Height: 40})
	mm = updated.(ModulesModel)
	if mm.offset != 0 {
		t.Fatalf("offset = %d after resize, want 0 (all rows fit)", mm.offset)
	}
	if mm.cursor != 11 {
		t.Fatalf("cursor = %d after resize, want 11", mm.cursor)
	}
}

func TestModulesView_EmptyState(t *testing.T) {
	m := sizedModules(t, map[string][]*challenge.Challenge{}, 80, 24)
	view := m.View()
	for _, want := range []string{"暂无训练模块", "Esc 返回"} {
		if !strings.Contains(view, want) {
			t.Fatalf("empty modules view missing %q\n%s", want, view)
		}
	}
}

func TestModulesView_CompactWidthKeepsActionWordVisible(t *testing.T) {
	m := sizedModules(t, sampleCategories(), 60, 20)
	view := m.View()
	for i, line := range strings.Split(view, "\n") {
		if got := lipgloss.Width(line); got > 60 {
			t.Fatalf("line %d width = %d > 60: %q", i, got, line)
		}
	}
	// The adaptive bar width must leave room for the action word even on the
	// narrowest supported terminal.
	if !strings.Contains(view, "未开始") {
		t.Fatalf("compact view lost the action word\n%s", view)
	}
}
