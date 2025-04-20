package server

import (
	"Classroom/Gateway/internal/lessons"
	"Classroom/Gateway/pkg/logger"
	"log/slog"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateLessonHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[lessons.CreateLessonRequest](r.Context())

	resp, err := s.Lessons.CreateLesson(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler lessons.CreateLesson error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.AlreadyExists:
				AlreadyExists(w)
			}
		} else {
			ServiceUnavailable(w)
		}
	}

	WriteJSON(w, resp, http.StatusCreated)
}

func (s *Server) GetLessonHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[lessons.GetLessonRequest](r.Context())

	resp, err := s.Lessons.GetLesson(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler lessons.GetLesson error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				NotFound(w)
			}
		} else {
			ServiceUnavailable(w)
		}
	}

	WriteJSON(w, resp, http.StatusOK)
}

func (s *Server) GetLessonsHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[lessons.GetLessonsRequest](r.Context())

	resp, err := s.Lessons.GetLessons(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler lessons.GetLessons error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				NotFound(w)
			}
		} else {
			ServiceUnavailable(w)
		}
	}

	WriteJSON(w, resp, http.StatusOK)
}

func (s *Server) UpdateLessonHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[lessons.UpdateLessonRequest](r.Context())

	resp, err := s.Lessons.UpdateLesson(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler lessons.UpdateLesson error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				NotFound(w)
			}
		} else {
			ServiceUnavailable(w)
		}
	}

	WriteJSON(w, resp, http.StatusOK)
}

func (s *Server) DeleteLessonHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[lessons.DeleteLessonRequest](r.Context())

	resp, err := s.Lessons.DeleteLesson(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler lessons.DeleteLesson error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				NotFound(w)
			}
		} else {
			ServiceUnavailable(w)
		}
	}

	WriteJSON(w, resp, http.StatusOK)
}
