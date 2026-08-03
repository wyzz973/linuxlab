import {mkdtemp, mkdir, writeFile} from 'node:fs/promises';
import {tmpdir} from 'node:os';
import path from 'node:path';
import {describe, expect, test} from 'vitest';
import {loadLinuxLabData} from './data.js';

describe('loadLinuxLabData', () => {
	test('loads challenges, progress and command references from disk', async () => {
		const root = await mkdtemp(path.join(tmpdir(), 'linuxlab-react-tui-'));
		const challengeDir = path.join(root, 'challenges', 'linux-basics', 'ls-basic');
		const referencesDir = path.join(root, 'references');
		await mkdir(challengeDir, {recursive: true});
		await mkdir(referencesDir, {recursive: true});

		await writeFile(
			path.join(challengeDir, 'challenge.yaml'),
			[
				'id: ls-basic',
				'title: "ls 基础"',
				'difficulty: 1',
				'category: linux-basics',
				'subcategory: navigation',
				'tags: [ls]',
				'description: "列出目录内容"',
				'hints:',
				'  - level: 1',
				'    text: "试试 ls -la"',
				'verify: []',
				'',
			].join('\n'),
		);
		await writeFile(
			path.join(root, 'progress.json'),
			JSON.stringify({
				skills: {},
				challenges: {
					'ls-basic': {status: 'passed', attempts: 1, hints_used: 0, last_attempt: '2026-05-09'},
				},
			}),
		);
		await writeFile(
			path.join(referencesDir, 'commands.yaml'),
			[
				'commands:',
				'  - name: ls',
				'    brief: 列出目录内容',
				'    examples:',
				'      - desc: 列出所有文件',
				'        cmd: ls -la',
				'',
			].join('\n'),
		);

		const data = await loadLinuxLabData({
			rootDir: root,
			progressPath: path.join(root, 'progress.json'),
		});

		expect(data.categories).toHaveLength(1);
		expect(data.categories[0]?.label).toBe('Linux 基础命令');
		expect(data.categories[0]?.passed).toBe(1);
		expect(data.categories[0]?.challenges[0]?.title).toBe('ls 基础');
		expect(data.references.commands[0]?.name).toBe('ls');
	});
});
