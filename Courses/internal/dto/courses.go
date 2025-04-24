package dto

import "time"

type CreateCourseDTO struct {
	TeacherID   string `validate:"required,uuid"`
	Title       string `validate:"required"`
	Description string `validate:"required"`
	Visibility  bool
	StartTime   *time.Time
	EndTime     *time.Time
}

type UpdateCourseDTO struct {
	ID          string `validate:"required,uuid"`
	Title       *string
	Description *string
	Visibility  *bool
	StartTime   *time.Time
	EndTime     *time.Time
}
