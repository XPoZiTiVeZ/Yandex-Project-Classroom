package service

import "log/slog"

type TaskRepo interface {
	//
}

type taskService struct {
	logger *slog.Logger // Для дебага и информации, ошибки логируются в контроллере
	tasks  TaskRepo
}

func NewTaskService(logger *slog.Logger, tasks TaskRepo) *taskService {
	return &taskService{logger: logger, tasks: tasks}
}
