package server

import (
	"Classroom/Gateway/internal/tasks"
	"Classroom/Gateway/pkg/logger"
	"log/slog"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateTaskHandler создает новую задачу
// @Summary Создание задачи
// @Description Создает новую задачу в курсе
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body tasks.CreateTaskRequest true "Данные для создания задачи"
// @Success 201 {object} tasks.CreateTaskResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /tasks/create [post]
func (s *Server) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[tasks.CreateTaskRequest](r.Context())

	resp, err := s.Tasks.CreateTask(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler tasks.CreateTask error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				BadRequest(w, e.Message())
			case codes.Unavailable:
				ServiceUnavailable(w)
			}
		} else {
			InternalError(w)
		}
	}

	WriteJSON(w, resp, http.StatusCreated)
}

// GetTaskHandler возвращает информацию о задаче
// @Summary Получение задачи
// @Description Возвращает детальную информацию о задаче в рамках курса
// @Tags Tasks
// @Produce json
// @Security BearerAuth
// @Param course_id path string true "ID курса" example("68b75705-9722-466e-96a7-0afc3b5ef22f")
// @Param task_id path string true "ID задачи" example("5a430d16-851d-45a9-b55b-15838785adea")
// @Success 200 {object} tasks.GetTaskResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 401 {object} ErrorResponse "Требуется авторизация"
// @Failure 403 {object} ErrorResponse "Доступ запрещен"
// @Failure 404 {object} ErrorResponse "Задача не найдены"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router  /tasks/task [get]
func (s *Server) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[tasks.GetTaskRequest](r.Context())

	resp, err := s.Tasks.GetTask(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler tasks.GetTask error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				BadRequest(w, e.Message())
			case codes.NotFound:
				NotFound(w)
			case codes.Unavailable:
				ServiceUnavailable(w)
			}
		} else {
			InternalError(w)
		}
	}

	WriteJSON(w, resp, http.StatusOK)
}

// GetStudentStatuses возвращает статусы студентов для задачи
// @Summary Получение статусов студентов
// @Description Возвращает статусы выполнения задачи для студентов
// @Tags Tasks
// @Produce json
// @Security BearerAuth
// @Param course_id path string true "ID курса" example("29f14724-aa33-48d1-872b-fdaddb212e40")
// @Param task_id path string true "ID задачи" example("21dad0c3-dcea-4c19-b501-fb2fe888f683")
// @Success 200 {object} tasks.GetStudentStatusesResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 401 {object} ErrorResponse "Требуется авторизация"
// @Failure 403 {object} ErrorResponse "Доступ запрещен"
// @Failure 404 {object} ErrorResponse "Данные не найдены"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router  /tasks/student-statuses [get]
func (s *Server) GetStudentStatuses(w http.ResponseWriter, r *http.Request) {
	body := GetBody[tasks.GetStudentStatusesRequest](r.Context())

	resp, err := s.Tasks.GetStudentStatuses(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler tasks.GetStudentStatuses error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				BadRequest(w, e.Message())
			case codes.NotFound:
				NotFound(w)
			case codes.Unavailable:
				ServiceUnavailable(w)
			}
		} else {
			InternalError(w)
		}
	}

	WriteJSON(w, resp, http.StatusOK)
}

// GetTasksHandler возвращает список задач
// @Summary Получение списка задач
// @Description Возвращает список задач курса
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body tasks.GetTasksRequest true "Идентификатор курса"
// @Success 200 {object} tasks.GetTasksResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 404 {object} ErrorResponse "Курс не найден"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /tasks/tasks [post]
func (s *Server) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[tasks.GetTasksRequest](r.Context())

	resp, err := s.Tasks.GetTasks(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler tasks.GetTasks error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				BadRequest(w, e.Message())
			case codes.NotFound:
				NotFound(w)
			case codes.Unavailable:
				ServiceUnavailable(w)
			}
		} else {
			InternalError(w)
		}
	}

	WriteJSON(w, resp, http.StatusOK)
}

// GetTasksForStudentHandler возвращает задачи для студента
// @Summary Получение задач студента
// @Description Возвращает список задач для конкретного студента
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body tasks.GetTasksForStudentRequest true "Идентификатор студента"
// @Success 200 {object} tasks.GetTasksForStudentResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 404 {object} ErrorResponse "Студент не найден"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /tasks/student-tasks [post]
func (s *Server) GetTasksForStudentHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[tasks.GetTasksForStudentRequest](r.Context())

	resp, err := s.Tasks.GetTasksForStudent(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler tasks.GetTasksForStudent error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				BadRequest(w, e.Message())
			case codes.NotFound:
				NotFound(w)
			case codes.Unavailable:
				ServiceUnavailable(w)
			}
		} else {
			InternalError(w)
		}
	}

	WriteJSON(w, resp, http.StatusOK)
}

// UpdateTaskHandler обновляет информацию о задаче
// @Summary Обновление задачи
// @Description Обновляет информацию о задаче
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body tasks.UpdateTaskRequest true "Данные для обновления"
// @Success 200 {object} tasks.UpdateTaskResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 404 {object} ErrorResponse "Задача не найдена"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /tasks/task/update [put]
func (s *Server) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[tasks.UpdateTaskRequest](r.Context())

	resp, err := s.Tasks.UpdateTask(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler tasks.UpdateTask error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				BadRequest(w, e.Message())
			case codes.NotFound:
				NotFound(w)
			case codes.Unavailable:
				ServiceUnavailable(w)
			}
		} else {
			InternalError(w)
		}
	}

	WriteJSON(w, resp, http.StatusOK)
}

// ChangeStatusTaskHandler изменяет статус задачи
// @Summary Изменение статуса задачи
// @Description Обновляет статус выполнения задачи студентом
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body tasks.ChangeStatusTaskRequest true "Данные для изменения статуса"
// @Success 200 {object} tasks.ChangeStatusTaskResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 404 {object} ErrorResponse "Задача не найдена"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /tasks/task/status [patch]
func (s *Server) ChangeStatusTaskHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[tasks.ChangeStatusTaskRequest](r.Context())

	resp, err := s.Tasks.ChangeStatusTask(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler tasks.ChangeStatusTask error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				BadRequest(w, e.Message())
			case codes.NotFound:
				NotFound(w)
			case codes.Unavailable:
				ServiceUnavailable(w)
			}
		} else {
			InternalError(w)
		}
	}

	WriteJSON(w, resp, http.StatusOK)
}

// DeleteTaskHandler удаляет задачу
// @Summary Удаление задачи
// @Description Удаляет задачу по идентификатору
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body tasks.DeleteTaskRequest true "Идентификатор задачи"
// @Success 200 {object} tasks.DeleteTaskResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 404 {object} ErrorResponse "Задача не найдена"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /tasks/task/delete [delete]
func (s *Server) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[tasks.DeleteTaskRequest](r.Context())

	resp, err := s.Tasks.DeleteTask(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler tasks.DeleteTask error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				BadRequest(w, e.Message())
			case codes.NotFound:
				NotFound(w)
			case codes.Unavailable:
				ServiceUnavailable(w)
			}
		} else {
			InternalError(w)
		}
	}

	WriteJSON(w, resp, http.StatusOK)
}
