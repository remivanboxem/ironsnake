import { writable } from 'svelte/store';
import type { User } from '$lib/types';

interface AuthState {
	user: User | null;
	isAuthenticated: boolean;
	isLoading: boolean;
}

function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>({
		user: null,
		isAuthenticated: false,
		isLoading: true
	});

	return {
		subscribe,
		setUser: (user: User | null) => {
			set({
				user,
				isAuthenticated: !!user,
				isLoading: false
			});
		},
		logout: () => {
			set({
				user: null,
				isAuthenticated: false,
				isLoading: false
			});
		},
		setLoading: (loading: boolean) => {
			update((state) => ({ ...state, isLoading: loading }));
		}
	};
}

export const authStore = createAuthStore();
