import React, {useEffect, useMemo, useReducer} from 'react';
import {Box, Text, useApp, useInput} from 'ink';
import {createLayoutSpec} from './app/breakpoints.js';
import {filterCommands} from './app/commands.js';
import {moveCursor, updateQueryFromInput} from './app/input.js';
import {createInitialState, reduceAppState} from './app/state.js';
import {
	filterCategories,
	filterChallenges,
	filterReferences,
	recommendChallenges,
	summarizeCategories,
} from './domain/selectors.js';
import {AppShell} from './components/shell/AppShell.js';
import {navItems} from './components/shell/Navigation.js';
import {HelpOverlay} from './components/modals/HelpOverlay.js';
import {CommandPalette} from './components/modals/CommandPalette.js';
import {ChallengeDetailScreen} from './components/screens/ChallengeDetailScreen.js';
import {ChallengeListScreen} from './components/screens/ChallengeListScreen.js';
import {HomeScreen} from './components/screens/HomeScreen.js';
import {ModulesScreen} from './components/screens/ModulesScreen.js';
import {RecommendScreen} from './components/screens/RecommendScreen.js';
import {ReferenceScreen} from './components/screens/ReferenceScreen.js';
import {ResultScreen} from './components/screens/ResultScreen.js';
import {SkillMapScreen} from './components/screens/SkillMapScreen.js';
import {theme} from './theme/theme.js';
import {fitText} from './utils/text.js';
import {runChallenge} from './runtime/goExecutor.js';
import {useTerminalSize} from './runtime/terminal.js';
import type {Category, Challenge, CommandRef, LinuxLabData, ScreenID, TerminalSize} from './types.js';

type AppProps = {
	data: LinuxLabData;
	initialScreen?: ScreenID;
	initialQuery?: string;
	terminalSize?: TerminalSize;
};

export function App({data, initialScreen = 'menu', initialQuery = '', terminalSize}: AppProps) {
	const {exit} = useApp();
	const detectedSize = useTerminalSize();
	const layout = createLayoutSpec(terminalSize ?? detectedSize);
	const [state, dispatch] = useReducer(
		reduceAppState,
		createInitialState({initialScreen, initialQuery}),
	);

	const totals = useMemo(() => summarizeCategories(data.categories), [data.categories]);
	const filteredCategories = useMemo(
		() => filterCategories(data.categories, state.screen === 'modules' ? state.query : ''),
		[data.categories, state.query, state.screen],
	);

	const selectedCategory = findCategory(data.categories, state.selectedCategoryID) ?? data.categories[0];
	const challengeSource = selectedCategory?.challenges ?? [];
	const filteredChallenges = useMemo(
		() => filterChallenges(challengeSource, state.screen === 'challenges' ? state.query : ''),
		[challengeSource, state.query, state.screen],
	);
	const selectedChallenge =
		findChallenge(challengeSource, state.selectedChallengeID) ??
		filteredChallenges[Math.min(state.challengeCursor, Math.max(0, filteredChallenges.length - 1))] ??
		challengeSource[0];

	const filteredRefs = useMemo(
		() => filterReferences(data.references.commands, state.screen === 'reference' || state.modal === 'commandPalette' ? state.query : ''),
		[data.references.commands, state.query, state.screen, state.modal],
	);
	const selectedRef = filteredRefs[Math.min(state.referenceCursor, Math.max(0, filteredRefs.length - 1))];
	const recommendations = useMemo(
		() => recommendChallenges(data.categories, data.progress.challenges),
		[data.categories, data.progress.challenges],
	);
	const relatedCommands = useMemo(
		() => selectedChallenge ? commandsForChallenge(data.references.commands, selectedChallenge) : [],
		[data.references.commands, selectedChallenge],
	);

	useEffect(() => {
		if (!state.runningChallengeID) {
			return;
		}

		const challenge = findChallengeByID(data.categories, state.runningChallengeID);
		if (!challenge) {
			dispatch({type: 'failChallenge', challengeID: state.runningChallengeID, message: `未找到题目: ${state.runningChallengeID}`});
			return;
		}

		let cancelled = false;
		const timer = setTimeout(() => {
			void runChallenge({challengeID: challenge.id})
				.then(result => {
					if (!cancelled) {
						dispatch({type: 'finishChallenge', result});
					}
				})
				.catch(error => {
					if (!cancelled) {
						const message = error instanceof Error ? error.message : String(error);
						dispatch({type: 'failChallenge', challengeID: challenge.id, message});
					}
				});
		}, 0);

		return () => {
			cancelled = true;
			clearTimeout(timer);
		};
	}, [data.categories, state.runningChallengeID]);

	useInput((input, key) => {
		if (key.ctrl && input === 'c') {
			exit();
			return;
		}

		if (state.modal) {
			handleModalInput(input, key);
			return;
		}

		if (input === '?') {
			dispatch({type: 'openModal', modal: 'help'});
			return;
		}

		if (key.escape) {
			if (state.query !== '' || state.searchActive) {
				dispatch({type: 'clearQuery'});
				return;
			}
			dispatch({type: 'back'});
			return;
		}

		if (state.searchActive && isSearchableScreen(state.screen) && (key.backspace || key.delete || isPrintable(input))) {
			handleSearchInput(input, key);
			return;
		}

		if (input === 'q') {
			if (state.screen === 'menu') {
				exit();
				return;
			}
			dispatch({type: 'back'});
			return;
		}

		if (input === '/') {
			if (state.screen === 'menu') {
				dispatch({type: 'openModal', modal: 'commandPalette'});
			} else if (isSearchableScreen(state.screen)) {
				dispatch({type: 'openSearch'});
			}
			return;
		}

		const navNumber = Number(input);
		if (state.screen === 'menu' && Number.isInteger(navNumber) && navNumber >= 1 && navNumber <= navItems.length) {
			const target = navItems[navNumber - 1]?.id;
			if (target) {
				dispatch({type: 'navigate', screen: target});
			}
			return;
		}

		switch (state.screen) {
			case 'modules':
				handleListInput(input, key, filteredCategories.length, state.moduleCursor, 'moduleCursor', Math.max(1, layout.mainHeight - 4), index => {
					const category = filteredCategories[index];
					if (!category) {
						return;
					}
					dispatch({type: 'selectCategory', categoryID: category.id, firstChallengeID: category.challenges[0]?.id ?? ''});
					dispatch({type: 'clearQuery'});
					dispatch({type: 'navigate', screen: 'challenges'});
				});
				break;
			case 'challenges':
				handleListInput(input, key, filteredChallenges.length, state.challengeCursor, 'challengeCursor', Math.max(1, layout.mainHeight - 4), index => {
					const challenge = filteredChallenges[index];
					if (!challenge) {
						return;
					}
					dispatch({type: 'selectChallenge', challengeID: challenge.id});
					dispatch({type: 'clearQuery'});
					dispatch({type: 'navigate', screen: 'detail'});
				});
				break;
			case 'reference':
				handleSearchableListInput(input, key, filteredRefs.length, state.referenceCursor, 'referenceCursor', Math.max(1, layout.mainHeight - 5));
				break;
			case 'recommend':
				handleListInput(input, key, recommendations.length, state.challengeCursor, 'challengeCursor', Math.max(1, layout.mainHeight - 4), index => {
					const challenge = recommendations[index];
					if (!challenge) {
						return;
					}
					dispatch({type: 'selectCategory', categoryID: challenge.category, firstChallengeID: challenge.id});
					dispatch({type: 'selectChallenge', challengeID: challenge.id});
					dispatch({type: 'navigate', screen: 'detail'});
				});
				break;
			case 'skillmap':
				handleListInput(input, key, data.categories.length, state.moduleCursor, 'moduleCursor', Math.max(1, layout.mainHeight - 4), () => {});
				break;
			case 'detail':
				if (key.return && selectedChallenge && !state.runningChallengeID) {
					dispatch({type: 'startChallenge', challengeID: selectedChallenge.id});
				}
				break;
			case 'result':
				if (input === 'r' && selectedChallenge) {
					dispatch({type: 'startChallenge', challengeID: selectedChallenge.id});
				}
				if ((input === 'n' || key.return) && selectedCategory && selectedChallenge) {
					const next = nextChallenge(selectedCategory, selectedChallenge);
					if (next) {
						dispatch({type: 'selectChallenge', challengeID: next.id});
						dispatch({type: 'navigate', screen: 'detail'});
					}
				}
				break;
			case 'menu':
				handleListInput(input, key, navItems.length, state.navCursor, 'navCursor', 5, index => {
					const target = navItems[index]?.id;
					if (target) {
						dispatch({type: 'navigate', screen: target});
					}
				});
				break;
		}
	}, {isActive: state.runningChallengeID === undefined});

	function handleModalInput(
		input: string,
		key: {escape?: boolean; upArrow?: boolean; downArrow?: boolean; return?: boolean; backspace?: boolean; delete?: boolean},
	) {
		if (key.escape || input === 'q') {
			dispatch({type: 'closeModal'});
			return;
		}
		if (state.modal === 'help') {
			return;
		}

		const commands = filterCommands(state.query);
		if (key.upArrow || input === 'k') {
			dispatch({type: 'setCursor', name: 'commandCursor', value: moveCursor(state.commandCursor, commands.length, 'up', 4)});
			return;
		}
		if (key.downArrow || input === 'j') {
			dispatch({type: 'setCursor', name: 'commandCursor', value: moveCursor(state.commandCursor, commands.length, 'down', 4)});
			return;
		}
		if (key.return) {
			const command = commands[Math.min(state.commandCursor, Math.max(0, commands.length - 1))];
			if (command?.id === 'help') {
				dispatch({type: 'openModal', modal: 'help'});
				return;
			}
			if (command?.target) {
				dispatch({type: 'navigate', screen: command.target});
			}
			return;
		}
		if (key.backspace || key.delete || isPrintable(input)) {
			dispatch({type: 'setQuery', query: updateQueryFromInput(state.query, input, key)});
		}
	}

	function handleListInput(
		input: string,
		key: {upArrow?: boolean; downArrow?: boolean; return?: boolean; pageUp?: boolean; pageDown?: boolean; leftArrow?: boolean; rightArrow?: boolean; backspace?: boolean; delete?: boolean},
		count: number,
		cursor: number,
		cursorName: 'navCursor' | 'moduleCursor' | 'challengeCursor' | 'referenceCursor',
		pageSize: number,
		onSelect: (index: number) => void,
	) {
		if (isSearchableScreen(state.screen) && state.searchActive && (key.backspace || key.delete || isPrintable(input))) {
			handleSearchInput(input, key);
			return;
		}
		if (key.upArrow || input === 'k') {
			dispatch({type: 'setCursor', name: cursorName, value: moveCursor(cursor, count, 'up', pageSize)});
			return;
		}
		if (key.downArrow || input === 'j') {
			dispatch({type: 'setCursor', name: cursorName, value: moveCursor(cursor, count, 'down', pageSize)});
			return;
		}
		if (input === 'g') {
			dispatch({type: 'setCursor', name: cursorName, value: moveCursor(cursor, count, 'first', pageSize)});
			return;
		}
		if (input === 'G') {
			dispatch({type: 'setCursor', name: cursorName, value: moveCursor(cursor, count, 'last', pageSize)});
			return;
		}
		if (key.pageUp) {
			dispatch({type: 'setCursor', name: cursorName, value: moveCursor(cursor, count, 'pageUp', pageSize)});
			return;
		}
		if (key.pageDown) {
			dispatch({type: 'setCursor', name: cursorName, value: moveCursor(cursor, count, 'pageDown', pageSize)});
			return;
		}
		if (key.return || key.rightArrow) {
			onSelect(Math.min(cursor, Math.max(0, count - 1)));
		}
	}

	function handleSearchableListInput(
		input: string,
		key: {upArrow?: boolean; downArrow?: boolean; return?: boolean; pageUp?: boolean; pageDown?: boolean; backspace?: boolean; delete?: boolean},
		count: number,
		cursor: number,
		cursorName: 'referenceCursor',
		pageSize: number,
	) {
		handleListInput(input, key, count, cursor, cursorName, pageSize, () => {});
	}

	function handleSearchInput(input: string, key: {backspace?: boolean; delete?: boolean}) {
		dispatch({type: 'setSearchQuery', query: updateQueryFromInput(state.query, input, key)});
	}

	const screenContent = renderScreen({
		data,
		layout,
		state,
		selectedCategory,
		filteredCategories,
		filteredChallenges,
		selectedChallenge,
		filteredRefs,
		selectedRef,
		recommendations,
		relatedCommands,
	});

	const inspector = renderInspector(
		state.screen === 'detail' ? selectedChallenge : undefined,
		state.screen === 'reference' ? selectedRef : undefined,
	);
	const modal = state.modal === 'help'
		? <HelpOverlay layout={layout} screen={state.screen} />
		: state.modal === 'commandPalette'
			? <CommandPalette layout={layout} query={state.query} cursor={state.commandCursor} />
			: undefined;

	return (
		<AppShell
			layout={layout}
			screen={state.screen}
			title="LinuxLab"
			progressText={`${totals.passed}/${totals.total}`}
			notice={state.notice || statusText(state.screen)}
			statusText="React preview"
			inspector={inspector}
			modal={modal}
		>
			{screenContent}
		</AppShell>
	);
}

function renderScreen({
	data,
	layout,
	state,
	selectedCategory,
	filteredCategories,
	filteredChallenges,
	selectedChallenge,
	filteredRefs,
	selectedRef,
	recommendations,
	relatedCommands,
}: {
	data: LinuxLabData;
	layout: ReturnType<typeof createLayoutSpec>;
	state: ReturnType<typeof createInitialState>;
	selectedCategory?: Category;
	filteredCategories: Category[];
	filteredChallenges: Challenge[];
	selectedChallenge?: Challenge;
	filteredRefs: CommandRef[];
	selectedRef?: CommandRef;
	recommendations: Challenge[];
	relatedCommands: CommandRef[];
}) {
	switch (state.screen) {
		case 'menu':
			return <HomeScreen data={data} layout={layout} />;
		case 'modules':
			return <ModulesScreen categories={filteredCategories} cursor={state.moduleCursor} query={state.query} searchActive={state.searchActive} layout={layout} />;
		case 'challenges':
			return selectedCategory
				? <ChallengeListScreen category={{...selectedCategory, challenges: filteredChallenges}} cursor={state.challengeCursor} query={state.query} searchActive={state.searchActive} progress={data.progress.challenges} layout={layout} />
				: <Text color={theme.dim}>暂无模块</Text>;
		case 'detail':
			return selectedChallenge
				? <ChallengeDetailScreen challenge={selectedChallenge} relatedCommands={relatedCommands} progress={data.progress.challenges} layout={layout} running={state.runningChallengeID === selectedChallenge.id} />
				: <Text color={theme.dim}>暂无题目</Text>;
		case 'reference':
			return <ReferenceScreen commands={filteredRefs} selected={selectedRef} cursor={state.referenceCursor} query={state.query} searchActive={state.searchActive} layout={layout} />;
		case 'recommend':
			return <RecommendScreen challenges={recommendations} cursor={state.challengeCursor} progress={data.progress.challenges} layout={layout} />;
		case 'skillmap':
			return <SkillMapScreen categories={data.categories} cursor={state.moduleCursor} layout={layout} />;
		case 'result':
			return <ResultScreen result={state.lastResult} layout={layout} />;
	}
}

function renderInspector(challenge?: Challenge, command?: CommandRef) {
	if (challenge) {
		return (
			<Box flexDirection="column">
				<Text color={theme.dim}>{fitText(`标签: ${challenge.tags.join(', ') || '无'}`, 34)}</Text>
				<Text color={theme.dim}>提示: {challenge.hints.length}</Text>
				<Text color={theme.dim}>检查: {challenge.verify.length}</Text>
				<Text color={theme.dim}>{fitText(`ID: ${challenge.id}`, 34)}</Text>
			</Box>
		);
	}
	if (command) {
		return (
			<Box flexDirection="column">
				<Text color={theme.warning}>{fitText(command.name, 34)}</Text>
				<Text color={theme.dim}>{fitText(command.brief, 34)}</Text>
				{command.examples.slice(0, 4).map(example => (
					<Text key={`${example.desc}-${example.cmd}`}>{fitText(`${example.desc}: ${example.cmd}`, 34)}</Text>
				))}
			</Box>
		);
	}
	return undefined;
}

function statusText(screen: ScreenID) {
	if (screen === 'reference') {
		return '/ 搜索 · ↑↓/j/k 选择 · Esc 清空/返回';
	}
	if (screen === 'detail') {
		return 'Enter 开始挑战 · ? 帮助 · q 返回';
	}
	if (screen === 'result') {
		return 'r 重试 · n 下一题 · q 返回';
	}
	if (screen === 'modules' || screen === 'challenges') {
		return '↑↓/j/k 选择 · Enter 确认 · / 搜索 · ? 帮助';
	}
	return '↑↓/j/k 选择 · Enter 确认 · / 命令 · ? 帮助';
}

function findCategory(categories: Category[], id: string) {
	return categories.find(category => category.id === id);
}

function findChallenge(challenges: Challenge[], id: string) {
	return challenges.find(challenge => challenge.id === id);
}

function findChallengeByID(categories: Category[], id: string) {
	for (const category of categories) {
		const challenge = findChallenge(category.challenges, id);
		if (challenge) {
			return challenge;
		}
	}
	return undefined;
}

function nextChallenge(category: Category, current: Challenge) {
	const index = category.challenges.findIndex(challenge => challenge.id === current.id);
	return index >= 0 ? category.challenges[index + 1] : undefined;
}

function commandsForChallenge(commands: CommandRef[], challenge: Challenge) {
	const tags = new Set(challenge.tags.map(tag => tag.toLowerCase()));
	return commands.filter(command =>
		tags.has(command.name.toLowerCase()) ||
		command.related_challenges?.includes(challenge.id),
	);
}

function isSearchableScreen(screen: ScreenID) {
	return screen === 'modules' || screen === 'challenges' || screen === 'reference';
}

function isPrintable(input: string) {
	return input.length > 0 && input !== '\r' && input !== '\n' && input !== '\u001b';
}
