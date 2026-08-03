import {describe, expect, it} from 'vitest';
import {moveCursor, normalizePastedSearch, updateQueryFromInput} from './input.js';

describe('input helpers', () => {
	it('normalizes pasted multi-line search text', () => {
		expect(normalizePastedSearch('  grep\n error\tapp.log')).toBe('grep error app.log');
	});

	it('appends printable search input', () => {
		expect(updateQueryFromInput('gr', 'ep', {})).toBe('grep');
	});

	it('removes one code unit on backspace', () => {
		expect(updateQueryFromInput('grep', '', {backspace: true})).toBe('gre');
	});

	it('moves cursor by single and page actions', () => {
		expect(moveCursor(2, 10, 'up', 4)).toBe(1);
		expect(moveCursor(2, 10, 'down', 4)).toBe(3);
		expect(moveCursor(5, 10, 'first', 4)).toBe(0);
		expect(moveCursor(5, 10, 'last', 4)).toBe(9);
		expect(moveCursor(5, 10, 'pageUp', 4)).toBe(1);
		expect(moveCursor(5, 10, 'pageDown', 4)).toBe(9);
	});
});
