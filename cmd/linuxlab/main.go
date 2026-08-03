package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sd3/linuxlab/internal/challenge"
	"github.com/sd3/linuxlab/internal/cli"
	"github.com/sd3/linuxlab/internal/progress"
	"github.com/sd3/linuxlab/internal/reference"
	"github.com/sd3/linuxlab/internal/tui"
)

func main() {
	// Determine challenges directory
	challengesDir := "challenges"
	if envDir := os.Getenv("LINUXLAB_CHALLENGES"); envDir != "" {
		challengesDir = envDir
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "无法获取用户目录: %v\n", err)
		os.Exit(1)
	}
	progressDir := filepath.Join(homeDir, ".linuxlab")
	if err := os.MkdirAll(progressDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "无法创建进度目录: %v\n", err)
		os.Exit(1)
	}
	progressPath := filepath.Join(progressDir, "progress.json")
	refsPath := "references/commands.yaml"

	// CLI path: Ctrl+C / SIGTERM cancel the context so the preparation phase
	// (image pull, init.sh) unwinds through the normal cancellation path and
	// sandbox cleanup can run, instead of the process being killed outright.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	code := cli.Run(ctx, cli.Env{
		Args:          os.Args[1:],
		Stdout:        os.Stdout,
		Stderr:        os.Stderr,
		Stdin:         os.Stdin,
		ChallengesDir: challengesDir,
		RefsPath:      refsPath,
		ProgressPath:  progressPath,
	})
	// Restore default signal behavior before the TUI takes over the terminal,
	// keeping the TUI path unchanged.
	stop()
	if code >= 0 {
		os.Exit(code)
	}

	// Load challenges
	byCategory, err := challenge.LoadAllByCategory(challengesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载题库失败: %v\n", err)
		fmt.Fprintf(os.Stderr, "请确保 %s 目录存在且包含题目\n", challengesDir)
		os.Exit(1)
	}

	if len(byCategory) == 0 {
		fmt.Fprintf(os.Stderr, "题库为空，请在 %s 目录下添加题目\n", challengesDir)
		os.Exit(1)
	}

	store, err := progress.NewStore(progressPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载进度失败: %v\n", err)
		os.Exit(1)
	}

	// Load command references (optional)
	var refs *reference.ReferenceData
	if loadedRefs, err := reference.LoadReferences(refsPath); err == nil {
		refs = loadedRefs
	}

	// Start TUI
	app := tui.NewAppModel(byCategory, store, refs)
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "运行失败: %v\n", err)
		os.Exit(1)
	}
}
