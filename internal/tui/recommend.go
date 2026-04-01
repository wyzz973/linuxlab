package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sd3/linuxlab/internal/challenge"
)

// RecommendModel shows weakness-based challenge recommendations.
type RecommendModel struct {
	challenges []*challenge.Challenge
	cursor     int
}

// NewRecommendModel creates a new recommendation screen.
func NewRecommendModel(challenges []*challenge.Challenge) tea.Model {
	return RecommendModel{challenges: challenges}
}

func (m RecommendModel) Init() tea.Cmd { return nil }

func (m RecommendModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.Type == tea.KeyUp || (msg.Type == tea.KeyRunes && string(msg.Runes) == "k"):
			if len(m.challenges) > 0 {
				m.cursor--
				if m.cursor < 0 {
					m.cursor = len(m.challenges) - 1
				}
			}
		case msg.Type == tea.KeyDown || (msg.Type == tea.KeyRunes && string(msg.Runes) == "j"):
			if len(m.challenges) > 0 {
				m.cursor++
				if m.cursor >= len(m.challenges) {
					m.cursor = 0
				}
			}
		case msg.Type == tea.KeyEnter:
			if len(m.challenges) > 0 {
				ch := m.challenges[m.cursor]
				return m, func() tea.Msg { return ChallengeSelectedMsg{Challenge: ch} }
			}
		case msg.Type == tea.KeyEsc || (msg.Type == tea.KeyRunes && string(msg.Runes) == "q"):
			return m, func() tea.Msg { return GoBackMsg{} }
		}
	}
	return m, nil
}

func (m RecommendModel) View() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render("  薄弱推荐"))
	b.WriteString("\n\n")

	if len(m.challenges) == 0 {
		b.WriteString(DimStyle.Render("  暂无推荐。完成一些题目后再来看看！"))
		b.WriteString("\n\n")
		b.WriteString(HelpStyle.Render("Esc 返回"))
		return BoxStyle.Render(b.String())
	}

	for i, ch := range m.challenges {
		cursor := "  "
		style := DimStyle
		if i == m.cursor {
			cursor = CurrentIcon + " "
			style = SelectedStyle
		}
		cat := CategoryLabel(ch.Category)
		stars := DifficultyStars(ch.Difficulty)
		b.WriteString(fmt.Sprintf("%s%s  %s  %s  %s\n", cursor, style.Render(ch.Title), stars, DimStyle.Render(cat), DimStyle.Render(ch.Subcategory)))
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("↑/k 上移 · ↓/j 下移 · Enter 选择 · Esc 返回"))

	return BoxStyle.Render(b.String())
}
