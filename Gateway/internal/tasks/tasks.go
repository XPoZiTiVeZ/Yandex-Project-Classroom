package tasks

import (
	pb "Classroom/Gateway/pkg/api/tasks"
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TasksServiceClient struct {
	Conn           *grpc.ClientConn
	Client         *pb.TasksServiceClient
	DefaultTimeout time.Duration
}

func NewTasksServiceClient(address string, port int, DefaultTimeout *time.Duration) (*TasksServiceClient, error) {
	var opts []grpc.DialOption
	opts = append(
		opts, grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", address, port), opts...)
	if err != nil {
		slog.Error("fail to dial: %v", slog.Any("error", err))
		return nil, err
	}

	state := conn.GetState()
	slog.Info("Connected to grpc Tasks", slog.String("address", address), slog.Int("port", port), slog.String("state", state.String()))

	client := pb.NewTasksServiceClient(conn)

	timeout := 10 * time.Second
	if DefaultTimeout != nil {
		timeout = *DefaultTimeout
	}

	return &TasksServiceClient{
		Conn:           conn,
		Client:         &client,
		DefaultTimeout: timeout,
	}, nil
}

func (s *TasksServiceClient) CreateTask(ctx context.Context, req CreateTaskRequest) (CreateTaskResponse, error) {
	slog.Debug("creating task", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := (*s.Client).CreateTask(ctx, NewCreateTaskRequest(req))
	if err != nil {
		return CreateTaskResponse{}, err
	}

	slog.Debug("tasks.CreateTask succeed")
	return NewCreateTaskResponse(resp), nil
}

func (s *TasksServiceClient) GetTask(ctx context.Context, req GetTaskRequest) (GetTaskResponse, error) {
	slog.Debug("getting task", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := (*s.Client).GetTask(ctx, NewGetTaskRequest(req))
	if err != nil {
		return GetTaskResponse{}, err
	}

	slog.Debug("tasks.GetTask succeed")
	return NewGetTaskResponse(resp), nil
}

func (s *TasksServiceClient) GetTasks(ctx context.Context, req GetTasksRequest) (GetTasksResponse, error) {
	slog.Debug("getting tasks", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := (*s.Client).GetTasks(ctx, NewGetTasksRequest(req))
	if err != nil {
		return GetTasksResponse{}, err
	}

	slog.Debug("tasks.GetTasks succeed")
	return NewGetTasksResponse(resp), nil
}

func (s *TasksServiceClient) UpdateTask(ctx context.Context, req UpdateTaskRequest) (UpdateTaskResponse, error) {
	slog.Debug("updating task", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := (*s.Client).UpdateTask(ctx, NewUpdateTaskRequest(req))
	if err != nil {
		return UpdateTaskResponse{}, err
	}

	slog.Debug("tasks.UpdateTask succeed")
	return NewUpdateTaskResponse(resp), nil
}

func (s *TasksServiceClient) ChangeStatusTask(ctx context.Context, req ChangeStatusTaskRequest) (ChangeStatusTaskResponse, error) {
	slog.Debug("changing task status", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := (*s.Client).ChangeStatusTask(ctx, NewChangeStatusTaskRequest(req))
	if err != nil {
		return ChangeStatusTaskResponse{}, err
	}

	slog.Debug("tasks.ChangeStatusTask succeed")
	return NewChangeStatusTaskResponse(resp), nil
}

func (s *TasksServiceClient) DeleteTask(ctx context.Context, req DeleteTaskRequest) (DeleteTaskResponse, error) {
	slog.Debug("deleting task", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := (*s.Client).DeleteTask(ctx, NewDeleteTaskRequest(req))
	if err != nil {
		return DeleteTaskResponse{}, err
	}

	slog.Debug("tasks.DeleteTask succeed")
	return NewDeleteTaskResponse(resp), nil
}