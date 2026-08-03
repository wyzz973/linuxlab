package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sd3/linuxlab/internal/challenge"
	"github.com/sd3/linuxlab/internal/progress"
)

func splashApp(t *testing.T, w, h int) AppModel {
	t.Helper()
	app := buildAppAt(t, screenMenu, w, h)
	app.showSplash = true
	return app
}

// The cover is what the app opens with, and the very first keystroke reveals
// the menu underneath it.
func TestSplash_ShownOnOpenAndDismissedByAnyKey(t *testing.T) {
	var m tea.Model = NewAppModel(sampleCategories(), tmpStore(t), sampleRefs())
	m, _ = m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	app := m.(AppModel)
	if !app.showSplash {
		t.Fatal("the app should open on the cover")
	}
	if view := app.View(); !strings.Contains(view, "按任意键开始") {
		t.Fatalf("cover missing its start hint\n%s", view)
	}
	if app.screen != screenMenu {
		t.Fatalf("screen = %v, want screenMenu (the cover is a layer, not a screen)", app.screen)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
	app = m.(AppModel)
	if app.showSplash {
		t.Fatal("any key should dismiss the cover")
	}
	if view := app.View(); !strings.Contains(view, "LinuxLab 训练控制台") {
		t.Fatalf("menu should be revealed after dismissal\n%s", view)
	}
}

func TestSplash_QuitsOnQ(t *testing.T) {
	app := splashApp(t, 100, 30)
	_, cmd := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd == nil {
		t.Fatal("q on the cover should quit")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatalf("expected tea.QuitMsg, got %T", cmd())
	}
}

// Reaching a screen programmatically must not leave the cover on top of it.
func TestSplash_NavigationMessageDismissesCover(t *testing.T) {
	var m tea.Model = NewAppModel(sampleCategories(), tmpStore(t), sampleRefs())
	m, _ = m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m, _ = m.Update(MenuChoiceMsg{Choice: "practice"})

	app := m.(AppModel)
	if app.showSplash {
		t.Fatal("navigation should dismiss the cover")
	}
	if !strings.Contains(app.View(), "选择训练模块") {
		t.Fatalf("modules screen should render, not the cover\n%s", app.View())
	}
}

// The wordmark is chosen from its own measured width, not a hardcoded
// breakpoint, and the block form is dropped on short terminals.
func TestSplash_WordmarkVariantBySize(t *testing.T) {
	large := artWidth(splashWordmarkLarge)
	if got := pickWordmark(large); len(got) != len(splashWordmarkLarge) {
		t.Fatalf("width %d should fit the block wordmark", large)
	}
	if got := pickWordmark(large - 1); len(got) != 1 {
		t.Fatalf("width %d should fall back to the single-line wordmark", large-1)
	}

	tall := renderSplash(100, 30, splashInfo{})
	if !strings.Contains(tall, "█") {
		t.Fatal("a tall terminal should render the block wordmark")
	}
	short := renderSplash(100, 19, splashInfo{})
	if strings.Contains(short, "█") {
		t.Fatalf("a short terminal should drop the block wordmark\n%s", short)
	}
	if !strings.Contains(short, "L I N U X L A B") {
		t.Fatalf("short cover missing the fallback wordmark\n%s", short)
	}
}

func TestSplash_DockerStateWording(t *testing.T) {
	cases := []struct {
		name  string
		info  splashInfo
		want  string
		avoid string
	}{
		{"探测中", splashInfo{}, "检测中", "未运行"},
		{"就绪", splashInfo{dockerProbed: true, dockerOK: true}, "Docker 就绪", "未运行"},
		{"未运行", splashInfo{dockerProbed: true}, "未运行", "就绪"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			view := renderSplash(100, 30, tc.info)
			if !strings.Contains(view, tc.want) {
				t.Fatalf("cover missing %q\n%s", tc.want, view)
			}
			if strings.Contains(view, tc.avoid) {
				t.Fatalf("cover should not contain %q\n%s", tc.avoid, view)
			}
		})
	}
}

func TestSplash_ShowsOverallProgress(t *testing.T) {
	info := splashInfo{stats: MenuStats{TotalChallenges: 278, TotalModules: 5, PassedChallenges: 12}}
	view := renderSplash(100, 30, info)
	for _, want := range []string{"总体进度", "12/278", "4%"} {
		if !strings.Contains(view, want) {
			t.Fatalf("cover missing %q\n%s", want, view)
		}
	}

	fresh := renderSplash(100, 30, splashInfo{stats: MenuStats{TotalChallenges: 278, TotalModules: 5}})
	if !strings.Contains(fresh, "0/278") {
		t.Fatalf("a fresh install should still show the ratio\n%s", fresh)
	}
	if !strings.Contains(fresh, "欢迎使用") {
		t.Fatalf("a first run should be greeted differently\n%s", fresh)
	}
}

// The cover reports where the session left off and what to pick up next,
// instead of internal details like the progress file's location.
func TestSplash_ShowsSessionPointers(t *testing.T) {
	last := &challenge.Challenge{ID: "chmod-numeric", Title: "使用数字模式修改权限", Difficulty: 2}
	next := &challenge.Challenge{ID: "chmod-recursive", Title: "递归修改目录权限", Difficulty: 3}

	cases := []struct {
		name  string
		entry *progress.ChallengeEntry
		icon  string
		state string
	}{
		{"未通过", &progress.ChallengeEntry{Status: "failed", Attempts: 2, LastAttempt: "2026-07-29"}, IconFail, "未通过"},
		{"已通过", &progress.ChallengeEntry{Status: "passed", Attempts: 1, LastAttempt: "2026-07-29"}, IconPass, "已通过"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			view := renderSplash(100, 30, splashInfo{stats: MenuStats{
				TotalChallenges: 278, PassedChallenges: 12,
				Last: last, LastEntry: tc.entry, Next: next,
			}})
			for _, want := range []string{"上次练到", last.Title, tc.icon, tc.state, "建议下一题", next.Title} {
				if !strings.Contains(view, want) {
					t.Fatalf("cover missing %q\n%s", want, view)
				}
			}
			if tc.entry.Attempts > 1 && !strings.Contains(view, "2 次尝试") {
				t.Fatalf("cover should report the attempt count\n%s", view)
			}
			if !strings.Contains(view, tc.entry.LastAttempt) {
				t.Fatalf("cover should report the date %q\n%s", tc.entry.LastAttempt, view)
			}
			if strings.Contains(view, "progress.json") {
				t.Fatalf("cover should not expose the progress file path\n%s", view)
			}
		})
	}
}

// Nothing attempted yet: the session lines are omitted rather than shown empty.
func TestSplash_OmitsSessionPointersWhenAbsent(t *testing.T) {
	view := renderSplash(100, 30, splashInfo{stats: MenuStats{TotalChallenges: 278}})
	for _, avoid := range []string{"上次练到", "建议下一题"} {
		if strings.Contains(view, avoid) {
			t.Fatalf("cover should omit %q with no history\n%s", avoid, view)
		}
	}
}

func TestSplash_HonorsFrameInvariants(t *testing.T) {
	for _, sz := range [][2]int{{60, 18}, {62, 19}, {80, 24}, {100, 30}, {140, 34}, {200, 55}} {
		app := splashApp(t, sz[0], sz[1])
		lines := strings.Split(app.View(), "\n")
		if len(lines) != sz[1] {
			t.Fatalf("%dx%d: %d lines, want %d", sz[0], sz[1], len(lines), sz[1])
		}
		for i, line := range lines {
			if got := displayWidth(line); got > sz[0] {
				t.Fatalf("%dx%d: line %d width %d > %d: %q", sz[0], sz[1], i, got, sz[0], line)
			}
		}
	}
}

// Renders must be byte-identical across calls so the frame diff stays cheap.
func TestSplash_FrameIsStable(t *testing.T) {
	app := splashApp(t, 100, 30)
	if a, b := app.View(), app.View(); a != b {
		t.Fatal("cover render is not stable across calls")
	}
}
