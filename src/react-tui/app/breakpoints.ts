import type {LayoutSpec, TerminalSize} from '../types.js';

const MIN_COLUMNS = 60;
const MIN_ROWS = 18;
const HEADER_HEIGHT = 1;
const FOOTER_HEIGHT = 1;

function clamp(value: number, min: number, max: number) {
	return Math.max(min, Math.min(max, value));
}

export function createLayoutSpec(size: TerminalSize): LayoutSpec {
	const columns = Math.max(1, Math.floor(size.columns || 80));
	const rows = Math.max(1, Math.floor(size.rows || 24));
	const mainHeight = Math.max(1, rows - HEADER_HEIGHT - FOOTER_HEIGHT - 4);

	if (columns < MIN_COLUMNS || rows < MIN_ROWS) {
		return {
			mode: 'unsupported',
			columns,
			rows,
			contentWidth: columns,
			mainWidth: columns,
			navWidth: 0,
			inspectorWidth: 0,
			headerHeight: HEADER_HEIGHT,
			footerHeight: FOOTER_HEIGHT,
			mainHeight,
			showSidebar: false,
			showInspector: false,
			showCompactTabs: false,
		};
	}

	if (columns < 80 || rows < 24) {
		return {
			mode: 'compact',
			columns,
			rows,
			contentWidth: columns,
			mainWidth: columns,
			navWidth: 0,
			inspectorWidth: 0,
			headerHeight: HEADER_HEIGHT,
			footerHeight: FOOTER_HEIGHT,
			mainHeight,
			showSidebar: false,
			showInspector: false,
			showCompactTabs: true,
		};
	}

	const navWidth = clamp(Math.floor(columns * 0.2), 18, 28);
	const isWide = columns >= 120 && rows >= 30;
	const inspectorWidth = isWide ? clamp(Math.floor(columns * 0.24), 28, 40) : 0;
	const gutters = isWide ? 4 : 2;
	const mainWidth = Math.max(20, columns - navWidth - inspectorWidth - gutters);

	return {
		mode: isWide ? 'wide' : 'standard',
		columns,
		rows,
		contentWidth: columns,
		mainWidth,
		navWidth,
		inspectorWidth,
		headerHeight: HEADER_HEIGHT,
		footerHeight: FOOTER_HEIGHT,
		mainHeight,
		showSidebar: true,
		showInspector: isWide,
		showCompactTabs: false,
	};
}
