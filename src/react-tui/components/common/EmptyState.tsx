import React from 'react';
import {Box, Text} from 'ink';
import {theme} from '../../theme/theme.js';
import {fitText} from '../../utils/text.js';

type EmptyStateProps = {
	title: string;
	action: string;
	width: number;
};

export function EmptyState({title, action, width}: EmptyStateProps) {
	return (
		<Box flexDirection="column">
			<Text color={theme.dim}>{fitText(title, width)}</Text>
			<Text color={theme.warning}>{fitText(action, width)}</Text>
		</Box>
	);
}
