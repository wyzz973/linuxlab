package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Legacy color aliases. New code should use the semantic tokens in theme.go;
// these names remain so screens can migrate incrementally.
var (
	ColorPrimary   = currentTheme.Primary
	ColorSecondary = currentTheme.Secondary
	ColorGreen     = currentTheme.Success
	ColorYellow    = currentTheme.Warning
	ColorRed       = currentTheme.Error
	ColorDim       = currentTheme.FgSubtle
	ColorText      = currentTheme.FgBase
	ColorSubtle    = currentTheme.Border
	ColorBg        = currentTheme.BgBase
	ColorPanel     = currentTheme.BgPanel
	ColorFocusBg   = currentTheme.BgFocus
)

// Legacy style aliases pointing at the centralized style set (theme.go).
var (
	TitleStyle    = S.Title
	SelectedStyle = S.Selected
	DimStyle      = S.Dim
	SubtleStyle   = S.Subtle
	TextStyle     = S.Text
	ErrorStyle    = S.Error
	SuccessStyle  = S.Success
	WarningStyle  = S.Warning
	HelpStyle     = S.Help
	KeyStyle      = S.Key
	BadgeStyle    = S.Badge
	MetaStyle     = S.Meta
)

// Status icons (pre-rendered).
var (
	PassedIcon  = S.IconPass
	FailedIcon  = S.IconFail
	CurrentIcon = S.IconCursor
	PendingIcon = S.IconTodo
)

// Progress bar characters (pre-rendered).
var (
	ProgressFull  = S.BarFull
	ProgressEmpty = S.BarEmpty
)

// minBoxWidth is the safe lower bound for a rendered box: border (2) +
// padding (4) + at least 2 columns of content. Anything narrower would make
// strings.Repeat counts negative.
const minBoxWidth = 8

// boxWidth calculates the content box width. It prefers a readable 50-90 column
// box, but respects narrow terminals instead of forcing horizontal overflow.
// The returned width is never smaller than minBoxWidth.
func boxWidth(termWidth int) int {
	if termWidth <= 0 {
		return 80
	}
	w := termWidth - 2 // 1 char margin each side
	if w > 90 {
		w = 90
	}
	if termWidth < 52 {
		if w < 20 {
			return maxInt(termWidth, minBoxWidth)
		}
		return w
	}
	if w < 50 {
		w = 50
	}
	return w
}

// truncateToWidth is a legacy alias for truncateWidth (textutil.go).
func truncateToWidth(s string, maxW int) string {
	return truncateWidth(s, maxW)
}

// contentBox renders a rounded-border box with a title in the top border.
// The title is rendered as a custom top line to avoid ANSI escape code
// corruption (manual borders per CLAUDE.md convention). Every produced line
// has the exact display width of the box: overlong titles and body lines are
// truncated, short body lines padded.
func contentBox(title string, body string, termWidth, termHeight int, rightLabel string) string {
	w := maxInt(boxWidth(termWidth), minBoxWidth)
	innerW := maxInt(w-6, 0) // border (2) + padding (4)

	// Build the top border line manually: ╭─ title ─────── rightLabel ─╮
	topLeft := DimStyle.Render("╭─")
	topRight := DimStyle.Render("─╮")
	rightStr := ""
	if rightLabel != "" {
		rightStr = " " + DimStyle.Render(rightLabel) + " "
	}

	topLeftW := lipgloss.Width(topLeft)
	topRightW := lipgloss.Width(topRight)
	rightW := lipgloss.Width(rightStr)

	// Drop the right label entirely when it alone would overflow the box.
	if rightStr != "" && topLeftW+rightW+topRightW > w {
		rightStr = ""
		rightW = 0
	}

	// Truncate an overlong title (display-width aware) so the top border line
	// always renders exactly w columns wide.
	titleStr := ""
	if title != "" {
		availTitle := w - topLeftW - rightW - topRightW - 2 // 2 spaces around title
		if availTitle >= 1 {
			titleStr = " " + TitleStyle.Render(truncateWidth(title, availTitle)) + " "
		}
	}
	titleW := lipgloss.Width(titleStr)

	fillW := maxInt(w-topLeftW-titleW-rightW-topRightW, 0)

	topLine := topLeft + titleStr + DimStyle.Render(strings.Repeat("─", fillW)) + rightStr + topRight

	// Build body with side borders. Every body line is truncated to the inner
	// width first so an overlong line can never break through the right border.
	bodyLines := strings.Split(body, "\n")
	var middle strings.Builder
	for _, line := range bodyLines {
		line = truncateWidth(line, innerW)
		lineW := lipgloss.Width(line)
		pad := maxInt(innerW-lineW, 0)
		middle.WriteString(DimStyle.Render("│") + "  " + line + strings.Repeat(" ", pad) + "  " + DimStyle.Render("│") + "\n")
	}

	// Empty padding line top and bottom inside box
	emptyLine := DimStyle.Render("│") + strings.Repeat(" ", maxInt(innerW+4, 0)) + DimStyle.Render("│")

	// Bottom border
	bottomLine := DimStyle.Render("╰" + strings.Repeat("─", maxInt(w-2, 0)) + "╯")

	rendered := topLine + "\n" + emptyLine + "\n" + middle.String() + emptyLine + "\n" + bottomLine

	// Horizontal centering — reference the box width, not the rendered top
	// line, so an anomalous line can never shift the whole box.
	leftPad := 0
	if termWidth > w {
		leftPad = (termWidth - w) / 2
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

// sectionTitle renders a labeled section divider: "── 标题 ────────────"
func sectionTitle(title string, width int) string {
	// Match contentBox's body width: border (2) + padding (4). Overshooting
	// here made the rule overflow and get truncated with an ellipsis.
	inner := maxInt(boxWidth(width)-6, 0)
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

	var b strings.Builder
	for i := 0; i < level; i++ {
		b.WriteString(S.StarOn)
	}
	for i := level; i < 5; i++ {
		b.WriteString(S.StarOff)
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

// padDisplayWidth is a legacy alias for padRight (textutil.go).
func padDisplayWidth(s string, width int) string {
	return padRight(s, width)
}

// boxLines renders a bordered box at an exact outer width and returns
// its lines without any horizontal centering — the wide layout needs to place
// two boxes side by side, which contentBox (terminal-centered, single box)
// cannot do. Visual language matches contentBox: manual borders, title in the
// top border, right-aligned label, 2-space inner padding. The design's main
// box keeps a blank padding row above and below the body (padded=true); the
// inspector box is content-tight (padded=false).
func boxLines(title string, body []string, w int, rightLabel string, padded bool) []string {
	if w < minBoxWidth {
		w = minBoxWidth
	}
	innerW := w - 6

	topLeft := DimStyle.Render("╭─")
	topRight := DimStyle.Render("─╮")
	rightStr := ""
	if rightLabel != "" {
		rightStr = " " + DimStyle.Render(rightLabel) + " "
	}
	topLeftW := lipgloss.Width(topLeft)
	topRightW := lipgloss.Width(topRight)
	rightW := lipgloss.Width(rightStr)
	if rightStr != "" && topLeftW+rightW+topRightW > w {
		rightStr = ""
		rightW = 0
	}

	titleStr := ""
	if title != "" {
		availTitle := w - topLeftW - rightW - topRightW - 2
		if availTitle >= 1 {
			titleStr = " " + TitleStyle.Render(truncateWidth(title, availTitle)) + " "
		}
	}
	titleW := lipgloss.Width(titleStr)
	fillW := maxInt(w-topLeftW-titleW-rightW-topRightW, 0)

	lines := make([]string, 0, len(body)+4)
	lines = append(lines, topLeft+titleStr+DimStyle.Render(strings.Repeat("─", fillW))+rightStr+topRight)

	empty := DimStyle.Render("│") + strings.Repeat(" ", maxInt(innerW+4, 0)) + DimStyle.Render("│")
	if padded {
		lines = append(lines, empty)
	}
	for _, line := range body {
		line = truncateWidth(line, innerW)
		pad := maxInt(innerW-lipgloss.Width(line), 0)
		lines = append(lines, DimStyle.Render("│")+"  "+line+strings.Repeat(" ", pad)+"  "+DimStyle.Render("│"))
	}
	if padded {
		lines = append(lines, empty)
	}
	lines = append(lines, DimStyle.Render("╰"+strings.Repeat("─", maxInt(w-2, 0))+"╯"))
	return lines
}
