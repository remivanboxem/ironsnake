package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// CourseResponse represents the JSON response structure for a course
type CourseResponse struct {
	ID           uuid.UUID `json:"id"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	AcademicYear string    `json:"academicYear"`
	CreatedBy    uuid.UUID `json:"createdBy"`
	CreatedAt    string    `json:"createdAt"`
}

func getCoursesHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var courses []Course
	if err := DB.Find(&courses).Error; err != nil {
		http.Error(w, "Failed to fetch courses", http.StatusInternalServerError)
		log.Printf("Error fetching courses: %v", err)
		return
	}

	// Map to response format
	response := make([]CourseResponse, len(courses))
	for i, course := range courses {
		response[i] = CourseResponse{
			ID:           course.ID,
			Code:         course.Code,
			Name:         course.Name,
			Description:  course.Description,
			AcademicYear: course.AcademicYear,
			CreatedBy:    course.CreatedBy,
			CreatedAt:    course.CreatedAt.Format(time.RFC3339),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
		return
	}
}
