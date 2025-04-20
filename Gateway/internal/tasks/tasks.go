package tasks

import (
	pb "Classroom/Gateway/pkg/api/tasks"
	"Classroom/Gateway/pkg/config"
	"Classroom/Gateway/pkg/logger"
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TasksServiceClient struct {
	Conn           *grpc.ClientConn
	Client         pb.TasksServiceClient
	DefaultTimeout time.Duration
}

func NewTasksServiceClient(ctx context.Context, config *config.Config) (*TasksServiceClient, error) {
	address, port := config.Courses.Address, config.Courses.Port
	timeout := config.Common.Timeout

	var opts []grpc.DialOption
	opts = append(
		opts, grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", address, port), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	state := conn.GetState()
	// if !conn.WaitForStateChange(ctx, state) {
	// 	return nil, fmt.Errorf("failed to wait for state change")
	// }
	// state = conn.GetState()

	logger.Info(ctx, "Connected to gRPC Tasks", slog.String("address", address), slog.Int("port", port), slog.String("state", state.String()))

	client := pb.NewTasksServiceClient(conn)

	return &TasksServiceClient{
		Conn:           conn,
		Client:         client,
		DefaultTimeout: timeout,
	}, nil
}

func (s *TasksServiceClient) CreateTask(ctx context.Context, req CreateTaskRequest) (CreateTaskResponse, error) {
	logger.Debug(ctx, "Creating task", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := s.Client.CreateTask(ctx, NewCreateTaskRequest(req))
	if err != nil {
		return CreateTaskResponse{}, err
	}

	logger.Debug(ctx, "Tasks.CreateTask succeed")
	return NewCreateTaskResponse(resp), nil
}

func (s *TasksServiceClient) GetTask(ctx context.Context, req GetTaskRequest) (GetTaskResponse, error) {
	logger.Debug(ctx, "Getting task", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := s.Client.GetTask(ctx, NewGetTaskRequest(req))
	if err != nil {
		return GetTaskResponse{}, err
	}

	logger.Debug(ctx, "Tasks.GetTask succeed")
	return NewGetTaskResponse(resp), nil
}

func (s *TasksServiceClient) GetTasks(ctx context.Context, req GetTasksRequest) (GetTasksResponse, error) {
	logger.Debug(ctx, "Getting tasks", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := s.Client.GetTasks(ctx, NewGetTasksRequest(req))
	if err != nil {
		return GetTasksResponse{}, err
	}

	logger.Debug(ctx, "Tasks.GetTasks succeed")
	return NewGetTasksResponse(resp), nil
}

func (s *TasksServiceClient) UpdateTask(ctx context.Context, req UpdateTaskRequest) (UpdateTaskResponse, error) {
	logger.Debug(ctx, "Updating task", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := s.Client.UpdateTask(ctx, NewUpdateTaskRequest(req))
	if err != nil {
		return UpdateTaskResponse{}, err
	}

	logger.Debug(ctx, "Tasks.UpdateTask succeed")
	return NewUpdateTaskResponse(resp), nil
}

func (s *TasksServiceClient) ChangeStatusTask(ctx context.Context, req ChangeStatusTaskRequest) (ChangeStatusTaskResponse, error) {
	logger.Debug(ctx, "Changing task status", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := s.Client.ChangeStatusTask(ctx, NewChangeStatusTaskRequest(req))
	if err != nil {
		return ChangeStatusTaskResponse{}, err
	}

	logger.Debug(ctx, "Tasks.ChangeStatusTask succeed")
	return NewChangeStatusTaskResponse(resp), nil
}

func (s *TasksServiceClient) DeleteTask(ctx context.Context, req DeleteTaskRequest) (DeleteTaskResponse, error) {
	logger.Debug(ctx, "Deleting task", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := s.Client.DeleteTask(ctx, NewDeleteTaskRequest(req))
	if err != nil {
		return DeleteTaskResponse{}, err
	}

	logger.Debug(ctx, "Tasks.DeleteTask succeed")
	return NewDeleteTaskResponse(resp), nil
}
