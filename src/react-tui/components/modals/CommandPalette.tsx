import React from 'react';
import {Box, Text} from 'ink';
import {filterCommands} from '../../app/commands.js';
import {theme} from '../../theme/theme.js';
import type {LayoutSpec} from '../../types.js';
import {Modal} from '../common/Modal.js';
import {SearchInput} from '../common/SearchInput.js';
import {ScrollList} from '../common/ScrollList.js';

type CommandPaletteProps = {
	layout: LayoutSpec;
	query: string;
	cursor: number;
};

export function CommandPalette({layout, query, cursor}: CommandPaletteProps) {
	const commands = filterCommands(query);
	const width = layout.mode === 'compact' ? layout.contentWidth - 4 : Math.min(68, layout.contentWidth - 12);
	return (
		<Modal layout={layout} title="命令面板">
			<Box flexDirection="column" marginTop={1}>
				<SearchInput query={query} width={width} placeholder="输入命令或页面名" />
				<ScrollList
					items={commands}
					cursor={cursor}
					height={Math.min(6, Math.max(3, layout.mainHeight - 4))}
					width={width}
					emptyTitle="没有匹配的命令"
					emptyAction="换关键词或按 Esc 关闭"
					renderItem={(command, active) => `${active ? '›' : ' '} ${command.label.padEnd(10)} ${command.description} ${command.keys ? `(${command.keys})` : ''}`}
				/>
				<Text color={theme.dim}>Enter 执行 · Esc 关闭</Text>
			</Box>
		</Modal>
	);
}
