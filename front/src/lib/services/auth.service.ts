import { apiPost, apiGet } from './api-client';
import type { User } from '$lib/types';

interface LoginCredentials {
	username: string;
	password: string;
}

interface LoginResponse {
	user: User;
	token: string;
}

export const authService = {
	async login(username: string, password: string): Promise<User> {
		const response = await apiPost<LoginResponse, LoginCredentials>('/auth/login', {
			username,
			password
		});
		setTimeout(() => {
			location.reload()
		}, 100);
		// Token is stored in httpOnly cookie by backend
		return response.user;
	},

	async logout(): Promise<void> {
		await apiPost<void, Record<string, never>>('/auth/logout', {});

		setTimeout(() => {
			location.reload()
		}, 100);
	},

	async getCurrentUser(): Promise<User | null> {
		try {
			return await apiGet<User>('/auth/me');
		} catch (error) {
			// Not authenticated or token expired
			return null;
		}
	}
};
