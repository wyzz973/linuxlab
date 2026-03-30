package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// MenuChoiceMsg is sent when the user selects a menu item.
type MenuChoiceMsg struct {
	Choice string
}

type menuItem struct {
	label string
	key   string
}

// MenuModel is the main menu screen.
type MenuModel struct {
	items  []menuItem
	cursor int
}

// NewMenuModel creates a new main menu.
func NewMenuModel() tea.Model {
	return MenuModel{
		items: []menuItem{
			{label: "开始练习", key: "practice"},
			{label: "能力图谱", key: "skillmap"},
			{label: "薄弱推荐", key: "recommend"},
			{label: "命令速查", key: "reference"},
		},
		cursor: 0,
	}
}

func (m MenuModel) Init() tea.Cmd { return nil }

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.Type == tea.KeyUp || (msg.Type == tea.KeyRunes && string(msg.Runes) == "k"):
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.items) - 1
			}
		case msg.Type == tea.KeyDown || (msg.Type == tea.KeyRunes && string(msg.Runes) == "j"):
			m.cursor++
			if m.cursor >= len(m.items) {
				m.cursor = 0
			}
		case msg.Type == tea.KeyEnter:
			choice := m.items[m.cursor].key
			return m, func() tea.Msg { return MenuChoiceMsg{Choice: choice} }
		case msg.Type == tea.KeyRunes && string(msg.Runes) == "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m MenuModel) View() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render("  LinuxLab - Linux 命令行学习"))
	b.WriteString("\n\n")

	for i, item := range m.items {
		cursor := "  "
		style := DimStyle
		if i == m.cursor {
			cursor = CurrentIcon + " "
			style = SelectedStyle
		}
		b.WriteString(fmt.Sprintf("%s%s\n", cursor, style.Render(item.label)))
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("↑/k 上移 · ↓/j 下移 · Enter 选择 · q 退出"))

	return BoxStyle.Render(b.String())
}
