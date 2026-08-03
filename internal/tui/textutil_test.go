package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestDisplayWidth(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"", 0},
		{"abc", 3},
		{"挑战", 4},
		{"Vim 操作", 8},
		{S.Title.Render("abc"), 3}, // ANSI sequences are zero-width
		{S.Selected.Render("容器与部署"), 10}, // styled CJK
	}
	for _, c := range cases {
		if got := displayWidth(c.in); got != c.want {
			t.Errorf("displayWidth(%q) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestPadRight(t *testing.T) {
	cases := []string{
		"Vim 操作",
		"Linux 基础命令",
		S.Selected.Render("容器与部署"),
		"plain",
		"",
	}
	for _, s := range cases {
		if got := lipgloss.Width(padRight(s, 18)); got != 18 {
			t.Errorf("padRight(%q, 18) display width = %d, want 18", s, got)
		}
	}

	// Never truncates strings already at or beyond the target width.
	wide := strings.Repeat("宽", 12)
	if got := padRight(wide, 18); got != wide {
		t.Errorf("over-wide string should be unchanged, got %q", got)
	}
	if got := padRight("abcd", 4); got != "abcd" {
		t.Errorf("exact-width string should be unchanged, got %q", got)
	}
	if got := padRight("abcd", 0); got != "abcd" {
		t.Errorf("padRight with 0 width should be a no-op, got %q", got)
	}
}

func TestTruncateWidth(t *testing.T) {
	// CJK: never cuts a wide rune in half, appends an ellipsis.
	got := truncateWidth("容器与部署", 6)
	if w := lipgloss.Width(got); w > 6 {
		t.Fatalf("truncated width = %d, want <= 6 (%q)", w, got)
	}
	if !strings.HasSuffix(got, "…") {
		t.Fatalf("expected ellipsis suffix, got %q", got)
	}

	// Short strings pass through unchanged.
	if got := truncateWidth("abc", 10); got != "abc" {
		t.Fatalf("short string changed: %q", got)
	}

	// Degenerate widths.
	if got := truncateWidth("anything", 0); got != "" {
		t.Fatalf("w=0 should return empty, got %q", got)
	}
	if got := truncateWidth("挑战", -3); got != "" {
		t.Fatalf("negative width should return empty, got %q", got)
	}
	// w=1 degrades to a plain cut (no room for the ellipsis).
	if got := truncateWidth("abcdef", 1); lipgloss.Width(got) > 1 {
		t.Fatalf("w=1 width = %d, want <= 1 (%q)", lipgloss.Width(got), got)
	}
	// w=1 with a 2-cell rune: nothing fits.
	if got := truncateWidth("挑战", 1); lipgloss.Width(got) > 1 {
		t.Fatalf("w=1 CJK width = %d, want <= 1 (%q)", lipgloss.Width(got), got)
	}

	// ANSI-styled input keeps the escape sequences intact and fits the budget.
	styled := S.Error.Render("很长很长很长的一段中文错误信息")
	got = truncateWidth(styled, 10)
	if w := lipgloss.Width(got); w > 10 {
		t.Fatalf("styled truncation width = %d, want <= 10", w)
	}
}

func TestWrapWidth(t *testing.T) {
	// English wraps at word boundaries.
	lines := wrapWidth("hello world", 8)
	if len(lines) != 2 || lines[0] != "hello" || lines[1] != "world" {
		t.Fatalf("word wrap = %q, want [hello world]", lines)
	}

	// CJK breaks between any two characters; every line fits the width.
	lines = wrapWidth("这是一段没有空格的很长的中文句子用于测试折行", 8)
	if len(lines) < 2 {
		t.Fatalf("expected CJK text to wrap, got %q", lines)
	}
	for i, line := range lines {
		if w := lipgloss.Width(line); w > 8 {
			t.Fatalf("line %d width = %d, want <= 8 (%q)", i, w, line)
		}
	}

	// Empty input stays a single empty line.
	if lines := wrapWidth("", 10); len(lines) != 1 || lines[0] != "" {
		t.Fatalf("empty wrap = %q", lines)
	}

	// Non-positive width leaves the text unwrapped instead of panicking.
	if lines := wrapWidth("abc", 0); len(lines) != 1 || lines[0] != "abc" {
		t.Fatalf("w=0 wrap = %q", lines)
	}
}
