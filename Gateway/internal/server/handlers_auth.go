package server

import (
	"Classroom/Gateway/internal/auth"
	he "Classroom/Gateway/internal/errors"
	"Classroom/Gateway/pkg/logger"
	"log/slog"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[auth.RegisterRequest](r.Context())

	resp, err := s.Auth.Register(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler auth.Register error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				BadRequest(w, "invalid arguments")
			case codes.AlreadyExists:
				AlreadyExists(w, "email already exists")
			case codes.Unavailable:
				ServiceUnavailable(w)
			}
		} else {
			InternalError(w)
		}
		return
	}

	WriteJSON(w, resp, http.StatusCreated)
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[auth.LoginRequest](r.Context())

	resp, err := s.Auth.Login(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler auth.Login error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				BadRequest(w, "invalid arguments")
			case codes.Unauthenticated:
				Unauthorized(w, "invalid credentials")
			case codes.Unavailable:
				ServiceUnavailable(w)
			}
		} else {
			InternalError(w)
		}
		return
	}

	WriteJSON(w, resp, http.StatusOK)
}

func (s *Server) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[auth.RefreshRequest](r.Context())

	resp, err := s.Auth.Refresh(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler auth.Refresh error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Unauthenticated:
				Unauthorized(w, "invalid refresh token")
			case codes.Unavailable:
				ServiceUnavailable(w)
			}
		} else {
			InternalError(w)
		}
		return
	}

	WriteJSON(w, resp, http.StatusOK)
}

func (s *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[auth.LogoutRequest](r.Context())

	resp, err := s.Auth.Logout(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler auth.Logout error", slog.Any("error", err))

		he.ServiceUnavailable(w)
	}
	WriteJSON(w, resp, http.StatusOK)
}

func (s *Server) GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[auth.GetUserInfoRequest](r.Context())

	resp, err := s.Auth.GetUserInfo(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler auth.GetUserInfo error", slog.Any("error", err))

		he.ServiceUnavailable(w)
	}

	WriteJSON(w, resp, http.StatusOK)
}
