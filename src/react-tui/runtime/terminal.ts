import {useStdout} from 'ink';
import type {TerminalSize} from '../types.js';

export function useTerminalSize(fallback: TerminalSize = {columns: 80, rows: 24}): TerminalSize {
	const {stdout} = useStdout();
	const envColumns = Number(process.env.COLUMNS);
	const envRows = Number(process.env.LINES);
	return {
		columns: envColumns || stdout.columns || fallback.columns,
		rows: envRows || stdout.rows || fallback.rows,
	};
}
