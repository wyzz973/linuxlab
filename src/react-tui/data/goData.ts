import {execFile} from 'node:child_process';
import {promisify} from 'node:util';
import type {LinuxLabData} from '../types.js';

const execFileAsync = promisify(execFile);

export type DoctorStatus = {
	docker: boolean;
	challenges: number;
	references: number;
};

export async function loadLinuxLabDataFromGo(binaryPath = process.env.LINUXLAB_BIN ?? './linuxlab'): Promise<LinuxLabData> {
	const {stdout} = await execFileAsync(binaryPath, ['data', 'dump', '--json'], {
		cwd: process.cwd(),
		maxBuffer: 10 * 1024 * 1024,
	});
	return JSON.parse(stdout) as LinuxLabData;
}

// loadDoctorStatus reads the Go engine's environment health (Docker
// availability, challenge/reference counts) for the header status badge.
// Rejects when the binary is missing, so callers can fall back gracefully.
export async function loadDoctorStatus(binaryPath = process.env.LINUXLAB_BIN ?? './linuxlab'): Promise<DoctorStatus> {
	const {stdout} = await execFileAsync(binaryPath, ['doctor', '--json'], {
		cwd: process.cwd(),
		maxBuffer: 1 * 1024 * 1024,
	});
	return JSON.parse(stdout) as DoctorStatus;
}
