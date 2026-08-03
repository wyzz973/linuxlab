import {spawnSync} from 'node:child_process';
import type {ChallengeRunEvent, ChallengeRunResult} from '../types.js';

export type RunChallengeOptions = {
	binaryPath?: string;
	challengeID: string;
	cwd?: string;
	onEvent?: (event: ChallengeRunEvent) => void;
};

export function parseRunEventLine(line: string): ChallengeRunEvent {
	return JSON.parse(line) as ChallengeRunEvent;
}

export async function runChallenge(options: RunChallengeOptions): Promise<ChallengeRunResult> {
	const binary = options.binaryPath ?? process.env.LINUXLAB_BIN ?? './linuxlab';
	const stdin = process.stdin as NodeJS.ReadStream & {setRawMode?: (mode: boolean) => void};
	const restoreRawMode = Boolean(stdin.isTTY && stdin.setRawMode);

	if (restoreRawMode) {
		stdin.setRawMode?.(false);
	}

	try {
		const child = spawnSync(binary, ['challenge', 'run', options.challengeID, '--json'], {
			cwd: options.cwd ?? process.cwd(),
			stdio: ['inherit', 'pipe', 'inherit'],
			encoding: 'utf8',
		});

		if (child.error) {
			throw child.error;
		}

		const result = resultFromStdout(options.challengeID, child.stdout ?? '', options.onEvent);
		if ((child.status ?? 0) !== 0 && !result.seenResult) {
			throw new Error(`挑战执行失败，退出码 ${child.status}`);
		}
		return result.result;
	} finally {
		if (restoreRawMode) {
			stdin.setRawMode?.(true);
		}
	}
}

function resultFromStdout(
	challengeID: string,
	stdout: string,
	onEvent?: (event: ChallengeRunEvent) => void,
): {result: ChallengeRunResult; seenResult: boolean} {
	const events: ChallengeRunEvent[] = [];
	let result: ChallengeRunResult | undefined;

	for (const line of stdout.split('\n')) {
		if (line.trim() === '') {
			continue;
		}
		try {
			const event = parseRunEventLine(line);
			events.push(event);
			onEvent?.(event);
			if (event.type === 'result') {
				result = {
					challengeID,
					passed: event.passed,
					hintsUsed: event.hintsUsed,
					results: event.results,
					events,
				};
			}
		} catch {
			// Ignore non-JSON terminal output; the Go CLI keeps machine events on stdout.
		}
	}

	return {
		seenResult: result !== undefined,
		result: result ?? {
			challengeID,
			passed: false,
			hintsUsed: 0,
			results: [{passed: false, message: '未收到检测结果'}],
			events,
		},
	};
}
