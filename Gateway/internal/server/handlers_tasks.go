package server

import (
	he "Classroom/Gateway/internal/errors"
	"Classroom/Gateway/internal/tasks"
	"Classroom/Gateway/pkg/logger"
	"log/slog"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[tasks.CreateTaskRequest](r.Context())

	resp, err := s.Tasks.CreateTask(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler tasks.CreateTask error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.AlreadyExists:
				he.AlreadyExists(w)
			}
		} else {
			he.ServiceUnavailable(w)
		}
	}

	WriteJSON(w, resp, http.StatusOK)
}

func (s *Server) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[tasks.GetTaskRequest](r.Context())

	resp, err := s.Tasks.GetTask(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler tasks.GetTask error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				he.NotFound(w)
			}
		} else {
			he.ServiceUnavailable(w)
		}
	}

	WriteJSON(w, resp, http.StatusOK)
}

func (s *Server) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[tasks.GetTasksRequest](r.Context())

	resp, err := s.Tasks.GetTasks(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler tasks.GetTasks error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				he.NotFound(w)
			}
		} else {
			he.ServiceUnavailable(w)
		}
	}

	WriteJSON(w, resp, http.StatusOK)
}

func (s *Server) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[tasks.UpdateTaskRequest](r.Context())

	resp, err := s.Tasks.UpdateTask(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler tasks.UpdateTask error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				he.NotFound(w)
			}
		} else {
			he.ServiceUnavailable(w)
		}
	}

	WriteJSON(w, resp, http.StatusOK)
}

func (s *Server) ChangeStatusTaskHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[tasks.ChangeStatusTaskRequest](r.Context())

	resp, err := s.Tasks.ChangeStatusTask(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler tasks.ChangeStatusTask error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				he.NotFound(w)
			}
		} else {
			he.ServiceUnavailable(w)
		}
	}

	WriteJSON(w, resp, http.StatusOK)
}

func (s *Server) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[tasks.DeleteTaskRequest](r.Context())

	resp, err := s.Tasks.DeleteTask(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler tasks.DeleteTask error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				he.NotFound(w)
			}
		} else {
			he.ServiceUnavailable(w)
		}
	}

	WriteJSON(w, resp, http.StatusOK)
}
