import type {ScreenID} from '../types.js';

export type KeyHint = {
	keys: string;
	label: string;
};

export const globalHints: KeyHint[] = [
	{keys: '↑↓', label: '选择'},
	{keys: 'Enter', label: '确认'},
	{keys: '/', label: '搜索'},
	{keys: '?', label: '帮助'},
	{keys: 'q', label: '返回'},
];

export function hintsForScreen(screen: ScreenID, compact: boolean): KeyHint[] {
	const screenSpecific: Record<ScreenID, KeyHint[]> = {
		menu: [{keys: '1-5', label: '跳转'}],
		modules: [{keys: 'Enter', label: '进入模块'}],
		challenges: [{keys: 'Enter', label: '查看题目'}],
		detail: [{keys: 'Enter', label: '开始挑战'}],
		reference: [{keys: 'Esc', label: '清空'}],
		recommend: [{keys: 'Enter', label: '查看推荐'}],
		skillmap: [{keys: 'g/G', label: '首尾'}],
		result: [{keys: 'r', label: '重试'}, {keys: 'n', label: '下一题'}],
	};
	const hints = [...screenSpecific[screen], ...globalHints];
	return compact ? hints.slice(0, 3) : hints;
}
