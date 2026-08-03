import React from 'react';
import {renderToString} from 'ink';
import {describe, expect, it} from 'vitest';
import {createLayoutSpec} from '../../app/breakpoints.js';
import {fixtureData} from '../../test/fixtures.js';
import {filterChallenges, filterReferences} from '../../domain/selectors.js';
import {ChallengeDetailScreen} from './ChallengeDetailScreen.js';
import {ChallengeListScreen} from './ChallengeListScreen.js';
import {HomeScreen} from './HomeScreen.js';
import {ModulesScreen} from './ModulesScreen.js';
import {ReferenceScreen} from './ReferenceScreen.js';

describe('responsive screens', () => {
	it('renders home at compact width', () => {
		const layout = createLayoutSpec({columns: 60, rows: 20});
		const output = renderToString(<HomeScreen data={fixtureData} layout={layout} />);
		expect(output).toContain('训练控制台');
	});

	it('renders module progress and selected marker', () => {
		const layout = createLayoutSpec({columns: 80, rows: 24});
		const output = renderToString(<ModulesScreen categories={fixtureData.categories} cursor={0} query="" layout={layout} />);
		expect(output).toContain('整体进度');
		expect(output).toContain('›');
	});

	it('limits visible challenge rows', () => {
		const layout = {...createLayoutSpec({columns: 100, rows: 30}), mainHeight: 7};
		const category = fixtureData.categories[0]!;
		const output = renderToString(<ChallengeListScreen category={category} cursor={2} query="" progress={fixtureData.progress.challenges} layout={layout} />);
		expect(output).toContain('chmod 权限');
	});

	it('wraps long challenge description', () => {
		const layout = createLayoutSpec({columns: 80, rows: 24});
		const challenge = fixtureData.categories[0]!.challenges[0]!;
		const output = renderToString(<ChallengeDetailScreen challenge={challenge} relatedCommands={[]} progress={fixtureData.progress.challenges} layout={layout} />);
		expect(output).toContain('任务描述');
	});

	it('renders reference empty state after filtering', () => {
		const layout = createLayoutSpec({columns: 80, rows: 24});
		const commands = filterReferences(fixtureData.references.commands, 'not-found');
		const output = renderToString(<ReferenceScreen commands={commands} cursor={0} query="not-found" layout={layout} />);
		expect(output).toContain('没有匹配的命令');
	});

	it('filters challenges by query', () => {
		expect(filterChallenges(fixtureData.categories[0]!.challenges, 'grep')).toHaveLength(1);
	});
});
