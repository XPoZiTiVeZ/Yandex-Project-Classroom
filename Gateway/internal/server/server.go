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
	Config  *config.Config
	Server  *http.Server
	Auth    *auth.AuthServiceClient
	Courses *courses.CoursesServiceClient
	Lessons *lessons.LessonsServiceClient
	Tasks   *tasks.TasksServiceClient
	logger  *slog.Logger
}

func (s *Server) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong!"))
}

func (s *Server) RegisterMux(mux *http.ServeMux) {
	mux.HandleFunc("/api/ping", s.Ping)

	// Auth handlers
	if s.Config.Auth.Enabled {
		mux.HandleFunc("/api/auth/register", HandlerWrapper[auth.RegisterRequest](s.RegisterHandler))
		mux.HandleFunc("/api/auth/login", HandlerWrapper[auth.LoginRequest](s.LoginHandler))
		mux.HandleFunc("/api/auth/refresh", HandlerWrapper[auth.RefreshRequest](s.RefreshHandler))
		mux.HandleFunc("/api/auth/logout", s.IsAuthenticated(HandlerWrapper[auth.LogoutRequest](s.LogoutHandler)))
		mux.HandleFunc("/api/auth/user-info", s.IsAuthenticated(HandlerWrapper[auth.GetUserInfoRequest](s.GetUserInfoHandler)))
	}

	// Courses handlers
	// if s.Config.Auth.Enabled && s.Config.Courses.Enabled {
	// 	mux.HandleFunc("/api/courses/create", s.IsAuthenticated(HandlerWrapper[courses.CreateCourseRequest](s.CreateCourseHandler)))
	// 	mux.HandleFunc("/api/courses/course", s.IsMember(HandlerWrapper[courses.GetCourseRequest](s.GetCourseHandler)))
	// 	mux.HandleFunc("/api/courses/courses", s.IsAuthenticated(HandlerWrapper[courses.GetCoursesRequest](s.GetCoursesHandler)))
	// 	mux.HandleFunc("/api/courses/student-courses", s.IsAuthenticated(HandlerWrapper[courses.GetCoursesByStudentRequest](s.GetCoursesByStudentHandler)))
	// 	mux.HandleFunc("/api/courses/teacher-courses", s.IsAuthenticated(HandlerWrapper[courses.GetCoursesByTeacherRequest](s.GetCoursesByTeacherHandler)))
	// 	mux.HandleFunc("/api/courses/course/update", s.IsTeacher(HandlerWrapper[courses.UpdateCourseRequest](s.UpdateCourseHandler)))
	// 	mux.HandleFunc("/api/courses/course/delete", s.IsTeacher(HandlerWrapper[courses.DeleteCourseRequest](s.DeleteCourseHandler)))
	// 	mux.HandleFunc("/api/courses/course/enroll", s.IsTeacher(HandlerWrapper[courses.EnrollUserRequest](s.EnrollUserHandler)))
	// 	mux.HandleFunc("/api/courses/course/expel", s.IsTeacher(HandlerWrapper[courses.ExpelUserRequest](s.ExpelUserHandler)))
	// 	mux.HandleFunc("/api/courses/course/students", s.IsMember(HandlerWrapper[courses.GetCourseStudentsRequest](s.GetCourseStudentsHandler)))
	// }

	// Lessons handlers
	// if s.Config.Auth.Enabled && s.Config.Courses.Enabled && s.Config.Lessons.Enabled {
	// 	mux.HandleFunc("/api/lessons/create", s.IsTeacher(HandlerWrapper[lessons.CreateLessonRequest](s.CreateLessonHandler)))
	// 	mux.HandleFunc("/api/lessons/lesson", s.IsMember(HandlerWrapper[lessons.GetLessonRequest](s.GetLessonHandler)))
	// 	mux.HandleFunc("/api/lessons/lessons", s.IsMember(HandlerWrapper[lessons.GetLessonsRequest](s.GetLessonsHandler)))
	// 	mux.HandleFunc("/api/lessons/lesson/update", s.IsTeacher(HandlerWrapper[lessons.UpdateLessonRequest](s.UpdateLessonHandler)))
	// 	mux.HandleFunc("/api/lessons/lesson/delete", s.IsTeacher(HandlerWrapper[lessons.DeleteLessonRequest](s.DeleteLessonHandler)))
	// }

	// Tasks handlers
	// if s.Config.Auth.Enabled && s.Config.Courses.Enabled && s.Config.Tasks.Enabled {
	// 	mux.HandleFunc("/api/tasks/create", s.IsTeacher(HandlerWrapper[tasks.CreateTaskRequest](s.CreateTaskHandler)))
	// 	mux.HandleFunc("/api/tasks/task", s.IsMember(HandlerWrapper[tasks.GetTaskRequest](s.GetTaskHandler)))
	// 	mux.HandleFunc("/api/tasks/tasks", s.IsMember(HandlerWrapper[tasks.GetTasksRequest](s.GetTasksHandler)))
	// 	mux.HandleFunc("/api/tasks/task/update", s.IsTeacher(HandlerWrapper[tasks.UpdateTaskRequest](s.UpdateTaskHandler)))
	// 	mux.HandleFunc("/api/tasks/task/delete", s.IsTeacher(HandlerWrapper[tasks.DeleteTaskRequest](s.DeleteTaskHandler)))
	// 	mux.HandleFunc("/api/tasks/task/changestatus", s.IsMember(HandlerWrapper[tasks.ChangeStatusTaskRequest](s.ChangeStatusTaskHandler)))
	// }
}

func NewServer(ctx context.Context, cfg *config.Config) (*Server, error) {
	var server Server
	server.Config = cfg

	// TODO: add real logger
	server.logger = slog.Default()

	mux := http.NewServeMux()
	server.RegisterMux(mux)

	server.Server = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", cfg.Host.Port),
		Handler: mux,
	}

	return &server, nil
}

func (s *Server) Run() {
	if s.Config.Auth.Enabled {
		auth, err := auth.NewAuthServiceClient(slog.Default(), s.Config.Auth.Address, s.Config.Auth.Port, nil)
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
}
