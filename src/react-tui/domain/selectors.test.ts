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
		expect(filterChallenges(categories[0]!.challenges, 'grep')).toHaveLength(1);
		expect(filterChallenges(categories[0]!.challenges, 'files')).toHaveLength(1);
	});

	it('filters references by command name and brief', () => {
		expect(filterReferences(refs, '列出')).toHaveLength(1);
		expect(filterReferences(refs, 'grep')).toHaveLength(1);
	});
});
