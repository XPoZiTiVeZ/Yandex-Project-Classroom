package repo

import (
	"Classroom/Lessons/internal/domain"
	"time"
)

type Lesson struct {
	ID        string    `db:"lesson_id"`
	CourseID  string    `db:"course_id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}

func (l Lesson) ToEntity() domain.Lesson {
	return domain.Lesson{
		ID:        l.ID,
		CourseID:  l.CourseID,
		Title:     l.Title,
		Content:   l.Content,
		CreatedAt: l.CreatedAt,
	}
}
