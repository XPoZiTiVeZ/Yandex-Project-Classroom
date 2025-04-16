package lessons

import (
	"time"

	pb "Classroom/Gateway/pkg/api/lessons"
)

type Lesson struct {
	LessonID    string    `json:"lesson_id"`
	CourseID    string    `json:"course_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateLessonRequest struct {
	CourseID string `json:"course_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

func NewCreateLessonRequest(req CreateLessonRequest) *pb.CreateLessonRequest {
	return &pb.CreateLessonRequest{
		CourseId: req.CourseID,
		Title:    req.Title,
		Content:  req.Content,
	}
}

type CreateLessonResponse struct {
	LessonID string `json:"lesson_id"`
}

func NewCreateLessonResponse(resp *pb.CreateLessonResponse) CreateLessonResponse {
	return CreateLessonResponse{
		LessonID: resp.GetLessonId(),
	}
}

type GetLessonRequest struct {
	LessonID string `json:"lesson_id"`
}

func NewGetLessonRequest(req GetLessonRequest) *pb.GetLessonRequest {
	return &pb.GetLessonRequest{
		LessonId: req.LessonID,
	}
}

type GetLessonResponse struct {
	Lesson Lesson `json:"lesson"`
}

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

type GetLessonsRequest struct {
	CourseID string `json:"course_id"`
}

func NewGetLessonsRequest(req GetLessonsRequest) *pb.GetLessonsRequest {
	return &pb.GetLessonsRequest{
		CourseId: req.CourseID,
	}
}

type GetLessonsResponse struct {
	Lessons []Lesson `json:"lessons"`
}

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

type UpdateLessonRequest struct {
	LessonID string  `json:"lesson_id"`
	Title    *string `json:"title,omitempty"`
	Content  *string `json:"description,omitempty"`
}

func NewUpdateLessonRequest(req UpdateLessonRequest) *pb.UpdateLessonRequest {
	return &pb.UpdateLessonRequest{
		LessonId: req.LessonID,
		Title:    req.Title,
		Content:  req.Content,
	}
}

type UpdateLessonResponse struct {
}

func NewUpdateLessonResponse(resp *pb.UpdateLessonResponse) UpdateLessonResponse {
	return UpdateLessonResponse{}
}

type DeleteLessonRequest struct {
	LessonID string `json:"lesson_id"`
}

func NewDeleteLessonRequest(req DeleteLessonRequest) *pb.DeleteLessonRequest {
	return &pb.DeleteLessonRequest{
		LessonId: req.LessonID,
	}
}

type DeleteLessonResponse struct {
}

func NewDeleteLessonResponse(resp *pb.DeleteLessonResponse) DeleteLessonResponse {
	return DeleteLessonResponse{}
}
