import React from 'react';
import {Text} from 'ink';
import {theme} from '../../theme/theme.js';

type ProgressBarProps = {
	value: number;
	total: number;
	width: number;
};

export function progressString(value: number, total: number, width: number): string {
	const safeWidth = Math.max(1, Math.floor(width));
	const ratio = total <= 0 ? 0 : Math.max(0, Math.min(1, value / total));
	const filled = Math.round(safeWidth * ratio);
	return `${'█'.repeat(filled)}${'░'.repeat(safeWidth - filled)}`;
}

export function ProgressBar({value, total, width}: ProgressBarProps) {
	const bar = progressString(value, total, width);
	const filled = bar.indexOf('░') === -1 ? bar.length : bar.indexOf('░');
	return (
		<Text>
			<Text color={theme.success}>{bar.slice(0, filled)}</Text>
			<Text color={theme.border}>{bar.slice(filled)}</Text>
		</Text>
	);
}
