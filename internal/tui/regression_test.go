package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sd3/linuxlab/internal/challenge"
)

// --- reference search-first semantics ---

func TestReferenceModel_NavLettersGoIntoQuery(t *testing.T) {
	m := NewReferenceModel(sampleRefs())

	// "grep" starts with 'g' which used to jump the cursor to Home instead of
	// entering the query.
	for _, r := range "grep" {
		m, _ = m.Update(runeKey(r))
	}
	rm := m.(ReferenceModel)
	if rm.query != "grep" {
		t.Fatalf("query = %q, want grep", rm.query)
	}
	if len(rm.filtered) != 1 || rm.filtered[0].Name != "grep" {
		t.Fatalf("expected grep to be found, got %+v", rm.filtered)
	}
}

func TestReferenceModel_SingleNavLetterDoesNotMoveCursor(t *testing.T) {
	m := NewReferenceModel(sampleRefs())
	m, _ = m.Update(runeKey('G')) // used to jump to End
	rm := m.(ReferenceModel)
	if rm.query != "G" {
		t.Fatalf("query = %q, want G", rm.query)
	}
	if rm.cursor != 0 {
		t.Fatalf("cursor = %d, want 0", rm.cursor)
	}
}

func TestReferenceModel_QAppendsToQueryInsteadOfGoingBack(t *testing.T) {
	m := NewReferenceModel(sampleRefs())
	updated, cmd := m.Update(runeKey('q'))
	if cmd != nil {
		t.Fatal("q should append to the query, not send GoBackMsg")
	}
	rm := updated.(ReferenceModel)
	if rm.query != "q" {
		t.Fatalf("query = %q, want q", rm.query)
	}
}

func TestReferenceModel_SpaceAppendsToQuery(t *testing.T) {
	m := NewReferenceModel(sampleRefs())
	m, _ = m.Update(runeKey('l'))
	m, _ = m.Update(runeKey('s'))
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeySpace})
	rm := m.(ReferenceModel)
	if rm.query != "ls " {
		t.Fatalf("query = %q, want %q", rm.query, "ls ")
	}
}

func TestReferenceModel_LeftGoesBackWithEmptyQuery(t *testing.T) {
	m := NewReferenceModel(sampleRefs())
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyLeft})
	if cmd == nil {
		t.Fatal("expected Left to go back when query is empty")
	}
	if _, ok := cmd().(GoBackMsg); !ok {
		t.Fatal("expected GoBackMsg from Left")
	}
}

// --- contentBox tiny width robustness ---

func TestContentBox_TinyWidthDoesNotPanic(t *testing.T) {
	// Width 1 used to hit strings.Repeat with a negative count and panic.
	for width := 0; width <= 10; width++ {
		_ = contentBox("标题", "内容", width, 24, "")
	}
	_ = contentBox("标题", "内容", 1, 24, "1-2/3")
}

func TestMenuModel_ViewAtWidthOneDoesNotPanic(t *testing.T) {
	m := NewMenuModelWithStats(42, 5).(MenuModel)
	m.width = 1
	m.height = 1
	if view := m.View(); view == "" {
		t.Fatal("expected non-empty view at width 1")
	}
}

func TestBoxWidth_HasSafeLowerBound(t *testing.T) {
	for width := -1; width <= 21; width++ {
		if got := boxWidth(width); got < minBoxWidth {
			t.Fatalf("boxWidth(%d) = %d, want >= %d", width, got, minBoxWidth)
		}
	}
}

// --- contentBox long title truncation ---

func TestContentBox_LongTitleTruncatedToBoxWidth(t *testing.T) {
	title := "使用 grep 在大文件中搜索（模拟 less 搜索）"
	box := contentBox(title, "内容\n第二行", 52, 24, "1-20/38")

	lines := strings.Split(box, "\n")
	want := lipgloss.Width(lines[len(lines)-1]) // bottom border defines box width
	for i, line := range lines {
		if got := lipgloss.Width(line); got != want {
			t.Fatalf("line %d width = %d, want %d\n%s", i, got, want, box)
		}
	}
	if !strings.Contains(box, "…") {
		t.Fatal("expected truncated title to end with an ellipsis")
	}
}

func TestTruncateToWidth_CJKAware(t *testing.T) {
	got := truncateToWidth("容器与部署", 6)
	if w := lipgloss.Width(got); w > 6 {
		t.Fatalf("truncated width = %d, want <= 6 (%q)", w, got)
	}
	if !strings.HasSuffix(got, "…") {
		t.Fatalf("expected ellipsis suffix, got %q", got)
	}
	if full := truncateToWidth("abc", 10); full != "abc" {
		t.Fatalf("short string should be unchanged, got %q", full)
	}
}

// --- display-width alignment (CJK labels) ---

func TestPadDisplayWidth(t *testing.T) {
	cases := []string{
		"Vim 操作",
		"Linux 基础命令",
		SelectedStyle.Render("容器与部署"),
		"plain",
	}
	for _, s := range cases {
		if got := lipgloss.Width(padDisplayWidth(s, 18)); got != 18 {
			t.Fatalf("padDisplayWidth(%q, 18) display width = %d, want 18", s, got)
		}
	}
	// Never truncates strings already wider than the target.
	wide := strings.Repeat("宽", 12)
	if got := padDisplayWidth(wide, 18); got != wide {
		t.Fatalf("over-wide string should be unchanged, got %q", got)
	}
}

func TestModulesModel_ViewAlignsCJKLabels(t *testing.T) {
	cats := map[string][]*challenge.Challenge{
		"containers":   {{ID: "c1", Title: "题", Category: "containers", Difficulty: 1}},
		"linux-basics": {{ID: "l1", Title: "题", Category: "linux-basics", Difficulty: 1}},
		"vim":          {{ID: "v1", Title: "题", Category: "vim", Difficulty: 1}},
	}
	store := tmpStore(t)
	m := NewModulesModel(cats, store).(ModulesModel)
	m.width = 100
	m.height = 30

	view := m.View()
	// Progress bars of module rows must start at the same display column even
	// though the CJK labels have different display widths. Compare the two
	// non-cursor rows (the cursor row has a differently padded key style).
	var barCols []int
	for _, line := range strings.Split(view, "\n") {
		if !strings.Contains(line, "Linux 基础命令") && !strings.Contains(line, "Vim 操作") {
			continue
		}
		idx := strings.IndexAny(line, "░█")
		if idx < 0 {
			t.Fatalf("module row missing progress bar: %q", line)
		}
		barCols = append(barCols, lipgloss.Width(line[:idx]))
	}
	if len(barCols) != 2 {
		t.Fatalf("expected 2 non-cursor module rows, got %d", len(barCols))
	}
	if barCols[0] != barCols[1] {
		t.Fatalf("progress bars misaligned: columns %v\n%s", barCols, view)
	}
}

// --- resize broadcast to inactive sub-models ---

func TestAppModel_ResizeWhileInSubScreenUpdatesMenu(t *testing.T) {
	store := tmpStore(t)
	m := NewAppModel(sampleCategories(), store, nil)
	m, _ = m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m, _ = m.Update(MenuChoiceMsg{Choice: "practice"})

	// Resize while the modules screen is active, then go back to the menu.
	m, _ = m.Update(tea.WindowSizeMsg{Width: 60, Height: 20})
	m, _ = m.Update(GoBackMsg{})

	app := m.(AppModel)
	if app.screen != screenMenu {
		t.Fatalf("screen = %d, want screenMenu", app.screen)
	}
	menu := app.menu.(MenuModel)
	if menu.width != 60 || menu.height != 20 {
		t.Fatalf("menu size = %dx%d, want 60x20", menu.width, menu.height)
	}
}

// --- offset re-clamped on resize ---

func TestChallengesModel_ResizeKeepsCursorVisible(t *testing.T) {
	store := tmpStore(t)
	m := NewChallengesModel("linux-basics", manyChallenges(12), store).(ChallengesModel)
	m.height = 20 // 8 visible rows

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnd})
	m = updated.(ChallengesModel)
	if m.cursor != 11 {
		t.Fatalf("cursor = %d, want 11", m.cursor)
	}

	// Shrink the terminal: the old offset would leave the cursor off-screen.
	updated, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 14}) // 5 visible rows
	m = updated.(ChallengesModel)
	maxVis := m.maxVisible()
	if m.cursor < m.offset || m.cursor >= m.offset+maxVis {
		t.Fatalf("cursor %d not visible in window [%d,%d)", m.cursor, m.offset, m.offset+maxVis)
	}
}

// --- retry keeps hints used ---

func TestAppModel_ResultRetryCarriesHintsUsed(t *testing.T) {
	store := tmpStore(t)
	cats := sampleCategories()
	m := NewAppModel(cats, store, nil).(AppModel)
	m.screen = screenResult
	m.currentCat = "linux-basics"
	m.currentChallenge = cats["linux-basics"][0]
	m.lastResult = &ChallengeResultMsg{Passed: false, HintsUsed: 2}

	_, cmd := m.Update(runeKey('r'))
	if cmd == nil {
		t.Fatal("expected r to launch the current challenge")
	}
	launch, ok := cmd().(LaunchChallengeMsg)
	if !ok {
		t.Fatal("expected LaunchChallengeMsg")
	}
	if launch.HintsUsed != 2 {
		t.Fatalf("HintsUsed = %d, want 2", launch.HintsUsed)
	}
}

// --- go-back preserves list state ---

// step feeds a message to the app model and, when the resulting command emits
// a follow-up message, feeds that too (mirroring the bubbletea loop).
func step(t *testing.T, m tea.Model, msg tea.Msg) tea.Model {
	t.Helper()
	m, cmd := m.Update(msg)
	if cmd != nil {
		if next := cmd(); next != nil {
			m, _ = m.Update(next)
		}
	}
	return m
}

func TestAppModel_GoBackFromDetailKeepsChallengeCursor(t *testing.T) {
	store := tmpStore(t)
	m := tea.Model(NewAppModel(sampleCategories(), store, nil))
	m = step(t, m, tea.WindowSizeMsg{Width: 100, Height: 30})
	m = step(t, m, MenuChoiceMsg{Choice: "practice"})
	m = step(t, m, ModuleSelectedMsg{Category: "linux-basics"})
	m = step(t, m, tea.KeyMsg{Type: tea.KeyDown}) // cursor -> 1 (cp-files)
	m = step(t, m, tea.KeyMsg{Type: tea.KeyEnter})

	app := m.(AppModel)
	if app.screen != screenDetail {
		t.Fatalf("screen = %d, want screenDetail", app.screen)
	}

	m = step(t, m, runeKey('q')) // back to the challenge list
	app = m.(AppModel)
	if app.screen != screenChallenges {
		t.Fatalf("screen = %d, want screenChallenges", app.screen)
	}
	cm := app.challenges.(ChallengesModel)
	if cm.cursor != 1 {
		t.Fatalf("cursor = %d, want 1 (position should survive go-back)", cm.cursor)
	}
}

func TestAppModel_GoBackFromChallengesKeepsModuleCursor(t *testing.T) {
	store := tmpStore(t)
	m := tea.Model(NewAppModel(sampleCategories(), store, nil))
	m = step(t, m, tea.WindowSizeMsg{Width: 100, Height: 30})
	m = step(t, m, MenuChoiceMsg{Choice: "practice"})
	m = step(t, m, tea.KeyMsg{Type: tea.KeyDown}) // modules cursor -> 1
	m = step(t, m, tea.KeyMsg{Type: tea.KeyEnter})

	app := m.(AppModel)
	if app.screen != screenChallenges {
		t.Fatalf("screen = %d, want screenChallenges", app.screen)
	}

	m = step(t, m, runeKey('q')) // back to the module list
	app = m.(AppModel)
	if app.screen != screenModules {
		t.Fatalf("screen = %d, want screenModules", app.screen)
	}
	mm := app.modules.(ModulesModel)
	if mm.cursor != 1 {
		t.Fatalf("modules cursor = %d, want 1 (position should survive go-back)", mm.cursor)
	}
}

func TestModulesModel_ViewReadsProgressLiveFromStore(t *testing.T) {
	cats := sampleCategories()
	store := tmpStore(t)
	m := NewModulesModel(cats, store).(ModulesModel)
	m.width = 100
	m.height = 30

	if view := m.View(); !strings.Contains(view, "0/2") {
		t.Fatalf("expected 0/2 before any attempt\n%s", view)
	}

	// The same (reused) model must reflect new progress without a rebuild.
	store.RecordAttempt("ls-basic", "linux-basics", "navigation", true, 0)
	if view := m.View(); !strings.Contains(view, "1/2") {
		t.Fatalf("expected 1/2 after passing ls-basic\n%s", view)
	}
}
