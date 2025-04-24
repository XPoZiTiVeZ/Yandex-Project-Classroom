package domain

import "time"

type Course struct {
	ID          string
	TeacherID   string
	Title       string
	Description string
	Visibility  bool
	StartTime   *time.Time
	EndTime     *time.Time
	CreatedAt   time.Time
}
