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

	// Escape single quotes in strings for shell embedding
	escapeShell := func(s string) string {
		return strings.ReplaceAll(s, "'", "'\"'\"'")
	}

	// Build hints array for the hint command
	hintsArray := ""
	for i, h := range ch.Hints {
		hintsArray += fmt.Sprintf("_HINTS[%d]='%s'\n", i, escapeShell(h.Text))
	}

	// Build learn content from reference data matched to challenge tags
	learnContent := ""
	if m.refs != nil {
		tagSet := make(map[string]bool)
		for _, t := range ch.Tags {
			tagSet[strings.ToLower(t)] = true
		}
		for _, cmd := range m.refs.Commands {
			if !tagSet[strings.ToLower(cmd.Name)] {
				continue
			}
			learnContent += fmt.Sprintf(`
  echo ""
  echo -e "\033[1;36m  ┌─ %s ─────────────────────────────────────\033[0m"
  echo -e "\033[1;36m  │\033[0m  %s"
  echo -e "\033[1;36m  │\033[0m"
`, escapeShell(cmd.Name), escapeShell(cmd.Brief))
			for _, ex := range cmd.Examples {
				learnContent += fmt.Sprintf(`  echo -e "\033[1;36m  │\033[0m  \033[2m%s:\033[0m"
  echo -e "\033[1;36m  │\033[0m    \033[1;32m$ %s\033[0m"
  echo -e "\033[1;36m  │\033[0m"
`, escapeShell(ex.Desc), escapeShell(ex.Cmd))
			}
			learnContent += `  echo -e "\033[1;36m  └──────────────────────────────────────────\033[0m"` + "\n"
		}
	}

	// Inject helper commands (.bashrc) into the sandbox
	desc := strings.TrimSpace(ch.Description)
	bashrc := fmt.Sprintf(`
# LinuxLab challenge helpers
_TITLE='%s'
_DESC='%s'
_HINT_COUNT=%d
_HINT_SHOWN=0
%s

task() {
  echo ""
  echo -e "\033[1;34m══════════════════════════════════════════════════════\033[0m"
  echo -e "\033[1;34m  $_TITLE\033[0m"
  echo -e "\033[1;34m══════════════════════════════════════════════════════\033[0m"
  echo ""
  echo -e "  $_DESC" | sed 's/^/  /'
  echo ""
  echo -e "\033[2m  learn 命令详解 · hint 提示 · help 帮助 · exit 完成\033[0m"
  echo -e "\033[1;34m══════════════════════════════════════════════════════\033[0m"
  echo ""
}

hint() {
  if [ $_HINT_COUNT -eq 0 ]; then
    echo -e "\033[33m  本题没有提示\033[0m"
    return
  fi
  if [ $_HINT_SHOWN -ge $_HINT_COUNT ]; then
    echo -e "\033[33m  已显示全部提示 ($_HINT_COUNT/$_HINT_COUNT)\033[0m"
    echo ""
    for i in $(seq 0 $((_HINT_COUNT-1))); do
      echo -e "\033[33m  $((i+1)). ${_HINTS[$i]}\033[0m"
    done
    return
  fi
  echo -e "\033[33m  提示 $((_HINT_SHOWN+1))/$_HINT_COUNT: ${_HINTS[$_HINT_SHOWN]}\033[0m"
  _HINT_SHOWN=$((_HINT_SHOWN+1))
}

learn() {
  echo ""
  echo -e "\033[1;36m  ═══ 本题涉及的命令详解 ═══\033[0m"
%s
  if [ -z "$1" ]; then
    echo ""
    echo -e "\033[2m  掌握了吗？输入 task 回顾任务，hint 查看提示\033[0m"
  fi
  echo ""
}

help() {
  echo ""
  echo -e "\033[1m  可用命令:\033[0m"
  echo -e "    \033[1;34mtask\033[0m      查看任务描述"
  echo -e "    \033[1;36mlearn\033[0m     查看本题涉及的命令详解和示例"
  echo -e "    \033[1;33mhint\033[0m      查看下一条提示 (共 $_HINT_COUNT 条)"
  echo -e "    \033[1mhelp\033[0m      显示此帮助"
  echo -e "    \033[1mexit\033[0m      完成挑战并检测结果"
  echo ""
}

# Custom prompt
export PS1='\[\033[1;34m\][linuxlab]\[\033[0m\] \w\$ '

# Show task on entry
task
`,
		escapeShell(ch.Title),
		escapeShell(desc),
		len(ch.Hints),
		hintsArray,
		learnContent,
	)

	// Write .bashrc into the sandbox
	sb.Exec(ctx, "cat > /tmp/.linuxlab_bashrc << 'LINUXLAB_EOF'\n"+bashrc+"\nLINUXLAB_EOF")

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
	var body strings.Builder

	if m.lastResult == nil {
		body.WriteString(DimStyle.Render("无结果"))
		body.WriteString("\n")

		box := contentBox("检测结果", body.String(), m.width, m.height, "")
		status := statusBar("", "Esc 返回 · q 退出", m.width)
		return verticalCenter(box, status, m.height)
	}

	if m.lastResult.Passed {
		body.WriteString(SuccessStyle.Render("● 挑战通过！"))
	} else {
		body.WriteString(ErrorStyle.Render("● 挑战未通过"))
	}
	body.WriteString("\n\n")

	for i, r := range m.lastResult.Results {
		icon := PassedIcon
		if !r.Passed {
			icon = FailedIcon
		}
		body.WriteString(fmt.Sprintf("%s  检查 %d: %s\n", icon, i+1, r.Message))
	}

	if m.lastResult.HintsUsed > 0 {
		body.WriteString(fmt.Sprintf("\n%s\n", DimStyle.Render(fmt.Sprintf("使用提示: %d", m.lastResult.HintsUsed))))
	}

	box := contentBox("检测结果", body.String(), m.width, m.height, "")
	status := statusBar("", "Esc 返回 · q 退出", m.width)

	return verticalCenter(box, status, m.height)
}
