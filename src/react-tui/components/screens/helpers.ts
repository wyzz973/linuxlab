import {progressString} from '../common/ProgressBar.js';

export function stars(level: number) {
	const safeLevel = Math.max(1, Math.min(5, level));
	return `${'★'.repeat(safeLevel)}${'☆'.repeat(5 - safeLevel)}`;
}

export function compactProgress(passed: number, total: number, width = 8) {
	return `${progressString(passed, total, width)} ${passed}/${total}`;
}
