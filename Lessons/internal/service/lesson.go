package service

import (
	"Classroom/Lessons/internal/domain"
	"Classroom/Lessons/internal/dto"
	"context"
	"fmt"
	"log/slog"
)

type LessonRepo interface {
	Create(ctx context.Context, dto dto.CreateLessonDTO) (domain.Lesson, error)
	GetByID(ctx context.Context, id string) (domain.Lesson, error)
	ListByCourseID(ctx context.Context, courseID string) ([]domain.Lesson, error)
	Update(ctx context.Context, dto dto.UpdateLessonDTO) (domain.Lesson, error)
	Delete(ctx context.Context, id string) error
}

type lessonService struct {
	logger  *slog.Logger // Для дебага и информации, ошибки логируются в контроллере
	lessons LessonRepo
}

func NewLessonService(logger *slog.Logger, lessons LessonRepo) *lessonService {
	return &lessonService{logger: logger, lessons: lessons}
}

func (s *lessonService) Create(ctx context.Context, dto dto.CreateLessonDTO) (domain.Lesson, error) {
	lesson, err := s.lessons.Create(ctx, dto)
	if err != nil {
		return domain.Lesson{}, fmt.Errorf("failed to create lesson: %w", err)
	}
	s.logger.Info("lesson created", "id", lesson.ID, "title", lesson.Title)
	return lesson, nil
}

func (s *lessonService) GetByID(ctx context.Context, id string) (domain.Lesson, error) {
	return s.lessons.GetByID(ctx, id)
}

func (s *lessonService) ListByCourseID(ctx context.Context, courseID string) ([]domain.Lesson, error) {
	return s.lessons.ListByCourseID(ctx, courseID)
}

func (s *lessonService) Update(ctx context.Context, dto dto.UpdateLessonDTO) (domain.Lesson, error) {
	return s.lessons.Update(ctx, dto)
}

func (s *lessonService) Delete(ctx context.Context, id string) error {
	return s.lessons.Delete(ctx, id)
}
