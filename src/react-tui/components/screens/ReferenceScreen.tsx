import React from 'react';
import {Box, Text} from 'ink';
import type {CommandRef, LayoutSpec} from '../../types.js';
import {theme} from '../../theme/theme.js';
import {fitText} from '../../utils/text.js';
import {SearchInput} from '../common/SearchInput.js';
import {ScrollList} from '../common/ScrollList.js';
import {Viewport} from '../common/Viewport.js';

type ReferenceScreenProps = {
	commands: CommandRef[];
	selected?: CommandRef;
	cursor: number;
	query: string;
	searchActive?: boolean;
	layout: LayoutSpec;
};

export function ReferenceScreen({commands, selected, cursor, query, searchActive = false, layout}: ReferenceScreenProps) {
	const width = layout.mainWidth;
	const previewHeight = layout.mode === 'compact' ? 3 : 5;
	const listHeight = Math.max(3, layout.mainHeight - previewHeight - 4);

	return (
		<Box flexDirection="column" width={width}>
			<Text bold color={theme.accent}>命令速查</Text>
			<SearchInput query={query} active={searchActive} width={width} placeholder="按 / 搜索命令名、说明或示例" />
			<ScrollList
				items={commands}
				cursor={cursor}
				height={listHeight}
				width={width}
				emptyTitle="没有匹配的命令"
				renderItem={(command, active) => `${active ? '›' : ' '} ${command.name.padEnd(12)} ${command.brief}  ${command.examples.length} 例`}
			/>
			{selected && layout.mode !== 'wide' && (
				<Box marginTop={1} flexDirection="column">
					<Text color={theme.warning}>{fitText(`${selected.name}: ${selected.brief}`, width)}</Text>
					<Viewport
						text={selected.examples.map(example => `${example.desc}: ${example.cmd}`)}
						width={width}
						height={previewHeight}
					/>
				</Box>
			)}
		</Box>
	);
}
