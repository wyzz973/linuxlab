import React from 'react';
import {Box, Text} from 'ink';
import {theme} from '../../theme/theme.js';
import {fitText} from '../../utils/text.js';
import {EmptyState} from './EmptyState.js';

type ScrollListProps<T> = {
	items: T[];
	cursor: number;
	height: number;
	width: number;
	emptyTitle?: string;
	emptyAction?: string;
	renderItem: (item: T, active: boolean, index: number) => string;
};

function visibleStart(cursor: number, count: number, height: number): number {
	if (count <= height) {
		return 0;
	}
	if (cursor >= count - height) {
		return count - height;
	}
	return Math.max(0, cursor);
}

export function ScrollList<T>({
	items,
	cursor,
	height,
	width,
	emptyTitle = '没有可显示内容',
	emptyAction = '换关键词或按 Esc 清空',
	renderItem,
}: ScrollListProps<T>) {
	const safeHeight = Math.max(1, Math.floor(height));
	const safeWidth = Math.max(1, Math.floor(width));
	const safeCursor = Math.max(0, Math.min(items.length - 1, cursor));

	if (items.length === 0) {
		return <EmptyState title={emptyTitle} action={emptyAction} width={safeWidth} />;
	}

	const start = visibleStart(safeCursor, items.length, safeHeight);
	const visible = items.slice(start, start + safeHeight);

	return (
		<Box flexDirection="column">
			{visible.map((item, offset) => {
				const index = start + offset;
				const active = index === safeCursor;
				return (
					<Text key={index} color={active ? theme.accent : undefined}>
						{fitText(renderItem(item, active, index), safeWidth)}
					</Text>
				);
			})}
			{items.length > safeHeight && (
				<Text color={theme.dim}>{fitText(`${start + 1}-${Math.min(start + safeHeight, items.length)}/${items.length}`, safeWidth)}</Text>
			)}
		</Box>
	);
}
