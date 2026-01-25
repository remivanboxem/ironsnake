package courseparser

import (
	"os"
	"path/filepath"
)

// CourseLoader loads and parses course files from the filesystem
type CourseLoader struct{}

// NewCourseLoader creates a new CourseLoader
func NewCourseLoader() *CourseLoader {
	return &CourseLoader{}
}

// LoadCourse loads a complete course from a directory path
func (l *CourseLoader) LoadCourse(dirPath string) (*ParsedCourse, error) {
	courseID := filepath.Base(dirPath)

	course := &ParsedCourse{
		DirPath:  dirPath,
		CourseID: courseID,
		Tasks:    make(map[string]TaskConfig),
	}

	// Load config.yaml
	configPath := filepath.Join(dirPath, "config.yaml")
	config, err := ParseCourseConfig(configPath)
	if err != nil {
		return nil, &CourseLoadError{
			CourseID: courseID,
			Message:  "failed to load config.yaml",
			Err:      err,
		}
	}
	course.Config = *config

	// Load access.yaml
	accessPath := filepath.Join(dirPath, "access.yaml")
	access, err := ParseAccessConfig(accessPath)
	if err != nil {
		return nil, &CourseLoadError{
			CourseID: courseID,
			Message:  "failed to load access.yaml",
			Err:      err,
		}
	}
	course.Access = *access

	// Load tasks
	tasksDir := filepath.Join(dirPath, "tasks")
	if err := l.loadTasks(tasksDir, course); err != nil {
		return nil, &CourseLoadError{
			CourseID: courseID,
			Message:  "failed to load tasks",
			Err:      err,
		}
	}

	// Load syllabus (optional)
	syllabusDir := filepath.Join(dirPath, "syllabus")
	if _, err := os.Stat(syllabusDir); err == nil {
		syllabus, err := ParseSyllabus(syllabusDir)
		if err != nil {
			return nil, &CourseLoadError{
				CourseID: courseID,
				Message:  "failed to load syllabus",
				Err:      err,
			}
		}
		course.Syllabus = syllabus
	}

	return course, nil
}

// loadTasks loads all tasks from the tasks directory
func (l *CourseLoader) loadTasks(tasksDir string, course *ParsedCourse) error {
	entries, err := os.ReadDir(tasksDir)
	if err != nil {
		return &ParseError{
			File:    tasksDir,
			Message: "failed to read tasks directory",
			Err:     err,
		}
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		taskID := entry.Name()
		taskPath := filepath.Join(tasksDir, taskID, "task.yaml")

		// Check if task.yaml exists
		if _, err := os.Stat(taskPath); os.IsNotExist(err) {
			continue // Skip directories without task.yaml
		}

		task, err := ParseTaskConfig(taskPath)
		if err != nil {
			return &ParseError{
				File:    taskPath,
				Message: "failed to parse task",
				Err:     err,
			}
		}

		course.Tasks[taskID] = *task
	}

	return nil
}

// LoadAllCourses loads all courses from a courses directory
func (l *CourseLoader) LoadAllCourses(coursesDir string) ([]*ParsedCourse, error) {
	entries, err := os.ReadDir(coursesDir)
	if err != nil {
		return nil, &ParseError{
			File:    coursesDir,
			Message: "failed to read courses directory",
			Err:     err,
		}
	}

	var courses []*ParsedCourse
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		coursePath := filepath.Join(coursesDir, entry.Name())

		// Check if this looks like a course directory (has config.yaml)
		configPath := filepath.Join(coursePath, "config.yaml")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			continue
		}

		course, err := l.LoadCourse(coursePath)
		if err != nil {
			return nil, err
		}

		courses = append(courses, course)
	}

	return courses, nil
}
