package dto

type CreateTaskDTO struct {
	Title    string `validate:"required"`
	Content  string `validate:"required"`
	CourseID string `validate:"required,uuid"`
}

type UpdateTaskDTO struct {
	TaskID  string `validate:"required,uuid"`
	Title   *string
	Content *string
}
