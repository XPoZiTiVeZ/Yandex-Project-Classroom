package domain

import "time"

type Enrollment struct {
	StudentID  string
	CourseID   string
	EnrolledAt time.Time
}
