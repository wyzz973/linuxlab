package tui

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	totalChallenges := 0
	for _, chs := range cats {
		totalChallenges += len(chs)
	}
	return AppModel{
		screen:     screenMenu,
		categories: cats,
		store:      store,
		refs:       refs,
		menu:       NewMenuModelWithStats(totalChallenges, len(cats)),
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
	case screenResult:
		// Handle keys directly for result screen (no sub-model)
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch {
			case keyMsg.Type == tea.KeyEsc || (keyMsg.Type == tea.KeyRunes && string(keyMsg.Runes) == "q"):
				return m, func() tea.Msg { return GoBackMsg{} }
			}
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
		return m.launchVim(ch, hintsUsed)
	}
	return m.launchSandbox(ch, hintsUsed)
}

func (m AppModel) launchVim(ch *challenge.Challenge, hintsUsed int) tea.Cmd {
	if len(ch.SetupFiles) == 0 {
		return func() tea.Msg {
			return ChallengeResultMsg{
				Passed:  false,
				Results: []verify.Result{{Message: "题目缺少 setup_files 配置"}},
			}
		}
	}

	// Prepare files in temp dir
	runner := &sandbox.VimRunner{WorkDir: os.TempDir()}
	paths, err := runner.PrepareFiles(ch.SetupFiles)
	if err != nil {
		return func() tea.Msg {
			return ChallengeResultMsg{Passed: false, Results: []verify.Result{{Message: fmt.Sprintf("准备文件失败: %v", err)}}}
		}
	}

	// Build verify rules with absolute paths
	rules := make([]challenge.VerifyRule, len(ch.Verify))
	copy(rules, ch.Verify)
	for i := range rules {
		if rules[i].Path != "" && len(paths) > 0 {
			// Map relative path to the prepared file path
			for _, p := range paths {
				if strings.HasSuffix(p, rules[i].Path) {
					rules[i].Path = p
					break
				}
			}
		}
	}

	c := exec.Command("vim", paths[0])
	return tea.ExecProcess(c, func(err error) tea.Msg {
		results := verify.RunAll(rules)
		passed := verify.AllPassed(results)
		return ChallengeResultMsg{Passed: passed, Results: results, HintsUsed: hintsUsed}
	})
}

func (m AppModel) launchSandbox(ch *challenge.Challenge, hintsUsed int) tea.Cmd {
	ctx := context.Background()
	sb, err := sandbox.NewSandbox(ctx, ch)
	if err != nil {
		return func() tea.Msg {
			return ChallengeResultMsg{
				Passed:    false,
				Results:   []verify.Result{{Passed: false, Message: fmt.Sprintf("沙盒创建失败: %v", err)}},
				HintsUsed: hintsUsed,
			}
		}
	}

	// Run init.sh if it exists
	initScript := filepath.Join(ch.Dir, "init.sh")
	if data, err := os.ReadFile(initScript); err == nil {
		sb.Exec(ctx, string(data))
	}

	// Build the shell command with a welcome banner
	banner := fmt.Sprintf(
		`echo ""
echo "══════════════════════════════════════════════════════"
echo "  %s"
echo "══════════════════════════════════════════════════════"
echo ""
echo "%s"
echo ""
echo "  输入 exit 完成挑战并检测结果"
echo "══════════════════════════════════════════════════════"
echo ""`,
		ch.Title,
		strings.ReplaceAll(strings.TrimSpace(ch.Description), "\n", "\n  "),
	)

	// Execute banner inside sandbox, then launch interactive shell
	sb.Exec(ctx, banner)

	args := sb.InteractiveShellArgs()
	c := exec.Command(args[0], args[1:]...)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		defer sb.Destroy(ctx)

		// Run verification INSIDE the sandbox, not on the host
		var results []verify.Result
		for _, rule := range ch.Verify {
			switch rule.Type {
			case "file_content":
				out, _, execErr := sb.Exec(ctx, "cat '"+rule.Path+"' 2>/dev/null")
				if execErr != nil {
					results = append(results, verify.Result{Passed: false, Message: fmt.Sprintf("无法读取文件 %s: %v", rule.Path, execErr)})
				} else {
					actual := strings.TrimSpace(out)
					expected := strings.TrimSpace(rule.Expect)
					if actual == expected {
						results = append(results, verify.Result{Passed: true, Message: "文件内容匹配"})
					} else {
						results = append(results, verify.Result{Passed: false, Message: fmt.Sprintf("文件内容不匹配\n期望: %s\n实际: %s", expected, actual)})
					}
				}
			case "file_exists":
				_, code, _ := sb.Exec(ctx, "test -e '"+rule.Path+"'")
				if code == 0 {
					results = append(results, verify.Result{Passed: true, Message: fmt.Sprintf("路径存在: %s", rule.Path)})
				} else {
					results = append(results, verify.Result{Passed: false, Message: fmt.Sprintf("路径不存在: %s", rule.Path)})
				}
			case "command_output":
				out, _, execErr := sb.Exec(ctx, rule.Command)
				if execErr != nil {
					results = append(results, verify.Result{Passed: false, Message: fmt.Sprintf("命令执行失败: %v", execErr)})
				} else {
					actual := strings.TrimSpace(out)
					expected := strings.TrimSpace(rule.Expect)
					if actual == expected {
						results = append(results, verify.Result{Passed: true, Message: "命令输出匹配"})
					} else {
						results = append(results, verify.Result{Passed: false, Message: fmt.Sprintf("命令输出不匹配\n期望: %s\n实际: %s", expected, actual)})
					}
				}
			case "exit_code":
				_, code, _ := sb.Exec(ctx, rule.Command)
				expected := strings.TrimSpace(rule.Expect)
				if fmt.Sprintf("%d", code) == expected {
					results = append(results, verify.Result{Passed: true, Message: fmt.Sprintf("退出码匹配: %s", expected)})
				} else {
					results = append(results, verify.Result{Passed: false, Message: fmt.Sprintf("退出码不匹配\n期望: %s\n实际: %d", expected, code)})
				}
			case "permissions":
				out, _, _ := sb.Exec(ctx, "stat -c '%a' '"+rule.Path+"' 2>/dev/null || stat -f '%Lp' '"+rule.Path+"' 2>/dev/null")
				actual := strings.TrimSpace(out)
				expected := strings.TrimSpace(rule.Expect)
				if actual == expected {
					results = append(results, verify.Result{Passed: true, Message: fmt.Sprintf("权限匹配: %s", expected)})
				} else {
					results = append(results, verify.Result{Passed: false, Message: fmt.Sprintf("权限不匹配\n期望: %s\n实际: %s", expected, actual)})
				}
			case "script":
				scriptPath := rule.Path
				if !filepath.IsAbs(scriptPath) {
					scriptPath = filepath.Join(ch.Dir, scriptPath)
				}
				if scriptData, readErr := os.ReadFile(scriptPath); readErr == nil {
					out, code, _ := sb.Exec(ctx, string(scriptData))
					if code == 0 {
						results = append(results, verify.Result{Passed: true, Message: "脚本检测通过"})
					} else {
						results = append(results, verify.Result{Passed: false, Message: fmt.Sprintf("脚本检测未通过: %s", strings.TrimSpace(out))})
					}
				} else {
					results = append(results, verify.Result{Passed: false, Message: fmt.Sprintf("无法读取脚本: %v", readErr)})
				}
			default:
				results = append(results, verify.Result{Passed: false, Message: fmt.Sprintf("未知验证类型: %s", rule.Type)})
			}
		}

		passed := true
		for _, r := range results {
			if !r.Passed {
				passed = false
				break
			}
		}
		return ChallengeResultMsg{Passed: passed, Results: results, HintsUsed: hintsUsed}
	})
}

func (m AppModel) resultView() string {
	if m.lastResult == nil {
		header := headerView("LinuxLab · 检测结果", m.width)
		footer := footerView("Esc 返回列表 · q 退出", m.width)
		contentHeight := maxInt(1, m.height-6)
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

	contentHeight := maxInt(1, m.height-6)
	content := fillContent(b.String(), m.width, contentHeight)

	return header + "\n" + content + "\n" + footer
}
