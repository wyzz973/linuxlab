import {execFile} from 'node:child_process';
import {promisify} from 'node:util';
import type {LinuxLabData} from '../types.js';

const execFileAsync = promisify(execFile);

export async function loadLinuxLabDataFromGo(binaryPath = process.env.LINUXLAB_BIN ?? './linuxlab'): Promise<LinuxLabData> {
	const {stdout} = await execFileAsync(binaryPath, ['data', 'dump', '--json'], {
		cwd: process.cwd(),
		maxBuffer: 10 * 1024 * 1024,
	});
	return JSON.parse(stdout) as LinuxLabData;
}
