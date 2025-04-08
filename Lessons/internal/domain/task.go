package domain

import "time"

type Task struct {
	ID          string    // Уникальный идентификатор задания
	CourseID    string    // Идентификатор курса, к которому относится задание
	Title       string    // Название задания
	Description string    // Описание задания
	CreatedAt   time.Time // Время создания задания
}
