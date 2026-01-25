import type { Course, CourseDetail } from '$lib/types';
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
	 * Get a specific course by ID (returns full details including tasks)
	 */
	async getCourseById(id: string): Promise<CourseDetail> {
		return apiGet<CourseDetail>(`/courses/${id}`);
	}
};
