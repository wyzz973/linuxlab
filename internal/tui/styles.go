package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Theme colors.
var (
	ColorPrimary = lipgloss.Color("#7aa2f7")
	ColorGreen   = lipgloss.Color("#9ece6a")
	ColorYellow  = lipgloss.Color("#e0af68")
	ColorRed     = lipgloss.Color("#f7768e")
	ColorDim     = lipgloss.Color("#565f89")
	ColorBg      = lipgloss.Color("#1a1b26")
)

// Reusable styles.
var (
	TitleStyle    = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
	SelectedStyle = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
	DimStyle      = lipgloss.NewStyle().Foreground(ColorDim)
	ErrorStyle    = lipgloss.NewStyle().Foreground(ColorRed)
	WarningStyle  = lipgloss.NewStyle().Foreground(ColorYellow)
	HelpStyle     = lipgloss.NewStyle().Foreground(ColorDim)
)

// Status icons.
var (
	PassedIcon  = lipgloss.NewStyle().Foreground(ColorGreen).Render("✓")
	FailedIcon  = lipgloss.NewStyle().Foreground(ColorRed).Render("✗")
	CurrentIcon = lipgloss.NewStyle().Foreground(ColorYellow).Render("▸")
	PendingIcon = lipgloss.NewStyle().Foreground(ColorDim).Render("○")
)

// Progress bar characters.
var (
	ProgressFull  = lipgloss.NewStyle().Foreground(ColorGreen).Render("█")
	ProgressEmpty = lipgloss.NewStyle().Foreground(ColorDim).Render("░")
)

// headerView renders a full-width header bar with colored background.
func headerView(title string, width int) string {
	titleStr := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorBg).
		Background(ColorPrimary).
		Padding(0, 1).
		Render(title)

	titleWidth := lipgloss.Width(titleStr)
	gap := ""
	if width > titleWidth {
		gap = lipgloss.NewStyle().
			Background(ColorPrimary).
			Render(strings.Repeat(" ", width-titleWidth))
	}
	return titleStr + gap
}

// footerView renders a full-width footer bar.
func footerView(help string, width int) string {
	helpStr := lipgloss.NewStyle().
		Foreground(ColorDim).
		Padding(0, 1).
		Render(help)

	helpWidth := lipgloss.Width(helpStr)
	gap := ""
	if width > helpWidth {
		gap = strings.Repeat(" ", width-helpWidth)
	}
	return helpStr + gap
}

// contentStyle returns a style for the content area with padding.
var contentPadding = lipgloss.NewStyle().Padding(1, 2)

// fillContent renders content in a padded area that fills the available height.
func fillContent(content string, width, height int) string {
	return contentPadding.
		Width(width).
		Height(height).
		Render(content)
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
