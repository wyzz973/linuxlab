import {describe, expect, it} from 'vitest';
import {parseRunEventLine} from './goExecutor.js';

describe('parseRunEventLine', () => {
	it('parses result events', () => {
		const event = parseRunEventLine('{"type":"result","passed":true,"hintsUsed":0,"results":[]}');
		expect(event.type).toBe('result');
		if (event.type === 'result') {
			expect(event.passed).toBe(true);
		}
	});
});
