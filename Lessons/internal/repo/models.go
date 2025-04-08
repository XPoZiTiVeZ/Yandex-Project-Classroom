package repo

import "Classroom/Lessons/internal/domain"

type Task struct {
	ID       string `db:"task_id"`
	CourseID string `db:"course_id"`
	Title    string `db:"title"`
	Content  string `db:"content"`
}

func (t Task) ToEntity() domain.Task {
	return domain.Task{
		ID:       t.ID,
		CourseID: t.CourseID,
		Title:    t.Title,
		Content:  t.Content,
	}
}
