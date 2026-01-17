<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Avatar, AvatarFallback } from '$lib/components/ui/avatar';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { Search, Settings, BookOpen, Library, Activity, LogOut } from 'lucide-svelte';

	// Mock logged-in state - you can make this dynamic later
	let isLoggedIn = true;
</script>

<header class="w-full border-b bg-white">
	<div class="container mx-auto flex h-16 items-center justify-between gap-4 px-4">
		{#if isLoggedIn}
			<!-- Logo -->
			<div class="flex items-center justify-center w-12 h-12 bg-primary text-primary-foreground rounded-lg font-bold text-xl">
				IS
			</div>

			<!-- Search Bar -->
			<div class="flex-1 max-w-2xl relative">
				<Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
				<Input 
					type="search" 
					placeholder="Search bar" 
					class="w-full pl-10"
				/>
			</div>

			<!-- User Profile -->
			<DropdownMenu.Root>
				<DropdownMenu.Trigger class="cursor-pointer">
					<Avatar>
						<AvatarFallback class="bg-blue-600 text-white">U</AvatarFallback>
					</Avatar>
				</DropdownMenu.Trigger>
				<DropdownMenu.Content align="end" class="w-48 bg-gray-50 z-50">
					<DropdownMenu.Item class="cursor-pointer">
						<Settings class="mr-2 h-4 w-4" />
						Settings
					</DropdownMenu.Item>
					<DropdownMenu.Item class="cursor-pointer">
						<BookOpen class="mr-2 h-4 w-4" />
						My courses
					</DropdownMenu.Item>
					<DropdownMenu.Item class="cursor-pointer">
						<Library class="mr-2 h-4 w-4" />
						All classes
					</DropdownMenu.Item>
					<DropdownMenu.Item class="cursor-pointer">
						<Activity class="mr-2 h-4 w-4" />
						Service status
					</DropdownMenu.Item>
					<DropdownMenu.Separator />
					<DropdownMenu.Item class="cursor-pointer text-destructive">
						<LogOut class="mr-2 h-4 w-4" />
						Logout
					</DropdownMenu.Item>
				</DropdownMenu.Content>
			</DropdownMenu.Root>
		{:else}
			<h1 class="text-2xl font-bold text-primary">IronSnake</h1>
			<Button>Login</Button>
		{/if}
	</div>
</header>
