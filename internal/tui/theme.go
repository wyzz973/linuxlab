package tui

import "github.com/charmbracelet/lipgloss"

// Theme is the semantic color token set for the whole TUI. All colors are
// adaptive: the Dark values are the project's canonical palette (matching
// docs/tui-preview.html) and the Light values are readable equivalents for
// light-background terminals.
type Theme struct {
	// Brand
	Primary       lipgloss.AdaptiveColor
	PrimarySubtle lipgloss.AdaptiveColor
	Secondary     lipgloss.AdaptiveColor

	// Foreground, three tiers
	FgBase   lipgloss.AdaptiveColor
	FgSubtle lipgloss.AdaptiveColor // secondary text (was ColorDim)
	FgMuted  lipgloss.AdaptiveColor // weaker than FgSubtle, e.g. footer key hints

	// Background
	BgBase  lipgloss.AdaptiveColor
	BgPanel lipgloss.AdaptiveColor
	BgFocus lipgloss.AdaptiveColor // unified selected-row background highlight

	// Border
	Border      lipgloss.AdaptiveColor
	BorderFocus lipgloss.AdaptiveColor

	// Status
	Success lipgloss.AdaptiveColor
	Warning lipgloss.AdaptiveColor
	Error   lipgloss.AdaptiveColor
}

// DefaultTheme returns the canonical linuxlab theme.
func DefaultTheme() Theme {
	return Theme{
		Primary:       lipgloss.AdaptiveColor{Light: "#0f766e", Dark: "#4cc9c1"},
		PrimarySubtle: lipgloss.AdaptiveColor{Light: "#5eaea7", Dark: "#2a6f6a"},
		Secondary:     lipgloss.AdaptiveColor{Light: "#3056b3", Dark: "#8db7ff"},

		FgBase:   lipgloss.AdaptiveColor{Light: "#1f2933", Dark: "#d7e1e8"},
		FgSubtle: lipgloss.AdaptiveColor{Light: "#5b6b76", Dark: "#81919c"},
		FgMuted:  lipgloss.AdaptiveColor{Light: "#8a99a3", Dark: "#5f6f79"},

		BgBase:  lipgloss.AdaptiveColor{Light: "#f7f9fa", Dark: "#101418"},
		BgPanel: lipgloss.AdaptiveColor{Light: "#eef1f3", Dark: "#151b21"},
		BgFocus: lipgloss.AdaptiveColor{Light: "#d8ebe9", Dark: "#22343a"},

		Border:      lipgloss.AdaptiveColor{Light: "#c3ced4", Dark: "#33424c"},
		BorderFocus: lipgloss.AdaptiveColor{Light: "#0f766e", Dark: "#4cc9c1"},

		Success: lipgloss.AdaptiveColor{Light: "#15803d", Dark: "#7ecf8f"},
		Warning: lipgloss.AdaptiveColor{Light: "#b45309", Dark: "#f3bc5f"},
		Error:   lipgloss.AdaptiveColor{Light: "#dc2626", Dark: "#ff7a7a"},
	}
}

// Icon constants. These are the only non-ASCII symbols allowed in
// alignment-critical rows; their width must always be measured through
// displayWidth / lipgloss.Width, never assumed.
const (
	IconPass    = "✓"
	IconFail    = "✗"
	IconTodo    = "○"
	IconCursor  = "›"
	StarOn      = "★"
	StarOff     = "☆"
	BarFull     = "█"
	BarEmpty    = "░"
	ArrowFold   = "▶"
	ArrowUnfold = "▼"
	ArrowUp     = "▲"
	ArrowDown   = "▼"
)

// Styles is the centralized, pre-built style set. All screens must take
// styles from here (or from the legacy aliases in styles.go that point here);
// calling lipgloss.NewStyle() inside View is forbidden.
type Styles struct {
	// Title / text
	Title    lipgloss.Style
	Text     lipgloss.Style
	Subtle   lipgloss.Style
	Dim      lipgloss.Style
	Muted    lipgloss.Style
	Selected lipgloss.Style

	// List
	SelectedRow lipgloss.Style // whole-row background highlight

	// Status
	Success lipgloss.Style
	Warning lipgloss.Style
	Error   lipgloss.Style

	// Box / borders
	Border      lipgloss.Style
	BorderFocus lipgloss.Style

	// Help / footer / header
	Help       lipgloss.Style
	HelpKey    lipgloss.Style
	HelpDesc   lipgloss.Style
	Notice     lipgloss.Style
	Breadcrumb lipgloss.Style
	Crumb      lipgloss.Style

	// Badge / meta
	Key   lipgloss.Style
	Badge lipgloss.Style
	Meta  lipgloss.Style

	// Pre-rendered icons
	IconPass   string
	IconFail   string
	IconTodo   string
	IconCursor string
	BarFull    string
	BarEmpty   string
	StarOn     string
	StarOff    string
}

// NewStyles builds the style set from a theme. Called once at startup.
func NewStyles(t Theme) Styles {
	return Styles{
		Title:    lipgloss.NewStyle().Foreground(t.Primary).Bold(true),
		Text:     lipgloss.NewStyle().Foreground(t.FgBase),
		Subtle:   lipgloss.NewStyle().Foreground(t.Border),
		Dim:      lipgloss.NewStyle().Foreground(t.FgSubtle),
		Muted:    lipgloss.NewStyle().Foreground(t.FgMuted),
		Selected: lipgloss.NewStyle().Foreground(t.Primary).Bold(true),

		SelectedRow: lipgloss.NewStyle().Background(t.BgFocus),

		Success: lipgloss.NewStyle().Foreground(t.Success),
		Warning: lipgloss.NewStyle().Foreground(t.Warning),
		Error:   lipgloss.NewStyle().Foreground(t.Error),

		Border:      lipgloss.NewStyle().Foreground(t.Border),
		BorderFocus: lipgloss.NewStyle().Foreground(t.BorderFocus),

		Help:       lipgloss.NewStyle().Foreground(t.FgSubtle),
		HelpKey:    lipgloss.NewStyle().Foreground(t.FgSubtle),
		HelpDesc:   lipgloss.NewStyle().Foreground(t.FgMuted),
		Notice:     lipgloss.NewStyle().Foreground(t.Warning),
		Breadcrumb: lipgloss.NewStyle().Foreground(t.FgSubtle),
		Crumb:      lipgloss.NewStyle().Foreground(t.FgBase),

		Key:   lipgloss.NewStyle().Foreground(t.BgBase).Background(t.Primary).Bold(true).Padding(0, 1),
		Badge: lipgloss.NewStyle().Foreground(t.Primary).Bold(true),
		Meta:  lipgloss.NewStyle().Foreground(t.FgSubtle),

		IconPass:   lipgloss.NewStyle().Foreground(t.Success).Render(IconPass),
		IconFail:   lipgloss.NewStyle().Foreground(t.Error).Render(IconFail),
		IconTodo:   lipgloss.NewStyle().Foreground(t.FgSubtle).Render(IconTodo),
		IconCursor: lipgloss.NewStyle().Foreground(t.Primary).Bold(true).Render(IconCursor),
		BarFull:    lipgloss.NewStyle().Foreground(t.Success).Render(BarFull),
		BarEmpty:   lipgloss.NewStyle().Foreground(t.Border).Render(BarEmpty),
		StarOn:     lipgloss.NewStyle().Foreground(t.Warning).Render(StarOn),
		StarOff:    lipgloss.NewStyle().Foreground(t.FgSubtle).Render(StarOff),
	}
}

// currentTheme and S are the package-wide theme and style set, pre-built at
// startup (iron rule: no lipgloss.NewStyle() inside View).
var (
	currentTheme = DefaultTheme()
	S            = NewStyles(currentTheme)
)
