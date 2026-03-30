package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sd3/linuxlab/internal/progress"
)

// SkillMapModel is the skill map view with expandable categories.
type SkillMapModel struct {
	store    *progress.Store
	skillMap *progress.SkillMap
	cursor   int
	expanded map[int]bool
}

// NewSkillMapModel creates a new skill map view.
func NewSkillMapModel(store *progress.Store) tea.Model {
	sm := progress.BuildSkillMap(store)
	return SkillMapModel{
		store:    store,
		skillMap: sm,
		cursor:   0,
		expanded: make(map[int]bool),
	}
}

func (m SkillMapModel) Init() tea.Cmd { return nil }

func (m SkillMapModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.Type == tea.KeyUp || (msg.Type == tea.KeyRunes && string(msg.Runes) == "k"):
			m.cursor--
			if m.cursor < 0 {
				count := len(m.skillMap.Categories)
				if count == 0 {
					m.cursor = 0
				} else {
					m.cursor = count - 1
				}
			}
		case msg.Type == tea.KeyDown || (msg.Type == tea.KeyRunes && string(msg.Runes) == "j"):
			m.cursor++
			if m.cursor >= len(m.skillMap.Categories) {
				m.cursor = 0
			}
		case msg.Type == tea.KeyEnter:
			if len(m.skillMap.Categories) > 0 {
				// Copy map to avoid shared mutation
				newExpanded := make(map[int]bool)
				for k, v := range m.expanded {
					newExpanded[k] = v
				}
				newExpanded[m.cursor] = !newExpanded[m.cursor]
				m.expanded = newExpanded
			}
		case msg.Type == tea.KeyEsc:
			return m, func() tea.Msg { return GoBackMsg{} }
		}
	}
	return m, nil
}

func (m SkillMapModel) View() string {
	var b strings.Builder

	b.WriteString(TitleStyle.Render("  能力图谱"))
	b.WriteString("\n\n")

	if m.skillMap.TotalCount > 0 {
		pct := m.skillMap.OverallScore
		bar := ProgressBar(pct, 20)
		b.WriteString(fmt.Sprintf("总进度: %s %.0f%%  (%d/%d)\n\n",
			bar, pct*100, m.skillMap.TotalPassed, m.skillMap.TotalCount))
	} else {
		b.WriteString(DimStyle.Render("暂无学习记录"))
		b.WriteString("\n\n")
	}

	for i, cat := range m.skillMap.Categories {
		cursor := "  "
		style := DimStyle
		if i == m.cursor {
			cursor = CurrentIcon + " "
			style = SelectedStyle
		}

		arrow := "▶"
		if m.expanded[i] {
			arrow = "▼"
		}

		label, ok := categoryLabels[cat.Name]
		if !ok {
			label = cat.Name
		}

		bar := ProgressBar(cat.Score, 10)
		b.WriteString(fmt.Sprintf("%s%s %s  %s  %.0f%%\n", cursor, arrow, style.Render(label), bar, cat.Score*100))

		if m.expanded[i] {
			for _, sub := range cat.Subcategories {
				subBar := ProgressBar(sub.Score, 8)
				b.WriteString(fmt.Sprintf("      %s  %s  %d/%d\n",
					DimStyle.Render(sub.Name), subBar, sub.Passed, sub.Total))
			}
		}
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("↑/k 上移 · ↓/j 下移 · Enter 展开/折叠 · Esc 返回"))

	return BoxStyle.Render(b.String())
}
