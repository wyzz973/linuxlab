package tui

import (
	"fmt"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sd3/linuxlab/internal/challenge"
	"github.com/sd3/linuxlab/internal/progress"
	"github.com/sd3/linuxlab/internal/reference"
)

// ReferenceModel is the command reference TUI screen with search and detail modes.
type ReferenceModel struct {
	commands    []reference.CommandRef
	filtered    []reference.CommandRef
	query       string
	cursor      int
	offset      int
	showDetail  bool
	detailIndex int
	width       int
	height      int

	// Optional challenge lookup for the detail view's related-challenge list.
	// When absent the list falls back to bare IDs.
	challenges map[string]*challenge.Challenge
	store      *progress.Store
}

// NewReferenceModel creates a new reference browser model.
func NewReferenceModel(refs *reference.ReferenceData) tea.Model {
	cmds := refs.Commands
	filtered := make([]reference.CommandRef, len(cmds))
	copy(filtered, cmds)
	return ReferenceModel{
		commands: cmds,
		filtered: filtered,
	}
}

// NewReferenceModelWithChallenges creates a reference browser that can resolve
// related challenge IDs to titles and live completion status.
func NewReferenceModelWithChallenges(refs *reference.ReferenceData, cats map[string][]*challenge.Challenge, store *progress.Store) tea.Model {
	m := NewReferenceModel(refs).(ReferenceModel)
	m.challenges = make(map[string]*challenge.Challenge)
	for _, chs := range cats {
		for _, ch := range chs {
			m.challenges[ch.ID] = ch
		}
	}
	m.store = store
	return m
}

func (m ReferenceModel) Init() tea.Cmd { return nil }

func (m ReferenceModel) maxVisible() int {
	return visibleListRows(m.height)
}

func (m ReferenceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Re-clamp scroll in model state so Update and View share the same
		// viewport instead of View recomputing a throwaway offset.
		l := ScrollList{Cursor: m.cursor, Offset: m.offset}
		l.Clamp(len(m.filtered), m.maxVisible())
		m.cursor, m.offset = l.Cursor, l.Offset
		return m, nil
	case tea.KeyMsg:
		if m.showDetail {
			return m.updateDetail(msg)
		}
		return m.updateList(msg)
	}
	return m, nil
}

// updateList handles keys in list mode. Search input has priority (locked
// behavior id=12): printable runes (and space) always go into the query, so
// command names containing j/k/g/G/q (grep, kill, jobs, ...) can be typed.
// List navigation is limited to arrow keys, PgUp/PgDn, Home/End; going back is
// Esc/Left only.
func (m ReferenceModel) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.Type == tea.KeyEsc:
		if m.query != "" {
			m.clearQuery()
			return m, nil
		}
		return m, func() tea.Msg { return GoBackMsg{} }

	case m.query == "" && msg.Type == tea.KeyLeft:
		return m, func() tea.Msg { return GoBackMsg{} }

	case msg.Type == tea.KeyCtrlU:
		if m.query != "" {
			m.clearQuery()
		}

	case isSelectKey(msg):
		if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
			m.showDetail = true
			m.detailIndex = m.cursor
		}

	case msg.Type == tea.KeyBackspace:
		if len(m.query) > 0 {
			// Remove last rune
			runes := []rune(m.query)
			m.query = string(runes[:len(runes)-1])
			m.filtered = reference.Search(m.query, m.commands)
			m.cursor = 0
			m.offset = 0
		}

	case msg.Type == tea.KeyRunes:
		m.appendQuery(string(msg.Runes))

	case msg.Type == tea.KeySpace:
		m.appendQuery(" ")

	default:
		l := ScrollList{Cursor: m.cursor, Offset: m.offset}
		if l.Move(msg, len(m.filtered), m.maxVisible(), false) {
			m.cursor, m.offset = l.Cursor, l.Offset
		}
	}

	return m, nil
}

func (m ReferenceModel) updateDetail(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if isBackKey(msg) {
		m.showDetail = false
		return m, nil
	}
	return m, nil
}

func (m *ReferenceModel) clearQuery() {
	m.query = ""
	m.filtered = reference.Search("", m.commands)
	m.cursor = 0
	m.offset = 0
}

// appendQuery appends text to the search query and refreshes the result list.
func (m *ReferenceModel) appendQuery(text string) {
	m.query += text
	m.filtered = reference.Search(m.query, m.commands)
	m.cursor = 0
	m.offset = 0
}

// refCursorBlock is the input-cursor block shown at the end of the query.
const refCursorBlock = "▌"

// refNameColW is the fixed column width for command names in the list.
const refNameColW = 12

// Pre-built private styles for the reference screen (iron rule: never call
// lipgloss.NewStyle() inside View). The selected row carries the BgFocus
// background on every segment so the highlight covers the whole inner row.
var (
	refPromptStyle = lipgloss.NewStyle().Foreground(currentTheme.Primary).Bold(true)
	refBlockStyle  = lipgloss.NewStyle().Foreground(currentTheme.Primary)
	refNameStyle   = lipgloss.NewStyle().Foreground(currentTheme.Secondary)
	refSelPrefix   = lipgloss.NewStyle().Foreground(currentTheme.Primary).Bold(true).Background(currentTheme.BgFocus)
	refSelName     = lipgloss.NewStyle().Foreground(currentTheme.Primary).Bold(true).Background(currentTheme.BgFocus)
	refSelBrief    = lipgloss.NewStyle().Foreground(currentTheme.FgBase).Background(currentTheme.BgFocus)
)

var placeholderStyle = lipgloss.NewStyle().Foreground(ColorYellow).Bold(true)
var placeholderRe = regexp.MustCompile(`\{\{[^}]+\}\}`)

func highlightPlaceholders(s string) string {
	return placeholderRe.ReplaceAllStringFunc(s, func(match string) string {
		return placeholderStyle.Render(match)
	})
}

func (m ReferenceModel) View() string {
	if m.showDetail {
		return m.detailView()
	}
	return m.listView()
}

// refInnerWidth is the content width inside the contentBox borders/padding,
// mirroring the box geometry so full-width elements (divider, selected-row
// background) line up with the box exactly.
func (m ReferenceModel) refInnerWidth() int {
	return maxInt(boxWidth(m.width)-6, 0)
}

// queryLine renders the search input row: "/ query▌  输入即搜索 · ctrl+u 清空".
func (m ReferenceModel) queryLine() string {
	var b strings.Builder
	b.WriteString(refPromptStyle.Render("/ "))
	if m.query != "" {
		b.WriteString(S.Text.Render(m.query))
	}
	b.WriteString(refBlockStyle.Render(refCursorBlock))
	b.WriteString("  ")
	b.WriteString(S.Dim.Render("输入即搜索 · ctrl+u 清空"))
	return b.String()
}

// renderRefRow renders one command row. The selected row gets the unified
// BgFocus whole-row highlight; padding is measured on the raw text first
// (width discipline) and the trailing gap keeps the background running to the
// inner edge of the box.
func renderRefRow(cmd reference.CommandRef, selected bool, inner int) string {
	name := padRight(cmd.Name, refNameColW)
	if !selected {
		return "  " + refNameStyle.Render(name) + " " + S.Dim.Render(cmd.Brief)
	}
	row := refSelPrefix.Render(IconCursor+" ") + refSelName.Render(name) + refSelBrief.Render(" "+cmd.Brief)
	if gap := inner - displayWidth(row); gap > 0 {
		row += S.SelectedRow.Render(strings.Repeat(" ", gap))
	}
	return row
}

func (m ReferenceModel) listView() string {
	inner := m.refInnerWidth()
	var body strings.Builder

	// Search input row with cursor block, then a full-width divider.
	body.WriteString(m.queryLine())
	body.WriteString("\n")
	body.WriteString(S.Subtle.Render(strings.Repeat("─", inner)))
	body.WriteString("\n")

	if len(m.filtered) == 0 {
		body.WriteString("\n")
		if m.query != "" {
			body.WriteString(EmptyState(fmt.Sprintf("没有匹配「%s」的命令", m.query), "试试更短的关键词。", "Ctrl+U 清空 · Esc 返回"))
		} else {
			body.WriteString(EmptyState("暂无命令数据", "命令速查表为空。", "Esc 返回"))
		}
		body.WriteString("\n")
		return contentBox("命令速查", body.String(), m.width, m.height, m.countLabel())
	}

	// The visible window and ▲/▼ overflow indicators come from the shared
	// ScrollList component state (Update keeps cursor/offset clamped, so View
	// reads the same viewport the key handlers use).
	l := ScrollList{Cursor: m.cursor, Offset: m.offset}
	start, end := l.Window(len(m.filtered), m.maxVisible())

	if start > 0 {
		body.WriteString(S.Dim.Render(fmt.Sprintf("%s 还有 %d 个命令", ArrowUp, start)))
		body.WriteString("\n")
	}
	for i := start; i < end; i++ {
		body.WriteString(renderRefRow(m.filtered[i], i == m.cursor, inner))
		body.WriteString("\n")
	}
	if end < len(m.filtered) {
		body.WriteString(S.Dim.Render(fmt.Sprintf("%s 还有 %d 个命令", ArrowDown, len(m.filtered)-end)))
		body.WriteString("\n")
	}

	return contentBox("命令速查", body.String(), m.width, m.height, m.countLabel())
}

// countLabel is the top-border right slot: "matched/total".
func (m ReferenceModel) countLabel() string {
	return fmt.Sprintf("%d/%d", len(m.filtered), len(m.commands))
}

// renderRelatedChallenge renders one entry of the detail view's "相关挑战"
// list: completion icon + challenge title, truncated to the box width. When
// the challenge is unknown (no lookup wired, or a stale ID in the reference
// data) it degrades to the bare ID.
func (m ReferenceModel) renderRelatedChallenge(id string) string {
	ch, known := m.challenges[id]
	if !known {
		return S.Dim.Render(id)
	}

	icon, label := S.IconTodo, S.Dim
	if m.store != nil {
		if entry, ok := m.store.Data.Challenges[id]; ok && entry != nil {
			switch entry.Status {
			case "passed":
				icon, label = S.IconPass, S.Text
			case "failed":
				icon, label = S.IconFail, S.Text
			}
		}
	}

	// contentBox geometry: border (2) + padding (4), minus the icon prefix.
	avail := maxInt(boxWidth(m.width)-6-2, 1)
	return icon + " " + label.Render(truncateWidth(ch.Title, avail))
}

func (m ReferenceModel) detailView() string {
	if m.detailIndex >= len(m.filtered) {
		return ""
	}

	cmd := m.filtered[m.detailIndex]

	var body strings.Builder

	body.WriteString(S.Text.Render(cmd.Brief))
	body.WriteString("\n\n")

	if len(cmd.Examples) > 0 {
		body.WriteString(sectionTitle("示例", m.width))
		body.WriteString("\n")
		for _, ex := range cmd.Examples {
			body.WriteString(S.Dim.Render("▸ " + ex.Desc))
			body.WriteString("\n")
			body.WriteString("  " + S.Subtle.Render("$") + " " + highlightPlaceholders(ex.Cmd))
			body.WriteString("\n")
		}
		body.WriteString("\n")
	}

	if len(cmd.RelatedChallenges) > 0 {
		body.WriteString(sectionTitle("相关挑战", m.width))
		body.WriteString("\n")
		for _, id := range cmd.RelatedChallenges {
			body.WriteString(m.renderRelatedChallenge(id))
			body.WriteString("\n")
		}
	}

	return contentBox(cmd.Name, body.String(), m.width, m.height, "命令详情")
}
