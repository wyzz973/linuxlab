package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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

	var body strings.Builder

	// Metadata
	body.WriteString(fmt.Sprintf("难度  %s\n", DifficultyStars(ch.Difficulty)))

	if len(ch.Tags) > 0 {
		body.WriteString(fmt.Sprintf("标签  %s\n", DimStyle.Render(strings.Join(ch.Tags, ", "))))
	}

	if ch.RequiresDocker && !m.dockerAvailable {
		body.WriteString("\n")
		body.WriteString(WarningStyle.Render("!! 需要 Docker (当前不可用，将使用本地模式)"))
		body.WriteString("\n")
	}

	// Task description section
	body.WriteString("\n")
	body.WriteString(sectionTitle("任务描述", m.width))
	body.WriteString("\n\n")
	body.WriteString(strings.ReplaceAll(ch.Description, "\n", "\n"))
	body.WriteString("\n")

	// Hints section
	body.WriteString("\n")
	body.WriteString(sectionTitle(fmt.Sprintf("提示 (%d/%d)", m.hintLevel, len(ch.Hints)), m.width))
	body.WriteString("\n\n")

	if m.hintLevel > 0 {
		for i := 0; i < m.hintLevel && i < len(ch.Hints); i++ {
			body.WriteString(fmt.Sprintf("%d. %s\n", i+1, ch.Hints[i].Text))
		}
	} else if len(ch.Hints) > 0 {
		body.WriteString(DimStyle.Render("按 h 解锁下一条提示"))
		body.WriteString("\n")
	}

	box := contentBox(ch.Title, body.String(), m.width, m.height, "")

	// Status bar
	rightHelp := "Enter 开始挑战 · h 查看提示 · Esc 返回"
	if m.hintLevel >= len(ch.Hints) && len(ch.Hints) > 0 {
		rightHelp = "Enter 开始挑战 · Esc 返回"
	}
	status := statusBar("", rightHelp, m.width)

	return verticalCenter(box, status, m.height)
}
