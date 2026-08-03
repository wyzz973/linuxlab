import React from 'react';
import {Text} from 'ink';
import {theme} from '../../theme/theme.js';
import {displayWidth, fitText} from '../../utils/text.js';

type SearchInputProps = {
	query: string;
	width: number;
	placeholder?: string;
	label?: string;
	active?: boolean;
};

export function SearchInput({query, width, placeholder = '按 / 搜索', label = '搜索', active = false}: SearchInputProps) {
	const value = query === '' ? placeholder : query;
	const color = query === '' ? theme.dim : theme.accent;
	const displayLabel = active ? `${label}*` : label;
	return (
		<Text>
			<Text color={active ? theme.warning : theme.dim}>{displayLabel}: </Text>
			<Text color={color}>{fitText(value, Math.max(1, width - displayWidth(displayLabel) - 2))}</Text>
		</Text>
	);
}
