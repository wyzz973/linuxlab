package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func runeMatches(msg tea.KeyMsg, values ...string) bool {
	if msg.Type != tea.KeyRunes {
		return false
	}
	got := string(msg.Runes)
	for _, value := range values {
		if got == value {
			return true
		}
	}
	return false
}

// isBackKey reports whether msg is a "go back" key (Esc/Left/q). The
// semantics are defined by globalKeys.Back in keymap.go so key matching and
// the generated help text can never drift apart.
func isBackKey(msg tea.KeyMsg) bool {
	return key.Matches(msg, globalKeys.Back)
}

// isSelectKey reports whether msg is a "confirm" key (Enter/Right); see
// globalKeys.Select in keymap.go.
func isSelectKey(msg tea.KeyMsg) bool {
	return key.Matches(msg, globalKeys.Select)
}

func visibleListRows(height int) int {
	v := height - 12
	if v < 5 {
		return 5
	}
	return v
}

func normalizeVisible(visible int) int {
	if visible < 1 {
		return 1
	}
	return visible
}

func ensureVisible(cursor, offset, count, visible int) int {
	if count <= 0 {
		return 0
	}
	visible = normalizeVisible(visible)
	if cursor < 0 {
		cursor = 0
	}
	if cursor >= count {
		cursor = count - 1
	}
	if count <= visible {
		return 0
	}
	if offset < 0 {
		offset = 0
	}
	maxOffset := count - visible
	if offset > maxOffset {
		offset = maxOffset
	}
	if cursor < offset {
		return cursor
	}
	if cursor >= offset+visible {
		return cursor - visible + 1
	}
	return offset
}

func moveListCursor(cursor, offset, count, visible int, msg tea.KeyMsg, allowRuneNav bool) (int, int, bool) {
	if count <= 0 {
		return 0, 0, false
	}
	visible = normalizeVisible(visible)
	moved := true

	switch {
	case msg.Type == tea.KeyUp || (allowRuneNav && runeMatches(msg, "k")):
		cursor--
		if cursor < 0 {
			cursor = count - 1
		}
	case msg.Type == tea.KeyDown || (allowRuneNav && runeMatches(msg, "j")):
		cursor++
		if cursor >= count {
			cursor = 0
		}
	case msg.Type == tea.KeyPgUp || msg.Type == tea.KeyCtrlB:
		cursor -= visible
		if cursor < 0 {
			cursor = 0
		}
	case msg.Type == tea.KeyPgDown || msg.Type == tea.KeyCtrlF:
		cursor += visible
		if cursor >= count {
			cursor = count - 1
		}
	case msg.Type == tea.KeyHome || msg.Type == tea.KeyCtrlA || (allowRuneNav && runeMatches(msg, "g")):
		cursor = 0
	case msg.Type == tea.KeyEnd || msg.Type == tea.KeyCtrlE || (allowRuneNav && runeMatches(msg, "G")):
		cursor = count - 1
	default:
		moved = false
	}

	if !moved {
		return cursor, offset, false
	}
	return cursor, ensureVisible(cursor, offset, count, visible), true
}

func clampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
