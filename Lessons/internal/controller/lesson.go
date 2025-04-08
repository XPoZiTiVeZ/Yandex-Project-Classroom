package controller

import (
	"log/slog"

	pb "Classroom/Lessons/pkg/api/lesson"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
)

type LessonService interface {
	//
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
