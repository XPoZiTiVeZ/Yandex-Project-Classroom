package domain

import "time"

type Task struct {
	ID        string    // Уникальный идентификатор задания
	CourseID  string    // Идентификатор курса, к которому относится задание
	Title     string    // Название задания
	Content   string    // Содержание задания
	CreatedAt time.Time // Дата создания задания
}

type StudentTask struct {
	ID        string    // Уникальный идентификатор задания
	CourseID  string    // Идентификатор курса, к которому относится задание
	Title     string    // Название задания
	Content   string    // Содержание задания
	Completed bool      // Флаг, указывающий на выполненность задания
	CreatedAt time.Time // Дата создания задания
}

type TaskStatus struct {
	UserID      string // Идентификатор пользователя
	TaskID      string // Идентификатор задания
	IsCompleted bool   // Флаг, указывающий на выполненность задания
}
