package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

// splash.go renders the cover shown when the app opens, in the spirit of
// Claude Code's welcome banner: a wordmark, one line of positioning, a boxed
// snapshot of the current session, and a hint for the first keystroke.
//
// Sizing follows gemini-cli's approach: several wordmark variants are kept and
// the widest one that actually fits is chosen from its measured width, rather
// than from hardcoded breakpoints.

// Gradient endpoints for the wordmark, resolved once at startup: AdaptiveColor
// cannot be interpolated, so the concrete variant for the detected background
// is picked here (iron rule: no terminal IO inside View).
var splashGradientFrom, splashGradientTo = splashGradientColors()

func splashGradientColors() (colorful.Color, colorful.Color) {
	fromHex, toHex := currentTheme.Primary.Dark, currentTheme.Secondary.Dark
	if !lipgloss.HasDarkBackground() {
		fromHex, toHex = currentTheme.Primary.Light, currentTheme.Secondary.Light
	}
	from, err := colorful.Hex(fromHex)
	if err != nil {
		from = colorful.Color{R: 0.3, G: 0.79, B: 0.76}
	}
	to, err := colorful.Hex(toHex)
	if err != nil {
		to = from
	}
	return from, to
}

// splashWordmarkLarge is the 5-row block rendering of "LINUXLAB" (47 columns).
var splashWordmarkLarge = []string{
	"█     █████ █   █ █   █ █   █ █      ███  ████ ",
	"█       █   ██  █ █   █  █ █  █     █   █ █   █",
	"█       █   █ █ █ █   █   █   █     █████ ████ ",
	"█       █   █  ██ █   █  █ █  █     █   █ █   █",
	"█████ █████ █   █  ███  █   █ █████ █   █ ████ ",
}

// splashWordmarkSmall is the single-line fallback for narrow terminals.
var splashWordmarkSmall = []string{"L I N U X L A B"}

// artWidth measures the display width of the widest line of ASCII art.
func artWidth(art []string) int {
	w := 0
	for _, line := range art {
		if lw := displayWidth(line); lw > w {
			w = lw
		}
	}
	return w
}

// pickWordmark returns the largest wordmark that fits the available width.
func pickWordmark(avail int) []string {
	if artWidth(splashWordmarkLarge) <= avail {
		return splashWordmarkLarge
	}
	return splashWordmarkSmall
}

// gradientArt colors each column of the art by interpolating from→to across
// the art's full width, so the wordmark reads as one continuous sweep. Blending
// happens in Luv space, which keeps the midpoint from going muddy.
func gradientArt(art []string, from, to colorful.Color) []string {
	width := artWidth(art)
	if width <= 0 {
		return art
	}

	// Pre-build one style per column: the art is far taller than it is deep in
	// distinct colors, so this avoids rebuilding styles per cell.
	styles := make([]lipgloss.Style, width)
	for i := 0; i < width; i++ {
		ratio := 0.0
		if width > 1 {
			ratio = float64(i) / float64(width-1)
		}
		styles[i] = lipgloss.NewStyle().Foreground(lipgloss.Color(from.BlendLuv(to, ratio).Hex()))
	}

	out := make([]string, 0, len(art))
	for _, line := range art {
		var b strings.Builder
		col := 0
		for _, r := range line {
			if r == ' ' {
				b.WriteRune(' ')
			} else {
				b.WriteString(styles[minInt(col, width-1)].Render(string(r)))
			}
			col += displayWidth(string(r))
		}
		out = append(out, b.String())
	}
	return out
}

// splashInfo is the session snapshot shown inside the cover's box.
type splashInfo struct {
	stats MenuStats
	// dockerProbed distinguishes "probe still running" from "Docker is down",
	// so the cover never accuses a healthy daemon of being absent.
	dockerProbed bool
	dockerOK     bool
}

// splashTips are the getting-started lines. They are dropped first when the
// terminal is too short for the full cover.
var splashTips = []string{
	"从「开始练习」按模块推进，完成的题目会自动记入进度",
	"卡住时按 4 查命令速查，不会打断当前挑战",
}

// renderSplash builds the cover for the given terminal size. It degrades in
// two steps — tips first, then the block wordmark — so the cover still reads
// correctly down to the minimum supported terminal.
func renderSplash(width, height int, info splashInfo) string {
	// The cover is drawn without the app shell, so it owns the full height.
	boxW := maxInt(boxWidth(width), minBoxWidth)
	innerW := maxInt(boxW-6, 0)

	wordmark := pickWordmark(width - 4)
	if height < 20 {
		wordmark = splashWordmarkSmall
	}
	if len(wordmark) > 1 {
		wordmark = gradientArt(wordmark, splashGradientFrom, splashGradientTo)
	} else {
		wordmark = []string{S.Title.Render(wordmark[0])}
	}

	stats := info.stats
	greeting := "✳ 欢迎回来"
	if stats.PassedChallenges == 0 && stats.Last == nil {
		greeting = "✳ 欢迎使用 LinuxLab"
	}

	var body strings.Builder
	body.WriteString(S.Selected.Render(greeting))
	body.WriteString("\n\n")

	// Overall progress, as a bar rather than a bare count.
	body.WriteString(truncateWidth(splashProgressLine(stats, innerW), innerW))
	body.WriteString("\n")

	// Where the session left off, and what to pick up next.
	for _, line := range splashSessionLines(stats, innerW) {
		body.WriteString(truncateWidth(line, innerW))
		body.WriteString("\n")
	}

	docker := S.Dim.Render("Docker 检测中…")
	if info.dockerProbed {
		docker = S.Warning.Render("Docker 未运行 · 将以本地降级模式执行")
		if info.dockerOK {
			docker = S.Success.Render("Docker 就绪")
		}
	}
	body.WriteString(truncateWidth(S.Dim.Render(padRight("运行环境", splashLabelW))+docker, innerW))

	box := contentBox("", body.String(), width, height, "")

	blocks := []string{
		centerBlock(strings.Join(wordmark, "\n"), width),
		"",
		centerBlock(S.Dim.Render("真实 Shell / Vim / 运维场景训练"), width),
		"",
		box,
	}

	// Tips fit only when there is vertical room to spare.
	used := 0
	for _, b := range blocks {
		used += lipgloss.Height(b)
	}
	if height-used >= len(splashTips)+4 {
		tips := make([]string, 0, len(splashTips))
		for _, tip := range splashTips {
			tips = append(tips, truncateWidth(S.Dim.Render("▸ ")+S.Dim.Render(tip), width-4))
		}
		blocks = append(blocks, "", centerBlock(strings.Join(tips, "\n"), width))
	}
	blocks = append(blocks, "", centerBlock(S.Help.Render("按任意键开始")+"  "+S.Dim.Render("· q 退出"), width))

	return strings.Join(blocks, "\n")
}

// splashLabelW is the width of the label column inside the cover's box. It is
// wider than the longest label ("建议下一题" = 10 columns) so every value lines
// up on the same left edge with a visible gap.
const splashLabelW = 12

// splashProgressLine renders "总体进度  ███░░░░  12/278 · 4%".
func splashProgressLine(stats MenuStats, innerW int) string {
	ratio := 0.0
	if stats.TotalChallenges > 0 {
		ratio = float64(stats.PassedChallenges) / float64(stats.TotalChallenges)
	}
	suffix := fmt.Sprintf("  %d/%d · %d%%", stats.PassedChallenges, stats.TotalChallenges, int(ratio*100))
	barW := clampInt(innerW-splashLabelW-displayWidth(suffix), 8, 28)
	return S.Dim.Render(padRight("总体进度", splashLabelW)) + ProgressBar(ratio, barW) + S.Dim.Render(suffix)
}

// splashSessionLines describe where the user left off and what to do next.
// Both lines are omitted when there is nothing true to say about them.
func splashSessionLines(stats MenuStats, innerW int) []string {
	var lines []string

	if stats.Last != nil {
		icon, state := S.IconTodo, S.Dim.Render("进行中")
		if stats.LastEntry != nil {
			switch stats.LastEntry.Status {
			case "passed":
				icon, state = S.IconPass, S.Success.Render("已通过")
			case "failed":
				icon, state = S.IconFail, S.Error.Render("未通过")
			}
		}
		meta := " · " + state
		if stats.LastEntry != nil {
			if stats.LastEntry.Attempts > 1 {
				meta += S.Dim.Render(fmt.Sprintf(" · %d 次尝试", stats.LastEntry.Attempts))
			}
			if day := stats.LastEntry.LastAttempt; day != "" {
				meta += S.Dim.Render(" · " + day)
			}
		}
		title := truncateWidth(stats.Last.Title,
			maxInt(innerW-splashLabelW-displayWidth(meta)-2, 8))
		lines = append(lines, S.Dim.Render(padRight("上次练到", splashLabelW))+icon+" "+S.Text.Render(title)+meta)
	}

	if stats.Next != nil {
		stars := DifficultyStars(stats.Next.Difficulty)
		title := truncateWidth(stats.Next.Title,
			maxInt(innerW-splashLabelW-displayWidth(stars)-3, 8))
		lines = append(lines, S.Dim.Render(padRight("建议下一题", splashLabelW))+S.Text.Render(title)+"  "+stars)
	}

	return lines
}

// centerBlock horizontally centers every line of a block within width, using
// the block's widest line as the reference so its internal alignment holds.
func centerBlock(block string, width int) string {
	lines := strings.Split(block, "\n")
	blockW := 0
	for _, line := range lines {
		if w := displayWidth(line); w > blockW {
			blockW = w
		}
	}
	pad := maxInt((width-blockW)/2, 0)
	if pad == 0 {
		return block
	}
	prefix := strings.Repeat(" ", pad)
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}
