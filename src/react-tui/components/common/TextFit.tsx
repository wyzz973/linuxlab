import React from 'react';
import {Text} from 'ink';
import {fitText} from '../../utils/text.js';

type TextFitProps = {
	children: string;
	width: number;
	color?: string;
	bold?: boolean;
	dimColor?: boolean;
};

export function TextFit({children, width, color, bold, dimColor}: TextFitProps) {
	return (
		<Text color={color} bold={bold} dimColor={dimColor}>
			{fitText(children, width)}
		</Text>
	);
}
