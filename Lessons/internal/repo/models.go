package repo

import (
	"Classroom/Lessons/internal/domain"
	"time"
)

type Lesson struct {
	ID          string    `db:"lesson_id"`
	CourseID    string    `db:"course_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}

func (l Lesson) ToEntity() domain.Lesson {
	return domain.Lesson{
		ID:          l.ID,
		CourseID:    l.CourseID,
		Title:       l.Title,
		Description: l.Description,
		CreatedAt:   l.CreatedAt,
	}
}
