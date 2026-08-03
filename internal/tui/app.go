package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sd3/linuxlab/internal/challenge"
	"github.com/sd3/linuxlab/internal/progress"
	"github.com/sd3/linuxlab/internal/reference"
	"github.com/sd3/linuxlab/internal/runner"
	"github.com/sd3/linuxlab/internal/sandbox"
	"github.com/sd3/linuxlab/internal/verify"
)

type screenID int

const (
	screenMenu screenID = iota
	screenModules
	screenChallenges
	screenDetail
	screenSkillMap
	screenRecommend
	screenReference
	screenResult
)

const (
	// ctrlCWindow is how long the first Ctrl+C keeps the quit armed.
	ctrlCWindow = time.Second
	// noticeTTL is how long a runtime notice stays in the footer.
	noticeTTL = 5 * time.Second

	ctrlCNotice = "再按一次 Ctrl+C 退出"

	// menuWeakLimit is how many weak-spot recommendations the recommend screen
	// shows, and therefore what the menu's "N 项薄弱" counter reports.
	menuWeakLimit = 5
)

// ChallengeResultMsg is sent after a challenge attempt completes.
type ChallengeResultMsg struct {
	Passed    bool
	Results   []verify.Result
	HintsUsed int
}

// noticeExpireMsg clears the footer notice set by the matching armNotice call.
type noticeExpireMsg struct{ seq int }

// dockerStatusMsg carries the async Docker availability probe result.
type dockerStatusMsg struct{ ok bool }

// probeDockerCmd checks Docker availability off the event loop (the sandbox
// package caches the probe, so repeated calls are cheap).
func probeDockerCmd() tea.Cmd {
	return func() tea.Msg {
		return dockerStatusMsg{ok: sandbox.DockerAvailable()}
	}
}

// AppModel is the root TUI model: a message router plus the app shell
// (header / active screen / footer) compositor.
type AppModel struct {
	screen screenID
	stack  []screenID // navigation history for GoBackMsg (push on forward nav)

	categories       map[string][]*challenge.Challenge
	store            *progress.Store
	refs             *reference.ReferenceData
	currentCat       string
	currentChallenge *challenge.Challenge
	totalChallenges  int

	width  int
	height int
	ready  bool // no size-dependent layout before the first WindowSizeMsg

	dockerProbed bool
	dockerOK     bool

	// showSplash covers the app with the opening screen until the first
	// keystroke. It is a layer over the menu rather than a screen of its own,
	// so the navigation state machine starts at screenMenu as before.
	showSplash bool

	help       help.Model
	notice     string
	noticeSeq  int
	ctrlCArmed bool

	menu       tea.Model
	modules    tea.Model
	challenges tea.Model
	detail     tea.Model
	skillmap   tea.Model
	recommend  tea.Model
	refModel   tea.Model

	lastResult *ChallengeResultMsg
}

// NewAppModel creates the root app model.
func NewAppModel(cats map[string][]*challenge.Challenge, store *progress.Store, refs *reference.ReferenceData) tea.Model {
	totalChallenges := 0
	for _, chs := range cats {
		totalChallenges += len(chs)
	}
	return AppModel{
		screen:          screenMenu,
		categories:      cats,
		store:           store,
		refs:            refs,
		totalChallenges: totalChallenges,
		help:            help.New(),
		menu:            NewMenuModelWithData(computeMenuStats(cats, store, refs)),
		showSplash:      true,
	}
}

// splashInfo snapshots what the cover displays. Counters come from the cached
// menu model rather than being recomputed, so View stays cheap.
func (m AppModel) splashInfo() splashInfo {
	info := splashInfo{dockerProbed: m.dockerProbed, dockerOK: m.dockerOK}
	if menu, ok := m.menu.(MenuModel); ok {
		info.stats = menu.stats
	}
	return info
}

// dismissSplash clears the cover. Navigation messages call it too, so a screen
// reached programmatically is never rendered underneath the cover.
func (m AppModel) dismissSplash() AppModel {
	m.showSplash = false
	return m
}

// allChallenges flattens the category map. Order is category-map dependent, so
// callers that render the result must sort or otherwise stabilize it.
func allChallenges(cats map[string][]*challenge.Challenge) []*challenge.Challenge {
	var all []*challenge.Challenge
	for _, chs := range cats {
		all = append(all, chs...)
	}
	return all
}

// computeMenuStats derives the menu's status-word counters. It walks the whole
// progress set and builds a skill map, so it is called on menu transitions
// only — never from View.
func computeMenuStats(cats map[string][]*challenge.Challenge, store *progress.Store, refs *reference.ReferenceData) MenuStats {
	all := allChallenges(cats)
	stats := MenuStats{TotalChallenges: len(all), TotalModules: len(cats)}
	if store != nil {
		for _, ch := range all {
			if entry, ok := store.Data.Challenges[ch.ID]; ok && entry != nil && entry.Status == "passed" {
				stats.PassedChallenges++
			}
		}
		recs := progress.RecommendMultiple(store, all, menuWeakLimit)
		stats.WeakCount = len(recs)
		if len(recs) > 0 {
			stats.Next = recs[0]
		}
		stats.Last, stats.LastEntry = progress.LastAttempted(store, all)
	}
	if refs != nil {
		stats.ReferenceCount = len(refs.Commands)
	}
	return stats
}

// refreshMenuStats recomputes the menu counters in place, keeping the cached
// menu model (and its cursor) alive.
func (m *AppModel) refreshMenuStats() {
	menu, ok := m.menu.(MenuModel)
	if !ok {
		return
	}
	m.menu = menu.withStats(computeMenuStats(m.categories, m.store, m.refs))
}

func (m AppModel) Init() tea.Cmd {
	return m.menu.Init()
}

// sizeModel sends a WindowSizeMsg to a sub-model so it knows the terminal dimensions.
func (m AppModel) sizeModel(sub tea.Model) tea.Model {
	sized, _ := sub.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
	return sized
}

// pushScreen records the current screen on the navigation stack and switches
// to next.
func (m *AppModel) pushScreen(next screenID) {
	m.stack = append(m.stack, m.screen)
	m.screen = next
}

// armNotice sets the footer notice and returns the command that clears it
// after ttl (stale timers are ignored via the sequence number).
func (m *AppModel) armNotice(text string, ttl time.Duration) tea.Cmd {
	m.notice = text
	m.noticeSeq++
	seq := m.noticeSeq
	return tea.Tick(ttl, func(time.Time) tea.Msg { return noticeExpireMsg{seq: seq} })
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Any navigation means the session has started, so the cover comes down
	// even when a screen is reached programmatically rather than by keystroke.
	switch msg.(type) {
	case MenuChoiceMsg, ModuleSelectedMsg, ChallengeSelectedMsg,
		LaunchChallengeMsg, ChallengeResultMsg, GoBackMsg:
		m = m.dismissSplash()
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		firstSize := !m.ready
		m.ready = true
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		// Broadcast to every live sub-model so screens revisited later (e.g.
		// going back to the menu) render with the current terminal size.
		// Only the active sub-model's command is kept.
		var cmd tea.Cmd
		resize := func(sub tea.Model, active bool) tea.Model {
			if sub == nil {
				return nil
			}
			updated, c := sub.Update(msg)
			if active {
				cmd = c
			}
			return updated
		}
		m.menu = resize(m.menu, m.screen == screenMenu)
		m.modules = resize(m.modules, m.screen == screenModules)
		m.challenges = resize(m.challenges, m.screen == screenChallenges)
		m.detail = resize(m.detail, m.screen == screenDetail)
		m.skillmap = resize(m.skillmap, m.screen == screenSkillMap)
		m.recommend = resize(m.recommend, m.screen == screenRecommend)
		m.refModel = resize(m.refModel, m.screen == screenReference)
		if firstSize {
			return m, tea.Batch(cmd, probeDockerCmd())
		}
		return m, cmd

	case dockerStatusMsg:
		m.dockerProbed = true
		m.dockerOK = msg.ok
		return m, nil

	case noticeExpireMsg:
		if msg.seq == m.noticeSeq {
			m.notice = ""
			m.ctrlCArmed = false
		}
		return m, nil

	case tea.KeyMsg:
		// Ctrl+C requires a double press to quit (the menu's q stays an
		// immediate, explicit exit).
		if msg.Type == tea.KeyCtrlC {
			if m.ctrlCArmed {
				return m, tea.Quit
			}
			m.ctrlCArmed = true
			return m, m.armNotice(ctrlCNotice, ctrlCWindow)
		}
		if m.ctrlCArmed {
			// Any other key disarms the pending quit.
			m.ctrlCArmed = false
			if m.notice == ctrlCNotice {
				m.notice = ""
			}
		}
		// The cover consumes the first keystroke: q still quits from it, and
		// anything else simply reveals the menu. It only ever covers the menu,
		// so a screen entered by other means is never intercepted.
		if m.showSplash && m.screen == screenMenu {
			if runeMatches(msg, "q") {
				return m, tea.Quit
			}
			return m.dismissSplash(), nil
		}
		if key.Matches(msg, globalKeys.Help) && m.helpToggleAllowed() {
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		}

	case MenuChoiceMsg:
		switch msg.Choice {
		case "practice":
			m.pushScreen(screenModules)
			m.modules = m.sizeModel(NewModulesModel(m.categories, m.store))
			return m, nil
		case "skillmap":
			m.pushScreen(screenSkillMap)
			m.skillmap = m.sizeModel(NewSkillMapModel(m.store))
			return m, nil
		case "recommend":
			recs := progress.RecommendMultiple(m.store, allChallenges(m.categories), menuWeakLimit)
			m.pushScreen(screenRecommend)
			m.recommend = m.sizeModel(NewRecommendModelWithStore(recs, m.store))
			return m, nil
		case "reference":
			if m.refs != nil {
				m.pushScreen(screenReference)
				m.refModel = m.sizeModel(NewReferenceModelWithChallenges(m.refs, m.categories, m.store))
				return m, nil
			}
			return m, nil
		}
		return m, nil

	case ModuleSelectedMsg:
		m.pushScreen(screenChallenges)
		m.currentCat = msg.Category
		m.challenges = m.sizeModel(NewChallengesModel(msg.Category, m.categories[msg.Category], m.store))
		return m, nil

	case ChallengeSelectedMsg:
		m.pushScreen(screenDetail)
		m.currentChallenge = msg.Challenge
		if msg.Challenge != nil {
			m.currentCat = msg.Challenge.Category
		}
		m.detail = m.sizeModel(NewDetailModel(msg.Challenge))
		return m, nil

	case LaunchChallengeMsg:
		return m, m.launchChallenge(msg)

	case ChallengeResultMsg:
		return m.handleChallengeResult(msg)

	case GoBackMsg:
		return m.goBack()
	}

	// Delegate to active sub-model
	var cmd tea.Cmd
	switch m.screen {
	case screenMenu:
		m.menu, cmd = m.menu.Update(msg)
	case screenModules:
		if m.modules != nil {
			m.modules, cmd = m.modules.Update(msg)
		}
	case screenChallenges:
		if m.challenges != nil {
			m.challenges, cmd = m.challenges.Update(msg)
		}
	case screenDetail:
		if m.detail != nil {
			m.detail, cmd = m.detail.Update(msg)
		}
	case screenSkillMap:
		if m.skillmap != nil {
			m.skillmap, cmd = m.skillmap.Update(msg)
		}
	case screenRecommend:
		if m.recommend != nil {
			m.recommend, cmd = m.recommend.Update(msg)
		}
	case screenReference:
		if m.refModel != nil {
			m.refModel, cmd = m.refModel.Update(msg)
		}
	case screenResult:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			return m.updateResult(keyMsg)
		}
	}
	return m, cmd
}

// helpToggleAllowed reports whether '?' may toggle the full help on the
// current screen. On the reference list '?' belongs to the search query
// (search-first semantics, locked behavior).
func (m AppModel) helpToggleAllowed() bool {
	if m.screen == screenReference {
		if rm, ok := m.refModel.(ReferenceModel); ok && !rm.showDetail {
			return false
		}
	}
	return true
}

// handleChallengeResult records the attempt and shows the result screen. The
// result replaces the detail screen (no push), so going back lands on the
// list the challenge was opened from. After tea.Exec the terminal may have
// been resized by the external session, so the size is re-queried and the
// Docker probe refreshed.
func (m AppModel) handleChallengeResult(msg ChallengeResultMsg) (tea.Model, tea.Cmd) {
	m.lastResult = &msg
	cmds := []tea.Cmd{tea.WindowSize(), probeDockerCmd()}
	ch := m.currentChallenge
	if ch != nil {
		m.store.RecordAttempt(ch.ID, ch.Category, ch.Subcategory, msg.Passed, msg.HintsUsed)
		m.store.Save() // best-effort persist
		cmds = append(cmds, m.armNotice("已保存进度", noticeTTL))
	}
	m.screen = screenResult
	return m, tea.Batch(cmds...)
}

// goBack pops the navigation stack. Sub-model caches are restored so cursor
// and scroll positions survive (locked behavior id=19).
func (m AppModel) goBack() (tea.Model, tea.Cmd) {
	if len(m.stack) == 0 {
		m.screen = screenMenu
		m.refreshMenuStats()
		return m, nil
	}
	target := m.stack[len(m.stack)-1]
	m.stack = m.stack[:len(m.stack)-1]

	switch target {
	case screenMenu:
		// Counters (passed / weak spots) may have changed while away.
		m.refreshMenuStats()
	case screenModules:
		// Reuse the existing modules model to keep the cursor position;
		// progress counts are read live from the store in View.
		if m.modules == nil {
			m.modules = m.sizeModel(NewModulesModel(m.categories, m.store))
		}
	case screenChallenges:
		m.challenges = m.restoredChallenges()
	}
	m.screen = target
	return m, nil
}

// updateResult handles keys on the result screen (no sub-model).
func (m AppModel) updateResult(keyMsg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case isBackKey(keyMsg):
		return m, func() tea.Msg { return GoBackMsg{} }
	case runeMatches(keyMsg, "r"):
		if m.currentChallenge != nil {
			ch := m.currentChallenge
			// Carry over the hints already unlocked in the previous
			// attempt so the retry result keeps showing them.
			hints := 0
			if m.lastResult != nil {
				hints = m.lastResult.HintsUsed
			}
			return m, func() tea.Msg { return LaunchChallengeMsg{Challenge: ch, HintsUsed: hints} }
		}
	case keyMsg.Type == tea.KeyEnter || keyMsg.Type == tea.KeyRight || runeMatches(keyMsg, "n"):
		if next := m.nextChallenge(); next != nil {
			// The next detail replaces the current result: the stack top is
			// still the list screen, so back returns there.
			m.screen = screenDetail
			m.currentChallenge = next
			m.detail = m.sizeModel(NewDetailModel(next))
			return m, nil
		}
		return m, func() tea.Msg { return GoBackMsg{} }
	}
	return m, nil
}

// View composes the app shell: 1-line header, the active screen centered in
// the main slot, 1-line footer. The output always has exactly m.height lines
// and every line's display width is at most m.width (iron rules #2/#3).
func (m AppModel) View() string {
	if !m.ready || m.width <= 0 || m.height <= 0 {
		return ""
	}
	spec := computeLayout(m.width, m.height)
	if spec.Mode == ModeUnsupported {
		return m.unsupportedView()
	}
	if m.showSplash && m.screen == screenMenu {
		cover := renderSplash(m.width, m.height, m.splashInfo())
		return finalizeFrame(fitToSlot(cover, m.width, m.height), m.width, m.height)
	}

	header := renderHeader(m.breadcrumbs(), m.headerRight(), m.width)
	footer := m.footerView()
	mainH := m.height - lipgloss.Height(header) - lipgloss.Height(footer)

	frame := make([]string, 0, m.height)
	frame = append(frame, strings.Split(header, "\n")...)
	frame = append(frame, fitToSlot(m.mainView(), m.width, mainH)...)
	frame = append(frame, strings.Split(footer, "\n")...)
	return finalizeFrame(frame, m.width, m.height)
}

// mainView renders the active screen's content (without shell chrome).
func (m AppModel) mainView() string {
	switch m.screen {
	case screenMenu:
		return m.menu.View()
	case screenModules:
		if m.modules != nil {
			return m.modules.View()
		}
	case screenChallenges:
		if m.challenges != nil {
			return m.challenges.View()
		}
	case screenDetail:
		if m.detail != nil {
			return m.detail.View()
		}
	case screenSkillMap:
		if m.skillmap != nil {
			return m.skillmap.View()
		}
	case screenRecommend:
		if m.recommend != nil {
			return m.recommend.View()
		}
	case screenReference:
		if m.refModel != nil {
			return m.refModel.View()
		}
	case screenResult:
		return m.resultView()
	}
	return ""
}

// breadcrumbs builds the header breadcrumb for the current screen.
func (m AppModel) breadcrumbs() []string {
	crumbs := []string{"LinuxLab"}
	appendCat := func() {
		if label := CategoryLabel(m.currentCat); m.currentCat != "" && label != "" {
			crumbs = append(crumbs, label)
		}
	}
	switch m.screen {
	case screenMenu:
		crumbs = append(crumbs, "主菜单")
	case screenModules:
		crumbs = append(crumbs, "练习")
	case screenChallenges:
		crumbs = append(crumbs, "练习")
		appendCat()
	case screenDetail:
		crumbs = append(crumbs, "练习")
		appendCat()
		if m.currentChallenge != nil {
			crumbs = append(crumbs, m.currentChallenge.Title)
		}
	case screenSkillMap:
		crumbs = append(crumbs, "能力图谱")
	case screenRecommend:
		crumbs = append(crumbs, "薄弱推荐")
	case screenReference:
		crumbs = append(crumbs, "命令速查")
	case screenResult:
		crumbs = append(crumbs, "练习")
		if m.currentChallenge != nil {
			crumbs = append(crumbs, m.currentChallenge.Title)
		}
		crumbs = append(crumbs, "检测结果")
	}
	return crumbs
}

// headerRight builds the right header slot: Docker status + overall progress.
func (m AppModel) headerRight() string {
	right := S.Meta.Render(fmt.Sprintf("%d/%d", m.passedTotal(), m.totalChallenges))
	if m.dockerProbed {
		icon := S.IconFail
		if m.dockerOK {
			icon = S.IconPass
		}
		right = S.Meta.Render("Docker ") + icon + S.Meta.Render(" · ") + right
	}
	return right
}

// passedTotal counts passed challenges across all categories (map iteration
// order does not matter for a count, so the output stays frame-stable).
func (m AppModel) passedTotal() int {
	if m.store == nil {
		return 0
	}
	passed := 0
	for _, chs := range m.categories {
		for _, ch := range chs {
			if entry, ok := m.store.Data.Challenges[ch.ID]; ok && entry.Status == "passed" {
				passed++
			}
		}
	}
	return passed
}

// footerView renders the shell footer: notice left, short key help right;
// with ShowAll toggled the bubbles/help full view replaces the single line.
func (m AppModel) footerView() string {
	km := m.currentKeyMap()
	if m.help.ShowAll {
		return m.help.View(km)
	}
	return renderFooter(m.notice, km, m.width)
}

// currentKeyMap returns the active screen's keymap for footer help.
func (m AppModel) currentKeyMap() help.KeyMap {
	switch m.screen {
	case screenModules:
		return modulesKeys
	case screenChallenges:
		if cm, ok := m.challenges.(ChallengesModel); ok && cm.filtering {
			return challengesFilterKeys
		}
		return challengesKeys
	case screenDetail:
		return detailKeys
	case screenSkillMap:
		return skillmapKeys
	case screenRecommend:
		return recommendKeys
	case screenReference:
		return referenceKeys
	case screenResult:
		km := resultKeys
		km.Next.SetEnabled(m.nextChallenge() != nil)
		return km
	default:
		return menuKeys
	}
}

// unsupportedView renders the "terminal too small" screen, still honoring the
// exact-height/width invariants.
func (m AppModel) unsupportedView() string {
	body := fmt.Sprintf("终端窗口太小\n请调整到至少 %d×%d", minTermWidth, minTermHeight)
	box := contentBox("", body, m.width, m.height, "")
	return finalizeFrame(fitToSlot(box, m.width, m.height), m.width, m.height)
}

// fitToSlot vertically centers content inside a slot of exactly `height`
// lines: taller content is clipped at the bottom, shorter content is padded.
// Every line is truncated to the given width.
func fitToSlot(content string, width, height int) []string {
	if height < 0 {
		height = 0
	}
	lines := strings.Split(content, "\n")
	for i := range lines {
		lines[i] = truncateWidth(lines[i], width)
	}
	if len(lines) > height {
		lines = lines[:height]
	}
	top := (height - len(lines)) / 2
	out := make([]string, 0, height)
	for i := 0; i < top; i++ {
		out = append(out, "")
	}
	out = append(out, lines...)
	for len(out) < height {
		out = append(out, "")
	}
	return out
}

// finalizeFrame enforces the frame invariants: exactly `height` lines, every
// line at most `width` display columns, trailing spaces trimmed (keeps the
// renderer's line diff small).
func finalizeFrame(lines []string, width, height int) string {
	if height < 1 {
		return ""
	}
	if len(lines) > height {
		lines = lines[:height]
	}
	for len(lines) < height {
		lines = append(lines, "")
	}
	for i := range lines {
		lines[i] = strings.TrimRight(truncateWidth(lines[i], width), " ")
	}
	return strings.Join(lines, "\n")
}

// restoredChallenges returns the challenge list model to show when navigating
// back from the detail/result screens. The existing model is reused so cursor
// and scroll position survive (pass/fail status is read live from the store in
// View); it is only rebuilt when missing or when the category changed (e.g.
// the detail screen was entered from the recommend list).
func (m AppModel) restoredChallenges() tea.Model {
	sub := m.challenges
	cm, ok := sub.(ChallengesModel)
	if !ok || (m.currentCat != "" && cm.category != m.currentCat) {
		if m.currentCat == "" {
			return sub
		}
		sub = m.sizeModel(NewChallengesModel(m.currentCat, m.categories[m.currentCat], m.store))
		cm, ok = sub.(ChallengesModel)
	}
	// Point the cursor at the challenge just visited (it may have advanced via
	// "next challenge" on the result screen).
	if ok && m.currentChallenge != nil {
		sub = cm.focusChallenge(m.currentChallenge.ID)
	}
	return sub
}

func (m AppModel) launchChallenge(msg LaunchChallengeMsg) tea.Cmd {
	cmd := &runner.Command{
		Ctx: context.Background(),
		Opts: runner.Options{
			Challenge: msg.Challenge,
			Refs:      m.refs,
			HintsUsed: msg.HintsUsed,
		},
	}
	return tea.Exec(cmd, func(err error) tea.Msg {
		if err != nil {
			return ChallengeResultMsg{
				Passed:    false,
				Results:   []verify.Result{{Passed: false, Message: fmt.Sprintf("挑战执行失败: %v", err)}},
				HintsUsed: msg.HintsUsed,
			}
		}
		return ChallengeResultMsg{
			Passed:    cmd.Result.Passed,
			Results:   cmd.Result.Results,
			HintsUsed: cmd.Result.HintsUsed,
		}
	})
}

func (m AppModel) nextChallenge() *challenge.Challenge {
	if m.currentCat == "" {
		return nil
	}
	chs := m.categories[m.currentCat]
	if len(chs) == 0 {
		return nil
	}
	if m.currentChallenge == nil {
		return chs[0]
	}
	for i, ch := range chs {
		if ch.ID == m.currentChallenge.ID && i+1 < len(chs) {
			return chs[i+1]
		}
	}
	return nil
}

func (m AppModel) resultView() string {
	var body strings.Builder

	if m.lastResult == nil {
		body.WriteString(DimStyle.Render("无结果"))
		body.WriteString("\n")
		return contentBox("检测结果", body.String(), m.width, m.height, "")
	}

	if m.lastResult.Passed {
		body.WriteString(SuccessStyle.Render("● 挑战通过！"))
	} else {
		body.WriteString(ErrorStyle.Render("● 挑战未通过"))
	}
	body.WriteString("\n\n")

	for i, r := range m.lastResult.Results {
		icon := PassedIcon
		if !r.Passed {
			icon = FailedIcon
		}
		body.WriteString(fmt.Sprintf("%s  检查 %d: %s\n", icon, i+1, r.Message))
	}

	if m.lastResult.HintsUsed > 0 {
		body.WriteString(fmt.Sprintf("\n%s\n", DimStyle.Render(fmt.Sprintf("使用提示: %d", m.lastResult.HintsUsed))))
	}

	return contentBox("检测结果", body.String(), m.width, m.height, "")
}
