function charWidth(char: string): number {
	const codePoint = char.codePointAt(0) ?? 0;
	if (codePoint === 0) {
		return 0;
	}
	if (codePoint < 32 || codePoint === 127) {
		return 0;
	}
	if (
		(codePoint >= 0x1100 && codePoint <= 0x115f) ||
		(codePoint >= 0x2329 && codePoint <= 0x232a) ||
		(codePoint >= 0x2e80 && codePoint <= 0xa4cf) ||
		(codePoint >= 0xac00 && codePoint <= 0xd7a3) ||
		(codePoint >= 0xf900 && codePoint <= 0xfaff) ||
		(codePoint >= 0xfe10 && codePoint <= 0xfe19) ||
		(codePoint >= 0xfe30 && codePoint <= 0xfe6f) ||
		(codePoint >= 0xff00 && codePoint <= 0xff60) ||
		(codePoint >= 0xffe0 && codePoint <= 0xffe6)
	) {
		return 2;
	}
	return 1;
}

export function displayWidth(input: string): number {
	return Array.from(input).reduce((width, char) => width + charWidth(char), 0);
}

function takeByWidth(input: string, width: number): string {
	let used = 0;
	let output = '';
	for (const char of Array.from(input)) {
		const next = used + charWidth(char);
		if (next > width) {
			break;
		}
		output += char;
		used = next;
	}
	return output;
}

export function fitText(input: string, width: number): string {
	const safeWidth = Math.max(0, Math.floor(width));
	if (safeWidth === 0) {
		return '';
	}
	if (displayWidth(input) <= safeWidth) {
		return input;
	}
	if (safeWidth === 1) {
		return '…';
	}
	return `${takeByWidth(input, safeWidth - 1)}…`;
}

export function wrapText(input: string, width: number): string[] {
	const safeWidth = Math.max(1, Math.floor(width));
	const normalized = input.replace(/\t/g, '    ').replace(/\r\n?/g, '\n');
	const wrapped: string[] = [];

	for (const paragraph of normalized.split('\n')) {
		const words = paragraph.replace(/\s+/g, ' ').trim().split(' ').filter(Boolean);
		if (words.length === 0) {
			wrapped.push('');
			continue;
		}

		let current = '';
		for (const word of words) {
			if (displayWidth(word) > safeWidth) {
				if (current !== '') {
					wrapped.push(current);
					current = '';
				}
				let rest = word;
				while (displayWidth(rest) > safeWidth) {
					const chunk = takeByWidth(rest, safeWidth);
					wrapped.push(chunk);
					rest = Array.from(rest).slice(Array.from(chunk).length).join('');
				}
				current = rest;
				continue;
			}

			const candidate = current === '' ? word : `${current} ${word}`;
			if (displayWidth(candidate) <= safeWidth) {
				current = candidate;
			} else {
				wrapped.push(current);
				current = word;
			}
		}

		if (current !== '') {
			wrapped.push(current);
		}
	}

	return wrapped.length > 0 ? wrapped : [''];
}
