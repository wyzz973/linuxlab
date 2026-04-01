package tui

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sd3/linuxlab/internal/challenge"
	"github.com/sd3/linuxlab/internal/progress"
	"github.com/sd3/linuxlab/internal/reference"
	"github.com/sd3/linuxlab/internal/sandbox"
	"github.com/sd3/linuxlab/internal/verify"
)

type screenID int

const (
	screenMenu screenID = iota
	screenModules
	screenChallenges
	screenDetail
	screenSkillMap
	screenRecommend
	screenReference
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
	refs           *reference.ReferenceData
	currentCat     string
	currentChallenge *challenge.Challenge

	width  int
	height int

	menu       tea.Model
	modules    tea.Model
	challenges tea.Model
	detail     tea.Model
	skillmap   tea.Model
	recommend  tea.Model
	refModel   tea.Model

	lastResult *ChallengeResultMsg
}

// NewAppModel creates the root app model.
func NewAppModel(cats map[string][]*challenge.Challenge, store *progress.Store, refs *reference.ReferenceData) tea.Model {
	return AppModel{
		screen:     screenMenu,
		categories: cats,
		store:      store,
		refs:       refs,
		menu:       NewMenuModel(),
	}
}

func (m AppModel) Init() tea.Cmd {
	return m.menu.Init()
}

// sizeModel sends a WindowSizeMsg to a sub-model so it knows the terminal dimensions.
func (m AppModel) sizeModel(sub tea.Model) tea.Model {
	sized, _ := sub.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
	return sized
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Pass to active sub-model
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
		case screenRecommend:
			if m.recommend != nil {
				m.recommend, cmd = m.recommend.Update(msg)
			}
		case screenReference:
			if m.refModel != nil {
				m.refModel, cmd = m.refModel.Update(msg)
			}
		}
		return m, cmd

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

	case MenuChoiceMsg:
		switch msg.Choice {
		case "practice":
			m.screen = screenModules
			m.modules = m.sizeModel(NewModulesModel(m.categories, m.store))
			return m, nil
		case "skillmap":
			m.screen = screenSkillMap
			m.skillmap = m.sizeModel(NewSkillMapModel(m.store))
			return m, nil
		case "recommend":
			var all []*challenge.Challenge
			for _, chs := range m.categories {
				all = append(all, chs...)
			}
			recs := progress.RecommendMultiple(m.store, all, 5)
			m.screen = screenRecommend
			m.recommend = m.sizeModel(NewRecommendModel(recs))
			return m, nil
		case "reference":
			if m.refs != nil {
				m.screen = screenReference
				m.refModel = m.sizeModel(NewReferenceModel(m.refs))
				return m, nil
			}
			return m, nil
		}
		return m, nil

	case ModuleSelectedMsg:
		m.screen = screenChallenges
		m.currentCat = msg.Category
		m.challenges = m.sizeModel(NewChallengesModel(msg.Category, m.categories[msg.Category], m.store))
		return m, nil

	case ChallengeSelectedMsg:
		m.screen = screenDetail
		m.currentChallenge = msg.Challenge
		m.detail = m.sizeModel(NewDetailModel(msg.Challenge))
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
			m.modules = m.sizeModel(NewModulesModel(m.categories, m.store))
		case screenDetail:
			m.screen = screenChallenges
			if m.currentCat != "" {
				m.challenges = m.sizeModel(NewChallengesModel(m.currentCat, m.categories[m.currentCat], m.store))
			}
		case screenSkillMap:
			m.screen = screenMenu
		case screenRecommend:
			m.screen = screenMenu
		case screenReference:
			m.screen = screenMenu
		case screenResult:
			m.screen = screenChallenges
			if m.currentCat != "" {
				m.challenges = m.sizeModel(NewChallengesModel(m.currentCat, m.categories[m.currentCat], m.store))
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
	case screenRecommend:
		if m.recommend != nil {
			m.recommend, cmd = m.recommend.Update(msg)
		}
	case screenReference:
		if m.refModel != nil {
			m.refModel, cmd = m.refModel.Update(msg)
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
	case screenRecommend:
		if m.recommend != nil {
			return m.recommend.View()
		}
	case screenReference:
		if m.refModel != nil {
			return m.refModel.View()
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

	// Create sandbox via factory (Docker, Compose, or Local fallback)
	ctx := context.Background()
	sb, err := sandbox.NewSandbox(ctx, ch)
	if err != nil {
		return func() tea.Msg {
			return ChallengeResultMsg{
				Passed:    false,
				Results:   []verify.Result{{Passed: false, Message: fmt.Sprintf("sandbox error: %v", err)}},
				HintsUsed: hintsUsed,
			}
		}
	}

	args := sb.InteractiveShellArgs()
	c := exec.Command(args[0], args[1:]...)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		defer sb.Destroy(ctx)
		results := verify.RunAll(ch.Verify)
		passed := verify.AllPassed(results)
		return ChallengeResultMsg{Passed: passed, Results: results, HintsUsed: hintsUsed}
	})
}

func (m AppModel) resultView() string {
	if m.lastResult == nil {
		header := headerView("LinuxLab · 检测结果", m.width)
		footer := footerView("Esc 返回列表 · q 退出", m.width)
		contentHeight := maxInt(1, m.height-2)
		content := fillContent("无结果", m.width, contentHeight)
		return header + "\n" + content + "\n" + footer
	}

	header := headerView("LinuxLab · 检测结果", m.width)
	footer := footerView("Esc 返回列表 · q 退出", m.width)

	var b strings.Builder

	if m.lastResult.Passed {
		b.WriteString(TitleStyle.Render(PassedIcon + " 挑战通过!"))
	} else {
		b.WriteString(ErrorStyle.Render(FailedIcon + " 挑战未通过"))
	}
	b.WriteString("\n\n")

	for i, r := range m.lastResult.Results {
		icon := PassedIcon
		if !r.Passed {
			icon = FailedIcon
		}
		b.WriteString(fmt.Sprintf("%s 检查 %d: %s\n", icon, i+1, r.Message))
	}

	if m.lastResult.HintsUsed > 0 {
		b.WriteString(fmt.Sprintf("\n使用提示: %d\n", m.lastResult.HintsUsed))
	}

	contentHeight := maxInt(1, m.height-2)
	content := fillContent(b.String(), m.width, contentHeight)

	return header + "\n" + content + "\n" + footer
}
