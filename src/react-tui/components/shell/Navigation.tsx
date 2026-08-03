import React from 'react';
import {Box, Text} from 'ink';
import type {ScreenID} from '../../types.js';
import {theme} from '../../theme/theme.js';
import {fitText} from '../../utils/text.js';

export type NavItem = {
	id: ScreenID;
	label: string;
	hint: string;
};

export const navItems: NavItem[] = [
	{id: 'menu', label: '总览', hint: '训练控制台'},
	{id: 'modules', label: '开始练习', hint: '按主题练习'},
	{id: 'recommend', label: '推荐', hint: '薄弱题'},
	{id: 'reference', label: '命令速查', hint: '命令参考'},
	{id: 'skillmap', label: '能力图谱', hint: '掌握度'},
];

type NavigationProps = {
	screen: ScreenID;
	width: number;
	compact?: boolean;
};

export function Navigation({screen, width, compact = false}: NavigationProps) {
	if (compact) {
		const segments = navItems.map((item, index) => `${item.id === screen ? '›' : ' '} ${index + 1}.${item.label}`);
		return (
			<Box width={width}>
				<Text color={theme.dim}>{fitText(segments.join('  '), width)}</Text>
			</Box>
		);
	}

	return (
		<Box flexDirection="column" width={width} borderStyle="single" borderColor={theme.border} paddingX={1}>
			<Text bold color={theme.accent}>导航</Text>
			{navItems.map((item, index) => {
				const active = item.id === screen;
				return (
					<Box key={item.id} flexDirection="column" marginTop={index === 0 ? 1 : 0}>
						<Text color={active ? theme.accent : undefined}>{active ? '› ' : '  '}{index + 1}. {fitText(item.label, Math.max(1, width - 7))}</Text>
						<Text color={theme.dim}>{fitText(`   ${item.hint}`, Math.max(1, width - 4))}</Text>
					</Box>
				);
			})}
		</Box>
	);
}
