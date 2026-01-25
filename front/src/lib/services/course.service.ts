import type {
	Course,
	CourseDetail,
	TaskDetail,
	MCQSubmissionRequest,
	MCQSubmissionResponse
} from '$lib/types';
import { apiGet, apiPost } from './api-client';

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
	},

	/**
	 * Get a specific task by course ID and task ID
	 */
	async getTaskById(courseId: string, taskId: string): Promise<TaskDetail> {
		return apiGet<TaskDetail>(`/courses/${courseId}/tasks/${taskId}`);
	},

	/**
	 * Submit MCQ answers for grading
	 */
	async submitMCQ(
		courseId: string,
		taskId: string,
		submission: MCQSubmissionRequest
	): Promise<MCQSubmissionResponse> {
		return apiPost<MCQSubmissionResponse, MCQSubmissionRequest>(
			`/courses/${courseId}/tasks/${taskId}`,
			submission
		);
	}
};
