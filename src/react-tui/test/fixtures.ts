import type {LinuxLabData} from '../types.js';

export const fixtureData: LinuxLabData = {
	categories: [
		{
			id: 'linux-basics',
			label: 'Linux 基础命令',
			total: 3,
			passed: 1,
			challenges: [
				{
					id: 'ls-basic',
					title: 'ls 基础',
					difficulty: 1,
					category: 'linux-basics',
					subcategory: 'files',
					tags: ['ls'],
					description: '使用 ls 查看当前目录内容，并理解隐藏文件和详细列表。这个描述故意稍长，用来检查终端换行和裁剪是否稳定。',
					hints: [{level: 1, text: '先运行 ls'}],
					verify: [{type: 'script', path: 'check.sh'}],
				},
				{
					id: 'grep-basic',
					title: 'grep 搜索长标题用于测试裁剪',
					difficulty: 2,
					category: 'linux-basics',
					subcategory: 'text',
					tags: ['grep'],
					description: '使用 grep 从文本中搜索指定模式。',
					hints: [],
					verify: [],
				},
				{
					id: 'chmod-basic',
					title: 'chmod 权限',
					difficulty: 3,
					category: 'linux-basics',
					subcategory: 'permissions',
					tags: ['chmod'],
					description: '修改文件权限。',
					hints: [],
					verify: [],
				},
			],
		},
	],
	progress: {
		skills: {},
		challenges: {
			'ls-basic': {status: 'passed', attempts: 1, hints_used: 0, last_attempt: '2026-05-09T00:00:00Z'},
		},
	},
	references: {
		commands: [
			{name: 'ls', brief: '列出目录内容', examples: [{desc: '查看详细信息', cmd: 'ls -la'}]},
			{name: 'grep', brief: '搜索文本', examples: [{desc: '搜索 error', cmd: 'grep error app.log'}]},
		],
	},
};
