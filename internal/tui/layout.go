package tui

// LayoutMode classifies the terminal size into the breakpoints fixed by the
// responsive-TUI PRD (docs/superpowers/specs).
type LayoutMode int

const (
	// ModeUnsupported: terminal smaller than 60x18 — render only a centered
	// "terminal too small" notice and skip all other layout.
	ModeUnsupported LayoutMode = iota
	// ModeCompact: w < 80 or h < 24.
	ModeCompact
	// ModeStandard: 80 <= w < 120 (and h >= 24), or wide-but-short terminals.
	ModeStandard
	// ModeWide: w >= 120 and h >= 30 — the roomiest tier; screens use it for
	// extra breathing room, never for a second column (side panels float).
	ModeWide
)

// Minimum supported terminal size.
const (
	minTermWidth  = 60
	minTermHeight = 18
)

// LayoutSpec is the pure result of computeLayout: the app shell's slot
// budget for a given terminal size. It replaces scattered height-12 /
// height-8 magic numbers as the single source of truth for the frame layout.
// Screens must still measure their own rendered headers with lipgloss.Height
// when subdividing MainH.
type LayoutSpec struct {
	Mode          LayoutMode
	Width, Height int

	HeaderH, FooterH int // shell header/footer heights (1 each when supported)
	MainW, MainH     int // budget for the active screen's content
}

// computeLayout derives the layout spec from the terminal size. Pure function:
// no IO, no state.
func computeLayout(w, h int) LayoutSpec {
	spec := LayoutSpec{Width: w, Height: h}

	switch {
	case w < minTermWidth || h < minTermHeight:
		spec.Mode = ModeUnsupported
		return spec
	case w < 80 || h < 24:
		spec.Mode = ModeCompact
	case w >= 120 && h >= 30:
		spec.Mode = ModeWide
	default:
		spec.Mode = ModeStandard
	}

	spec.HeaderH = 1
	spec.FooterH = 1
	spec.MainW = w
	spec.MainH = h - spec.HeaderH - spec.FooterH
	return spec
}
