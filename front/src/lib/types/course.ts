export interface Course {
	id: string;
	code: string;
	name: string;
	accessible: boolean;
	admins: string[];
	tutors: string[];
	taskCount: number;
}

export interface Problem {
	id: string;
	type: string;
	name: string;
	header: string;
}

export interface Task {
	id: string;
	name: string;
	author: string;
	environmentType: string;
	problems: Problem[];
}

export interface SummaryEntry {
	title: string;
	path?: string;
	children?: SummaryEntry[];
}

export interface Syllabus {
	title: string;
	author: string;
	summary: SummaryEntry[];
}

export interface CourseDetail extends Course {
	tasks: Task[];
	syllabus?: Syllabus;
}
