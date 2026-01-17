package main

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Username     string    `gorm:"type:varchar(255);not null;unique"`
	Email        string    `gorm:"type:varchar(255);not null;unique"`
	PasswordHash string    `gorm:"type:varchar(255);not null"`
	FirstName    string    `gorm:"type:varchar(255)"`
	LastName     string    `gorm:"type:varchar(255)"`
	CreatedAt    time.Time `gorm:"type:timestamp;default:now()"`
	Tasks        []Task    `gorm:"foreignKey:UserID"`
}

type Task struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	User      User      `gorm:"foreignKey:UserID"`
	Name      string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"type:timestamp;default:now()"`
}

type Course struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Code         string    `gorm:"type:varchar(255);not null"`
	Name         string    `gorm:"type:varchar(255);not null"`
	Description  string    `gorm:"type:text"`
	AcademicYear string    `gorm:"type:varchar(255);not null"`
	CreatedBy    uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt    time.Time `gorm:"type:timestamp;default:now()"`
}

type CourseTeacher struct {
	CourseID  uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
	TeacherID uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
	RoleID    uuid.UUID `gorm:"type:uuid;not null"`
	Course    Course    `gorm:"foreignKey:CourseID"`
	Teacher   User      `gorm:"foreignKey:TeacherID"`
	Role      Role      `gorm:"foreignKey:RoleID"`
}

type Role struct {
	ID   uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name string    `gorm:"type:varchar(255);not null;unique"`
}
