package server

import (
	"Classroom/Gateway/internal/auth"
	"Classroom/Gateway/internal/courses"
	"Classroom/Gateway/internal/lessons"
	"Classroom/Gateway/internal/tasks"
	"Classroom/Gateway/pkg/config"
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

type Server struct {
	CtxStop context.CancelFunc
	Config  config.Config
	Server  *http.Server
	Auth    *auth.AuthServiceClient
	Courses *courses.CoursesServiceClient
	Lessons *lessons.LessonsServiceClient
	Tasks   *tasks.TasksServiceClient
}

func (s *Server) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong!"))
}

func (s *Server) RegisterMux(mux *http.ServeMux) {
	mux.HandleFunc("/api/ping", s.Ping)

	// Auth handlers
	if s.Config.Auth.Enabled {
		mux.HandleFunc("/api/auth/register", HandlerWrapper(s.RegisterHandler))
		mux.HandleFunc("/api/auth/login", HandlerWrapper(s.LoginHandler))
		mux.HandleFunc("/api/auth/refresh", HandlerWrapper(s.RefreshHandler))
		mux.HandleFunc("/api/auth/logout", s.IsAuthenticated(HandlerWrapper(s.LogoutHandler)))
	}

	// Courses handlers
	if s.Config.Auth.Enabled && s.Config.Courses.Enabled {
		mux.HandleFunc("/api/courses/enroll", s.IsTeacher(HandlerWrapper(s.EnrollUserHandler)))
		mux.HandleFunc("/api/courses/create", s.IsAuthenticated(HandlerWrapper(s.CreateCourseHandler)))
		mux.HandleFunc("/api/courses/course", s.IsAuthenticated(HandlerWrapper(s.GetCourseHandler)))
		mux.HandleFunc("/api/courses/courses", s.IsAuthenticated(HandlerWrapper(s.GetCoursesHandler)))
		mux.HandleFunc("/api/courses/course/update", s.IsTeacher(HandlerWrapper(s.UpdateCourseHandler)))
		mux.HandleFunc("/api/courses/course/delete", s.IsTeacher(HandlerWrapper(s.DeleteCourseHandler)))
	}

	// Lessons handlers
	if s.Config.Auth.Enabled && s.Config.Courses.Enabled && s.Config.Lessons.Enabled {
		mux.HandleFunc("/api/lessons/create", s.IsTeacher(HandlerWrapper(s.CreateLessonHandler)))
		mux.HandleFunc("/api/lessons/lesson", s.IsAuthenticated(HandlerWrapper(s.GetLessonHandler)))
		mux.HandleFunc("/api/lessons/lessons", s.IsAuthenticated(HandlerWrapper(s.GetLessonsHandler)))
		mux.HandleFunc("/api/lessons/lesson/update", s.IsTeacher(HandlerWrapper(s.UpdateLessonHandler)))
		mux.HandleFunc("/api/lessons/lesson/delete", s.IsTeacher(HandlerWrapper(s.DeleteLessonHandler)))
	}

	// Tasks handlers
	if s.Config.Auth.Enabled && s.Config.Courses.Enabled && s.Config.Tasks.Enabled {
		mux.HandleFunc("/api/tasks/create", s.IsTeacher(HandlerWrapper(s.CreateTaskHandler)))
		mux.HandleFunc("/api/tasks/task", s.IsAuthenticated(HandlerWrapper(s.GetTaskHandler)))
		mux.HandleFunc("/api/tasks/tasks", s.IsAuthenticated(HandlerWrapper(s.GetTasksHandler)))
		mux.HandleFunc("/api/tasks/task/update", s.IsTeacher(HandlerWrapper(s.UpdateTaskHandler)))
		mux.HandleFunc("/api/tasks/task/delete", s.IsTeacher(HandlerWrapper(s.DeleteTaskHandler)))
		mux.HandleFunc("/api/tasks/task/changestatus", s.IsAuthenticated(HandlerWrapper(s.ChangeStatusTaskHandler)))
	}
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
	if s.Config.Auth.Enabled {
		auth, err := auth.NewAuthServiceClient(s.Config.Auth.Address, s.Config.Auth.Port, nil)
		if err != nil {
			slog.Error("Auth service failed with error", slog.Any("error", err))
			s.CtxStop()
			return
		}
		defer auth.Conn.Close()
		s.Auth = auth
	}

	if s.Config.Auth.Enabled && s.Config.Courses.Enabled {
		courses, err := courses.NewCoursesServiceClient(s.Config.Courses.Address, s.Config.Courses.Port, nil)
		if err != nil {
			slog.Error("courses service failed with error", slog.Any("error", err))
			s.CtxStop()
			return
		}
		defer courses.Conn.Close()
		s.Courses = courses
	}

	if s.Config.Auth.Enabled && s.Config.Courses.Enabled && s.Config.Lessons.Enabled {
		lessons, err := lessons.NewLessonsServiceClient(s.Config.Lessons.Address, s.Config.Lessons.Port, nil)
		if err != nil {
			slog.Error("lessons service failed with error", slog.Any("error", err))
			s.CtxStop()
			return
		}
		defer lessons.Conn.Close()
		s.Lessons = lessons
	}

	if s.Config.Auth.Enabled && s.Config.Courses.Enabled && s.Config.Tasks.Enabled {
		tasks, err := tasks.NewTasksServiceClient(s.Config.Tasks.Address, s.Config.Tasks.Port, nil)
		if err != nil {
			slog.Error("tasks service failed with error", slog.Any("error", err))
			s.CtxStop()
			return
		}
		defer tasks.Conn.Close()
		s.Tasks = tasks
	}

	err := s.Server.ListenAndServe()
	if err == http.ErrServerClosed {
		slog.Info("HTTP Server closed")
		return
	}

	if err != nil {
		slog.Error("Server failed with error", slog.Any("error", err))
		s.CtxStop()
	}

	return
}
