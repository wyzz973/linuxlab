import {describe, expect, it} from 'vitest';
import {createLayoutSpec} from './breakpoints.js';

describe('createLayoutSpec', () => {
	it('uses unsupported mode below minimum terminal size', () => {
		const spec = createLayoutSpec({columns: 50, rows: 16});
		expect(spec.mode).toBe('unsupported');
		expect(spec.contentWidth).toBe(50);
		expect(spec.mainHeight).toBeGreaterThanOrEqual(1);
	});

	it('uses compact mode for small integrated terminals', () => {
		const spec = createLayoutSpec({columns: 60, rows: 20});
		expect(spec.mode).toBe('compact');
		expect(spec.showSidebar).toBe(false);
		expect(spec.showCompactTabs).toBe(true);
		expect(spec.showInspector).toBe(false);
		expect(spec.mainWidth).toBeLessThanOrEqual(60);
	});

	it('uses standard mode around 100 columns', () => {
		const spec = createLayoutSpec({columns: 100, rows: 30});
		expect(spec.mode).toBe('standard');
		expect(spec.showSidebar).toBe(true);
		expect(spec.navWidth).toBeGreaterThanOrEqual(18);
		expect(spec.showInspector).toBe(false);
		expect(spec.mainWidth + spec.navWidth).toBeLessThanOrEqual(100);
	});

	it('uses wide mode and exposes inspector budget', () => {
		const spec = createLayoutSpec({columns: 140, rows: 40});
		expect(spec.mode).toBe('wide');
		expect(spec.showSidebar).toBe(true);
		expect(spec.showInspector).toBe(true);
		expect(spec.inspectorWidth).toBeGreaterThanOrEqual(28);
		expect(spec.mainWidth + spec.navWidth + spec.inspectorWidth).toBeLessThanOrEqual(140);
	});
});
