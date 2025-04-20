package server

import (
	he "Classroom/Gateway/internal/errors"
	"Classroom/Gateway/internal/tasks"
	"log/slog"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateTaskHandler(w http.ResponseWriter, r *http.Request) (any, error) {
	var body tasks.CreateTaskRequest = r.Context().Value("body").(tasks.CreateTaskRequest)

	resp, err := s.Tasks.CreateTask(r.Context(), body)
	if err != nil {
		slog.Error("handler tasks.CreateTask error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.AlreadyExists:
				he.AlreadyExists(w)
			}
		} else {
			he.ServiceUnavailable(w)
		}
	}

	return resp, err
}

func (s *Server) GetTaskHandler(w http.ResponseWriter, r *http.Request) (any, error) {
	var body tasks.GetTaskRequest = r.Context().Value("body").(tasks.GetTaskRequest)

	resp, err := s.Tasks.GetTask(r.Context(), body)
	if err != nil {
		slog.Error("handler tasks.GetTask error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				he.NotFound(w)
			}
		} else {
			he.ServiceUnavailable(w)
		}
	}

	return resp, err
}

func (s *Server) GetTasksHandler(w http.ResponseWriter, r *http.Request) (any, error) {
	var body tasks.GetTasksRequest = r.Context().Value("body").(tasks.GetTasksRequest)

	resp, err := s.Tasks.GetTasks(r.Context(), body)
	if err != nil {
		slog.Error("handler tasks.GetTasks error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				he.NotFound(w)
			}
		} else {
			he.ServiceUnavailable(w)
		}
	}

	return resp, err
}

func (s *Server) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) (any, error) {
	var body tasks.UpdateTaskRequest = r.Context().Value("body").(tasks.UpdateTaskRequest)

	resp, err := s.Tasks.UpdateTask(r.Context(), body)
	if err != nil {
		slog.Error("handler tasks.UpdateTask error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				he.NotFound(w)
			}
		} else {
			he.ServiceUnavailable(w)
		}
	}

	return resp, err
}

func (s *Server) ChangeStatusTaskHandler(w http.ResponseWriter, r *http.Request) (any, error) {
	var body tasks.ChangeStatusTaskRequest = r.Context().Value("body").(tasks.ChangeStatusTaskRequest)

	resp, err := s.Tasks.ChangeStatusTask(r.Context(), body)
	if err != nil {
		slog.Error("handler tasks.ChangeStatusTask error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				he.NotFound(w)
			}
		} else {
			he.ServiceUnavailable(w)
		}
	}

	return resp, err
}

func (s *Server) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) (any, error) {
	var body tasks.DeleteTaskRequest = r.Context().Value("body").(tasks.DeleteTaskRequest)

	resp, err := s.Tasks.DeleteTask(r.Context(), body)
	if err != nil {
		slog.Error("handler tasks.DeleteTask error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				he.NotFound(w)
			}
		} else {
			he.ServiceUnavailable(w)
		}
	}

	return resp, err
}
