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

type moduleEntry struct {
	category   string
	label      string
	total      int
	passed     int
}

// ModulesModel is the module selection screen.
type ModulesModel struct {
	modules []moduleEntry
	cursor  int
	store   *progress.Store
}

// NewModulesModel creates a new module selection screen.
func NewModulesModel(cats map[string][]*challenge.Challenge, store *progress.Store) tea.Model {
	var modules []moduleEntry
	for cat, challenges := range cats {
		label, ok := categoryLabels[cat]
		if !ok {
			label = cat
		}
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
	var b strings.Builder
	b.WriteString(TitleStyle.Render("  选择学习模块"))
	b.WriteString("\n\n")

	for i, mod := range m.modules {
		cursor := "  "
		style := DimStyle
		if i == m.cursor {
			cursor = CurrentIcon + " "
			style = SelectedStyle
		}

		pct := 0.0
		if mod.total > 0 {
			pct = float64(mod.passed) / float64(mod.total)
		}
		bar := ProgressBar(pct, 10)
		b.WriteString(fmt.Sprintf("%s%s  %s  %d/%d\n", cursor, style.Render(mod.label), bar, mod.passed, mod.total))
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("↑/k 上移 · ↓/j 下移 · Enter 选择 · Esc 返回"))

	return BoxStyle.Render(b.String())
}
