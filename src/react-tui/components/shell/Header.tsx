import React from 'react';
import {Box, Text} from 'ink';
import type {ScreenID} from '../../types.js';
import {theme} from '../../theme/theme.js';
import {fitText} from '../../utils/text.js';

const routeLabels: Record<ScreenID, string> = {
	menu: '总览',
	modules: '模块',
	challenges: '练习列表',
	detail: '题目详情',
	reference: '命令速查',
	recommend: '薄弱推荐',
	skillmap: '能力图谱',
	result: '检测结果',
};

type HeaderProps = {
	title: string;
	screen: ScreenID;
	progressText: string;
	statusText?: string;
	width: number;
};

export function Header({title, screen, progressText, statusText = 'React preview', width}: HeaderProps) {
	const right = `${statusText} · ${progressText}`;
	const leftWidth = Math.max(10, width - right.length - 3);
	return (
		<Box width={width} justifyContent="space-between">
			<Text>
				<Text bold color={theme.accent}>{title}</Text>
				<Text color={theme.dim}> / {fitText(routeLabels[screen], Math.max(1, leftWidth - title.length - 3))}</Text>
			</Text>
			<Text color={theme.dim}>{fitText(right, Math.max(1, width - leftWidth - 1))}</Text>
		</Box>
	);
}
