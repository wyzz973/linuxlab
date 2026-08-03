import React from 'react';
import {Box, Text} from 'ink';
import type {Category, LayoutSpec} from '../../types.js';
import {theme} from '../../theme/theme.js';
import {ScrollList} from '../common/ScrollList.js';
import {compactProgress} from './helpers.js';

type SkillMapScreenProps = {
	categories: Category[];
	cursor: number;
	layout: LayoutSpec;
};

export function SkillMapScreen({categories, cursor, layout}: SkillMapScreenProps) {
	const width = layout.mainWidth;
	return (
		<Box flexDirection="column" width={width}>
			<Text bold color={theme.accent}>能力图谱</Text>
			<Text color={theme.dim}>按模块查看掌握度，后续可细分到 subcategory。</Text>
			<ScrollList
				items={categories}
				cursor={cursor}
				height={Math.max(3, layout.mainHeight - 3)}
				width={width}
				renderItem={(category, active) => `${active ? '›' : ' '} ${category.label.padEnd(16)} ${compactProgress(category.passed, category.total, 12)}`}
			/>
		</Box>
	);
}
