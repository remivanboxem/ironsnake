package courseparser

import (
	"testing"
)

func TestLoadCourse(t *testing.T) {
	loader := NewCourseLoader()
	course, err := loader.LoadCourse("../../courses/CS01")
	if err != nil {
		t.Fatalf("failed to load course: %v", err)
	}

	// Verify course ID
	if course.CourseID != "CS01" {
		t.Errorf("expected CourseID 'CS01', got %q", course.CourseID)
	}

	// Verify config
	if course.Config.Name != "[CS01] Introduction to Computer Science" {
		t.Errorf("unexpected course name: %q", course.Config.Name)
	}
	if len(course.Config.Admins) != 2 {
		t.Errorf("expected 2 admins, got %d", len(course.Config.Admins))
	}

	// Verify tasks loaded
	if len(course.Tasks) == 0 {
		t.Error("no tasks loaded")
	}

	// Verify task01 (code problem)
	task01, ok := course.Tasks["task01"]
	if !ok {
		t.Error("task01 not found")
	} else {
		if task01.Name != "Convertisseur binaire vers Base64" {
			t.Errorf("unexpected task01 name: %q", task01.Name)
		}
		if !task01.IsDocker() {
			t.Error("task01 should be a docker task")
		}
		if len(task01.Problems) != 1 {
			t.Errorf("expected 1 problem in task01, got %d", len(task01.Problems))
		}

		// Verify code problem type
		problem, ok := task01.Problems["binary_to_base64"]
		if !ok {
			t.Error("problem binary_to_base64 not found")
		} else {
			if problem.GetType() != "code" {
				t.Errorf("expected type 'code', got %q", problem.GetType())
			}
			codeProblem, ok := problem.(*CodeProblem)
			if !ok {
				t.Error("problem should be a CodeProblem")
			} else {
				if codeProblem.Language != "python" {
					t.Errorf("expected language 'python', got %q", codeProblem.Language)
				}
			}
		}
	}

	// Verify task05 (MCQ with multiple problem types)
	task05, ok := course.Tasks["task05"]
	if !ok {
		t.Error("task05 not found")
	} else {
		if !task05.IsMCQ() {
			t.Error("task05 should be an MCQ task")
		}
		if len(task05.Problems) != 5 {
			t.Errorf("expected 5 problems in task05, got %d", len(task05.Problems))
		}

		// Verify multiple_choice problem
		q1, ok := task05.Problems["Q1"]
		if !ok {
			t.Error("Q1 not found")
		} else {
			if q1.GetType() != "multiple_choice" {
				t.Errorf("expected type 'multiple_choice', got %q", q1.GetType())
			}
			mcProblem, ok := q1.(*MultipleChoiceProblem)
			if !ok {
				t.Error("Q1 should be a MultipleChoiceProblem")
			} else {
				if len(mcProblem.Choices) != 6 {
					t.Errorf("expected 6 choices in Q1, got %d", len(mcProblem.Choices))
				}
			}
		}

		// Verify match problem
		q4, ok := task05.Problems["Q4"]
		if !ok {
			t.Error("Q4 not found")
		} else {
			if q4.GetType() != "match" {
				t.Errorf("expected type 'match', got %q", q4.GetType())
			}
			matchProblem, ok := q4.(*MatchProblem)
			if !ok {
				t.Error("Q4 should be a MatchProblem")
			} else {
				if matchProblem.Answer != "O(1)" {
					t.Errorf("expected answer 'O(1)', got %q", matchProblem.Answer)
				}
			}
		}
	}

	// Verify access config
	task01Access, ok := course.Access.DispenserData.Config["task01"]
	if !ok {
		t.Error("task01 access config not found")
	} else {
		if task01Access.Accessibility.IsBoolean {
			t.Error("task01 accessibility should be a date range, not boolean")
		}
		if task01Access.Accessibility.DateRange == nil {
			t.Error("task01 accessibility date range is nil")
		}
		if task01Access.EvaluationMode != "best" {
			t.Errorf("expected evaluation_mode 'best', got %q", task01Access.EvaluationMode)
		}
	}

	task02Access, ok := course.Access.DispenserData.Config["task02"]
	if !ok {
		t.Error("task02 access config not found")
	} else {
		if !task02Access.Accessibility.IsBoolean {
			t.Error("task02 accessibility should be boolean")
		}
		if !task02Access.Accessibility.BoolValue {
			t.Error("task02 accessibility should be true")
		}
	}

	// Verify syllabus
	if course.Syllabus == nil {
		t.Error("syllabus is nil")
	} else {
		if course.Syllabus.Book.Book.Title != "Introduction à la programmation" {
			t.Errorf("unexpected book title: %q", course.Syllabus.Book.Book.Title)
		}
		if len(course.Syllabus.Summary) == 0 {
			t.Error("summary is empty")
		}
	}
}

func TestLoadAllCourses(t *testing.T) {
	loader := NewCourseLoader()
	courses, err := loader.LoadAllCourses("../../courses")
	if err != nil {
		t.Fatalf("failed to load courses: %v", err)
	}

	if len(courses) != 1 {
		t.Errorf("expected 1 course, got %d", len(courses))
	}

	if len(courses) > 0 && courses[0].CourseID != "CS01" {
		t.Errorf("expected course ID 'CS01', got %q", courses[0].CourseID)
	}
}
