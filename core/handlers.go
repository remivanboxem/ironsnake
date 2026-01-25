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

// CourseDetailResponse represents the full course detail
type CourseDetailResponse struct {
	CourseResponse
	Tasks    []TaskResponse `json:"tasks"`
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
