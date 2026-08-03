package tui

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sd3/linuxlab/internal/reference"
)

// manyRefs builds a reference set with count commands (cmd-01 ... cmd-NN).
func manyRefs(count int) *reference.ReferenceData {
	refs := &reference.ReferenceData{}
	for i := 0; i < count; i++ {
		refs.Commands = append(refs.Commands, reference.CommandRef{
			Name:  fmt.Sprintf("cmd-%02d", i+1),
			Brief: fmt.Sprintf("第 %02d 个命令的中文简介", i+1),
		})
	}
	return refs
}

func TestReferenceModel_ListViewShowsCursorBlockAndHint(t *testing.T) {
	m := NewReferenceModel(sampleRefs())
	m, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	view := m.View()
	for _, want := range []string{refCursorBlock, "输入即搜索", "ctrl+u 清空", "3/3"} {
		if !strings.Contains(view, want) {
			t.Fatalf("list view missing %q\n%s", want, view)
		}
	}
}

func TestReferenceModel_ListViewShowsQueryAndMatchCount(t *testing.T) {
	m := NewReferenceModel(sampleRefs())
	m, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m, _ = m.Update(runeKey('l'))
	m, _ = m.Update(runeKey('s'))

	view := m.View()
	if !strings.Contains(view, "ls") {
		t.Fatalf("view should echo the query\n%s", view)
	}
	if !strings.Contains(view, "1/3") {
		t.Fatalf("count label should show matched/total 1/3\n%s", view)
	}
	// The query text keeps the cursor block right after it.
	if !strings.Contains(view, "ls"+refCursorBlock) {
		t.Fatalf("cursor block should follow the query text\n%s", view)
	}
}

func TestReferenceModel_SelectedRowUsesCursorIcon(t *testing.T) {
	m := NewReferenceModel(sampleRefs())
	m, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	view := m.View()
	if !strings.Contains(view, IconCursor) {
		t.Fatalf("selected row should carry the cursor icon %q\n%s", IconCursor, view)
	}

	// Moving the cursor moves the icon onto the second command's row.
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	rm := m.(ReferenceModel)
	row := renderRefRow(rm.filtered[1], true, 60)
	if !strings.Contains(row, IconCursor) || !strings.Contains(row, rm.filtered[1].Name) {
		t.Fatalf("selected row render missing cursor or name: %q", row)
	}
}

func TestReferenceModel_OverflowIndicatorsUseTriangles(t *testing.T) {
	m := NewReferenceModel(manyRefs(15))
	m, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 14}) // 5 visible rows

	view := m.View()
	if !strings.Contains(view, ArrowDown+" 还有 10 个命令") {
		t.Fatalf("expected bottom overflow indicator\n%s", view)
	}
	if strings.Contains(view, ArrowUp+" 还有") {
		t.Fatalf("no top indicator expected at offset 0\n%s", view)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnd})
	view = m.View()
	if !strings.Contains(view, ArrowUp+" 还有 10 个命令") {
		t.Fatalf("expected top overflow indicator after End\n%s", view)
	}
	if strings.Contains(view, ArrowDown+" 还有") {
		t.Fatalf("no bottom indicator expected at the end\n%s", view)
	}
}

func TestReferenceModel_NoResultsShowsEmptyState(t *testing.T) {
	m := NewReferenceModel(sampleRefs())
	m, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	for _, r := range "zzz" {
		m, _ = m.Update(runeKey(r))
	}

	rm := m.(ReferenceModel)
	if len(rm.filtered) != 0 {
		t.Fatalf("expected no matches for zzz, got %d", len(rm.filtered))
	}
	view := m.View()
	for _, want := range []string{"没有匹配「zzz」的命令", "Ctrl+U 清空", "0/3"} {
		if !strings.Contains(view, want) {
			t.Fatalf("empty state missing %q\n%s", want, view)
		}
	}
}

func TestReferenceModel_DetailViewSectionsAndPlaceholders(t *testing.T) {
	m := NewReferenceModel(sampleRefs())
	m, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // open "ls"

	view := m.View()
	for _, want := range []string{"命令详情", "示例", "{{目录}}", "相关挑战", "ls-basic", "▸ 列出所有文件"} {
		if !strings.Contains(view, want) {
			t.Fatalf("detail view missing %q\n%s", want, view)
		}
	}
}

func TestReferenceModel_ViewWidthInvariantAndFrameStable(t *testing.T) {
	for _, size := range []struct{ w, h int }{{60, 20}, {80, 24}, {100, 30}} {
		m := NewReferenceModel(manyRefs(15))
		m, _ = m.Update(tea.WindowSizeMsg{Width: size.w, Height: size.h})

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

// Search-first (locked behavior id=12): printable runes always edit the
// query, so j/k/g/q navigate nothing and quit nothing in list mode.
func TestReferenceModel_RunesGoToQueryNotNavigation(t *testing.T) {
	m := NewReferenceModel(manyRefs(15))
	m, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	for _, r := range "jqgk" {
		var cmd tea.Cmd
		m, cmd = m.Update(runeKey(r))
		if cmd != nil {
			t.Fatalf("typing %q must not emit a command", r)
		}
	}
	rm := m.(ReferenceModel)
	if rm.query != "jqgk" {
		t.Fatalf("query = %q, want %q", rm.query, "jqgk")
	}
	if rm.cursor != 0 {
		t.Fatalf("cursor moved to %d while typing", rm.cursor)
	}
}

func TestReferenceModel_BackspaceEditsQuery(t *testing.T) {
	m := NewReferenceModel(sampleRefs())
	m, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	for _, r := range "lsx" {
		m, _ = m.Update(runeKey(r))
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	rm := m.(ReferenceModel)
	if rm.query != "ls" {
		t.Fatalf("query after backspace = %q, want %q", rm.query, "ls")
	}
	if len(rm.filtered) != 1 || rm.filtered[0].Name != "ls" {
		t.Fatalf("backspace should refresh the filtered list, got %d matches", len(rm.filtered))
	}
}

func TestReferenceModel_LeftOnlyGoesBackWithEmptyQuery(t *testing.T) {
	m := NewReferenceModel(sampleRefs())
	m, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m, _ = m.Update(runeKey('l'))

	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyLeft})
	if cmd != nil {
		t.Fatal("Left with a non-empty query must not go back")
	}
	if rm := m.(ReferenceModel); rm.query != "l" {
		t.Fatalf("query = %q, want %q", rm.query, "l")
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyCtrlU})
	m, cmd = m.Update(tea.KeyMsg{Type: tea.KeyLeft})
	if cmd == nil {
		t.Fatal("Left with an empty query should go back")
	}
	if _, ok := cmd().(GoBackMsg); !ok {
		t.Fatal("expected GoBackMsg from Left on an empty query")
	}
}

func TestReferenceModel_CJKQueryKeepsWidthInvariant(t *testing.T) {
	m := NewReferenceModel(sampleRefs())
	m, _ = m.Update(tea.WindowSizeMsg{Width: 60, Height: 20})
	for _, r := range "目录啊" {
		m, _ = m.Update(runeKey(r))
	}

	view := m.View()
	if !strings.Contains(view, "目录啊"+refCursorBlock) {
		t.Fatalf("cursor block should follow the CJK query\n%s", view)
	}
	for i, line := range strings.Split(view, "\n") {
		if got := lipgloss.Width(line); got > 60 {
			t.Fatalf("line %d width = %d > 60: %q", i, got, line)
		}
	}
}

func TestReferenceModel_ResizeReclampsThroughScrollList(t *testing.T) {
	m := NewReferenceModel(manyRefs(15))
	m, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 20}) // 8 visible
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnd})

	// Shrink: the offset must be re-clamped so the cursor stays visible.
	m, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 14}) // 5 visible
	rm := m.(ReferenceModel)
	vis := rm.maxVisible()
	if rm.cursor < rm.offset || rm.cursor >= rm.offset+vis {
		t.Fatalf("cursor %d not visible in window [%d,%d)", rm.cursor, rm.offset, rm.offset+vis)
	}
}
