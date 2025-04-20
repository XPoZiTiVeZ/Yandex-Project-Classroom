package server

import (
	"Classroom/Gateway/internal/courses"
	he "Classroom/Gateway/internal/errors"
	"log/slog"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateCourseHandler(w http.ResponseWriter, r *http.Request) (any, error) {
	var body courses.CreateCourseRequest = r.Context().Value("body").(courses.CreateCourseRequest)

	resp, err := s.Courses.CreateCourse(r.Context(), body)
	if err != nil {
		slog.Error("handler courses.CreateCourse error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.AlreadyExists:
				he.AlreadyExists(w)
			}

		} else {
			he.ServiceUnavailable(w)
		}
	}

	return resp, err
}

func (s *Server) GetCourseHandler(w http.ResponseWriter, r *http.Request) (any, error) {
	var body courses.GetCourseRequest = r.Context().Value("body").(courses.GetCourseRequest)

	resp, err := s.Courses.GetCourse(r.Context(), body)
	if err != nil {
		slog.Error("handler courses.GetCourse error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				he.NotFound(w)
			}

		} else {
			he.ServiceUnavailable(w)
		}
	}

	return resp, err
}

func (s *Server) GetCoursesHandler(w http.ResponseWriter, r *http.Request) (any, error) {
	var body courses.GetCoursesRequest = r.Context().Value("body").(courses.GetCoursesRequest)

	resp, err := s.Courses.GetCourses(r.Context(), body)
	if err != nil {
		slog.Error("handler courses.GetCourses error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				he.NotFound(w)
			}

		} else {
			he.ServiceUnavailable(w)
		}
	}

	return resp, err
}

func (s *Server) GetCoursesByStudentHandler(w http.ResponseWriter, r *http.Request) (any, error) {
	var body courses.GetCoursesByStudentRequest = r.Context().Value("body").(courses.GetCoursesByStudentRequest)

	resp, err := s.Courses.GetCoursesByStudent(r.Context(), body)
	if err != nil {
		slog.Error("handler courses.GetCoursesByStudent error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				he.NotFound(w)
			}

		} else {
			he.ServiceUnavailable(w)
		}
	}

	return resp, err
}

func (s *Server) GetCoursesByTeacherHandler(w http.ResponseWriter, r *http.Request) (any, error) {
	var body courses.GetCoursesByTeacherRequest = r.Context().Value("body").(courses.GetCoursesByTeacherRequest)

	resp, err := s.Courses.GetCoursesByTeacher(r.Context(), body)
	if err != nil {
		slog.Error("handler courses.GetCoursesByTeacher error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				he.NotFound(w)
			}

		} else {
			he.ServiceUnavailable(w)
		}
	}

	return resp, err
}

func (s *Server) UpdateCourseHandler(w http.ResponseWriter, r *http.Request) (any, error) {
	var body courses.UpdateCourseRequest = r.Context().Value("body").(courses.UpdateCourseRequest)

	resp, err := s.Courses.UpdateCourse(r.Context(), body)
	if err != nil {
		slog.Error("handler courses.UpdateCourse error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				he.NotFound(w)
			}

		} else {
			he.ServiceUnavailable(w)
		}
	}

	return resp, err
}

func (s *Server) DeleteCourseHandler(w http.ResponseWriter, r *http.Request) (any, error) {
	var body courses.DeleteCourseRequest = r.Context().Value("body").(courses.DeleteCourseRequest)

	resp, err := s.Courses.DeleteCourse(r.Context(), body)
	if err != nil {
		slog.Error("handler courses.DeleteCourse error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				he.NotFound(w)
			}

		} else {
			he.ServiceUnavailable(w)
		}
	}

	return resp, err
}

func (s *Server) EnrollUserHandler(w http.ResponseWriter, r *http.Request) (any, error) {
	var body courses.EnrollUserRequest = r.Context().Value("body").(courses.EnrollUserRequest)

	resp, err := s.Courses.EnrollUser(r.Context(), body)
	if err != nil {
		slog.Error("handler courses.EnrollUser error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				he.NotFound(w)
			}

		} else {
			he.ServiceUnavailable(w)
		}
	}

	return resp, err
}

func (s *Server) ExpelUserHandler(w http.ResponseWriter, r *http.Request) (any, error) {
	var body courses.ExpelUserRequest = r.Context().Value("body").(courses.ExpelUserRequest)

	resp, err := s.Courses.ExpelUser(r.Context(), body)
	if err != nil {
		slog.Error("handler courses.EnrollUser error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				he.NotFound(w)
			}

		} else {
			he.ServiceUnavailable(w)
		}
	}

	return resp, err
}

func (s *Server) GetCourseStudentsHandler(w http.ResponseWriter, r *http.Request) (any, error) {
	var body courses.GetCourseStudentsRequest = r.Context().Value("body").(courses.GetCourseStudentsRequest)

	resp, err := s.Courses.GetCourseStudents(r.Context(), body)
	if err != nil {
		slog.Error("handler courses.GetCourseStudents error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				he.NotFound(w)
			}

		} else {
			he.ServiceUnavailable(w)
		}
	}

	return resp, err
}
