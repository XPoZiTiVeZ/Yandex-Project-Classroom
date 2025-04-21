package server

import (
	"Classroom/Gateway/internal/courses"
	"Classroom/Gateway/pkg/logger"
	"log/slog"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateCourseHandler создает новый курс
// @Summary Создание курса
// @Description Создает новый курс
// @Tags Courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body courses.CreateCourseRequest true "Данные для создания курса"
// @Success 201 {object} courses.CreateCourseResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /courses/create [post]
func (s *Server) CreateCourseHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[courses.CreateCourseRequest](r.Context())

	resp, err := s.Courses.CreateCourse(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler courses.CreateCourse error", slog.Any("error", err))

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

// GetCourseHandler возвращает информацию о курсе
// @Summary Получение курса
// @Description Возвращает детальную информацию о курсе
// @Tags Courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body courses.GetCourseRequest true "Идентификатор курса"
// @Success 200 {object} courses.GetCourseResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 404 {object} ErrorResponse "Курс не найден"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /courses/course [post]
func (s *Server) GetCourseHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[courses.GetCourseRequest](r.Context())

	resp, err := s.Courses.GetCourse(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler courses.GetCourse error", slog.Any("error", err))

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
	}

	WriteJSON(w, resp, http.StatusOK)
}

// GetCoursesHandler возвращает список курсов
// @Summary Получение списка курсов
// @Description Возвращает список курсов с возможностью фильтрации
// @Tags Courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body courses.GetCoursesRequest true "Параметры фильтрации"
// @Success 200 {object} courses.GetCoursesResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /courses/courses [post]
func (s *Server) GetCoursesHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[courses.GetCoursesRequest](r.Context())

	resp, err := s.Courses.GetCourses(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler courses.GetCourses error", slog.Any("error", err))

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

	WriteJSON(w, resp, http.StatusOK)
}

// GetCoursesByStudentHandler возвращает курсы студента
// @Summary Получение курсов студента
// @Description Возвращает список курсов, на которые записан студент
// @Tags Courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body courses.GetCoursesByStudentRequest true "Идентификатор студента"
// @Success 200 {object} courses.GetCoursesResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /courses/student-courses [post]
func (s *Server) GetCoursesByStudentHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[courses.GetCoursesByStudentRequest](r.Context())

	resp, err := s.Courses.GetCoursesByStudent(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "handler courses.GetCoursesByStudent error", slog.Any("error", err))

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

	WriteJSON(w, resp, http.StatusOK)
}

// GetCoursesByTeacherHandler возвращает курсы преподавателя
// @Summary Получение курсов преподавателя
// @Description Возвращает список курсов, которые ведет преподаватель
// @Tags Courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body courses.GetCoursesByTeacherRequest true "Идентификатор преподавателя"
// @Success 200 {object} courses.GetCoursesResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /courses/teacher-courses [post]
func (s *Server) GetCoursesByTeacherHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[courses.GetCoursesByTeacherRequest](r.Context())

	resp, err := s.Courses.GetCoursesByTeacher(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler courses.GetCoursesByTeacher error", slog.Any("error", err))

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

	WriteJSON(w, resp, http.StatusOK)
}

// UpdateCourseHandler обновляет информацию о курсе
// @Summary Обновление курса
// @Description Обновляет информацию о курсе
// @Tags Courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body courses.UpdateCourseRequest true "Данные для обновления"
// @Success 200 {object} courses.UpdateCourseResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 404 {object} ErrorResponse "Курс не найден"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /courses/course/update [patch]
func (s *Server) UpdateCourseHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[courses.UpdateCourseRequest](r.Context())

	resp, err := s.Courses.UpdateCourse(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler courses.UpdateCourse error", slog.Any("error", err))

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

// DeleteCourseHandler удаляет курс
// @Summary Удаление курса
// @Description Удаляет курс по идентификатору
// @Tags Courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body courses.DeleteCourseRequest true "Идентификатор курса"
// @Success 200 {object} courses.DeleteCourseResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 404 {object} ErrorResponse "Курс не найден"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /courses/course/delete [delete]
func (s *Server) DeleteCourseHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[courses.DeleteCourseRequest](r.Context())

	resp, err := s.Courses.DeleteCourse(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler courses.DeleteCourse error", slog.Any("error", err))

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

// EnrollUserHandler записывает пользователя на курс
// @Summary Запись на курс
// @Description Записывает пользователя на указанный курс
// @Tags Courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body courses.EnrollUserRequest true "Данные для записи"
// @Success 200 {object} courses.EnrollUserResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 404 {object} ErrorResponse "Курс или пользователь не найден"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /courses/course/enroll [post]
func (s *Server) EnrollUserHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[courses.EnrollUserRequest](r.Context())

	resp, err := s.Courses.EnrollUser(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler courses.EnrollUser error", slog.Any("error", err))

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

// ExpelUserHandler отчисляет пользователя с курса
// @Summary Отчисление с курса
// @Description Отчисляет пользователя с указанного курса
// @Tags Courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body courses.ExpelUserRequest true "Данные для отчисления"
// @Success 200 {object} courses.ExpelUserResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 404 {object} ErrorResponse "Курс или пользователь не найден"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /courses/course/expel [post]
func (s *Server) ExpelUserHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[courses.ExpelUserRequest](r.Context())

	resp, err := s.Courses.ExpelUser(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler courses.EnrollUser error", slog.Any("error", err))

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

// GetCourseStudentsHandler возвращает студентов курса
// @Summary Получение студентов курса
// @Description Возвращает список студентов, записанных на курс
// @Tags Courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body courses.GetCourseStudentsRequest true "Идентификатор курса"
// @Success 200 {object} courses.GetCourseStudentsResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 404 {object} ErrorResponse "Курс не найден"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /courses/course/students [post]
func (s *Server) GetCourseStudentsHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[courses.GetCourseStudentsRequest](r.Context())

	resp, err := s.Courses.GetCourseStudents(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler courses.GetCourseStudents error", slog.Any("error", err))

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
