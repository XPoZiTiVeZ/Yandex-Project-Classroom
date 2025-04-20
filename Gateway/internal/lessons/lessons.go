package lessons

import (
	pb "Classroom/Gateway/pkg/api/lessons"
	"Classroom/Gateway/pkg/config"
	"Classroom/Gateway/pkg/logger"
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LessonsServiceClient struct {
	Conn           *grpc.ClientConn
	Client         pb.LessonsServiceClient
	DefaultTimeout time.Duration
}

func NewLessonsServiceClient(ctx context.Context, config *config.Config) (*LessonsServiceClient, error) {
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

	logger.Info(ctx, "Connected to gRPC Lessons", slog.String("address", address), slog.Int("port", port), slog.String("state", state.String()))

	client := pb.NewLessonsServiceClient(conn)

	return &LessonsServiceClient{
		Conn:           conn,
		Client:         client,
		DefaultTimeout: timeout,
	}, nil
}

func (s *LessonsServiceClient) CreateLesson(ctx context.Context, req CreateLessonRequest) (CreateLessonResponse, error) {
	logger.Debug(ctx, "Creating lesson", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := s.Client.CreateLesson(ctx, NewCreateLessonRequest(req))
	if err != nil {
		return CreateLessonResponse{}, err
	}

	logger.Debug(ctx, "Lessons.CreateLesson succeed")
	return NewCreateLessonResponse(resp), nil
}

func (s *LessonsServiceClient) GetLesson(ctx context.Context, req GetLessonRequest) (GetLessonResponse, error) {
	logger.Debug(ctx, "Getting lesson", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := s.Client.GetLesson(ctx, NewGetLessonRequest(req))
	if err != nil {
		return GetLessonResponse{}, err
	}

	logger.Debug(ctx, "Lessons.GetLesson succeed")
	return NewGetLessonResponse(resp), nil
}

func (s *LessonsServiceClient) GetLessons(ctx context.Context, req GetLessonsRequest) (GetLessonsResponse, error) {
	logger.Debug(ctx, "Getting lessons", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := s.Client.GetLessons(ctx, NewGetLessonsRequest(req))
	if err != nil {
		return GetLessonsResponse{}, err
	}

	logger.Debug(ctx, "Lessons.GetLessons succeed")
	return NewGetLessonsResponse(resp), nil
}

func (s *LessonsServiceClient) UpdateLesson(ctx context.Context, req UpdateLessonRequest) (UpdateLessonResponse, error) {
	logger.Debug(ctx, "Updating lesson", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := s.Client.UpdateLesson(ctx, NewUpdateLessonRequest(req))
	if err != nil {
		return UpdateLessonResponse{}, err
	}

	logger.Debug(ctx, "Lessons.UpdateLesson succeed")
	return NewUpdateLessonResponse(resp), nil
}

func (s *LessonsServiceClient) DeleteLesson(ctx context.Context, req DeleteLessonRequest) (DeleteLessonResponse, error) {
	logger.Debug(ctx, "Deleting lesson", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := s.Client.DeleteLesson(ctx, NewDeleteLessonRequest(req))
	if err != nil {
		return DeleteLessonResponse{}, err
	}

	logger.Debug(ctx, "Lessons.DeleteLesson succeed")
	return NewDeleteLessonResponse(resp), nil
}
