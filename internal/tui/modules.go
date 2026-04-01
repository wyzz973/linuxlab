package tui

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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

type moduleEntry struct {
	category string
	label    string
	total    int
	passed   int
}

// ModulesModel is the module selection screen.
type ModulesModel struct {
	modules []moduleEntry
	cursor  int
	store   *progress.Store
	width   int
	height  int
}

// NewModulesModel creates a new module selection screen.
func NewModulesModel(cats map[string][]*challenge.Challenge, store *progress.Store) tea.Model {
	var modules []moduleEntry
	for cat, challenges := range cats {
		label := CategoryLabel(cat)
		passed := 0
		for _, ch := range challenges {
			if entry, exists := store.Data.Challenges[ch.ID]; exists && entry.Status == "passed" {
				passed++
			}
		}
		modules = append(modules, moduleEntry{
			category: cat,
			label:    label,
			total:    len(challenges),
			passed:   passed,
		})
	}
	sort.Slice(modules, func(i, j int) bool {
		return modules[i].category < modules[j].category
	})
	return ModulesModel{modules: modules, cursor: 0, store: store}
}

func (m ModulesModel) Init() tea.Cmd { return nil }

func (m ModulesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				m.cursor = len(m.modules) - 1
			}
		case msg.Type == tea.KeyDown || (msg.Type == tea.KeyRunes && string(msg.Runes) == "j"):
			m.cursor++
			if m.cursor >= len(m.modules) {
				m.cursor = 0
			}
		case msg.Type == tea.KeyEnter:
			cat := m.modules[m.cursor].category
			return m, func() tea.Msg { return ModuleSelectedMsg{Category: cat} }
		case msg.Type == tea.KeyEsc:
			return m, func() tea.Msg { return GoBackMsg{} }
		}
	}
	return m, nil
}

func (m ModulesModel) View() string {
	var body strings.Builder

	for i, mod := range m.modules {
		pct := 0.0
		if mod.total > 0 {
			pct = float64(mod.passed) / float64(mod.total)
		}

		// Module name line
		if i == m.cursor {
			body.WriteString(fmt.Sprintf("%s %s\n", CurrentIcon, SelectedStyle.Render(mod.label)))
		} else {
			body.WriteString(fmt.Sprintf("  %s\n", DimStyle.Render(mod.label)))
		}

		// Progress bar line (indented under name)
		bar := ProgressBar(pct, 30)
		countStr := fmt.Sprintf("%d/%d", mod.passed, mod.total)
		pctStr := fmt.Sprintf("%3.0f%%", pct*100)
		body.WriteString(fmt.Sprintf("  %s  %s  %s\n", bar, DimStyle.Render(countStr), DimStyle.Render(pctStr)))

		// Spacing between modules
		if i < len(m.modules)-1 {
			body.WriteString("\n")
		}
	}

	box := contentBox("选择学习模块", body.String(), m.width, m.height, "")
	status := statusBar("选择学习模块", "↑↓ 选择 · Enter 确认 · Esc 返回", m.width)

	return verticalCenter(box, status, m.height)
}
