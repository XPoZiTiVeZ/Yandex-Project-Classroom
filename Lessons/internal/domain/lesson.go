package domain

import "time"

// Lesson представляет доменную модель урока
type Lesson struct {
	ID          string    // Уникальный идентификатор урока
	CourseID    string    // Идентификатор курса, к которому относится урок
	Title       string    // Название урока
	Description string    // Описание урока
	CreatedAt   time.Time // Время создания урока
}
