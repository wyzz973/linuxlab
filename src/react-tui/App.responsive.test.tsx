import React from 'react';
import {describe, expect, it} from 'vitest';
import {App} from './App.js';
import {fixtureData} from './test/fixtures.js';
import {renderTui} from './test/render.js';

describe('App responsive rendering', () => {
	it.each([
		['compact', 60, 20],
		['standard-min', 80, 24],
		['standard', 100, 30],
		['wide', 140, 40],
	])('renders %s layout without obvious overflow', (_name, columns, rows) => {
		const output = renderTui(<App data={fixtureData} terminalSize={{columns, rows}} />, columns);
		expect(output).toContain('LinuxLab');
		for (const line of output.split('\n')) {
			expect(line.length).toBeLessThanOrEqual(columns + 4);
		}
		expect(output).toMatchSnapshot();
	});
});
