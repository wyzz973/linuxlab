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
)

// Reusable styles.
var (
	BoxStyle      = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(ColorPrimary).Padding(1, 2)
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
