export type VerifyRule = {
	type: string;
	path?: string;
	expect?: string;
	command?: string;
};

export type Hint = {
	level: number;
	text: string;
};

export type Challenge = {
	id: string;
	title: string;
	difficulty: number;
	category: string;
	subcategory: string;
	tags: string[];
	description: string;
	hints: Hint[];
	verify: VerifyRule[];
	setup_files?: Array<{path: string; content: string}>;
	compose_file?: string;
	requires_docker?: boolean;
	dir?: string;
};

export type ChallengeEntry = {
	status: 'passed' | 'failed' | string;
	attempts: number;
	hints_used: number;
	last_attempt: string;
};

export type SkillEntry = {
	total: number;
	passed: number;
	score: number;
};

export type ProgressData = {
	skills: Record<string, SkillEntry>;
	challenges: Record<string, ChallengeEntry>;
};

export type CommandExample = {
	desc: string;
	cmd: string;
};

export type CommandRef = {
	name: string;
	brief: string;
	examples: CommandExample[];
	related_challenges?: string[];
};

export type ReferenceData = {
	commands: CommandRef[];
};

export type Category = {
	id: string;
	label: string;
	total: number;
	passed: number;
	challenges: Challenge[];
};

export type LinuxLabData = {
	categories: Category[];
	progress: ProgressData;
	references: ReferenceData;
};

export type VerifyResult = {
	Passed?: boolean;
	passed?: boolean;
	Message?: string;
	message?: string;
};

export type ChallengeRunEvent =
	| {type: 'setup'; message: string}
	| {type: 'handoff'; mode?: string; message: string}
	| {type: 'result'; passed: boolean; hintsUsed: number; results: VerifyResult[]}
	| {type: 'error'; message: string};

export type ChallengeRunResult = {
	challengeID: string;
	passed: boolean;
	hintsUsed: number;
	results: VerifyResult[];
	events: ChallengeRunEvent[];
};
