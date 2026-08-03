import React from 'react';
import {Box, Text} from 'ink';
import type {ScreenID} from '../../types.js';
import {hintsForScreen} from '../../app/keymap.js';
import {theme} from '../../theme/theme.js';
import {fitText} from '../../utils/text.js';
import {KeyHints} from '../common/KeyHints.js';

type FooterProps = {
	screen: ScreenID;
	notice: string;
	width: number;
	compact: boolean;
};

export function Footer({screen, notice, width, compact}: FooterProps) {
	const hintWidth = compact ? width : Math.max(10, Math.floor(width * 0.55));
	const noticeWidth = compact ? 0 : Math.max(8, width - hintWidth - 2);
	return (
		<Box width={width} justifyContent={compact ? 'flex-start' : 'space-between'}>
			{!compact && <Text color={theme.warning}>{fitText(notice, noticeWidth)}</Text>}
			<KeyHints hints={hintsForScreen(screen, compact)} width={compact ? width : hintWidth} />
		</Box>
	);
}
