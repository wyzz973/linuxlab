package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sd3/linuxlab/internal/challenge"
	"github.com/sd3/linuxlab/internal/sandbox"
)

// LaunchChallengeMsg is sent when the user wants to start a challenge.
type LaunchChallengeMsg struct {
	Challenge *challenge.Challenge
	HintsUsed int
}

// DetailModel is the challenge detail screen with progressive hints.
type DetailModel struct {
	challenge       *challenge.Challenge
	hintLevel       int
	dockerAvailable bool
	width           int
	height          int
}

// NewDetailModel creates a new challenge detail screen.
func NewDetailModel(ch *challenge.Challenge) tea.Model {
	return DetailModel{
		challenge:       ch,
		hintLevel:       0,
		dockerAvailable: sandbox.DockerAvailable(),
	}
}

func (m DetailModel) Init() tea.Cmd { return nil }

func (m DetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch {
		case msg.Type == tea.KeyEnter:
			ch := m.challenge
			hints := m.hintLevel
			return m, func() tea.Msg {
				return LaunchChallengeMsg{Challenge: ch, HintsUsed: hints}
			}
		case msg.Type == tea.KeyEsc:
			return m, func() tea.Msg { return GoBackMsg{} }
		case msg.Type == tea.KeyRunes && string(msg.Runes) == "h":
			if m.hintLevel < len(m.challenge.Hints) {
				m.hintLevel++
			}
		}
	}
	return m, nil
}

func (m DetailModel) View() string {
	ch := m.challenge
	var b strings.Builder

	b.WriteString(TitleStyle.Render("  " + ch.Title))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("难度: %s\n", DifficultyStars(ch.Difficulty)))

	if len(ch.Tags) > 0 {
		b.WriteString(fmt.Sprintf("标签: %s\n", DimStyle.Render(strings.Join(ch.Tags, ", "))))
	}

	if ch.RequiresDocker && !m.dockerAvailable {
		b.WriteString("\n")
		b.WriteString(WarningStyle.Render("!! 需要 Docker (当前不可用，将使用本地模式)"))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(ch.Description)
	b.WriteString("\n")

	if m.hintLevel > 0 {
		b.WriteString("\n")
		b.WriteString(WarningStyle.Render("提示:"))
		b.WriteString("\n")
		for i := 0; i < m.hintLevel && i < len(ch.Hints); i++ {
			b.WriteString(fmt.Sprintf("  %d. %s\n", i+1, ch.Hints[i].Text))
		}
	}

	b.WriteString("\n")
	hintInfo := ""
	if m.hintLevel < len(ch.Hints) {
		hintInfo = fmt.Sprintf(" · h 查看提示 (%d/%d)", m.hintLevel, len(ch.Hints))
	} else if len(ch.Hints) > 0 {
		hintInfo = fmt.Sprintf(" · 已显示全部提示 (%d/%d)", m.hintLevel, len(ch.Hints))
	}
	b.WriteString(HelpStyle.Render("Enter 开始挑战 · Esc 返回" + hintInfo))

	boxWidth := responsiveBoxWidth(m.width)
	content := BoxStyle.Width(boxWidth).Render(b.String())
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}
