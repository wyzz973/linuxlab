package tui

import (
	"fmt"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sd3/linuxlab/internal/reference"
)

// ReferenceModel is the command reference TUI screen with search and detail modes.
type ReferenceModel struct {
	commands    []reference.CommandRef
	filtered    []reference.CommandRef
	query       string
	cursor      int
	showDetail  bool
	detailIndex int
	width       int
	height      int
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

func (m ReferenceModel) Init() tea.Cmd { return nil }

func (m ReferenceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		if m.showDetail {
			return m.updateDetail(msg)
		}
		return m.updateList(msg)
	}
	return m, nil
}

func (m ReferenceModel) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.Type == tea.KeyEsc:
		return m, func() tea.Msg { return GoBackMsg{} }

	case msg.Type == tea.KeyUp || (msg.Type == tea.KeyRunes && string(msg.Runes) == "k" && m.query == ""):
		m.cursor--
		if m.cursor < 0 {
			if len(m.filtered) > 0 {
				m.cursor = len(m.filtered) - 1
			} else {
				m.cursor = 0
			}
		}

	case msg.Type == tea.KeyDown || (msg.Type == tea.KeyRunes && string(msg.Runes) == "j" && m.query == ""):
		m.cursor++
		if m.cursor >= len(m.filtered) {
			m.cursor = 0
		}

	case msg.Type == tea.KeyEnter:
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
		}

	case msg.Type == tea.KeyRunes:
		m.query += string(msg.Runes)
		m.filtered = reference.Search(m.query, m.commands)
		m.cursor = 0
	}

	return m, nil
}

func (m ReferenceModel) updateDetail(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.Type == tea.KeyEsc {
		m.showDetail = false
		return m, nil
	}
	return m, nil
}

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

func (m ReferenceModel) listView() string {
	var body strings.Builder

	// Search box
	queryDisplay := m.query
	if queryDisplay == "" {
		queryDisplay = DimStyle.Render("输入命令名搜索...")
	}
	body.WriteString(fmt.Sprintf("/ %s\n", queryDisplay))
	body.WriteString(DimStyle.Render(strings.Repeat("─", 30)) + "\n\n")

	// Command list
	maxVisible := m.height - 12
	if maxVisible < 5 {
		maxVisible = 5
	}
	count := len(m.filtered)
	if count > maxVisible {
		count = maxVisible
	}

	for i := 0; i < count; i++ {
		cmd := m.filtered[i]
		cursor := "  "
		style := DimStyle
		if i == m.cursor {
			cursor = CurrentIcon + " "
			style = SelectedStyle
		}
		body.WriteString(fmt.Sprintf("%s%-12s %s\n", cursor, style.Render(cmd.Name), DimStyle.Render(cmd.Brief)))
	}

	if len(m.filtered) > maxVisible {
		body.WriteString(fmt.Sprintf("\n%s\n", DimStyle.Render(fmt.Sprintf("... 还有 %d 个命令", len(m.filtered)-maxVisible))))
	}

	if len(m.filtered) == 0 && m.query != "" {
		body.WriteString(DimStyle.Render("未找到匹配的命令") + "\n")
	}

	box := contentBox("命令速查", body.String(), m.width, m.height, "")
	status := statusBar("输入搜索 · ↑↓ 选择", "Enter 查看详情 · Esc 返回", m.width)

	return verticalCenter(box, status, m.height)
}

func (m ReferenceModel) detailView() string {
	if m.detailIndex >= len(m.filtered) {
		return ""
	}

	cmd := m.filtered[m.detailIndex]

	var body strings.Builder

	body.WriteString(cmd.Brief)
	body.WriteString("\n\n")

	if len(cmd.Examples) > 0 {
		body.WriteString(sectionTitle("示例", m.width))
		body.WriteString("\n\n")
		for _, ex := range cmd.Examples {
			body.WriteString(fmt.Sprintf("%s %s\n", DimStyle.Render("▸"), ex.Desc))
			body.WriteString(fmt.Sprintf("  $ %s\n\n", highlightPlaceholders(ex.Cmd)))
		}
	}

	if len(cmd.RelatedChallenges) > 0 {
		body.WriteString(sectionTitle("相关挑战", m.width))
		body.WriteString("\n\n")
		body.WriteString(DimStyle.Render(strings.Join(cmd.RelatedChallenges, ", ")))
		body.WriteString("\n")
	}

	box := contentBox(cmd.Name, body.String(), m.width, m.height, "")
	status := statusBar("", "Esc 返回列表", m.width)

	return verticalCenter(box, status, m.height)
}
