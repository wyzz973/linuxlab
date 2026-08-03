import React from 'react';
import {Text} from 'ink';
import {theme} from '../../theme/theme.js';

type StatusBadgeProps = {
	status?: string;
};

export function statusLabel(status?: string) {
	if (status === 'passed') {
		return '已通过';
	}
	if (status === 'failed') {
		return '失败过';
	}
	return '未完成';
}

export function statusIcon(status?: string) {
	if (status === 'passed') {
		return '✓';
	}
	if (status === 'failed') {
		return '✗';
	}
	return '○';
}

export function statusColor(status?: string) {
	if (status === 'passed') {
		return theme.success;
	}
	if (status === 'failed') {
		return theme.error;
	}
	return theme.dim;
}

export function StatusBadge({status}: StatusBadgeProps) {
	return <Text color={statusColor(status)}>{statusIcon(status)} {statusLabel(status)}</Text>;
}
