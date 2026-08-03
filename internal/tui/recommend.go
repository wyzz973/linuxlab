package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sd3/linuxlab/internal/challenge"
	"github.com/sd3/linuxlab/internal/progress"
)

// RecommendModel shows weakness-based challenge recommendations. Rows are
// rendered through ScrollList with the shared challenge row style set
// (selected row gets the unified BgFocus background); with no practice
// history it shows an EmptyState that points at the next action.
type RecommendModel struct {
	challenges []*challenge.Challenge
	store      *progress.Store
	cursor     int
	offset     int
	width      int
	height     int

	// Floating preview panel, drawn over the list without changing its layout.
	preview panelAnim
}

// NewRecommendModel creates a new recommendation screen.
func NewRecommendModel(challenges []*challenge.Challenge) tea.Model {
	return RecommendModel{challenges: challenges}
}

// NewRecommendModelWithStore creates a recommendation screen whose preview
// panel can also report each challenge's completion state.
func NewRecommendModelWithStore(challenges []*challenge.Challenge, store *progress.Store) tea.Model {
	return RecommendModel{challenges: challenges, store: store}
}

func (m RecommendModel) Init() tea.Cmd { return nil }

func (m RecommendModel) maxVisible() int {
	return visibleListRows(m.height)
}

func (m RecommendModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Re-clamp scroll so the cursor stays inside the resized viewport.
		m.offset = ensureVisible(m.cursor, m.offset, len(m.challenges), m.maxVisible())
		return m, nil
	case panelTickMsg:
		var cmd tea.Cmd
		m.preview, cmd = m.preview.advance()
		return m, cmd
	case tea.KeyMsg:
		switch {
		case runeMatches(msg, "p"):
			var cmd tea.Cmd
			m.preview, cmd = m.preview.toggle()
			return m, cmd
		case isBackKey(msg):
			// Esc closes the preview first, then leaves the screen.
			if m.preview.opening() {
				var cmd tea.Cmd
				m.preview, cmd = m.preview.toggle()
				return m, cmd
			}
			return m, func() tea.Msg { return GoBackMsg{} }
		case isSelectKey(msg):
			if len(m.challenges) > 0 {
				ch := m.challenges[m.cursor]
				return m, func() tea.Msg { return ChallengeSelectedMsg{Challenge: ch} }
			}
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

// recommendTitleWidth computes the fixed title column width: the widest title
// (capped) that still leaves room for the cursor, stars and the category
// label column inside the box.
func recommendTitleWidth(chs []*challenge.Challenge, innerW int) int {
	maxTitle := 0
	maxCat := 0
	for _, ch := range chs {
		if w := displayWidth(ch.Title); w > maxTitle {
			maxTitle = w
		}
		if w := displayWidth(CategoryLabel(ch.Category)); w > maxCat {
			maxCat = w
		}
	}
	if maxTitle > 34 {
		maxTitle = 34
	}
	// Fixed overhead: cursor 2 + gap 2 + stars 5 + gap 2 + category column.
	if budget := innerW - 11 - maxCat; budget < maxTitle {
		maxTitle = budget
	}
	if maxTitle < 4 {
		maxTitle = 4
	}
	return maxTitle
}

// recommendSpaced reports whether the rendered rows should be separated by
// blank lines (the airy layout of the design mock). Spacing is only used when
// the spaced list still fits the row budget, so long lists fall back to the
// compact one-line-per-row layout deterministically.
func recommendSpaced(rows, maxVisible int) bool {
	return rows > 0 && 2*rows-1 <= maxVisible
}

// renderRow renders one recommendation row:
//
//	› 标题…………………………  ★★★☆☆  Linux 基础命令
//
// A focused row is padded to the full inner width so its background highlight
// spans the whole line.
func (m RecommendModel) renderRow(ch *challenge.Challenge, selected bool, titleW, rowW int) string {
	rs := chRowNormal
	if selected {
		rs = chRowFocus
	}

	title := truncateWidth(ch.Title, titleW)

	var b strings.Builder
	b.WriteString(rs.CursorSlot)
	b.WriteString(rs.Text.Render(title))
	b.WriteString(rs.pad(titleW - displayWidth(title)))
	b.WriteString(rs.pad(2))
	b.WriteString(rs.stars(ch.Difficulty))
	b.WriteString(rs.pad(2))
	b.WriteString(rs.Category.Render(CategoryLabel(ch.Category)))
	row := b.String()

	if selected {
		row += rs.pad(rowW - displayWidth(row))
	}
	return row
}

func (m RecommendModel) View() string {
	if len(m.challenges) == 0 {
		body := EmptyState("还没有练习记录", "先完成一个挑战，这里会推荐薄弱环节的练习。", "Esc 返回") + "\n"
		return contentBox("薄弱推荐", body, m.width, m.height, "")
	}

	innerW := chInnerWidth(m.width)
	titleW := recommendTitleWidth(m.challenges, innerW)
	count := len(m.challenges)

	list := ScrollList{Cursor: m.cursor, Offset: m.offset}
	rows := list.RenderRows(count, m.maxVisible(), innerW, func(i int, selected bool) string {
		return m.renderRow(m.challenges[i], selected, titleW, innerW)
	})

	var body strings.Builder
	body.WriteString(S.Dim.Render("根据你的练习记录，优先补强以下题目："))
	body.WriteString("\n")

	spaced := recommendSpaced(len(rows), m.maxVisible())
	if !spaced {
		body.WriteString("\n")
	}
	for _, row := range rows {
		if spaced {
			body.WriteString("\n")
		}
		body.WriteString(row)
		body.WriteString("\n")
	}

	base := contentBox("薄弱推荐", body.String(), m.width, m.height, fmt.Sprintf("%d 题", count))
	return m.withPreviewOverlay(base)
}

// recListFirstRow is the line index of the first recommendation row inside the
// rendered box: top border, padding row, lead-in line, blank line.
const recListFirstRow = 4

// withPreviewOverlay composites the floating preview onto the rendered list.
// The list is laid out exactly as if the panel did not exist.
func (m RecommendModel) withPreviewOverlay(base string) string {
	if !m.preview.visible() || len(m.challenges) == 0 {
		return base
	}
	lines := strings.Split(base, "\n")
	boxW := maxInt(boxWidth(m.width), minBoxWidth)
	panelW := clampInt(boxW/2, chPreviewMinW, chPreviewMaxW)
	if panelW+8 > boxW || len(lines) < 6 {
		return base
	}

	// Rows may be spaced with a blank line between them; the anchor has to
	// follow the same stride as the list it points at.
	stride := 1
	if recommendSpaced(minInt(len(m.challenges), m.maxVisible()), m.maxVisible()) {
		stride = 2
	}

	descLimit := maxInt(m.maxVisible()-6, 3)
	full := challengePreviewLines(m.challenges[m.cursor], m.store, maxInt(panelW-6, 0), descLimit)
	panel := boxLines("预览", full[:m.preview.reveal(len(full))], panelW, "", false)

	leftPad := maxInt((m.width-boxW)/2, 0)
	x := leftPad + boxW - panelW - 2
	y := clampInt(recListFirstRow+stride*(m.cursor-m.offset), 1, maxInt(len(lines)-len(panel)-1, 1))
	return placeOverlay(base, strings.Join(panel, "\n"), x, y)
}
