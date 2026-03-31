package tui

import (
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sd3/linuxlab/internal/challenge"
	"github.com/sd3/linuxlab/internal/progress"
	"github.com/sd3/linuxlab/internal/verify"
)

type screenID int

const (
	screenMenu screenID = iota
	screenModules
	screenChallenges
	screenDetail
	screenSkillMap
	screenResult
)

// ChallengeResultMsg is sent after a challenge attempt completes.
type ChallengeResultMsg struct {
	Passed   bool
	Results  []verify.Result
	HintsUsed int
}

// AppModel is the root TUI model that manages screen navigation.
type AppModel struct {
	screen         screenID
	categories     map[string][]*challenge.Challenge
	store          *progress.Store
	currentCat     string
	currentChallenge *challenge.Challenge

	menu       tea.Model
	modules    tea.Model
	challenges tea.Model
	detail     tea.Model
	skillmap   tea.Model

	lastResult *ChallengeResultMsg
}

// NewAppModel creates the root app model.
func NewAppModel(cats map[string][]*challenge.Challenge, store *progress.Store) tea.Model {
	return AppModel{
		screen:     screenMenu,
		categories: cats,
		store:      store,
		menu:       NewMenuModel(),
	}
}

func (m AppModel) Init() tea.Cmd {
	return m.menu.Init()
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

	case MenuChoiceMsg:
		switch msg.Choice {
		case "practice":
			m.screen = screenModules
			m.modules = NewModulesModel(m.categories, m.store)
			return m, nil
		case "skillmap":
			m.screen = screenSkillMap
			m.skillmap = NewSkillMapModel(m.store)
			return m, nil
		case "recommend":
			// Flatten challenges for recommendation
			var all []*challenge.Challenge
			for _, chs := range m.categories {
				all = append(all, chs...)
			}
			rec := progress.RecommendWeakest(m.store, all)
			if rec != nil {
				m.screen = screenDetail
				m.currentChallenge = rec
				m.currentCat = rec.Category
				m.detail = NewDetailModel(rec)
			}
			return m, nil
		}
		return m, nil

	case ModuleSelectedMsg:
		m.screen = screenChallenges
		m.currentCat = msg.Category
		m.challenges = NewChallengesModel(msg.Category, m.categories[msg.Category], m.store)
		return m, nil

	case ChallengeSelectedMsg:
		m.screen = screenDetail
		m.currentChallenge = msg.Challenge
		m.detail = NewDetailModel(msg.Challenge)
		return m, nil

	case LaunchChallengeMsg:
		return m, m.launchChallenge(msg)

	case ChallengeResultMsg:
		m.lastResult = &msg
		ch := m.currentChallenge
		if ch != nil {
			m.store.RecordAttempt(ch.ID, ch.Category, ch.Subcategory, msg.Passed, msg.HintsUsed)
			m.store.Save() // best-effort persist
		}
		m.screen = screenResult
		return m, nil

	case GoBackMsg:
		switch m.screen {
		case screenModules:
			m.screen = screenMenu
		case screenChallenges:
			m.screen = screenModules
			m.modules = NewModulesModel(m.categories, m.store)
		case screenDetail:
			m.screen = screenChallenges
			if m.currentCat != "" {
				m.challenges = NewChallengesModel(m.currentCat, m.categories[m.currentCat], m.store)
			}
		case screenSkillMap:
			m.screen = screenMenu
		case screenResult:
			m.screen = screenChallenges
			if m.currentCat != "" {
				m.challenges = NewChallengesModel(m.currentCat, m.categories[m.currentCat], m.store)
			}
		}
		return m, nil
	}

	// Delegate to active sub-model
	var cmd tea.Cmd
	switch m.screen {
	case screenMenu:
		m.menu, cmd = m.menu.Update(msg)
	case screenModules:
		if m.modules != nil {
			m.modules, cmd = m.modules.Update(msg)
		}
	case screenChallenges:
		if m.challenges != nil {
			m.challenges, cmd = m.challenges.Update(msg)
		}
	case screenDetail:
		if m.detail != nil {
			m.detail, cmd = m.detail.Update(msg)
		}
	case screenSkillMap:
		if m.skillmap != nil {
			m.skillmap, cmd = m.skillmap.Update(msg)
		}
	}
	return m, cmd
}

func (m AppModel) View() string {
	switch m.screen {
	case screenMenu:
		return m.menu.View()
	case screenModules:
		if m.modules != nil {
			return m.modules.View()
		}
	case screenChallenges:
		if m.challenges != nil {
			return m.challenges.View()
		}
	case screenDetail:
		if m.detail != nil {
			return m.detail.View()
		}
	case screenSkillMap:
		if m.skillmap != nil {
			return m.skillmap.View()
		}
	case screenResult:
		return m.resultView()
	}
	return ""
}

func (m AppModel) launchChallenge(msg LaunchChallengeMsg) tea.Cmd {
	ch := msg.Challenge
	hintsUsed := msg.HintsUsed

	if ch.Category == "vim" {
		// Vim challenge - launch vim
		var filePath string
		if len(ch.SetupFiles) > 0 {
			filePath = ch.SetupFiles[0].Path
		}
		c := exec.Command("vim", filePath)
		return tea.ExecProcess(c, func(err error) tea.Msg {
			results := verify.RunAll(ch.Verify)
			passed := verify.AllPassed(results)
			return ChallengeResultMsg{Passed: passed, Results: results, HintsUsed: hintsUsed}
		})
	}

	// Docker challenge - use interactive shell
	// containerID would be set up elsewhere; for now use a placeholder
	c := exec.Command("docker", "exec", "-it", "linuxlab-sandbox", "/bin/bash")
	return tea.ExecProcess(c, func(err error) tea.Msg {
		results := verify.RunAll(ch.Verify)
		passed := verify.AllPassed(results)
		return ChallengeResultMsg{Passed: passed, Results: results, HintsUsed: hintsUsed}
	})
}

func (m AppModel) resultView() string {
	var b strings.Builder

	if m.lastResult == nil {
		return BoxStyle.Render("无结果")
	}

	if m.lastResult.Passed {
		b.WriteString(TitleStyle.Render("  " + PassedIcon + " 挑战通过!"))
	} else {
		b.WriteString(ErrorStyle.Render("  " + FailedIcon + " 挑战未通过"))
	}
	b.WriteString("\n\n")

	for i, r := range m.lastResult.Results {
		icon := PassedIcon
		if !r.Passed {
			icon = FailedIcon
		}
		b.WriteString(fmt.Sprintf("  %s 检查 %d: %s\n", icon, i+1, r.Message))
	}

	if m.lastResult.HintsUsed > 0 {
		b.WriteString(fmt.Sprintf("\n使用提示: %d\n", m.lastResult.HintsUsed))
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("Esc 返回列表 · q 退出"))

	return BoxStyle.Render(b.String())
}
