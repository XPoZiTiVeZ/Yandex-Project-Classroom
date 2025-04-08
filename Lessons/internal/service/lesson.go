package service

import "log/slog"

type LessonRepo interface {
	//
}

type lessonService struct {
	logger  *slog.Logger // Для дебага и информации, ошибки логируются в контроллере
	lessons LessonRepo
}

func NewLessonService(logger *slog.Logger, lessons LessonRepo) *lessonService {
	return &lessonService{logger: logger, lessons: lessons}
}
