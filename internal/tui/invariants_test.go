package tui

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sd3/linuxlab/internal/challenge"
	"github.com/sd3/linuxlab/internal/progress"
	"github.com/sd3/linuxlab/internal/verify"
)

// Rendering invariants (spec §1.2/§1.3): at every supported terminal size,
// every line of AppModel.View() is at most the terminal width and the number
// of lines equals the terminal height exactly. Checked for every screen at
// four sizes with CJK fixtures (including an overlong CJK title).

var invariantSizes = []struct{ w, h int }{
	{60, 20},
	{80, 24},
	{100, 30},
	{140, 40},
}

var invariantScreens = []struct {
	name   string
	target screenID
}{
	{"menu", screenMenu},
	{"modules", screenModules},
	{"challenges", screenChallenges},
	{"detail", screenDetail},
	{"skillmap", screenSkillMap},
	{"recommend", screenRecommend},
	{"reference", screenReference},
	{"result", screenResult},
}

func invariantCategories() map[string][]*challenge.Challenge {
	chs := make([]*challenge.Challenge, 0, 13)
	for i := 0; i < 12; i++ {
		chs = append(chs, &challenge.Challenge{
			ID:          fmt.Sprintf("inv-%02d", i+1),
			Title:       fmt.Sprintf("挑战 %02d", i+1),
			Category:    "linux-basics",
			Subcategory: "导航与文件",
			Difficulty:  (i % 5) + 1,
			Description: "学习使用 ls 命令列出目录内容，并理解权限位的含义。\n第二行说明。",
			Hints:       []challenge.Hint{{Level: 1, Text: "试试 ls -la，注意隐藏文件与权限列的区别。"}},
		})
	}
	chs = append(chs, &challenge.Challenge{
		ID:          "inv-long",
		Title:       "使用 grep 在超大日志文件中搜索指定关键字并统计出现次数（超长标题用于截断测试）",
		Category:    "linux-basics",
		Subcategory: "文本处理",
		Difficulty:  3,
		Description: strings.Repeat("这一行用于制造较长的任务说明，检验换行与滚动在各尺寸下是否稳定。", 6),
		Hints:       []challenge.Hint{{Level: 1, Text: "grep -c 关键字 文件名"}},
	})
	return map[string][]*challenge.Challenge{
		"linux-basics": chs,
		"vim": {{
			ID: "vim-move", Title: "Vim 光标移动", Category: "vim",
			Subcategory: "移动", Difficulty: 1, Description: "使用 hjkl 移动光标。",
		}},
	}
}

func invariantStore(t *testing.T) *progress.Store {
	t.Helper()
	store := tmpStore(t)
	store.RecordAttempt("inv-01", "linux-basics", "导航与文件", true, 0)
	store.RecordAttempt("inv-02", "linux-basics", "导航与文件", false, 1)
	store.RecordAttempt("vim-move", "vim", "移动", false, 2)
	return store
}

// buildAppAt navigates a fresh app model to the given screen at the given
// terminal size. Commands returned along the way are intentionally not run:
// only the state machine matters here.
func buildAppAt(t *testing.T, target screenID, w, h int) AppModel {
	t.Helper()
	cats := invariantCategories()
	m := tea.Model(NewAppModel(cats, invariantStore(t), sampleRefs()))
	feed := func(msg tea.Msg) { m, _ = m.Update(msg) }

	feed(tea.WindowSizeMsg{Width: w, Height: h})
	switch target {
	case screenModules:
		feed(MenuChoiceMsg{Choice: "practice"})
	case screenChallenges:
		feed(MenuChoiceMsg{Choice: "practice"})
		feed(ModuleSelectedMsg{Category: "linux-basics"})
	case screenDetail:
		feed(MenuChoiceMsg{Choice: "practice"})
		feed(ModuleSelectedMsg{Category: "linux-basics"})
		feed(ChallengeSelectedMsg{Challenge: cats["linux-basics"][12]}) // overlong CJK title
	case screenSkillMap:
		feed(MenuChoiceMsg{Choice: "skillmap"})
	case screenRecommend:
		feed(MenuChoiceMsg{Choice: "recommend"})
	case screenReference:
		feed(MenuChoiceMsg{Choice: "reference"})
	case screenResult:
		feed(MenuChoiceMsg{Choice: "practice"})
		feed(ModuleSelectedMsg{Category: "linux-basics"})
		feed(ChallengeSelectedMsg{Challenge: cats["linux-basics"][0]})
		feed(ChallengeResultMsg{
			Passed:    false,
			Results:   []verify.Result{{Passed: true, Message: "文件已创建"}, {Passed: false, Message: "权限不正确，期望 644"}},
			HintsUsed: 1,
		})
	}

	app := m.(AppModel)
	if app.screen != target {
		t.Fatalf("screen = %d, want %d", app.screen, target)
	}
	return app
}

func assertFrameInvariants(t *testing.T, view string, w, h int) {
	t.Helper()
	lines := strings.Split(view, "\n")
	if len(lines) != h {
		t.Fatalf("view has %d lines, want exactly %d\n%s", len(lines), h, view)
	}
	for i, line := range lines {
		if got := lipgloss.Width(line); got > w {
			t.Fatalf("line %d width = %d > terminal width %d: %q", i, got, w, line)
		}
	}
}

func TestViewInvariants_AllScreensAllSizes(t *testing.T) {
	for _, sc := range invariantScreens {
		for _, size := range invariantSizes {
			t.Run(fmt.Sprintf("%s_%dx%d", sc.name, size.w, size.h), func(t *testing.T) {
				app := buildAppAt(t, sc.target, size.w, size.h)
				assertFrameInvariants(t, app.View(), size.w, size.h)
			})
		}
	}
}

func TestViewInvariants_FrameStableAcrossRenders(t *testing.T) {
	// Same state must render byte-for-byte identically (spec §1.5).
	app := buildAppAt(t, screenChallenges, 100, 30)
	first := app.View()
	for i := 0; i < 3; i++ {
		if got := app.View(); got != first {
			t.Fatalf("render %d differs from first render", i+1)
		}
	}
}

func TestViewInvariants_UnsupportedSize(t *testing.T) {
	app := buildAppAt(t, screenMenu, 100, 30)
	m, _ := app.Update(tea.WindowSizeMsg{Width: 59, Height: 17})
	small := m.(AppModel)
	view := small.View()
	assertFrameInvariants(t, view, 59, 17)
	if !strings.Contains(view, "终端窗口太小") {
		t.Fatalf("unsupported view missing notice:\n%s", view)
	}
}

func TestViewInvariants_FullHelpKeepsFrame(t *testing.T) {
	app := buildAppAt(t, screenChallenges, 80, 24)
	m, _ := app.Update(runeKey('?'))
	toggled := m.(AppModel)
	if !toggled.help.ShowAll {
		t.Fatal("expected ? to toggle full help")
	}
	assertFrameInvariants(t, toggled.View(), 80, 24)

	m, _ = toggled.Update(runeKey('?'))
	back := m.(AppModel)
	if back.help.ShowAll {
		t.Fatal("expected second ? to hide full help")
	}
}

// --- ready gate -------------------------------------------------------------

func TestAppModel_ReadyGate_EmptyViewBeforeFirstSize(t *testing.T) {
	m := NewAppModel(invariantCategories(), tmpStore(t), sampleRefs())
	if view := m.(AppModel).View(); view != "" {
		t.Fatalf("expected empty view before the first WindowSizeMsg, got %q", view)
	}
}

// --- Ctrl+C double press ----------------------------------------------------

func TestAppModel_CtrlCFirstPressArmsSecondQuits(t *testing.T) {
	app := buildAppAt(t, screenChallenges, 80, 24)

	m, cmd := app.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	armed := m.(AppModel)
	if !armed.ctrlCArmed {
		t.Fatal("first Ctrl+C should arm the quit")
	}
	if armed.notice != ctrlCNotice {
		t.Fatalf("notice = %q, want %q", armed.notice, ctrlCNotice)
	}
	if cmd == nil {
		t.Fatal("first Ctrl+C should schedule the disarm tick")
	}
	if !strings.Contains(armed.View(), "再按一次") {
		t.Fatal("footer should show the confirm-quit notice")
	}

	_, cmd = armed.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatal("second Ctrl+C should quit")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatalf("expected tea.QuitMsg from second Ctrl+C, got %T", cmd())
	}
}

func TestAppModel_CtrlCOtherKeyDisarms(t *testing.T) {
	app := buildAppAt(t, screenChallenges, 80, 24)
	m, _ := app.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	disarmed := m.(AppModel)
	if disarmed.ctrlCArmed {
		t.Fatal("any other key should disarm the pending quit")
	}
	if disarmed.notice != "" {
		t.Fatalf("notice should be cleared, got %q", disarmed.notice)
	}
}

func TestAppModel_CtrlCTimeoutDisarms(t *testing.T) {
	app := buildAppAt(t, screenChallenges, 80, 24)
	m, _ := app.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	armed := m.(AppModel)
	m, _ = armed.Update(noticeExpireMsg{seq: armed.noticeSeq})
	expired := m.(AppModel)
	if expired.ctrlCArmed {
		t.Fatal("expired tick should disarm")
	}
	if expired.notice != "" {
		t.Fatalf("notice should be cleared, got %q", expired.notice)
	}
}

func TestAppModel_StaleNoticeExpireIgnored(t *testing.T) {
	app := buildAppAt(t, screenChallenges, 80, 24)
	m, _ := app.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	armed := m.(AppModel)
	m, _ = armed.Update(noticeExpireMsg{seq: armed.noticeSeq - 1})
	still := m.(AppModel)
	if !still.ctrlCArmed || still.notice == "" {
		t.Fatal("a stale expire message must not clear a newer notice")
	}
}

// --- route stack ------------------------------------------------------------

func TestAppModel_BackFromDetailReturnsToRecommend(t *testing.T) {
	// Entering a challenge from the recommend list must return to the
	// recommend list, not to the challenges screen (route stack semantics).
	cats := invariantCategories()
	m := tea.Model(NewAppModel(cats, invariantStore(t), sampleRefs()))
	m, _ = m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m, _ = m.Update(MenuChoiceMsg{Choice: "recommend"})
	m, _ = m.Update(ChallengeSelectedMsg{Challenge: cats["linux-basics"][0]})

	app := m.(AppModel)
	if app.screen != screenDetail {
		t.Fatalf("screen = %d, want screenDetail", app.screen)
	}

	m, _ = m.Update(GoBackMsg{})
	app = m.(AppModel)
	if app.screen != screenRecommend {
		t.Fatalf("back from recommend-entered detail landed on %d, want screenRecommend", app.screen)
	}

	m, _ = m.Update(GoBackMsg{})
	app = m.(AppModel)
	if app.screen != screenMenu {
		t.Fatalf("back from recommend landed on %d, want screenMenu", app.screen)
	}
}

func TestAppModel_GoBackOnEmptyStackFallsBackToMenu(t *testing.T) {
	app := buildAppAt(t, screenMenu, 80, 24)
	app.screen = screenResult // simulate an orphan screen with no history
	m, _ := app.Update(GoBackMsg{})
	if got := m.(AppModel).screen; got != screenMenu {
		t.Fatalf("screen = %d, want screenMenu", got)
	}
}
