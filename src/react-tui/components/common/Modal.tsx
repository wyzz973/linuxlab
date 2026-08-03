import React from 'react';
import {Box, Text} from 'ink';
import type {LayoutSpec} from '../../types.js';
import {theme} from '../../theme/theme.js';

type ModalProps = {
	layout: LayoutSpec;
	title: string;
	children: React.ReactNode;
};

export function Modal({layout, title, children}: ModalProps) {
	const width = layout.mode === 'compact' ? layout.contentWidth : Math.min(72, Math.max(20, layout.contentWidth - 8));
	return (
		<Box width={width} borderStyle="single" borderColor={theme.accent} paddingX={1} flexDirection="column">
			<Text bold color={theme.accent}>{title}</Text>
			{children}
		</Box>
	);
}
