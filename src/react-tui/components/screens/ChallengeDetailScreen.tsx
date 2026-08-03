import React from 'react';
import {Box, Text} from 'ink';
import type {Challenge, CommandRef, LayoutSpec, ProgressData} from '../../types.js';
import {theme} from '../../theme/theme.js';
import {Viewport} from '../common/Viewport.js';
import {StatusBadge} from '../common/StatusBadge.js';
import {stars} from './helpers.js';

type ChallengeDetailScreenProps = {
	challenge: Challenge;
	relatedCommands: CommandRef[];
	progress: ProgressData['challenges'];
	layout: LayoutSpec;
	running?: boolean;
};

export function ChallengeDetailScreen({challenge, relatedCommands, progress, layout, running}: ChallengeDetailScreenProps) {
	const width = layout.mainWidth;
	const descriptionHeight = Math.max(4, layout.mainHeight - (layout.mode === 'compact' ? 7 : 9));
	const status = progress[challenge.id]?.status;

	return (
		<Box flexDirection="column" width={width}>
			<Text bold color={theme.accent}>{challenge.title}</Text>
			<Text color={theme.dim}>难度 {stars(challenge.difficulty)} · {challenge.tags.join(', ') || '无标签'}</Text>
			<StatusBadge status={status} />
			<Box marginTop={1} flexDirection="column">
				<Text bold>任务描述</Text>
				<Viewport text={challenge.description.trim() || '暂无描述'} width={width} height={descriptionHeight} />
			</Box>
			{layout.mode !== 'wide' && (
				<Box marginTop={1} flexDirection="column">
					<Text color={theme.dim}>提示 {challenge.hints.length} 条 · 检查 {challenge.verify.length} 项</Text>
					<Text color={theme.dim}>相关命令 {relatedCommands.map(command => command.name).join(', ') || '无'}</Text>
				</Box>
			)}
			<Box marginTop={1}>
				<Text color={running ? theme.warning : theme.accent}>{running ? '正在启动...' : 'Enter 开始挑战'}</Text>
			</Box>
		</Box>
	);
}
