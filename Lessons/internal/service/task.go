package service

import (
	"Classroom/Lessons/internal/domain"
	"Classroom/Lessons/internal/dto"
	"context"
	"fmt"
	"log/slog"
)

type TaskRepo interface {
	Create(ctx context.Context, payload dto.CreateTaskDTO) (domain.Task, error)
	GetTaskByID(ctx context.Context, id string) (domain.Task, error)
	ListByCourseID(ctx context.Context, course_id string) ([]domain.Task, error)
}

type taskService struct {
	logger *slog.Logger // Для дебага и информации, ошибки логируются в контроллере
	tasks  TaskRepo
}

func NewTaskService(logger *slog.Logger, tasks TaskRepo) *taskService {
	return &taskService{logger: logger, tasks: tasks}
}

func (s *taskService) Create(ctx context.Context, payload dto.CreateTaskDTO) (string, error) {
	task, err := s.tasks.Create(ctx, payload)
	if err != nil {
		return "", fmt.Errorf("failed to create task: %w", err)
	}

	s.logger.Info("task created", "id", task.ID, "title", task.Title)
	return task.ID, nil
}

func (s *taskService) GetTaskByID(ctx context.Context, id string) (domain.Task, error) {
	return s.tasks.GetTaskByID(ctx, id)
}

func (s *taskService) ListByCourseID(ctx context.Context, course_id string) ([]domain.Task, error) {
	return s.tasks.ListByCourseID(ctx, course_id)
}
