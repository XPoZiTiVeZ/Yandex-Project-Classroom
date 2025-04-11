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
	mux.HandleFunc("/ping", s.Ping)

	// Auth handlers
	mux.HandleFunc("/auth/register", HandlerWrapper(s.RegisterHandler))
	mux.HandleFunc("/auth/login", HandlerWrapper(s.LoginHandler))
	mux.HandleFunc("/auth/refresh", HandlerWrapper(s.RefreshHandler))
	mux.HandleFunc("/auth/logout", s.IsAuthenticated(HandlerWrapper(s.LogoutHandler)))

	// Courses handlers
	mux.HandleFunc("/courses/enroll", s.IsTeacher(HandlerWrapper(s.EnrollUserHandler)))
	mux.HandleFunc("/courses/create", s.IsAuthenticated(HandlerWrapper(s.CreateCourseHandler)))
	mux.HandleFunc("/courses/course", s.IsAuthenticated(HandlerWrapper(s.GetCourseHandler)))
	mux.HandleFunc("/courses/courses", s.IsAuthenticated(HandlerWrapper(s.GetCoursesHandler)))
	mux.HandleFunc("/courses/course/update", s.IsTeacher(HandlerWrapper(s.UpdateCourseHandler)))
	mux.HandleFunc("/courses/course/delete", s.IsTeacher(HandlerWrapper(s.DeleteCourseHandler)))

	// Lessons handlers
	mux.HandleFunc("/lessons/create", s.IsTeacher(HandlerWrapper(s.CreateLessonHandler)))
	mux.HandleFunc("/lessons/lesson", s.IsAuthenticated(HandlerWrapper(s.GetLessonHandler)))
	mux.HandleFunc("/lessons/lessons", s.IsAuthenticated(HandlerWrapper(s.GetLessonsHandler)))
	mux.HandleFunc("/lessons/lesson/update", s.IsTeacher(HandlerWrapper(s.UpdateLessonHandler)))
	mux.HandleFunc("/lessons/lesson/delete", s.IsTeacher(HandlerWrapper(s.DeleteLessonHandler)))

	// Tasks handlers
	mux.HandleFunc("/tasks/create", s.IsTeacher(HandlerWrapper(s.CreateTaskHandler)))
	mux.HandleFunc("/tasks/task", s.IsAuthenticated(HandlerWrapper(s.GetTaskHandler)))
	mux.HandleFunc("/tasks/tasks", s.IsAuthenticated(HandlerWrapper(s.GetTasksHandler)))
	mux.HandleFunc("/tasks/task/update", s.IsTeacher(HandlerWrapper(s.UpdateTaskHandler)))
	mux.HandleFunc("/tasks/task/delete", s.IsTeacher(HandlerWrapper(s.DeleteTaskHandler)))
	mux.HandleFunc("/tasks/task/changestatus", s.IsAuthenticated(HandlerWrapper(s.ChangeStatusTaskHandler)))
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
	auth, err := auth.NewAuthServiceClient(s.Config.Auth.Address, s.Config.Auth.Port, nil)
	if err != nil {
		slog.Error("Auth service failed with error", slog.Any("error", err))
		s.CtxStop()
		return
	}
	defer auth.Conn.Close()

	courses, err := courses.NewCoursesServiceClient(s.Config.Courses.Address, s.Config.Courses.Port, nil)
	if err != nil {
		slog.Error("courses service failed with error", slog.Any("error", err))
		s.CtxStop()
		return
	}
	defer courses.Conn.Close()

	lessons, err := lessons.NewLessonsServiceClient(s.Config.Lessons.Address, s.Config.Lessons.Port, nil)
	if err != nil {
		slog.Error("lessons service failed with error", slog.Any("error", err))
		s.CtxStop()
		return
	}

	tasks, err := tasks.NewTasksServiceClient(s.Config.Lessons.Address, s.Config.Lessons.Port, nil)
	if err != nil {
		slog.Error("tasks service failed with error", slog.Any("error", err))
		s.CtxStop()
		return
	}

	s.Auth = auth
	s.Courses = courses
	s.Lessons = lessons
	s.Tasks = tasks

	err = s.Server.ListenAndServe()
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
