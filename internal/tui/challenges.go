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
	// Box has ~2 lines border + 2 lines padding, status bar 1 line, vertical centering
	// Each challenge is 1 line; leave room for scroll indicators
	v := m.height - 12
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
	label := CategoryLabel(m.category)

	// Count passed
	passedCount := 0
	for _, ch := range m.challenges {
		if entry, exists := m.store.Data.Challenges[ch.ID]; exists && entry.Status == "passed" {
			passedCount++
		}
	}

	var body strings.Builder

	maxVis := m.maxVisible()
	end := m.offset + maxVis
	if end > len(m.challenges) {
		end = len(m.challenges)
	}

	// Calculate max title width for right-alignment of stars
	maxTitleW := 0
	for i := m.offset; i < end; i++ {
		tw := lipgloss.Width(m.challenges[i].Title)
		if tw > maxTitleW {
			maxTitleW = tw
		}
	}
	if maxTitleW > 40 {
		maxTitleW = 40
	}

	for i := m.offset; i < end; i++ {
		ch := m.challenges[i]

		// Status icon
		icon := PendingIcon
		if entry, exists := m.store.Data.Challenges[ch.ID]; exists {
			switch entry.Status {
			case "passed":
				icon = PassedIcon
			case "failed":
				icon = FailedIcon
			}
		}

		// Cursor and styling
		cursor := "  "
		style := DimStyle
		if i == m.cursor {
			cursor = CurrentIcon + " "
			style = SelectedStyle
			icon = CurrentIcon
		}

		title := style.Render(ch.Title)
		titleW := lipgloss.Width(ch.Title)
		pad := ""
		if maxTitleW > titleW {
			pad = strings.Repeat(" ", maxTitleW-titleW)
		}

		stars := DifficultyStars(ch.Difficulty)
		body.WriteString(fmt.Sprintf("%s %s  %s%s  %s\n", cursor, icon, title, pad, stars))
	}

	// Scroll indicator
	if end < len(m.challenges) {
		remaining := len(m.challenges) - end
		body.WriteString(DimStyle.Render(fmt.Sprintf("%s%s", strings.Repeat(" ", maxTitleW-2), fmt.Sprintf("▼ 还有 %d 题", remaining))))
		body.WriteString("\n")
	}

	rightLabel := fmt.Sprintf("%d/%d", passedCount, len(m.challenges))
	box := contentBox(label, body.String(), m.width, m.height, rightLabel)
	status := statusBar(label, "↑↓ 选择 · Enter 确认 · Esc 返回", m.width)

	return verticalCenter(box, status, m.height)
}
