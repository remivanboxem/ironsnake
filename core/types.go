package main

// CourseResponse represents the JSON response structure for a course
type CourseResponse struct {
	ID         string   `json:"id"`
	Code       string   `json:"code"`
	Name       string   `json:"name"`
	Accessible bool     `json:"accessible"`
	Admins     []string `json:"admins"`
	Tutors     []string `json:"tutors"`
	TaskCount  int      `json:"taskCount"`
}

// TaskResponse represents a task in the API response
type TaskResponse struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Author          string            `json:"author"`
	EnvironmentType string            `json:"environmentType"`
	Problems        []ProblemResponse `json:"problems"`
}

// ProblemResponse represents a problem in the API response
type ProblemResponse struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Header string `json:"header"`
}

// ProblemDetailResponse includes full problem details
type ProblemDetailResponse struct {
	ProblemResponse
	Language string   `json:"language,omitempty"` // for code problems
	Default  string   `json:"default,omitempty"`  // for code problems
	Choices  []Choice `json:"choices,omitempty"`  // for multiple choice
	Limit    int      `json:"limit,omitempty"`    // for multiple choice
	// Note: Answer field is intentionally not included to prevent exposing correct answers
}

// Choice represents a multiple choice option (for API response - does not expose correct answer)
type Choice struct {
	Text string `json:"text"`
}

// EnvironmentLimits represents resource limits
type EnvironmentLimits struct {
	Time     string `json:"time"`
	HardTime string `json:"hardTime"`
	Memory   string `json:"memory"`
}

// TaskDetailResponse represents full task details
type TaskDetailResponse struct {
	ID              string                  `json:"id"`
	CourseID        string                  `json:"courseId"`
	Name            string                  `json:"name"`
	Author          string                  `json:"author"`
	ContactURL      string                  `json:"contactUrl"`
	Context         string                  `json:"context"`
	EnvironmentID   string                  `json:"environmentId"`
	EnvironmentType string                  `json:"environmentType"`
	Limits          *EnvironmentLimits      `json:"limits,omitempty"`
	NetworkGrading  bool                    `json:"networkGrading"`
	Problems        []ProblemDetailResponse `json:"problems"`
}

// CourseDetailResponse represents the full course detail
type CourseDetailResponse struct {
	CourseResponse
	Tasks    []TaskResponse    `json:"tasks"`
	Syllabus *SyllabusResponse `json:"syllabus,omitempty"`
}

// SyllabusResponse represents the syllabus in the API response
type SyllabusResponse struct {
	Title   string         `json:"title"`
	Author  string         `json:"author"`
	Summary []SummaryEntry `json:"summary"`
}

// SummaryEntry represents a syllabus entry
type SummaryEntry struct {
	Title    string         `json:"title"`
	Path     string         `json:"path,omitempty"`
	Children []SummaryEntry `json:"children,omitempty"`
}

// RunCodeRequest represents the request body for code execution
type RunCodeRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

// RunCodeResponse represents the response from code execution
type RunCodeResponse struct {
	Output   string `json:"output"`
	Error    string `json:"error,omitempty"`
	ExitCode int    `json:"exitCode"`
}

// MCQSubmissionRequest represents a student's MCQ submission
type MCQSubmissionRequest struct {
	// Answers maps problem ID to selected choice indices (for multiple_choice)
	// or text answer (for match problems)
	Answers map[string]MCQAnswer `json:"answers"`
}

// MCQAnswer represents an answer to a single problem
type MCQAnswer struct {
	// SelectedIndices contains the indices of selected choices (for multiple_choice)
	SelectedIndices []int `json:"selectedIndices,omitempty"`
	// TextAnswer contains the text answer (for match problems)
	TextAnswer string `json:"textAnswer,omitempty"`
}

// MCQSubmissionResponse represents the result of an MCQ submission
type MCQSubmissionResponse struct {
	Score   float64                   `json:"score"`   // Score as percentage (0-100)
	Results map[string]ProblemResult  `json:"results"` // Results per problem
	Total   int                       `json:"total"`   // Total number of problems
	Correct int                       `json:"correct"` // Number of correct answers
}

// ProblemResult represents the result for a single problem
type ProblemResult struct {
	Correct bool `json:"correct"`
}
