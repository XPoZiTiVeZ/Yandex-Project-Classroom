package server

import (
	he "Classroom/Gateway/internal/errors"
	"Classroom/Gateway/internal/lessons"
	"log/slog"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateLessonHandler(w http.ResponseWriter, r *http.Request) (error, any) {
	var body lessons.CreateLessonRequest = r.Context().Value("body").(lessons.CreateLessonRequest)

	resp, err := s.Lessons.CreateLesson(r.Context(), body)
	if err != nil {
		slog.Error("handler lessons.CreateLesson error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.AlreadyExists:
				he.AlreadyExists(w)
			}
		} else {
			he.ServiceUnavailable(w)
		}
	}

	return err, resp
}

func (s *Server) GetLessonHandler(w http.ResponseWriter, r *http.Request) (error, any) {
	var body lessons.GetLessonRequest = r.Context().Value("body").(lessons.GetLessonRequest)

	resp, err := s.Lessons.GetLesson(r.Context(), body)
	if err != nil {
		slog.Error("handler lessons.GetLesson error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				he.NotFound(w)
			}
		} else {
			he.ServiceUnavailable(w)
		}
	}

	return err, resp
}

func (s *Server) GetLessonsHandler(w http.ResponseWriter, r *http.Request) (error, any) {
	var body lessons.GetLessonsRequest = r.Context().Value("body").(lessons.GetLessonsRequest)

	resp, err := s.Lessons.GetLessons(r.Context(), body)
	if err != nil {
		slog.Error("handler lessons.GetLessons error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				he.NotFound(w)
			}
		} else {
			he.ServiceUnavailable(w)
		}
	}

	return err, resp
}

func (s *Server) UpdateLessonHandler(w http.ResponseWriter, r *http.Request) (error, any) {
	var body lessons.UpdateLessonRequest = r.Context().Value("body").(lessons.UpdateLessonRequest)

	resp, err := s.Lessons.UpdateLesson(r.Context(), body)
	if err != nil {
		slog.Error("handler lessons.UpdateLesson error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				he.NotFound(w)
			}
		} else {
			he.ServiceUnavailable(w)
		}
	}

	return err, resp
}

func (s *Server) DeleteLessonHandler(w http.ResponseWriter, r *http.Request) (error, any) {
	var body lessons.DeleteLessonRequest = r.Context().Value("body").(lessons.DeleteLessonRequest)

	resp, err := s.Lessons.DeleteLesson(r.Context(), body)
	if err != nil {
		slog.Error("handler lessons.DeleteLesson error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				he.NotFound(w)
			}
		} else {
			he.ServiceUnavailable(w)
		}
	}

	return err, resp
}