package domain

import "time"

// Lesson представляет доменную модель урока
type Lesson struct {
	ID        string    // Уникальный идентификатор урока
	CourseID  string    // Идентификатор курса, к которому относится урок
	Title     string    // Название урока
	Content   string    // Содержание урока
	CreatedAt time.Time // Время создания урока
}
