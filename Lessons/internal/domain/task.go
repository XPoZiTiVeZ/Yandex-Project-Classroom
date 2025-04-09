package domain

import "time"

type Task struct {
	ID        string    // Уникальный идентификатор задания
	CourseID  string    // Идентификатор курса, к которому относится задание
	Title     string    // Название задания
	Content   string    // Содержание задания
	CreatedAt time.Time // Дата создания задания
}
