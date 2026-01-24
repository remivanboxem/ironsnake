<script lang="ts">
  import { onMount } from 'svelte';
	import { courseService, ApiError } from '$lib/services';
	import type { Course } from '../../../lib/types';
  import { page } from '$app/state';

  let course: Course | null = null;
  let error: string | null = null;
  let loading = true;
  
	onMount(async () => {
    if(!page.params.id) {
      error = 'Course ID is missing in the URL';
      loading = false;
      return;
    }

		try {
			course = await courseService.getCourseById(page.params.id);
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

{#if loading}
  <p>Loading course details...</p>
{:else if error}
  <p class="text-red-500">Error: {error}</p>
{:else if course}
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-4">{course.name}</h1>
    <p class="mb-6">{course.description}</p>
    <h2 class="text-2xl font-semibold mb-4">Author</h2>
    <p>{course.author.firstName} {course.author.lastName}</p>
    <!-- Add more course details as needed -->
  </div>
{:else}
  <p>Course not found.</p>
{/if}