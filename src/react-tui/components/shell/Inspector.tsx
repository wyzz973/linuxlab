import React from 'react';
import {Box, Text} from 'ink';
import {theme} from '../../theme/theme.js';
import {fitText} from '../../utils/text.js';

type InspectorProps = {
	width: number;
	children?: React.ReactNode;
};

export function Inspector({width, children}: InspectorProps) {
	return (
		<Box flexDirection="column" width={width} borderStyle="single" borderColor={theme.border} paddingX={1}>
			<Text bold color={theme.accent}>上下文</Text>
			<Box marginTop={1} flexDirection="column">
				{children ?? <Text color={theme.dim}>{fitText('选择题目或命令后显示更多信息', Math.max(1, width - 4))}</Text>}
			</Box>
		</Box>
	);
}
