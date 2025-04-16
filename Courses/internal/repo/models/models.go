package repository

import "time"

type Course struct {
	CourseID    string
	TeacherID   string
	Title       string
	Description string
	Visibility  bool
	StartTime   time.Time
	EndTime     time.Time
}

type UserCourse struct {
	UserID   string
	CourseID string
}

type Member struct {
    UserID    string
    Email     string
    FirstName string
    LastName  string
}