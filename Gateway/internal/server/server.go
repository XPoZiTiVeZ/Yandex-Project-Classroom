package server

import (
	"Classroom/Gateway/internal/auth"
	"Classroom/Gateway/internal/courses"
	"Classroom/Gateway/internal/lessons"
	"Classroom/Gateway/pkg/config"
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

type Server struct {
	CtxStop context.CancelFunc
	Config  *config.Config
	Server  *http.Server
	Auth    *auth.AuthServiceClient
	Courses *courses.CourseServiceClient
	Lessons *lessons.LessonServiceClient
}

func (s *Server) RegisterMux(mux *http.ServeMux) {
	mux.HandleFunc("/api/ping", s.Ping)
	mux.HandleFunc("/api/auth/register", s.RegisterHandler)
	mux.HandleFunc("/api/auth/register", s.LoginHandler)
}

func NewServer(address string, port int, ctx context.Context) (*Server, error) {
	var server Server

	mux := http.NewServeMux()

	server.RegisterMux(mux)

	server.Server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", address, port),
		Handler: mux,
	}

	return &server, nil
}

func (s *Server) Run() {
	if s.Config == nil {
		err := fmt.Errorf("config is nil")
		slog.Error("no server config provided", slog.Any("error", err))
		s.CtxStop()
		return
	}

	auth, err := auth.NewAuthServiceClient(s.Config.Auth.Address, s.Config.Auth.Port)
	if err != nil {
		slog.Error("auth service failed with error", slog.Any("error", err))
		s.CtxStop()
		return
	}
	courses, err := courses.NewCoursesServiceClient(s.Config.Courses.Address, s.Config.Courses.Port)
	if err != nil {
		slog.Error("courses service failed with error", slog.Any("error", err))
		s.CtxStop()
		return
	}
	lessons, err := lessons.NewLessonsServiceClient(s.Config.Lessons.Address, s.Config.Lessons.Port)
	if err != nil {
		slog.Error("lessons service failed with error", slog.Any("error", err))
		s.CtxStop()
		return
	}

	s.Auth = auth
	s.Courses = courses
	s.Lessons = lessons

	err = s.Server.ListenAndServe()
	if err != nil {
		slog.Error("server failed with error", slog.Any("error", err))
		s.CtxStop()
	}

	return
}
