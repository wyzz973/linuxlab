import {renderToString} from 'ink';
import type {ReactElement} from 'react';

export function renderTui(element: ReactElement, columns = 100) {
	return renderToString(element, {columns}).replace(/\x1B\[[0-?]*[ -/]*[@-~]/g, '');
}
