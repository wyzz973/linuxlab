import React from 'react';
import {Box, Text} from 'ink';
import type {LayoutSpec, LinuxLabData} from '../../types.js';
import {countFailedChallenges, summarizeCategories} from '../../domain/selectors.js';
import {theme} from '../../theme/theme.js';
import {fitText} from '../../utils/text.js';
import {compactProgress} from './helpers.js';

type HomeScreenProps = {
	data: LinuxLabData;
	layout: LayoutSpec;
};

export function HomeScreen({data, layout}: HomeScreenProps) {
	const totals = summarizeCategories(data.categories);
	const failed = countFailedChallenges(data.categories, data.progress.challenges);
	const action = totals.passed > 0 ? '继续练习' : '开始第一题';
	const width = layout.mainWidth;

	return (
		<Box flexDirection="column" width={width}>
			<Text bold color={theme.accent}>训练控制台</Text>
			<Text color={theme.dim}>{fitText('真实终端练习、命令速查和能力复盘集中在同一个工作区。', width)}</Text>
			<Box marginTop={1} flexDirection={layout.mode === 'compact' ? 'column' : 'row'}>
				<Metric label="题库" value={`${totals.total} 题`} width={Math.floor(width / 2)} tone={theme.accent} />
				<Metric label="通过" value={`${totals.passed} 题`} width={Math.floor(width / 2)} tone={theme.success} />
			</Box>
			<Box flexDirection={layout.mode === 'compact' ? 'column' : 'row'}>
				<Metric label="待复盘" value={`${failed} 题`} width={Math.floor(width / 2)} tone={failed > 0 ? theme.warning : theme.dim} />
				<Metric label="速查" value={`${data.references.commands.length} 条`} width={Math.floor(width / 2)} tone={theme.accent} />
			</Box>
			<Box marginTop={1} flexDirection="column">
				<Text color={theme.warning}>{action}</Text>
				<Text>{fitText(`整体进度 ${compactProgress(totals.passed, totals.total, layout.mode === 'compact' ? 10 : 18)}`, width)}</Text>
			</Box>
			{layout.mode === 'wide' && (
				<Box marginTop={1} flexDirection="column">
					<Text bold color={theme.accent}>最近可做</Text>
					{data.categories.flatMap(category => category.challenges).slice(0, 4).map(challenge => (
						<Text key={challenge.id}>{fitText(`○ ${challenge.title}`, width)}</Text>
					))}
				</Box>
			)}
		</Box>
	);
}

function Metric({label, value, width, tone}: {label: string; value: string; width: number; tone: string}) {
	return (
		<Box width={Math.max(16, width)}>
			<Text color={theme.dim}>{label} </Text>
			<Text bold color={tone}>{value}</Text>
		</Box>
	);
}
