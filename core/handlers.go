package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// AuthorResponse represents the author/creator information
type AuthorResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
}

// CourseResponse represents the JSON response structure for a course
type CourseResponse struct {
	ID           uuid.UUID       `json:"id"`
	Code         string          `json:"code"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	AcademicYear string          `json:"academicYear"`
	CreatedBy    uuid.UUID       `json:"createdBy"`
	Author       AuthorResponse  `json:"author"`
	CreatedAt    string          `json:"createdAt"`
}

func getCoursesHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Query courses with author information
	type CourseWithAuthor struct {
		Course
		AuthorID        uuid.UUID
		AuthorUsername  string
		AuthorFirstName string
		AuthorLastName  string
	}

	var coursesWithAuthors []CourseWithAuthor
	if err := DB.Table("courses").
		Select("courses.*, users.id as author_id, users.username as author_username, users.first_name as author_first_name, users.last_name as author_last_name").
		Joins("LEFT JOIN users ON users.id = courses.created_by").
		Find(&coursesWithAuthors).Error; err != nil {
		http.Error(w, "Failed to fetch courses", http.StatusInternalServerError)
		log.Printf("Error fetching courses: %v", err)
		return
	}

	// Map to response format
	response := make([]CourseResponse, len(coursesWithAuthors))
	for i, cwa := range coursesWithAuthors {
		response[i] = CourseResponse{
			ID:           cwa.ID,
			Code:         cwa.Code,
			Name:         cwa.Name,
			Description:  cwa.Description,
			AcademicYear: cwa.AcademicYear,
			CreatedBy:    cwa.CreatedBy,
			Author: AuthorResponse{
				ID:        cwa.AuthorID,
				Username:  cwa.AuthorUsername,
				FirstName: cwa.AuthorFirstName,
				LastName:  cwa.AuthorLastName,
			},
			CreatedAt: cwa.CreatedAt.Format(time.RFC3339),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
		return
	}
}
