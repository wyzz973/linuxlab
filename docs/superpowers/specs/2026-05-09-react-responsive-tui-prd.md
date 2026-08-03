# LinuxLab React Responsive TUI PRD

## 1. 背景

LinuxLab 目前有两套终端界面：

- 生产可用的 Go/Bubble Tea TUI：路径在 `internal/tui/`，已经接入 Docker、Compose、LocalSandbox、Vim handoff、验证和进度记录。
- 试验中的 React/Ink TUI：路径在 `src/react-tui/`，已经可以加载真实题库、进度和命令速查数据，但当前布局固定宽度，不能自适应 terminal 尺寸，挑战执行尚未接入。

用户希望 React TUI 的体验参考 Codex、Claude Code、opencode、Hermes 这类现代 agent CLI。核心问题不是“再做漂亮一点”，而是需要一套响应式终端 UI 架构：组件根据 `stdout.columns` / `stdout.rows` 自适应布局、保持输入区稳定、滚动区可控、窄屏不溢出、宽屏有足够信息密度。

## 2. 参考来源与结论

### 2.1 Codex CLI

公开资料显示 Codex TUI 是 fullscreen terminal interface，使用 Ratatui + Crossterm。其渲染流程是事件采集、状态更新、构建 widget tree、渲染 frame。这个模式说明 LinuxLab React TUI 也应该有明确的中央状态和纯渲染组件，而不是在组件里散落业务状态。

Codex 的样式约束强调语义颜色：cyan 用于用户输入、选择、状态指标；green 用于成功；red 用于错误和删除；并避免不可控的硬编码颜色。它还强调文本 wrapping utility 和 snapshot testing。这对 LinuxLab 的直接启发是：

- 颜色只表达状态，不做装饰。
- 所有文案都必须经过宽度裁剪/换行。
- UI 改动必须有 snapshot 或 render-to-string 测试覆盖。
- terminal resize 是一等事件，不能靠固定 `width={96}`。

来源：

- https://www.mintlify.com/openai/codex/architecture/tui
- https://github.com/openai/codex

### 2.2 opencode / OpenTUI

opencode 的 TUI 使用 OpenTUI + SolidJS，运行在 CLI 进程内，并通过 SDK client 与后端 HTTP/SSE 同步。它有 route system，区分 home 和 session view；session view 又拆为 header、sidebar、footer、prompt、dialog、command palette 等组件。它的 prompt 组件支持 autocomplete、history、stash、attachments、slash commands。

LinuxLab 不需要照搬 agent chat，但应该借鉴这些结构：

- `Route` / `Screen` 是显式状态，不用散落条件判断。
- `Header`、`Sidebar`、`MainPane`、`Footer`、`ModalLayer`、`PromptBar` 是稳定布局槽位。
- 宽屏显示 sidebar，窄屏折叠为 top tabs。
- 命令速查、题目选择、提示、帮助等应走 command palette / modal overlay，而不是跳全屏。
- 后端执行应被抽象为 client/service，不让 React 直接复制 Go sandbox 逻辑。

来源：

- https://deepwiki.com/baumaxt/opencode/4-user-interfaces
- https://github.com/anomalyco/opentui

### 2.3 Hermes Agent

Hermes TUI 的公开文档强调：TUI 是现代前端，但复用同一个 Python runtime、同一套 session 和 slash commands；Node TUI 作为 subprocess 启动。它的关键体验点包括 instant first frame、non-blocking input、modal overlays、live session panel、mouse-friendly selection、alternate-screen differential rendering、paste safety、文件/图片附件规范化。

LinuxLab 的直接结论：

- React TUI 应该复用 Go 执行引擎，而不是重写 Docker/Vim/verify。
- 启动必须先画 skeleton，再异步加载题库/进度/引用数据。
- 挑战启动、提示、确认、错误都应该用 modal 或 bottom notice，不打断主布局。
- 粘贴长文本、路径、命令时要有折叠和安全处理。

来源：

- https://hermes-agent.nousresearch.com/docs/user-guide/tui

### 2.4 Claude Code

Claude Code 官方 CLI 源码不是稳定的公开开源参考。PRD 不引用、不依赖任何泄露源码。这里只把它作为可观察产品行为参考：紧凑对话/任务流、底部 composer、slash commands、清晰的执行状态、可中断的长任务、少量语义色。

### 2.5 Charmbracelet 组件库

用户补充的 Bubbles、Harmonica、BubbleZone、Glamour 都属于 Go/Bubble Tea 生态。LinuxLab 的既有生产 TUI 也基于 Bubble Tea，因此这些库有两层价值：

- 作为现有 Go TUI 的成熟组件参考：Bubbles 提供 list、viewport、textinput、textarea、table、progress、spinner、help、keybinding 等组件；这些正好对应 LinuxLab 目前手写列表、详情滚动、帮助栏、搜索框和进度条的薄弱点。
- 作为 React/Ink 重构的行为基准：即使 React/Ink 不直接引入这些 Go 包，也应该实现同等组件能力，而不是只做静态文本布局。

具体结论：

- Bubbles `viewport` / `list` 证明滚动和分页应该是组件级能力，LinuxLab React TUI 需要 `ScrollList`、`Viewport`、`Paginator` 的明确边界，不能在每个 screen 内手写 `slice()`。
- Bubbles `key` / `help` 的模式说明快捷键定义和底部帮助应该由同一份 keymap 生成，避免快捷键实际行为和 footer 文案不一致。
- Bubbles `textinput` / `textarea` 支持 unicode、粘贴和超宽内容内部滚动；React TUI 的搜索框、command palette、未来 composer 也必须支持宽度限制、粘贴折叠和光标不撑破布局。
- Glamour 的 markdown 渲染能力适合保留在 Go 边界：挑战说明、hint、reference 可以先由 Go 子命令输出预处理后的 ANSI-safe text，React TUI 再负责排版；如果 React 端渲染 markdown，也必须用同等的宽度感知渲染器。
- BubbleZone 是鼠标区域管理参考，但 LinuxLab 第一优先级仍是键盘完整可用。鼠标点击 sidebar/list 可以作为增强项，不得成为唯一交互入口。
- Harmonica 的动画理念只适合很克制地用于 progress/spinner，不能让训练工具出现持续大面积动画；低帧率、可关闭、不得影响输入延迟。

来源：

- https://github.com/charmbracelet/bubbles
- https://github.com/charmbracelet/harmonica
- https://github.com/lrstanley/bubblezone
- https://github.com/charmbracelet/glamour

## 3. 产品目标

### 3.1 总目标

把 LinuxLab React TUI 重构为一个响应式、可扩展、现代 agent-style 终端界面：

- 在不同 terminal 尺寸下自动选择布局。
- 先保留 Go 执行引擎，React 负责显示和交互编排。
- 用组件化架构承载后续挑战执行、命令速查、题目搜索、帮助、提示和结果复盘。
- 达到 Codex/opencode/Hermes 类工具的结构清晰度，而不是仅有静态菜单。

### 3.2 成功指标

- 在 `60x20`、`80x24`、`100x30`、`140x40` 四种 terminal 尺寸下无水平溢出。
- 所有主要屏幕在 resize 后重新布局，不需要重启。
- 主菜单、模块列表、挑战列表、详情、命令速查、帮助、结果页均有 render-to-string snapshot 覆盖。
- React TUI 能通过 Go 子命令启动真实挑战，并在退出后显示检测结果。
- React TUI 不复制 Docker/Vim/verify 业务逻辑，只调用 Go 执行边界。
- `make react-tui`、`make test-react`、`make typecheck-react`、`go test ./... -v` 通过。

## 4. 非目标

- 不重写 Docker、Compose、LocalSandbox、Vim runner 和 verify 引擎。
- 不把项目变成纯 Node/TypeScript 项目。
- 不做鼠标优先 UI；鼠标支持可以有，但键盘必须完整可用。
- 不依赖 Claude Code 泄露源码。
- 不在第一阶段实现完整 agent chat；LinuxLab 仍然是训练工具，不是 AI coding agent。

## 5. 用户与场景

### 5.1 主要用户

开发者或运维学习者，在真实 terminal 中练习 Linux/Vim/Shell/运维/容器能力。

### 5.2 场景

- 用户在宽屏 terminal 中希望左侧看模块/进度，右侧看题目详情。
- 用户在 VS Code/JetBrains 内嵌 terminal 中只有 80 列左右，界面必须自动折叠。
- 用户在小窗口中查看命令速查，搜索框不能挤爆布局。
- 用户完成挑战后，需要立即重试、下一题、查看失败原因。
- 用户粘贴命令、路径、长文本时，界面不能乱跳。

## 6. 响应式布局规格

### 6.1 断点

基于 terminal columns/rows，而不是 CSS viewport。

| 名称 | 条件 | 布局 |
|---|---|---|
| Compact | `columns < 80` 或 `rows < 24` | 单栏布局；顶部 tabs；隐藏 sidebar；footer 只显示核心快捷键 |
| Standard | `80 <= columns < 120` 且 `rows >= 24` | 左侧窄导航 + 右侧主工作区；隐藏非关键 meta |
| Wide | `columns >= 120` 且 `rows >= 30` | 左侧导航 + 中间主内容 + 右侧 inspector |
| Tall | `rows >= 40` | 列表显示更多行；详情页显示更多提示/引用 |

### 6.2 布局槽位

所有屏幕共用以下 shell：

```text
┌ Header: product, current route, status badges ┐
├ Navigation: sidebar or top tabs                ┤
├ Main: current screen content                   ┤
├ Inspector: wide-only contextual panel          ┤
├ Footer: key hints, notices, runtime state      ┤
└ ModalLayer: command palette, help, confirm      ┘
```

### 6.3 Header

Header 必须一行内完成，窄屏裁剪右侧状态。

内容优先级：

1. `LinuxLab`
2. 当前路径，例如 `模块 / Linux 基础 / ls 基础`
3. 运行状态，例如 `Docker ok`、`Local fallback`、`React preview`
4. 总进度，例如 `12/278`

### 6.4 Navigation

Standard/Wide：

- 左侧 sidebar 宽度随 terminal 调整，范围 `18-28 columns`。
- 只显示 5 个主入口：总览、开始练习、推荐、命令速查、能力图谱。
- 每项一行，选中项使用 `›` + cyan，不使用背景块作为唯一焦点。

Compact：

- 改为顶部 tabs：`[1 总览] [2 练习] [3 速查]`。
- 不显示描述文案。

### 6.5 Main

Main 区域负责当前屏幕：

- 总览：训练摘要、建议下一步、最近失败题、快捷命令。
- 模块：模块列表、进度、状态。
- 挑战列表：题目列表、难度、状态、标签；支持滚动。
- 详情：题目描述、提示、相关命令、启动按钮/快捷键。
- 命令速查：搜索输入、结果列表、示例预览。
- 结果：通过/失败、检查项、重试/下一题。

所有 Main 内容必须使用 `fitText()` / `wrapText()`，不得直接渲染未限制宽度的长字符串。

### 6.6 Inspector

Wide-only 右侧栏，宽度 `28-40 columns`。

可显示：

- 当前题目的标签、验证规则、提示数量。
- 当前模块的弱项。
- 命令速查的例子。
- Docker/Compose 状态。

Compact/Standard 不显示 Inspector，其内容进入 modal 或详情页下半部分。

### 6.7 Footer

Footer 固定在底部，不因列表滚动而移动。

内容分两类：

- 左：runtime notice，例如 `Docker ready`、`已保存进度`、`检测失败`
- 右：当前屏幕快捷键，例如 `↑↓ 选择 · Enter 进入 · / 搜索 · ? 帮助`

Compact 下只保留最关键三项。

## 7. 交互规格

### 7.1 全局快捷键

| 快捷键 | 行为 |
|---|---|
| `q` | 返回上一级；在总览退出 |
| `Esc` | 关闭 modal/search；无 modal 时返回 |
| `?` | 打开帮助 overlay |
| `Cmd/Ctrl+C` | 安全退出 |
| `1-5` | 跳转主导航 |
| `/` | 当前屏幕搜索；总览中打开 command palette |
| `g/G` | 列表首尾 |
| `PgUp/PgDn` | 列表翻页 |

### 7.2 Command Palette

使用 `/` 或 `Cmd/Ctrl+K` 打开。

第一期命令：

- `开始练习`
- `命令速查`
- `继续上一题`
- `查看失败题`
- `打开帮助`
- `切换主题`

### 7.3 搜索

模块、挑战、命令速查都应支持搜索。

- 搜索输入固定在列表上方。
- 搜索结果为空时显示可操作空状态。
- Esc 第一次清空搜索，第二次返回。

### 7.4 挑战启动

React TUI 不直接实现 Docker/Vim；它调用 Go 子命令。

目标命令：

```bash
linuxlab challenge run <challenge-id> --json
```

输出 NDJSON：

```json
{"type":"setup","message":"创建 Docker 沙盒"}
{"type":"handoff","mode":"docker","message":"进入容器，exit 后自动检测"}
{"type":"result","passed":true,"hintsUsed":0,"results":[{"passed":true,"message":"脚本检测通过"}]}
```

React TUI 启动子进程后进入 handoff 状态：

- Go 子进程接管终端。
- 退出后 React TUI 恢复 alternate screen。
- 结果页显示检测结果。

### 7.5 Modal

用于：

- 帮助
- 命令面板
- 确认退出挑战
- Docker 不可用说明
- 检测失败详情

Modal 必须适配宽度：

- Compact：全屏。
- Standard/Wide：居中 panel，最大宽度 72。

## 8. 视觉语言

### 8.1 颜色语义

遵循 Codex 式语义色：

- Cyan：焦点、输入、可操作状态。
- Green：通过、成功、已完成。
- Red：失败、错误、危险。
- Yellow/Amber：提示、警告、推荐。
- Dim：次要元信息。

避免：

- 大面积彩色背景。
- 紫蓝渐变。
- 多种近似颜色表达同一种状态。
- 在未知终端背景上依赖黑/白前景。

### 8.2 密度

LinuxLab 是训练工具，应偏高密度：

- 列表一行一个实体。
- 描述放 inspector 或详情页，不在列表里堆长文。
- 宽屏多列，窄屏单列。
- 只在关键状态使用边框。

### 8.3 空状态

空状态必须说明下一步：

- 没有推荐：提示“完成或失败几道题后会生成推荐”，并提供“开始练习”。
- 搜索无结果：提示“换关键词或按 Esc 清空”。
- Docker 不可用：提示“将使用本地模式；容器题不可运行”。

## 9. 技术架构

### 9.1 前端运行时

继续使用 Ink + React 19 + TypeScript。

原因：

- 已经接入项目。
- 对当前 Node 环境可用。
- `renderToString()` 适合 snapshot 测试。
- 组件模型满足需求。

注意：

- 当前 `tsx -e` 与 Ink/yoga-layout ESM top-level await 不兼容；运行入口必须用脚本文件或编译后 ESM，不能用 `tsx -e` 做正式验证。

### 9.2 推荐目录结构

```text
src/react-tui/
  index.tsx                 # 启动入口
  App.tsx                   # provider 装配，不放复杂 UI
  app/
    state.ts                # AppState, reducers/actions
    breakpoints.ts          # terminal size -> layout mode
    navigation.ts           # route stack, selected ids
  data/
    load.ts                 # 题库/进度/速查读取
    normalize.ts            # YAML/JSON -> domain model
  domain/
    types.ts                # Challenge, Category, Progress, Result
    selectors.ts            # progress summaries, recommendations
  components/
    shell/
      AppShell.tsx
      Header.tsx
      Navigation.tsx
      Footer.tsx
      Inspector.tsx
    common/
      TextFit.tsx
      ProgressBar.tsx
      StatusBadge.tsx
      ScrollList.tsx
      EmptyState.tsx
      Modal.tsx
      KeyHints.tsx
    screens/
      HomeScreen.tsx
      ModulesScreen.tsx
      ChallengeListScreen.tsx
      ChallengeDetailScreen.tsx
      ReferenceScreen.tsx
      ResultScreen.tsx
  runtime/
    terminal.ts             # useStdout size, resize event bridge
    goExecutor.ts           # spawn Go child process
  theme/
    theme.ts                # semantic colors
    render.ts               # color helpers
  test/
    fixtures.ts
    render.tsx
```

### 9.3 Go 边界

新增 Go 子命令，而不是让 React 调内部包。

目标：

```bash
linuxlab challenge run <challenge-id> --json
linuxlab data dump --json
linuxlab doctor --json
```

第一期只必须有：

- `linuxlab challenge run <challenge-id> --json`

后续可把 React 数据加载也改成 `linuxlab data dump --json`，避免 TS/YAML 与 Go/YAML 行为差异。

### 9.4 组件能力映射

React/Ink 重构不直接使用 Bubbles 等 Go 包，但必须达到这些组件的行为能力。

| 能力 | Go 生态参考 | React/Ink 目标组件 |
|---|---|---|
| 纵向滚动内容 | Bubbles `viewport` | `components/common/Viewport.tsx` |
| 可搜索列表 | Bubbles `list` | `components/common/ScrollList.tsx` + `app/search.ts` |
| 快捷键与帮助同步 | Bubbles `key` / `help` | `app/keymap.ts` + `components/common/KeyHints.tsx` |
| 单行输入 | Bubbles `textinput` | `components/common/SearchInput.tsx` |
| 多行输入/粘贴 | Bubbles `textarea` | 后续 `components/common/Composer.tsx` |
| Markdown/说明渲染 | Glamour | Go 预渲染或 React 宽度感知 markdown renderer |
| 鼠标点击区域 | BubbleZone | 后续 screen item hit zones；键盘优先 |
| 轻量动画 | Harmonica | spinner/progress，低频、可关闭 |

## 10. 测试策略

### 10.1 React 单元测试

使用 Vitest + Ink `renderToString()`：

- breakpoint 计算测试。
- selectors 测试。
- 每个 screen 的 compact/standard/wide snapshot。
- 空状态测试。
- 搜索过滤测试。

### 10.2 Go 测试

- 子命令解析测试。
- `challenge run --json` 输出事件测试。
- Docker 不可用时 skip 或 fallback。
- 原有 `go test ./... -v` 必须持续通过。

### 10.3 手动验证矩阵

```bash
COLUMNS=60 LINES=20 npm run tui:react
COLUMNS=80 LINES=24 npm run tui:react
COLUMNS=100 LINES=30 npm run tui:react
COLUMNS=140 LINES=40 npm run tui:react
```

需要检查：

- 无水平溢出。
- footer 始终可见。
- 搜索输入不覆盖列表。
- 详情描述换行。
- modal 在窄屏和宽屏都合理。

## 11. 分阶段交付

### Phase 1：响应式 shell

- 引入 terminal size manager。
- 实现 compact/standard/wide layout。
- 拆出 AppShell、Header、Navigation、Footer。
- 建立 snapshot 测试。

### Phase 2：列表和详情组件

- ScrollList。
- Viewport。
- TextFit/wrap。
- ChallengeList、Reference、Detail 自适应。
- EmptyState。
- Keymap 驱动的 KeyHints。

### Phase 3：命令面板和帮助 overlay

- ModalLayer。
- CommandPalette。
- HelpOverlay。
- 搜索状态管理。

### Phase 4：Go 执行边界

- 新增 `linuxlab challenge run <id> --json`。
- React 调子进程执行挑战。
- 结果页接入真实验证结果。

### Phase 5：主题和 polish

- 语义主题。
- 终端背景探测或保守 ANSI 颜色。
- Paste handling。
- 宽屏 inspector。

## 12. 风险

| 风险 | 影响 | 缓解 |
|---|---|---|
| Ink/yoga ESM 运行限制 | 启动失败 | 使用脚本入口和 ESM 配置；禁止 `tsx -e` 作为验证 |
| TS YAML 与 Go YAML 行为不同 | 数据加载不一致 | 长期改为 Go `data dump --json`；短期保留 sanitizer |
| React 与 Go 执行状态同步复杂 | 结果页不稳定 | 用 NDJSON 事件协议，先少量事件 |
| 小 terminal 太窄 | UI 不可用 | Compact 模式最低支持 60 列；低于 60 显示 unsupported screen |
| 复制/选择被 alternate screen 影响 | 用户体验差 | 模态/选择背景保守，支持退出后保留结果摘要 |

## 13. 验收标准

- `npm run test:react` 通过。
- `npm run typecheck:react` 通过。
- `go test ./... -v` 通过。
- `go build -buildvcs=false -o /tmp/linuxlab-build-check ./cmd/linuxlab` 通过。
- 4 个 terminal 尺寸 snapshot 更新并审查。
- `make react-tui` 在当前 terminal 中不溢出。
- React TUI 能从详情页启动至少一个 Linux 基础题并显示真实检测结果。
