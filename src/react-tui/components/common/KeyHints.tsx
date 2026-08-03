import React from 'react';
import {Text} from 'ink';
import type {KeyHint} from '../../app/keymap.js';
import {theme} from '../../theme/theme.js';
import {fitText} from '../../utils/text.js';

type KeyHintsProps = {
	hints: KeyHint[];
	width: number;
};

export function KeyHints({hints, width}: KeyHintsProps) {
	const content = hints.map(hint => `${hint.keys} ${hint.label}`).join(' · ');
	return <Text color={theme.dim}>{fitText(content, width)}</Text>;
}
