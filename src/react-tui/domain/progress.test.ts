import {describe, expect, it} from 'vitest';
import type {ChallengeRunResult, LinuxLabData} from '../types.js';
import {applyRunResultToData} from './progress.js';

const baseData: LinuxLabData = {
	categories: [
		{
			id: 'linux-basics',
			label: 'Linux 基础命令',
			total: 2,
			passed: 1,
			challenges: [
				{id: 'ls-basic', title: 'ls 基础', difficulty: 1, category: 'linux-basics', subcategory: 'files', tags: ['ls'], description: '列出文件', hints: [], verify: []},
				{id: 'grep-basic', title: 'grep 搜索', difficulty: 2, category: 'linux-basics', subcategory: 'text', tags: ['grep'], description: '搜索文本', hints: [], verify: []},
			],
		},
	],
	progress: {
		skills: {},
		challenges: {
			'ls-basic': {status: 'passed', attempts: 1, hints_used: 0, last_attempt: '2026-05-09T00:00:00Z'},
		},
	},
	references: {commands: []},
};

function runResult(challengeID: string, passed: boolean, hintsUsed: number): ChallengeRunResult {
	return {challengeID, passed, hintsUsed, results: [], events: []};
}

describe('applyRunResultToData', () => {
	it('records a new passed attempt with hints used', () => {
		const next = applyRunResultToData(baseData, runResult('grep-basic', true, 2));
		expect(next.progress.challenges['grep-basic']?.status).toBe('passed');
		expect(next.progress.challenges['grep-basic']?.attempts).toBe(1);
		expect(next.progress.challenges['grep-basic']?.hints_used).toBe(2);
		expect(next.progress.challenges['grep-basic']?.last_attempt).toBeTruthy();
		expect(next.categories[0]?.passed).toBe(2);
	});

	it('increments attempts for a repeated challenge', () => {
		const next = applyRunResultToData(baseData, runResult('ls-basic', false, 1));
		expect(next.progress.challenges['ls-basic']?.status).toBe('failed');
		expect(next.progress.challenges['ls-basic']?.attempts).toBe(2);
		expect(next.progress.challenges['ls-basic']?.hints_used).toBe(1);
		expect(next.categories[0]?.passed).toBe(0);
	});

	it('leaves unrelated challenges untouched', () => {
		const next = applyRunResultToData(baseData, runResult('grep-basic', true, 0));
		expect(next.progress.challenges['ls-basic']).toEqual(baseData.progress.challenges['ls-basic']);
		expect(next.categories).toHaveLength(1);
	});

	it('is pure: does not mutate the input data', () => {
		const before = JSON.stringify(baseData);
		applyRunResultToData(baseData, runResult('grep-basic', true, 0));
		expect(JSON.stringify(baseData)).toBe(before);
	});
});
