package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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

func (m ChallengesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.Type == tea.KeyUp || (msg.Type == tea.KeyRunes && string(msg.Runes) == "k"):
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.challenges) - 1
			}
		case msg.Type == tea.KeyDown || (msg.Type == tea.KeyRunes && string(msg.Runes) == "j"):
			m.cursor++
			if m.cursor >= len(m.challenges) {
				m.cursor = 0
			}
		case msg.Type == tea.KeyEnter:
			ch := m.challenges[m.cursor]
			return m, func() tea.Msg { return ChallengeSelectedMsg{Challenge: ch} }
		case msg.Type == tea.KeyEsc:
			return m, func() tea.Msg { return GoBackMsg{} }
		}
	}
	return m, nil
}

func (m ChallengesModel) View() string {
	var b strings.Builder

	label, ok := categoryLabels[m.category]
	if !ok {
		label = m.category
	}
	b.WriteString(TitleStyle.Render("  " + label))
	b.WriteString("\n\n")

	for i, ch := range m.challenges {
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

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("↑/k 上移 · ↓/j 下移 · Enter 选择 · Esc 返回"))

	return BoxStyle.Render(b.String())
}
