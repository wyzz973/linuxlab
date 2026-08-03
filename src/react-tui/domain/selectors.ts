import type {Category, Challenge, CommandRef} from './types.js';

export function summarizeCategories(categories: Category[]) {
	return {
		total: categories.reduce((sum, category) => sum + category.total, 0),
		passed: categories.reduce((sum, category) => sum + category.passed, 0),
		failed: 0,
	};
}

export function countFailedChallenges(categories: Category[], progress: Record<string, {status?: string}>): number {
	return categories
		.flatMap(category => category.challenges)
		.filter(challenge => progress[challenge.id]?.status === 'failed').length;
}

export function filterCategories(categories: Category[], query: string): Category[] {
	const normalized = query.trim().toLowerCase();
	if (normalized === '') {
		return categories;
	}
	return categories.filter(category => {
		const haystack = [category.id, category.label, ...category.challenges.map(challenge => challenge.title)]
			.join(' ')
			.toLowerCase();
		return haystack.includes(normalized);
	});
}

export function filterChallenges(challenges: Challenge[], query: string): Challenge[] {
	const normalized = query.trim().toLowerCase();
	if (normalized === '') {
		return challenges;
	}
	return challenges.filter(challenge => {
		const haystack = [
			challenge.id,
			challenge.title,
			challenge.subcategory,
			challenge.description,
			...challenge.tags,
		].join(' ').toLowerCase();
		return haystack.includes(normalized);
	});
}

export function filterReferences(commands: CommandRef[], query: string): CommandRef[] {
	const normalized = query.trim().toLowerCase();
	if (normalized === '') {
		return commands;
	}
	return commands.filter(command => {
		const haystack = [
			command.name,
			command.brief,
			...command.examples.flatMap(example => [example.desc, example.cmd]),
		].join(' ').toLowerCase();
		return haystack.includes(normalized);
	});
}

export function recommendChallenges(categories: Category[], progress: Record<string, {status?: string}>): Challenge[] {
	const all = categories.flatMap(category => category.challenges);
	const failed = all.filter(challenge => progress[challenge.id]?.status === 'failed');
	const incomplete = all.filter(challenge => progress[challenge.id]?.status !== 'passed');
	return [...failed, ...incomplete.filter(challenge => !failed.includes(challenge))].slice(0, 8);
}
