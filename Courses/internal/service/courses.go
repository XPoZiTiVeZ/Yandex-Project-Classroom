package service

import (
	"Classroom/Courses/internal/domain"
	"Classroom/Courses/internal/dto"
	pb "Classroom/Courses/pkg/api/courses"
	"Classroom/Courses/pkg/events"
	"context"
	"errors"
	"log/slog"
	"slices"
	"time"

	"github.com/go-playground/validator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type CourseRepo interface {
	Create(ctx context.Context, dto dto.CreateCourseDTO) (domain.Course, error)
	GetByID(ctx context.Context, courseID string) (domain.Course, error)
	Update(ctx context.Context, dto dto.UpdateCourseDTO) (domain.Course, error)
	Delete(ctx context.Context, courseID string) (domain.Course, error)

	// Делает фильтрацию по visibility, start_time, end_time и teacher_id
	ListByStudentID(ctx context.Context, teacherID string) ([]domain.Course, error)
	ListByTeacherID(ctx context.Context, teacherID string) ([]domain.Course, error)

	ListCourseStudents(ctx context.Context, courseID string, index, limit int32) ([]domain.Student, int32, error)
	EnrollUser(ctx context.Context, courseID, studentID string) (domain.Enrollment, error)
	ExpelUser(ctx context.Context, courseID, studentID string) (domain.Enrollment, error)

	IsTeacher(ctx context.Context, courseID, teacherID string) (bool, error)
	IsMember(ctx context.Context, courseID, userID string) (bool, error)
}

type Producer interface {
	PublishUserEnrolled(event events.UserEnrolled) error
	PublishUserExpelled(event events.UserExpelled) error
}

type CoursesService struct {
	pb.UnimplementedCoursesServiceServer
	logger   *slog.Logger
	validate *validator.Validate
	repo     CourseRepo
	producer Producer
}

func NewCoursesService(logger *slog.Logger, repo CourseRepo, producer Producer) *CoursesService {
	validate := validator.New()
	return &CoursesService{repo: repo, logger: logger, validate: validate, producer: producer}
}

func (s *CoursesService) CreateCourse(ctx context.Context, req *pb.CreateCourseRequest) (*pb.CreateCourseResponse, error) {
	dto := dto.CreateCourseDTO{
		TeacherID:   req.UserId,
		Title:       req.Title,
		Description: req.Description,
		Visibility:  req.Visibility,
		StartTime:   timestampToTime(req.StartTime),
		EndTime:     timestampToTime(req.EndTime),
	}

	if err := s.validate.Struct(dto); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err.Error())
	}

	course, err := s.repo.Create(ctx, dto)
	if err != nil {
		s.logger.Error("failed to create course", "error", err)
		return nil, status.Error(codes.Internal, "failed to create course")
	}

	s.logger.Info("course created", "id", course.ID, "title", course.Title)
	return &pb.CreateCourseResponse{Course: courseToPb(course)}, nil
}

func (s *CoursesService) GetCourse(ctx context.Context, req *pb.GetCourseRequest) (*pb.GetCourseResponse, error) {
	if err := s.validate.Var(req.CourseId, "required,uuid"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid course id")
	}
	if err := s.validate.Var(req.UserId, "required,uuid"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	course, err := s.repo.GetByID(ctx, req.CourseId)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "course not found")
	}
	if err != nil {
		s.logger.Error("failed to get course", "error", err)
		return nil, status.Error(codes.Internal, "failed to get course")
	}

	if course.TeacherID == req.UserId {
		return &pb.GetCourseResponse{Course: courseToPb(course)}, nil
	}
	if !course.Visibility {
		return nil, status.Error(codes.NotFound, "course is hidden")
	}
	if course.StartTime != nil && course.StartTime.Before(time.Now()) {
		return nil, status.Error(codes.NotFound, "course is not started yet")
	}
	if course.EndTime != nil && course.EndTime.After(time.Now()) {
		return nil, status.Error(codes.NotFound, "course is over")
	}

	return &pb.GetCourseResponse{Course: courseToPb(course)}, nil
}

func (s *CoursesService) UpdateCourse(ctx context.Context, req *pb.UpdateCourseRequest) (*pb.UpdateCourseResponse, error) {
	dto := dto.UpdateCourseDTO{
		ID:          req.CourseId,
		Title:       req.Title,
		Description: req.Description,
		Visibility:  req.Visibility,
		StartTime:   timestampToTime(req.StartTime),
		EndTime:     timestampToTime(req.EndTime),
	}

	if err := s.validate.Struct(dto); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err.Error())
	}

	course, err := s.repo.Update(ctx, dto)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "course not found")
	}
	if err != nil {
		s.logger.Error("failed to update course", "error", err)
		return nil, status.Error(codes.Internal, "failed to update course")
	}

	s.logger.Info("course updated", "id", course.ID, "title", course.Title)
	return &pb.UpdateCourseResponse{Course: courseToPb(course)}, nil
}

func (s *CoursesService) DeleteCourse(ctx context.Context, req *pb.DeleteCourseRequest) (*pb.DeleteCourseResponse, error) {
	if err := s.validate.Var(req.CourseId, "required,uuid"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid course id")
	}

	course, err := s.repo.Delete(ctx, req.CourseId)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "course not found")
	}
	if err != nil {
		s.logger.Error("failed to delete course", "error", err)
		return nil, status.Error(codes.Internal, "failed to delete course")
	}

	s.logger.Info("course deleted", "id", course.ID, "title", course.Title)
	return &pb.DeleteCourseResponse{Course: courseToPb(course)}, nil
}

func (s *CoursesService) EnrollUser(ctx context.Context, req *pb.EnrollUserRequest) (*pb.EnrollUserResponse, error) {
	if err := s.validate.Var(req.CourseId, "required,uuid"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid course id")
	}
	if err := s.validate.Var(req.UserId, "required,uuid"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}
	enrollment, err := s.repo.EnrollUser(ctx, req.CourseId, req.UserId)
	if err != nil {
		s.logger.Error("failed to enroll user", "error", err)
		return nil, status.Error(codes.Internal, "failed to enroll user")
	}

	err = s.producer.PublishUserEnrolled(events.UserEnrolled{
		CourseID: enrollment.CourseID,
		UserID:   enrollment.StudentID,
	})

	if err != nil {
		s.logger.Error("failed to publish user enrolled event", "error", err)
		// не возвращаем ошибку потому что действие и так было выполнено в бд, фикс будет через логи
	}

	s.logger.Info("user enrolled", "course_id", enrollment.CourseID, "student_id", enrollment.StudentID)
	return &pb.EnrollUserResponse{
		Enrollment: &pb.Enrollment{
			CourseId:   enrollment.CourseID,
			StudentId:  enrollment.StudentID,
			EnrolledAt: timestamppb.New(enrollment.EnrolledAt),
		},
	}, nil
}

func (s *CoursesService) ExpelUser(ctx context.Context, req *pb.ExpelUserRequest) (*pb.ExpelUserResponse, error) {
	if err := s.validate.Var(req.CourseId, "required,uuid"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid course id")
	}
	if err := s.validate.Var(req.UserId, "required,uuid"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	enrollment, err := s.repo.ExpelUser(ctx, req.CourseId, req.UserId)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "enrollment not found")
	}
	if err != nil {
		s.logger.Error("failed to expel user", "error", err)
		return nil, status.Error(codes.Internal, "failed to expel user")
	}

	err = s.producer.PublishUserExpelled(events.UserExpelled{
		CourseID: enrollment.CourseID,
		UserID:   enrollment.StudentID,
	})

	if err != nil {
		s.logger.Error("failed to publish user expelled event", "error", err)
		// не возвращаем ошибку потому что действие и так было выполнено в бд, фикс будет через логи
	}

	s.logger.Info("user expelled", "course_id", enrollment.CourseID, "student_id", enrollment.StudentID)
	return &pb.ExpelUserResponse{
		Enrollment: &pb.Enrollment{
			CourseId:   enrollment.CourseID,
			StudentId:  enrollment.StudentID,
			EnrolledAt: timestamppb.New(enrollment.EnrolledAt),
		},
	}, nil
}

func (s *CoursesService) IsTeacher(ctx context.Context, req *pb.IsTeacherRequest) (*pb.IsTeacherResponse, error) {
	if err := s.validate.Var(req.CourseId, "required,uuid"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid course id")
	}
	if err := s.validate.Var(req.UserId, "required,uuid"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	isTeacher, err := s.repo.IsTeacher(ctx, req.CourseId, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to check teacher")
	}

	return &pb.IsTeacherResponse{IsTeacher: isTeacher}, nil
}

func (s *CoursesService) IsMember(ctx context.Context, req *pb.IsMemberRequest) (*pb.IsMemberResponse, error) {
	if err := s.validate.Var(req.CourseId, "required,uuid"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid course id")
	}
	if err := s.validate.Var(req.UserId, "required,uuid"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	isMember, err := s.repo.IsMember(ctx, req.CourseId, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to check member")
	}

	if !isMember {
		isTeacher, err := s.repo.IsTeacher(ctx, req.CourseId, req.UserId)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to check teacher")
		}

		if isTeacher {
			return &pb.IsMemberResponse{IsMember: true}, nil
		}
	}

	return &pb.IsMemberResponse{IsMember: isMember}, nil
}

func (s *CoursesService) GetCourses(ctx context.Context, req *pb.GetCoursesRequest) (*pb.GetCoursesResponse, error) {
	if err := s.validate.Var(req.UserId, "required,uuid"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}
	studentCourses, err := s.repo.ListByStudentID(ctx, req.UserId)
	if err != nil {
		s.logger.Error("failed to get courses by student", "error", err)
		return nil, status.Error(codes.Internal, "failed to get courses by student")
	}

	teacherCourses, err := s.repo.ListByTeacherID(ctx, req.UserId)
	if err != nil {
		s.logger.Error("failed to get courses by teacher", "error", err)
		return nil, status.Error(codes.Internal, "failed to get courses by teacher")
	}

	pbCourses := make([]*pb.Course, 0, len(teacherCourses)+len(studentCourses))
	for _, c := range teacherCourses {
		pbCourses = append(pbCourses, courseToPb(c))
	}
	for _, c := range studentCourses {
		pbCourses = append(pbCourses, courseToPb(c))
	}

	slices.SortFunc(pbCourses, func(a, b *pb.Course) int {
		if a.CreatedAt.AsTime().Before(b.CreatedAt.AsTime()) {
			return 1
		}
		return -1
	})

	return &pb.GetCoursesResponse{Courses: pbCourses}, nil
}

func (s *CoursesService) GetCoursesByStudent(ctx context.Context, req *pb.GetCoursesByStudentRequest) (*pb.GetCoursesResponse, error) {
	if err := s.validate.Var(req.StudentId, "required,uuid"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid student id")
	}
	courses, err := s.repo.ListByStudentID(ctx, req.StudentId)
	if err != nil {
		s.logger.Error("failed to get courses by student", "error", err)
		return nil, status.Error(codes.Internal, "failed to get courses by student")
	}

	pbCourses := make([]*pb.Course, len(courses))
	for i, c := range courses {
		pbCourses[i] = courseToPb(c)
	}

	return &pb.GetCoursesResponse{Courses: pbCourses}, nil
}

func (s *CoursesService) GetCoursesByTeacher(ctx context.Context, req *pb.GetCoursesByTeacherRequest) (*pb.GetCoursesResponse, error) {
	if err := s.validate.Var(req.TeacherId, "required,uuid"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid teacher id")
	}
	courses, err := s.repo.ListByTeacherID(ctx, req.TeacherId)
	if err != nil {
		s.logger.Error("failed to get courses by teacher", "error", err)
		return nil, status.Error(codes.Internal, "failed to get courses by teacher")
	}

	pbCourses := make([]*pb.Course, len(courses))
	for i, c := range courses {
		pbCourses[i] = courseToPb(c)
	}

	return &pb.GetCoursesResponse{Courses: pbCourses}, nil
}

func (s *CoursesService) GetCourseStudents(ctx context.Context, req *pb.GetCourseStudentsRequest) (*pb.GetCourseStudentsResponse, error) {
	if err := s.validate.Var(req.CourseId, "required,uuid"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid course id")
	}

	students, total, err := s.repo.ListCourseStudents(ctx, req.CourseId, req.Index, req.Limit)
	if err != nil {
		s.logger.Error("failed to get course students", "error", err)
		return nil, status.Error(codes.Internal, "failed to get course students")
	}

	pbStudents := make([]*pb.Student, len(students))
	for i, s := range students {
		pbStudents[i] = &pb.Student{
			UserId:    s.UserID,
			Email:     s.Email,
			FirstName: s.FirstName,
			LastName:  s.LastName,
		}
	}
	return &pb.GetCourseStudentsResponse{Students: pbStudents, Total: total, Index: req.Index + int32(len(pbStudents))}, nil
}

func courseToPb(c domain.Course) *pb.Course {
	return &pb.Course{
		CourseId:    c.ID,
		TeacherId:   c.TeacherID,
		Title:       c.Title,
		Description: c.Description,
		Visibility:  c.Visibility,
		StartTime:   timeToTimestamp(c.StartTime),
		EndTime:     timeToTimestamp(c.EndTime),
		CreatedAt:   timestamppb.New(c.CreatedAt),
	}
}

func timestampToTime(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}

func timeToTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}
