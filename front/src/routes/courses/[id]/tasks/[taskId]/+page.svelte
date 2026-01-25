<script lang="ts">
	import { onMount } from 'svelte';
	import { courseService, codeService, ApiError, type RunCodeResponse } from '$lib/services';
	import type { TaskDetail } from '$lib/types';
	import { page } from '$app/state';
	import * as Card from '$lib/components/ui/card';
	import { Markdown } from '$lib/components/ui/markdown';
	import { EditorView, basicSetup } from 'codemirror';
	import { python } from '@codemirror/lang-python';
	import { oneDark } from '@codemirror/theme-one-dark';
	import { EditorState } from '@codemirror/state';
	import { SvelteMap } from 'svelte/reactivity';

	let task: TaskDetail | null = $state(null);
	let error: string | null = $state(null);
	let loading = $state(true);

	// Store editor instances and code values per problem
	let editors: Map<string, EditorView> = new SvelteMap();
	let codeValues: Map<string, string> = new SvelteMap();

	// Store execution state per problem
	let runningProblems: Set<string> = $state(new Set());
	let outputs: Map<string, RunCodeResponse> = $state(new Map());

	function initEditor(element: HTMLElement, params: [string, string]) {
		const [problemId, initialCode] = params;
		if (editors.has(problemId)) return;

		codeValues.set(problemId, initialCode);

		const state = EditorState.create({
			doc: initialCode,
			extensions: [
				basicSetup,
				python(),
				oneDark,
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
			}
		};
	}

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
				[problemId, { output: '', error: err instanceof Error ? err.message : 'Unknown error', exitCode: 1 }]
			]);
		} finally {
			runningProblems = new Set([...runningProblems].filter((id) => id !== problemId));
		}
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
			task = await courseService.getTaskById(courseId, taskId);
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
		<nav class="text-sm mb-4">
			<a href="/" class="text-primary hover:underline">Home</a>
			<span class="mx-2 text-muted-foreground">/</span>
			<a href={`/courses/${task.courseId}`} class="text-primary hover:underline">{task.courseId}</a>
			<span class="mx-2 text-muted-foreground">/</span>
			<span class="text-muted-foreground">{task.id}</span>
		</nav>

		<!-- Header -->
		<div class="flex items-start justify-between mb-6">
			<div>
				<h1 class="text-3xl font-bold">{task.name}</h1>
				<p class="text-muted-foreground mt-1">by {task.author}</p>
			</div>
			<span class="px-3 py-1 text-sm rounded-full {task.environmentType === 'docker' ? 'bg-blue-100 text-blue-800' : 'bg-purple-100 text-purple-800'}">
				{task.environmentType === 'docker' ? 'Code Task' : 'Quiz'}
			</span>
		</div>

		<!-- Task Info -->
		<div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
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
					<p class="text-xs text-muted-foreground">problem{task.problems.length !== 1 ? 's' : ''} in this task</p>
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
			<h2 class="text-2xl font-semibold mb-4">Problems</h2>
			<div class="space-y-4">
				{#each task.problems as problem, index (problem.id)}
					<Card.Root>
						<Card.Header>
							<div class="flex items-center justify-between">
								<div class="flex items-center gap-3">
									<span class="text-lg font-bold text-muted-foreground">#{index + 1}</span>
									<Card.Title>{problem.name}</Card.Title>
								</div>
								<span class="px-2 py-1 text-xs rounded-full {getProblemTypeColor(problem.type)}">
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
									<div class="flex items-center justify-between mb-2">
										<p class="text-sm font-medium">Code Editor ({problem.language}):</p>
										<button
											onclick={() => runCode(problem.id, problem.language || 'python')}
											disabled={isRunning}
											class="flex items-center gap-2 px-4 py-2 bg-green-600 hover:bg-green-700 disabled:bg-green-800 disabled:cursor-not-allowed text-white text-sm font-medium rounded-md transition-colors"
										>
											{#if isRunning}
												<svg class="animate-spin" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
													<path d="M21 12a9 9 0 1 1-6.219-8.56"/>
												</svg>
												Running...
											{:else}
												<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
													<path d="M8 5v14l11-7z"/>
												</svg>
												Run
											{/if}
										</button>
									</div>
									<div
										class="rounded-md overflow-hidden border border-border"
										use:initEditor={[problem.id, problem.default]}
									></div>

									<!-- Output display -->
									{#if output}
										<div class="mt-4">
											<p class="text-sm font-medium mb-2">Output:</p>
											<div class="bg-zinc-900 text-zinc-100 p-4 rounded-md font-mono text-sm overflow-x-auto">
												{#if output.error}
													<div class="text-red-400">{output.error}</div>
												{/if}
												{#if output.output}
													<pre class="whitespace-pre-wrap">{output.output}</pre>
												{/if}
												<div class="mt-2 pt-2 border-t border-zinc-700 text-xs text-zinc-500">
													Exit code: {output.exitCode}
												</div>
											</div>
										</div>
									{/if}
								</div>
							{/if}

							{#if problem.type === 'multiple_choice' && problem.choices}
								<div class="mt-4">
									<p class="text-sm font-medium mb-2">Choices:</p>
									<ul class="space-y-2">
										{#each problem.choices as choice, choiceIndex (choiceIndex)}
											<li class="flex items-center gap-2 p-2 rounded-md bg-muted/50">
												<span class="w-6 h-6 rounded-full bg-muted flex items-center justify-center text-xs font-medium">
													{String.fromCharCode(65 + choiceIndex)}
												</span>
												<span>{choice.text}</span>
											</li>
										{/each}
									</ul>
								</div>
							{/if}

							{#if problem.type === 'match'}
								<div class="mt-4 p-3 bg-muted/50 rounded-md">
									<p class="text-sm text-muted-foreground">This is a fill-in-the-blank question.</p>
								</div>
							{/if}
						</Card.Content>
					</Card.Root>
				{/each}
			</div>
		</section>

		<!-- Contact -->
		{#if task.contactUrl}
			<div class="mt-8 p-4 bg-muted rounded-lg">
				<p class="text-sm">
					<span class="font-medium">Need help?</span>
					<a href={task.contactUrl} class="text-primary hover:underline ml-2">Contact the author</a>
				</p>
			</div>
		{/if}
	</div>
{:else}
	<div class="container mx-auto px-4 py-8">
		<p>Task not found.</p>
	</div>
{/if}
