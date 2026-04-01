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
	desc  string
	key   string
}

// MenuModel is the main menu screen.
type MenuModel struct {
	items           []menuItem
	cursor          int
	width           int
	height          int
	totalChallenges int
	totalModules    int
}

// NewMenuModel creates a new main menu.
func NewMenuModel() tea.Model {
	return MenuModel{
		items: []menuItem{
			{label: "开始练习", desc: "选择模块开始学习", key: "practice"},
			{label: "能力图谱", desc: "查看各项技能掌握程度", key: "skillmap"},
			{label: "薄弱推荐", desc: "针对薄弱环节推荐练习", key: "recommend"},
			{label: "命令速查", desc: "快速查阅常用命令用法", key: "reference"},
		},
		cursor: 0,
	}
}

// NewMenuModelWithStats creates a menu model with challenge/module counts.
func NewMenuModelWithStats(totalChallenges, totalModules int) tea.Model {
	return MenuModel{
		items: []menuItem{
			{label: "开始练习", desc: "选择模块开始学习", key: "practice"},
			{label: "能力图谱", desc: "查看各项技能掌握程度", key: "skillmap"},
			{label: "薄弱推荐", desc: "针对薄弱环节推荐练习", key: "recommend"},
			{label: "命令速查", desc: "快速查阅常用命令用法", key: "reference"},
		},
		cursor:          0,
		totalChallenges: totalChallenges,
		totalModules:    totalModules,
	}
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
	var body strings.Builder

	// Welcome message
	body.WriteString(DimStyle.Render("欢迎来到 LinuxLab！在真实终端中掌握 Linux 技能。"))
	body.WriteString("\n\n")

	// Menu items with descriptions
	for i, item := range m.items {
		if i == m.cursor {
			body.WriteString(fmt.Sprintf("%s %s", CurrentIcon, SelectedStyle.Render(item.label)))
			body.WriteString("          ")
			body.WriteString(DimStyle.Render(item.desc))
			body.WriteString("\n")
		} else {
			body.WriteString(fmt.Sprintf("  %s", DimStyle.Render(item.label)))
			body.WriteString("          ")
			body.WriteString(SubtleStyle.Render(item.desc))
			body.WriteString("\n")
		}
	}

	box := contentBox("LinuxLab", body.String(), m.width, m.height, "")

	// Stats in status bar
	leftStatus := ""
	if m.totalChallenges > 0 {
		leftStatus = fmt.Sprintf("%d 道挑战 · %d 个模块", m.totalChallenges, m.totalModules)
	}
	status := statusBar(leftStatus, "↑↓ 选择 · Enter 确认 · q 退出", m.width)

	return verticalCenter(box, status, m.height)
}
