import type {ChallengeRunResult, ModalID, ScreenID} from '../types.js';

export type CursorName = 'navCursor' | 'moduleCursor' | 'challengeCursor' | 'referenceCursor' | 'commandCursor';

export type AppState = {
	screen: ScreenID;
	history: ScreenID[];
	modal?: ModalID;
	query: string;
	searchActive: boolean;
	navCursor: number;
	moduleCursor: number;
	challengeCursor: number;
	referenceCursor: number;
	commandCursor: number;
	selectedCategoryID: string;
	selectedChallengeID: string;
	notice: string;
	hintLevel: number;
	runningChallengeID?: string;
	lastResult?: ChallengeRunResult;
};

export type AppAction =
	| {type: 'navigate'; screen: ScreenID}
	| {type: 'back'}
	| {type: 'openModal'; modal: ModalID}
	| {type: 'closeModal'}
	| {type: 'setQuery'; query: string}
	| {type: 'setSearchQuery'; query: string}
	| {type: 'openSearch'}
	| {type: 'clearQuery'}
	| {type: 'setNotice'; notice: string}
	| {type: 'setCursor'; name: CursorName; value: number}
	| {type: 'selectCategory'; categoryID: string; firstChallengeID: string}
	| {type: 'selectChallenge'; challengeID: string}
	| {type: 'revealHint'; max: number}
	| {type: 'setHintLevel'; level: number}
	| {type: 'startChallenge'; challengeID: string}
	| {type: 'finishChallenge'; result: ChallengeRunResult}
	| {type: 'failChallenge'; challengeID: string; message: string};

type InitialOptions = {
	initialScreen?: ScreenID;
	initialQuery?: string;
};

export function createInitialState(options: InitialOptions = {}): AppState {
	return {
		screen: options.initialScreen ?? 'menu',
		history: [],
		query: options.initialQuery ?? '',
		searchActive: (options.initialQuery ?? '') !== '',
		navCursor: 0,
		moduleCursor: 0,
		challengeCursor: 0,
		referenceCursor: 0,
		commandCursor: 0,
		selectedCategoryID: '',
		selectedChallengeID: '',
		notice: '',
		hintLevel: 0,
	};
}

export function reduceAppState(state: AppState, action: AppAction): AppState {
	switch (action.type) {
		case 'navigate':
			return {
				...state,
				screen: action.screen,
				history: state.screen === action.screen ? state.history : [...state.history, state.screen],
				modal: undefined,
				query: '',
				searchActive: false,
				notice: '',
			};
		case 'back': {
			const previous = state.history.at(-1);
			if (!previous) {
				return state.screen === 'menu' ? state : {...state, screen: 'menu', history: [], modal: undefined, searchActive: false, notice: ''};
			}
			return {...state, screen: previous, history: state.history.slice(0, -1), modal: undefined, searchActive: false, notice: ''};
		}
		case 'openModal':
			return {...state, modal: action.modal, query: action.modal === 'commandPalette' ? '' : state.query, commandCursor: 0};
		case 'closeModal':
			return {...state, modal: undefined};
		case 'setQuery':
			return {...state, query: action.query};
		case 'setSearchQuery':
			return {
				...state,
				query: action.query,
				searchActive: true,
				moduleCursor: 0,
				challengeCursor: 0,
				referenceCursor: 0,
			};
		case 'openSearch':
			return {...state, searchActive: true, notice: '搜索模式：输入关键词，Esc 退出'};
		case 'clearQuery':
			return {...state, query: '', searchActive: false, notice: '', referenceCursor: 0, challengeCursor: 0, moduleCursor: 0, commandCursor: 0};
		case 'setNotice':
			return {...state, notice: action.notice};
		case 'setCursor':
			return {...state, [action.name]: Math.max(0, action.value)};
		case 'selectCategory':
			return {
				...state,
				selectedCategoryID: action.categoryID,
				selectedChallengeID: action.firstChallengeID,
				challengeCursor: 0,
			};
		case 'selectChallenge':
			return {...state, selectedChallengeID: action.challengeID, hintLevel: 0};
		case 'revealHint':
			return {...state, hintLevel: Math.min(action.max, state.hintLevel + 1), notice: ''};
		case 'setHintLevel':
			return {...state, hintLevel: Math.max(0, action.level)};
		case 'startChallenge':
			return {...state, runningChallengeID: action.challengeID, notice: '正在启动挑战...'};
		case 'finishChallenge':
			return {
				...state,
				runningChallengeID: undefined,
				lastResult: action.result,
				screen: 'result',
				history: pushHistory(state.history, state.screen, 'result'),
				notice: action.result.passed ? '检测通过' : '检测未通过',
			};
		case 'failChallenge':
			return {
				...state,
				runningChallengeID: undefined,
				lastResult: {
					challengeID: action.challengeID,
					passed: false,
					hintsUsed: 0,
					results: [{passed: false, message: action.message}],
					events: [{type: 'error', message: action.message}],
				},
				screen: 'result',
				history: pushHistory(state.history, state.screen, 'result'),
				notice: '挑战执行失败',
			};
	}
}

function pushHistory(history: ScreenID[], current: ScreenID, next: ScreenID) {
	if (current === next || history.at(-1) === current) {
		return history;
	}
	return [...history, current];
}
