package events

// Сообщение о том, что пользователь был зачислен на курс
type UserEnrolled struct {
	CourseID string `json:"course_id"`
	UserID   string `json:"user_id"`
}

// Сообщение о том, что пользователь был отчислен с курса
type UserExpelled struct {
	CourseID string `json:"course_id"`
	UserID   string `json:"user_id"`
}
