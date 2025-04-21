// +swaggo

package server

import (
	"Classroom/Gateway/internal/auth"
	"Classroom/Gateway/internal/courses"
	"Classroom/Gateway/internal/lessons"
	"Classroom/Gateway/internal/tasks"
	"Classroom/Gateway/pkg/config"
	"Classroom/Gateway/pkg/logger"
	"context"
	"fmt"
	"log/slog"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	CtxStop context.CancelFunc
	Config  *config.Config
	Server  *http.Server
	Auth    *auth.AuthServiceClient
	Courses *courses.CoursesServiceClient
	Lessons *lessons.LessonsServiceClient
	Tasks   *tasks.TasksServiceClient
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем запросы с любого origin (можно указать конкретные домены)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Пропускаем OPTIONS запросы (preflight)
		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Pong - структура ответа для проверки работоспособности API
// @Description Используется для health-check и проверки доступности сервера
type Pong struct {
	// Сообщение-ответ сервера
	Message string `json:"msg" example:"Pong!" extensions:"x-order=0"`
	// HTTP статус код ответа
	Status int `json:"status" example:"200" extensions:"x-order=1"`
} // @name Pong

// Ping обрабатывает запрос проверки работоспособности сервера
// @Summary Проверка работоспособности
// @Description Возвращает ответ "Pong!" для проверки доступности сервера
// @Tags Service
// @Produce json
// @Success 200 {object} server.Pong "Успешный ответ"
// @Router /ping [get]
func (s *Server) Ping(w http.ResponseWriter, r *http.Request) {
	pong := Pong{
		Message: "Pong!",
		Status:  200,
	}

	WriteJSON(w, pong, http.StatusOK)
}

func (s *Server) RegisterMux(mux *http.ServeMux) {
	mux.Handle("/api/swagger/", httpSwagger.Handler(
		// TODO: можно вынести в конфиг
		httpSwagger.URL("http://localhost/files/swagger.json"),
	))

	mux.HandleFunc("GET /api/ping", s.Ping)

	// Auth handlers
	if s.Config.Auth.Enabled {

		mux.HandleFunc("POST /api/auth/register", HandlerWrapper[auth.RegisterRequest](s.RegisterHandler))
		mux.HandleFunc("POST /api/auth/login", HandlerWrapper[auth.LoginRequest](s.LoginHandler))
		mux.HandleFunc("POST /api/auth/refresh", HandlerWrapper[auth.RefreshRequest](s.RefreshHandler))
		mux.HandleFunc("POST /api/auth/logout", s.IsAuthenticated(HandlerWrapper[auth.LogoutRequest](s.LogoutHandler)))
		mux.HandleFunc("POST /api/auth/user-info", s.IsAuthenticated(HandlerWrapper[auth.GetUserInfoRequest](s.GetUserInfoHandler)))
	}

	// Courses handlers
	if s.Config.Auth.Enabled && s.Config.Courses.Enabled {
		mux.HandleFunc("POST /api/courses/create", s.IsAuthenticated(HandlerWrapper[courses.CreateCourseRequest](s.CreateCourseHandler)))
		mux.HandleFunc("POST /api/courses/course", s.IsMember(HandlerWrapper[courses.GetCourseRequest](s.GetCourseHandler)))
		mux.HandleFunc("POST /api/courses/courses", s.IsAuthenticated(HandlerWrapper[courses.GetCoursesRequest](s.GetCoursesHandler)))
		mux.HandleFunc("POST /api/courses/student-courses", s.IsAuthenticated(HandlerWrapper[courses.GetCoursesByStudentRequest](s.GetCoursesByStudentHandler)))
		mux.HandleFunc("POST /api/courses/teacher-courses", s.IsAuthenticated(HandlerWrapper[courses.GetCoursesByTeacherRequest](s.GetCoursesByTeacherHandler)))
		mux.HandleFunc("PUT /api/courses/course/update", s.IsTeacher(HandlerWrapper[courses.UpdateCourseRequest](s.UpdateCourseHandler)))
		mux.HandleFunc("DELETE /api/courses/course/delete", s.IsTeacher(HandlerWrapper[courses.DeleteCourseRequest](s.DeleteCourseHandler)))
		mux.HandleFunc("POST /api/courses/course/enroll", s.IsTeacher(HandlerWrapper[courses.EnrollUserRequest](s.EnrollUserHandler)))
		mux.HandleFunc("POST /api/courses/course/expel", s.IsTeacher(HandlerWrapper[courses.ExpelUserRequest](s.ExpelUserHandler)))
		mux.HandleFunc("POST /api/courses/course/students", s.IsMember(HandlerWrapper[courses.GetCourseStudentsRequest](s.GetCourseStudentsHandler)))
	}

	// Lessons handlers
	if s.Config.Auth.Enabled && s.Config.Courses.Enabled && s.Config.Lessons.Enabled {
		mux.HandleFunc("POST /api/lessons/create", s.IsTeacher(HandlerWrapper[lessons.CreateLessonRequest](s.CreateLessonHandler)))
		mux.HandleFunc("POST /api/lessons/lesson", s.IsMember(HandlerWrapper[lessons.GetLessonRequest](s.GetLessonHandler)))
		mux.HandleFunc("POST /api/lessons/lessons", s.IsMember(HandlerWrapper[lessons.GetLessonsRequest](s.GetLessonsHandler)))
		mux.HandleFunc("PUT /api/lessons/lesson/update", s.IsTeacher(HandlerWrapper[lessons.UpdateLessonRequest](s.UpdateLessonHandler)))
		mux.HandleFunc("DELET /api/lessons/lesson/delete", s.IsTeacher(HandlerWrapper[lessons.DeleteLessonRequest](s.DeleteLessonHandler)))
	}

	// Tasks handlers
	if s.Config.Auth.Enabled && s.Config.Courses.Enabled && s.Config.Tasks.Enabled {
		mux.HandleFunc("POST /api/tasks/create", s.IsTeacher(HandlerWrapper[tasks.CreateTaskRequest](s.CreateTaskHandler)))
		mux.HandleFunc("POST /api/tasks/task", s.IsMember(HandlerWrapper[tasks.GetTaskRequest](s.GetTaskHandler)))
		mux.HandleFunc("POST /api/tasks/student-tasks", s.IsStudent(HandlerWrapper[tasks.GetTasksRequest](s.GetTasksHandler)))
		mux.HandleFunc("POST /api/tasks/teacher-tasks", s.IsTeacher(HandlerWrapper[tasks.GetTasksRequest](s.GetTasksHandler)))
		mux.HandleFunc("POST /api/tasks/tasks-statuses", s.IsTeacher(HandlerWrapper[tasks.GetTasksRequest](s.GetTasksHandler)))
		mux.HandleFunc("PUT /api/tasks/task/update", s.IsTeacher(HandlerWrapper[tasks.UpdateTaskRequest](s.UpdateTaskHandler)))
		mux.HandleFunc("DELETE /api/tasks/task/delete", s.IsTeacher(HandlerWrapper[tasks.DeleteTaskRequest](s.DeleteTaskHandler)))
		mux.HandleFunc("PUT /api/tasks/task/changestatus", s.IsStudent(HandlerWrapper[tasks.ChangeStatusTaskRequest](s.ChangeStatusTaskHandler)))
	}
}

func NewServer(ctx context.Context, cfg *config.Config) (*Server, error) {
	var server Server

	server.Config = cfg

	mux := http.NewServeMux()
	server.RegisterMux(mux)

	server.Server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Host.Address, cfg.Host.Port),
		Handler: enableCORS(mux),
	}

	return &server, nil
}

func (s *Server) Run(ctx context.Context) {
	if s.Config.Auth.Enabled {
		auth, err := auth.NewAuthServiceClient(ctx, s.Config)
		if err != nil {
			logger.Error(ctx, "Auth service failed with error", slog.Any("error", err))
			s.CtxStop()
			return
		}
		defer auth.Conn.Close()
		s.Auth = auth
	}

	if s.Config.Auth.Enabled && s.Config.Courses.Enabled {
		courses, err := courses.NewCoursesServiceClient(ctx, s.Config)
		if err != nil {
			logger.Error(ctx, "Courses service failed with error", slog.Any("error", err))
			s.CtxStop()
			return
		}
		defer courses.Conn.Close()
		s.Courses = courses
	}

	if s.Config.Auth.Enabled && s.Config.Courses.Enabled && s.Config.Lessons.Enabled {
		lessons, err := lessons.NewLessonsServiceClient(ctx, s.Config)
		if err != nil {
			logger.Error(ctx, "Lessons service failed with error", slog.Any("error", err))
			s.CtxStop()
			return
		}
		defer lessons.Conn.Close()
		s.Lessons = lessons
	}

	if s.Config.Auth.Enabled && s.Config.Courses.Enabled && s.Config.Tasks.Enabled {
		tasks, err := tasks.NewTasksServiceClient(ctx, s.Config)
		if err != nil {
			logger.Error(ctx, "Tasks service failed with error", slog.Any("error", err))
			s.CtxStop()
			return
		}
		defer tasks.Conn.Close()
		s.Tasks = tasks
	}

	err := s.Server.ListenAndServe()
	if err == http.ErrServerClosed {
		logger.Info(ctx, "HTTP Server closed")
		return
	}

	if err != nil {
		logger.Error(ctx, "Server failed with error", slog.Any("error", err))
		s.CtxStop()
	}
}
