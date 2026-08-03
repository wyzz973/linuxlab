import type {ScreenID} from '../types.js';

export type TuiCommand = {
	id: string;
	label: string;
	description: string;
	keys?: string;
	target?: ScreenID;
};

export const tuiCommands: TuiCommand[] = [
	{id: 'practice', label: '开始练习', description: '进入模块列表', keys: '2', target: 'modules'},
	{id: 'reference', label: '命令速查', description: '搜索 Linux 常用命令', keys: '4', target: 'reference'},
	{id: 'recommend', label: '查看推荐', description: '根据历史进度查看薄弱题', keys: '3', target: 'recommend'},
	{id: 'skillmap', label: '能力图谱', description: '查看分类掌握度', keys: '5', target: 'skillmap'},
	{id: 'help', label: '打开帮助', description: '查看当前快捷键', keys: '?'},
];

export function filterCommands(query: string): TuiCommand[] {
	const normalized = query.trim().toLowerCase();
	if (normalized === '') {
		return tuiCommands;
	}
	return tuiCommands.filter(command =>
		[command.id, command.label, command.description, command.keys ?? ''].join(' ').toLowerCase().includes(normalized),
	);
}
