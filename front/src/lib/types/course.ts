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

// Task detail types
export interface Choice {
	text: string;
	valid: boolean;
}

export interface ProblemDetail extends Problem {
	language?: string; // for code problems
	default?: string; // for code problems
	choices?: Choice[]; // for multiple choice
	answer?: string; // for match problems
	limit?: number; // for multiple choice
}

export interface EnvironmentLimits {
	time: string;
	hardTime: string;
	memory: string;
}

export interface TaskDetail {
	id: string;
	courseId: string;
	name: string;
	author: string;
	contactUrl: string;
	context: string;
	environmentId: string;
	environmentType: string;
	limits?: EnvironmentLimits;
	networkGrading: boolean;
	problems: ProblemDetail[];
}
