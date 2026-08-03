import {describe, expect, it} from 'vitest';
import {displayWidth, fitText, wrapText} from './text.js';

describe('fitText', () => {
	it('truncates long text with ellipsis inside width', () => {
		expect(fitText('abcdefghijklmnopqrstuvwxyz', 10)).toBe('abcdefghi…');
	});

	it('returns empty string for non-positive width', () => {
		expect(fitText('abc', 0)).toBe('');
	});

	it('keeps short Chinese text unchanged', () => {
		expect(fitText('命令速查', 20)).toBe('命令速查');
	});
});

describe('wrapText', () => {
	it('wraps long text into bounded lines', () => {
		const lines = wrapText('alpha beta gamma delta', 8);
		expect(lines).toEqual(['alpha', 'beta', 'gamma', 'delta']);
	});

	it('splits a single overlong token', () => {
		const lines = wrapText('abcdefghijkl', 5);
		expect(lines).toEqual(['abcde', 'fghij', 'kl']);
	});

	it('counts Chinese characters as wide terminal cells', () => {
		expect(displayWidth('命令')).toBe(4);
		expect(fitText('命令速查面板', 8)).toBe('命令速…');
	});
});
