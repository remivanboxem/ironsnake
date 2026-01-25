<script lang="ts">
	import { onMount } from 'svelte';
	import * as Card from '$lib/components/ui/card';
	import type { Course } from '$lib/types';
	import { courseService, ApiError } from '$lib/services';

	// Mock data for tasks/recently accessed
	const recentTasks = [
		{ id: 1, title: 'Apprendre JS', description: 'Learn JavaScript' },
		{ id: 2, title: 'Faire une DB', description: 'Create a Database' },
		{ id: 3, title: 'Coder proprement', description: 'Code properly' },
		{ id: 4, title: 'Math 1', description: 'Mathematics course 1' }
	];

	let courses: Course[] = [];
	let loading = true;
	let error: string | null = null;

	// Function to get initials from a username
	function getInitials(username: string): string {
		if (!username) return '?';
		return username.slice(0, 2).toUpperCase();
	}

	onMount(async () => {
		try {
			courses = await courseService.getAllCourses();
			error = null;
		} catch (err) {
			if (err instanceof ApiError) {
				error = `Failed to fetch courses: ${err.message}`;
			} else {
				error = err instanceof Error ? err.message : 'An error occurred while fetching courses';
			}
			console.error('Error fetching courses:', err);
		} finally {
			loading = false;
		}
	});
</script>

<div class="container mx-auto px-4 py-8">
	<!-- Tasks to do / Recently accessed Tasks -->
	<section class="mb-12">
		<h2 class="text-2xl font-semibold mb-6">Tasks to do/Recently accessed Tasks</h2>
		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
			{#each recentTasks as task (task.id)}
				<Card.Root class="hover:shadow-lg transition-shadow cursor-pointer">
					<Card.Content class="p-6">
						<h3 class="font-semibold text-lg text-center">{task.title}</h3>
					</Card.Content>
				</Card.Root>
			{/each}
		</div>
	</section>

	<!-- My courses -->
	<section>
		<h2 class="text-2xl font-semibold mb-6">My courses</h2>

		{#if loading}
			<div class="text-center py-8">
				<p class="text-muted-foreground">Loading courses...</p>
			</div>
		{:else if error}
			<div class="text-center py-8">
				<p class="text-red-500">Error: {error}</p>
			</div>
		{:else if courses.length === 0}
			<div class="text-center py-8">
				<p class="text-muted-foreground">No courses available</p>
			</div>
		{:else}
			<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
				{#each courses as course (course.id)}
					<Card.Root class="hover:shadow-lg transition-shadow cursor-pointer">
						<Card.Content class="p-6">
							<div class="flex items-center justify-between mb-2">
								<h3 class="font-bold text-xl">{course.code}</h3>
								{#if course.accessible}
									<span class="px-2 py-1 text-xs rounded-full bg-green-100 text-green-800">Open</span>
								{:else}
									<span class="px-2 py-1 text-xs rounded-full bg-gray-100 text-gray-600">Closed</span>
								{/if}
							</div>
							<p class="text-sm text-muted-foreground">{course.name}</p>

							<!-- Task count -->
							<p class="text-xs text-muted-foreground mt-2">
								{course.taskCount} task{course.taskCount !== 1 ? 's' : ''}
							</p>

							<!-- Admins -->
							{#if course.admins && course.admins.length > 0}
								<div class="flex items-center mt-4 gap-2">
									<div class="flex -space-x-2">
										{#each course.admins.slice(0, 3) as admin (admin)}
											<div
												class="w-8 h-8 rounded-full bg-primary flex items-center justify-center text-primary-foreground text-xs font-semibold border-2 border-background"
												title={admin}
											>
												{getInitials(admin)}
											</div>
										{/each}
										{#if course.admins.length > 3}
											<div class="w-8 h-8 rounded-full bg-muted flex items-center justify-center text-muted-foreground text-xs font-semibold border-2 border-background">
												+{course.admins.length - 3}
											</div>
										{/if}
									</div>
									<span class="text-xs text-muted-foreground">
										{course.admins.slice(0, 2).join(', ')}{course.admins.length > 2 ? '...' : ''}
									</span>
								</div>
							{/if}
						</Card.Content>
						<Card.Footer class="p-4 border-t text-center text-sm text-muted-foreground">
							<a href={`/courses/${course.id}`} class="text-primary hover:underline">View Course Details</a>
						</Card.Footer>
					</Card.Root>
				{/each}
			</div>
		{/if}
	</section>
</div>
