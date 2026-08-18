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
	hintLevel?: number;
};

export function ChallengeDetailScreen({challenge, relatedCommands, progress, layout, running, hintLevel = 0}: ChallengeDetailScreenProps) {
	const width = layout.mainWidth;
	const hintRows = hintLevel > 0 ? Math.min(hintLevel, 3) * 2 + 2 : 0;
	const descriptionHeight = Math.max(4, layout.mainHeight - (layout.mode === 'compact' ? 7 : 9) - hintRows);
	const status = progress[challenge.id]?.status;
	const unlocked = challenge.hints.slice(0, Math.min(hintLevel, challenge.hints.length));

	return (
		<Box flexDirection="column" width={width}>
			<Text bold color={theme.accent}>{challenge.title}</Text>
			<Text color={theme.dim}>难度 {stars(challenge.difficulty)} · {challenge.tags.join(', ') || '无标签'}</Text>
			<StatusBadge status={status} />
			<Box marginTop={1} flexDirection="column">
				<Text bold>任务描述</Text>
				<Viewport text={challenge.description.trim() || '暂无描述'} width={width} height={descriptionHeight} />
			</Box>
			<Box marginTop={1} flexDirection="column">
				<Text bold>提示</Text>
				{challenge.hints.length === 0 ? (
					<Text color={theme.dim}>本题没有提示</Text>
				) : hintLevel === 0 ? (
					<Text color={theme.dim}>按 h 查看第一条提示</Text>
				) : (
					unlocked.map((hint, index) => (
						<Viewport key={hint.level} text={`${index + 1}. ${hint.text}`} width={width} height={2} />
					))
				)}
				{challenge.hints.length > 0 && (
					<Text color={theme.dim}>
						{hintLevel < challenge.hints.length
							? `h 解锁下一条提示（影响得分） · 已用 ${hintLevel}/${challenge.hints.length}`
							: `已解锁全部提示（${hintLevel}/${challenge.hints.length}）`}
					</Text>
				)}
			</Box>
			{layout.mode !== 'wide' && (
				<Box marginTop={1} flexDirection="column">
					<Text color={theme.dim}>检查 {challenge.verify.length} 项 · 相关命令 {relatedCommands.map(command => command.name).join(', ') || '无'}</Text>
				</Box>
			)}
			<Box marginTop={1}>
				<Text color={running ? theme.warning : theme.accent}>{running ? '正在启动...' : 'Enter 开始挑战'}</Text>
			</Box>
		</Box>
	);
}
