import type { Course } from '$lib/types';
import { apiGet } from './api-client';

/**
 * Course service for managing course-related API calls
 */
export const courseService = {
	/**
	 * Get all courses
	 */
	async getAllCourses(): Promise<Course[]> {
		return apiGet<Course[]>('/courses');
	},

	/**
	 * Get a specific course by ID
	 */
	async getCourseById(id: string): Promise<Course> {
		return apiGet<Course>(`/courses/${id}`);
	}
};
