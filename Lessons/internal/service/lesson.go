package service

import (
	"Classroom/Lessons/internal/domain"
	"Classroom/Lessons/internal/dto"
	"Classroom/Lessons/pkg/events"
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
	CourseExists(ctx context.Context, courseID string) (bool, error)
}

type Producer interface {
	PublishLessonCreated(msg events.LessonCreated) error
}

type lessonService struct {
	logger   *slog.Logger // Для дебага и информации, ошибки логируются в контроллере
	lessons  LessonRepo
	producer Producer
}

func NewLessonService(logger *slog.Logger, lessons LessonRepo, producer Producer) *lessonService {
	return &lessonService{logger: logger, lessons: lessons, producer: producer}
}

func (s *lessonService) Create(ctx context.Context, dto dto.CreateLessonDTO) (domain.Lesson, error) {
	courseExists, err := s.lessons.CourseExists(ctx, dto.CourseID)
	if err != nil {
		return domain.Lesson{}, fmt.Errorf("failed to check course exists: %w", err)
	}
	if !courseExists {
		return domain.Lesson{}, domain.ErrNotFound
	}

	lesson, err := s.lessons.Create(ctx, dto)
	if err != nil {
		return domain.Lesson{}, fmt.Errorf("failed to create lesson: %w", err)
	}
	s.logger.Info("lesson created", "id", lesson.ID, "title", lesson.Title)

	msg := events.LessonCreated{
		LessonID: lesson.ID,
		CourseID: lesson.CourseID,
	}
	if err := s.producer.PublishLessonCreated(msg); err != nil {
		s.logger.Error("failed to publish lesson created event", "err", err)
	}
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
