package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sd3/linuxlab/internal/progress"
)

// Bar widths for the three progress-bar tiers (design bundle skillmap.html).
const (
	skillmapOverallBarW = 26
	skillmapCatBarW     = 18
	skillmapSubBarW     = 12
	skillmapLabelW      = 16
	skillmapSubLabelW   = 14
)

// Pre-built private styles for the selected category row (BgFocus highlight
// on the arrow+label segment; bars keep their own colors).
var (
	skillmapArrowSel = lipgloss.NewStyle().Foreground(currentTheme.Primary).Bold(true).Background(currentTheme.BgFocus)
	skillmapLabelSel = lipgloss.NewStyle().Foreground(currentTheme.FgBase).Bold(true).Background(currentTheme.BgFocus)
)

// skillmapCatRow holds the cursor-independent render fragments of one
// category row: the padded raw label (styled per-frame, depending on
// selection) plus the pre-rendered bar/percent strings and subcategory lines.
type skillmapCatRow struct {
	label string   // raw label, padded to skillmapLabelW
	bar   string   // rendered progress bar
	pct   string   // rendered "  10%" (dim)
	subs  []string // fully rendered subcategory lines
}

// skillmapCache memoizes the heavy parts of the skill-map render, keyed by
// (width, progress version). Cursor movement and fold toggles reassemble
// cached fragments without recomputing bars or widths; the cache only
// invalidates when the terminal width or the progress data changes
// (spec §3, opencode PartCache simplified). The cache pointer is shared by
// the value copies bubbletea creates, which is safe: it is a pure memo and
// the event loop is single-goroutine.
type skillmapCache struct {
	sm      *progress.SkillMap
	version int

	rendered bool
	width    int
	overall  string
	cats     []skillmapCatRow

	builds int // number of cache rebuilds (observed by tests)
}

// skillmapProgressVersion derives the progress data version from the store:
// the total number of recorded attempts. RecordAttempt always increments an
// attempt counter, so any progress change bumps the version. Summing over the
// map is order-independent, keeping renders frame-stable.
func skillmapProgressVersion(store *progress.Store) int {
	v := 0
	for _, entry := range store.Data.Challenges {
		if entry != nil {
			v += entry.Attempts
		}
	}
	return v
}

// SkillMapModel is the skill map view with expandable categories.
type SkillMapModel struct {
	store    *progress.Store
	cursor   int
	expanded map[int]bool
	width    int
	height   int
	cache    *skillmapCache
}

// NewSkillMapModel creates a new skill map view.
func NewSkillMapModel(store *progress.Store) tea.Model {
	return SkillMapModel{
		store:    store,
		cursor:   0,
		expanded: make(map[int]bool),
		cache: &skillmapCache{
			sm:      progress.BuildSkillMap(store),
			version: skillmapProgressVersion(store),
		},
	}
}

func (m SkillMapModel) Init() tea.Cmd { return nil }

// skillMap returns the current skill-map snapshot, rebuilding it when the
// progress data changed since it was built (pure in-memory work, no IO).
func (m SkillMapModel) skillMap() *progress.SkillMap {
	if v := skillmapProgressVersion(m.store); v != m.cache.version {
		m.cache.sm = progress.BuildSkillMap(m.store)
		m.cache.version = v
		m.cache.rendered = false
	}
	return m.cache.sm
}

func (m SkillMapModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		sm := m.skillMap()
		switch {
		case isBackKey(msg):
			return m, func() tea.Msg { return GoBackMsg{} }
		case msg.Type == tea.KeyEnter || msg.Type == tea.KeyRight || msg.Type == tea.KeySpace:
			if len(sm.Categories) > 0 {
				// Copy map to avoid shared mutation
				newExpanded := make(map[int]bool)
				for k, v := range m.expanded {
					newExpanded[k] = v
				}
				newExpanded[m.cursor] = !newExpanded[m.cursor]
				m.expanded = newExpanded
			}
		default:
			cursor, _, moved := moveListCursor(m.cursor, 0, len(sm.Categories), len(sm.Categories), msg, true)
			if moved {
				m.cursor = cursor
			}
		}
	}
	return m, nil
}

// ensureRendered fills the cache's render fragments for the current
// (width, progress version) key.
func (m SkillMapModel) ensureRendered(sm *progress.SkillMap) {
	c := m.cache
	if c.rendered && c.width == m.width {
		return
	}
	c.builds++
	c.rendered = true
	c.width = m.width

	c.overall = ""
	if sm.TotalCount > 0 {
		c.overall = ProgressSummary("总体掌握", sm.TotalPassed, sm.TotalCount, skillmapOverallBarW)
	}

	c.cats = c.cats[:0]
	for _, cat := range sm.Categories {
		row := skillmapCatRow{
			label: padRight(CategoryLabel(cat.Name), skillmapLabelW),
			bar:   ProgressBar(cat.Score, skillmapCatBarW),
			pct:   S.Dim.Render(fmt.Sprintf("%3.0f%%", cat.Score*100)),
		}
		for _, sub := range cat.Subcategories {
			row.subs = append(row.subs, "    "+
				ProgressSummary(S.Dim.Render(padRight(sub.Name, skillmapSubLabelW)), sub.Passed, sub.Total, skillmapSubBarW))
		}
		c.cats = append(c.cats, row)
	}
}

func (m SkillMapModel) View() string {
	sm := m.skillMap()

	if sm.TotalCount == 0 {
		var body strings.Builder
		body.WriteString(EmptyState("还没有练习记录", "完成一个挑战后，这里会展示各分类的掌握度。", "Esc 返回"))
		body.WriteString("\n")
		return contentBox("能力图谱", body.String(), m.width, m.height, "")
	}

	m.ensureRendered(sm)
	c := m.cache

	var body strings.Builder
	body.WriteString(c.overall)
	body.WriteString("\n\n")

	for i, row := range c.cats {
		arrow := ArrowFold
		if m.expanded[i] {
			arrow = ArrowUnfold
		}

		if i == m.cursor {
			body.WriteString(skillmapArrowSel.Render(arrow + " "))
			body.WriteString(skillmapLabelSel.Render(row.label))
		} else {
			body.WriteString(S.Dim.Render(arrow + " "))
			body.WriteString(S.Text.Render(row.label))
		}
		body.WriteString("  ")
		body.WriteString(row.bar)
		body.WriteString("  ")
		body.WriteString(row.pct)
		body.WriteString("\n")

		if m.expanded[i] {
			for _, sub := range row.subs {
				body.WriteString(sub)
				body.WriteString("\n")
			}
		}
	}

	right := fmt.Sprintf("%d 个分类", len(c.cats))
	return contentBox("能力图谱", body.String(), m.width, m.height, right)
}
