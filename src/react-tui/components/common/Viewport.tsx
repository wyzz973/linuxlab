import React from 'react';
import {Box, Text} from 'ink';
import {theme} from '../../theme/theme.js';
import {fitText, wrapText} from '../../utils/text.js';

type ViewportProps = {
	text: string | string[];
	width: number;
	height: number;
};

export function Viewport({text, width, height}: ViewportProps) {
	const safeWidth = Math.max(1, Math.floor(width));
	const safeHeight = Math.max(1, Math.floor(height));
	const lines = Array.isArray(text) ? text.flatMap(line => wrapText(line, safeWidth)) : wrapText(text, safeWidth);
	const visible = lines.slice(0, safeHeight);
	const overflow = lines.length > safeHeight;

	if (overflow && visible.length > 0) {
		visible[visible.length - 1] = `${fitText(visible[visible.length - 1] ?? '', Math.max(1, safeWidth - 1))}…`;
	}

	return (
		<Box flexDirection="column">
			{visible.map((line, index) => (
				<Text key={index} color={overflow && index === visible.length - 1 ? theme.dim : undefined}>
					{fitText(line, safeWidth)}
				</Text>
			))}
		</Box>
	);
}
