<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Avatar, AvatarFallback } from '$lib/components/ui/avatar';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { Search, Settings, BookOpen, Library, Activity, LogOut } from 'lucide-svelte';
	import SunIcon from '@lucide/svelte/icons/sun';
	import MoonIcon from '@lucide/svelte/icons/moon';

	import { toggleMode } from 'mode-watcher';

	// Mock logged-in state - you can make this dynamic later
	let isLoggedIn = true;
</script>

<header class="w-full border-b">
	<div class="container mx-auto flex h-16 items-center justify-between gap-4 px-4">
		{#if isLoggedIn}
			<!-- Logo -->
			 <a href="/">
				<div
					class="flex h-12 w-12 items-center justify-center rounded-lg bg-primary text-xl font-bold text-primary-foreground"
				>
					IS
				</div>
			</a>

			<!-- Search Bar -->
			<div class="relative max-w-2xl flex-1">
				<Search class="absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
				<Input type="search" placeholder="Search bar" class="w-full pl-10" />
			</div>

			<Button onclick={toggleMode} variant="outline" size="icon">
				<SunIcon
					class="h-[1.2rem] w-[1.2rem] scale-100 rotate-0 transition-all! dark:scale-0 dark:-rotate-90"
				/>
				<MoonIcon
					class="absolute h-[1.2rem] w-[1.2rem] scale-0 rotate-90 transition-all! dark:scale-100 dark:rotate-0"
				/>
				<span class="sr-only">Toggle theme</span>
			</Button>

			<!-- User Profile -->
			<DropdownMenu.Root>
				<DropdownMenu.Trigger class="cursor-pointer">
					<Avatar>
						<AvatarFallback>U</AvatarFallback>
					</Avatar>
				</DropdownMenu.Trigger>
				<DropdownMenu.Content align="end" class="z-50 w-48">
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
