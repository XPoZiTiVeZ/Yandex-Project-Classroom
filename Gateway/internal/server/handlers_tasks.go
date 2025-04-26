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
// @Failure 401 {object} ErrorResponse "Требуется авторизация"
// @Failure 403 {object} ErrorResponse "Доступ запрещён"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /tasks/create [post]
func (s *Server) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[tasks.CreateTaskRequest](r.Context())

	isTeacher, err := s.IsTeacher(r.Context(), body.CourseID)
	if err != nil {
		logger.Error(r.Context(), "Handler courses.IsTeacher error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				BadRequest(w, e.Message())
			case codes.NotFound:
				NotFound(w, e.Message())
			case codes.Unavailable:
				ServiceUnavailable(w)
			}
		} else {
			InternalError(w)
		}
		return
	}

	if !isTeacher {
		Forbidden(w)
		return
	}

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
		return
	}

	WriteJSON(w, resp, http.StatusCreated)
}

// GetTaskHandler возвращает информацию о задаче
// @Summary Получение задачи
// @Description Возвращает детальную информацию о задаче в рамках курса
// @Tags Tasks
// @Produce json
// @Security BearerAuth
// @Param task_id query string true "ID задачи" example("5a430d16-851d-45a9-b55b-15838785adea")
// @Success 200 {object} tasks.GetTaskResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 401 {object} ErrorResponse "Требуется авторизация"
// @Failure 403 {object} ErrorResponse "Доступ запрещен"
// @Failure 404 {object} ErrorResponse "Задача не найдена"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router  /tasks/task [get]
func (s *Server) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[tasks.GetTaskRequest](r.Context())

	resp, err := s.Tasks.GetTask(r.Context(), body)

	isMember, err := s.IsMember(r.Context(), resp.Task.CourseID)
	if err != nil {
		logger.Error(r.Context(), "Handler courses.IsMember error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				BadRequest(w, e.Message())
			case codes.NotFound:
				NotFound(w, e.Message())
			case codes.Unavailable:
				ServiceUnavailable(w)
			}
		} else {
			InternalError(w)
		}
		return
	}

	if !isMember {
		Forbidden(w)
		return
	}

	WriteJSON(w, resp, http.StatusOK)
}

// GetStudentStatuses возвращает статусы студентов для задачи
// @Summary Получение статусов студентов
// @Description Возвращает статусы выполнения задачи для студентов
// @Tags Tasks
// @Produce json
// @Security BearerAuth
// @Param task_id query string true "ID задачи" example("21dad0c3-dcea-4c19-b501-fb2fe888f683")
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

	body1 := tasks.GetTaskRequest{
		TaskID: body.TaskID,
	}
	resp1, err := s.Tasks.GetTask(r.Context(), body1)
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
		return
	}

	isTeacher, err := s.IsTeacher(r.Context(), resp1.Task.CourseID)
	if err != nil {
		logger.Error(r.Context(), "Handler courses.IsTeacher error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				BadRequest(w, e.Message())
			case codes.NotFound:
				NotFound(w, e.Message())
			case codes.Unavailable:
				ServiceUnavailable(w)
			}
		} else {
			InternalError(w)
		}
		return
	}

	if !isTeacher {
		Forbidden(w)
		return
	}

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
		return
	}

	WriteJSON(w, resp, http.StatusOK)
}

// GetTasksHandler возвращает задачи для учителя
// @Summary Получение списка задач для учителя
// @Description Возвращает список задач курса для учителя
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param course_id query string true "Идентификатор курса" example(44e7f029-82cc-46f5-83e8-34b7d056ce32)
// @Success 200 {object} tasks.GetTasksResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 401 {object} ErrorResponse "Требуется авторизация"
// @Failure 404 {object} ErrorResponse "Курс не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /tasks/teacher-tasks [get]
func (s *Server) GetTasksForTeacherHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[tasks.GetTasksRequest](r.Context())

	isTeacher, err := s.IsMember(r.Context(), body.CourseID)
	if err != nil {
		logger.Error(r.Context(), "Handler courses.IsTeacher error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				BadRequest(w, e.Message())
			case codes.NotFound:
				NotFound(w, e.Message())
			case codes.Unavailable:
				ServiceUnavailable(w)
			}
		} else {
			InternalError(w)
		}
		return
	}

	if !isTeacher {
		Forbidden(w)
		return
	}

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
		return
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
// @Param course_id query string true "Идентификатор курса" example(44e7f029-82cc-46f5-83e8-34b7d056ce32)
// @Success 200 {object} tasks.GetTasksForStudentResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 401 {object} ErrorResponse "Требуется авторизация"
// @Failure 404 {object} ErrorResponse "Студент не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /tasks/student-tasks [get]
func (s *Server) GetTasksForStudentHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[tasks.GetTasksForStudentRequest](r.Context())
	claims, _ := GetClaims(r.Context())
	body.StudentID = claims.UserID

	isStudent, err := s.IsStudent(r.Context(), body.CourseID)
	if err != nil {
		logger.Error(r.Context(), "Handler courses.IsStudent error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				BadRequest(w, e.Message())
			case codes.NotFound:
				NotFound(w, e.Message())
			case codes.Unavailable:
				ServiceUnavailable(w)
			}
		} else {
			InternalError(w)
		}
		return
	}

	if !isStudent {
		Forbidden(w)
		return
	}

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
		return
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
// @Failure 401 {object} ErrorResponse "Требуется авторизация"
// @Failure 404 {object} ErrorResponse "Задача не найдена"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /tasks/task/update [put]
func (s *Server) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[tasks.UpdateTaskRequest](r.Context())

	body1 := tasks.GetTaskRequest{
		TaskID: body.TaskID,
	}
	resp1, err := s.Tasks.GetTask(r.Context(), body1)
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
		return
	}

	isTeacher, err := s.IsTeacher(r.Context(), resp1.Task.CourseID)
	if err != nil {
		logger.Error(r.Context(), "Handler courses.IsTeacher error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				BadRequest(w, e.Message())
			case codes.NotFound:
				NotFound(w, e.Message())
			case codes.Unavailable:
				ServiceUnavailable(w)
			}
		} else {
			InternalError(w)
		}
		return
	}

	if !isTeacher {
		Forbidden(w)
		return
	}

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
		return
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
// @Failure 401 {object} ErrorResponse "Требуется авторизация"
// @Failure 404 {object} ErrorResponse "Задача не найдена"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /tasks/task/changestatus [patch]
func (s *Server) ChangeStatusTaskHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[tasks.ChangeStatusTaskRequest](r.Context())

	body1 := tasks.GetTaskRequest{
		TaskID: body.TaskID,
	}
	resp1, err := s.Tasks.GetTask(r.Context(), body1)
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
		return
	}

	isStudent, err := s.IsStudent(r.Context(), resp1.Task.CourseID)
	if err != nil {
		logger.Error(r.Context(), "Handler courses.IsStudent error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				BadRequest(w, e.Message())
			case codes.NotFound:
				NotFound(w, e.Message())
			case codes.Unavailable:
				ServiceUnavailable(w)
			}
		} else {
			InternalError(w)
		}
		return
	}

	if !isStudent {
		Forbidden(w)
		return
	}

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
		return
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
// @Failure 401 {object} ErrorResponse "Требуется авторизация"
// @Failure 404 {object} ErrorResponse "Задача не найдена"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /tasks/task/delete [delete]
func (s *Server) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[tasks.DeleteTaskRequest](r.Context())

	body1 := tasks.GetTaskRequest{
		TaskID: body.TaskID,
	}
	resp1, err := s.Tasks.GetTask(r.Context(), body1)
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
		return
	}

	isTeacher, err := s.IsTeacher(r.Context(), resp1.Task.CourseID)
	if err != nil {
		logger.Error(r.Context(), "Handler courses.IsTeacher error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				BadRequest(w, e.Message())
			case codes.NotFound:
				NotFound(w, e.Message())
			case codes.Unavailable:
				ServiceUnavailable(w)
			}
		} else {
			InternalError(w)
		}
		return
	}

	if !isTeacher {
		Forbidden(w)
		return
	}

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
		return
	}

	WriteJSON(w, resp, http.StatusOK)
}
