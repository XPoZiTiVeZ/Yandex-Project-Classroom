package server

import (
	"Classroom/Gateway/internal/courses"
	"Classroom/Gateway/pkg/logger"
	"log/slog"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
