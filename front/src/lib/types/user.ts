export interface User {
	id: string;
	username: string;
	email: string;
	firstName: string;
	lastName: string;
	authProvider: 'ldap' | 'local';
}
