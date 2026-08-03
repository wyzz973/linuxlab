package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sd3/linuxlab/internal/challenge"
	"github.com/sd3/linuxlab/internal/sandbox"
)

// LaunchChallengeMsg is sent when the user wants to start a challenge.
type LaunchChallengeMsg struct {
	Challenge *challenge.Challenge
	HintsUsed int
}

// Column geometry for the wide (>=120x30) two-column layout: main content box
// on the left, inspector metadata box on the right. The inspector width comes
// from computeLayout (28-40 columns); the main box is capped so prose lines
// stay readable.
const (
	detailInspectorMinW = 28
	detailInspectorMaxW = 44
	detailLabelW        = 6 // inspector label column ("检查项" = 6 cells)
)

// Pre-built private styles (iron rule: no lipgloss.NewStyle() inside View).
// These cover two gaps in the shared Styles set: an accent style for tag/ID
// values (theme Secondary has no ready-made style) and a warning block style
// with a soft background for the Docker notice (theme has no BgWarning token).
var (
	detailAccentStyle = lipgloss.NewStyle().Foreground(currentTheme.Secondary)

	detailWarnBlockStyle = lipgloss.NewStyle().
				Foreground(currentTheme.Warning).
				Bold(true).
				Background(lipgloss.AdaptiveColor{Light: "#f6e8c8", Dark: "#2a2213"})
)

// dockerWarningText is the actionable degraded-mode notice (spec §4.4).
const dockerWarningText = "Docker 未运行，将以本地降级模式执行；启动 Docker 后重新进入本页。"

// DetailModel is the challenge detail screen: viewport-scrolled description
// and progressive hints, plus an inspector column in wide mode.
type DetailModel struct {
	challenge       *challenge.Challenge
	hintLevel       int
	scroll          int // mirror of vp.YOffset (read by locked tests)
	dockerAvailable bool
	width           int
	height          int

	vp viewport.Model

	// Floating metadata panel, drawn over the content box without changing it.
	inspector panelAnim

	// Content cache signature: SetContent runs only when one of these changed
	// (spec §3 — viewport windows rendering, content rebuild stays rare).
	contentReady bool
	contentWrapW int
	contentHints int
	contentDock  bool
	contentTotal int
}

// NewDetailModel creates a new challenge detail screen.
func NewDetailModel(ch *challenge.Challenge) tea.Model {
	return DetailModel{
		challenge:       ch,
		hintLevel:       0,
		dockerAvailable: sandbox.DockerAvailable(),
		vp:              viewport.New(0, 0),
	}
}

func (m DetailModel) Init() tea.Cmd { return nil }

func (m DetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.syncViewport() // re-clamps the scroll offset after a resize
		return m, nil
	case panelTickMsg:
		var cmd tea.Cmd
		m.inspector, cmd = m.inspector.advance()
		return m, cmd
	case tea.KeyMsg:
		switch {
		case runeMatches(msg, "i"):
			var cmd tea.Cmd
			m.inspector, cmd = m.inspector.toggle()
			return m, cmd
		case msg.Type == tea.KeyEnter:
			if m.challenge == nil {
				return m, nil
			}
			ch := m.challenge
			hints := m.hintLevel
			return m, func() tea.Msg {
				return LaunchChallengeMsg{Challenge: ch, HintsUsed: hints}
			}
		case isBackKey(msg):
			// Esc closes the inspector first, then leaves the screen.
			if m.inspector.opening() {
				var cmd tea.Cmd
				m.inspector, cmd = m.inspector.toggle()
				return m, cmd
			}
			return m, func() tea.Msg { return GoBackMsg{} }
		case runeMatches(msg, "h"):
			if m.challenge != nil && m.hintLevel < len(m.challenge.Hints) {
				m.hintLevel++
				m.syncViewport()
				// The new hint is appended at the very bottom of the
				// scrollable content — jump there so unlocking is visible.
				m.vp.GotoBottom()
				m.scroll = m.vp.YOffset
			}
		default:
			if m.applyScrollKey(msg) {
				return m, nil
			}
		}
	}
	return m, nil
}

// applyScrollKey maps the locked scroll key set onto viewport operations.
// Key semantics are unchanged from the pre-viewport implementation:
// ↑/k, ↓/j, PgUp/Ctrl+B, PgDn/Ctrl+F, Home/Ctrl+A/g, End/Ctrl+E/G.
func (m *DetailModel) applyScrollKey(msg tea.KeyMsg) bool {
	m.syncViewport()
	switch {
	case msg.Type == tea.KeyUp || runeMatches(msg, "k"):
		m.vp.ScrollUp(1)
	case msg.Type == tea.KeyDown || runeMatches(msg, "j"):
		m.vp.ScrollDown(1)
	case msg.Type == tea.KeyPgUp || msg.Type == tea.KeyCtrlB:
		m.vp.ViewUp()
	case msg.Type == tea.KeyPgDown || msg.Type == tea.KeyCtrlF:
		m.vp.ViewDown()
	case msg.Type == tea.KeyHome || msg.Type == tea.KeyCtrlA || runeMatches(msg, "g"):
		m.vp.GotoTop()
	case msg.Type == tea.KeyEnd || msg.Type == tea.KeyCtrlE || runeMatches(msg, "G"):
		m.vp.GotoBottom()
	default:
		return false
	}
	m.scroll = m.vp.YOffset
	return true
}

func (m DetailModel) View() string {
	if m.challenge == nil {
		return ""
	}
	m.syncViewport()
	return m.withInspectorOverlay(m.standardView())
}

// --- layout ------------------------------------------------------------------

// mainBoxWidth is the outer width of the content box. The inspector floats
// above it, so the box keeps the same width whether the panel is open or not.
func (m DetailModel) mainBoxWidth() int {
	return maxInt(boxWidth(m.width), minBoxWidth)
}

// inspectorWidth is the floating panel's outer width, or 0 when the box is too
// narrow to host it.
func (m DetailModel) inspectorWidth() int {
	boxW := m.mainBoxWidth()
	w := clampInt(boxW/2, detailInspectorMinW, detailInspectorMaxW)
	if w+8 > boxW {
		return 0
	}
	return w
}

// wrapW is the word-aware wrap width for description and hint prose.
func (m DetailModel) wrapW() int {
	w := m.mainBoxWidth() - 8 // border (2) + padding (4) + breathing room (2)
	if w < 20 {
		return 20
	}
	return w
}

// visibleBudget is the number of scrollable content rows the viewport may use.
// Box chrome takes 4 rows, the app shell 2, plus the pinned hint CTA rows in
// the single-column layout.
func (m DetailModel) visibleBudget() int {
	v := m.height - 8
	if m.ctaLine() != "" {
		v -= 2 // blank + CTA row pinned under the viewport
	}
	if v < 6 {
		return 6
	}
	return v
}

// syncViewport brings the viewport dimensions, content and scroll clamp in
// line with the current model state. SetContent only runs when the content
// signature (wrap width, unlocked hints, docker state) changed.
func (m *DetailModel) syncViewport() {
	if m.challenge == nil {
		return
	}
	wrapW := m.wrapW()
	if !m.contentReady || m.contentWrapW != wrapW || m.contentHints != m.hintLevel ||
		m.contentDock != m.dockerAvailable {
		lines := m.scrollContent(wrapW, false)
		m.vp.SetContent(strings.Join(lines, "\n"))
		m.contentTotal = len(lines)
		m.contentReady = true
		m.contentWrapW = wrapW
		m.contentHints = m.hintLevel
		m.contentDock = m.dockerAvailable
	}
	m.vp.Width = wrapW
	h := m.visibleBudget()
	if h > m.contentTotal {
		h = m.contentTotal // short content keeps the box content-sized
	}
	if h < 1 {
		h = 1
	}
	m.vp.Height = h
	m.vp.SetYOffset(m.vp.YOffset) // re-clamp (locked behavior: resize re-clamps offset)
	m.scroll = m.vp.YOffset
}

// --- content -----------------------------------------------------------------

// scrollContent builds the scrollable body lines. In wide mode the metadata
// line is omitted (it lives in the inspector column instead).
func (m DetailModel) scrollContent(wrapW int, wide bool) []string {
	ch := m.challenge
	var lines []string

	if !wide {
		lines = append(lines, truncateWidth(m.metaLine(), wrapW), "")
	}
	if ch.RequiresDocker && !m.dockerAvailable {
		lines = append(lines, m.dockerWarning(wrapW)...)
		lines = append(lines, "")
	}

	lines = append(lines, detailSectionLine("任务描述", wrapW), "")
	for _, raw := range strings.Split(strings.TrimRight(ch.Description, "\n"), "\n") {
		lines = append(lines, wrapWidth(raw, wrapW)...)
	}

	lines = append(lines, "", detailSectionLine(fmt.Sprintf("提示 %d/%d", m.hintLevel, len(ch.Hints)), wrapW), "")
	switch {
	case len(ch.Hints) == 0:
		lines = append(lines, S.Dim.Render("本挑战没有提示"))
	case m.hintLevel == 0:
		lines = append(lines, S.Dim.Render("尚未解锁提示，按 h 查看第一条"))
	default:
		for i := 0; i < m.hintLevel && i < len(ch.Hints); i++ {
			if i > 0 {
				lines = append(lines, "")
			}
			lines = append(lines, renderHintBlock(ch.Hints[i].Text, wrapW)...)
		}
	}
	return lines
}

// metaLine renders the single-column metadata row:
// 难度 ★★☆☆☆   标签 权限 · chmod   分类 Linux 基础 / 权限管理
func (m DetailModel) metaLine() string {
	ch := m.challenge
	var b strings.Builder
	b.WriteString(S.Dim.Render("难度 "))
	b.WriteString(DifficultyStars(ch.Difficulty))
	if len(ch.Tags) > 0 {
		b.WriteString("   ")
		b.WriteString(S.Dim.Render("标签 "))
		b.WriteString(joinTags(ch.Tags))
	}
	if label := m.categoryLabel(); label != "" {
		b.WriteString("   ")
		b.WriteString(S.Dim.Render("分类 "))
		b.WriteString(S.Text.Render(label))
	}
	return b.String()
}

// joinTags joins tag values in the design's tag rendering: values accented,
// separators dimmed ("权限 · chmod").
func joinTags(tags []string) string {
	parts := make([]string, len(tags))
	for i, t := range tags {
		parts[i] = detailAccentStyle.Render(t)
	}
	return strings.Join(parts, S.Subtle.Render(" · "))
}

func (m DetailModel) categoryLabel() string {
	ch := m.challenge
	label := CategoryLabel(ch.Category)
	if label == "" {
		return ""
	}
	if ch.Subcategory != "" {
		return label + " / " + ch.Subcategory
	}
	return label
}

// dockerWarning renders the degraded-mode notice as full-width warning block
// lines (raw width measured before styling, per the width discipline).
func (m DetailModel) dockerWarning(wrapW int) []string {
	var out []string
	for _, l := range wrapWidth(dockerWarningText, wrapW) {
		out = append(out, detailWarnBlockStyle.Render(padRight(l, wrapW)))
	}
	return out
}

// ctaLine is the pinned hint call-to-action row shown under the viewport in
// the single-column layout (the inspector carries this info in wide mode).
func (m DetailModel) ctaLine() string {
	ch := m.challenge
	if ch == nil || len(ch.Hints) == 0 {
		return ""
	}
	if m.hintLevel < len(ch.Hints) {
		return S.Badge.Render("h") +
			S.Dim.Render(fmt.Sprintf(" 解锁下一条提示（影响得分） · 已用 %d/%d", m.hintLevel, len(ch.Hints)))
	}
	return S.Dim.Render(fmt.Sprintf("已解锁全部提示（%d/%d）", m.hintLevel, len(ch.Hints)))
}

// renderHintBlock wraps one hint with a "▸ " marker and hanging indent. The
// marker's width is measured, never assumed (it is East-Asian-ambiguous).
func renderHintBlock(text string, wrapW int) []string {
	prefix := "▸ "
	prefixW := displayWidth(prefix)
	indent := strings.Repeat(" ", prefixW)
	bodyW := wrapW - prefixW
	if bodyW < 10 {
		bodyW = 10
	}
	var out []string
	for _, raw := range strings.Split(strings.TrimRight(text, "\n"), "\n") {
		for _, l := range wrapWidth(raw, bodyW) {
			if len(out) == 0 {
				out = append(out, S.Warning.Render(prefix+l))
			} else {
				out = append(out, S.Warning.Render(indent+l))
			}
		}
	}
	if len(out) == 0 {
		out = append(out, S.Warning.Render(prefix))
	}
	return out
}

// detailSectionLine renders a section divider at an exact display width:
// ─ 标题 ────────────
func detailSectionLine(title string, w int) string {
	prefix := "─ " + title + " "
	rest := w - displayWidth(prefix)
	if rest < 0 {
		return S.Subtle.Render(truncateWidth(prefix, w))
	}
	return S.Subtle.Render(prefix + strings.Repeat("─", rest))
}

// scrollLabel is the top-border scroll indicator ("1-16/22"), shown only when
// the content overflows the viewport.
func (m DetailModel) scrollLabel() string {
	if m.contentTotal <= m.vp.Height {
		return ""
	}
	start := m.vp.YOffset
	end := start + m.vp.Height
	if end > m.contentTotal {
		end = m.contentTotal
	}
	return fmt.Sprintf("%d-%d/%d", start+1, end, m.contentTotal)
}

// --- rendering ---------------------------------------------------------------

// standardView renders the single-column layout (Compact/Standard): the
// viewport window plus the pinned hint CTA inside one content box.
func (m DetailModel) standardView() string {
	body := m.vp.View()
	if cta := m.ctaLine(); cta != "" {
		body += "\n\n" + cta
	}
	return contentBox(m.challenge.Title, body, m.width, m.height, m.scrollLabel())
}

// inspectorLines builds the inspector rows (label column + value), each
// truncated to the inspector's inner width.
func (m DetailModel) inspectorLines(w int) []string {
	ch := m.challenge
	row := func(label, value string) string {
		return truncateWidth(S.Dim.Render(padRight(label, detailLabelW))+" "+value, w)
	}

	lines := []string{
		row("ID", detailAccentStyle.Render(ch.ID)),
		row("分类", S.Text.Render(CategoryLabel(ch.Category))),
	}
	if ch.Subcategory != "" {
		lines = append(lines, row("子类", S.Text.Render(ch.Subcategory)))
	}
	lines = append(lines, row("难度", DifficultyStars(ch.Difficulty)))

	tagVal := S.Dim.Render("—")
	if len(ch.Tags) > 0 {
		tagVal = joinTags(ch.Tags)
	}
	lines = append(lines, row("标签", tagVal), "")

	hintVal := S.Dim.Render("无")
	if len(ch.Hints) > 0 {
		hintVal = S.Warning.Render(fmt.Sprintf("已用 %d/%d", m.hintLevel, len(ch.Hints)))
	}
	lines = append(lines,
		row("提示", hintVal),
		row("检查项", S.Text.Render(fmt.Sprintf("%d 条", len(ch.Verify)))),
		row("状态", m.sandboxStatus()),
	)
	return lines
}

// sandboxStatus describes how the challenge will run in the current
// environment (the progress store is not reachable from this screen, so the
// row reports sandbox readiness rather than pass/fail history).
func (m DetailModel) sandboxStatus() string {
	if !m.challenge.RequiresDocker {
		return S.Dim.Render("无需 Docker")
	}
	if m.dockerAvailable {
		return S.Success.Render("Docker 就绪")
	}
	return S.Warning.Render("本地降级模式")
}

// withInspectorOverlay composites the floating metadata panel onto the
// rendered content box. The box itself is laid out as if the panel did not
// exist, so opening it never reflows the prose underneath.
func (m DetailModel) withInspectorOverlay(base string) string {
	if !m.inspector.visible() {
		return base
	}
	panelW := m.inspectorWidth()
	if panelW == 0 {
		return base
	}
	lines := strings.Split(base, "\n")
	if len(lines) < 6 {
		return base
	}

	full := m.inspectorLines(maxInt(panelW-6, 0))
	panel := boxLines("元信息", full[:m.inspector.reveal(len(full))], panelW, "", false)

	boxW := m.mainBoxWidth()
	leftPad := maxInt((m.width-boxW)/2, 0)
	x := leftPad + boxW - panelW - 2
	y := clampInt(2, 1, maxInt(len(lines)-len(panel)-1, 1))
	return placeOverlay(base, strings.Join(panel, "\n"), x, y)
}
