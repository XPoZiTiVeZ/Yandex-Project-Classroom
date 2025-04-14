package repo

import (
	"Classroom/Tasks/internal/domain"
	"time"
)

type StudentTask struct {
	ID        string    `db:"task_id"`
	CourseID  string    `db:"course_id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	Completed bool      `db:"completed"`
	CreatedAt time.Time `db:"created_at"`
}

func (t StudentTask) ToEntity() domain.StudentTask {
	return domain.StudentTask{
		ID:        t.ID,
		CourseID:  t.CourseID,
		Title:     t.Title,
		Content:   t.Content,
		Completed: t.Completed,
		CreatedAt: t.CreatedAt,
	}
}

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
	UserID    string `db:"student_id"`
	TaskID    string `db:"task_id"`
	Completed bool   `db:"completed"`
}

func (t TaskStatus) ToEntity() domain.TaskStatus {
	return domain.TaskStatus{
		UserID:    t.UserID,
		TaskID:    t.TaskID,
		Completed: t.Completed,
	}
}
