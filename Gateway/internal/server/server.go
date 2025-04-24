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

	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	CtxStop context.CancelFunc
	Config  *config.Config
	Server  *http.Server
	Redis   *redis.Client
	Auth    *auth.AuthServiceClient
	Courses *courses.CoursesServiceClient
	Lessons *lessons.LessonsServiceClient
	Tasks   *tasks.TasksServiceClient
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем запросы с любого origin (можно указать конкретные домены)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
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
	// Статус код ответа сервера
	Code int `json:"code" example:"200" extensions:"x-orders=0"`
	// Сообщение-ответ сервера
	Message string `json:"message" example:"Pong!" extensions:"x-order=1"`
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
		Code: 200,
		Message: "Pong!",
	}

	WriteJSON(w, pong, http.StatusOK)
}

func (s *Server) RegisterMux(mux *http.ServeMux) {
	mux.Handle("/api/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://127.0.0.1/files/swagger.json"),
	))

	mux.HandleFunc("GET /api/ping", s.Ping)

	// Auth handlers
	if s.Config.Auth.Enabled {

		mux.HandleFunc("POST /api/auth/register", JSONHandlerWrapper[auth.RegisterRequest](s.RegisterHandler))
		mux.HandleFunc("POST /api/auth/login", JSONHandlerWrapper[auth.LoginRequest](s.LoginHandler))
		mux.HandleFunc("POST /api/auth/refresh", JSONHandlerWrapper[auth.RefreshRequest](s.RefreshHandler))
		mux.HandleFunc("POST /api/auth/logout", s.IsAuthenticated(JSONHandlerWrapper[auth.LogoutRequest](s.LogoutHandler)))
		mux.HandleFunc("GET /api/auth/user-info", s.IsAuthenticated(QueryHandlerWrapper[auth.GetUserInfoRequest](s.GetUserInfoHandler)))
	}

	// Courses handlers
	if s.Config.Auth.Enabled && s.Config.Courses.Enabled {
		mux.HandleFunc("POST /api/courses/create", s.IsAuthenticated(JSONHandlerWrapper[courses.CreateCourseRequest](s.CreateCourseHandler)))
		mux.HandleFunc("GET /api/courses/course", s.IsAuthenticated(QueryHandlerWrapper[courses.GetCourseRequest](s.GetCourseHandler)))
		mux.HandleFunc("GET /api/courses/courses", s.IsAuthenticated(QueryHandlerWrapper[courses.GetCoursesRequest](s.GetCoursesHandler)))
		mux.HandleFunc("GET /api/courses/student-courses", s.IsAuthenticated(QueryHandlerWrapper[courses.GetCoursesByStudentRequest](s.GetCoursesByStudentHandler)))
		mux.HandleFunc("GET /api/courses/teacher-courses", s.IsAuthenticated(QueryHandlerWrapper[courses.GetCoursesByTeacherRequest](s.GetCoursesByTeacherHandler)))
		mux.HandleFunc("PUT /api/courses/course/update", s.IsAuthenticated(JSONHandlerWrapper[courses.UpdateCourseRequest](s.UpdateCourseHandler)))
		mux.HandleFunc("DELETE /api/courses/course/delete", s.IsAuthenticated(JSONHandlerWrapper[courses.DeleteCourseRequest](s.DeleteCourseHandler)))
		mux.HandleFunc("PUT /api/courses/course/enroll", s.IsAuthenticated(JSONHandlerWrapper[courses.EnrollUserRequest](s.EnrollUserHandler)))
		mux.HandleFunc("PUT /api/courses/course/expel", s.IsAuthenticated(JSONHandlerWrapper[courses.ExpelUserRequest](s.ExpelUserHandler)))
		mux.HandleFunc("GET /api/courses/course/students", s.IsAuthenticated(QueryHandlerWrapper[courses.GetCourseStudentsRequest](s.GetCourseStudentsHandler)))
	}

	// Lessons handlers
	if s.Config.Auth.Enabled && s.Config.Courses.Enabled && s.Config.Lessons.Enabled {
		mux.HandleFunc("POST /api/lessons/create", s.IsAuthenticated(JSONHandlerWrapper[lessons.CreateLessonRequest](s.CreateLessonHandler)))
		mux.HandleFunc("POST /api/lessons/lesson", s.IsAuthenticated(JSONHandlerWrapper[lessons.GetLessonRequest](s.GetLessonHandler)))
		mux.HandleFunc("POST /api/lessons/lessons", s.IsAuthenticated(JSONHandlerWrapper[lessons.GetLessonsRequest](s.GetLessonsHandler)))
		mux.HandleFunc("PUT /api/lessons/lesson/update", s.IsAuthenticated(JSONHandlerWrapper[lessons.UpdateLessonRequest](s.UpdateLessonHandler)))
		mux.HandleFunc("DELETE /api/lessons/lesson/delete", s.IsAuthenticated(JSONHandlerWrapper[lessons.DeleteLessonRequest](s.DeleteLessonHandler)))
	}

	// Tasks handlers
	if s.Config.Auth.Enabled && s.Config.Courses.Enabled && s.Config.Tasks.Enabled {
		mux.HandleFunc("POST /api/tasks/create", JSONHandlerWrapper[tasks.CreateTaskRequest](s.CreateTaskHandler))
		mux.HandleFunc("GET /api/tasks/task", QueryHandlerWrapper[tasks.GetTaskRequest](s.GetTaskHandler))
		mux.HandleFunc("GET /api/tasks/student-tasks", QueryHandlerWrapper[tasks.GetTasksRequest](s.GetTasksForStudentHandler))
		mux.HandleFunc("GET /api/tasks/teacher-tasks", QueryHandlerWrapper[tasks.GetTasksRequest](s.GetTasksForTeacherHandler))
		mux.HandleFunc("GET /api/tasks/tasks-statuses", JSONHandlerWrapper[tasks.GetTasksRequest](s.GetStudentStatuses))
		mux.HandleFunc("PUT /api/tasks/task/update", JSONHandlerWrapper[tasks.UpdateTaskRequest](s.UpdateTaskHandler))
		mux.HandleFunc("DELETE /api/tasks/task/delete", JSONHandlerWrapper[tasks.DeleteTaskRequest](s.DeleteTaskHandler))
		mux.HandleFunc("PATCH /api/tasks/task/changestatus", JSONHandlerWrapper[tasks.ChangeStatusTaskRequest](s.ChangeStatusTaskHandler))
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
