package tui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/sd3/linuxlab/internal/challenge"
	"github.com/sd3/linuxlab/internal/progress"
)

// overlay.go composites a floating panel over an already-rendered screen,
// following opencode's PlaceOverlay model: the background keeps its own
// rendering and the panel is spliced into a rectangle of it, so a popup never
// forces the screen underneath to change its layout.
//
// Slicing is display-width based and ANSI aware (ansi.Truncate / ansi.Cut), so
// cutting through styled or CJK text neither splits a wide glyph nor leaks a
// colour into the rest of the line.

// placeOverlay draws fg over bg with its top-left corner at (x, y), measured in
// display columns and rows. Lines of bg outside the panel's rectangle are
// returned untouched; the result always has bg's line count.
func placeOverlay(bg, fg string, x, y int) string {
	if fg == "" {
		return bg
	}
	bgLines := strings.Split(bg, "\n")
	fgLines := strings.Split(fg, "\n")

	fgW := 0
	for _, line := range fgLines {
		if w := lipgloss.Width(line); w > fgW {
			fgW = w
		}
	}
	if fgW == 0 {
		return bg
	}
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	for i, fgLine := range fgLines {
		row := y + i
		if row >= len(bgLines) {
			break
		}
		bgLines[row] = spliceLine(bgLines[row], padRight(fgLine, fgW), x, fgW)
	}
	return strings.Join(bgLines, "\n")
}

// spliceLine replaces the [x, x+width) column range of line with panel.
func spliceLine(line, panel string, x, width int) string {
	lineW := lipgloss.Width(line)

	left := ""
	if x > 0 {
		left = ansi.Truncate(line, x, "")
		// The background may be shorter than the panel's offset, or may end
		// with a wide glyph that cannot be split; pad out to the exact column.
		if gap := x - lipgloss.Width(left); gap > 0 {
			left += strings.Repeat(" ", gap)
		}
	}

	// The tail is taken by dropping columns from the left rather than by an end
	// index: a wide glyph straddling the boundary makes end-index slicing round
	// outward and the line comes back a column or two too long.
	right := ""
	if want := lineW - x - width; want > 0 {
		right = ansi.TruncateLeft(line, x+width, "")
		switch got := lipgloss.Width(right); {
		case got > want:
			// A split wide glyph was rounded in; drop it and keep the columns.
			right = ansi.TruncateLeft(right, got-want, "")
			if pad := want - lipgloss.Width(right); pad > 0 {
				right = strings.Repeat(" ", pad) + right
			}
		case got < want:
			right = strings.Repeat(" ", want-got) + right
		}
	}
	return left + panel + right
}

// --- floating panel animation -------------------------------------------------

// panelTickMsg advances an open/close animation by one frame.
type panelTickMsg struct{}

// panelFrame is how long one animation frame lasts; panelSteps frames make the
// whole reveal, so a panel opens in roughly 100ms.
const (
	panelFrame = 25 * time.Millisecond
	panelSteps = 4
)

func panelTickCmd() tea.Cmd {
	return tea.Tick(panelFrame, func(time.Time) tea.Msg { return panelTickMsg{} })
}

// panelAnim is the open/close state of a floating panel, shared by every
// screen that has one. step is the frame currently drawn and target the frame
// it is moving toward, so closing animates out instead of blinking away.
type panelAnim struct {
	step   int
	target int
}

// toggle flips the panel's direction, starting the animation when idle.
func (a panelAnim) toggle() (panelAnim, tea.Cmd) {
	running := a.step != a.target
	if a.target > 0 {
		a.target = 0
	} else {
		a.target = panelSteps
	}
	if running {
		// A tick is already in flight and will pick up the new target.
		return a, nil
	}
	return a, panelTickCmd()
}

// advance moves one frame toward the target and schedules the next.
func (a panelAnim) advance() (panelAnim, tea.Cmd) {
	switch {
	case a.step < a.target:
		a.step++
	case a.step > a.target:
		a.step--
	default:
		return a, nil
	}
	if a.step == a.target {
		return a, nil
	}
	return a, panelTickCmd()
}

// opening reports whether the panel is open or on its way there, which is what
// Esc should close before it leaves the screen.
func (a panelAnim) opening() bool { return a.target > 0 }

// visible reports whether anything should be drawn at all.
func (a panelAnim) visible() bool { return a.step > 0 }

// reveal returns how many of total body rows this frame shows, so the panel
// grows with its borders intact rather than appearing half-drawn.
func (a panelAnim) reveal(total int) int {
	return clampInt(total*a.step/panelSteps, 0, total)
}

// challengePreviewLines builds the body of a challenge preview panel: title,
// completion state, difficulty, subcategory, tags and the opening of the
// description. store may be nil on screens that do not track progress, in
// which case the state row is omitted rather than guessed.
func challengePreviewLines(ch *challenge.Challenge, store *progress.Store, w, descLimit int) []string {
	if ch == nil {
		return []string{S.Dim.Render("没有可预览的挑战")}
	}

	row := func(label, value string) string {
		return truncateWidth(S.Dim.Render(padRight(label, 4))+" "+value, w)
	}
	lines := []string{truncateWidth(S.Text.Render(ch.Title), w), ""}

	if store != nil {
		status, statusStyle := "未完成", S.Dim
		if entry, ok := store.Data.Challenges[ch.ID]; ok && entry != nil {
			switch entry.Status {
			case "passed":
				status, statusStyle = "已通过", S.Success
			case "failed":
				status, statusStyle = "未通过", S.Error
			}
		}
		lines = append(lines, row("状态", statusStyle.Render(status)))
	}
	lines = append(lines, row("难度", DifficultyStars(ch.Difficulty)))
	if ch.Subcategory != "" {
		lines = append(lines, row("子类", S.Text.Render(ch.Subcategory)))
	}
	if len(ch.Tags) > 0 {
		lines = append(lines, row("标签", S.Dim.Render(strings.Join(ch.Tags, " · "))))
	}
	if desc := strings.TrimSpace(ch.Description); desc != "" {
		lines = append(lines, "", S.Dim.Render("── 描述 "+strings.Repeat("─", maxInt(w-8, 0))))
		wrapped := wrapWidth(desc, w)
		if descLimit > 0 && len(wrapped) > descLimit {
			wrapped = append(wrapped[:descLimit:descLimit], S.Dim.Render("…"))
		}
		lines = append(lines, wrapped...)
	}
	return lines
}
