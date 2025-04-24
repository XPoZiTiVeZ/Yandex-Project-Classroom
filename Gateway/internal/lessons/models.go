package lessons

import (
	"time"

	pb "Classroom/Gateway/pkg/api/lessons"
)

// Lesson - информация о занятии
// @Description Полная информация о занятии в курсе
type Lesson struct {
    // Уникальный идентификатор занятия
    LessonID string `json:"lesson_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
    // Идентификатор курса
    CourseID string `json:"course_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=1"`
    // Название занятия
    Title string `json:"title" example:"Введение в программирование" extensions:"x-order=2"`
    // Содержание занятия
    Description string `json:"description" example:"Базовые понятия и термины" extensions:"x-order=3"`
    // Дата создания
    CreatedAt time.Time `json:"created_at" example:"2023-01-15T10:00:00Z" extensions:"x-order=4"`
} // @name Lesson

// CreateLessonRequest - запрос на создание занятия
// @Description Параметры для создания нового занятия в курсе
type CreateLessonRequest struct {
    // ID курса
    CourseID string `json:"course_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
    // Название занятия
    Title string `json:"title" example:"Основы алгоритмов" extensions:"x-order=1"`
    // Содержание занятия
    Content string `json:"content" example:"Подробное описание занятия..." extensions:"x-order=2"`
} // @name CreateLessonRequest

func NewCreateLessonRequest(req CreateLessonRequest) *pb.CreateLessonRequest {
	return &pb.CreateLessonRequest{
		CourseId: req.CourseID,
		Title:    req.Title,
		Content:  req.Content,
	}
}

// CreateLessonResponse - ответ после создания занятия
// @Description Возвращает ID созданного занятия
type CreateLessonResponse struct {
    // ID созданного занятия
    LessonID string `json:"lesson_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
} // @name CreateLessonResponse

func NewCreateLessonResponse(resp *pb.CreateLessonResponse) CreateLessonResponse {
	return CreateLessonResponse{
		LessonID: resp.GetLessonId(),
	}
}

// GetLessonRequest - запрос информации о занятии
// @Description Требует ID курса и ID занятия для получения данных
type GetLessonRequest struct {
    // ID занятия
    LessonID string `json:"lesson_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
} // @name GetLessonRequest

func NewGetLessonRequest(req GetLessonRequest) *pb.GetLessonRequest {
	return &pb.GetLessonRequest{
		LessonId: req.LessonID,
	}
}

// GetLessonResponse - информация о занятии
// @Description Возвращает полные данные занятия
type GetLessonResponse struct {
    // Объект занятия
    Lesson Lesson `json:"lesson" extensions:"x-order=0"`
} // @name GetLessonResponse

func NewGetLessonResponse(resp *pb.GetLessonResponse) GetLessonResponse {
	return GetLessonResponse{
		Lesson: Lesson{
			LessonID:    resp.GetLesson().GetLessonId(),
			CourseID:    resp.GetLesson().GetCourseId(),
			Title:       resp.GetLesson().GetTitle(),
			Description: resp.GetLesson().GetContent(),
			CreatedAt:   resp.GetLesson().GetCreatedAt().AsTime(),
		},
	}
}

// GetLessonsRequest - запрос списка занятий
// @Description Возвращает все занятия в указанном курсе
type GetLessonsRequest struct {
    // ID курса
    CourseID string `json:"course_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
} // @name GetLessonsRequest

func NewGetLessonsRequest(req GetLessonsRequest) *pb.GetLessonsRequest {
	return &pb.GetLessonsRequest{
		CourseId: req.CourseID,
	}
}

// GetLessonsResponse - список занятий
// @Description Содержит массив занятий в курсе
type GetLessonsResponse struct {
    // Массив занятий
    Lessons []Lesson `json:"lessons" extensions:"x-order=0"`
} // @name GetLessonsResponse

func NewGetLessonsResponse(resp *pb.GetLessonsResponse) GetLessonsResponse {
	return GetLessonsResponse{
		Lessons: func() []Lesson {
			lessons := make([]Lesson, 0, len(resp.GetLessons()))
			for _, lesson := range resp.GetLessons() {
				lessons = append(lessons, Lesson{
					LessonID:    lesson.GetLessonId(),
					CourseID:    lesson.GetCourseId(),
					Title:       lesson.GetTitle(),
					Description: lesson.GetContent(),
					CreatedAt:   lesson.GetCreatedAt().AsTime(),
				})
			}
			return lessons
		}(),
	}
}

// UpdateLessonRequest - запрос обновления занятия
// @Description Позволяет частично обновить данные занятия
type UpdateLessonRequest struct {
    // ID занятия
    LessonID string `json:"lesson_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
    // Новое название (опционально)
    Title *string `json:"title,omitempty" example:"Обновленное название" extensions:"x-order=1"`
    // Новое содержание (опционально)
    Content *string `json:"description,omitempty" example:"Обновленное содержание" extensions:"x-order=2"`
} // @name UpdateLessonRequest

func NewUpdateLessonRequest(req UpdateLessonRequest) *pb.UpdateLessonRequest {
	return &pb.UpdateLessonRequest{
		LessonId: req.LessonID,
		Title:    req.Title,
		Content:  req.Content,
	}
}

// UpdateLessonResponse - результат обновления
// @Description Пустой ответ при успешном обновлении
type UpdateLessonResponse struct {
	Lesson Lesson `json:"lesson" extensions:"x-order=0"`
} // @name UpdateLessonResponse

func NewUpdateLessonResponse(resp *pb.UpdateLessonResponse) UpdateLessonResponse {
	return UpdateLessonResponse{}
}

// DeleteLessonRequest - запрос удаления занятия
// @Description Требует ID курса и ID занятия для удаления
type DeleteLessonRequest struct {
    // ID занятия
    LessonID string `json:"lesson_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
} // @name DeleteLessonRequest

func NewDeleteLessonRequest(req DeleteLessonRequest) *pb.DeleteLessonRequest {
	return &pb.DeleteLessonRequest{
		LessonId: req.LessonID,
	}
}

// DeleteLessonResponse - результат удаления
// @Description Пустой ответ при успешном удалении
type DeleteLessonResponse struct {
} // @name DeleteLessonResponse

func NewDeleteLessonResponse(resp *pb.DeleteLessonResponse) DeleteLessonResponse {
	return DeleteLessonResponse{}
}
