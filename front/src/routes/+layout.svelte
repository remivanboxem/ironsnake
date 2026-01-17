<script lang="ts">
	import './layout.css';
	import Header from '$lib/components/Header.svelte';

	let { children } = $props();
</script>

<div class="min-h-screen">
	<Header />
	<main>
		{@render children()}
	</main>
</div>
