import React from 'react';
import {Box, Text} from 'ink';
import type {Challenge, LayoutSpec, ProgressData} from '../../types.js';
import {theme} from '../../theme/theme.js';
import {ScrollList} from '../common/ScrollList.js';
import {statusIcon, statusLabel} from '../common/StatusBadge.js';
import {stars} from './helpers.js';

type RecommendScreenProps = {
	challenges: Challenge[];
	cursor: number;
	progress: ProgressData['challenges'];
	layout: LayoutSpec;
};

export function RecommendScreen({challenges, cursor, progress, layout}: RecommendScreenProps) {
	const width = layout.mainWidth;
	return (
		<Box flexDirection="column" width={width}>
			<Text bold color={theme.accent}>薄弱推荐</Text>
			<Text color={theme.dim}>优先显示失败过和未完成题目。</Text>
			<ScrollList
				items={challenges}
				cursor={cursor}
				height={Math.max(3, layout.mainHeight - 3)}
				width={width}
				emptyTitle="暂无推荐"
				emptyAction="完成或失败几道题后会生成推荐"
				renderItem={(challenge, active) => {
					const status = progress[challenge.id]?.status;
					return `${active ? '›' : ' '} ${statusIcon(status)} ${challenge.title.padEnd(24)} ${stars(challenge.difficulty)} ${statusLabel(status)}`;
				}}
			/>
		</Box>
	);
}
