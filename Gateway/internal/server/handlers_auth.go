package server

import (
	"Classroom/Gateway/internal/auth"
	he "Classroom/Gateway/internal/errors"
	"log/slog"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) RegisterHandler(w http.ResponseWriter, r *http.Request) (error, any) {
	var body auth.RegisterRequest = r.Context().Value("body").(auth.RegisterRequest)
	
	resp, err := s.Auth.Register(r.Context(), body)
	if err != nil {
		slog.Error("auth.Register error", slog.Any("error", err))

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

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) (error, any) {
	var body auth.LoginRequest = r.Context().Value("body").(auth.LoginRequest)
	
	resp, err := s.Auth.Login(r.Context(), body)
	if err != nil {
		slog.Error("auth.Login error", slog.Any("error", err))

		he.ServiceUnavailable(w)
	}

	return err, resp
}

func (s *Server) RefreshHandler(w http.ResponseWriter, r *http.Request) (error, any) {
	var body auth.RefreshRequest = r.Context().Value("body").(auth.RefreshRequest)

	resp, err := s.Auth.Refresh(r.Context(), body)
	if err != nil {
		slog.Error("auth.Refresh error", slog.Any("error", err))

		he.ServiceUnavailable(w)
	}

	return err, resp
}

func (s *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) (error, any) {
	var body auth.LogoutRequest = r.Context().Value("body").(auth.LogoutRequest)

	resp, err := s.Auth.Logout(r.Context(), body)
	if err != nil {
		slog.Error("auth.Logout error", slog.Any("error", err))

		he.ServiceUnavailable(w)
	}

	return err, resp
}
