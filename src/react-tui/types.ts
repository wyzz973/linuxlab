export type ScreenID =
	| 'menu'
	| 'modules'
	| 'challenges'
	| 'detail'
	| 'reference'
	| 'recommend'
	| 'skillmap'
	| 'result';

export type ModalID = 'help' | 'commandPalette' | 'confirmExit' | 'errorDetails';

export type TerminalSize = {
	columns: number;
	rows: number;
};

export type LayoutMode = 'unsupported' | 'compact' | 'standard' | 'wide';

export type LayoutSpec = {
	mode: LayoutMode;
	columns: number;
	rows: number;
	contentWidth: number;
	mainWidth: number;
	navWidth: number;
	inspectorWidth: number;
	headerHeight: number;
	footerHeight: number;
	mainHeight: number;
	showSidebar: boolean;
	showInspector: boolean;
	showCompactTabs: boolean;
};

export type {
	Category,
	Challenge,
	ChallengeEntry,
	ChallengeRunEvent,
	ChallengeRunResult,
	CommandExample,
	CommandRef,
	LinuxLabData,
	ProgressData,
	ReferenceData,
	SkillEntry,
	VerifyResult,
	VerifyRule,
} from './domain/types.js';
