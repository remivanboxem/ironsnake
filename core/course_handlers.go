package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"

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
		problems := make([]ProblemResponse, 0, len(task.Problems))
		for problemID, problem := range task.Problems {
			problems = append(problems, ProblemResponse{
				ID:     problemID,
				Type:   problem.GetType(),
				Name:   problem.GetName(),
				Header: problem.GetHeader(),
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

	// Build problem responses with full details
	problems := make([]ProblemDetailResponse, 0, len(task.Problems))
	for problemID, problem := range task.Problems {
		problemResp := ProblemDetailResponse{
			ProblemResponse: ProblemResponse{
				ID:     problemID,
				Type:   problem.GetType(),
				Name:   problem.GetName(),
				Header: problem.GetHeader(),
			},
		}

		// Add type-specific fields
		switch p := problem.(type) {
		case *courseparser.CodeProblem:
			problemResp.Language = p.Language
			problemResp.Default = p.Default
		case *courseparser.MultipleChoiceProblem:
			problemResp.Choices = make([]Choice, len(p.Choices))
			for i, c := range p.Choices {
				problemResp.Choices[i] = Choice{
					Text:  c.Text,
					Valid: c.Valid,
				}
			}
			problemResp.Limit = p.Limit
		case *courseparser.MatchProblem:
			problemResp.Answer = p.Answer
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
