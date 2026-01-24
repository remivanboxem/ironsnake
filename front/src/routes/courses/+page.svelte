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

	// Function to get initials from first and last name
	function getInitials(firstName: string, lastName: string): string {
		const firstInitial = firstName?.charAt(0).toUpperCase() || '';
		const lastInitial = lastName?.charAt(0).toUpperCase() || '';
		return firstInitial + lastInitial;
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
					<Card.Root class="hover:shadow-lg transition-shadow cursor-pointer" >
						<Card.Content class="p-6">
							<h3 class="font-bold text-xl text-center">{course.code}</h3>
							<p class="text-sm text-muted-foreground text-center mt-2">{course.name}</p>
							
							<!-- Author Avatar and Info -->
							<div class="flex items-center justify-center mt-4 gap-2">
								<div class="w-8 h-8 rounded-full bg-primary flex items-center justify-center text-primary-foreground text-sm font-semibold">
									{getInitials(course.author.firstName, course.author.lastName)}
								</div>
								<span class="text-xs text-muted-foreground">
									{course.author.firstName} {course.author.lastName}
								</span>
							</div>
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
