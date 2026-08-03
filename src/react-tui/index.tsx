import React from 'react';
import {render} from 'ink';
import {App} from './App.js';
import {loadLinuxLabData} from './data.js';

try {
	const data = await loadLinuxLabData();
	render(<App data={data} />, {
		exitOnCtrlC: true,
		alternateScreen: true,
	});
} catch (error) {
	const message = error instanceof Error ? error.message : String(error);
	console.error(`加载 React TUI 失败: ${message}`);
	process.exit(1);
}
