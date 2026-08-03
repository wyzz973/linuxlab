import React from 'react';
import {Box, Text} from 'ink';
import type {Category, LayoutSpec, ProgressData} from '../../types.js';
import {theme} from '../../theme/theme.js';
import {SearchInput} from '../common/SearchInput.js';
import {ScrollList} from '../common/ScrollList.js';
import {statusIcon, statusLabel} from '../common/StatusBadge.js';
import {compactProgress, stars} from './helpers.js';

type ChallengeListScreenProps = {
	category: Category;
	cursor: number;
	query: string;
	searchActive?: boolean;
	progress: ProgressData['challenges'];
	layout: LayoutSpec;
};

export function ChallengeListScreen({category, cursor, query, searchActive = false, progress, layout}: ChallengeListScreenProps) {
	const width = layout.mainWidth;
	const listHeight = Math.max(3, layout.mainHeight - 4);

	return (
		<Box flexDirection="column" width={width}>
			<Text bold color={theme.accent}>{category.label}</Text>
			<Text color={theme.dim}>完成进度 {compactProgress(category.passed, category.total, layout.mode === 'compact' ? 8 : 14)}</Text>
			<SearchInput query={query} active={searchActive} width={width} placeholder="按 / 搜索题目、标签、说明" />
			<ScrollList
				items={category.challenges}
				cursor={cursor}
				height={listHeight}
				width={width}
				emptyTitle="没有匹配的题目"
				renderItem={(challenge, active) => {
					const status = progress[challenge.id]?.status;
					const marker = active ? '›' : ' ';
					if (layout.mode === 'compact') {
						return `${marker} ${statusIcon(status)} ${challenge.title} ${stars(challenge.difficulty).slice(0, challenge.difficulty)}`;
					}
					return `${marker} ${statusIcon(status)} ${challenge.title.padEnd(24)} ${stars(challenge.difficulty)} ${statusLabel(status)}`;
				}}
			/>
		</Box>
	);
}
