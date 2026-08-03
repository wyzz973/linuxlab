package tui

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// This file holds the three text primitives every layout calculation must go
// through. Width discipline (iron rule #1): display width is always measured
// via ansi.StringWidth / lipgloss.Width (grapheme-aware, ANSI-aware); len(),
// utf8.RuneCount and fmt %-Ns must never participate in layout.

// displayWidth returns the display width of s, ANSI-escape aware and
// CJK / grapheme-cluster correct.
func displayWidth(s string) int {
	return ansi.StringWidth(s)
}

// padRight pads s with trailing spaces until its display width reaches w.
// Strings already at or beyond w are returned unchanged (never truncated).
func padRight(s string, w int) string {
	gap := w - ansi.StringWidth(s)
	if gap <= 0 {
		return s
	}
	return s + strings.Repeat(" ", gap)
}

// truncateWidth truncates s so its display width fits w, appending an
// ellipsis when content is cut. For w <= 1 there is no room for an ellipsis,
// so it degrades to a plain cut. ANSI sequences are preserved and wide runes
// are never split in half.
func truncateWidth(s string, w int) string {
	if w <= 0 {
		return ""
	}
	if ansi.StringWidth(s) <= w {
		return s
	}
	if w <= 1 {
		return ansi.Truncate(s, w, "")
	}
	return ansi.Truncate(s, w, "…")
}

// wrapWidth wraps s to the given display width and returns the resulting
// lines. Wrapping is word-aware for ASCII text and breaks anywhere between
// CJK characters (ansi.Wrap), so Chinese prose without spaces wraps correctly
// and English words are not cut mid-word unless a single word exceeds w.
func wrapWidth(s string, w int) []string {
	if w < 1 {
		return strings.Split(s, "\n")
	}
	return strings.Split(ansi.Wrap(s, w, ""), "\n")
}
