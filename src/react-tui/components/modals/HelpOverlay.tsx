import React from 'react';
import {Box, Text} from 'ink';
import {hintsForScreen} from '../../app/keymap.js';
import {theme} from '../../theme/theme.js';
import type {LayoutSpec, ScreenID} from '../../types.js';
import {Modal} from '../common/Modal.js';

type HelpOverlayProps = {
	layout: LayoutSpec;
	screen: ScreenID;
};

export function HelpOverlay({layout, screen}: HelpOverlayProps) {
	return (
		<Modal layout={layout} title="帮助">
			<Box flexDirection="column" marginTop={1}>
				{hintsForScreen(screen, false).map(hint => (
					<Text key={`${hint.keys}-${hint.label}`}>
						<Text color={theme.accent}>{hint.keys.padEnd(8)}</Text>
						<Text>{hint.label}</Text>
					</Text>
				))}
				<Text color={theme.dim}>Esc 关闭帮助</Text>
			</Box>
		</Modal>
	);
}
