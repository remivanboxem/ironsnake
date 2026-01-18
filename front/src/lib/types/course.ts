export interface Author {
	id: string;
	username: string;
	firstName: string;
	lastName: string;
}

export interface Course {
	id: string;
	code: string;
	name: string;
	description: string;
	academicYear: string;
	createdBy: string;
	author: Author;
	createdAt: string;
}
