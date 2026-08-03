import {readdir, readFile} from 'node:fs/promises';
import {existsSync} from 'node:fs';
import {homedir} from 'node:os';
import path from 'node:path';
import YAML from 'yaml';
import type {Category, Challenge, LinuxLabData, ProgressData, ReferenceData} from './types.js';
import {loadLinuxLabDataFromGo} from './data/goData.js';

const categoryLabels: Record<string, string> = {
	'linux-basics': 'Linux 基础命令',
	vim: 'Vim 操作',
	'shell-scripting': 'Shell 脚本',
	ops: '运维实战',
	containers: '容器与部署',
};

type LoadOptions = {
	rootDir?: string;
	challengesDir?: string;
	progressPath?: string;
	referencesPath?: string;
};

export async function loadLinuxLabData(options: LoadOptions = {}): Promise<LinuxLabData> {
	if (process.env.LINUXLAB_USE_GO_DATA === '1' && !options.rootDir && !options.challengesDir && !options.progressPath && !options.referencesPath) {
		try {
			return await loadLinuxLabDataFromGo();
		} catch {
			// Fall back to the TypeScript loader so the preview remains usable before the Go binary is built.
		}
	}

	const rootDir = options.rootDir ?? process.cwd();
	const challengesDir = options.challengesDir ?? path.join(rootDir, 'challenges');
	const progressPath = options.progressPath ?? path.join(homedir(), '.linuxlab', 'progress.json');
	const referencesPath = options.referencesPath ?? path.join(rootDir, 'references', 'commands.yaml');

	const [categories, progress, references] = await Promise.all([
		loadCategories(challengesDir, progressPath),
		loadProgress(progressPath),
		loadReferences(referencesPath),
	]);

	return {categories, progress, references};
}

async function loadCategories(challengesDir: string, progressPath: string): Promise<Category[]> {
	const progress = await loadProgress(progressPath);
	const categoryNames = await readdir(challengesDir, {withFileTypes: true});
	const categories: Category[] = [];

	for (const categoryDir of categoryNames) {
		if (!categoryDir.isDirectory()) {
			continue;
		}

		const categoryID = categoryDir.name;
		const challengeRoot = path.join(challengesDir, categoryID);
		const challengeDirs = await readdir(challengeRoot, {withFileTypes: true});
		const challenges: Challenge[] = [];

		for (const challengeDir of challengeDirs) {
			if (!challengeDir.isDirectory()) {
				continue;
			}

			const yamlPath = path.join(challengeRoot, challengeDir.name, 'challenge.yaml');
			if (!existsSync(yamlPath)) {
				continue;
			}

			const raw = await readFile(yamlPath, 'utf8');
			const parsed = YAML.parse(sanitizeYamlEscapes(raw)) as Partial<Challenge>;
			challenges.push(normalizeChallenge(parsed, path.dirname(yamlPath), categoryID));
		}

		challenges.sort((a, b) => a.id.localeCompare(b.id));
		const passed = challenges.filter(challenge => progress.challenges[challenge.id]?.status === 'passed').length;
		categories.push({
			id: categoryID,
			label: categoryLabels[categoryID] ?? categoryID,
			total: challenges.length,
			passed,
			challenges,
		});
	}

	categories.sort((a, b) => a.id.localeCompare(b.id));
	return categories;
}

function normalizeChallenge(parsed: Partial<Challenge>, dir: string, fallbackCategory: string): Challenge {
	return {
		id: parsed.id ?? path.basename(dir),
		title: parsed.title ?? path.basename(dir),
		difficulty: parsed.difficulty ?? 1,
		category: parsed.category ?? fallbackCategory,
		subcategory: parsed.subcategory ?? 'default',
		tags: parsed.tags ?? [],
		description: parsed.description ?? '',
		hints: parsed.hints ?? [],
		verify: parsed.verify ?? [],
		setup_files: parsed.setup_files,
		compose_file: parsed.compose_file,
		requires_docker: parsed.requires_docker,
		dir,
	};
}

async function loadProgress(progressPath: string): Promise<ProgressData> {
	if (!existsSync(progressPath)) {
		return {skills: {}, challenges: {}};
	}

	const raw = await readFile(progressPath, 'utf8');
	const parsed = JSON.parse(raw) as Partial<ProgressData>;
	return {
		skills: parsed.skills ?? {},
		challenges: parsed.challenges ?? {},
	};
}

async function loadReferences(referencesPath: string): Promise<ReferenceData> {
	if (!existsSync(referencesPath)) {
		return {commands: []};
	}

	const raw = await readFile(referencesPath, 'utf8');
	const parsed = YAML.parse(sanitizeYamlEscapes(raw)) as Partial<ReferenceData>;
	return {commands: parsed.commands ?? []};
}

function sanitizeYamlEscapes(raw: string) {
	// gopkg.in/yaml accepts some non-standard backslash escapes in quoted
	// strings. The JS YAML parser is stricter, so preserve those literals.
	return raw.replace(/(\\+)([^0abtnvfre "\\/N_LPuxU\r\n])/g, (match, slashes: string, char: string) => {
		if (slashes.length % 2 === 0) {
			return match;
		}
		return `${slashes}\\${char}`;
	});
}
