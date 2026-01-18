/**
 * Base API client configuration and utilities
 */

const API_BASE_URL = '/api';

export class ApiError extends Error {
	constructor(
		message: string,
		public status: number,
		public statusText: string
	) {
		super(message);
		this.name = 'ApiError';
	}
}

async function handleResponse<T>(response: Response): Promise<T> {
	if (!response.ok) {
		throw new ApiError(
			`API request failed: ${response.statusText}`,
			response.status,
			response.statusText
		);
	}

	return response.json();
}

export async function apiGet<T>(endpoint: string): Promise<T> {
	const response = await fetch(`${API_BASE_URL}${endpoint}`, {
		credentials: 'include'
	});
	return handleResponse<T>(response);
}

export async function apiPost<T, D = unknown>(endpoint: string, data: D): Promise<T> {
	const response = await fetch(`${API_BASE_URL}${endpoint}`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify(data),
		credentials: 'include'
	});
	return handleResponse<T>(response);
}

export async function apiPut<T, D = unknown>(endpoint: string, data: D): Promise<T> {
	const response = await fetch(`${API_BASE_URL}${endpoint}`, {
		method: 'PUT',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify(data),
		credentials: 'include'
	});
	return handleResponse<T>(response);
}

export async function apiDelete<T>(endpoint: string): Promise<T> {
	const response = await fetch(`${API_BASE_URL}${endpoint}`, {
		method: 'DELETE',
		credentials: 'include'
	});
	return handleResponse<T>(response);
}
