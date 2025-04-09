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
	Update(ctx context.Context, task domain.Task) error
	Delete(ctx context.Context, id string) error
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

func (s *taskService) Update(ctx context.Context, dto dto.UpdateTaskDTO) (domain.Task, error) {
	task, err := s.tasks.GetTaskByID(ctx, dto.TaskID)
	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to get task: %w", err)
	}

	if dto.Title != nil {
		task.Title = *dto.Title
	}
	if dto.Content != nil {
		task.Content = *dto.Content
	}

	if err = s.tasks.Update(ctx, task); err != nil {
		return domain.Task{}, fmt.Errorf("failed to update task: %w", err)
	}

	return task, nil
}

func (s *taskService) Delete(ctx context.Context, id string) error {
	return s.tasks.Delete(ctx, id)
}
