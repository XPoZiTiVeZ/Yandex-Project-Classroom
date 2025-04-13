package service

import (
	"Classroom/Tasks/internal/domain"
	"Classroom/Tasks/internal/dto"
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type TaskRepo interface {
	Create(ctx context.Context, payload dto.CreateTaskDTO) (domain.Task, error)
	GetTaskByID(ctx context.Context, id string) (domain.Task, error)
	ListByCourseID(ctx context.Context, course_id string) ([]domain.Task, error)
	Update(ctx context.Context, task domain.Task) error
	Delete(ctx context.Context, id string) error
	GetTaskStatus(ctx context.Context, taskID, userID string) (domain.TaskStatus, error)
	UpdateTaskStatus(ctx context.Context, status domain.TaskStatus) error
	CreateTaskStatus(ctx context.Context, status domain.TaskStatus) error
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

func (s *taskService) ListByCourseID(ctx context.Context, courseID string) ([]domain.Task, error) {
	return s.tasks.ListByCourseID(ctx, courseID)
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

	s.logger.Info("task updated", "id", task.ID, "title", task.Title, "content", task.Content)
	return task, nil
}

func (s *taskService) Delete(ctx context.Context, id string) error {
	return s.tasks.Delete(ctx, id)
}

func (s *taskService) UpdateTaskStatus(ctx context.Context, taskID, userID string) (domain.TaskStatus, error) {
	currentStatus, err := s.tasks.GetTaskStatus(ctx, taskID, userID)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return domain.TaskStatus{}, fmt.Errorf("failed to get task status: %w", err)
	}

	// Если записи еще нет, создаем и ставим true
	if errors.Is(err, domain.ErrNotFound) {
		currentStatus = domain.TaskStatus{
			UserID:      userID,
			TaskID:      taskID,
			IsCompleted: true,
		}
		if err := s.tasks.CreateTaskStatus(ctx, currentStatus); err != nil {
			return domain.TaskStatus{}, fmt.Errorf("failed to create task status: %w", err)
		}

		s.logger.Info("task status created", "task_id", taskID, "user_id", userID, "status", currentStatus.IsCompleted)
		return currentStatus, nil
	}

	// Меняем статус на противоположный
	currentStatus.IsCompleted = !currentStatus.IsCompleted
	if err := s.tasks.UpdateTaskStatus(ctx, currentStatus); err != nil {
		return domain.TaskStatus{}, fmt.Errorf("failed to update task status: %w", err)
	}

	s.logger.Info("task status updated", "task_id", taskID, "user_id", userID, "status", currentStatus.IsCompleted)
	return currentStatus, nil
}
