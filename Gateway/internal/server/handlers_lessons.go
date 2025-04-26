package server

import (
	"Classroom/Gateway/internal/lessons"
	"Classroom/Gateway/internal/redis"
	"Classroom/Gateway/pkg/logger"
	"log/slog"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateLessonHandler создает новый урок
// @Summary Создание урока
// @Description Создает новый урок в курсе
// @Tags Lessons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body lessons.CreateLessonRequest true "Данные для создания урока"
// @Success 201 {object} lessons.CreateLessonResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 401 {object} ErrorResponse "Требуется авторизация"
// @Failure 403 {object} ErrorResponse "Доступ запрещён"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /lessons/create [post]
func (s *Server) CreateLessonHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[lessons.CreateLessonRequest](r.Context())

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

	resp, err := s.Lessons.CreateLesson(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler lessons.CreateLesson error", slog.Any("error", err))

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

// GetLessonHandler возвращает информацию об уроке
// @Summary Получение урока
// @Description Возвращает детальную информацию об уроке
// @Tags Lessons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param lesson_id query string true "ID урока" example("94f9a22f-3a83-4591-a988-7aa3f0ec6eb0")
// @Success 200 {object} lessons.GetLessonResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 401 {object} ErrorResponse "Требуется авторизация"
// @Failure 403 {object} ErrorResponse "Доступ запрещен"
// @Failure 404 {object} ErrorResponse "Урок не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Failure 503 {object} ErrorResponse "Сервис временно недоступен"
// @Router /lessons/lesson [get]
func (s *Server) GetLessonHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[lessons.GetLessonRequest](r.Context())

	resp, err := s.Lessons.GetLesson(r.Context(), body)

	isMember, err := s.IsTeacher(r.Context(), resp.Lesson.CourseID)
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

// GetLessonsHandler возвращает список уроков
// @Summary Получение списка уроков
// @Description Возвращает список уроков с возможностью фильтрации по курсу
// @Tags Lessons
// @Produce json
// @Security BearerAuth
// @Param course_id query string true "ID курса" example("6994aefe-6815-476b-bdc0-2ae5c4d0c18e")
// @Success 200 {object} lessons.GetLessonsResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 401 {object} ErrorResponse "Требуется авторизация"
// @Failure 404 {object} ErrorResponse "Курс не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /lessons/lessons [get]
func (s *Server) GetLessonsHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[lessons.GetLessonsRequest](r.Context())

	isMember, err := s.IsMember(r.Context(), body.CourseID)
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

	resp, err := s.Lessons.GetLessons(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler lessons.GetLessons error", slog.Any("error", err))

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

	WriteJSON(w, resp, http.StatusOK)
}

// UpdateLessonHandler обновляет информацию об уроке
// @Summary Обновление урока
// @Description Обновляет информацию об уроке
// @Tags Lessons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body lessons.UpdateLessonRequest true "Данные для обновления"
// @Success 200 {object} lessons.UpdateLessonResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 401 {object} ErrorResponse "Требуется авторизация"
// @Failure 404 {object} ErrorResponse "Урок не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /lessons/lesson/update [put]
func (s *Server) UpdateLessonHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[lessons.UpdateLessonRequest](r.Context())

	body1 := lessons.GetLessonRequest{
		LessonID: body.LessonID,
	}
	resp1, err := s.Lessons.GetLesson(r.Context(), body1)
	if err != nil {
		logger.Error(r.Context(), "Handler lessons.GetLesson error", slog.Any("error", err))

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

	isTeacher, err := s.IsTeacher(r.Context(), resp1.Lesson.CourseID)
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

	resp, err := s.Lessons.UpdateLesson(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler lessons.UpdateLesson error", slog.Any("error", err))

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

	err = redis.Delete(s.Redis, r.Context(), "Lessons.GetLesson", body.LessonID)
	logger.Debug(r.Context(), "Lessons.GetLesson uncached", slog.Any("error", err))

	WriteJSON(w, resp, http.StatusOK)
}

// DeleteLessonHandler удаляет урок
// @Summary Удаление урока
// @Description Удаляет урок по идентификатору
// @Tags Lessons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body lessons.DeleteLessonRequest true "Идентификатор урока"
// @Success 200 {object} lessons.DeleteLessonResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 401 {object} ErrorResponse "Требуется авторизация"
// @Failure 404 {object} ErrorResponse "Урок не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /lessons/lesson/delete [delete]
func (s *Server) DeleteLessonHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[lessons.DeleteLessonRequest](r.Context())

	body1 := lessons.GetLessonRequest{
		LessonID: body.LessonID,
	}
	resp1, err := s.Lessons.GetLesson(r.Context(), body1)
	if err != nil {
		logger.Error(r.Context(), "Handler lessons.GetLesson error", slog.Any("error", err))

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

	isTeacher, err := s.IsTeacher(r.Context(), resp1.Lesson.CourseID)
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

	resp, err := s.Lessons.DeleteLesson(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler lessons.DeleteLesson error", slog.Any("error", err))

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

	err = redis.Delete(s.Redis, r.Context(), "Lessons.GetLesson", body.LessonID)
	logger.Debug(r.Context(), "Lessons.GetLesson uncached", slog.Any("error", err))

	WriteJSON(w, resp, http.StatusOK)
}
