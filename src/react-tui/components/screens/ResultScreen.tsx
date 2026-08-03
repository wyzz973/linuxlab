import React from 'react';
import {Box, Text} from 'ink';
import type {ChallengeRunResult, LayoutSpec, VerifyResult} from '../../types.js';
import {theme} from '../../theme/theme.js';
import {ScrollList} from '../common/ScrollList.js';

type ResultScreenProps = {
	result?: ChallengeRunResult;
	layout: LayoutSpec;
};

function resultPassed(result: VerifyResult) {
	return result.passed ?? result.Passed ?? false;
}

function resultMessage(result: VerifyResult) {
	return result.message ?? result.Message ?? '无消息';
}

export function ResultScreen({result, layout}: ResultScreenProps) {
	const width = layout.mainWidth;
	if (!result) {
		return (
			<Box flexDirection="column" width={width}>
				<Text bold color={theme.accent}>检测结果</Text>
				<Text color={theme.dim}>暂无结果</Text>
			</Box>
		);
	}

	return (
		<Box flexDirection="column" width={width}>
			<Text bold color={result.passed ? theme.success : theme.error}>{result.passed ? '挑战通过' : '挑战未通过'}</Text>
			<Text color={theme.dim}>题目 {result.challengeID} · 使用提示 {result.hintsUsed}</Text>
			<Box marginTop={1}>
				<ScrollList
					items={result.results}
					cursor={0}
					height={Math.max(3, layout.mainHeight - 4)}
					width={width}
					emptyTitle="未收到检查项"
					emptyAction="返回题目后可重试"
					renderItem={(item, active, index) => `${active ? '›' : ' '} ${resultPassed(item) ? '✓' : '✗'} 检查 ${index + 1}: ${resultMessage(item)}`}
				/>
			</Box>
			<Text color={theme.dim}>r 重试 · n 下一题 · q 返回</Text>
		</Box>
	);
}
