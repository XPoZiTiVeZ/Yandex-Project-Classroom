package repo

import (
	"Classroom/Lessons/internal/domain"
	"time"
)

type Task struct {
	ID        string    `db:"task_id"`
	CourseID  string    `db:"course_id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}

func (t Task) ToEntity() domain.Task {
	return domain.Task{
		ID:        t.ID,
		CourseID:  t.CourseID,
		Title:     t.Title,
		Content:   t.Content,
		CreatedAt: t.CreatedAt,
	}
}

type TaskStatus struct {
	UserID      string `db:"student_id"`
	TaskID      string `db:"task_id"`
	IsCompleted bool   `db:"completed"`
}

func (t TaskStatus) ToEntity() domain.TaskStatus {
	return domain.TaskStatus{
		UserID:      t.UserID,
		TaskID:      t.TaskID,
		IsCompleted: t.IsCompleted,
	}
}

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
