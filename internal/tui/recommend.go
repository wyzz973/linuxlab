package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sd3/linuxlab/internal/challenge"
)

// RecommendModel shows weakness-based challenge recommendations.
type RecommendModel struct {
	challenges []*challenge.Challenge
	cursor     int
	offset     int
	width      int
	height     int
}

// NewRecommendModel creates a new recommendation screen.
func NewRecommendModel(challenges []*challenge.Challenge) tea.Model {
	return RecommendModel{challenges: challenges}
}

func (m RecommendModel) Init() tea.Cmd { return nil }

func (m RecommendModel) maxVisible() int {
	v := m.height - 12
	if v < 5 {
		v = 5
	}
	return v
}

func (m RecommendModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch {
		case msg.Type == tea.KeyUp || (msg.Type == tea.KeyRunes && string(msg.Runes) == "k"):
			if len(m.challenges) > 0 {
				m.cursor--
				if m.cursor < 0 {
					m.cursor = len(m.challenges) - 1
					maxVis := m.maxVisible()
					if len(m.challenges) > maxVis {
						m.offset = len(m.challenges) - maxVis
					}
				}
			}
		case msg.Type == tea.KeyDown || (msg.Type == tea.KeyRunes && string(msg.Runes) == "j"):
			if len(m.challenges) > 0 {
				m.cursor++
				if m.cursor >= len(m.challenges) {
					m.cursor = 0
					m.offset = 0
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

		// Adjust scroll offset
		if len(m.challenges) > 0 {
			maxVis := m.maxVisible()
			if m.cursor < m.offset {
				m.offset = m.cursor
			}
			if m.cursor >= m.offset+maxVis {
				m.offset = m.cursor - maxVis + 1
			}
		}
	}
	return m, nil
}

func (m RecommendModel) View() string {
	var body strings.Builder

	if len(m.challenges) == 0 {
		body.WriteString(DimStyle.Render("暂无推荐。完成一些题目后再来看看！"))
		body.WriteString("\n")

		box := contentBox("薄弱推荐", body.String(), m.width, m.height, "")
		status := statusBar("薄弱推荐", "Esc 返回", m.width)
		return verticalCenter(box, status, m.height)
	}

	maxVis := m.maxVisible()
	end := m.offset + maxVis
	if end > len(m.challenges) {
		end = len(m.challenges)
	}

	// Calculate max title width for alignment
	maxTitleW := 0
	for i := m.offset; i < end; i++ {
		tw := lipgloss.Width(m.challenges[i].Title)
		if tw > maxTitleW {
			maxTitleW = tw
		}
	}
	if maxTitleW > 30 {
		maxTitleW = 30
	}

	for i := m.offset; i < end; i++ {
		ch := m.challenges[i]
		cursor := "  "
		style := DimStyle
		if i == m.cursor {
			cursor = CurrentIcon + " "
			style = SelectedStyle
		}

		title := style.Render(ch.Title)
		titleW := lipgloss.Width(ch.Title)
		pad := ""
		if maxTitleW > titleW {
			pad = strings.Repeat(" ", maxTitleW-titleW)
		}

		cat := CategoryLabel(ch.Category)
		stars := DifficultyStars(ch.Difficulty)
		body.WriteString(fmt.Sprintf("%s%s%s  %s  %s\n", cursor, title, pad, stars, DimStyle.Render(cat)))
	}

	if end < len(m.challenges) {
		body.WriteString(DimStyle.Render(fmt.Sprintf("  ▼ 还有 %d 题", len(m.challenges)-end)))
		body.WriteString("\n")
	}

	box := contentBox("薄弱推荐", body.String(), m.width, m.height, "")
	status := statusBar("薄弱推荐", "↑↓ 选择 · Enter 确认 · Esc 返回", m.width)

	return verticalCenter(box, status, m.height)
}
