package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sd3/linuxlab/internal/challenge"
	"github.com/sd3/linuxlab/internal/progress"
)

// MenuChoiceMsg is sent when the user selects a menu item.
type MenuChoiceMsg struct {
	Choice string
}

// menuStatusKind selects the color of a menu item's right-hand status word.
type menuStatusKind int

const (
	menuStatusNone menuStatusKind = iota
	menuStatusSuccess
	menuStatusWarning
	menuStatusMeta
)

type menuItem struct {
	label  string
	desc   string
	key    string
	status string         // right-hand status word, "" for none
	kind   menuStatusKind // color of the status word
}

// MenuStats carries the live counters the menu turns into right-hand status
// words. The zero value renders the static fallback words, which is what
// NewMenuModel (no data available) relies on.
type MenuStats struct {
	TotalChallenges  int
	TotalModules     int
	PassedChallenges int
	WeakCount        int // challenges recommended as weak spots
	ReferenceCount   int // command reference entries

	// Session pointers used by the opening cover. Computed alongside the
	// counters so the cover never has to walk the progress set itself.
	Last      *challenge.Challenge // most recently attempted
	LastEntry *progress.ChallengeEntry
	Next      *challenge.Challenge // top weak-spot recommendation
}

// practiceStatus reports progress through the whole challenge set.
func (s MenuStats) practiceStatus() (string, menuStatusKind) {
	switch {
	case s.TotalChallenges == 0:
		return "继续训练", menuStatusSuccess
	case s.PassedChallenges == 0:
		return "未开始", menuStatusMeta
	case s.PassedChallenges >= s.TotalChallenges:
		return "已完成", menuStatusSuccess
	default:
		return fmt.Sprintf("%d/%d", s.PassedChallenges, s.TotalChallenges), menuStatusSuccess
	}
}

func (s MenuStats) recommendStatus() (string, menuStatusKind) {
	if s.WeakCount == 0 {
		return "查漏补缺", menuStatusWarning
	}
	return fmt.Sprintf("%d 项薄弱", s.WeakCount), menuStatusWarning
}

func (s MenuStats) referenceStatus() (string, menuStatusKind) {
	if s.ReferenceCount == 0 {
		return "速查手册", menuStatusMeta
	}
	return fmt.Sprintf("%d 条", s.ReferenceCount), menuStatusMeta
}

// menuItems returns the canonical entry list: 4 items in fixed order, number
// keys 1-4 jump directly (order and count locked by menu_test). Status words
// are derived from stats; the zero value keeps the static fallback wording.
func menuItems(stats MenuStats) []menuItem {
	practice, practiceKind := stats.practiceStatus()
	recommend, recommendKind := stats.recommendStatus()
	ref, refKind := stats.referenceStatus()
	return []menuItem{
		{label: "开始练习", desc: "选择模块开始学习", key: "practice", status: practice, kind: practiceKind},
		{label: "能力图谱", desc: "查看各项技能掌握程度", key: "skillmap"},
		{label: "薄弱推荐", desc: "针对薄弱环节推荐练习", key: "recommend", status: recommend, kind: recommendKind},
		{label: "命令速查", desc: "快速查阅常用命令用法", key: "reference", status: ref, kind: refKind},
	}
}

// listRowStyles bundles the per-segment styles of one selectable list row
// (shared by the menu and modules screens). The focus variant carries the
// BgFocus background on every segment so the selected row highlights as one
// full-width block (spec §4.1). Pre-built at startup — no lipgloss.NewStyle()
// inside View.
type listRowStyles struct {
	cursor  lipgloss.Style
	badge   lipgloss.Style
	label   lipgloss.Style
	desc    lipgloss.Style
	success lipgloss.Style
	warning lipgloss.Style
	meta    lipgloss.Style
	fill    lipgloss.Style

	// Pre-rendered progress-bar glyphs (same palette as ProgressBar, plus the
	// row background in the focus variant) for in-row bars (modules screen).
	barFull  string
	barEmpty string
}

func newListRowStyles(focused bool) listRowStyles {
	t := currentTheme
	base := lipgloss.NewStyle()
	if focused {
		base = base.Background(t.BgFocus)
	}
	return listRowStyles{
		cursor:  base.Foreground(t.Primary).Bold(true),
		badge:   base.Foreground(t.Primary),
		label:   base.Foreground(t.FgBase).Bold(focused),
		desc:    base.Foreground(t.FgSubtle),
		success: base.Foreground(t.Success),
		warning: base.Foreground(t.Warning),
		meta:    base.Foreground(t.FgSubtle),
		fill:    base,

		barFull:  base.Foreground(t.Success).Render(BarFull),
		barEmpty: base.Foreground(t.Border).Render(BarEmpty),
	}
}

var (
	listRowPlain = newListRowStyles(false)
	listRowFocus = newListRowStyles(true)
)

// bar renders a progress bar of the given width from the pre-rendered glyphs,
// so a focused row's bar keeps the row background behind it.
func (s listRowStyles) bar(ratio float64, width int) string {
	if width < 1 {
		return ""
	}
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	filled := int(ratio * float64(width))
	if filled > width {
		filled = width
	}
	var b strings.Builder
	for i := 0; i < filled; i++ {
		b.WriteString(s.barFull)
	}
	for i := filled; i < width; i++ {
		b.WriteString(s.barEmpty)
	}
	return b.String()
}

func (s listRowStyles) statusStyle(kind menuStatusKind) lipgloss.Style {
	switch kind {
	case menuStatusSuccess:
		return s.success
	case menuStatusWarning:
		return s.warning
	default:
		return s.meta
	}
}

// MenuModel is the main menu screen.
type MenuModel struct {
	items  []menuItem
	cursor int
	width  int
	height int
	stats  MenuStats
}

// NewMenuModel creates a new main menu without live counters.
func NewMenuModel() tea.Model {
	return NewMenuModelWithData(MenuStats{})
}

// NewMenuModelWithStats creates a menu model with challenge/module counts.
func NewMenuModelWithStats(totalChallenges, totalModules int) tea.Model {
	return NewMenuModelWithData(MenuStats{TotalChallenges: totalChallenges, TotalModules: totalModules})
}

// NewMenuModelWithData creates a menu model from the full counter set.
func NewMenuModelWithData(stats MenuStats) tea.Model {
	return MenuModel{items: menuItems(stats), cursor: 0, stats: stats}
}

// withStats returns a copy carrying refreshed counters, preserving the cursor
// and cached terminal size so returning to the menu keeps its position.
func (m MenuModel) withStats(stats MenuStats) MenuModel {
	m.stats = stats
	m.items = menuItems(stats)
	if m.cursor >= len(m.items) {
		m.cursor = len(m.items) - 1
	}
	return m
}

func (m MenuModel) Init() tea.Cmd { return nil }

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch {
		case runeMatches(msg, "q"):
			return m, tea.Quit
		case msg.Type == tea.KeyRunes:
			if choice, ok := m.choiceByNumber(string(msg.Runes)); ok {
				return m, func() tea.Msg { return MenuChoiceMsg{Choice: choice} }
			}
			cursor, _, moved := moveListCursor(m.cursor, 0, len(m.items), len(m.items), msg, true)
			if moved {
				m.cursor = cursor
			}
		case isSelectKey(msg):
			choice := m.items[m.cursor].key
			return m, func() tea.Msg { return MenuChoiceMsg{Choice: choice} }
		default:
			cursor, _, moved := moveListCursor(m.cursor, 0, len(m.items), len(m.items), msg, true)
			if moved {
				m.cursor = cursor
			}
		}
	}
	return m, nil
}

func (m MenuModel) choiceByNumber(value string) (string, bool) {
	if len(value) != 1 {
		return "", false
	}
	index := int(value[0] - '1')
	if index < 0 || index >= len(m.items) {
		return "", false
	}
	return m.items[index].key, true
}

// menuColumnWidths derives the label/description column budgets from the
// widest entry (display width, CJK aware), so the columns stay aligned.
func menuColumnWidths(items []menuItem) (labelW, descW int) {
	for _, it := range items {
		labelW = maxInt(labelW, displayWidth(it.label))
		descW = maxInt(descW, displayWidth(it.desc))
	}
	return labelW + 2, descW + 4
}

// renderItem renders one menu row:
//
//	› [1] 开始练习  选择模块开始学习        继续训练
//
// The selected row is padded to the full inner width so the BgFocus
// background reads as a whole-row highlight.
func (m MenuModel) renderItem(i, innerW, labelW, descW int) string {
	item := m.items[i]
	selected := i == m.cursor
	st := listRowPlain
	if selected {
		st = listRowFocus
	}

	var b strings.Builder
	if selected {
		b.WriteString(st.cursor.Render(IconCursor + " "))
	} else {
		b.WriteString(st.fill.Render("  "))
	}
	b.WriteString(st.badge.Render(fmt.Sprintf("[%d] ", i+1)))
	b.WriteString(st.label.Render(padRight(item.label, labelW)))
	b.WriteString(st.desc.Render(padRight(item.desc, descW)))
	if item.status != "" {
		b.WriteString(st.statusStyle(item.kind).Render(item.status))
	}

	row := b.String()
	if selected {
		if gap := innerW - displayWidth(row); gap > 0 {
			row += st.fill.Render(strings.Repeat(" ", gap))
		}
	}
	return row
}

func (m MenuModel) View() string {
	// Mirror contentBox geometry: border (2) + padding (4).
	innerW := maxInt(boxWidth(m.width)-6, 0)
	labelW, descW := menuColumnWidths(m.items)

	lines := make([]string, 0, 2*len(m.items)+4)
	lines = append(lines, DimStyle.Render("真实 Shell / Vim / 运维场景训练"), "")
	for i := range m.items {
		if i > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, m.renderItem(i, innerW, labelW, descW))
	}
	if m.stats.TotalChallenges > 0 {
		lines = append(lines, "", DimStyle.Render(fmt.Sprintf("%d 道挑战 · %d 个模块", m.stats.TotalChallenges, m.stats.TotalModules)))
	}

	return contentBox("LinuxLab 训练控制台", strings.Join(lines, "\n"), m.width, m.height, "")
}
