<script lang="ts">
	import * as Card from '$lib/components/ui/card';

	// Mock data for tasks/recently accessed
	const recentTasks = [
		{ id: 1, title: 'Apprendre JS', description: 'Learn JavaScript' },
		{ id: 2, title: 'Faire une DB', description: 'Create a Database' },
		{ id: 3, title: 'Coder proprement', description: 'Code properly' },
		{ id: 4, title: 'Math 1', description: 'Mathematics course 1' }
	];

	// Mock data for courses
	const courses = [
		{ id: 1, code: 'LINFO2255', title: 'Software Engineering' },
		{ id: 2, code: 'LINFO1006', title: 'Introduction to Computer Science' },
		{ id: 3, code: 'LEPL1401', title: 'Programming' }
	];
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
		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
			{#each courses as course (course.id)}
				<Card.Root class="hover:shadow-lg transition-shadow cursor-pointer">
					<Card.Content class="p-6">
						<h3 class="font-bold text-xl text-center">{course.code}</h3>
						<p class="text-sm text-muted-foreground text-center mt-2">{course.title}</p>
					</Card.Content>
				</Card.Root>
			{/each}
		</div>
	</section>
</div>
