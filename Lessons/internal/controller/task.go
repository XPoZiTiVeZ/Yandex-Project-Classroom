package controller

import (
	"context"
	"log/slog"

	"Classroom/Lessons/internal/dto"
	pb "Classroom/Lessons/pkg/api/task"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TaskService interface {
	Create(ctx context.Context, dto dto.CreateTaskDTO) (string, error)
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

func (c *taskController) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	dto := dto.CreateTaskDTO{
		Title:    req.Title,
		Content:  req.Description,
		CourseID: req.CourseId,
	}

	if err := c.validate.Struct(dto); err != nil {
		c.logger.Debug("invalid request", "err", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	taskID, err := c.svc.Create(ctx, dto)
	if err != nil {
		c.logger.Error("failed to create task", "err", err)
		return nil, status.Error(codes.Internal, "failed to create task")
	}

	return &pb.CreateTaskResponse{TaskId: taskID}, nil
}

func (c *taskController) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTask not implemented")
}

func (c *taskController) GetTasks(ctx context.Context, req *pb.GetTasksRequest) (*pb.GetTasksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTasks not implemented")
}

func (c *taskController) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTask not implemented")
}

func (c *taskController) ChangeStatusTask(ctx context.Context, req *pb.ChangeStatusTaskRequest) (*pb.ChangeStatusTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeStatusTask not implemented")
}

func (c *taskController) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTask not implemented")
}
