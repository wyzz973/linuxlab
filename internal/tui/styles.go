package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Theme colors — semantic, high-contrast palette.
var (
	ColorPrimary   = lipgloss.Color("#7aa2f7") // Blue — titles, selected items
	ColorSecondary = lipgloss.Color("#bb9af7") // Purple — categories
	ColorGreen     = lipgloss.Color("#9ece6a") // Green — passed, success
	ColorYellow    = lipgloss.Color("#e0af68") // Yellow — stars, hints, warnings
	ColorRed       = lipgloss.Color("#f7768e") // Red — failed, errors
	ColorDim       = lipgloss.Color("#565f89") // Gray — borders, secondary text
	ColorText      = lipgloss.Color("#c0caf5") // Light — primary text
	ColorSubtle    = lipgloss.Color("#414868") // Dark gray — very subtle
	ColorBg        = lipgloss.Color("#1a1b26")
)

// Reusable styles.
var (
	TitleStyle    = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
	SelectedStyle = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
	DimStyle      = lipgloss.NewStyle().Foreground(ColorDim)
	SubtleStyle   = lipgloss.NewStyle().Foreground(ColorSubtle)
	TextStyle     = lipgloss.NewStyle().Foreground(ColorText)
	ErrorStyle    = lipgloss.NewStyle().Foreground(ColorRed)
	SuccessStyle  = lipgloss.NewStyle().Foreground(ColorGreen)
	WarningStyle  = lipgloss.NewStyle().Foreground(ColorYellow)
	HelpStyle     = lipgloss.NewStyle().Foreground(ColorDim)
)

// Status icons.
var (
	PassedIcon  = lipgloss.NewStyle().Foreground(ColorGreen).Render("✓")
	FailedIcon  = lipgloss.NewStyle().Foreground(ColorRed).Render("✗")
	CurrentIcon = lipgloss.NewStyle().Foreground(ColorPrimary).Render("›")
	PendingIcon = lipgloss.NewStyle().Foreground(ColorDim).Render("○")
)

// Progress bar characters.
var (
	ProgressFull  = lipgloss.NewStyle().Foreground(ColorGreen).Render("█")
	ProgressEmpty = lipgloss.NewStyle().Foreground(ColorSubtle).Render("░")
)

// boxWidth calculates the content box width, clamped to [50, 90].
func boxWidth(termWidth int) int {
	w := termWidth - 2 // 1 char margin each side
	if w > 90 {
		w = 90
	}
	if w < 50 {
		w = 50
	}
	return w
}

// contentBox renders a rounded-border box with a title in the top border.
// The title is rendered as a custom top line to avoid ANSI escape code corruption.
func contentBox(title string, body string, termWidth, termHeight int, rightLabel string) string {
	w := boxWidth(termWidth)
	innerW := w - 6 // border (2) + padding (4)

	// Build the top border line manually: ╭─ title ─────── rightLabel ─╮
	topLeft := DimStyle.Render("╭─")
	topRight := DimStyle.Render("─╮")
	titleStr := ""
	if title != "" {
		titleStr = " " + TitleStyle.Render(title) + " "
	}
	rightStr := ""
	if rightLabel != "" {
		rightStr = " " + DimStyle.Render(rightLabel) + " "
	}

	titleW := lipgloss.Width(titleStr)
	rightW := lipgloss.Width(rightStr)
	topLeftW := lipgloss.Width(topLeft)
	topRightW := lipgloss.Width(topRight)
	fillW := w - topLeftW - titleW - rightW - topRightW
	if fillW < 0 {
		fillW = 0
	}

	topLine := topLeft + titleStr + DimStyle.Render(strings.Repeat("─", fillW)) + rightStr + topRight

	// Build body with side borders
	bodyLines := strings.Split(body, "\n")
	var middle strings.Builder
	for _, line := range bodyLines {
		lineW := lipgloss.Width(line)
		pad := innerW - lineW
		if pad < 0 {
			pad = 0
		}
		middle.WriteString(DimStyle.Render("│") + "  " + line + strings.Repeat(" ", pad) + "  " + DimStyle.Render("│") + "\n")
	}

	// Empty padding line top and bottom inside box
	emptyLine := DimStyle.Render("│") + strings.Repeat(" ", innerW+4) + DimStyle.Render("│")

	// Bottom border
	bottomLine := DimStyle.Render("╰" + strings.Repeat("─", w-2) + "╯")

	rendered := topLine + "\n" + emptyLine + "\n" + middle.String() + emptyLine + "\n" + bottomLine

	// Horizontal centering
	renderedWidth := lipgloss.Width(topLine) // use top line as width reference
	leftPad := 0
	if termWidth > renderedWidth {
		leftPad = (termWidth - renderedWidth) / 2
	}
	if leftPad > 0 {
		padStr := strings.Repeat(" ", leftPad)
		lines := strings.Split(rendered, "\n")
		for i, line := range lines {
			lines[i] = padStr + line
		}
		rendered = strings.Join(lines, "\n")
	}

	return rendered
}

// statusBar renders a bottom bar outside the box: left-aligned context, right-aligned keys.
func statusBar(left, right string, termWidth int) string {
	w := boxWidth(termWidth)
	leftPad := 0
	if termWidth > w {
		leftPad = (termWidth - w) / 2
	}

	leftRendered := DimStyle.Render(left)
	rightRendered := DimStyle.Render(right)

	leftW := lipgloss.Width(leftRendered)
	rightW := lipgloss.Width(rightRendered)
	gap := w - leftW - rightW
	if gap < 1 {
		gap = 1
	}

	padding := ""
	if leftPad > 0 {
		padding = strings.Repeat(" ", leftPad)
	}

	return padding + leftRendered + strings.Repeat(" ", gap) + rightRendered
}

// verticalCenter wraps the box + status bar in vertical padding to center on screen.
func verticalCenter(box, status string, termHeight int) string {
	boxH := lipgloss.Height(box)
	totalH := boxH + 1 // +1 for status bar
	topPad := (termHeight - totalH) / 2
	if topPad < 0 {
		topPad = 0
	}
	return strings.Repeat("\n", topPad) + box + "\n" + status
}

// sectionTitle renders a labeled section divider: "── 标题 ────────────"
func sectionTitle(title string, width int) string {
	inner := boxWidth(width) - 4 // account for padding
	prefix := fmt.Sprintf("── %s ", title)
	prefixW := lipgloss.Width(prefix)
	remaining := maxInt(0, inner-prefixW)
	return DimStyle.Render(prefix + strings.Repeat("─", remaining))
}

// DifficultyStars returns a star string for the given difficulty level (1-5).
func DifficultyStars(level int) string {
	if level < 1 {
		level = 1
	}
	if level > 5 {
		level = 5
	}
	filled := lipgloss.NewStyle().Foreground(ColorYellow).Render("★")
	empty := lipgloss.NewStyle().Foreground(ColorDim).Render("☆")

	var b strings.Builder
	for i := 0; i < level; i++ {
		b.WriteString(filled)
	}
	for i := level; i < 5; i++ {
		b.WriteString(empty)
	}
	return b.String()
}

// ProgressBar returns a progress bar string of the given width for the given ratio (0.0-1.0).
func ProgressBar(ratio float64, width int) string {
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	if width < 1 {
		return ""
	}

	filled := int(ratio * float64(width))
	if filled > width {
		filled = width
	}

	var b strings.Builder
	for i := 0; i < filled; i++ {
		b.WriteString(ProgressFull)
	}
	for i := filled; i < width; i++ {
		b.WriteString(ProgressEmpty)
	}
	return b.String()
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
