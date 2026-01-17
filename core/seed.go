package main

import (
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// SeedDatabase seeds the database with fake data if it's empty
func SeedDatabase() error {
	// Check if data already exists
	var userCount int64
	DB.Model(&User{}).Count(&userCount)
	if userCount > 0 {
		log.Println("Database already contains data, skipping seed")
		return nil
	}

	log.Println("Seeding database with fake data...")

	// Create roles
	roles := []Role{
		{ID: uuid.New(), Name: "Professor"},
		{ID: uuid.New(), Name: "Assistant Professor"},
		{ID: uuid.New(), Name: "Teaching Assistant"},
		{ID: uuid.New(), Name: "Lecturer"},
	}
	for _, role := range roles {
		if err := DB.Create(&role).Error; err != nil {
			return err
		}
	}
	log.Printf("Created %d roles", len(roles))

	// Create users
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	users := []User{
		{
			ID:           uuid.New(),
			Username:     "jdoe",
			Email:        "john.doe@university.edu",
			PasswordHash: string(passwordHash),
			FirstName:    "John",
			LastName:     "Doe",
			CreatedAt:    time.Now().AddDate(0, -6, 0),
		},
		{
			ID:           uuid.New(),
			Username:     "asmith",
			Email:        "alice.smith@university.edu",
			PasswordHash: string(passwordHash),
			FirstName:    "Alice",
			LastName:     "Smith",
			CreatedAt:    time.Now().AddDate(0, -5, 0),
		},
		{
			ID:           uuid.New(),
			Username:     "bjones",
			Email:        "bob.jones@university.edu",
			PasswordHash: string(passwordHash),
			FirstName:    "Bob",
			LastName:     "Jones",
			CreatedAt:    time.Now().AddDate(0, -4, 0),
		},
		{
			ID:           uuid.New(),
			Username:     "cwilliams",
			Email:        "carol.williams@university.edu",
			PasswordHash: string(passwordHash),
			FirstName:    "Carol",
			LastName:     "Williams",
			CreatedAt:    time.Now().AddDate(0, -3, 0),
		},
		{
			ID:           uuid.New(),
			Username:     "dbrown",
			Email:        "david.brown@university.edu",
			PasswordHash: string(passwordHash),
			FirstName:    "David",
			LastName:     "Brown",
			CreatedAt:    time.Now().AddDate(0, -2, 0),
		},
	}
	for _, user := range users {
		if err := DB.Create(&user).Error; err != nil {
			return err
		}
	}
	log.Printf("Created %d users", len(users))

	// Create courses
	courses := []Course{
		{
			ID:           uuid.New(),
			Code:         "CS101",
			Name:         "Introduction to Computer Science",
			Description:  "Fundamental concepts of computer science including algorithms, data structures, and programming.",
			AcademicYear: "2025-2026",
			CreatedBy:    users[0].ID,
			CreatedAt:    time.Now().AddDate(0, -2, 0),
		},
		{
			ID:           uuid.New(),
			Code:         "CS201",
			Name:         "Data Structures and Algorithms",
			Description:  "Advanced study of data structures, algorithm design, and complexity analysis.",
			AcademicYear: "2025-2026",
			CreatedBy:    users[1].ID,
			CreatedAt:    time.Now().AddDate(0, -2, 0),
		},
		{
			ID:           uuid.New(),
			Code:         "CS301",
			Name:         "Database Systems",
			Description:  "Design and implementation of database systems, SQL, and data modeling.",
			AcademicYear: "2025-2026",
			CreatedBy:    users[0].ID,
			CreatedAt:    time.Now().AddDate(0, -1, -15),
		},
		{
			ID:           uuid.New(),
			Code:         "CS401",
			Name:         "Machine Learning",
			Description:  "Introduction to machine learning algorithms, neural networks, and AI applications.",
			AcademicYear: "2025-2026",
			CreatedBy:    users[2].ID,
			CreatedAt:    time.Now().AddDate(0, -1, -10),
		},
		{
			ID:           uuid.New(),
			Code:         "MATH201",
			Name:         "Linear Algebra",
			Description:  "Vector spaces, matrices, linear transformations, and eigenvalues.",
			AcademicYear: "2025-2026",
			CreatedBy:    users[3].ID,
			CreatedAt:    time.Now().AddDate(0, -1, -5),
		},
	}
	for _, course := range courses {
		if err := DB.Create(&course).Error; err != nil {
			return err
		}
	}
	log.Printf("Created %d courses", len(courses))

	// Create course-teacher relationships
	courseTeachers := []CourseTeacher{
		{CourseID: courses[0].ID, TeacherID: users[0].ID, RoleID: roles[0].ID}, // CS101 - John Doe - Professor
		{CourseID: courses[0].ID, TeacherID: users[4].ID, RoleID: roles[2].ID}, // CS101 - David Brown - TA
		{CourseID: courses[1].ID, TeacherID: users[1].ID, RoleID: roles[0].ID}, // CS201 - Alice Smith - Professor
		{CourseID: courses[2].ID, TeacherID: users[0].ID, RoleID: roles[0].ID}, // CS301 - John Doe - Professor
		{CourseID: courses[3].ID, TeacherID: users[2].ID, RoleID: roles[1].ID}, // CS401 - Bob Jones - Assistant Professor
		{CourseID: courses[3].ID, TeacherID: users[1].ID, RoleID: roles[3].ID}, // CS401 - Alice Smith - Lecturer
		{CourseID: courses[4].ID, TeacherID: users[3].ID, RoleID: roles[0].ID}, // MATH201 - Carol Williams - Professor
	}
	for _, ct := range courseTeachers {
		if err := DB.Create(&ct).Error; err != nil {
			return err
		}
	}
	log.Printf("Created %d course-teacher relationships", len(courseTeachers))

	// Create tasks for users
	tasks := []Task{
		{
			ID:        uuid.New(),
			UserID:    users[0].ID,
			Name:      "Prepare lecture slides for CS101",
			CreatedAt: time.Now().AddDate(0, 0, -5),
		},
		{
			ID:        uuid.New(),
			UserID:    users[0].ID,
			Name:      "Grade midterm exams",
			CreatedAt: time.Now().AddDate(0, 0, -3),
		},
		{
			ID:        uuid.New(),
			UserID:    users[1].ID,
			Name:      "Update course syllabus",
			CreatedAt: time.Now().AddDate(0, 0, -4),
		},
		{
			ID:        uuid.New(),
			UserID:    users[1].ID,
			Name:      "Review homework submissions",
			CreatedAt: time.Now().AddDate(0, 0, -2),
		},
		{
			ID:        uuid.New(),
			UserID:    users[2].ID,
			Name:      "Prepare ML project guidelines",
			CreatedAt: time.Now().AddDate(0, 0, -1),
		},
		{
			ID:        uuid.New(),
			UserID:    users[3].ID,
			Name:      "Schedule office hours",
			CreatedAt: time.Now().AddDate(0, 0, -6),
		},
		{
			ID:        uuid.New(),
			UserID:    users[4].ID,
			Name:      "Hold TA session for students",
			CreatedAt: time.Now().AddDate(0, 0, -1),
		},
		{
			ID:        uuid.New(),
			UserID:    users[4].ID,
			Name:      "Prepare quiz questions",
			CreatedAt: time.Now(),
		},
	}
	for _, task := range tasks {
		if err := DB.Create(&task).Error; err != nil {
			return err
		}
	}
	log.Printf("Created %d tasks", len(tasks))

	log.Println("Database seeding completed successfully!")
	return nil
}
