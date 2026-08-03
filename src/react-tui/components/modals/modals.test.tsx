import React from 'react';
import {renderToString, Text} from 'ink';
import {describe, expect, it} from 'vitest';
import {createLayoutSpec} from '../../app/breakpoints.js';
import {filterCommands} from '../../app/commands.js';
import {Modal} from '../common/Modal.js';
import {CommandPalette} from './CommandPalette.js';
import {HelpOverlay} from './HelpOverlay.js';

describe('modal layer', () => {
	it('renders compact modal', () => {
		const output = renderToString(
			<Modal layout={createLayoutSpec({columns: 60, rows: 20})} title="帮助">
				<Text>内容</Text>
			</Modal>,
		);
		expect(output).toContain('帮助');
		expect(output).toContain('内容');
	});

	it('renders help overlay with current hints', () => {
		const output = renderToString(<HelpOverlay layout={createLayoutSpec({columns: 100, rows: 30})} screen="detail" />);
		expect(output).toContain('开始挑战');
	});

	it('filters command palette by Chinese label and id', () => {
		expect(filterCommands('速查')[0]?.id).toBe('reference');
		expect(filterCommands('practice')[0]?.label).toBe('开始练习');
	});

	it('renders command palette', () => {
		const output = renderToString(<CommandPalette layout={createLayoutSpec({columns: 100, rows: 30})} query="练习" cursor={0} />);
		expect(output).toContain('命令面板');
		expect(output).toContain('开始练习');
	});
});
