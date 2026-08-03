export type CursorAction = 'up' | 'down' | 'first' | 'last' | 'pageUp' | 'pageDown';

export function normalizePastedSearch(input: string): string {
	return input.replace(/\s+/g, ' ').trimStart();
}

export function updateQueryFromInput(query: string, input: string, key: {backspace?: boolean; delete?: boolean}) {
	if (key.backspace || key.delete) {
		return query.slice(0, -1);
	}
	return `${query}${normalizePastedSearch(input)}`;
}

export function moveCursor(current: number, count: number, action: CursorAction, pageSize: number) {
	if (count <= 0) {
		return 0;
	}

	switch (action) {
		case 'up':
			return Math.max(0, current - 1);
		case 'down':
			return Math.min(count - 1, current + 1);
		case 'first':
			return 0;
		case 'last':
			return count - 1;
		case 'pageUp':
			return Math.max(0, current - Math.max(1, pageSize));
		case 'pageDown':
			return Math.min(count - 1, current + Math.max(1, pageSize));
	}
}
