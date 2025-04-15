package dto

type CreateLessonDTO struct {
	Title    string `validate:"required"`
	Content  string `validate:"required"`
	CourseID string `validate:"required,uuid"`
}

type UpdateLessonDTO struct {
	LessonID string `validate:"required,uuid"`
	Title    *string
	Content  *string
}
