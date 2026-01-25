import { apiPost } from './api-client';

export interface RunCodeRequest {
	code: string;
	language: string;
}

export interface RunCodeResponse {
	output: string;
	error?: string;
	exitCode: number;
}

/**
 * Code execution service for running code
 */
export const codeService = {
	/**
	 * Run code and get the output
	 */
	async runCode(code: string, language: string): Promise<RunCodeResponse> {
		return apiPost<RunCodeResponse, RunCodeRequest>('/run', { code, language });
	}
};
