package tui

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sd3/linuxlab/internal/challenge"
	"github.com/sd3/linuxlab/internal/progress"
)

// ModuleSelectedMsg is sent when the user selects a module/category.
type ModuleSelectedMsg struct {
	Category string
}

// GoBackMsg is sent when the user wants to go back.
type GoBackMsg struct{}

var categoryLabels = map[string]string{
	"linux-basics":    "Linux 基础命令",
	"vim":             "Vim 操作",
	"shell-scripting": "Shell 脚本",
	"ops":             "运维实战",
	"containers":      "容器与部署",
}

// CategoryLabel returns the display label for a category key.
func CategoryLabel(category string) string {
	if label, ok := categoryLabels[category]; ok {
		return label
	}
	return category
}

const (
	// modulesOverallBarWidth is the width of the overall progress bar.
	modulesOverallBarWidth = 24
	// moduleRowBarWidth is the preferred per-module bar width (design
	// mockup); narrow terminals shrink it down to moduleRowBarMinWidth so the
	// action word stays visible instead of being truncated away.
	moduleRowBarWidth    = 18
	moduleRowBarMinWidth = 8
)

type moduleEntry struct {
	category   string
	label      string
	challenges []*challenge.Challenge
}

// ModulesModel is the module selection screen.
type ModulesModel struct {
	modules []moduleEntry
	cursor  int
	offset  int
	store   *progress.Store
	width   int
	height  int
}

// NewModulesModel creates a new module selection screen.
func NewModulesModel(cats map[string][]*challenge.Challenge, store *progress.Store) tea.Model {
	var modules []moduleEntry
	for cat, challenges := range cats {
		modules = append(modules, moduleEntry{
			category:   cat,
			label:      CategoryLabel(cat),
			challenges: challenges,
		})
	}
	sort.Slice(modules, func(i, j int) bool {
		return modules[i].category < modules[j].category
	})
	return ModulesModel{modules: modules, cursor: 0, store: store}
}

// passedCount reads the number of passed challenges for a module live from the
// store, so a reused model always shows up-to-date progress.
func (m ModulesModel) passedCount(mod moduleEntry) int {
	passed := 0
	for _, ch := range mod.challenges {
		if entry, exists := m.store.Data.Challenges[ch.ID]; exists && entry.Status == "passed" {
			passed++
		}
	}
	return passed
}

// overallProgress sums passed/total across all modules.
func (m ModulesModel) overallProgress() (passed, total int) {
	for _, mod := range m.modules {
		total += len(mod.challenges)
		passed += m.passedCount(mod)
	}
	return passed, total
}

// scrollList adapts the model's cursor/offset fields (their names are locked
// by shared regression tests) to the ScrollList component.
func (m ModulesModel) scrollList() ScrollList {
	return ScrollList{Cursor: m.cursor, Offset: m.offset}
}

// rowsGapped reports whether module rows get a blank separator line between
// them (Standard/Wide layouts, per the design mockup); Compact terminals pack
// the rows to keep everything on screen.
func (m ModulesModel) rowsGapped() bool {
	mode := computeLayout(m.width, m.height).Mode
	return mode == ModeStandard || mode == ModeWide
}

// headerBlock is the fixed content above the module rows: the overall
// progress summary plus a separating blank line.
func (m ModulesModel) headerBlock() string {
	passed, total := m.overallProgress()
	return ProgressSummary("整体进度", passed, total, modulesOverallBarWidth) + "\n"
}

// maxVisible derives the module-row budget from the terminal height: shell
// header/footer (2) + contentBox chrome (4) + the measured header block, with
// gapped rows costing two lines each.
func (m ModulesModel) maxVisible() int {
	avail := m.height - 2 - 4 - lipgloss.Height(m.headerBlock())
	if m.rowsGapped() {
		avail = (avail + 1) / 2
	}
	if avail < 1 {
		return 1
	}
	return avail
}

func (m ModulesModel) Init() tea.Cmd { return nil }

func (m ModulesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Re-clamp the scroll offset so the cursor stays visible (id=18).
		l := m.scrollList()
		l.Clamp(len(m.modules), m.maxVisible())
		m.cursor, m.offset = l.Cursor, l.Offset
		return m, nil
	case tea.KeyMsg:
		switch {
		case isBackKey(msg):
			return m, func() tea.Msg { return GoBackMsg{} }
		case msg.Type == tea.KeyRunes:
			if cat, ok := m.categoryByNumber(string(msg.Runes)); ok {
				return m, func() tea.Msg { return ModuleSelectedMsg{Category: cat} }
			}
			return m.moveCursor(msg), nil
		case isSelectKey(msg):
			if len(m.modules) == 0 {
				return m, nil
			}
			cat := m.modules[m.cursor].category
			return m, func() tea.Msg { return ModuleSelectedMsg{Category: cat} }
		default:
			return m.moveCursor(msg), nil
		}
	}
	return m, nil
}

// moveCursor applies a navigation key through the ScrollList component
// (moveListCursor semantics, locked by ux_test).
func (m ModulesModel) moveCursor(msg tea.KeyMsg) ModulesModel {
	l := m.scrollList()
	if l.Move(msg, len(m.modules), m.maxVisible(), true) {
		m.cursor, m.offset = l.Cursor, l.Offset
	}
	return m
}

func (m ModulesModel) categoryByNumber(value string) (string, bool) {
	if len(value) != 1 {
		return "", false
	}
	index := int(value[0] - '1')
	if index < 0 || index >= len(m.modules) {
		return "", false
	}
	return m.modules[index].category, true
}

// moduleRowGeom is the shared column geometry of the module rows, derived
// once per View from display widths (CJK aware) so every row aligns.
type moduleRowGeom struct {
	innerW int // contentBox inner width
	labelW int // widest label + 2
	countW int // widest "passed/total" count
	barW   int // per-row bar width (shrinks on narrow terminals)
}

func (m ModulesModel) rowGeom() moduleRowGeom {
	g := moduleRowGeom{innerW: maxInt(boxWidth(m.width)-6, 0)}
	for _, mod := range m.modules {
		g.labelW = maxInt(g.labelW, displayWidth(mod.label))
		g.countW = maxInt(g.countW, displayWidth(fmt.Sprintf("%d/%d", m.passedCount(mod), len(mod.challenges))))
	}
	g.labelW += 2
	// Fixed row overhead besides label/count/bar: cursor (2) + badge (4) +
	// bar gap (2) + count gap (1) + percent (5) + widest action word (8).
	g.barW = clampInt(g.innerW-g.labelW-g.countW-22, moduleRowBarMinWidth, moduleRowBarWidth)
	return g
}

// renderModuleRow renders one module row:
//
//	› [1] Linux 基础命令  ██░░░░░░░░  10/98 10%  继续训练
//
// The selected row is padded to the full inner width so the BgFocus
// background reads as a whole-row highlight.
func (m ModulesModel) renderModuleRow(i int, selected bool, g moduleRowGeom) string {
	mod := m.modules[i]
	st := listRowPlain
	if selected {
		st = listRowFocus
	}

	modTotal := len(mod.challenges)
	modPassed := m.passedCount(mod)
	modPct := 0.0
	if modTotal > 0 {
		modPct = float64(modPassed) / float64(modTotal)
	}

	action := st.desc.Render("未开始")
	switch {
	case modTotal > 0 && modPassed == modTotal:
		action = st.success.Render("已完成")
	case modPassed > 0:
		action = st.success.Render("继续训练")
	}

	var b strings.Builder
	if selected {
		b.WriteString(st.cursor.Render(IconCursor + " "))
	} else {
		b.WriteString(st.fill.Render("  "))
	}
	b.WriteString(st.badge.Render(fmt.Sprintf("[%d] ", i+1)))
	b.WriteString(st.label.Render(padRight(mod.label, g.labelW)))
	b.WriteString(st.bar(modPct, g.barW))
	b.WriteString(st.fill.Render("  "))
	b.WriteString(st.desc.Render(padRight(fmt.Sprintf("%d/%d", modPassed, modTotal), g.countW) + " "))
	b.WriteString(st.desc.Render(padRight(fmt.Sprintf("%.0f%%", modPct*100), 5)))
	b.WriteString(action)

	row := b.String()
	if selected {
		if gap := g.innerW - displayWidth(row); gap > 0 {
			row += st.fill.Render(strings.Repeat(" ", gap))
		}
	}
	return row
}

func (m ModulesModel) View() string {
	title := "选择训练模块"
	if len(m.modules) == 0 {
		body := EmptyState("暂无训练模块。", "挑战目录为空，请检查安装是否完整。", "Esc 返回")
		return contentBox(title, body, m.width, m.height, "")
	}

	g := m.rowGeom()
	rows := m.scrollList().RenderRows(len(m.modules), m.maxVisible(), g.innerW,
		func(i int, selected bool) string {
			return m.renderModuleRow(i, selected, g)
		})

	rowSep := "\n"
	if m.rowsGapped() {
		rowSep = "\n\n"
	}

	var body strings.Builder
	body.WriteString(m.headerBlock())
	body.WriteString("\n")
	body.WriteString(strings.Join(rows, rowSep))

	rightLabel := fmt.Sprintf("%d 个模块", len(m.modules))
	return contentBox(title, body.String(), m.width, m.height, rightLabel)
}
