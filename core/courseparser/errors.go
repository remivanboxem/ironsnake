package courseparser

import "fmt"

// ParseError represents a parsing error with context
type ParseError struct {
	File    string // File path where error occurred
	Field   string // Field name (optional)
	Message string // Error description
	Err     error  // Underlying error
}

func (e *ParseError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: field %q: %s", e.File, e.Field, e.Message)
	}
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.File, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.File, e.Message)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

// CourseLoadError represents an error when loading a course
type CourseLoadError struct {
	CourseID string
	Message  string
	Err      error
}

func (e *CourseLoadError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("course %s: %s: %v", e.CourseID, e.Message, e.Err)
	}
	return fmt.Sprintf("course %s: %s", e.CourseID, e.Message)
}

func (e *CourseLoadError) Unwrap() error {
	return e.Err
}
