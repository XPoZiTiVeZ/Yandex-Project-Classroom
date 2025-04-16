package lessons

import (
	pb "Classroom/Gateway/pkg/api/lessons"
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LessonsServiceClient struct {
	Conn           *grpc.ClientConn
	Client         *pb.LessonsServiceClient
	DefaultTimeout time.Duration
}

func NewLessonsServiceClient(address string, port int, DefaultTimeout *time.Duration) (*LessonsServiceClient, error) {
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
	slog.Info("Connected to grpc Lessons", slog.String("address", address), slog.Int("port", port), slog.String("state", state.String()))

	client := pb.NewLessonsServiceClient(conn)

	timeout := 10 * time.Second
	if DefaultTimeout != nil {
		timeout = *DefaultTimeout
	}

	return &LessonsServiceClient{
		Conn:           conn,
		Client:         &client,
		DefaultTimeout: timeout,
	}, nil
}

func (s *LessonsServiceClient) CreateLesson(ctx context.Context, req CreateLessonRequest) (CreateLessonResponse, error) {
	slog.Debug("creating lesson", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := (*s.Client).CreateLesson(ctx, NewCreateLessonRequest(req))
	if err != nil {
		return CreateLessonResponse{}, err
	}

	slog.Debug("lessons.CreateLesson succeed")
	return NewCreateLessonResponse(resp), nil
}

func (s *LessonsServiceClient) GetLesson(ctx context.Context, req GetLessonRequest) (GetLessonResponse, error) {
	slog.Debug("getting lesson", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := (*s.Client).GetLesson(ctx, NewGetLessonRequest(req))
	if err != nil {
		return GetLessonResponse{}, err
	}

	slog.Debug("lessons.GetLesson succeed")
	return NewGetLessonResponse(resp), nil
}

func (s *LessonsServiceClient) GetLessons(ctx context.Context, req GetLessonsRequest) (GetLessonsResponse, error) {
	slog.Debug("getting lessons", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := (*s.Client).GetLessons(ctx, NewGetLessonsRequest(req))
	if err != nil {
		return GetLessonsResponse{}, err
	}

	slog.Debug("lessons.GetLessons succeed")
	return NewGetLessonsResponse(resp), nil
}

func (s *LessonsServiceClient) UpdateLesson(ctx context.Context, req UpdateLessonRequest) (UpdateLessonResponse, error) {
	slog.Debug("updating lesson", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := (*s.Client).UpdateLesson(ctx, NewUpdateLessonRequest(req))
	if err != nil {
		return UpdateLessonResponse{}, err
	}

	slog.Debug("lessons.UpdateLesson succeed")
	return NewUpdateLessonResponse(resp), nil
}

func (s *LessonsServiceClient) DeleteLesson(ctx context.Context, req DeleteLessonRequest) (DeleteLessonResponse, error) {
	slog.Debug("deleting lesson", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := (*s.Client).DeleteLesson(ctx, NewDeleteLessonRequest(req))
	if err != nil {
		return DeleteLessonResponse{}, err
	}

	slog.Debug("lessons.DeleteLesson succeed")
	return NewDeleteLessonResponse(resp), nil
}