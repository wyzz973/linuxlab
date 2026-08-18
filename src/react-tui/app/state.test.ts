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

	it('updates search query and resets searchable cursors in one reducer pass', () => {
		const state = {
			...createInitialState({initialScreen: 'challenges'}),
			moduleCursor: 3,
			challengeCursor: 5,
			referenceCursor: 2,
		};
		const next = reduceAppState(state, {type: 'setSearchQuery', query: 'grep'});
		expect(next.query).toBe('grep');
		expect(next.searchActive).toBe(true);
		expect(next.moduleCursor).toBe(0);
		expect(next.challengeCursor).toBe(0);
		expect(next.referenceCursor).toBe(0);
	});

	it('clears stale search when navigating to another screen', () => {
		const state = createInitialState({initialScreen: 'modules', initialQuery: 'grep'});
		const next = reduceAppState(state, {type: 'navigate', screen: 'reference'});
		expect(next.screen).toBe('reference');
		expect(next.query).toBe('');
		expect(next.searchActive).toBe(false);
	});

	it('opens and clears explicit search mode', () => {
		const opened = reduceAppState(createInitialState({initialScreen: 'reference'}), {type: 'openSearch'});
		expect(opened.searchActive).toBe(true);
		expect(opened.notice).toContain('搜索模式');
		const cleared = reduceAppState(opened, {type: 'clearQuery'});
		expect(cleared.searchActive).toBe(false);
		expect(cleared.query).toBe('');
	});

	it('does not push result screen into history when retry finishes from result', () => {
		const state = {
			...createInitialState({initialScreen: 'result'}),
			history: ['detail' as const],
			runningChallengeID: 'ls-basic',
		};
		const next = reduceAppState(state, {
			type: 'finishChallenge',
			result: {
				challengeID: 'ls-basic',
				passed: true,
				hintsUsed: 0,
				results: [],
				events: [],
			},
		});
		expect(next.screen).toBe('result');
		expect(next.history).toEqual(['detail']);
	});

	it('reveals hints progressively and caps at the hint count', () => {
		let state = createInitialState({initialScreen: 'detail'});
		expect(state.hintLevel).toBe(0);
		state = reduceAppState(state, {type: 'revealHint', max: 3});
		state = reduceAppState(state, {type: 'revealHint', max: 3});
		expect(state.hintLevel).toBe(2);
		state = reduceAppState(state, {type: 'revealHint', max: 3});
		state = reduceAppState(state, {type: 'revealHint', max: 3});
		expect(state.hintLevel).toBe(3);
	});

	it('resets hint level when selecting another challenge', () => {
		const state = {...createInitialState({initialScreen: 'detail'}), hintLevel: 2};
		const next = reduceAppState(state, {type: 'selectChallenge', challengeID: 'grep-basic'});
		expect(next.selectedChallengeID).toBe('grep-basic');
		expect(next.hintLevel).toBe(0);
	});

	it('restores hint level for retry from result', () => {
		const state = {...createInitialState({initialScreen: 'result'}), hintLevel: 0};
		const next = reduceAppState(state, {type: 'setHintLevel', level: 2});
		expect(next.hintLevel).toBe(2);
		expect(reduceAppState(next, {type: 'setHintLevel', level: -1}).hintLevel).toBe(0);
	});

	it('keeps hint level when starting and finishing a challenge', () => {
		const state = {...createInitialState({initialScreen: 'detail'}), hintLevel: 2};
		const started = reduceAppState(state, {type: 'startChallenge', challengeID: 'ls-basic'});
		expect(started.runningChallengeID).toBe('ls-basic');
		expect(started.hintLevel).toBe(2);
		const finished = reduceAppState(started, {
			type: 'finishChallenge',
			result: {
				challengeID: 'ls-basic',
				passed: true,
				hintsUsed: 2,
				results: [],
				events: [],
			},
		});
		expect(finished.hintLevel).toBe(2);
	});
});
