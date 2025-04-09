package controller

import (
	"context"
	"errors"
	"log/slog"

	"Classroom/Lessons/internal/domain"
	"Classroom/Lessons/internal/dto"
	pb "Classroom/Lessons/pkg/api/task"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TaskService interface {
	Create(ctx context.Context, dto dto.CreateTaskDTO) (string, error)
	GetTaskByID(ctx context.Context, id string) (domain.Task, error)
	ListByCourseID(ctx context.Context, course_id string) ([]domain.Task, error)
	Update(ctx context.Context, dto dto.UpdateTaskDTO) (domain.Task, error)
	Delete(ctx context.Context, id string) error
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
	if err := c.validate.Var(req.TaskId, "required,uuid"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid task id")
	}

	task, err := c.svc.GetTaskByID(ctx, req.TaskId)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "task not found")
	}
	if err != nil {
		c.logger.Error("failed to get task", "err", err, "id", req.TaskId)
		return nil, status.Error(codes.Internal, "failed to get task")
	}

	return &pb.GetTaskResponse{
		Task: &pb.Task{
			TaskId:    task.ID,
			Title:     task.Title,
			Content:   task.Content,
			CourseId:  task.CourseID,
			CreatedAt: timestamppb.New(task.CreatedAt),
		},
	}, nil
}

func (c *taskController) GetTasks(ctx context.Context, req *pb.GetTasksRequest) (*pb.GetTasksResponse, error) {
	if err := c.validate.Var(req.CourseId, "required,uuid"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid task id")
	}
	tasks, err := c.svc.ListByCourseID(ctx, req.CourseId)
	if err != nil {
		c.logger.Error("failed to get tasks", "err", err, "course_id", req.CourseId)
		return nil, status.Error(codes.Internal, "failed to get tasks")
	}

	pbTasks := make([]*pb.Task, len(tasks))
	for i, task := range tasks {
		pbTasks[i] = &pb.Task{
			TaskId:    task.ID,
			Title:     task.Title,
			Content:   task.Content,
			CourseId:  task.CourseID,
			CreatedAt: timestamppb.New(task.CreatedAt),
		}
	}
	return &pb.GetTasksResponse{Tasks: pbTasks}, nil
}

func (c *taskController) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {
	dto := dto.UpdateTaskDTO{
		Title:   req.Title,
		Content: req.Content,
		TaskID:  req.TaskId,
	}
	if err := c.validate.Struct(dto); err != nil {
		c.logger.Debug("invalid request", "err", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	task, err := c.svc.Update(ctx, dto)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "task not found")
	}
	if err != nil {
		c.logger.Error("failed to update task", "err", err, "id", req.TaskId)
		return nil, status.Error(codes.Internal, "failed to update task")
	}

	return &pb.UpdateTaskResponse{
		Task: &pb.Task{
			TaskId:    task.ID,
			Title:     task.Title,
			Content:   task.Content,
			CourseId:  task.CourseID,
			CreatedAt: timestamppb.New(task.CreatedAt),
		},
	}, nil
}

func (c *taskController) ChangeStatusTask(ctx context.Context, req *pb.ChangeStatusTaskRequest) (*pb.ChangeStatusTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeStatusTask not implemented")
}

func (c *taskController) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	if err := c.validate.Var(req.TaskId, "required,uuid"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid task id")
	}
	err := c.svc.Delete(ctx, req.TaskId)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "task not found")
	}
	if err != nil {
		c.logger.Error("failed to delete task", "err", err, "id", req.TaskId)
		return nil, status.Error(codes.Internal, "failed to delete task")
	}
	return &pb.DeleteTaskResponse{Success: true}, nil
}
