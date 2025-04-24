package events

// Сообщение о том, что на курсе было добавлено новое дз
type TaskCreated struct {
	CourseID string `json:"course_id"`
	TaskID   string `json:"task_id"`
}
