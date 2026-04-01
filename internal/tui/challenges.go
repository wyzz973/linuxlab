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
	v := m.height - 4 // header + footer + 2 padding lines
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
	header := headerView("LinuxLab · "+label, m.width)
	footer := footerView(fmt.Sprintf("%d/%d · ↑/k ↓/j 选择 · Enter 确认 · Esc 返回", m.cursor+1, len(m.challenges)), m.width)

	var b strings.Builder

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

	contentHeight := maxInt(1, m.height-2)
	content := fillContent(b.String(), m.width, contentHeight)

	return header + "\n" + content + "\n" + footer
}
