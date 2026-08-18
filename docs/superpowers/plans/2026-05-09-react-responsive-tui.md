# React Responsive TUI Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use `superpowers:executing-plans` to implement this plan task-by-task. If the work is split across agents, use `superpowers:subagent-driven-development` and assign disjoint file ownership before editing.

**Goal:** Rebuild LinuxLab's experimental React/Ink TUI into a responsive, Codex/opencode/Hermes-style terminal interface that adapts to terminal size and delegates challenge execution to the existing Go engine.

> **状态（2026-08-17）：** 本计划的 15 个任务已全部实现并验证（见 `2026-05-09-react-responsive-tui-verification.md`）。接入收尾补充：详情页渐进式提示（`h` 解锁、`--hints N` 透传、重试保留）、挑战结束后的进度即时刷新（乐观更新 + Go data dump 对账）、Header 后端状态徽标（doctor），并修复了 Compose 挑战的宿主执行语义（绝对 compose 路径、宿主 init/check/交互 shell）。

**Architecture:** React/Ink owns rendering, navigation, responsive layout, search, modals, and result presentation. Go remains the source of truth for challenge loading, Docker/Compose/LocalSandbox/Vim handoff, verification, and progress persistence, exposed through a small CLI boundary. Shared behavior is organized around app state, layout slots, semantic theme tokens, width-aware text utilities, keymaps, and screen-specific selectors.

**Tech Stack:** React 19, Ink 7, TypeScript, Vitest, Node `child_process`, Go, Bubble Tea-compatible runner extraction, existing `internal/challenge`, `internal/sandbox`, `internal/verify`, `internal/progress`, `internal/reference`.

---

## Source Documents

- Product requirements: `docs/superpowers/specs/2026-05-09-react-responsive-tui-prd.md`
- Existing Go TUI: `internal/tui/`
- Existing React prototype: `src/react-tui/`
- Existing commands: `package.json`, `Makefile`
- Build/test rules: `AGENTS.md`

## Non-Negotiable Constraints

- All user-facing strings remain Chinese.
- Do not rewrite Docker, Compose, LocalSandbox, Vim, verification, or progress logic in TypeScript.
- React components must never render unbounded challenge descriptions, command examples, paths, or error output.
- Every major screen must render at compact, standard, and wide sizes in tests.
- The current Go/Bubble Tea TUI must continue to run unless the final default command is intentionally switched after parity is proven.
- Existing user changes in the worktree must not be reverted.

## Target Layout Contract

```text
Header:       one row, always visible
Navigation:   sidebar on standard/wide, tabs on compact
Main:         route-specific content, scroll budget controlled by shell
Inspector:    wide-only contextual panel
Footer:       one row, always visible
ModalLayer:   command palette, help, confirm, error details
```

Responsive modes:

```ts
type LayoutMode = 'unsupported' | 'compact' | 'standard' | 'wide';

// unsupported: columns < 60 or rows < 18
// compact:     columns < 80 or rows < 24
// standard:    columns >= 80 and columns < 120 and rows >= 24
// wide:        columns >= 120 and rows >= 30
```

Manual validation sizes:

```bash
COLUMNS=60 LINES=20 npm run tui:react
COLUMNS=80 LINES=24 npm run tui:react
COLUMNS=100 LINES=30 npm run tui:react
COLUMNS=140 LINES=40 npm run tui:react
```

---

## Task 1: Establish Responsive Breakpoint Model

**Files:**
- Create: `src/react-tui/app/breakpoints.ts`
- Test: `src/react-tui/app/breakpoints.test.ts`
- Modify: `src/react-tui/types.ts`

**Goal:** Replace hard-coded terminal dimensions with a tested layout model.

**Step 1: Add domain types**

Add to `src/react-tui/types.ts`:

```ts
export type TerminalSize = {
	columns: number;
	rows: number;
};

export type LayoutMode = 'unsupported' | 'compact' | 'standard' | 'wide';

export type LayoutSpec = {
	mode: LayoutMode;
	columns: number;
	rows: number;
	contentWidth: number;
	mainWidth: number;
	navWidth: number;
	inspectorWidth: number;
	headerHeight: number;
	footerHeight: number;
	mainHeight: number;
	showSidebar: boolean;
	showInspector: boolean;
	showCompactTabs: boolean;
};
```

**Step 2: Write tests first**

Create `src/react-tui/app/breakpoints.test.ts`:

```ts
import {describe, expect, it} from 'vitest';
import {createLayoutSpec} from './breakpoints.js';

describe('createLayoutSpec', () => {
	it('uses unsupported mode below minimum terminal size', () => {
		const spec = createLayoutSpec({columns: 50, rows: 16});
		expect(spec.mode).toBe('unsupported');
		expect(spec.contentWidth).toBe(50);
		expect(spec.mainHeight).toBeGreaterThanOrEqual(1);
	});

	it('uses compact mode for small integrated terminals', () => {
		const spec = createLayoutSpec({columns: 60, rows: 20});
		expect(spec.mode).toBe('compact');
		expect(spec.showSidebar).toBe(false);
		expect(spec.showCompactTabs).toBe(true);
		expect(spec.showInspector).toBe(false);
		expect(spec.mainWidth).toBeLessThanOrEqual(60);
	});

	it('uses standard mode around 100 columns', () => {
		const spec = createLayoutSpec({columns: 100, rows: 30});
		expect(spec.mode).toBe('standard');
		expect(spec.showSidebar).toBe(true);
		expect(spec.navWidth).toBeGreaterThanOrEqual(18);
		expect(spec.showInspector).toBe(false);
		expect(spec.mainWidth + spec.navWidth).toBeLessThanOrEqual(100);
	});

	it('uses wide mode and exposes inspector budget', () => {
		const spec = createLayoutSpec({columns: 140, rows: 40});
		expect(spec.mode).toBe('wide');
		expect(spec.showSidebar).toBe(true);
		expect(spec.showInspector).toBe(true);
		expect(spec.inspectorWidth).toBeGreaterThanOrEqual(28);
		expect(spec.mainWidth + spec.navWidth + spec.inspectorWidth).toBeLessThanOrEqual(140);
	});
});
```

**Step 3: Run failing test**

```bash
npm run test:react -- src/react-tui/app/breakpoints.test.ts
```

Expected before implementation: failure because `breakpoints.ts` does not exist.

**Step 4: Implement breakpoint calculation**

Create `src/react-tui/app/breakpoints.ts`:

```ts
import type {LayoutSpec, TerminalSize} from '../types.js';

const MIN_COLUMNS = 60;
const MIN_ROWS = 18;
const HEADER_HEIGHT = 1;
const FOOTER_HEIGHT = 1;

function clamp(value: number, min: number, max: number) {
	return Math.max(min, Math.min(max, value));
}

export function createLayoutSpec(size: TerminalSize): LayoutSpec {
	const columns = Math.max(1, Math.floor(size.columns || 80));
	const rows = Math.max(1, Math.floor(size.rows || 24));
	const mainHeight = Math.max(1, rows - HEADER_HEIGHT - FOOTER_HEIGHT - 2);

	if (columns < MIN_COLUMNS || rows < MIN_ROWS) {
		return {
			mode: 'unsupported',
			columns,
			rows,
			contentWidth: columns,
			mainWidth: columns,
			navWidth: 0,
			inspectorWidth: 0,
			headerHeight: HEADER_HEIGHT,
			footerHeight: FOOTER_HEIGHT,
			mainHeight,
			showSidebar: false,
			showInspector: false,
			showCompactTabs: false,
		};
	}

	if (columns < 80 || rows < 24) {
		return {
			mode: 'compact',
			columns,
			rows,
			contentWidth: columns,
			mainWidth: columns,
			navWidth: 0,
			inspectorWidth: 0,
			headerHeight: HEADER_HEIGHT,
			footerHeight: FOOTER_HEIGHT,
			mainHeight,
			showSidebar: false,
			showInspector: false,
			showCompactTabs: true,
		};
	}

	const navWidth = clamp(Math.floor(columns * 0.22), 18, 28);
	const isWide = columns >= 120 && rows >= 30;
	const inspectorWidth = isWide ? clamp(Math.floor(columns * 0.26), 28, 40) : 0;
	const gutters = isWide ? 4 : 2;
	const mainWidth = Math.max(20, columns - navWidth - inspectorWidth - gutters);

	return {
		mode: isWide ? 'wide' : 'standard',
		columns,
		rows,
		contentWidth: columns,
		mainWidth,
		navWidth,
		inspectorWidth,
		headerHeight: HEADER_HEIGHT,
		footerHeight: FOOTER_HEIGHT,
		mainHeight,
		showSidebar: true,
		showInspector: isWide,
		showCompactTabs: false,
	};
}
```

**Step 5: Verify**

```bash
npm run test:react -- src/react-tui/app/breakpoints.test.ts
npm run typecheck-react
```

Expected: both pass.

---

## Task 2: Add Width-Aware Text Utilities

**Files:**
- Create: `src/react-tui/theme/theme.ts`
- Create: `src/react-tui/utils/text.ts`
- Create: `src/react-tui/components/common/TextFit.tsx`
- Test: `src/react-tui/utils/text.test.ts`

**Goal:** Every screen gets safe truncation and wrapping before layout work starts.

**Step 1: Define semantic colors**

Create `src/react-tui/theme/theme.ts`:

```ts
export const theme = {
	accent: '#4cc9c1',
	success: '#7ecf8f',
	warning: '#f3bc5f',
	error: '#ff7a7a',
	dim: '#81919c',
	border: '#33424c',
	text: undefined,
} as const;

export type ThemeColor = keyof typeof theme;
```

**Step 2: Write tests**

Create `src/react-tui/utils/text.test.ts`:

```ts
import {describe, expect, it} from 'vitest';
import {fitText, wrapText} from './text.js';

describe('fitText', () => {
	it('truncates long text with ellipsis inside width', () => {
		expect(fitText('abcdefghijklmnopqrstuvwxyz', 10)).toBe('abcdefghi…');
	});

	it('returns empty string for non-positive width', () => {
		expect(fitText('abc', 0)).toBe('');
	});

	it('keeps short Chinese text unchanged', () => {
		expect(fitText('命令速查', 20)).toBe('命令速查');
	});
});

describe('wrapText', () => {
	it('wraps long text into bounded lines', () => {
		const lines = wrapText('alpha beta gamma delta', 8);
		expect(lines).toEqual(['alpha', 'beta', 'gamma', 'delta']);
	});

	it('splits a single overlong token', () => {
		const lines = wrapText('abcdefghijkl', 5);
		expect(lines).toEqual(['abcde', 'fghij', 'kl']);
	});
});
```

**Step 3: Implement utilities**

Create `src/react-tui/utils/text.ts`:

```ts
export function fitText(input: string, width: number): string {
	const safeWidth = Math.max(0, Math.floor(width));
	if (safeWidth === 0) {
		return '';
	}
	if (input.length <= safeWidth) {
		return input;
	}
	if (safeWidth === 1) {
		return '…';
	}
	return `${input.slice(0, safeWidth - 1)}…`;
}

export function wrapText(input: string, width: number): string[] {
	const safeWidth = Math.max(1, Math.floor(width));
	const words = input.replace(/\s+/g, ' ').trim().split(' ').filter(Boolean);
	if (words.length === 0) {
		return [''];
	}

	const lines: string[] = [];
	let current = '';

	for (const word of words) {
		if (word.length > safeWidth) {
			if (current !== '') {
				lines.push(current);
				current = '';
			}
			for (let index = 0; index < word.length; index += safeWidth) {
				lines.push(word.slice(index, index + safeWidth));
			}
			continue;
		}

		const candidate = current === '' ? word : `${current} ${word}`;
		if (candidate.length <= safeWidth) {
			current = candidate;
		} else {
			lines.push(current);
			current = word;
		}
	}

	if (current !== '') {
		lines.push(current);
	}
	return lines;
}
```

**Step 4: Create React wrapper**

Create `src/react-tui/components/common/TextFit.tsx`:

```tsx
import React from 'react';
import {Text} from 'ink';
import {fitText} from '../../utils/text.js';

type TextFitProps = {
	children: string;
	width: number;
	color?: string;
	bold?: boolean;
	dimColor?: boolean;
};

export function TextFit({children, width, color, bold, dimColor}: TextFitProps) {
	return (
		<Text color={color} bold={bold} dimColor={dimColor}>
			{fitText(children, width)}
		</Text>
	);
}
```

**Step 5: Verify**

```bash
npm run test:react -- src/react-tui/utils/text.test.ts
npm run typecheck-react
```

Expected: both pass.

---

## Task 3: Split Data and Selectors from UI

**Files:**
- Create: `src/react-tui/domain/types.ts`
- Create: `src/react-tui/domain/selectors.ts`
- Create: `src/react-tui/domain/selectors.test.ts`
- Modify: `src/react-tui/types.ts`
- Modify: `src/react-tui/data.ts`

**Goal:** Make screen components render from derived view data, not raw ad hoc calculations.

**Step 1: Move domain types**

Move challenge, progress, reference, category, and result-related types from `src/react-tui/types.ts` to `src/react-tui/domain/types.ts`. Keep `ScreenID`, `TerminalSize`, `LayoutMode`, and `LayoutSpec` in `src/react-tui/types.ts`.

`src/react-tui/types.ts` should re-export domain types:

```ts
export type {
	Category,
	Challenge,
	ChallengeEntry,
	CommandExample,
	CommandRef,
	LinuxLabData,
	ProgressData,
	ReferenceData,
	SkillEntry,
	VerifyRule,
} from './domain/types.js';
```

**Step 2: Add selectors tests**

Create `src/react-tui/domain/selectors.test.ts`:

```ts
import {describe, expect, it} from 'vitest';
import {filterChallenges, filterReferences, summarizeCategories} from './selectors.js';
import type {Category, CommandRef} from './types.js';

const categories: Category[] = [{
	id: 'linux-basics',
	label: 'Linux 基础命令',
	total: 2,
	passed: 1,
	challenges: [
		{id: 'ls-basic', title: 'ls 基础', difficulty: 1, category: 'linux-basics', subcategory: 'files', tags: ['ls'], description: '查看文件', hints: [], verify: []},
		{id: 'grep-basic', title: 'grep 搜索', difficulty: 2, category: 'linux-basics', subcategory: 'text', tags: ['grep'], description: '搜索文本', hints: [], verify: []},
	],
}];

const refs: CommandRef[] = [
	{name: 'ls', brief: '列出文件', examples: []},
	{name: 'grep', brief: '搜索文本', examples: []},
];

describe('domain selectors', () => {
	it('summarizes categories', () => {
		expect(summarizeCategories(categories)).toEqual({total: 2, passed: 1, failed: 0});
	});

	it('filters challenges by title, id, and tag', () => {
		expect(filterChallenges(categories[0].challenges, 'grep')).toHaveLength(1);
		expect(filterChallenges(categories[0].challenges, 'files')).toHaveLength(1);
	});

	it('filters references by command name and brief', () => {
		expect(filterReferences(refs, '列出')).toHaveLength(1);
		expect(filterReferences(refs, 'grep')).toHaveLength(1);
	});
});
```

**Step 3: Implement selectors**

Create `src/react-tui/domain/selectors.ts`:

```ts
import type {Category, Challenge, CommandRef} from './types.js';

export function summarizeCategories(categories: Category[]) {
	return {
		total: categories.reduce((sum, category) => sum + category.total, 0),
		passed: categories.reduce((sum, category) => sum + category.passed, 0),
		failed: 0,
	};
}

export function filterChallenges(challenges: Challenge[], query: string): Challenge[] {
	const normalized = query.trim().toLowerCase();
	if (normalized === '') {
		return challenges;
	}
	return challenges.filter(challenge => {
		const haystack = [
			challenge.id,
			challenge.title,
			challenge.subcategory,
			challenge.description,
			...challenge.tags,
		].join(' ').toLowerCase();
		return haystack.includes(normalized);
	});
}

export function filterReferences(commands: CommandRef[], query: string): CommandRef[] {
	const normalized = query.trim().toLowerCase();
	if (normalized === '') {
		return commands;
	}
	return commands.filter(command => {
		const haystack = [
			command.name,
			command.brief,
			...command.examples.flatMap(example => [example.desc, example.cmd]),
		].join(' ').toLowerCase();
		return haystack.includes(normalized);
	});
}
```

**Step 4: Update imports**

Update `src/react-tui/data.ts` and `src/react-tui/App.tsx` imports so domain types come through `./types.js` or `./domain/types.js` consistently. Prefer `./types.js` for public app-level imports.

**Step 5: Verify**

```bash
npm run test:react -- src/react-tui/domain/selectors.test.ts
npm run test:react
npm run typecheck-react
```

Expected: all pass.

---

## Task 4: Build App State, Routes, and Keymap

**Files:**
- Create: `src/react-tui/app/state.ts`
- Create: `src/react-tui/app/keymap.ts`
- Create: `src/react-tui/app/state.test.ts`
- Modify: `src/react-tui/types.ts`
- Modify: `src/react-tui/App.tsx`

**Goal:** Replace scattered `useState` navigation with a reducer and a single source of truth for key behavior.

**Step 1: Extend screen IDs**

In `src/react-tui/types.ts`, update `ScreenID`:

```ts
export type ScreenID =
	| 'menu'
	| 'modules'
	| 'challenges'
	| 'detail'
	| 'reference'
	| 'recommend'
	| 'skillmap'
	| 'result';

export type ModalID = 'help' | 'commandPalette' | 'confirmExit' | 'errorDetails';
```

**Step 2: Add reducer tests**

Create `src/react-tui/app/state.test.ts`:

```ts
import {describe, expect, it} from 'vitest';
import {createInitialState, reduceAppState} from './state.js';

describe('app state reducer', () => {
	it('navigates to modules and records route history', () => {
		const state = createInitialState({initialScreen: 'menu'});
		const next = reduceAppState(state, {type: 'navigate', screen: 'modules'});
		expect(next.screen).toBe('modules');
		expect(next.history).toEqual(['menu']);
	});

	it('goes back to previous route', () => {
		const state = reduceAppState(createInitialState({initialScreen: 'menu'}), {type: 'navigate', screen: 'modules'});
		const next = reduceAppState(state, {type: 'back'});
		expect(next.screen).toBe('menu');
		expect(next.history).toEqual([]);
	});

	it('opens and closes modals without changing screen', () => {
		const state = createInitialState({initialScreen: 'reference'});
		const opened = reduceAppState(state, {type: 'openModal', modal: 'help'});
		expect(opened.modal).toBe('help');
		const closed = reduceAppState(opened, {type: 'closeModal'});
		expect(closed.screen).toBe('reference');
		expect(closed.modal).toBeUndefined();
	});
});
```

**Step 3: Implement reducer**

Create `src/react-tui/app/state.ts`:

```ts
import type {ModalID, ScreenID} from '../types.js';

export type AppState = {
	screen: ScreenID;
	history: ScreenID[];
	modal?: ModalID;
	query: string;
	navCursor: number;
	moduleCursor: number;
	challengeCursor: number;
	referenceCursor: number;
	selectedCategoryID: string;
	selectedChallengeID: string;
	notice: string;
};

export type AppAction =
	| {type: 'navigate'; screen: ScreenID}
	| {type: 'back'}
	| {type: 'openModal'; modal: ModalID}
	| {type: 'closeModal'}
	| {type: 'setQuery'; query: string}
	| {type: 'setNotice'; notice: string}
	| {type: 'setCursor'; name: 'navCursor' | 'moduleCursor' | 'challengeCursor' | 'referenceCursor'; value: number}
	| {type: 'selectCategory'; categoryID: string; firstChallengeID: string}
	| {type: 'selectChallenge'; challengeID: string};

type InitialOptions = {
	initialScreen?: ScreenID;
	initialQuery?: string;
};

export function createInitialState(options: InitialOptions = {}): AppState {
	return {
		screen: options.initialScreen ?? 'menu',
		history: [],
		query: options.initialQuery ?? '',
		navCursor: 0,
		moduleCursor: 0,
		challengeCursor: 0,
		referenceCursor: 0,
		selectedCategoryID: '',
		selectedChallengeID: '',
		notice: '',
	};
}

export function reduceAppState(state: AppState, action: AppAction): AppState {
	switch (action.type) {
		case 'navigate':
			return {...state, screen: action.screen, history: [...state.history, state.screen], modal: undefined};
		case 'back': {
			const previous = state.history.at(-1);
			if (!previous) {
				return state.screen === 'menu' ? state : {...state, screen: 'menu', history: []};
			}
			return {...state, screen: previous, history: state.history.slice(0, -1), modal: undefined};
		}
		case 'openModal':
			return {...state, modal: action.modal};
		case 'closeModal':
			return {...state, modal: undefined};
		case 'setQuery':
			return {...state, query: action.query};
		case 'setNotice':
			return {...state, notice: action.notice};
		case 'setCursor':
			return {...state, [action.name]: Math.max(0, action.value)};
		case 'selectCategory':
			return {...state, selectedCategoryID: action.categoryID, selectedChallengeID: action.firstChallengeID, challengeCursor: 0};
		case 'selectChallenge':
			return {...state, selectedChallengeID: action.challengeID};
	}
}
```

**Step 4: Define keymap**

Create `src/react-tui/app/keymap.ts`:

```ts
import type {ScreenID} from '../types.js';

export type KeyHint = {
	keys: string;
	label: string;
};

export const globalHints: KeyHint[] = [
	{keys: '↑↓', label: '选择'},
	{keys: 'Enter', label: '确认'},
	{keys: '/', label: '搜索'},
	{keys: '?', label: '帮助'},
	{keys: 'q', label: '返回'},
];

export function hintsForScreen(screen: ScreenID, compact: boolean): KeyHint[] {
	const screenSpecific: Record<ScreenID, KeyHint[]> = {
		menu: [{keys: '1-5', label: '跳转'}],
		modules: [{keys: 'Enter', label: '进入模块'}],
		challenges: [{keys: 'Enter', label: '查看题目'}],
		detail: [{keys: 'Enter', label: '开始挑战'}],
		reference: [{keys: 'Esc', label: '清空'}],
		recommend: [{keys: 'Enter', label: '查看推荐'}],
		skillmap: [{keys: 'g/G', label: '首尾'}],
		result: [{keys: 'r', label: '重试'}, {keys: 'n', label: '下一题'}],
	};
	const hints = [...screenSpecific[screen], ...globalHints];
	return compact ? hints.slice(0, 3) : hints;
}
```

**Step 5: Wire reducer into `App.tsx`**

Modify `src/react-tui/App.tsx` to call `useReducer(reduceAppState, createInitialState(...))`. Keep behavior equivalent before visual refactor. This step should remove most `useState` calls, but it can leave data selection variables as derived constants.

**Step 6: Verify**

```bash
npm run test:react -- src/react-tui/app/state.test.ts
npm run test:react
npm run typecheck-react
```

Expected: all pass.

---

## Task 5: Create Responsive App Shell

**Files:**
- Create: `src/react-tui/runtime/terminal.ts`
- Create: `src/react-tui/components/shell/AppShell.tsx`
- Create: `src/react-tui/components/shell/Header.tsx`
- Create: `src/react-tui/components/shell/Navigation.tsx`
- Create: `src/react-tui/components/shell/Footer.tsx`
- Create: `src/react-tui/components/shell/Inspector.tsx`
- Create: `src/react-tui/components/common/KeyHints.tsx`
- Test: `src/react-tui/components/shell/AppShell.test.tsx`
- Modify: `src/react-tui/App.tsx`

**Goal:** One shared shell controls layout for every screen.

**Step 1: Add terminal hook**

Create `src/react-tui/runtime/terminal.ts`:

```ts
import {useStdout} from 'ink';
import type {TerminalSize} from '../types.js';

export function useTerminalSize(fallback: TerminalSize = {columns: 80, rows: 24}): TerminalSize {
	const {stdout} = useStdout();
	return {
		columns: stdout.columns ?? Number(process.env.COLUMNS) || fallback.columns,
		rows: stdout.rows ?? Number(process.env.LINES) || fallback.rows,
	};
}
```

**Step 2: Add shell test**

Create `src/react-tui/components/shell/AppShell.test.tsx`:

```tsx
import React from 'react';
import {describe, expect, it} from 'vitest';
import {renderToString} from 'ink-testing-library';
import {createLayoutSpec} from '../../app/breakpoints.js';
import {AppShell} from './AppShell.js';

describe('AppShell', () => {
	it('renders compact tabs without sidebar', () => {
		const output = renderToString(
			<AppShell
				layout={createLayoutSpec({columns: 60, rows: 20})}
				screen="menu"
				title="LinuxLab"
				progressText="0/278"
				notice="Docker ready"
			>
				测试内容
			</AppShell>,
		);
		expect(output).toContain('总览');
		expect(output).not.toContain('导航');
	});

	it('renders wide inspector slot', () => {
		const output = renderToString(
			<AppShell
				layout={createLayoutSpec({columns: 140, rows: 40})}
				screen="detail"
				title="LinuxLab"
				progressText="1/278"
				notice="Docker ready"
				inspector={<span>验证规则</span>}
			>
				题目详情
			</AppShell>,
		);
		expect(output).toContain('导航');
		expect(output).toContain('验证规则');
	});
});
```

If `ink-testing-library` is not installed, use the existing project pattern in `src/react-tui/App.test.tsx`. Do not add another renderer if current tests already use Ink's `renderToString()` successfully.

**Step 3: Implement shell components**

`AppShell` responsibilities:

- Render unsupported message for `layout.mode === 'unsupported'`.
- Render `Header` at the top.
- Render compact tabs if `layout.showCompactTabs`.
- Render `Navigation` sidebar if `layout.showSidebar`.
- Render main children with `width={layout.mainWidth}` and fixed `height={layout.mainHeight}`.
- Render `Inspector` only when `layout.showInspector`.
- Render `Footer` always.

The main shell must not contain screen-specific logic.

**Step 4: Replace root fixed width**

In `src/react-tui/App.tsx`, remove:

```tsx
<Box flexDirection="column" width={96} ...>
```

Use:

```tsx
const terminalSize = useTerminalSize();
const layout = createLayoutSpec(terminalSize);

return (
	<AppShell
		layout={layout}
		screen={state.screen}
		title="LinuxLab"
		progressText={`${totals.passed}/${totals.total}`}
		notice={state.notice || statusText(state.screen)}
		inspector={inspector}
	>
		{screenContent}
	</AppShell>
);
```

**Step 5: Verify**

```bash
npm run test:react -- src/react-tui/components/shell/AppShell.test.tsx
npm run test:react
npm run typecheck-react
```

Expected: all pass. `App.tsx` no longer contains `width={96}`.

---

## Task 6: Implement Common TUI Components Inspired by Bubbles

**Files:**
- Create: `src/react-tui/components/common/ScrollList.tsx`
- Create: `src/react-tui/components/common/Viewport.tsx`
- Create: `src/react-tui/components/common/ProgressBar.tsx`
- Create: `src/react-tui/components/common/StatusBadge.tsx`
- Create: `src/react-tui/components/common/EmptyState.tsx`
- Create: `src/react-tui/components/common/SearchInput.tsx`
- Test: `src/react-tui/components/common/common.test.tsx`

**Goal:** Match the component capabilities of Bubbles `list`, `viewport`, `progress`, `textinput`, `key`, and `help` in React/Ink form.

**Step 1: Write tests**

Create `src/react-tui/components/common/common.test.tsx`:

```tsx
import React from 'react';
import {describe, expect, it} from 'vitest';
import {renderToString} from 'ink-testing-library';
import {ProgressBar} from './ProgressBar.js';
import {ScrollList} from './ScrollList.js';
import {Viewport} from './Viewport.js';

describe('common TUI components', () => {
	it('renders progress bar at requested width', () => {
		const output = renderToString(<ProgressBar value={3} total={10} width={10} />);
		expect(output.trim().length).toBe(10);
	});

	it('renders only visible rows in scroll list', () => {
		const output = renderToString(
			<ScrollList
				items={['one', 'two', 'three', 'four']}
				cursor={2}
				height={2}
				width={20}
				renderItem={(item, active) => `${active ? '›' : ' '} ${item}`}
			/>,
		);
		expect(output).toContain('three');
		expect(output).toContain('four');
		expect(output).not.toContain('one');
	});

	it('wraps viewport content within width and height', () => {
		const output = renderToString(<Viewport text="alpha beta gamma delta" width={8} height={2} />);
		expect(output.split('\n')).toHaveLength(2);
	});
});
```

**Step 2: Implement `ProgressBar`**

Use ASCII-safe glyphs:

```tsx
import React from 'react';
import {Text} from 'ink';
import {theme} from '../../theme/theme.js';

type ProgressBarProps = {
	value: number;
	total: number;
	width: number;
};

export function ProgressBar({value, total, width}: ProgressBarProps) {
	const safeWidth = Math.max(1, width);
	const ratio = total <= 0 ? 0 : Math.max(0, Math.min(1, value / total));
	const filled = Math.round(safeWidth * ratio);
	return (
		<Text>
			<Text color={theme.success}>{'█'.repeat(filled)}</Text>
			<Text color={theme.border}>{'░'.repeat(safeWidth - filled)}</Text>
		</Text>
	);
}
```

**Step 3: Implement `ScrollList`**

Required behavior:

- `height` controls visible row count.
- Cursor stays visible by calculating a start index.
- Each row is passed through `fitText(row, width)`.
- Empty list renders `EmptyState`.

**Step 4: Implement `Viewport`**

Required behavior:

- Accepts text or lines.
- Uses `wrapText()` for string content.
- Trims to `height`.
- Adds dim `…` marker on last row if content overflows.

**Step 5: Implement `SearchInput`**

Required behavior:

- Renders `搜索: <query>` within width.
- Shows Chinese empty-query hint text when query is empty.
- Never expands beyond width.

**Step 6: Verify**

```bash
npm run test:react -- src/react-tui/components/common/common.test.tsx
npm run test:react
npm run typecheck-react
```

Expected: all pass.

---

## Task 7: Refactor Screens into Responsive Components

**Files:**
- Create: `src/react-tui/components/screens/HomeScreen.tsx`
- Create: `src/react-tui/components/screens/ModulesScreen.tsx`
- Create: `src/react-tui/components/screens/ChallengeListScreen.tsx`
- Create: `src/react-tui/components/screens/ChallengeDetailScreen.tsx`
- Create: `src/react-tui/components/screens/ReferenceScreen.tsx`
- Create: `src/react-tui/components/screens/RecommendScreen.tsx`
- Create: `src/react-tui/components/screens/SkillMapScreen.tsx`
- Test: `src/react-tui/components/screens/screens.test.tsx`
- Modify: `src/react-tui/App.tsx`

**Goal:** Remove screen rendering functions from `App.tsx`; every screen receives `layout`, `data`, and state-derived props.

**Step 1: Write smoke tests**

Create `src/react-tui/components/screens/screens.test.tsx` using `test/fixtures.ts` from Task 8 if it already exists. If not, create local minimal fixtures in the test.

Required assertions:

- `HomeScreen` at 60 columns includes `训练控制台` and no hard overflow marker from tests.
- `ModulesScreen` at 80 columns includes module progress and selected row marker.
- `ChallengeListScreen` uses `height` to limit visible rows.
- `ChallengeDetailScreen` wraps long description.
- `ReferenceScreen` filters and shows empty state.

**Step 2: Implement `HomeScreen`**

Content:

- title: `训练控制台`
- metrics: total challenges, passed, failed, command refs
- recommended action: `继续练习` or `开始第一题`
- compact mode: metrics in two short rows
- wide mode: metrics plus recent failures preview

**Step 3: Implement `ModulesScreen`**

Use `ScrollList<Category>` and `ProgressBar`.

Row format:

```text
› Linux 基础命令         ███░░░  12/40
  Vim 操作              ░░░░░░   0/20
```

**Step 4: Implement `ChallengeListScreen`**

Use `filterChallenges()` and `ScrollList<Challenge>`.

Row format:

```text
› ✓ ls 基础                     ★    已通过
  ○ grep 搜索                   ★★   未开始
```

Compact mode omits tags and long status text.

**Step 5: Implement `ChallengeDetailScreen`**

Sections:

- title and status
- wrapped description
- difficulty/tags
- related commands
- hint count
- action hint: `Enter 开始挑战`

Wide mode should move tags/verify/hint metadata to inspector via `App.tsx`.

**Step 6: Implement `ReferenceScreen`**

Structure:

- `SearchInput`
- `ScrollList<CommandRef>`
- selected command preview

Compact mode:

- show only command name and brief
- preview goes below list if height allows

Wide mode:

- selected examples go to inspector

**Step 7: Implement `RecommendScreen` and `SkillMapScreen`**

Keep logic conservative:

- `RecommendScreen` shows failed/incomplete items based on progress.
- `SkillMapScreen` groups by category/subcategory and shows progress bars.
- Both must use width-aware rows.

**Step 8: Wire App**

`App.tsx` should become provider/state orchestration:

- load terminal layout
- derive selected category/challenge/reference
- dispatch keyboard actions
- choose screen component
- choose inspector component
- render `AppShell`

**Step 9: Verify**

```bash
npm run test:react -- src/react-tui/components/screens/screens.test.tsx
npm run test:react
npm run typecheck-react
```

Expected: all pass. `src/react-tui/App.tsx` should be under roughly 220 lines after this task.

---

## Task 8: Add Fixtures and Snapshot Matrix

**Files:**
- Create: `src/react-tui/test/fixtures.ts`
- Create: `src/react-tui/test/render.tsx`
- Create: `src/react-tui/App.responsive.test.tsx`
- Modify: `src/react-tui/App.test.tsx`

**Goal:** Codex-style visual regression coverage for terminal sizes.

**Step 1: Create stable fixture data**

Create `src/react-tui/test/fixtures.ts`:

```ts
import type {LinuxLabData} from '../types.js';

export const fixtureData: LinuxLabData = {
	categories: [
		{
			id: 'linux-basics',
			label: 'Linux 基础命令',
			total: 3,
			passed: 1,
			challenges: [
				{id: 'ls-basic', title: 'ls 基础', difficulty: 1, category: 'linux-basics', subcategory: 'files', tags: ['ls'], description: '使用 ls 查看当前目录内容，并理解隐藏文件和详细列表。', hints: [{level: 1, text: '先运行 ls'}], verify: [{type: 'script', path: 'check.sh'}]},
				{id: 'grep-basic', title: 'grep 搜索长标题用于测试裁剪', difficulty: 2, category: 'linux-basics', subcategory: 'text', tags: ['grep'], description: '使用 grep 从文本中搜索指定模式。', hints: [], verify: []},
				{id: 'chmod-basic', title: 'chmod 权限', difficulty: 3, category: 'linux-basics', subcategory: 'permissions', tags: ['chmod'], description: '修改文件权限。', hints: [], verify: []},
			],
		},
	],
	progress: {
		skills: {},
		challenges: {
			'ls-basic': {status: 'passed', attempts: 1, hints_used: 0, last_attempt: '2026-05-09T00:00:00Z'},
		},
	},
	references: {
		commands: [
			{name: 'ls', brief: '列出目录内容', examples: [{desc: '查看详细信息', cmd: 'ls -la'}]},
			{name: 'grep', brief: '搜索文本', examples: [{desc: '搜索 error', cmd: 'grep error app.log'}]},
		],
	},
};
```

**Step 2: Create render helper**

Create `src/react-tui/test/render.tsx`:

```tsx
import React from 'react';
import {renderToString} from 'ink-testing-library';
import type {ReactElement} from 'react';

export function renderTui(element: ReactElement) {
	return renderToString(element).replace(/\x1B\[[0-?]*[ -/]*[@-~]/g, '');
}
```

**Step 3: Add responsive app snapshots**

Create `src/react-tui/App.responsive.test.tsx`:

```tsx
import React from 'react';
import {describe, expect, it} from 'vitest';
import {App} from './App.js';
import {fixtureData} from './test/fixtures.js';
import {renderTui} from './test/render.js';

describe('App responsive rendering', () => {
	it.each([
		['compact', 60, 20],
		['standard', 100, 30],
		['wide', 140, 40],
	])('renders %s layout without obvious overflow', (_name, columns, rows) => {
		process.env.COLUMNS = String(columns);
		process.env.LINES = String(rows);
		const output = renderTui(<App data={fixtureData} />);
		expect(output).toContain('LinuxLab');
		for (const line of output.split('\n')) {
			expect(line.length).toBeLessThanOrEqual(columns + 4);
		}
	});
});
```

The `+ 4` allowance covers ANSI/rendering artifacts that remain after stripping and box border edge cases. If the helper strips all ANSI reliably, tighten this to `columns`.

**Step 4: Verify**

```bash
npm run test:react -- src/react-tui/App.responsive.test.tsx
npm run test:react
npm run typecheck-react
```

Expected: all pass.

---

## Task 9: Implement Modal Layer, Help, and Command Palette

**Files:**
- Create: `src/react-tui/components/common/Modal.tsx`
- Create: `src/react-tui/components/modals/HelpOverlay.tsx`
- Create: `src/react-tui/components/modals/CommandPalette.tsx`
- Create: `src/react-tui/app/commands.ts`
- Test: `src/react-tui/components/modals/modals.test.tsx`
- Modify: `src/react-tui/App.tsx`
- Modify: `src/react-tui/app/state.ts`

**Goal:** Add opencode/Hermes-style overlays without screen jumps.

**Step 1: Define commands**

Create `src/react-tui/app/commands.ts`:

```ts
import type {ScreenID} from '../types.js';

export type TuiCommand = {
	id: string;
	label: string;
	description: string;
	keys?: string;
	target?: ScreenID;
};

export const tuiCommands: TuiCommand[] = [
	{id: 'practice', label: '开始练习', description: '进入模块列表', keys: '2', target: 'modules'},
	{id: 'reference', label: '命令速查', description: '搜索 Linux 常用命令', keys: '4', target: 'reference'},
	{id: 'recommend', label: '查看推荐', description: '根据历史进度查看薄弱题', keys: '3', target: 'recommend'},
	{id: 'skillmap', label: '能力图谱', description: '查看分类掌握度', keys: '5', target: 'skillmap'},
	{id: 'help', label: '打开帮助', description: '查看当前快捷键', keys: '?'},
];
```

**Step 2: Write modal tests**

Required assertions:

- compact modal uses full width.
- standard/wide modal caps width at 72.
- help overlay includes current key hints.
- command palette filters by Chinese label and English command id.

**Step 3: Implement `Modal`**

Behavior:

- Inputs: `layout`, `title`, `children`.
- Compact: `width={layout.contentWidth}`.
- Standard/Wide: `width={Math.min(72, layout.contentWidth - 8)}`.
- Use border only for modal; no nested card inside card.

**Step 4: Implement overlays**

`HelpOverlay`:

- reads `hintsForScreen(screen, layout.mode === 'compact')`
- shows global plus screen-specific actions

`CommandPalette`:

- uses `SearchInput`
- uses `ScrollList`
- Enter dispatches selected command action

**Step 5: Update input handling**

`App.tsx` input rules:

- `?` opens help modal.
- `/` opens command palette on menu; starts search on list/search screens.
- `Esc` closes modal first.
- `Enter` inside command palette runs selected command.

**Step 6: Verify**

```bash
npm run test:react -- src/react-tui/components/modals/modals.test.tsx
npm run test:react
npm run typecheck-react
```

Expected: all pass.

---

## Task 10: Extract Go Challenge Runner

**Files:**
- Create: `internal/runner/runner.go`
- Create: `internal/runner/events.go`
- Create: `internal/runner/runner_test.go`
- Modify: `internal/tui/app.go`

**Goal:** Share challenge execution between Go TUI and React TUI CLI command without duplicating sandbox logic.

**Step 1: Define runner events**

Create `internal/runner/events.go`:

```go
package runner

import "github.com/sd3/linuxlab/internal/verify"

type Event struct {
	Type      string          `json:"type"`
	Message   string          `json:"message,omitempty"`
	Mode      string          `json:"mode,omitempty"`
	Passed    *bool           `json:"passed,omitempty"`
	HintsUsed int             `json:"hintsUsed,omitempty"`
	Results   []verify.Result `json:"results,omitempty"`
}
```

**Step 2: Define runner API**

Create `internal/runner/runner.go` with this public shape:

```go
package runner

import (
	"context"

	"github.com/sd3/linuxlab/internal/challenge"
	"github.com/sd3/linuxlab/internal/reference"
	"github.com/sd3/linuxlab/internal/verify"
)

type Options struct {
	Challenge *challenge.Challenge
	Refs      *reference.ReferenceData
	HintsUsed int
	Emit      func(Event)
}

type Result struct {
	Passed    bool
	Results   []verify.Result
	HintsUsed int
}

func RunInteractive(ctx context.Context, opts Options) (Result, error) {
	if opts.Challenge.Category == "vim" {
		return runVim(ctx, opts)
	}
	return runSandbox(ctx, opts)
}
```

In the same edit, add unexported `runVim(ctx, opts)` and `runSandbox(ctx, opts)` helpers by extracting the existing code from `internal/tui/app.go`:

- `launchVim()` maps to `runVim()`.
- `launchSandbox()` maps to `runSandbox()`.
- `ChallengeResultMsg` maps to `runner.Result`.
- calls that previously returned `func() tea.Msg` now return `Result, error`.
- calls that previously used `tea.ExecProcess` now run the same `exec.Command` with inherited `os.Stdin`, `os.Stdout`, and `os.Stderr`.
- emit `Event{Type: "setup"}` before setup work, `Event{Type: "handoff"}` before interactive shell handoff, and `Event{Type: "result"}` after verification.

**Step 3: Move logic from `internal/tui/app.go`**

Move these behaviors:

- Vim setup file preparation.
- Sandbox creation.
- `init.sh` execution.
- `.bashrc` helper injection.
- interactive shell command execution.
- in-sandbox verification.
- result aggregation.
- sandbox destroy.

Do not move:

- Bubble Tea screen transition logic.
- Progress store persistence.
- Result screen rendering.

**Step 4: Update Go TUI**

In `internal/tui/app.go`, keep `launchChallenge()` returning `tea.Cmd`, but call `runner.RunInteractive()` inside `tea.ExecProcess` compatible flow or inside a command that uses inherited terminal process behavior. The result must still become `ChallengeResultMsg`.

**Step 5: Add runner unit tests**

Create `internal/runner/runner_test.go`:

```go
package runner

import "testing"

func TestEventShapeHasStableJSONFields(t *testing.T) {
	passed := true
	event := Event{Type: "result", Passed: &passed, HintsUsed: 1}
	if event.Type != "result" {
		t.Fatalf("unexpected event type: %s", event.Type)
	}
}
```

Add deeper runner tests after CLI command is in place; avoid launching real Docker in this unit test.

**Step 6: Verify**

```bash
go test ./internal/runner ./internal/tui -v
go test ./... -v
```

Expected: all pass. If Docker is unavailable, Docker-specific tests must skip as existing conventions require.

---

## Task 11: Add Go CLI Subcommands for React Boundary

**Files:**
- Create: `internal/cli/cli.go`
- Create: `internal/cli/challenge_run.go`
- Create: `internal/cli/doctor.go`
- Create: `internal/cli/cli_test.go`
- Modify: `cmd/linuxlab/main.go`

**Goal:** Expose a stable process boundary: `linuxlab challenge run <id> --json`.

**Step 1: Add CLI dispatch**

Create `internal/cli/cli.go`:

```go
package cli

import (
	"context"
	"io"
)

type Env struct {
	Args          []string
	Stdout        io.Writer
	Stderr        io.Writer
	ChallengesDir string
	RefsPath      string
	ProgressPath  string
}

func Run(ctx context.Context, env Env) int {
	if len(env.Args) == 0 {
		return -1 // caller should launch default Bubble Tea TUI
	}
	switch env.Args[0] {
	case "challenge":
		return runChallenge(ctx, env)
	case "doctor":
		return runDoctor(ctx, env)
	default:
		return writeError(env.Stderr, "未知命令: "+env.Args[0])
	}
}
```

Use a small helper `writeError()` to write Chinese errors and return non-zero.

**Step 2: Implement `challenge run`**

Required behavior:

```bash
linuxlab challenge run ls-basic --json
```

Rules:

- Locate challenge by id from `LINUXLAB_CHALLENGES` or default `challenges`.
- Load references from `references/commands.yaml` if available.
- Emit NDJSON if `--json` is present:

```json
{"type":"setup","message":"准备挑战 ls-basic"}
{"type":"handoff","mode":"docker","message":"进入挑战环境，退出后自动检测"}
{"type":"result","passed":true,"hintsUsed":0,"results":[{"Passed":true,"Message":"脚本检测通过"}]}
```

- Without `--json`, print concise Chinese text for direct CLI use.
- Return `0` when runner completes, even if challenge verification fails. Verification failure is data, not CLI crash.
- Return non-zero for missing challenge, loader errors, or runner setup failure.

**Step 3: Implement `doctor --json`**

Minimal first version:

```json
{"docker":true,"challenges":278,"references":78}
```

This supports future header status badges. It must not block on long checks.

**Step 4: Update `cmd/linuxlab/main.go`**

At process start:

```go
if code := cli.Run(context.Background(), cli.Env{
	Args: os.Args[1:],
	Stdout: os.Stdout,
	Stderr: os.Stderr,
	ChallengesDir: challengesDir,
	RefsPath: "references/commands.yaml",
	ProgressPath: progressPath,
}); code >= 0 {
	os.Exit(code)
}
```

Keep the existing Bubble Tea launch path for no args.

**Step 5: Add CLI tests**

Create `internal/cli/cli_test.go` with tests for:

- unknown command returns non-zero and Chinese error.
- `doctor --json` returns parseable JSON.
- missing challenge id returns non-zero.

Do not launch interactive shell in unit tests.

**Step 6: Verify**

```bash
go test ./internal/cli -v
go test ./... -v
go build -buildvcs=false -o /tmp/linuxlab-build-check ./cmd/linuxlab
/tmp/linuxlab-build-check doctor --json
```

Expected: tests pass, build passes, doctor prints JSON.

---

## Task 12: Connect React to Go Runner

**Files:**
- Create: `src/react-tui/runtime/goExecutor.ts`
- Create: `src/react-tui/runtime/goExecutor.test.ts`
- Create: `src/react-tui/components/screens/ResultScreen.tsx`
- Modify: `src/react-tui/domain/types.ts`
- Modify: `src/react-tui/App.tsx`

**Goal:** Pressing Enter on a challenge detail starts the real Go challenge flow and returns to a result screen.

**Step 1: Add result types**

In `src/react-tui/domain/types.ts`:

```ts
export type VerifyResult = {
	Passed?: boolean;
	passed?: boolean;
	Message?: string;
	message?: string;
};

export type ChallengeRunEvent =
	| {type: 'setup'; message: string}
	| {type: 'handoff'; mode?: string; message: string}
	| {type: 'result'; passed: boolean; hintsUsed: number; results: VerifyResult[]}
	| {type: 'error'; message: string};

export type ChallengeRunResult = {
	challengeID: string;
	passed: boolean;
	hintsUsed: number;
	results: VerifyResult[];
	events: ChallengeRunEvent[];
};
```

**Step 2: Implement executor**

Create `src/react-tui/runtime/goExecutor.ts`:

```ts
import {spawn} from 'node:child_process';
import type {ChallengeRunEvent, ChallengeRunResult} from '../types.js';

export type RunChallengeOptions = {
	binaryPath?: string;
	challengeID: string;
	cwd?: string;
	onEvent?: (event: ChallengeRunEvent) => void;
};

export async function runChallenge(options: RunChallengeOptions): Promise<ChallengeRunResult> {
	const binary = options.binaryPath ?? process.env.LINUXLAB_BIN ?? './linuxlab';
	const events: ChallengeRunEvent[] = [];

	return await new Promise((resolve, reject) => {
		const child = spawn(binary, ['challenge', 'run', options.challengeID, '--json'], {
			cwd: options.cwd ?? process.cwd(),
			stdio: ['inherit', 'pipe', 'inherit'],
		});

		let buffer = '';
		let result: ChallengeRunResult | undefined;

		child.stdout.setEncoding('utf8');
		child.stdout.on('data', chunk => {
			buffer += chunk;
			const lines = buffer.split('\n');
			buffer = lines.pop() ?? '';
			for (const line of lines) {
				if (line.trim() === '') {
					continue;
				}
				const event = JSON.parse(line) as ChallengeRunEvent;
				events.push(event);
				options.onEvent?.(event);
				if (event.type === 'result') {
					result = {
						challengeID: options.challengeID,
						passed: event.passed,
						hintsUsed: event.hintsUsed,
						results: event.results,
						events,
					};
				}
			}
		});

		child.on('error', reject);
		child.on('close', code => {
			if (code !== 0 && !result) {
				reject(new Error(`挑战执行失败，退出码 ${code}`));
				return;
			}
			resolve(result ?? {
				challengeID: options.challengeID,
				passed: false,
				hintsUsed: 0,
				results: [{passed: false, message: '未收到检测结果'}],
				events,
			});
		});
	});
}
```

**Step 3: Test parser behavior**

`goExecutor.test.ts` should test event parsing by extracting a pure helper if needed:

```ts
import {describe, expect, it} from 'vitest';
import {parseRunEventLine} from './goExecutor.js';

describe('parseRunEventLine', () => {
	it('parses result events', () => {
		const event = parseRunEventLine('{"type":"result","passed":true,"hintsUsed":0,"results":[]}');
		expect(event.type).toBe('result');
	});
});
```

Add `parseRunEventLine()` to keep JSON parsing testable without spawning a process.

**Step 4: Implement `ResultScreen`**

Show:

- pass/fail headline
- each verify result
- hints used
- actions: `r 重试`, `n 下一题`, `q 返回`

**Step 5: Update `App.tsx`**

On Detail screen Enter:

- set notice `正在启动挑战...`
- call `runChallenge()`
- dispatch result into state
- navigate to `result`

Because Ink input handlers cannot be `async` directly in a clean way, wrap execution in a callback that starts a promise and updates React state when it resolves. Ensure repeated Enter while running is ignored.

**Step 6: Verify**

```bash
npm run test:react -- src/react-tui/runtime/goExecutor.test.ts
npm run test:react
npm run typecheck-react
go build -buildvcs=false -o ./linuxlab ./cmd/linuxlab
LINUXLAB_BIN=./linuxlab npm run tui:react
```

Manual check: open a simple Linux basics challenge detail, press Enter, exit challenge shell, confirm result screen appears.

---

## Task 13: Add Markdown and Reference Rendering Strategy

**Files:**
- Create: `internal/cli/data_dump.go`
- Create: `internal/cli/data_dump_test.go`
- Create: `src/react-tui/data/goData.ts`
- Modify: `src/react-tui/data.ts`
- Modify: `src/react-tui/components/screens/ChallengeDetailScreen.tsx`

**Goal:** Reduce TypeScript YAML drift and prepare Glamour-style markdown rendering without bloating React components.

**Step 1: Add `linuxlab data dump --json`**

CLI output shape:

```json
{
  "categories": [],
  "progress": {"skills": {}, "challenges": {}},
  "references": {"commands": []}
}
```

Use existing Go loaders:

- `challenge.LoadAllByCategory`
- `progress.NewStore`
- `reference.LoadReferences`

**Step 2: Add test**

`internal/cli/data_dump_test.go` should:

- create a temp challenge directory with one `challenge.yaml`
- call `Run(ctx, Env{Args: []string{"data", "dump", "--json"}, ...})`
- assert JSON contains one category and one challenge

Use `t.TempDir()`. Do not use hardcoded temp paths.

**Step 3: Add React Go data loader**

Create `src/react-tui/data/goData.ts`:

```ts
import {execFile} from 'node:child_process';
import {promisify} from 'node:util';
import type {LinuxLabData} from '../types.js';

const execFileAsync = promisify(execFile);

export async function loadLinuxLabDataFromGo(binaryPath = process.env.LINUXLAB_BIN ?? './linuxlab'): Promise<LinuxLabData> {
	const {stdout} = await execFileAsync(binaryPath, ['data', 'dump', '--json'], {
		cwd: process.cwd(),
		maxBuffer: 10 * 1024 * 1024,
	});
	return JSON.parse(stdout) as LinuxLabData;
}
```

**Step 4: Update `data.ts` fallback**

Behavior:

- Try Go data dump if `LINUXLAB_USE_GO_DATA=1`.
- Fall back to existing TypeScript loader if Go data command fails or env is absent.
- Keep `sanitizeYamlEscapes()` until Go data path becomes default.

**Step 5: Markdown/Glamour decision**

Do not add a JS markdown renderer in this task. Record in code comments and PRD that markdown-like challenge text is rendered as wrapped plain text in React, and future rich ANSI rendering should come from Go data dump using Glamour or a dedicated renderer.

**Step 6: Verify**

```bash
go test ./internal/cli -v
go test ./... -v
go build -buildvcs=false -o ./linuxlab ./cmd/linuxlab
LINUXLAB_USE_GO_DATA=1 LINUXLAB_BIN=./linuxlab npm run test:react
npm run typecheck-react
```

Expected: all pass.

---

## Task 14: Polish Interaction Details and Paste Safety

**Files:**
- Create: `src/react-tui/app/input.ts`
- Create: `src/react-tui/app/input.test.ts`
- Modify: `src/react-tui/App.tsx`
- Modify: `src/react-tui/components/common/SearchInput.tsx`
- Modify: `src/react-tui/components/modals/CommandPalette.tsx`

**Goal:** Make keyboard/search behavior consistent and keep pasted text from breaking layout.

**Step 1: Add input helper tests**

Test cases:

- printable character appends to query.
- backspace removes one code unit.
- Esc clears search first, then backs out.
- pasted multi-line input collapses whitespace for search.
- `g` and `G` move cursor to first/last for lists.
- `PgUp` / `PgDn` step by visible height.

**Step 2: Implement input helpers**

Create pure functions:

```ts
export function normalizePastedSearch(input: string): string {
	return input.replace(/\s+/g, ' ').trimStart();
}

export function updateQueryFromInput(query: string, input: string, key: {backspace?: boolean}) {
	if (key.backspace) {
		return query.slice(0, -1);
	}
	return `${query}${normalizePastedSearch(input)}`;
}

export function moveCursor(current: number, count: number, action: 'up' | 'down' | 'first' | 'last' | 'pageUp' | 'pageDown', pageSize: number) {
	if (count <= 0) {
		return 0;
	}
	switch (action) {
		case 'up':
			return Math.max(0, current - 1);
		case 'down':
			return Math.min(count - 1, current + 1);
		case 'first':
			return 0;
		case 'last':
			return count - 1;
		case 'pageUp':
			return Math.max(0, current - pageSize);
		case 'pageDown':
			return Math.min(count - 1, current + pageSize);
	}
}
```

**Step 3: Wire helpers**

`App.tsx` should use input helpers instead of inline cursor math.

**Step 4: Verify**

```bash
npm run test:react -- src/react-tui/app/input.test.ts
npm run test:react
npm run typecheck-react
```

Expected: all pass.

---

## Task 15: Documentation, Make Targets, and Final Verification

**Files:**
- Modify: `Makefile`
- Modify: `README.md`
- Modify: `docs/tui-preview.html`
- Create: `docs/superpowers/plans/2026-05-09-react-responsive-tui-verification.md`

**Goal:** Make the new TUI easy to run, verify, and review.

**Step 1: Update Makefile**

Ensure targets exist:

```makefile
react-tui:
	npm run tui:react

test-react:
	npm run test:react

typecheck-react:
	npm run typecheck-react

verify-react-tui:
	npm run test:react
	npm run typecheck-react
	go test ./... -v
	go build -buildvcs=false -o /tmp/linuxlab-build-check ./cmd/linuxlab
```

If targets already exist, update only what is missing.

**Step 2: Update README**

Add:

```md
### React TUI preview

```bash
npm install
go build -buildvcs=false -o ./linuxlab ./cmd/linuxlab
LINUXLAB_BIN=./linuxlab npm run tui:react
```

The React TUI is responsive to terminal size. It delegates challenge execution to the Go binary and keeps the existing Go/Bubble Tea TUI available.
```

**Step 3: Update HTML preview**

`docs/tui-preview.html` should include compact, standard, and wide mockups. This is a communication artifact only; source of truth remains tests and the TUI.

**Step 4: Create verification note**

Create `docs/superpowers/plans/2026-05-09-react-responsive-tui-verification.md` with:

- exact commands run
- pass/fail output summary
- terminal sizes manually checked
- any known limitations

**Step 5: Run full verification**

```bash
npm run test:react
npm run typecheck-react
go test ./... -v
go build -buildvcs=false -o /tmp/linuxlab-build-check ./cmd/linuxlab
```

Known local caveat:

```bash
make build
```

may fail if Go VCS stamping cannot read the local Git metadata. When that happens, use:

```bash
go build -buildvcs=false -o /tmp/linuxlab-build-check ./cmd/linuxlab
```

and record the caveat in the verification note.

---

## Parallelization Guidance

Safe parallel slices:

- Worker A: Tasks 1, 2, 6, 8 under `src/react-tui/app`, `src/react-tui/utils`, `src/react-tui/components/common`, `src/react-tui/test`.
- Worker B: Tasks 3, 4, 7, 9 under `src/react-tui/domain`, `src/react-tui/components/screens`, `src/react-tui/components/modals`, `src/react-tui/App.tsx`.
- Worker C: Tasks 10, 11, 13 Go boundary under `internal/runner`, `internal/cli`, `cmd/linuxlab/main.go`.
- Integrator: Task 12 and Task 15 after A/B/C complete.

Rules for parallel work:

- Workers must not edit each other's files.
- `App.tsx` should be owned by one worker at a time.
- Go runner extraction must land before React execution integration.
- Every worker runs the smallest relevant tests before handing off.

## Final Acceptance Checklist

- [x] No `width={96}` or equivalent fixed root width remains in React TUI.
- [x] `createLayoutSpec()` covers unsupported, compact, standard, and wide modes.
- [x] Header, navigation, main, inspector, footer, and modal layer are shared shell slots.
- [x] Sidebar collapses to tabs on compact terminals.
- [x] Footer remains visible in every tested screen.
- [x] Lists are height-bounded and cursor-visible.
- [x] Challenge details wrap long descriptions.
- [x] Reference search never grows beyond its container.
- [x] Help and command palette render as modals.
- [x] React TUI can launch at least one real challenge through Go.
- [x] Result screen displays real verification results.
- [x] Go/Bubble Tea TUI still launches with no args.
- [x] `npm run test:react` passes.
- [x] `npm run typecheck-react` passes.
- [x] `go test ./... -v` passes.
- [x] `go build -buildvcs=false -o /tmp/linuxlab-build-check ./cmd/linuxlab` passes.
- [x] Manual checks cover `60x20`, `80x24`, `100x30`, and `140x40`.
