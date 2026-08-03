import React from 'react';
import {Box, Text} from 'ink';
import type {LayoutSpec, ScreenID} from '../../types.js';
import {theme} from '../../theme/theme.js';
import {fitText} from '../../utils/text.js';
import {Footer} from './Footer.js';
import {Header} from './Header.js';
import {Inspector} from './Inspector.js';
import {Navigation} from './Navigation.js';

type AppShellProps = {
	layout: LayoutSpec;
	screen: ScreenID;
	title: string;
	progressText: string;
	notice: string;
	statusText?: string;
	inspector?: React.ReactNode;
	modal?: React.ReactNode;
	children: React.ReactNode;
};

export function AppShell({
	layout,
	screen,
	title,
	progressText,
	notice,
	statusText,
	inspector,
	modal,
	children,
}: AppShellProps) {
	if (layout.mode === 'unsupported') {
		return (
			<Box flexDirection="column" width={layout.contentWidth}>
				<Text bold color={theme.error}>LinuxLab 需要更大的终端窗口</Text>
				<Text color={theme.dim}>{fitText('最低支持 60 列 x 18 行；请放大窗口后继续。', layout.contentWidth)}</Text>
			</Box>
		);
	}

	return (
		<Box flexDirection="column" width={layout.contentWidth}>
			<Header title={title} screen={screen} progressText={progressText} statusText={statusText} width={layout.contentWidth} />
			<Text color={theme.border}>{'─'.repeat(layout.contentWidth)}</Text>
			{layout.showCompactTabs && <Navigation screen={screen} width={layout.contentWidth} compact />}
			<Box flexDirection="row" height={layout.mainHeight}>
				{layout.showSidebar && <Navigation screen={screen} width={layout.navWidth} />}
				<Box flexDirection="column" width={layout.mainWidth} marginLeft={layout.showSidebar ? 1 : 0}>
					{children}
				</Box>
				{layout.showInspector && <Box width={2}><Text>  </Text></Box>}
				{layout.showInspector && <Inspector width={layout.inspectorWidth}>{inspector}</Inspector>}
			</Box>
			{modal}
			<Text color={theme.border}>{'─'.repeat(layout.contentWidth)}</Text>
			<Footer screen={screen} notice={notice} width={layout.contentWidth} compact={layout.mode === 'compact'} />
		</Box>
	);
}
