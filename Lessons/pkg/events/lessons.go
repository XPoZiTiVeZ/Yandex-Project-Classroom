package events

// Сообщение о том, что на курсе был добавлен новый урок
type LessonCreated struct {
	CourseID string `json:"course_id"`
	LessonID string `json:"lesson_id"`
}

const LessonCreatedTopic = "lesson.created"
