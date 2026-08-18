import {describe, expect, it} from 'vitest';
import {buildRunArgs, parseRunEventLine} from './goExecutor.js';

describe('parseRunEventLine', () => {
	it('parses result events', () => {
		const event = parseRunEventLine('{"type":"result","passed":true,"hintsUsed":0,"results":[]}');
		expect(event.type).toBe('result');
		if (event.type === 'result') {
			expect(event.passed).toBe(true);
		}
	});
});

describe('buildRunArgs', () => {
	it('omits --hints when none were used', () => {
		expect(buildRunArgs('ls-basic')).toEqual(['challenge', 'run', 'ls-basic', '--json']);
		expect(buildRunArgs('ls-basic', 0)).toEqual(['challenge', 'run', 'ls-basic', '--json']);
	});

	it('appends --hints when hints were unlocked', () => {
		expect(buildRunArgs('ls-basic', 2)).toEqual(['challenge', 'run', 'ls-basic', '--json', '--hints', '2']);
	});
});
