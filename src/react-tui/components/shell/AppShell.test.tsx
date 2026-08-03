import React from 'react';
import {renderToString, Text} from 'ink';
import {describe, expect, it} from 'vitest';
import {createLayoutSpec} from '../../app/breakpoints.js';
import {AppShell} from './AppShell.js';

describe('AppShell', () => {
	it('renders compact tabs without sidebar', () => {
		const output = renderToString(
			<AppShell
				layout={createLayoutSpec({columns: 60, rows: 20})}
				screen="menu"
				title="LinuxLab"
				progressText="0/278"
				notice="Docker ready"
			>
				<Text>测试内容</Text>
			</AppShell>,
		);
		expect(output).toContain('总览');
		expect(output).not.toContain('导航');
	});

	it('renders wide inspector slot', () => {
		const output = renderToString(
			<AppShell
				layout={createLayoutSpec({columns: 140, rows: 40})}
				screen="detail"
				title="LinuxLab"
				progressText="1/278"
				notice="Docker ready"
				inspector={<Text>验证规则</Text>}
			>
				<Text>题目详情</Text>
			</AppShell>,
		);
		expect(output).toContain('导航');
		expect(output).toContain('验证规则');
	});
});
