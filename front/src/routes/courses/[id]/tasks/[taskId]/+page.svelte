<script lang="ts">
	import { onMount } from 'svelte';
	import { courseService, codeService, ApiError, type RunCodeResponse } from '$lib/services';
	import type { TaskDetail, CourseDetail, MCQSubmissionResponse, MCQAnswer } from '$lib/types';
	import { page } from '$app/state';
	import * as Card from '$lib/components/ui/card';
	import { Markdown } from '$lib/components/ui/markdown';
	import { EditorView, basicSetup } from 'codemirror';
	import { python } from '@codemirror/lang-python';
	import { oneDark } from '@codemirror/theme-one-dark';
	import { EditorState, StateEffect, Compartment } from '@codemirror/state';
	import { mode } from 'mode-watcher';
	import { SvelteMap } from 'svelte/reactivity';

	let task: TaskDetail | null = $state(null);
	let course: CourseDetail | null = $state(null);
	let error: string | null = $state(null);
	let loading = $state(true);

	// Store editor instances and code values per problem
	let editors: Map<string, EditorView> = new SvelteMap();
	let codeValues: Map<string, string> = new SvelteMap();
	let themeCompartments: Map<string, Compartment> = new SvelteMap();

	// Store execution state per problem
	let runningProblems: Set<string> = $state(new Set());
	let outputs: Map<string, RunCodeResponse> = $state(new Map());

	// MCQ state
	let mcqAnswers: Map<string, MCQAnswer> = $state(new SvelteMap());
	let mcqSubmitting = $state(false);
	let mcqResult: MCQSubmissionResponse | null = $state(null);
	let mcqError: string | null = $state(null);

	// Track the current mode for reactive updates
	let currentMode = $derived(mode.current);

	// Check if this task has MCQ problems
	let hasMCQProblems = $derived(
		task?.problems.some((p) => p.type === 'multiple_choice' || p.type === 'match') ?? false
	);

	// Check if this is an MCQ-only environment
	let isMCQEnvironment = $derived(task?.environmentType === 'mcq');

	// Get all MCQ problem IDs
	let mcqProblemIds = $derived(
		task?.problems
			.filter((p) => p.type === 'multiple_choice' || p.type === 'match')
			.map((p) => p.id) ?? []
	);

	// Check if all MCQ questions have been answered
	let allMCQAnswered = $derived(() => {
		if (!task) return false;
		for (const problemId of mcqProblemIds) {
			const answer = mcqAnswers.get(problemId);
			if (!answer) return false;
			const problem = task.problems.find((p) => p.id === problemId);
			if (!problem) return false;
			if (problem.type === 'multiple_choice') {
				if (!answer.selectedIndices || answer.selectedIndices.length === 0) return false;
			} else if (problem.type === 'match') {
				if (!answer.textAnswer || answer.textAnswer.trim() === '') return false;
			}
		}
		return true;
	});

	// Count answered questions
	let answeredCount = $derived(() => {
		if (!task) return 0;
		let count = 0;
		for (const problemId of mcqProblemIds) {
			const answer = mcqAnswers.get(problemId);
			if (!answer) continue;
			const problem = task.problems.find((p) => p.id === problemId);
			if (!problem) continue;
			if (problem.type === 'multiple_choice') {
				if (answer.selectedIndices && answer.selectedIndices.length > 0) count++;
			} else if (problem.type === 'match') {
				if (answer.textAnswer && answer.textAnswer.trim() !== '') count++;
			}
		}
		return count;
	});

	// Check if submission is allowed
	let canSubmit = $derived(isMCQEnvironment ? allMCQAnswered() : mcqAnswers.size > 0);

	function initEditor(element: HTMLElement, params: [string, string]) {
		const [problemId, initialCode] = params;
		if (editors.has(problemId)) return;

		codeValues.set(problemId, initialCode);

		// Create a compartment for the theme so it can be reconfigured
		const themeCompartment = new Compartment();
		themeCompartments.set(problemId, themeCompartment);

		// Get current theme based on mode - use oneDark for dark mode, default light theme for light mode
		const themeExtension = mode.current === 'dark' ? oneDark : [];

		const state = EditorState.create({
			doc: initialCode,
			extensions: [
				basicSetup,
				python(),
				themeCompartment.of(themeExtension),
				EditorView.updateListener.of((update) => {
					if (update.docChanged) {
						codeValues.set(problemId, update.state.doc.toString());
					}
				}),
				EditorView.theme({
					'&': { height: '300px' },
					'.cm-scroller': { overflow: 'auto' }
				})
			]
		});

		const editor = new EditorView({
			state,
			parent: element
		});

		editors.set(problemId, editor);

		return {
			destroy() {
				editor.destroy();
				editors.delete(problemId);
				themeCompartments.delete(problemId);
			}
		};
	}

	// Update all editors when theme changes
	$effect(() => {
		const themeExtension = currentMode === 'dark' ? oneDark : [];

		editors.forEach((editor, problemId) => {
			const compartment = themeCompartments.get(problemId);
			if (compartment) {
				editor.dispatch({
					effects: compartment.reconfigure(themeExtension)
				});
			}
		});
	});

	async function runCode(problemId: string, language: string) {
		const code = codeValues.get(problemId) || '';

		// Mark as running
		runningProblems = new Set([...runningProblems, problemId]);

		try {
			const result = await codeService.runCode(code, language);
			outputs = new Map([...outputs, [problemId, result]]);
		} catch (err) {
			outputs = new Map([
				...outputs,
				[
					problemId,
					{ output: '', error: err instanceof Error ? err.message : 'Unknown error', exitCode: 1 }
				]
			]);
		} finally {
			runningProblems = new Set([...runningProblems].filter((id) => id !== problemId));
		}
	}

	// MCQ functions
	function toggleChoice(problemId: string, choiceIndex: number) {
		const current = mcqAnswers.get(problemId) || { selectedIndices: [] };
		const indices = current.selectedIndices || [];
		const newIndices = indices.includes(choiceIndex)
			? indices.filter((i) => i !== choiceIndex)
			: [...indices, choiceIndex];
		mcqAnswers.set(problemId, { selectedIndices: newIndices });
		// Trigger reactivity
		mcqAnswers = new Map(mcqAnswers);
	}

	function isChoiceSelected(problemId: string, choiceIndex: number): boolean {
		const answer = mcqAnswers.get(problemId);
		return answer?.selectedIndices?.includes(choiceIndex) ?? false;
	}

	function updateMatchAnswer(problemId: string, text: string) {
		mcqAnswers.set(problemId, { textAnswer: text });
		mcqAnswers = new Map(mcqAnswers);
	}

	function getMatchAnswer(problemId: string): string {
		return mcqAnswers.get(problemId)?.textAnswer ?? '';
	}

	async function submitMCQ() {
		if (!task) return;

		mcqSubmitting = true;
		mcqError = null;

		try {
			const answers: Record<string, MCQAnswer> = {};
			mcqAnswers.forEach((answer, problemId) => {
				answers[problemId] = answer;
			});

			mcqResult = await courseService.submitMCQ(task.courseId, task.id, { answers });
		} catch (err) {
			if (err instanceof ApiError) {
				mcqError = `Failed to submit: ${err.message}`;
			} else {
				mcqError = err instanceof Error ? err.message : 'An error occurred while submitting';
			}
			console.error('Error submitting MCQ:', err);
		} finally {
			mcqSubmitting = false;
		}
	}

	function getProblemResult(problemId: string): boolean | null {
		if (!mcqResult) return null;
		return mcqResult.results[problemId]?.correct ?? null;
	}

	function getProblemTypeLabel(type: string): string {
		switch (type) {
			case 'code':
				return 'Code';
			case 'multiple_choice':
				return 'Multiple Choice';
			case 'match':
				return 'Fill-in';
			default:
				return type;
		}
	}

	function getProblemTypeColor(type: string): string {
		switch (type) {
			case 'code':
				return 'bg-blue-100 text-blue-800';
			case 'multiple_choice':
				return 'bg-purple-100 text-purple-800';
			case 'match':
				return 'bg-orange-100 text-orange-800';
			default:
				return 'bg-gray-100 text-gray-800';
		}
	}

	onMount(async () => {
		const courseId = page.params.id;
		const taskId = page.params.taskId;

		if (!courseId || !taskId) {
			error = 'Course ID or Task ID is missing in the URL';
			loading = false;
			return;
		}

		try {
			[task, course] = await Promise.all([
				courseService.getTaskById(courseId, taskId),
				courseService.getCourseById(courseId)
			]);
			error = null;
		} catch (err) {
			if (err instanceof ApiError) {
				error = `Failed to fetch task: ${err.message}`;
			} else {
				error = err instanceof Error ? err.message : 'An error occurred while fetching task';
			}
			console.error('Error fetching task:', err);
		} finally {
			loading = false;
		}
	});
</script>

{#if loading}
	<div class="container mx-auto px-4 py-8">
		<p class="text-muted-foreground">Loading task details...</p>
	</div>
{:else if error}
	<div class="container mx-auto px-4 py-8">
		<p class="text-red-500">Error: {error}</p>
	</div>
{:else if task}
	<div class="container mx-auto px-4 py-8">
		<!-- Breadcrumb -->
		<nav class="mb-4 text-sm">
			<a href="/" class="text-primary hover:underline">Home</a>
			<span class="mx-2 text-muted-foreground">/</span>
			<a href={`/courses/${task.courseId}`} class="text-primary hover:underline"
				>{course?.name ?? task.courseId}</a
			>
			<span class="mx-2 text-muted-foreground">/</span>
			<span class="text-muted-foreground">{task.name}</span>
		</nav>

		<!-- Header -->
		<div class="mb-6 flex items-start justify-between">
			<div>
				<h1 class="text-3xl font-bold">{task.name}</h1>
				<p class="mt-1 text-muted-foreground">by {task.author}</p>
			</div>
			<span
				class="rounded-full px-3 py-1 text-sm {task.environmentType === 'docker'
					? 'bg-blue-100 text-blue-800'
					: 'bg-purple-100 text-purple-800'}"
			>
				{task.environmentType === 'docker' ? 'Code Task' : 'Quiz'}
			</span>
		</div>

		<!-- Task Info -->
		<div class="mb-8 grid grid-cols-1 gap-4 md:grid-cols-3">
			<!-- Environment Info -->
			<Card.Root>
				<Card.Header class="pb-2">
					<Card.Title class="text-sm font-medium">Environment</Card.Title>
				</Card.Header>
				<Card.Content>
					<p class="text-lg font-semibold">{task.environmentId}</p>
					<p class="text-xs text-muted-foreground">{task.environmentType}</p>
				</Card.Content>
			</Card.Root>

			<!-- Limits -->
			{#if task.limits}
				<Card.Root>
					<Card.Header class="pb-2">
						<Card.Title class="text-sm font-medium">Resource Limits</Card.Title>
					</Card.Header>
					<Card.Content>
						<div class="space-y-1 text-sm">
							<p><span class="text-muted-foreground">Time:</span> {task.limits.time}s</p>
							<p><span class="text-muted-foreground">Memory:</span> {task.limits.memory} MB</p>
							{#if task.limits.hardTime}
								<p><span class="text-muted-foreground">Hard Time:</span> {task.limits.hardTime}s</p>
							{/if}
						</div>
					</Card.Content>
				</Card.Root>
			{/if}

			<!-- Stats -->
			<Card.Root>
				<Card.Header class="pb-2">
					<Card.Title class="text-sm font-medium">Problems</Card.Title>
				</Card.Header>
				<Card.Content>
					<p class="text-3xl font-bold">{task.problems.length}</p>
					<p class="text-xs text-muted-foreground">
						problem{task.problems.length !== 1 ? 's' : ''} in this task
					</p>
				</Card.Content>
			</Card.Root>
		</div>

		<!-- Context -->
		{#if task.context}
			<Card.Root class="mb-8">
				<Card.Header>
					<Card.Title>Context</Card.Title>
				</Card.Header>
				<Card.Content>
					<Markdown content={task.context} />
				</Card.Content>
			</Card.Root>
		{/if}

		<!-- Problems -->
		<section>
			<h2 class="mb-4 text-2xl font-semibold">Problems</h2>
			<div class="space-y-4">
				{#each task.problems as problem, index (problem.id)}
					<Card.Root>
						<Card.Header>
							<div class="flex items-center justify-between">
								<div class="flex items-center gap-3">
									<span class="text-lg font-bold text-muted-foreground">#{index + 1}</span>
									<Card.Title>{problem.name}</Card.Title>
								</div>
								<span class="rounded-full px-2 py-1 text-xs {getProblemTypeColor(problem.type)}">
									{getProblemTypeLabel(problem.type)}
								</span>
							</div>
						</Card.Header>
						<Card.Content>
							<!-- Problem header/description -->
							<div class="mb-4">
								<Markdown content={problem.header} />
							</div>

							<!-- Type-specific content -->
							{#if problem.type === 'code' && problem.default}
								{@const isRunning = runningProblems.has(problem.id)}
								{@const output = outputs.get(problem.id)}
								<div class="mt-4">
									<div class="mb-2 flex items-center justify-between">
										<p class="text-sm font-medium">Code Editor ({problem.language}):</p>
										<button
											onclick={() => runCode(problem.id, problem.language || 'python')}
											disabled={isRunning}
											class="flex items-center gap-2 rounded-md bg-green-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-green-700 disabled:cursor-not-allowed disabled:bg-green-800"
										>
											{#if isRunning}
												<svg
													class="animate-spin"
													xmlns="http://www.w3.org/2000/svg"
													width="16"
													height="16"
													viewBox="0 0 24 24"
													fill="none"
													stroke="currentColor"
													stroke-width="2"
												>
													<path d="M21 12a9 9 0 1 1-6.219-8.56" />
												</svg>
												Running...
											{:else}
												<svg
													xmlns="http://www.w3.org/2000/svg"
													width="16"
													height="16"
													viewBox="0 0 24 24"
													fill="currentColor"
												>
													<path d="M8 5v14l11-7z" />
												</svg>
												Run
											{/if}
										</button>
									</div>
									<div
										class="overflow-hidden rounded-md border border-border"
										use:initEditor={[problem.id, problem.default]}
									></div>

									<!-- Output display -->
									{#if output}
										<div class="mt-4">
											<p class="mb-2 text-sm font-medium">Output:</p>
											<div
												class="overflow-x-auto rounded-md bg-zinc-900 p-4 font-mono text-sm text-zinc-100"
											>
												{#if output.error}
													<div class="text-red-400">{output.error}</div>
												{/if}
												{#if output.output}
													<pre class="whitespace-pre-wrap">{output.output}</pre>
												{/if}
												<div class="mt-2 border-t border-zinc-700 pt-2 text-xs text-zinc-500">
													Exit code: {output.exitCode}
												</div>
											</div>
										</div>
									{/if}
								</div>
							{/if}

							{#if problem.type === 'multiple_choice' && problem.choices}
								{@const result = getProblemResult(problem.id)}
								<div class="mt-4">
									<p class="mb-2 text-sm font-medium">
										Select your answer{problem.limit === 0 || (problem.limit && problem.limit > 1)
											? 's'
											: ''}:
									</p>
									<ul class="space-y-2">
										{#each problem.choices as choice, choiceIndex (choiceIndex)}
											{@const isSelected = isChoiceSelected(problem.id, choiceIndex)}
											<li>
												<button
													type="button"
													onclick={() => toggleChoice(problem.id, choiceIndex)}
													class="flex w-full cursor-pointer items-center gap-3 rounded-md p-3 text-left transition-colors
														{isSelected ? 'bg-primary/20 ring-2 ring-primary' : 'bg-muted/50 hover:bg-muted'}"
												>
													<span
														class="flex h-6 w-6 shrink-0 items-center justify-center rounded-full text-xs font-medium
															{isSelected ? 'bg-primary text-primary-foreground' : 'bg-muted'}"
													>
														{String.fromCharCode(65 + choiceIndex)}
													</span>
													<span class="flex-1">{choice.text}</span>
													{#if isSelected}
														<svg
															xmlns="http://www.w3.org/2000/svg"
															width="20"
															height="20"
															viewBox="0 0 24 24"
															fill="none"
															stroke="currentColor"
															stroke-width="2"
															stroke-linecap="round"
															stroke-linejoin="round"
															class="text-primary"
														>
															<polyline points="20 6 9 17 4 12"></polyline>
														</svg>
													{/if}
												</button>
											</li>
										{/each}
									</ul>
									{#if result !== null}
										<div
											class="mt-3 flex items-center gap-2 rounded-md p-2 {result
												? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400'
												: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400'}"
										>
											{#if result}
												<svg
													xmlns="http://www.w3.org/2000/svg"
													width="20"
													height="20"
													viewBox="0 0 24 24"
													fill="none"
													stroke="currentColor"
													stroke-width="2"
													stroke-linecap="round"
													stroke-linejoin="round"
												>
													<path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
													<polyline points="22 4 12 14.01 9 11.01"></polyline>
												</svg>
												<span class="font-medium">Correct!</span>
											{:else}
												<svg
													xmlns="http://www.w3.org/2000/svg"
													width="20"
													height="20"
													viewBox="0 0 24 24"
													fill="none"
													stroke="currentColor"
													stroke-width="2"
													stroke-linecap="round"
													stroke-linejoin="round"
												>
													<circle cx="12" cy="12" r="10"></circle>
													<line x1="15" y1="9" x2="9" y2="15"></line>
													<line x1="9" y1="9" x2="15" y2="15"></line>
												</svg>
												<span class="font-medium">Incorrect</span>
											{/if}
										</div>
									{/if}
								</div>
							{/if}

							{#if problem.type === 'match'}
								{@const result = getProblemResult(problem.id)}
								<div class="mt-4">
									<label for="match-{problem.id}" class="mb-2 block text-sm font-medium">
										Your answer:
									</label>
									<input
										id="match-{problem.id}"
										type="text"
										value={getMatchAnswer(problem.id)}
										oninput={(e) => updateMatchAnswer(problem.id, e.currentTarget.value)}
										placeholder="Type your answer here..."
										class="w-full rounded-md border border-border bg-background px-3 py-2 text-sm focus:ring-2 focus:ring-primary focus:outline-none"
									/>
									{#if result !== null}
										<div
											class="mt-3 flex items-center gap-2 rounded-md p-2 {result
												? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400'
												: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400'}"
										>
											{#if result}
												<svg
													xmlns="http://www.w3.org/2000/svg"
													width="20"
													height="20"
													viewBox="0 0 24 24"
													fill="none"
													stroke="currentColor"
													stroke-width="2"
													stroke-linecap="round"
													stroke-linejoin="round"
												>
													<path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
													<polyline points="22 4 12 14.01 9 11.01"></polyline>
												</svg>
												<span class="font-medium">Correct!</span>
											{:else}
												<svg
													xmlns="http://www.w3.org/2000/svg"
													width="20"
													height="20"
													viewBox="0 0 24 24"
													fill="none"
													stroke="currentColor"
													stroke-width="2"
													stroke-linecap="round"
													stroke-linejoin="round"
												>
													<circle cx="12" cy="12" r="10"></circle>
													<line x1="15" y1="9" x2="9" y2="15"></line>
													<line x1="9" y1="9" x2="15" y2="15"></line>
												</svg>
												<span class="font-medium">Incorrect</span>
											{/if}
										</div>
									{/if}
								</div>
							{/if}
						</Card.Content>
					</Card.Root>
				{/each}
			</div>

			<!-- MCQ Submit Button and Results -->
			{#if hasMCQProblems}
				<div class="mt-8">
					{#if mcqResult}
						<!-- Results Summary -->
						<div
							class="rounded-lg border p-6 {mcqResult.score >= 50
								? 'border-green-500 bg-green-50 dark:bg-green-900/20'
								: 'border-red-500 bg-red-50 dark:bg-red-900/20'}"
						>
							<h3 class="mb-4 text-xl font-semibold">Results</h3>
							<div class="grid grid-cols-3 gap-4 text-center">
								<div>
									<p
										class="text-3xl font-bold {mcqResult.score >= 50
											? 'text-green-600 dark:text-green-400'
											: 'text-red-600 dark:text-red-400'}"
									>
										{mcqResult.score.toFixed(0)}%
									</p>
									<p class="text-sm text-muted-foreground">Score</p>
								</div>
								<div>
									<p class="text-3xl font-bold text-green-600 dark:text-green-400">
										{mcqResult.correct}
									</p>
									<p class="text-sm text-muted-foreground">Correct</p>
								</div>
								<div>
									<p class="text-3xl font-bold">{mcqResult.total}</p>
									<p class="text-sm text-muted-foreground">Total</p>
								</div>
							</div>
						</div>
					{/if}

					<!-- Submit Button (always visible for resubmission) -->
					<div class="mt-4">
						{#if isMCQEnvironment && !allMCQAnswered()}
							<p class="mb-2 text-center text-sm text-muted-foreground">
								Please answer all questions before submitting ({answeredCount()}/{mcqProblemIds.length})
							</p>
						{/if}
						<button
							onclick={submitMCQ}
							disabled={mcqSubmitting || !canSubmit}
							class="flex w-full items-center justify-center gap-2 rounded-md bg-primary px-6 py-3 text-lg font-medium text-primary-foreground transition-colors hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-50"
						>
							{#if mcqSubmitting}
								<svg
									class="animate-spin"
									xmlns="http://www.w3.org/2000/svg"
									width="20"
									height="20"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
								>
									<path d="M21 12a9 9 0 1 1-6.219-8.56" />
								</svg>
								Submitting...
							{:else if mcqResult}
								Resubmit Answers
							{:else}
								Submit Answers
							{/if}
						</button>
						{#if mcqError}
							<p class="mt-2 text-center text-sm text-red-500">{mcqError}</p>
						{/if}
					</div>
				</div>
			{/if}
		</section>

		<!-- Contact -->
		{#if task.contactUrl}
			<div class="mt-8 rounded-lg bg-muted p-4">
				<p class="text-sm">
					<span class="font-medium">Need help?</span>
					<a href={task.contactUrl} class="ml-2 text-primary hover:underline">Contact the author</a>
				</p>
			</div>
		{/if}
	</div>
{:else}
	<div class="container mx-auto px-4 py-8">
		<p>Task not found.</p>
	</div>
{/if}
