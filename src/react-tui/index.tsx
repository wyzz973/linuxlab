import React from 'react';
import {render} from 'ink';
import {App} from './App.js';
import {loadLinuxLabData} from './data.js';
import {loadDoctorStatus} from './data/goData.js';

try {
	// 数据与后端健康检查并行加载；doctor 失败时仅省略状态徽标，不阻塞启动。
	const [data, doctor] = await Promise.all([
		loadLinuxLabData(),
		loadDoctorStatus().catch(() => undefined),
	]);
	const backendStatus = doctor === undefined
		? undefined
		: doctor.docker ? 'Docker 就绪' : '本地模式';
	render(<App data={data} backendStatus={backendStatus} />, {
		exitOnCtrlC: true,
		alternateScreen: true,
	});
} catch (error) {
	const message = error instanceof Error ? error.message : String(error);
	console.error(`加载 React TUI 失败: ${message}`);
	process.exit(1);
}
