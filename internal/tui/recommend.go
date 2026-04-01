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
	v := m.height - 4 // header + footer + 2 padding lines
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
	header := headerView("LinuxLab · 薄弱推荐", m.width)
	footer := footerView("↑/k 上移 · ↓/j 下移 · Enter 选择 · Esc 返回", m.width)

	var b strings.Builder

	if len(m.challenges) == 0 {
		b.WriteString(DimStyle.Render("暂无推荐。完成一些题目后再来看看！"))
		b.WriteString("\n")

		contentHeight := maxInt(1, m.height-6)
		content := fillContent(b.String(), m.width, contentHeight)
		return header + "\n" + content + "\n" + footer
	}

	maxVis := m.maxVisible()
	end := m.offset + maxVis
	if end > len(m.challenges) {
		end = len(m.challenges)
	}

	if m.offset > 0 {
		b.WriteString(DimStyle.Render("  ▲ 上滑查看更多"))
		b.WriteString("\n")
	}

	for i := m.offset; i < end; i++ {
		ch := m.challenges[i]
		cursor := "  "
		style := DimStyle
		if i == m.cursor {
			cursor = CurrentIcon + " "
			style = SelectedStyle
		}
		cat := CategoryLabel(ch.Category)
		stars := DifficultyStars(ch.Difficulty)
		b.WriteString(fmt.Sprintf("%s%-24s %s  %s\n", cursor, style.Render(ch.Title), stars, DimStyle.Render(cat)))
	}

	if end < len(m.challenges) {
		b.WriteString(DimStyle.Render(fmt.Sprintf("  ▼ 下滑查看更多 (%d/%d)", end, len(m.challenges))))
		b.WriteString("\n")
	}

	contentHeight := maxInt(1, m.height-6)
	content := fillContent(b.String(), m.width, contentHeight)

	return header + "\n" + content + "\n" + footer
}
