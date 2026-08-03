package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sd3/linuxlab/internal/challenge"
)

// previewCats builds a category whose rows have every optional column populated.
func previewCats(n int) map[string][]*challenge.Challenge {
	chs := make([]*challenge.Challenge, 0, n)
	for i := 0; i < n; i++ {
		chs = append(chs, &challenge.Challenge{
			ID:          "c" + string(rune('a'+i%26)) + string(rune('0'+i%10)),
			Title:       "挑战标题 " + strings.Repeat("长", i%3) + string(rune('A'+i%26)),
			Category:    "linux-basics",
			Subcategory: "权限管理",
			Difficulty:  1 + i%5,
			Tags:        []string{"权限", "chmod"},
			Description: "把 /opt/deploy.sh 的权限设置为 750，然后用 ls -l 验证结果。",
		})
	}
	return map[string][]*challenge.Challenge{"linux-basics": chs}
}

func previewChallengesModel(t *testing.T, w, h, items int) ChallengesModel {
	t.Helper()
	cats := previewCats(items)
	m := NewChallengesModel("linux-basics", cats["linux-basics"], tmpStore(t))
	sized, _ := m.Update(tea.WindowSizeMsg{Width: w, Height: h})
	return sized.(ChallengesModel)
}

// Difficulty and status sit flush with the row's right edge at every width, so
// the list reads as a table instead of trailing into whitespace.
func TestChallengesRow_TailIsFlushRight(t *testing.T) {
	for _, rowW := range []int{60, 78, 90} {
		m := previewChallengesModel(t, 140, 34, 8)
		titleW := challengeTitleWidth(m.challenges, rowW)
		for i := range m.challenges {
			row := m.renderRow(m.challenges[i], false, titleW, rowW)
			if got := displayWidth(row); got != rowW {
				t.Fatalf("rowW=%d row %d width = %d, want exactly %d: %q", rowW, i, got, rowW, row)
			}
		}
	}
}

func TestChallengesRow_SubcategoryColumnDropsWhenNarrow(t *testing.T) {
	m := previewChallengesModel(t, 140, 34, 4)
	ch := m.challenges[0]

	wideRow := m.renderRow(ch, false, challengeTitleWidth(m.challenges, 90), 90)
	if !strings.Contains(wideRow, ch.Subcategory) {
		t.Fatalf("wide row should carry the subcategory column: %q", wideRow)
	}

	// Narrow enough that title + difficulty + status already consume the row.
	narrow := 34
	narrowRow := m.renderRow(ch, false, challengeTitleWidth(m.challenges, narrow), narrow)
	if strings.Contains(narrowRow, ch.Subcategory) {
		t.Fatalf("narrow row should drop the subcategory column: %q", narrowRow)
	}
	if got := displayWidth(narrowRow); got != narrow {
		t.Fatalf("narrow row width = %d, want %d", got, narrow)
	}
}

// --- floating preview panel ---------------------------------------------------

// openPreview presses 'p' and runs the reveal animation to completion.
func openPreview(t *testing.T, m ChallengesModel) ChallengesModel {
	t.Helper()
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")})
	if cmd == nil {
		t.Fatal("opening the preview should start the animation")
	}
	m = updated.(ChallengesModel)
	for i := 0; i < panelSteps; i++ {
		next, _ := m.Update(panelTickMsg{})
		m = next.(ChallengesModel)
	}
	if m.preview.step != panelSteps {
		t.Fatalf("previewStep = %d, want %d", m.preview.step, panelSteps)
	}
	return m
}

// The panel is an overlay: opening it must not change a single line of the
// list's own layout outside the panel's rectangle.
func TestChallengesPreview_DoesNotRelayoutTheList(t *testing.T) {
	m := previewChallengesModel(t, 140, 34, 30)
	before := strings.Split(m.View(), "\n")
	after := strings.Split(openPreview(t, m).View(), "\n")

	if len(before) != len(after) {
		t.Fatalf("preview changed the list height: %d -> %d", len(before), len(after))
	}
	changed := 0
	for i := range before {
		if before[i] != after[i] {
			changed++
		}
		if a, b := displayWidth(before[i]), displayWidth(after[i]); a != b {
			t.Fatalf("line %d width changed %d -> %d (the list must keep its layout)", i, a, b)
		}
	}
	if changed == 0 {
		t.Fatal("the preview did not draw anything")
	}
	if changed > panelSteps*6 {
		t.Fatalf("%d lines changed; the panel should only cover its own rectangle", changed)
	}
}

func TestChallengesPreview_AnimatesOpenAndClosed(t *testing.T) {
	m := previewChallengesModel(t, 140, 34, 30)
	if strings.Contains(m.View(), "预览") {
		t.Fatal("the panel should be hidden until it is opened")
	}

	// Opening reveals progressively: each frame draws at least as much as the
	// one before it.
	opening, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")})
	m = opening.(ChallengesModel)
	prev := 0
	for i := 0; i < panelSteps; i++ {
		next, cmd := m.Update(panelTickMsg{})
		m = next.(ChallengesModel)
		height := strings.Count(m.View(), "─╮")
		if m.preview.step < panelSteps && cmd == nil {
			t.Fatalf("frame %d should schedule the next tick", i)
		}
		if m.preview.step < prev {
			t.Fatal("the reveal went backwards")
		}
		prev = m.preview.step
		_ = height
	}
	if !strings.Contains(m.View(), "预览") {
		t.Fatalf("the panel should be visible once open\n%s", m.View())
	}

	// Closing animates back out and then disappears entirely.
	closing, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")})
	m = closing.(ChallengesModel)
	for i := 0; i < panelSteps; i++ {
		next, _ := m.Update(panelTickMsg{})
		m = next.(ChallengesModel)
	}
	if m.preview.step != 0 {
		t.Fatalf("previewStep = %d after closing, want 0", m.preview.step)
	}
	if strings.Contains(m.View(), "预览") {
		t.Fatalf("the panel should be gone after closing\n%s", m.View())
	}
}

// Esc closes the panel before it leaves the screen.
func TestChallengesPreview_EscClosesPanelFirst(t *testing.T) {
	m := openPreview(t, previewChallengesModel(t, 140, 34, 30))

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = updated.(ChallengesModel)
	if m.preview.target != 0 {
		t.Fatal("Esc should start closing the panel")
	}
	if cmd != nil {
		if _, ok := cmd().(GoBackMsg); ok {
			t.Fatal("Esc closed the panel and left the screen in one press")
		}
	}

	for i := 0; i < panelSteps; i++ {
		next, _ := m.Update(panelTickMsg{})
		m = next.(ChallengesModel)
	}
	_, cmd = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("Esc with the panel closed should leave the screen")
	}
	if _, ok := cmd().(GoBackMsg); !ok {
		t.Fatalf("expected GoBackMsg, got %T", cmd())
	}
}

func TestChallengesPreview_DescribesSelectedChallenge(t *testing.T) {
	m := previewChallengesModel(t, 140, 34, 30)
	moved, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = openPreview(t, moved.(ChallengesModel))

	view := m.View()
	for _, want := range []string{"预览", m.challenges[1].Title, "状态", "难度"} {
		if !strings.Contains(view, want) {
			t.Fatalf("panel missing %q\n%s", want, view)
		}
	}
}

// The overlay must never break the frame invariants at any size or frame.
func TestChallengesPreview_HonorsFrameInvariants(t *testing.T) {
	for _, sz := range [][2]int{{60, 18}, {80, 24}, {100, 30}, {140, 34}, {200, 55}} {
		cats := previewCats(60)
		var app tea.Model = NewAppModel(cats, tmpStore(t), sampleRefs())
		app, _ = app.Update(tea.WindowSizeMsg{Width: sz[0], Height: sz[1]})
		app, _ = app.Update(MenuChoiceMsg{Choice: "practice"})
		app, _ = app.Update(ModuleSelectedMsg{Category: "linux-basics"})
		app, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")})

		for frame := 0; frame <= panelSteps; frame++ {
			lines := strings.Split(app.(AppModel).View(), "\n")
			if len(lines) != sz[1] {
				t.Fatalf("%dx%d frame %d: %d lines, want %d", sz[0], sz[1], frame, len(lines), sz[1])
			}
			for i, line := range lines {
				if got := displayWidth(line); got > sz[0] {
					t.Fatalf("%dx%d frame %d: line %d width %d > %d: %q",
						sz[0], sz[1], frame, i, got, sz[0], line)
				}
			}
			app, _ = app.Update(panelTickMsg{})
		}
	}
}
