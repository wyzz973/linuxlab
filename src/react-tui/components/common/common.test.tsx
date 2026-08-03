import React from 'react';
import {renderToString} from 'ink';
import {describe, expect, it} from 'vitest';
import {ProgressBar} from './ProgressBar.js';
import {ScrollList} from './ScrollList.js';
import {Viewport} from './Viewport.js';

describe('common TUI components', () => {
	it('renders progress bar at requested width', () => {
		const output = renderToString(<ProgressBar value={3} total={10} width={10} />);
		expect(output.trim().length).toBe(10);
	});

	it('renders only visible rows in scroll list', () => {
		const output = renderToString(
			<ScrollList
				items={['one', 'two', 'three', 'four']}
				cursor={2}
				height={2}
				width={20}
				renderItem={(item, active) => `${active ? '›' : ' '} ${item}`}
			/>,
		);
		expect(output).toContain('three');
		expect(output).toContain('four');
		expect(output).not.toContain('one');
	});

	it('wraps viewport content within width and height', () => {
		const output = renderToString(<Viewport text="alpha beta gamma delta" width={8} height={2} />);
		expect(output.split('\n')).toHaveLength(2);
	});
});
