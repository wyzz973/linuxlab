import React from 'react';
import {Box, Text} from 'ink';
import type {Category, LayoutSpec} from '../../types.js';
import {theme} from '../../theme/theme.js';
import {SearchInput} from '../common/SearchInput.js';
import {ScrollList} from '../common/ScrollList.js';
import {compactProgress} from './helpers.js';

type ModulesScreenProps = {
	categories: Category[];
	cursor: number;
	query: string;
	searchActive?: boolean;
	layout: LayoutSpec;
};

export function ModulesScreen({categories, cursor, query, searchActive = false, layout}: ModulesScreenProps) {
	const width = layout.mainWidth;
	const listHeight = Math.max(3, layout.mainHeight - 4);
	const total = categories.reduce((sum, category) => sum + category.total, 0);
	const passed = categories.reduce((sum, category) => sum + category.passed, 0);

	return (
		<Box flexDirection="column" width={width}>
			<Text bold color={theme.accent}>模块</Text>
			<Text color={theme.dim}>整体进度 {compactProgress(passed, total, layout.mode === 'compact' ? 8 : 14)}</Text>
			<SearchInput query={query} active={searchActive} width={width} placeholder="按 / 搜索模块、题目名" />
			<ScrollList
				items={categories}
				cursor={cursor}
				height={listHeight}
				width={width}
				emptyTitle="没有匹配的模块"
				renderItem={(category, active) => {
					const marker = active ? '›' : ' ';
					return `${marker} ${category.label.padEnd(14)} ${compactProgress(category.passed, category.total, layout.mode === 'compact' ? 6 : 10)}`;
				}}
			/>
		</Box>
	);
}
