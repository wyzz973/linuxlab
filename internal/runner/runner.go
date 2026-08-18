package runner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/sd3/linuxlab/internal/challenge"
	"github.com/sd3/linuxlab/internal/reference"
	"github.com/sd3/linuxlab/internal/sandbox"
	"github.com/sd3/linuxlab/internal/verify"
)

type Options struct {
	Challenge *challenge.Challenge
	Refs      *reference.ReferenceData
	HintsUsed int
	Emit      func(Event)
	Stdin     io.Reader
	Stdout    io.Writer
	Stderr    io.Writer
}

type Result struct {
	Passed    bool
	Results   []verify.Result
	HintsUsed int
}

type Command struct {
	Ctx    context.Context
	Opts   Options
	Result Result
}

func (c *Command) SetStdin(reader io.Reader)  { c.Opts.Stdin = reader }
func (c *Command) SetStdout(writer io.Writer) { c.Opts.Stdout = writer }
func (c *Command) SetStderr(writer io.Writer) { c.Opts.Stderr = writer }

func (c *Command) Run() error {
	result, err := RunInteractive(c.Ctx, c.Opts)
	c.Result = result
	return err
}

func RunInteractive(ctx context.Context, opts Options) (Result, error) {
	if opts.Challenge == nil {
		return Result{}, fmt.Errorf("题目为空")
	}
	if opts.Challenge.Category == "vim" {
		return runVim(ctx, opts)
	}
	return runSandbox(ctx, opts)
}

func runVim(ctx context.Context, opts Options) (Result, error) {
	ch := opts.Challenge
	if len(ch.SetupFiles) == 0 {
		return Result{
			Passed: false,
			Results: []verify.Result{{
				Passed:  false,
				Message: "题目缺少 setup_files 配置",
			}},
			HintsUsed: opts.HintsUsed,
		}, nil
	}

	emit(opts, Event{Type: "setup", Message: "准备 Vim 练习文件"})
	// 每次运行使用独立的临时目录，避免多个实例互相覆盖练习文件，
	// 结束时（含错误路径）自动清理。
	workDir, err := os.MkdirTemp("", "linuxlab-vim-*")
	if err != nil {
		return Result{}, fmt.Errorf("创建练习目录失败: %w", err)
	}
	defer os.RemoveAll(workDir)

	vimRunner := &sandbox.VimRunner{WorkDir: workDir}
	paths, err := vimRunner.PrepareFiles(ch.SetupFiles)
	if err != nil {
		return Result{}, fmt.Errorf("准备文件失败: %w", err)
	}

	rules, unmatched := resolveVimRules(ch, workDir, paths)
	if unmatched != "" {
		return Result{
			Passed: false,
			Results: []verify.Result{{
				Passed:  false,
				Message: fmt.Sprintf("验证规则路径 %s 未匹配任何练习文件，题目配置有误", unmatched),
			}},
			HintsUsed: opts.HintsUsed,
		}, nil
	}

	emit(opts, Event{Type: "handoff", Mode: "vim", Message: "进入 Vim，保存退出后自动检测"})
	cmd := exec.CommandContext(ctx, "vim", paths[0])
	cmd.Stdin = stdin(opts)
	cmd.Stdout = stdout(opts)
	cmd.Stderr = stderr(opts)
	if err := cmd.Run(); err != nil {
		// 编辑器的非零退出码（如 :cq）不算致命错误，仍继续检测；
		// 只有无法启动等非 ExitError 才向上返回。
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			return Result{}, err
		}
	}

	// 检查脚本以练习目录为工作目录执行，脚本内的相对路径
	//（如 cat challenge.txt）指向被编辑的文件。
	results := verify.RunAllWithOptions(rules, verify.Options{WorkDir: workDir})
	passed := verify.AllPassed(results)
	result := Result{Passed: passed, Results: results, HintsUsed: opts.HintsUsed}
	emit(opts, Event{Type: "result", Passed: &passed, HintsUsed: opts.HintsUsed, Results: results})
	return result, nil
}

// resolveVimRules rewrites relative rule paths for a vim challenge: script
// rules resolve against the challenge directory, file rules map to the setup
// file with the exact same base name inside workDir. It returns the resolved
// rules and the path of the first file rule that matched no setup file
// (empty string when every rule resolved).
func resolveVimRules(ch *challenge.Challenge, workDir string, paths []string) ([]challenge.VerifyRule, string) {
	rules := make([]challenge.VerifyRule, len(ch.Verify))
	copy(rules, ch.Verify)
	for i := range rules {
		if rules[i].Path == "" || filepath.IsAbs(rules[i].Path) {
			continue
		}
		if rules[i].Type == "script" {
			rules[i].Path = filepath.Join(ch.Dir, rules[i].Path)
			continue
		}
		matched := false
		want := filepath.Join(workDir, rules[i].Path)
		for _, path := range paths {
			if path == want || filepath.Base(path) == rules[i].Path {
				rules[i].Path = path
				matched = true
				break
			}
		}
		if !matched {
			return nil, rules[i].Path
		}
	}
	return rules, ""
}

func runSandbox(ctx context.Context, opts Options) (Result, error) {
	ch := opts.Challenge
	emit(opts, Event{Type: "setup", Message: fmt.Sprintf("准备挑战 %s", ch.ID)})

	sb, err := sandbox.NewSandbox(ctx, ch)
	if err != nil {
		return Result{}, fmt.Errorf("沙盒创建失败: %w", err)
	}
	defer func() {
		// 清理必须使用独立的 context：调用方 ctx 可能已被取消或超时
		//（如 Ctrl+C 触发 signal.NotifyContext），复用它会导致 Destroy
		// 立即失败、容器泄漏。
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if derr := sb.Destroy(cleanupCtx); derr != nil {
			emit(opts, Event{Type: "cleanup", Message: fmt.Sprintf("清理沙盒失败: %v", derr)})
		}
	}()

	// Compose 挑战操作的是宿主级资源（docker compose up、宿主 docker CLI、
	// 宿主 /tmp 路径），其 init/check 脚本必须在宿主的题目目录里执行，而
	// 不是在服务容器内（容器里没有 docker 二进制）。交互 shell 同样是宿主
	// shell（以题目目录为工作目录），bashrc 注入也落在宿主上。
	isCompose := false
	execInChallenge := sb.Exec
	if _, isCompose = sb.(*sandbox.ComposeSandbox); isCompose {
		execInChallenge = func(ctx context.Context, command string) (string, int, error) {
			return runHostShell(ctx, ch.Dir, command)
		}
	}

	initScript := filepath.Join(ch.Dir, "init.sh")
	if data, err := os.ReadFile(initScript); err == nil {
		if isCompose {
			// 宿主脚本以文件方式执行：sh -c 会把 $0 置为 shell 路径
			//（如 /bin/sh），导致脚本内 `cd "$(dirname "$0")"` 解析错误。
			runHostShellFile(ctx, ch.Dir, initScript)
		} else {
			execInChallenge(ctx, string(data))
		}
	}

	bashrc := buildBashrc(ch, opts.Refs)
	execInChallenge(ctx, "cat > /tmp/.linuxlab_bashrc << 'LINUXLAB_EOF'\n"+bashrc+"\nLINUXLAB_EOF")

	mode := sandboxMode(sb)
	emit(opts, Event{Type: "handoff", Mode: mode, Message: "进入挑战环境，退出后自动检测"})
	args := sb.InteractiveShellArgs()
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdin = stdin(opts)
	cmd.Stdout = stdout(opts)
	cmd.Stderr = stderr(opts)
	if isCompose {
		cmd.Dir = ch.Dir
	}
	if err := cmd.Run(); err != nil {
		// 交互 shell 的非零退出码（如最后一条命令失败后 exit）不算致命
		// 错误，仍继续执行验证；只有无法启动等非 ExitError 才向上返回。
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			return Result{}, err
		}
	}

	results := verifyInSandbox(ctx, ch, execInChallenge)
	passed := verify.AllPassed(results)
	result := Result{Passed: passed, Results: results, HintsUsed: opts.HintsUsed}
	emit(opts, Event{Type: "result", Passed: &passed, HintsUsed: opts.HintsUsed, Results: results})
	return result, nil
}

// runHostShell executes a script on the host with the given working directory.
// Used for ComposeSandbox init/check scripts, which manage host-level docker
// resources and files.
func runHostShell(ctx context.Context, dir, script string) (string, int, error) {
	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", script)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return string(out), exitErr.ExitCode(), nil
		}
		return string(out), -1, err
	}
	return string(out), 0, nil
}

// runHostShellFile executes a script FILE on the host with the given working
// directory, so $0 resolves to the script path and `cd "$(dirname "$0")"`
// works as intended (unlike `sh -c`, where $0 is the shell path). The path is
// absolutized because a relative path would be re-resolved against cmd.Dir
// (which is itself relative, e.g. challenges/<category>/<id>) and double the
// directory prefix.
func runHostShellFile(ctx context.Context, dir, scriptPath string) (string, int, error) {
	abs, err := filepath.Abs(scriptPath)
	if err != nil {
		return "", -1, err
	}
	cmd := exec.CommandContext(ctx, "bash", abs)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return string(out), exitErr.ExitCode(), nil
		}
		return string(out), -1, err
	}
	return string(out), 0, nil
}

func verifyInSandbox(ctx context.Context, ch *challenge.Challenge, execFn func(ctx context.Context, command string) (string, int, error)) []verify.Result {
	// 空规则集视为题目配置错误，避免什么都没验证就判过。
	if len(ch.Verify) == 0 {
		return []verify.Result{{Passed: false, Message: "题目缺少验证规则，无法判定完成情况"}}
	}
	var results []verify.Result
	for _, rule := range ch.Verify {
		switch rule.Type {
		case "file_content":
			out, _, execErr := execFn(ctx, "cat '"+rule.Path+"' 2>/dev/null")
			if execErr != nil {
				results = append(results, verify.Result{Passed: false, Message: fmt.Sprintf("无法读取文件 %s: %v", rule.Path, execErr)})
				continue
			}
			actual := strings.TrimSpace(out)
			expected := strings.TrimSpace(rule.Expect)
			if actual == expected {
				results = append(results, verify.Result{Passed: true, Message: "文件内容匹配"})
			} else {
				results = append(results, verify.Result{Passed: false, Message: fmt.Sprintf("文件内容不匹配\n期望: %s\n实际: %s", expected, actual)})
			}
		case "file_exists":
			_, code, _ := execFn(ctx, "test -e '"+rule.Path+"'")
			if code == 0 {
				results = append(results, verify.Result{Passed: true, Message: fmt.Sprintf("路径存在: %s", rule.Path)})
			} else {
				results = append(results, verify.Result{Passed: false, Message: fmt.Sprintf("路径不存在: %s", rule.Path)})
			}
		case "command_output":
			out, _, execErr := execFn(ctx, rule.Command)
			if execErr != nil {
				results = append(results, verify.Result{Passed: false, Message: fmt.Sprintf("命令执行失败: %v", execErr)})
				continue
			}
			actual := strings.TrimSpace(out)
			expected := strings.TrimSpace(rule.Expect)
			if actual == expected {
				results = append(results, verify.Result{Passed: true, Message: "命令输出匹配"})
			} else {
				results = append(results, verify.Result{Passed: false, Message: fmt.Sprintf("命令输出不匹配\n期望: %s\n实际: %s", expected, actual)})
			}
		case "exit_code":
			_, code, _ := execFn(ctx, rule.Command)
			expected := strings.TrimSpace(rule.Expect)
			if fmt.Sprintf("%d", code) == expected {
				results = append(results, verify.Result{Passed: true, Message: fmt.Sprintf("退出码匹配: %s", expected)})
			} else {
				results = append(results, verify.Result{Passed: false, Message: fmt.Sprintf("退出码不匹配\n期望: %s\n实际: %d", expected, code)})
			}
		case "permissions":
			out, _, _ := execFn(ctx, "stat -c '%a' '"+rule.Path+"' 2>/dev/null || stat -f '%Lp' '"+rule.Path+"' 2>/dev/null")
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
				out, code, _ := execFn(ctx, string(scriptData))
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
	return results
}

func buildBashrc(ch *challenge.Challenge, refs *reference.ReferenceData) string {
	escapeShell := func(s string) string {
		return strings.ReplaceAll(s, "'", "'\"'\"'")
	}

	hintsArray := ""
	for i, h := range ch.Hints {
		hintsArray += fmt.Sprintf("_HINTS[%d]='%s'\n", i, escapeShell(h.Text))
	}

	learnContent := ""
	if refs != nil {
		tagSet := make(map[string]bool)
		for _, t := range ch.Tags {
			tagSet[strings.ToLower(t)] = true
		}
		for _, cmd := range refs.Commands {
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

	desc := strings.TrimSpace(ch.Description)
	return fmt.Sprintf(`
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

export PS1='\[\033[1;34m\][linuxlab]\[\033[0m\] \w\$ '
task
`,
		escapeShell(ch.Title),
		escapeShell(desc),
		len(ch.Hints),
		hintsArray,
		learnContent,
	)
}

func sandboxMode(sb sandbox.Sandbox) string {
	switch sb.(type) {
	case *sandbox.DockerSandbox:
		return "docker"
	case *sandbox.ComposeSandbox:
		return "compose"
	case *sandbox.LocalSandbox:
		return "local"
	default:
		return "unknown"
	}
}

func emit(opts Options, event Event) {
	if opts.Emit != nil {
		opts.Emit(event)
	}
}

func stdin(opts Options) io.Reader {
	if opts.Stdin != nil {
		return opts.Stdin
	}
	return os.Stdin
}

func stdout(opts Options) io.Writer {
	if opts.Stdout != nil {
		return opts.Stdout
	}
	return os.Stdout
}

func stderr(opts Options) io.Writer {
	if opts.Stderr != nil {
		return opts.Stderr
	}
	return os.Stderr
}
