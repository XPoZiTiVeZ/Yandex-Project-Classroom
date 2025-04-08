package dto

type CreateTaskDTO struct {
	Title    string `validate:"required"`
	Content  string `validate:"required"`
	CourseID string `validate:"required,uuid"`
}
