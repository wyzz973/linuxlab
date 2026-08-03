export const theme = {
	accent: '#4cc9c1',
	success: '#7ecf8f',
	warning: '#f3bc5f',
	error: '#ff7a7a',
	dim: '#81919c',
	border: '#33424c',
	text: undefined,
} as const;

export type ThemeColor = keyof typeof theme;
