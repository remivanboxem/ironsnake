<script lang="ts">
	import './layout.css';
	import { ModeWatcher } from "mode-watcher";
	import Header from '$lib/components/Header.svelte';

	let { children } = $props();
</script>

<ModeWatcher />
<div class="min-h-screen">
	<Header />
	<main>
		{@render children()}
	</main>
</div>
