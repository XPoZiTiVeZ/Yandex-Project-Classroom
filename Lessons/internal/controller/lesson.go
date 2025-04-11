package controller

import (
	"context"
	"errors"
	"log/slog"

	"Classroom/Lessons/internal/domain"
	"Classroom/Lessons/internal/dto"
	pb "Classroom/Lessons/pkg/api/lesson"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type LessonService interface {
	Create(ctx context.Context, dto dto.CreateLessonDTO) (domain.Lesson, error)
	GetByID(ctx context.Context, id string) (domain.Lesson, error)
	ListByCourseID(ctx context.Context, courseID string) ([]domain.Lesson, error)
	Update(ctx context.Context, dto dto.UpdateLessonDTO) (domain.Lesson, error)
	Delete(ctx context.Context, id string) error
}

type lessonController struct {
	svc      LessonService
	logger   *slog.Logger // Для логирования ошибок и дебага запросов
	validate *validator.Validate
	pb.UnimplementedLessonServiceServer
}

func NewLessonController(logger *slog.Logger, svc LessonService) *lessonController {
	validate := validator.New()
	return &lessonController{
		svc:      svc,
		logger:   logger,
		validate: validate,
	}
}

func (c *lessonController) Init(srv *grpc.Server) {
	pb.RegisterLessonServiceServer(srv, c)
}

func (c *lessonController) CreateLesson(ctx context.Context, req *pb.CreateLessonRequest) (*pb.CreateLessonResponse, error) {
	dto := dto.CreateLessonDTO{
		Title:       req.Title,
		Description: req.Description,
		CourseID:    req.CourseId,
	}
	if err := c.validate.Struct(dto); err != nil {
		c.logger.Debug("invalid request", "err", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	lesson, err := c.svc.Create(ctx, dto)
	if err != nil {
		c.logger.Error("failed to create lesson", "err", err)
		return nil, status.Error(codes.Internal, "failed to create lesson")
	}

	return &pb.CreateLessonResponse{LessonId: lesson.ID}, nil
}

func (c *lessonController) GetLesson(ctx context.Context, req *pb.GetLessonRequest) (*pb.GetLessonResponse, error) {
	if err := c.validate.Var(req.LessonId, "required,uuid"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid lesson id")
	}

	lesson, err := c.svc.GetByID(ctx, req.LessonId)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "lesson not found")
	}
	if err != nil {
		c.logger.Error("failed to get lesson", "err", err, "id", req.LessonId)
		return nil, status.Error(codes.Internal, "failed to get lesson")
	}

	return &pb.GetLessonResponse{
		Lesson: &pb.Lesson{
			LessonId:    lesson.ID,
			Title:       lesson.Title,
			Description: lesson.Description,
			CourseId:    lesson.CourseID,
			CreatedAt:   timestamppb.New(lesson.CreatedAt),
		},
	}, nil
}

func (c *lessonController) GetLessons(ctx context.Context, req *pb.GetLessonsRequest) (*pb.GetLessonsResponse, error) {
	if err := c.validate.Var(req.CourseId, "required,uuid"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid task id")
	}
	lessons, err := c.svc.ListByCourseID(ctx, req.CourseId)
	if err != nil {
		c.logger.Error("failed to get tasks", "err", err, "course_id", req.CourseId)
		return nil, status.Error(codes.Internal, "failed to get tasks")
	}

	pbLessons := make([]*pb.Lesson, len(lessons))
	for i, lesson := range lessons {
		pbLessons[i] = &pb.Lesson{
			LessonId:    lesson.ID,
			Title:       lesson.Title,
			Description: lesson.Description,
			CourseId:    lesson.CourseID,
			CreatedAt:   timestamppb.New(lesson.CreatedAt),
		}
	}
	return &pb.GetLessonsResponse{Lessons: pbLessons}, nil
}

func (c *lessonController) UpdateLesson(ctx context.Context, req *pb.UpdateLessonRequest) (*pb.UpdateLessonResponse, error) {
	dto := dto.UpdateLessonDTO{
		LessonID:    req.LessonId,
		Title:       req.Title,
		Description: req.Description,
	}
	if err := c.validate.Struct(dto); err != nil {
		c.logger.Debug("invalid request", "err", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	lesson, err := c.svc.Update(ctx, dto)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "lesson not found")
	}
	if err != nil {
		c.logger.Error("failed to update lesson", "err", err, "id", req.LessonId)
		return nil, status.Error(codes.Internal, "failed to update lesson")
	}

	return &pb.UpdateLessonResponse{
		Lesson: &pb.Lesson{
			LessonId:    lesson.ID,
			Title:       lesson.Title,
			Description: lesson.Description,
			CourseId:    lesson.CourseID,
			CreatedAt:   timestamppb.New(lesson.CreatedAt),
		},
	}, nil
}

func (c *lessonController) DeleteLesson(ctx context.Context, req *pb.DeleteLessonRequest) (*pb.DeleteLessonResponse, error) {
	if err := c.validate.Var(req.LessonId, "required,uuid"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid lesson id")
	}

	err := c.svc.Delete(ctx, req.LessonId)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "lesson not found")
	}
	if err != nil {
		c.logger.Error("failed to delete lesson", "err", err, "id", req.LessonId)
		return nil, status.Error(codes.Internal, "failed to delete lesson")
	}
	return &pb.DeleteLessonResponse{Success: true}, nil
}
