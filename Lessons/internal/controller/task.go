package controller

import (
	"log/slog"

	pb "Classroom/Lessons/pkg/api/task"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
)

type TaskService interface {
	//
}

type taskController struct {
	svc      TaskService
	logger   *slog.Logger // Для логирования ошибок и дебага запросов
	validate *validator.Validate
	pb.UnimplementedTasksServiceServer
}

func NewTaskController(logger *slog.Logger, svc TaskService) *taskController {
	validate := validator.New()
	return &taskController{
		svc:      svc,
		logger:   logger,
		validate: validate,
	}
}

func (c *taskController) Init(srv *grpc.Server) {
	pb.RegisterTasksServiceServer(srv, c)
}
