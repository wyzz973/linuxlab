import type {ChallengeRunResult, LinuxLabData} from '../types.js';

// applyRunResultToData patches progress data with a finished run, mirroring
// what the Go engine persists in ~/.linuxlab/progress.json. It keeps the UI
// consistent immediately after a challenge and before the background reload
// from `linuxlab data dump --json` reconciles everything.
export function applyRunResultToData(data: LinuxLabData, result: ChallengeRunResult): LinuxLabData {
	const previous = data.progress.challenges[result.challengeID];
	const patched = {
		...data.progress.challenges,
		[result.challengeID]: {
			status: result.passed ? 'passed' : 'failed',
			attempts: (previous?.attempts ?? 0) + 1,
			hints_used: result.hintsUsed,
			last_attempt: new Date().toISOString(),
		},
	};

	const categories = data.categories.map(category => ({
		...category,
		passed: category.challenges.reduce(
			(count, challenge) => count + (patched[challenge.id]?.status === 'passed' ? 1 : 0),
			0,
		),
	}));

	return {
		...data,
		progress: {...data.progress, challenges: patched},
		categories,
	};
}
