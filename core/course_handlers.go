package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"ironsnake/core/courseparser"
)

// CoursesDir is the path to the courses directory
var CoursesDir = "courses"

func getCoursesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	loader := courseparser.NewCourseLoader()
	courses, err := loader.LoadAllCourses(CoursesDir)
	if err != nil {
		http.Error(w, "Failed to load courses", http.StatusInternalServerError)
		log.Printf("Error loading courses: %v", err)
		return
	}

	response := make([]CourseResponse, len(courses))
	for i, course := range courses {
		response[i] = CourseResponse{
			ID:         course.CourseID,
			Code:       course.CourseID,
			Name:       course.Config.Name,
			Accessible: course.Config.Accessible,
			Admins:     course.Config.Admins,
			Tutors:     course.Config.Tutors,
			TaskCount:  len(course.Tasks),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

// submitMCQHandler handles POST requests for MCQ submissions
func submitMCQHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract IDs from URL path (format: /courses/:courseID/tasks/:taskID)
	path := r.URL.Path
	remaining := path[len("/courses/"):]

	// Find the /tasks/ separator
	tasksIdx := -1
	for i := 0; i < len(remaining)-6; i++ {
		if remaining[i:i+7] == "/tasks/" {
			tasksIdx = i
			break
		}
	}

	if tasksIdx == -1 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	courseID := remaining[:tasksIdx]
	taskID := remaining[tasksIdx+7:]

	if courseID == "" || taskID == "" {
		http.Error(w, "Course ID and Task ID are required", http.StatusBadRequest)
		return
	}

	// Parse the request body
	var submission MCQSubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Error parsing submission: %v", err)
		return
	}

	// Load the course
	loader := courseparser.NewCourseLoader()
	coursePath := filepath.Join(CoursesDir, courseID)
	course, err := loader.LoadCourse(coursePath)
	if err != nil {
		http.Error(w, "Course not found", http.StatusNotFound)
		log.Printf("Error loading course %s: %v", courseID, err)
		return
	}

	// Find the task
	task, ok := course.Tasks[taskID]
	if !ok {
		http.Error(w, "Task not found", http.StatusNotFound)
		log.Printf("Task %s not found in course %s", taskID, courseID)
		return
	}

	// Grade the submission
	results := make(map[string]ProblemResult)
	correctCount := 0
	totalProblems := 0

	for _, op := range task.Problems.Problems {
		problemID := op.ID
		problem := op.Problem
		answer, hasAnswer := submission.Answers[problemID]

		var isCorrect bool

		switch p := problem.(type) {
		case *courseparser.MultipleChoiceProblem:
			totalProblems++
			if hasAnswer {
				isCorrect = gradeMultipleChoice(p, answer.SelectedIndices)
			}
		case *courseparser.MatchProblem:
			totalProblems++
			if hasAnswer {
				isCorrect = gradeMatch(p, answer.TextAnswer)
			}
		default:
			// Skip non-gradable problems (e.g., code problems)
			continue
		}

		if isCorrect {
			correctCount++
		}
		results[problemID] = ProblemResult{Correct: isCorrect}
	}

	// Calculate score
	var score float64
	if totalProblems > 0 {
		score = float64(correctCount) / float64(totalProblems) * 100
	}

	// Build response
	response := MCQSubmissionResponse{
		Score:   score,
		Results: results,
		Total:   totalProblems,
		Correct: correctCount,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

// gradeMultipleChoice checks if the selected choices match the correct answers
func gradeMultipleChoice(problem *courseparser.MultipleChoiceProblem, selectedIndices []int) bool {
	if len(problem.Choices) == 0 {
		return false
	}

	// Build a set of correct indices
	correctIndices := make(map[int]bool)
	for i, choice := range problem.Choices {
		if choice.Valid {
			correctIndices[i] = true
		}
	}

	// Build a set of selected indices
	selectedSet := make(map[int]bool)
	for _, idx := range selectedIndices {
		// Validate index is within bounds
		if idx < 0 || idx >= len(problem.Choices) {
			return false
		}
		selectedSet[idx] = true
	}

	// Check if selected matches correct exactly
	if len(selectedSet) != len(correctIndices) {
		return false
	}

	for idx := range correctIndices {
		if !selectedSet[idx] {
			return false
		}
	}

	return true
}

// gradeMatch checks if the text answer matches the expected answer
func gradeMatch(problem *courseparser.MatchProblem, textAnswer string) bool {
	// Case-insensitive comparison, trimming whitespace
	expected := strings.TrimSpace(strings.ToLower(problem.Answer))
	given := strings.TrimSpace(strings.ToLower(textAnswer))
	return expected == given
}

func getCourseByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path (format: /courses/:id)
	path := r.URL.Path
	courseID := path[len("/courses/"):]

	if courseID == "" {
		http.Error(w, "Course ID is required", http.StatusBadRequest)
		return
	}

	// Load the specific course
	loader := courseparser.NewCourseLoader()
	coursePath := filepath.Join(CoursesDir, courseID)
	course, err := loader.LoadCourse(coursePath)
	if err != nil {
		http.Error(w, "Course not found", http.StatusNotFound)
		log.Printf("Error loading course %s: %v", courseID, err)
		return
	}

	// Build task responses
	tasks := make([]TaskResponse, 0, len(course.Tasks))
	for taskID, task := range course.Tasks {
		problems := make([]ProblemResponse, 0, task.Problems.Len())
		for _, op := range task.Problems.Problems {
			problems = append(problems, ProblemResponse{
				ID:     op.ID,
				Type:   op.Problem.GetType(),
				Name:   op.Problem.GetName(),
				Header: op.Problem.GetHeader(),
			})
		}

		tasks = append(tasks, TaskResponse{
			ID:              taskID,
			Name:            task.Name,
			Author:          task.Author,
			EnvironmentType: task.EnvironmentType,
			Problems:        problems,
		})
	}

	// Build response
	response := CourseDetailResponse{
		CourseResponse: CourseResponse{
			ID:         course.CourseID,
			Code:       course.CourseID,
			Name:       course.Config.Name,
			Accessible: course.Config.Accessible,
			Admins:     course.Config.Admins,
			Tutors:     course.Config.Tutors,
			TaskCount:  len(course.Tasks),
		},
		Tasks: tasks,
	}

	// Add syllabus if present
	if course.Syllabus != nil {
		response.Syllabus = &SyllabusResponse{
			Title:   course.Syllabus.Book.Book.Title,
			Author:  course.Syllabus.Book.Book.Author,
			Summary: convertSummary(course.Syllabus.Summary),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

func convertSummary(entries []courseparser.SummaryEntry) []SummaryEntry {
	result := make([]SummaryEntry, len(entries))
	for i, entry := range entries {
		result[i] = SummaryEntry{
			Title:    entry.Title,
			Path:     entry.Path,
			Children: convertSummary(entry.Children),
		}
	}
	return result
}

func getTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract IDs from URL path (format: /courses/:courseID/tasks/:taskID)
	path := r.URL.Path
	// Remove prefix "/courses/"
	remaining := path[len("/courses/"):]

	// Find the /tasks/ separator
	tasksIdx := -1
	for i := 0; i < len(remaining)-6; i++ {
		if remaining[i:i+7] == "/tasks/" {
			tasksIdx = i
			break
		}
	}

	if tasksIdx == -1 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	courseID := remaining[:tasksIdx]
	taskID := remaining[tasksIdx+7:] // +7 for "/tasks/"

	if courseID == "" || taskID == "" {
		http.Error(w, "Course ID and Task ID are required", http.StatusBadRequest)
		return
	}

	// Load the course
	loader := courseparser.NewCourseLoader()
	coursePath := filepath.Join(CoursesDir, courseID)
	course, err := loader.LoadCourse(coursePath)
	if err != nil {
		http.Error(w, "Course not found", http.StatusNotFound)
		log.Printf("Error loading course %s: %v", courseID, err)
		return
	}

	// Find the task
	task, ok := course.Tasks[taskID]
	if !ok {
		http.Error(w, "Task not found", http.StatusNotFound)
		log.Printf("Task %s not found in course %s", taskID, courseID)
		return
	}

	// Build problem responses with full details (preserving order)
	problems := make([]ProblemDetailResponse, 0, task.Problems.Len())
	for _, op := range task.Problems.Problems {
		problemResp := ProblemDetailResponse{
			ProblemResponse: ProblemResponse{
				ID:     op.ID,
				Type:   op.Problem.GetType(),
				Name:   op.Problem.GetName(),
				Header: op.Problem.GetHeader(),
			},
		}

		// Add type-specific fields
		switch p := op.Problem.(type) {
		case *courseparser.CodeProblem:
			problemResp.Language = p.Language
			problemResp.Default = p.Default
		case *courseparser.MultipleChoiceProblem:
			// Only expose choice text, not the correct answer (Valid field)
			problemResp.Choices = make([]Choice, len(p.Choices))
			for i, c := range p.Choices {
				problemResp.Choices[i] = Choice{
					Text: c.Text,
				}
			}
			problemResp.Limit = p.Limit
		case *courseparser.MatchProblem:
			// Do not expose the answer for match problems
		}

		problems = append(problems, problemResp)
	}

	// Build response
	response := TaskDetailResponse{
		ID:              taskID,
		CourseID:        courseID,
		Name:            task.Name,
		Author:          task.Author,
		ContactURL:      task.ContactURL,
		Context:         task.Context,
		EnvironmentID:   task.EnvironmentID,
		EnvironmentType: task.EnvironmentType,
		NetworkGrading:  task.NetworkGrading,
		Problems:        problems,
	}

	// Add limits if present
	if task.EnvironmentParameters.Limits != nil {
		response.Limits = &EnvironmentLimits{
			Time:     task.EnvironmentParameters.Limits.Time,
			HardTime: task.EnvironmentParameters.Limits.HardTime,
			Memory:   task.EnvironmentParameters.Limits.Memory,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}
