<script lang="ts">
	import { onMount } from 'svelte';
	import { courseService, ApiError } from '$lib/services';
	import type { CourseDetail } from '$lib/types';
	import { page } from '$app/state';
	import * as Card from '$lib/components/ui/card';

	let course: CourseDetail | null = null;
	let error: string | null = null;
	let loading = true;

	onMount(async () => {
		if (!page.params.id) {
			error = 'Course ID is missing in the URL';
			loading = false;
			return;
		}

		try {
			course = await courseService.getCourseById(page.params.id);
			error = null;
		} catch (err) {
			if (err instanceof ApiError) {
				error = `Failed to fetch course: ${err.message}`;
			} else {
				error = err instanceof Error ? err.message : 'An error occurred while fetching course';
			}
			console.error('Error fetching course:', err);
		} finally {
			loading = false;
		}
	});
</script>

{#if loading}
	<div class="container mx-auto px-4 py-8">
		<p class="text-muted-foreground">Loading course details...</p>
	</div>
{:else if error}
	<div class="container mx-auto px-4 py-8">
		<p class="text-red-500">Error: {error}</p>
	</div>
{:else if course}
	<div class="container mx-auto px-4 py-8">
		<!-- Header -->
		<div class="flex items-center justify-between mb-6">
			<div>
				<h1 class="text-3xl font-bold">{course.code}</h1>
				<p class="text-lg text-muted-foreground">{course.name}</p>
			</div>
			{#if course.accessible}
				<span class="px-3 py-1 text-sm rounded-full bg-green-100 text-green-800">Open</span>
			{:else}
				<span class="px-3 py-1 text-sm rounded-full bg-gray-100 text-gray-600">Closed</span>
			{/if}
		</div>

		<!-- Course Info -->
		<div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
			<!-- Admins -->
			<Card.Root>
				<Card.Header>
					<Card.Title>Admins</Card.Title>
				</Card.Header>
				<Card.Content>
					{#if course.admins && course.admins.length > 0}
						<ul class="space-y-1">
							{#each course.admins as admin (admin)}
								<li class="text-sm">{admin}</li>
							{/each}
						</ul>
					{:else}
						<p class="text-sm text-muted-foreground">No admins</p>
					{/if}
				</Card.Content>
			</Card.Root>

			<!-- Tutors -->
			<Card.Root>
				<Card.Header>
					<Card.Title>Tutors</Card.Title>
				</Card.Header>
				<Card.Content>
					{#if course.tutors && course.tutors.length > 0}
						<ul class="space-y-1">
							{#each course.tutors as tutor (tutor)}
								<li class="text-sm">{tutor}</li>
							{/each}
						</ul>
					{:else}
						<p class="text-sm text-muted-foreground">No tutors</p>
					{/if}
				</Card.Content>
			</Card.Root>
		</div>

		<!-- Syllabus -->
		{#if course.syllabus}
			<section class="mb-8">
				<h2 class="text-2xl font-semibold mb-4">Syllabus</h2>
				<Card.Root>
					<Card.Header>
						<Card.Title>{course.syllabus.title}</Card.Title>
						<Card.Description>by {course.syllabus.author}</Card.Description>
					</Card.Header>
					<Card.Content>
						{#if course.syllabus.summary && course.syllabus.summary.length > 0}
							<ul class="space-y-2">
								{#each course.syllabus.summary as entry (entry.title)}
									<li>
										<span class="font-medium">{entry.title}</span>
										{#if entry.children && entry.children.length > 0}
											<ul class="ml-4 mt-1 space-y-1">
												{#each entry.children as child (child.path || child.title)}
													<li class="text-sm text-muted-foreground">- {child.title}</li>
												{/each}
											</ul>
										{/if}
									</li>
								{/each}
							</ul>
						{/if}
					</Card.Content>
				</Card.Root>
			</section>
		{/if}

		<!-- Tasks -->
		<section>
			<h2 class="text-2xl font-semibold mb-4">Tasks ({course.taskCount})</h2>
			{#if course.tasks && course.tasks.length > 0}
				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
					{#each course.tasks as task (task.id)}
						<Card.Root class="hover:shadow-lg transition-shadow">
							<Card.Header>
								<Card.Title class="text-lg">{task.name}</Card.Title>
								<Card.Description>
									{task.environmentType === 'docker' ? 'Code' : 'Quiz'} - {task.problems.length} problem{task.problems.length !== 1 ? 's' : ''}
								</Card.Description>
							</Card.Header>
							<Card.Content>
								<p class="text-xs text-muted-foreground">Author: {task.author}</p>
							</Card.Content>
							<Card.Footer>
								<a href={`/courses/${course.id}/tasks/${task.id}`} class="text-primary hover:underline text-sm">
									View Task
								</a>
							</Card.Footer>
						</Card.Root>
					{/each}
				</div>
			{:else}
				<p class="text-muted-foreground">No tasks available</p>
			{/if}
		</section>
	</div>
{:else}
	<div class="container mx-auto px-4 py-8">
		<p>Course not found.</p>
	</div>
{/if}
