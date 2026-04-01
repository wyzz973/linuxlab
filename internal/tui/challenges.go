package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sd3/linuxlab/internal/challenge"
	"github.com/sd3/linuxlab/internal/progress"
)

// ChallengeSelectedMsg is sent when the user selects a challenge.
type ChallengeSelectedMsg struct {
	Challenge *challenge.Challenge
}

// ChallengesModel is the challenge list screen for a category.
type ChallengesModel struct {
	category   string
	challenges []*challenge.Challenge
	store      *progress.Store
	cursor     int
	offset     int
	width      int
	height     int
}

// NewChallengesModel creates a new challenge list screen.
func NewChallengesModel(category string, challenges []*challenge.Challenge, store *progress.Store) tea.Model {
	return ChallengesModel{
		category:   category,
		challenges: challenges,
		store:      store,
		cursor:     0,
	}
}

func (m ChallengesModel) Init() tea.Cmd { return nil }

func (m ChallengesModel) maxVisible() int {
	v := m.height - 10 // account for header, footer, box padding
	if v < 5 {
		v = 5
	}
	return v
}

func (m ChallengesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				m.cursor = len(m.challenges) - 1
				// Jump offset to show the last item
				maxVis := m.maxVisible()
				if len(m.challenges) > maxVis {
					m.offset = len(m.challenges) - maxVis
				}
			}
		case msg.Type == tea.KeyDown || (msg.Type == tea.KeyRunes && string(msg.Runes) == "j"):
			m.cursor++
			if m.cursor >= len(m.challenges) {
				m.cursor = 0
				m.offset = 0
			}
		case msg.Type == tea.KeyEnter:
			ch := m.challenges[m.cursor]
			return m, func() tea.Msg { return ChallengeSelectedMsg{Challenge: ch} }
		case msg.Type == tea.KeyEsc:
			return m, func() tea.Msg { return GoBackMsg{} }
		}

		// Adjust scroll offset
		maxVis := m.maxVisible()
		if m.cursor < m.offset {
			m.offset = m.cursor
		}
		if m.cursor >= m.offset+maxVis {
			m.offset = m.cursor - maxVis + 1
		}
	}
	return m, nil
}

func (m ChallengesModel) View() string {
	var b strings.Builder

	label := CategoryLabel(m.category)
	b.WriteString(TitleStyle.Render("  " + label))
	b.WriteString("\n\n")

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

		icon := PendingIcon
		if entry, exists := m.store.Data.Challenges[ch.ID]; exists {
			switch entry.Status {
			case "passed":
				icon = PassedIcon
			case "failed":
				icon = FailedIcon
			}
		}

		stars := DifficultyStars(ch.Difficulty)
		b.WriteString(fmt.Sprintf("%s%s %s  %s  %s\n", cursor, icon, style.Render(ch.Title), stars, DimStyle.Render(ch.ID)))
	}

	if end < len(m.challenges) {
		b.WriteString(DimStyle.Render(fmt.Sprintf("  ▼ 下滑查看更多 (%d/%d)", end, len(m.challenges))))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("↑/k 上移 · ↓/j 下移 · Enter 选择 · Esc 返回"))

	boxWidth := responsiveBoxWidth(m.width)
	content := BoxStyle.Width(boxWidth).Render(b.String())
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}
