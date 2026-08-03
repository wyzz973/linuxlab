package tui

// detail_view_test.go — stage-2 tests for the reworked detail screen (spec
// §4.4): bubbles/viewport scrolling, word-aware wrapping via wrapWidth, the
// wide-mode (>=120x30) inspector column, and the Warning-colored actionable
// Docker notice. The locked hint/launch semantics stay in detail_test.go.

import (
	"fmt"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sd3/linuxlab/internal/challenge"
)

// longDetailChallenge returns a challenge whose description far exceeds every
// tested viewport height, so scrolling is actually exercised.
func longDetailChallenge() *challenge.Challenge {
	ch := sampleChallenge()
	ch.Description = strings.TrimRight(
		strings.Repeat("这一行中文说明用于制造需要滚动的长任务描述。\n", 40), "\n")
	return ch
}

// detailAt builds a detail model and feeds it the terminal size.
func detailAt(t *testing.T, ch *challenge.Challenge, w, h int) DetailModel {
	t.Helper()
	m := NewDetailModel(ch).(DetailModel)
	return detailUpdate(t, m, tea.WindowSizeMsg{Width: w, Height: h})
}

func detailUpdate(t *testing.T, m DetailModel, msg tea.Msg) DetailModel {
	t.Helper()
	updated, _ := m.Update(msg)
	next, ok := updated.(DetailModel)
	if !ok {
		t.Fatalf("Update returned %T, want DetailModel", updated)
	}
	return next
}

// --- layout / rendering ------------------------------------------------------

func TestDetailView_StandardLayout(t *testing.T) {
	m := detailAt(t, sampleChallenge(), 100, 30)
	view := m.View()

	for _, want := range []string{
		"ls 基础", // box title
		"难度 ",   // meta row labels
		"标签 ",
		"Linux 基础命令 / navigation", // category / subcategory
		"任务描述",                    // section divider
		"学习使用 ls 命令列出目录内容。", // description body
		"提示 0/3", // hint section header
		"尚未解锁提示，按 h 查看第一条",       // hint empty state
		"解锁下一条提示（影响得分） · 已用 0/3", // pinned CTA
	} {
		if !strings.Contains(view, want) {
			t.Errorf("standard view missing %q\n%s", want, view)
		}
	}
	if strings.Contains(view, "元信息") {
		t.Fatalf("the inspector must stay hidden until it is opened\n%s", view)
	}
}

func TestDetailView_WidthInvariantAcrossSizes(t *testing.T) {
	ch := longDetailChallenge()
	ch.Title = "使用 grep 在超大日志文件中搜索指定关键字并统计出现次数（超长标题用于截断测试）"
	sizes := []struct{ w, h int }{{60, 20}, {80, 24}, {100, 30}, {120, 30}, {140, 34}}
	for _, size := range sizes {
		m := detailAt(t, ch, size.w, size.h)
		for i, line := range strings.Split(m.View(), "\n") {
			if got := lipgloss.Width(line); got > size.w {
				t.Errorf("%dx%d line %d width = %d > %d: %q",
					size.w, size.h, i, got, size.w, line)
			}
		}
	}
}

func TestDetailView_WordAwareWrap(t *testing.T) {
	// At 100 columns the wrap width is 82: the 78-char run fills the first
	// line and the following English word must move to the next line intact
	// instead of being cut mid-word (spec §4.4: wrapWidth is word-aware).
	ch := sampleChallenge()
	ch.Description = strings.Repeat("a", 78) + " journalctl"
	m := detailAt(t, ch, 100, 30)

	found := false
	for _, line := range strings.Split(m.View(), "\n") {
		if strings.Contains(line, "journalctl") {
			found = true
			continue
		}
		if strings.Contains(line, "jour") {
			t.Fatalf("English word split mid-word: %q", line)
		}
	}
	if !found {
		t.Fatal("wrapped description lost the word journalctl")
	}
}

func TestDetailView_FrameStable(t *testing.T) {
	m := detailAt(t, longDetailChallenge(), 100, 30)
	m = detailUpdate(t, m, tea.KeyMsg{Type: tea.KeyDown})
	first := m.View()
	for i := 0; i < 3; i++ {
		if got := m.View(); got != first {
			t.Fatalf("render %d differs from first render", i+1)
		}
	}
}

// --- viewport scrolling ------------------------------------------------------

func TestDetailView_ScrollIndicatorAndKeys(t *testing.T) {
	m := detailAt(t, longDetailChallenge(), 100, 30)
	total, vh := m.contentTotal, m.vp.Height
	if total <= vh {
		t.Fatalf("fixture must overflow the viewport: total=%d visible=%d", total, vh)
	}

	if want := fmt.Sprintf("%d-%d/%d", 1, vh, total); !strings.Contains(m.View(), want) {
		t.Fatalf("view missing scroll indicator %q\n%s", want, m.View())
	}

	// Line down.
	m = detailUpdate(t, m, tea.KeyMsg{Type: tea.KeyDown})
	if m.scroll != 1 {
		t.Fatalf("after ↓ scroll = %d, want 1", m.scroll)
	}
	if want := fmt.Sprintf("%d-%d/%d", 2, vh+1, total); !strings.Contains(m.View(), want) {
		t.Fatalf("view missing scroll indicator %q after ↓", want)
	}

	// Page down from offset 1 advances by one viewport height (clamped).
	m = detailUpdate(t, m, tea.KeyMsg{Type: tea.KeyPgDown})
	want := 1 + vh
	if max := total - vh; want > max {
		want = max
	}
	if m.scroll != want {
		t.Fatalf("after PgDn scroll = %d, want %d", m.scroll, want)
	}

	// G to the bottom, g back to the top (locked key set).
	m = detailUpdate(t, m, runeKey('G'))
	if m.scroll != total-vh {
		t.Fatalf("after G scroll = %d, want %d", m.scroll, total-vh)
	}
	m = detailUpdate(t, m, runeKey('g'))
	if m.scroll != 0 {
		t.Fatalf("after g scroll = %d, want 0", m.scroll)
	}

	// End / Home key variants.
	m = detailUpdate(t, m, tea.KeyMsg{Type: tea.KeyEnd})
	if m.scroll != total-vh {
		t.Fatalf("after End scroll = %d, want %d", m.scroll, total-vh)
	}
	m = detailUpdate(t, m, tea.KeyMsg{Type: tea.KeyHome})
	if m.scroll != 0 {
		t.Fatalf("after Home scroll = %d, want 0", m.scroll)
	}
}

func TestDetailView_ResizeReclampsScroll(t *testing.T) {
	m := detailAt(t, longDetailChallenge(), 100, 30)
	m = detailUpdate(t, m, tea.KeyMsg{Type: tea.KeyEnd})
	before := m.scroll
	if before == 0 {
		t.Fatal("fixture must scroll before the resize")
	}

	// A taller terminal shrinks the max offset; the offset must re-clamp
	// (locked behavior id=18).
	m = detailUpdate(t, m, tea.WindowSizeMsg{Width: 100, Height: 46})
	if maxOff := m.contentTotal - m.vp.Height; m.scroll != maxOff {
		t.Fatalf("after resize scroll = %d, want re-clamped %d", m.scroll, maxOff)
	}
	if m.scroll >= before {
		t.Fatalf("scroll %d not re-clamped below previous %d", m.scroll, before)
	}
}

func TestDetailView_HintUnlockScrollsToBottom(t *testing.T) {
	m := detailAt(t, longDetailChallenge(), 100, 30)
	m = detailUpdate(t, m, runeKey('h'))
	if m.hintLevel != 1 {
		t.Fatalf("hintLevel = %d, want 1", m.hintLevel)
	}
	if maxOff := m.contentTotal - m.vp.Height; m.scroll != maxOff {
		t.Fatalf("h must jump to the new hint at the bottom: scroll = %d, want %d", m.scroll, maxOff)
	}
	view := m.View()
	for _, want := range []string{"提示 1/3", "试试 ls -la", "已用 1/3"} {
		if !strings.Contains(view, want) {
			t.Errorf("view after h missing %q", want)
		}
	}
}

// --- hint CTA states ---------------------------------------------------------

func TestDetailView_AllHintsUnlockedCTA(t *testing.T) {
	m := detailAt(t, sampleChallenge(), 100, 30)
	for i := 0; i < 3; i++ {
		m = detailUpdate(t, m, runeKey('h'))
	}
	view := m.View()
	if !strings.Contains(view, "已解锁全部提示（3/3）") {
		t.Fatalf("view missing exhausted-hints CTA\n%s", view)
	}
	if strings.Contains(view, "解锁下一条提示") {
		t.Fatal("exhausted hints must not keep offering the next hint")
	}
}

func TestDetailView_NoHints(t *testing.T) {
	ch := sampleChallenge()
	ch.Hints = nil
	m := detailAt(t, ch, 100, 30)
	view := m.View()
	if !strings.Contains(view, "本挑战没有提示") {
		t.Fatalf("view missing no-hints notice\n%s", view)
	}
	if strings.Contains(view, "解锁下一条提示") {
		t.Fatal("challenge without hints must not render the hint CTA")
	}
}

// --- Docker warning ----------------------------------------------------------

func TestDetailView_DockerWarningWhenUnavailable(t *testing.T) {
	ch := sampleChallenge()
	ch.RequiresDocker = true
	m := DetailModel{challenge: ch, dockerAvailable: false, vp: viewport.New(0, 0)}
	m = detailUpdate(t, m, tea.WindowSizeMsg{Width: 100, Height: 30})

	view := m.View()
	for _, want := range []string{"Docker 未运行", "重新进入本页"} {
		if !strings.Contains(view, want) {
			t.Errorf("degraded-mode view missing %q\n%s", want, view)
		}
	}
}

func TestDetailView_NoDockerWarningWhenAvailable(t *testing.T) {
	ch := sampleChallenge()
	ch.RequiresDocker = true
	m := DetailModel{challenge: ch, dockerAvailable: true, vp: viewport.New(0, 0)}
	m = detailUpdate(t, m, tea.WindowSizeMsg{Width: 100, Height: 30})

	if view := m.View(); strings.Contains(view, "Docker 未运行") {
		t.Fatalf("warning must not render when Docker is available\n%s", view)
	}
}

// --- floating inspector panel -------------------------------------------------

// openInspector presses 'i' and runs the reveal animation to completion.
func openInspector(t *testing.T, m DetailModel) DetailModel {
	t.Helper()
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("i")})
	if cmd == nil {
		t.Fatal("opening the inspector should start the animation")
	}
	m = updated.(DetailModel)
	for i := 0; i < panelSteps; i++ {
		next, _ := m.Update(panelTickMsg{})
		m = next.(DetailModel)
	}
	return m
}

func TestDetailView_InspectorPanelContent(t *testing.T) {
	ch := sampleChallenge()
	ch.Verify = []challenge.VerifyRule{{Type: "script"}, {Type: "script"}}
	m := openInspector(t, detailAt(t, ch, 140, 34))
	view := m.View()

	for _, want := range []string{
		"元信息",              // panel title
		"ls-basic",         // ID row
		"分类", "Linux 基础命令", // category row
		"子类", "navigation", // subcategory row
		"难度", "标签", // stars + tags rows
		"已用 0/3",     // hints-used row
		"检查项", "2 条", // verify-rule count row
		"状态", // sandbox status row
	} {
		if !strings.Contains(view, want) {
			t.Errorf("inspector panel missing %q\n%s", want, view)
		}
	}
}

// The panel floats: opening it must not reflow the content box underneath.
func TestDetailView_InspectorDoesNotRelayout(t *testing.T) {
	m := detailAt(t, longDetailChallenge(), 140, 34)
	before := strings.Split(m.View(), "\n")
	after := strings.Split(openInspector(t, m).View(), "\n")

	if len(before) != len(after) {
		t.Fatalf("inspector changed the box height: %d -> %d", len(before), len(after))
	}
	for i := range before {
		if a, b := displayWidth(before[i]), displayWidth(after[i]); a != b {
			t.Fatalf("line %d width changed %d -> %d; the box must keep its layout", i, a, b)
		}
	}
	// The pinned CTA belongs to the box and stays put with the panel open.
	if !strings.Contains(strings.Join(after, "\n"), "解锁下一条提示") {
		t.Fatal("the box content should be unaffected by the floating panel")
	}
}

func TestDetailView_InspectorAnimatesAndEscCloses(t *testing.T) {
	m := openInspector(t, detailAt(t, sampleChallenge(), 140, 34))
	if !strings.Contains(m.View(), "元信息") {
		t.Fatal("the panel should be visible once open")
	}

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = updated.(DetailModel)
	if cmd != nil {
		if _, ok := cmd().(GoBackMsg); ok {
			t.Fatal("Esc closed the panel and left the screen in one press")
		}
	}
	for i := 0; i < panelSteps; i++ {
		next, _ := m.Update(panelTickMsg{})
		m = next.(DetailModel)
	}
	if strings.Contains(m.View(), "元信息") {
		t.Fatalf("the panel should be gone after closing\n%s", m.View())
	}

	_, cmd = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("Esc with the panel closed should leave the screen")
	}
	if _, ok := cmd().(GoBackMsg); !ok {
		t.Fatalf("expected GoBackMsg, got %T", cmd())
	}
}

// The panel is available at every supported size, and never overflows.
func TestDetailView_InspectorHonorsFrameInvariants(t *testing.T) {
	for _, sz := range [][2]int{{60, 18}, {80, 24}, {100, 30}, {140, 34}, {200, 55}} {
		m := detailAt(t, longDetailChallenge(), sz[0], sz[1])
		updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("i")})
		m = updated.(DetailModel)
		for frame := 0; frame <= panelSteps; frame++ {
			for i, line := range strings.Split(m.View(), "\n") {
				if got := displayWidth(line); got > sz[0] {
					t.Fatalf("%dx%d frame %d: line %d width %d > %d: %q",
						sz[0], sz[1], frame, i, got, sz[0], line)
				}
			}
			next, _ := m.Update(panelTickMsg{})
			m = next.(DetailModel)
		}
	}
}

func TestDetailView_InspectorSandboxStatusRow(t *testing.T) {
	cases := []struct {
		name     string
		requires bool
		docker   bool
		want     string
	}{
		{"no docker needed", false, false, "无需 Docker"},
		{"docker ready", true, true, "Docker 就绪"},
		{"degraded", true, false, "本地降级模式"},
	}
	for _, tc := range cases {
		ch := sampleChallenge()
		ch.RequiresDocker = tc.requires
		m := DetailModel{challenge: ch, dockerAvailable: tc.docker, vp: viewport.New(0, 0)}
		m = detailUpdate(t, m, tea.WindowSizeMsg{Width: 140, Height: 34})
		if view := openInspector(t, m).View(); !strings.Contains(view, tc.want) {
			t.Errorf("%s: inspector missing status %q\n%s", tc.name, tc.want, view)
		}
	}
}

func TestDetailView_InspectorEnterCarriesHintsUsed(t *testing.T) {
	m := detailAt(t, sampleChallenge(), 140, 34)
	m = detailUpdate(t, m, runeKey('h'))

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected a launch command")
	}
	launch, ok := cmd().(LaunchChallengeMsg)
	if !ok {
		t.Fatalf("expected LaunchChallengeMsg, got %T", cmd())
	}
	if launch.HintsUsed != 1 {
		t.Fatalf("HintsUsed = %d, want 1", launch.HintsUsed)
	}
}
