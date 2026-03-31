package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sd3/linuxlab/internal/challenge"
	"github.com/sd3/linuxlab/internal/progress"
	"github.com/sd3/linuxlab/internal/tui"
)

func main() {
	// Determine challenges directory
	challengesDir := "challenges"
	if envDir := os.Getenv("LINUXLAB_CHALLENGES"); envDir != "" {
		challengesDir = envDir
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

	// Load progress
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

	store, err := progress.NewStore(progressPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载进度失败: %v\n", err)
		os.Exit(1)
	}

	// Start TUI
	app := tui.NewAppModel(byCategory, store)
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "运行失败: %v\n", err)
		os.Exit(1)
	}
}
