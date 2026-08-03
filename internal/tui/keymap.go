package tui

import "github.com/charmbracelet/bubbles/key"

// keymap.go is the single presentation-layer source of truth for key
// bindings and their (Chinese) help texts. Actual key matching still goes
// through navigation.go (isBackKey / isSelectKey / moveListCursor), whose
// semantics are locked by tests; the bindings here are defined to match those
// semantics exactly, and help lines are generated from them instead of being
// hand-written per screen.

// GlobalKeyMap holds bindings shared by every screen.
type GlobalKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	Home     key.Binding
	End      key.Binding
	Select   key.Binding
	Back     key.Binding
	Help     key.Binding
	Quit     key.Binding
}

func newGlobalKeyMap() GlobalKeyMap {
	return GlobalKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "上移"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "下移"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "ctrl+b"),
			key.WithHelp("PgUp", "上翻页"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown", "ctrl+f"),
			key.WithHelp("PgDn", "下翻页"),
		),
		Home: key.NewBinding(
			key.WithKeys("home", "ctrl+a", "g"),
			key.WithHelp("g", "顶部"),
		),
		End: key.NewBinding(
			key.WithKeys("end", "ctrl+e", "G"),
			key.WithHelp("G", "底部"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter", "right"),
			key.WithHelp("Enter", "确认"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc", "left", "q"),
			key.WithHelp("q/Esc", "返回"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "帮助"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("Ctrl+C", "退出"),
		),
	}
}

// globalKeys is the canonical global keymap instance.
var globalKeys = newGlobalKeyMap()

// navHelp is a display-only binding summarizing list navigation.
func navHelp() key.Binding {
	return key.NewBinding(
		key.WithKeys("up", "down", "k", "j"),
		key.WithHelp("↑↓", "选择"),
	)
}

// pageHelp is a display-only binding summarizing paging.
func pageHelp() key.Binding {
	return key.NewBinding(
		key.WithKeys("pgup", "pgdown"),
		key.WithHelp("PgUp/Dn", "翻页"),
	)
}

// --- Per-screen keymaps -----------------------------------------------------
// Each implements bubbles/help.KeyMap. ShortHelp orders bindings with the
// most essential keys LAST — the footer collapses by dropping from the front
// when the terminal is narrow, so Esc/Enter survive the longest.

// MenuKeyMap is the main menu keymap.
type MenuKeyMap struct {
	Nav    key.Binding
	Jump   key.Binding
	Select key.Binding
	Quit   key.Binding
}

func newMenuKeyMap() MenuKeyMap {
	return MenuKeyMap{
		Nav:    navHelp(),
		Jump:   key.NewBinding(key.WithKeys("1", "2", "3", "4"), key.WithHelp("1-4", "跳转")),
		Select: globalKeys.Select,
		Quit:   key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "退出")),
	}
}

func (k MenuKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Jump, k.Nav, k.Select, k.Quit}
}

func (k MenuKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Nav, k.Jump},
		{k.Select, k.Quit},
		{globalKeys.Help, globalKeys.Quit},
	}
}

// ModulesKeyMap is the module list keymap.
type ModulesKeyMap struct {
	Nav    key.Binding
	Jump   key.Binding
	Select key.Binding
	Back   key.Binding
}

func newModulesKeyMap() ModulesKeyMap {
	return ModulesKeyMap{
		Nav:    navHelp(),
		Jump:   key.NewBinding(key.WithKeys("1", "2", "3", "4", "5", "6", "7", "8", "9"), key.WithHelp("数字", "跳转")),
		Select: globalKeys.Select,
		Back:   globalKeys.Back,
	}
}

func (k ModulesKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Jump, k.Nav, k.Select, k.Back}
}

func (k ModulesKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Nav, k.Jump},
		{k.Select, k.Back},
		{globalKeys.Help, globalKeys.Quit},
	}
}

// ChallengesKeyMap is the challenge list keymap.
type ChallengesKeyMap struct {
	Nav     key.Binding
	Page    key.Binding
	Filter  key.Binding
	Preview key.Binding
	Select  key.Binding
	Back    key.Binding
}

func newChallengesKeyMap() ChallengesKeyMap {
	return ChallengesKeyMap{
		Nav:     navHelp(),
		Page:    pageHelp(),
		Filter:  key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "过滤")),
		Preview: key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "预览")),
		Select:  globalKeys.Select,
		Back:    globalKeys.Back,
	}
}

// Paging is a universal convention and stays in the full help; the preview is
// specific to this app, so it takes the scarce footer slot instead.
func (k ChallengesKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Preview, k.Nav, k.Filter, k.Select, k.Back}
}

func (k ChallengesKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Nav, k.Page},
		{globalKeys.Home, globalKeys.End},
		{k.Preview, k.Filter},
		{k.Select, k.Back},
		{globalKeys.Help, globalKeys.Quit},
	}
}

// ChallengesFilterKeyMap is shown while the challenge list filter is active.
type ChallengesFilterKeyMap struct {
	Exit   key.Binding
	Nav    key.Binding
	Select key.Binding
}

func newChallengesFilterKeyMap() ChallengesFilterKeyMap {
	return ChallengesFilterKeyMap{
		Exit:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "退出过滤")),
		Nav:    navHelp(),
		Select: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "进入")),
	}
}

func (k ChallengesFilterKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Nav, k.Select, k.Exit}
}

func (k ChallengesFilterKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Nav, k.Select},
		{k.Exit},
		{globalKeys.Help, globalKeys.Quit},
	}
}

// DetailKeyMap is the challenge detail keymap.
type DetailKeyMap struct {
	Scroll    key.Binding
	Hint      key.Binding
	Inspector key.Binding
	Launch    key.Binding
	Back      key.Binding
}

func newDetailKeyMap() DetailKeyMap {
	return DetailKeyMap{
		Scroll:    key.NewBinding(key.WithKeys("up", "down", "k", "j"), key.WithHelp("↑↓", "滚动")),
		Hint:      key.NewBinding(key.WithKeys("h"), key.WithHelp("h", "查看提示")),
		Inspector: key.NewBinding(key.WithKeys("i"), key.WithHelp("i", "元信息")),
		Launch:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("Enter", "开始挑战")),
		Back:      globalKeys.Back,
	}
}

func (k DetailKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Scroll, k.Inspector, k.Hint, k.Launch, k.Back}
}

func (k DetailKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Scroll, pageHelp()},
		{k.Inspector, k.Hint},
		{k.Launch, k.Back},
		{globalKeys.Help, globalKeys.Quit},
	}
}

// SkillMapKeyMap is the skill map keymap.
type SkillMapKeyMap struct {
	Nav    key.Binding
	Toggle key.Binding
	Back   key.Binding
}

func newSkillMapKeyMap() SkillMapKeyMap {
	return SkillMapKeyMap{
		Nav:    navHelp(),
		Toggle: key.NewBinding(key.WithKeys("enter", " "), key.WithHelp("Enter/Space", "展开")),
		Back:   globalKeys.Back,
	}
}

func (k SkillMapKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Nav, k.Toggle, k.Back}
}

func (k SkillMapKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Nav, k.Toggle},
		{k.Back, globalKeys.Quit},
	}
}

// RecommendKeyMap is the recommendation list keymap.
type RecommendKeyMap struct {
	Nav     key.Binding
	Page    key.Binding
	Preview key.Binding
	Select  key.Binding
	Back    key.Binding
}

func newRecommendKeyMap() RecommendKeyMap {
	return RecommendKeyMap{
		Nav:     navHelp(),
		Page:    pageHelp(),
		Preview: key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "预览")),
		Select:  globalKeys.Select,
		Back:    globalKeys.Back,
	}
}

func (k RecommendKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Preview, k.Nav, k.Select, k.Back}
}

func (k RecommendKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Nav, k.Page},
		{k.Preview},
		{k.Select, k.Back},
		{globalKeys.Help, globalKeys.Quit},
	}
}

// ReferenceKeyMap is the command reference keymap (search-first: printable
// runes always go into the query, per locked behavior id=12).
type ReferenceKeyMap struct {
	Type   key.Binding
	Nav    key.Binding
	Select key.Binding
	Clear  key.Binding
	Back   key.Binding
}

func newReferenceKeyMap() ReferenceKeyMap {
	return ReferenceKeyMap{
		Type:   key.NewBinding(key.WithKeys("runes"), key.WithHelp("输入", "搜索")),
		Nav:    key.NewBinding(key.WithKeys("up", "down"), key.WithHelp("↑↓", "选择")),
		Select: key.NewBinding(key.WithKeys("enter"), key.WithHelp("Enter", "查看详情")),
		Clear:  key.NewBinding(key.WithKeys("ctrl+u"), key.WithHelp("Ctrl+U", "清空")),
		Back:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("Esc", "清空/返回")),
	}
}

func (k ReferenceKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Type, k.Nav, k.Clear, k.Select, k.Back}
}

func (k ReferenceKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Type, k.Clear},
		{k.Nav, pageHelp()},
		{k.Select, k.Back},
		{globalKeys.Quit},
	}
}

// ResultKeyMap is the challenge result keymap.
type ResultKeyMap struct {
	Retry key.Binding
	Next  key.Binding
	Back  key.Binding
}

func newResultKeyMap() ResultKeyMap {
	return ResultKeyMap{
		Retry: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "再试一次")),
		Next:  key.NewBinding(key.WithKeys("n", "enter"), key.WithHelp("n/Enter", "下一题")),
		Back:  globalKeys.Back,
	}
}

func (k ResultKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Retry, k.Next, k.Back}
}

func (k ResultKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Retry, k.Next},
		{k.Back, globalKeys.Quit},
	}
}

// Per-screen canonical keymap instances (pre-built once).
var (
	menuKeys             = newMenuKeyMap()
	modulesKeys          = newModulesKeyMap()
	challengesKeys       = newChallengesKeyMap()
	challengesFilterKeys = newChallengesFilterKeyMap()
	detailKeys           = newDetailKeyMap()
	skillmapKeys         = newSkillMapKeyMap()
	recommendKeys        = newRecommendKeyMap()
	referenceKeys        = newReferenceKeyMap()
	resultKeys           = newResultKeyMap()
)
