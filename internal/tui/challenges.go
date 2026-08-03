package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sd3/linuxlab/internal/challenge"
	"github.com/sd3/linuxlab/internal/progress"
)

// ChallengeSelectedMsg is sent when the user selects a challenge.
type ChallengeSelectedMsg struct {
	Challenge *challenge.Challenge
}

// chFilterCursor is the caret glyph shown at the end of the filter query.
// Its width is always measured through displayWidth, never assumed.
const chFilterCursor = "▌"

// --- shared list-row style sets ---------------------------------------------

// chRowStyles groups the pre-built styles used to paint one list row of the
// challenges / recommend screens. The focus variant carries the BgFocus
// background on every segment — including spacer runs, icons and stars —
// because wrapping an already-styled row in a single background style would
// lose the background at each inner ANSI reset.
type chRowStyles struct {
	CursorSlot string // "› " (focused) or two plain spaces
	Text       lipgloss.Style
	Dim        lipgloss.Style
	Success    lipgloss.Style
	Error      lipgloss.Style
	Category   lipgloss.Style
	Space      lipgloss.Style // spacer runs (carries the row background)
	IconPass   string
	IconFail   string
	IconTodo   string
	starOn     string
	starOff    string
}

// newChRowStyles builds a row style set from the theme. Called once at
// startup (iron rule: no lipgloss.NewStyle() inside View).
func newChRowStyles(t Theme, focused bool) chRowStyles {
	base := lipgloss.NewStyle()
	if focused {
		base = base.Background(t.BgFocus)
	}
	text := base.Foreground(t.FgBase)
	if focused {
		text = text.Bold(true)
	}
	rs := chRowStyles{
		Text:     text,
		Dim:      base.Foreground(t.FgSubtle),
		Success:  base.Foreground(t.Success),
		Error:    base.Foreground(t.Error),
		Category: base.Foreground(t.Secondary),
		Space:    base,
		IconPass: base.Foreground(t.Success).Render(IconPass),
		IconFail: base.Foreground(t.Error).Render(IconFail),
		IconTodo: base.Foreground(t.FgSubtle).Render(IconTodo),
		starOn:   base.Foreground(t.Warning).Render(StarOn),
		starOff:  base.Foreground(t.FgSubtle).Render(StarOff),
	}
	if focused {
		rs.CursorSlot = base.Foreground(t.Primary).Bold(true).Render(IconCursor + " ")
	} else {
		rs.CursorSlot = "  "
	}
	return rs
}

var (
	chRowNormal = newChRowStyles(currentTheme, false)
	chRowFocus  = newChRowStyles(currentTheme, true)
)

// stars renders the 5-glyph difficulty gauge in this row's style set.
func (rs chRowStyles) stars(level int) string {
	if level < 1 {
		level = 1
	}
	if level > 5 {
		level = 5
	}
	var b strings.Builder
	for i := 0; i < 5; i++ {
		if i < level {
			b.WriteString(rs.starOn)
		} else {
			b.WriteString(rs.starOff)
		}
	}
	return b.String()
}

// pad renders n spacer columns in this row's style set (so a focused row's
// background covers the gaps between segments too).
func (rs chRowStyles) pad(n int) string {
	if n <= 0 {
		return ""
	}
	return rs.Space.Render(strings.Repeat(" ", n))
}

// chInnerWidth mirrors contentBox's inner content width (box width minus
// border and padding) so rows can be laid out against the real budget.
func chInnerWidth(termWidth int) int {
	return maxInt(boxWidth(termWidth)-6, 0)
}

const (
	// Subcategory column bounds for the middle of a row.
	chSubcolMin = 6
	chSubcolMax = 18

	// The preview is a floating panel drawn over the list (opencode's overlay
	// model) rather than a second column: the list keeps its own layout and
	// logic no matter whether the preview is open.
	chPreviewMinW = 28
	chPreviewMaxW = 44

	// chListFirstRow is the line index of the first list row inside the
	// rendered content box: top border, padding row, progress line, separator.
	chListFirstRow = 4
)

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// --- model -------------------------------------------------------------------

// ChallengesModel is the challenge list screen for a category. Besides the
// plain list it has a '/' filter mode: typing filters by title + subcategory
// substring as you type, Esc leaves the filter and restores the full list,
// and Enter opens the highlighted challenge directly.
type ChallengesModel struct {
	category   string
	challenges []*challenge.Challenge
	store      *progress.Store
	cursor     int
	offset     int
	width      int
	height     int

	// '/' filter mode state. While filtering, cursor/offset index into
	// filtered (indices into challenges, in list order), not the full list.
	filtering bool
	query     string
	filtered  []int

	// Floating preview panel, drawn over the list without changing its layout.
	preview panelAnim
}

// NewChallengesModel creates a new challenge list screen.
func NewChallengesModel(category string, challenges []*challenge.Challenge, store *progress.Store) tea.Model {
	return ChallengesModel{
		category:   category,
		challenges: challenges,
		store:      store,
		cursor:     0,
	}
}

func (m ChallengesModel) Init() tea.Cmd { return nil }

// focusChallenge moves the cursor onto the challenge with the given ID (if
// present) and keeps it inside the visible window. The scroll position is
// otherwise preserved. Coming back from the detail screen always lands on the
// full list, so any active filter is dropped first.
func (m ChallengesModel) focusChallenge(id string) ChallengesModel {
	if m.filtering {
		m.filtering = false
		m.query = ""
		m.filtered = nil
		if m.cursor >= len(m.challenges) {
			m.cursor = 0
		}
	}
	for i, ch := range m.challenges {
		if ch.ID == id {
			m.cursor = i
			break
		}
	}
	m.offset = ensureVisible(m.cursor, m.offset, len(m.challenges), m.maxVisible())
	return m
}

func (m ChallengesModel) maxVisible() int {
	return visibleListRows(m.height)
}

// listLen returns the number of rows in the active view (filtered or full).
func (m ChallengesModel) listLen() int {
	if m.filtering {
		return len(m.filtered)
	}
	return len(m.challenges)
}

// itemAt resolves the i-th row of the active view to its challenge.
func (m ChallengesModel) itemAt(i int) *challenge.Challenge {
	if m.filtering {
		return m.challenges[m.filtered[i]]
	}
	return m.challenges[i]
}

func (m ChallengesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Re-clamp scroll so the cursor stays inside the resized viewport.
		m.offset = ensureVisible(m.cursor, m.offset, m.listLen(), m.maxVisible())
		return m, nil
	case panelTickMsg:
		var cmd tea.Cmd
		m.preview, cmd = m.preview.advance()
		return m, cmd
	case tea.KeyMsg:
		if m.filtering {
			return m.updateFilter(msg)
		}
		switch {
		case runeMatches(msg, "p"):
			return m.togglePreview()
		case isBackKey(msg):
			// Esc closes the preview first, then leaves the screen.
			if m.preview.opening() {
				return m.togglePreview()
			}
			return m, func() tea.Msg { return GoBackMsg{} }
		case runeMatches(msg, "/"):
			m.enterFilter()
			return m, nil
		case isSelectKey(msg):
			if len(m.challenges) == 0 {
				return m, nil
			}
			ch := m.challenges[m.cursor]
			return m, func() tea.Msg { return ChallengeSelectedMsg{Challenge: ch} }
		default:
			list := ScrollList{Cursor: m.cursor, Offset: m.offset}
			if list.Move(msg, len(m.challenges), m.maxVisible(), true) {
				m.cursor = list.Cursor
				m.offset = list.Offset
			}
		}
	}
	return m, nil
}

// updateFilter handles keys while the '/' filter is active. Query input has
// priority: printable runes (and space) always go into the query, so titles
// containing j/k/g/q remain typeable. Navigation is limited to arrow keys,
// PgUp/PgDn and Home/End; Esc leaves the filter, Enter opens the highlighted
// challenge.
func (m ChallengesModel) updateFilter(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.Type == tea.KeyEsc:
		m.exitFilter()
		return m, nil

	case isSelectKey(msg):
		if len(m.filtered) == 0 || m.cursor >= len(m.filtered) {
			return m, nil
		}
		ch := m.challenges[m.filtered[m.cursor]]
		// Commit: leave the filter so going back lands on the full list with
		// the cursor on the challenge just opened (locked behavior id=19).
		m.exitFilter()
		return m, func() tea.Msg { return ChallengeSelectedMsg{Challenge: ch} }

	case msg.Type == tea.KeyCtrlU:
		if m.query != "" {
			m.setQuery("")
		}

	case msg.Type == tea.KeyBackspace:
		if m.query != "" {
			runes := []rune(m.query)
			m.setQuery(string(runes[:len(runes)-1]))
		}

	case msg.Type == tea.KeyRunes:
		m.setQuery(m.query + string(msg.Runes))

	case msg.Type == tea.KeySpace:
		m.setQuery(m.query + " ")

	default:
		list := ScrollList{Cursor: m.cursor, Offset: m.offset}
		if list.Move(msg, len(m.filtered), m.maxVisible(), false) {
			m.cursor = list.Cursor
			m.offset = list.Offset
		}
	}
	return m, nil
}

// enterFilter switches into filter mode with an empty query (the full list
// stays visible, so cursor and scroll positions carry over unchanged).
func (m *ChallengesModel) enterFilter() {
	m.filtering = true
	m.query = ""
	m.filtered = make([]int, len(m.challenges))
	for i := range m.challenges {
		m.filtered[i] = i
	}
}

// exitFilter leaves filter mode and restores the full list. The cursor
// follows the challenge highlighted in the filtered view so the selection
// stays meaningful.
func (m *ChallengesModel) exitFilter() {
	if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
		m.cursor = m.filtered[m.cursor]
	} else if m.cursor >= len(m.challenges) {
		m.cursor = 0
	}
	m.filtering = false
	m.query = ""
	m.filtered = nil
	m.offset = ensureVisible(m.cursor, m.offset, len(m.challenges), m.maxVisible())
}

// setQuery replaces the filter query, recomputes the match set and resets the
// cursor to the first match (filter-as-you-type semantics).
func (m *ChallengesModel) setQuery(query string) {
	m.query = query
	m.filtered = filterChallenges(m.challenges, query)
	m.cursor = 0
	m.offset = 0
}

// filterChallenges returns the indices of the challenges whose title or
// subcategory contains the query (case-insensitive substring match; CJK text
// needs no fuzzying).
func filterChallenges(chs []*challenge.Challenge, query string) []int {
	q := strings.ToLower(strings.TrimSpace(query))
	out := make([]int, 0, len(chs))
	for i, ch := range chs {
		if q == "" ||
			strings.Contains(strings.ToLower(ch.Title), q) ||
			strings.Contains(strings.ToLower(ch.Subcategory), q) {
			out = append(out, i)
		}
	}
	return out
}

// passedCount reads the number of passed challenges live from the store, so a
// reused model always shows up-to-date progress.
func (m ChallengesModel) passedCount() int {
	passed := 0
	for _, ch := range m.challenges {
		if entry, exists := m.store.Data.Challenges[ch.ID]; exists && entry.Status == "passed" {
			passed++
		}
	}
	return passed
}

// challengeTitleWidth computes the fixed title column width: the widest title
// (capped) that still leaves room for the cursor, icon, stars and status
// columns inside the box.
func challengeTitleWidth(chs []*challenge.Challenge, innerW int) int {
	maxTitle := 0
	for _, ch := range chs {
		if w := displayWidth(ch.Title); w > maxTitle {
			maxTitle = w
		}
	}
	if maxTitle > 40 {
		maxTitle = 40
	}
	// Fixed overhead: cursor 2 + icon 2 + gap 2 + stars 5 + gap 2 + status 6.
	if budget := innerW - 19; budget < maxTitle {
		maxTitle = budget
	}
	if maxTitle < 4 {
		maxTitle = 4
	}
	return maxTitle
}

// renderRow renders one challenge list row:
//
//	› ✓ 标题…………………………          ★★☆☆☆  已通过
//
// The difficulty and status columns are flush with the row's right edge, so a
// wide list reads as a table instead of trailing off into whitespace. A
// focused row spans the full inner width so its background highlight does too.
func (m ChallengesModel) renderRow(ch *challenge.Challenge, selected bool, titleW, rowW int) string {
	rs := chRowNormal
	if selected {
		rs = chRowFocus
	}

	icon := rs.IconTodo
	status := rs.Dim.Render("未完成")
	if entry, exists := m.store.Data.Challenges[ch.ID]; exists {
		switch entry.Status {
		case "passed":
			icon = rs.IconPass
			status = rs.Success.Render("已通过")
		case "failed":
			icon = rs.IconFail
			status = rs.Error.Render("未通过")
		}
	}

	title := truncateWidth(ch.Title, titleW)

	head := rs.CursorSlot + icon + rs.pad(1) + rs.Text.Render(title) + rs.pad(titleW-displayWidth(title))
	tail := rs.stars(ch.Difficulty) + rs.pad(2) + status

	// Spare width becomes a subcategory column rather than dead space; the
	// column disappears entirely once the row is too narrow for it.
	spare := rowW - displayWidth(head) - displayWidth(tail)
	mid := ""
	if ch.Subcategory != "" && spare >= chSubcolMin+4 {
		sub := truncateWidth(ch.Subcategory, minInt(spare-4, chSubcolMax))
		mid = rs.pad(2) + rs.Dim.Render(sub)
	}

	// Push the tail to the right edge, keeping at least two columns of
	// separation when the row is too narrow to spread out.
	gap := maxInt(rowW-displayWidth(head)-displayWidth(mid)-displayWidth(tail), 2)
	row := head + mid + rs.pad(gap) + tail

	if selected {
		row += rs.pad(rowW - displayWidth(row))
	}
	return row
}

// filterInputLine renders the query row of the filter mode:
//
//	/ chmod▌   （3/98）
func (m ChallengesModel) filterInputLine(matches, total int) string {
	return S.Selected.Render("/ ") +
		S.Text.Render(m.query) +
		S.Selected.Render(chFilterCursor) +
		"   " +
		S.Dim.Render(fmt.Sprintf("（%d/%d）", matches, total))
}

func (m ChallengesModel) View() string {
	label := CategoryLabel(m.category)
	total := len(m.challenges)

	if total == 0 {
		body := EmptyState("当前模块暂无挑战。", "其他模块还有更多训练内容。", "Esc 返回") + "\n"
		return contentBox(label, body, m.width, m.height, "0/0")
	}

	body, rightLabel := m.listBody(chInnerWidth(m.width))
	base := contentBox(label, body, m.width, m.height, rightLabel)
	return m.withPreviewOverlay(base)
}

// togglePreview opens or closes the floating preview panel.
func (m ChallengesModel) togglePreview() (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.preview, cmd = m.preview.toggle()
	return m, cmd
}

// withPreviewOverlay composites the floating preview onto the rendered list.
// The list itself is rendered exactly as if the preview did not exist.
func (m ChallengesModel) withPreviewOverlay(base string) string {
	if !m.preview.visible() {
		return base
	}
	lines := strings.Split(base, "\n")
	boxW := maxInt(boxWidth(m.width), minBoxWidth)
	panelW := clampInt(boxW/2, chPreviewMinW, chPreviewMaxW)
	if panelW+8 > boxW || len(lines) < 6 {
		return base // no room to float anything meaningful
	}

	full := m.previewLines(maxInt(panelW-6, 0))
	panel := boxLines("预览", full[:m.preview.reveal(len(full))], panelW, "", false)

	// Anchor to the right of the box, next to the selected row.
	leftPad := maxInt((m.width-boxW)/2, 0)
	x := leftPad + boxW - panelW - 2
	y := clampInt(chListFirstRow+(m.cursor-m.offset), 1, maxInt(len(lines)-len(panel)-1, 1))
	return placeOverlay(base, strings.Join(panel, "\n"), x, y)
}

// previewLines describes the challenge under the cursor.
func (m ChallengesModel) previewLines(w int) []string {
	return challengePreviewLines(m.itemAt(m.cursor), m.store, w, m.previewDescLines())
}

// previewDescLines caps the description excerpt so the preview box never grows
// taller than the list box beside it.
func (m ChallengesModel) previewDescLines() int {
	return maxInt(m.maxVisible()-6, 3)
}

// listBody renders the list portion shared by both layouts.
func (m ChallengesModel) listBody(innerW int) (string, string) {
	total := len(m.challenges)
	titleW := challengeTitleWidth(m.challenges, innerW)
	count := m.listLen()

	list := ScrollList{Cursor: m.cursor, Offset: m.offset}
	rows := list.RenderRows(count, m.maxVisible(), innerW, func(i int, selected bool) string {
		return m.renderRow(m.itemAt(i), selected, titleW, innerW)
	})

	var body strings.Builder
	if m.filtering {
		body.WriteString(m.filterInputLine(count, total))
	} else {
		body.WriteString(ProgressSummary("完成进度", m.passedCount(), total, 24))
	}
	body.WriteString("\n")
	body.WriteString(S.Subtle.Render(strings.Repeat("─", innerW)))
	body.WriteString("\n")

	for _, row := range rows {
		body.WriteString(row)
		body.WriteString("\n")
	}
	if m.filtering && count == 0 {
		body.WriteString("\n")
		body.WriteString(EmptyState(fmt.Sprintf("没有匹配「%s」的挑战", m.query), "换个关键词试试。", "Esc 退出过滤"))
		body.WriteString("\n")
	}

	rightLabel := ""
	if m.filtering {
		rightLabel = fmt.Sprintf("%d/%d", count, total)
	} else {
		start, end := list.Window(count, m.maxVisible())
		rightLabel = fmt.Sprintf("%d-%d/%d", start+1, end, total)
	}
	return strings.TrimRight(body.String(), "\n"), rightLabel
}
