import React from 'react';
import {renderToString} from 'ink';
import {describe, expect, test} from 'vitest';
import {App} from './App.js';
import type {LinuxLabData} from './types.js';

const sampleData: LinuxLabData = {
	categories: [
		{
			id: 'linux-basics',
			label: 'Linux 基础命令',
			total: 2,
			passed: 1,
			challenges: [
				{
					id: 'ls-basic',
					title: 'ls 基础',
					category: 'linux-basics',
					subcategory: 'navigation',
					difficulty: 1,
					tags: ['ls'],
					description: '学习使用 ls 命令列出目录内容。',
					hints: [{level: 1, text: '试试 ls -la'}],
					verify: [],
				},
				{
					id: 'grep-basic',
					title: 'grep 搜索',
					category: 'linux-basics',
					subcategory: 'pipes',
					difficulty: 2,
					tags: ['grep'],
					description: '用 grep 过滤文本。',
					hints: [],
					verify: [],
				},
			],
		},
	],
	progress: {
		challenges: {
			'ls-basic': {status: 'passed', attempts: 1, hints_used: 0, last_attempt: '2026-05-09'},
		},
		skills: {},
	},
	references: {
		commands: [
			{
				name: 'ls',
				brief: '列出目录内容',
				examples: [{desc: '列出所有文件', cmd: 'ls -la'}],
				related_challenges: ['ls-basic'],
			},
		],
	},
};

describe('React TUI App', () => {
	test('renders the training console menu', () => {
		const output = renderToString(<App data={sampleData} initialScreen="menu" />, {columns: 100});

		expect(output).toContain('LinuxLab');
		expect(output).toContain('训练控制台');
		expect(output).toContain('开始练习');
		expect(output).toContain('命令速查');
	});

	test('renders module progress summary', () => {
		const output = renderToString(<App data={sampleData} initialScreen="modules" />, {columns: 100});

		expect(output).toContain('整体进度');
		expect(output).toContain('Linux 基础命令');
		expect(output).toContain('1/2');
	});

	test('renders reference search state', () => {
		const output = renderToString(<App data={sampleData} initialScreen="reference" initialQuery="ls" />, {columns: 100});

		expect(output).toContain('命令速查');
		expect(output).toContain('ls');
		expect(output).toContain('列出目录内容');
	});
});
