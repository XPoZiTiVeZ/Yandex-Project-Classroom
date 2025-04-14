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
	GetByID(ctx context.Context, id string) (domain.Task, error)
	ListByCourseID(ctx context.Context, courseID string) ([]domain.Task, error)
	ListByStudentID(ctx context.Context, studentID, courseID string) ([]domain.StudentTask, error)
	Update(ctx context.Context, task domain.Task) error
	Delete(ctx context.Context, id string) error
}

type StatusRepo interface {
	Create(ctx context.Context, status domain.TaskStatus) error
	Update(ctx context.Context, status domain.TaskStatus) error
	Get(ctx context.Context, taskID, userID string) (domain.TaskStatus, error)
	ListByTaskID(ctx context.Context, taskID string) ([]domain.TaskStatus, error)
}

type taskService struct {
	logger   *slog.Logger // Для дебага и информации, ошибки логируются в контроллере
	tasks    TaskRepo
	statuses StatusRepo
}

func NewTaskService(logger *slog.Logger, tasks TaskRepo, statuses StatusRepo) *taskService {
	return &taskService{logger: logger, tasks: tasks, statuses: statuses}
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
	return s.tasks.GetByID(ctx, id)
}

// Просто получение всех заданий с курса
func (s *taskService) ListByCourseID(ctx context.Context, courseID string) ([]domain.Task, error) {
	return s.tasks.ListByCourseID(ctx, courseID)
}

// Обновление задания
func (s *taskService) Update(ctx context.Context, dto dto.UpdateTaskDTO) (domain.Task, error) {
	task, err := s.tasks.GetByID(ctx, dto.TaskID)
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

func (s *taskService) ToggleTaskStatus(ctx context.Context, taskID, userID string) (domain.TaskStatus, error) {
	currentStatus, err := s.statuses.Get(ctx, taskID, userID)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return domain.TaskStatus{}, fmt.Errorf("failed to get task status: %w", err)
	}

	// Если записи еще нет, создаем и ставим true
	if errors.Is(err, domain.ErrNotFound) {
		currentStatus = domain.TaskStatus{
			UserID:    userID,
			TaskID:    taskID,
			Completed: true,
		}
		if err := s.statuses.Create(ctx, currentStatus); err != nil {
			return domain.TaskStatus{}, fmt.Errorf("failed to create task status: %w", err)
		}

		s.logger.Info("task status created", "task_id", taskID, "user_id", userID, "status", currentStatus.Completed)
		return currentStatus, nil
	}

	// Меняем статус на противоположный
	currentStatus.Completed = !currentStatus.Completed
	if err := s.statuses.Update(ctx, currentStatus); err != nil {
		return domain.TaskStatus{}, fmt.Errorf("failed to update task status: %w", err)
	}

	s.logger.Info("task status updated", "task_id", taskID, "user_id", userID, "status", currentStatus.Completed)
	return currentStatus, nil
}

func (s *taskService) ListTaskStatuses(ctx context.Context, taskID string) ([]domain.TaskStatus, error) {
	// Берутся все пользователи с курса, и если записи по статусу нет, то ставится false
	return s.statuses.ListByTaskID(ctx, taskID)
}

func (s *taskService) ListByStudentID(ctx context.Context, studentID, courseID string) ([]domain.StudentTask, error) {
	// Берутся все задания с курса, и если записи по статусу нет, то в статус ставится false
	return s.tasks.ListByStudentID(ctx, studentID, courseID)
}
