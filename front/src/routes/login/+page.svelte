<script lang="ts">
	import { goto } from '$app/navigation';
	import { authService } from '$lib/services';
	import { authStore } from '$lib/stores/auth.store';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { resolve } from '$app/paths';

	let username = '';
	let password = '';
	let loading = false;
	let error = '';

	async function handleLogin() {
		if (!username || !password) {
			error = 'Username and password are required';
			return;
		}

		loading = true;
		error = '';

		try {
			const user = await authService.login(username, password);
			authStore.setUser(user);
			goto(resolve('/'));
		} catch (err: unknown) {
			if (err instanceof Error) {
				error = err.message || 'Login failed';
			} else {
				error = 'Login failed';
			}
		} finally {
			loading = false;
		}
	}

	function handleSubmit(e: Event) {
		e.preventDefault();
		handleLogin();
	}
</script>

<div class="flex h-[90vh] items-center justify-center overflow-hidden bg-background p-4">
	<Card.Root class="w-full max-w-md">
		<Card.Header>
			<Card.Title class="text-2xl">Login to IronSnake</Card.Title>
			<Card.Description>Enter your credentials to access the system</Card.Description>
		</Card.Header>
		<Card.Content>
			<form on:submit={handleSubmit} class="space-y-4">
				<div class="space-y-2">
					<Label for="username">Username</Label>
					<Input
						id="username"
						type="text"
						placeholder="Enter your username"
						bind:value={username}
						disabled={loading}
						required
					/>
				</div>
				<div class="space-y-2">
					<Label for="password">Password</Label>
					<Input
						id="password"
						type="password"
						placeholder="Enter your password"
						bind:value={password}
						disabled={loading}
						required
					/>
				</div>
				{#if error}
					<div class="rounded-md bg-destructive/15 p-3 text-sm text-destructive">
						{error}
					</div>
				{/if}
				<Button type="submit" class="w-full" disabled={loading}>
					{loading ? 'Logging in...' : 'Login'}
				</Button>
			</form>
		</Card.Content>
	</Card.Root>
</div>
