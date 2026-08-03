package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// --- ScrollList -------------------------------------------------------------

// ScrollList extracts the shared cursor+offset+overflow-indicator mechanics of
// the challenge/recommend/reference lists. Movement delegates to
// moveListCursor / ensureVisible so the paging and wrap-around semantics
// locked by ux_test are preserved bit-for-bit.
type ScrollList struct {
	Cursor int
	Offset int
}

// Move applies a navigation key. Returns true when the cursor moved.
func (l *ScrollList) Move(msg tea.KeyMsg, count, visible int, allowRuneNav bool) bool {
	cursor, offset, moved := moveListCursor(l.Cursor, l.Offset, count, visible, msg, allowRuneNav)
	if moved {
		l.Cursor = cursor
		l.Offset = offset
	}
	return moved
}

// Clamp re-clamps the scroll offset so the cursor stays visible (call after a
// resize or after the item count changed).
func (l *ScrollList) Clamp(count, visible int) {
	if count <= 0 {
		l.Cursor = 0
		l.Offset = 0
		return
	}
	if l.Cursor >= count {
		l.Cursor = count - 1
	}
	l.Offset = ensureVisible(l.Cursor, l.Offset, count, visible)
}

// Window returns the half-open visible range [start, end).
func (l ScrollList) Window(count, visible int) (start, end int) {
	visible = normalizeVisible(visible)
	start = l.Offset
	if start > count {
		start = count
	}
	end = start + visible
	if end > count {
		end = count
	}
	return start, end
}

// RenderRows renders the visible window. Each row comes from the callback
// (which is responsible for selection styling) and is truncated to width.
// Overflow indicators ▲/▼ with a remaining-item count are appended above and
// below the window when content is scrolled out of view; they are extra lines
// beyond the `visible` row budget, matching the existing screens' layout.
func (l ScrollList) RenderRows(count, visible, width int, render func(i int, selected bool) string) []string {
	if count <= 0 {
		return nil
	}
	start, end := l.Window(count, visible)

	lines := make([]string, 0, end-start+2)
	if start > 0 {
		lines = append(lines, truncateWidth(S.Dim.Render(fmt.Sprintf("%s 还有 %d 项", ArrowUp, start)), width))
	}
	for i := start; i < end; i++ {
		lines = append(lines, truncateWidth(render(i, i == l.Cursor), width))
	}
	if end < count {
		lines = append(lines, truncateWidth(S.Dim.Render(fmt.Sprintf("%s 还有 %d 项", ArrowDown, count-end)), width))
	}
	return lines
}

// --- Header -----------------------------------------------------------------

// renderHeader renders the 1-line shell header: a breadcrumb on the left
// (truncated from the right when space runs out) and a status label on the
// right ("Docker ✓ · 12/278").
func renderHeader(crumbs []string, right string, width int) string {
	if width < 1 {
		return ""
	}

	sep := S.Breadcrumb.Render(" › ")
	var b strings.Builder
	for i, crumb := range crumbs {
		if i > 0 {
			b.WriteString(sep)
		}
		if i == len(crumbs)-1 {
			b.WriteString(S.Crumb.Render(crumb))
		} else {
			b.WriteString(S.Breadcrumb.Render(crumb))
		}
	}
	left := b.String()

	rightW := displayWidth(right)
	if right == "" || rightW+2 > width {
		// No room for the right slot: breadcrumb takes the whole line.
		return truncateWidth(left, width)
	}

	leftBudget := width - rightW - 2
	left = truncateWidth(left, leftBudget)
	gap := width - displayWidth(left) - rightW
	if gap < 1 {
		gap = 1
	}
	return truncateWidth(left+strings.Repeat(" ", gap)+right, width)
}

// --- Footer -----------------------------------------------------------------

// renderFooter renders the 1-line shell footer: a runtime notice on the left
// (may be empty) and the current screen's short key help on the right. When
// the width is insufficient the left slot is dropped first, then secondary
// key hints are dropped one by one; the last bindings in ShortHelp (the core
// Esc/Enter keys) survive the longest.
func renderFooter(notice string, km interface{ ShortHelp() []key.Binding }, width int) string {
	if width < 1 {
		return ""
	}

	right := collapseShortHelp(km.ShortHelp(), width)
	if notice == "" {
		gap := width - displayWidth(right)
		if gap < 0 {
			gap = 0
		}
		return strings.Repeat(" ", gap) + right
	}

	left := S.Notice.Render(notice)
	leftW := displayWidth(left)
	rightW := displayWidth(right)
	if leftW+2+rightW > width {
		// Drop the left slot first.
		gap := width - rightW
		if gap < 0 {
			gap = 0
		}
		return strings.Repeat(" ", gap) + right
	}
	gap := width - leftW - rightW
	return left + strings.Repeat(" ", gap) + right
}

// collapseShortHelp renders enabled bindings as "key desc" joined by " · ",
// dropping bindings from the FRONT of the list while the line is too wide, so
// callers should order ShortHelp with core keys last.
func collapseShortHelp(bindings []key.Binding, width int) string {
	segments := make([]string, 0, len(bindings))
	for _, b := range bindings {
		if !b.Enabled() {
			continue
		}
		h := b.Help()
		segments = append(segments, S.HelpKey.Render(h.Key)+" "+S.HelpDesc.Render(h.Desc))
	}
	if len(segments) == 0 {
		return ""
	}

	sep := S.Muted.Render(" · ")
	for len(segments) > 1 {
		line := strings.Join(segments, sep)
		if displayWidth(line) <= width {
			return line
		}
		segments = segments[1:]
	}
	return truncateWidth(segments[0], width)
}

// --- EmptyState -------------------------------------------------------------

// EmptyState renders a friendly empty-state block: what happened, what to do
// next, and which key does it. Empty states must always offer a next action.
func EmptyState(title, hint, actionKey string) string {
	var b strings.Builder
	b.WriteString(S.Text.Render(title))
	if hint != "" {
		b.WriteString("\n")
		b.WriteString(S.Dim.Render(hint))
	}
	if actionKey != "" {
		b.WriteString("\n\n")
		b.WriteString(S.Muted.Render("（" + actionKey + "）"))
	}
	return b.String()
}

// --- ProgressSummary --------------------------------------------------------

// ProgressSummary renders the unified progress-bar line:
// "标签  ████░░░░  50%  (3/6)".
func ProgressSummary(label string, passed, total, barWidth int) string {
	ratio := 0.0
	if total > 0 {
		ratio = float64(passed) / float64(total)
	}
	bar := ProgressBar(ratio, barWidth)
	if label == "" {
		return fmt.Sprintf("%s  %3.0f%%  (%d/%d)", bar, ratio*100, passed, total)
	}
	return fmt.Sprintf("%s  %s  %3.0f%%  (%d/%d)", label, bar, ratio*100, passed, total)
}
