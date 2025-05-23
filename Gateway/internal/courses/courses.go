package courses

import (
	rds "Classroom/Gateway/internal/redis"
	pb "Classroom/Gateway/pkg/api/courses"
	"Classroom/Gateway/pkg/config"
	"Classroom/Gateway/pkg/logger"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CoursesServiceClient struct {
	Conn           *grpc.ClientConn
	Client         pb.CoursesServiceClient
	DefaultTimeout time.Duration
}

func NewCoursesServiceClient(ctx context.Context, config *config.Config) (*CoursesServiceClient, error) {
	address, port := config.Courses.Address, config.Courses.Port
	timeout := config.Common.Timeout

	var opts []grpc.DialOption
	opts = append(
		opts, grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", address, port), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	state := conn.GetState()
	// if !conn.WaitForStateChange(ctx, state) {
	// 	return nil, fmt.Errorf("failed to wait for state change")
	// }
	// state = conn.GetState()

	logger.Info(ctx, "Connected to gRPC Courses", slog.String("address", address), slog.Int("port", port), slog.String("state", state.String()))

	client := pb.NewCoursesServiceClient(conn)

	return &CoursesServiceClient{
		Conn:           conn,
		Client:         client,
		DefaultTimeout: timeout,
	}, nil
}

func (s *CoursesServiceClient) CreateCourse(ctx context.Context, req CreateCourseRequest) (CreateCourseResponse, error) {
	logger.Debug(ctx, "Creating course", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := s.Client.CreateCourse(ctx, NewCreateCourseRequest(req))
	if err != nil {
		return CreateCourseResponse{}, err
	}

	logger.Debug(ctx, "Courses.CreateCourse succeed")
	return NewCreateCourseResponse(resp), nil
}

func (s *CoursesServiceClient) GetCourse(ctx context.Context, req GetCourseRequest) (GetCourseResponse, error) {
	logger.Debug(ctx, "Getting course", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	pbresp, err := s.Client.GetCourse(ctx, NewGetCourseRequest(req))
	if err != nil {
		return GetCourseResponse{}, err
	}

	logger.Debug(ctx, "Courses.GetCourse succeed")
	return NewGetCourseResponse(pbresp), nil
}

func (s *CoursesServiceClient) GetCourses(ctx context.Context, req GetCoursesRequest) (GetCoursesResponse, error) {
	logger.Debug(ctx, "Getting courses", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := s.Client.GetCourses(ctx, NewGetCoursesRequest(req))
	if err != nil {
		return GetCoursesResponse{}, err
	}

	logger.Debug(ctx, "Courses.GetCourses succeed")
	return NewGetCoursesResponse(resp), nil
}

func (s *CoursesServiceClient) GetCoursesByStudent(ctx context.Context, req GetCoursesByStudentRequest) (GetCoursesResponse, error) {
	logger.Debug(ctx, "Getting courses by student", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := s.Client.GetCoursesByStudent(ctx, NewGetCoursesByStudentRequest(req))
	if err != nil {
		return GetCoursesResponse{}, err
	}

	logger.Debug(ctx, "Courses.GetCoursesByStudent succeed")
	return NewGetCoursesResponse(resp), nil
}

func (s *CoursesServiceClient) GetCoursesByTeacher(ctx context.Context, req GetCoursesByTeacherRequest) (GetCoursesResponse, error) {
	logger.Debug(ctx, "Getting courses by teacher", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := s.Client.GetCoursesByTeacher(ctx, NewGetCoursesByTeacherRequest(req))
	if err != nil {
		return GetCoursesResponse{}, err
	}

	logger.Debug(ctx, "Courses.GetCoursesByTeacher succeed")
	return NewGetCoursesResponse(resp), nil
}

func (s *CoursesServiceClient) UpdateCourse(ctx context.Context, req UpdateCourseRequest) (UpdateCourseResponse, error) {
	logger.Debug(ctx, "Updating course", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := s.Client.UpdateCourse(ctx, NewUpdateCourseRequest(req))
	if err != nil {
		return UpdateCourseResponse{}, err
	}

	logger.Debug(ctx, "Courses.UpdateCourse succeed")
	return NewUpdateCourseResponse(resp), nil
}

func (s *CoursesServiceClient) DeleteCourse(ctx context.Context, req DeleteCourseRequest) (DeleteCourseResponse, error) {
	logger.Debug(ctx, "Deleting course", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := s.Client.DeleteCourse(ctx, NewDeleteCourseRequest(req))
	if err != nil {
		return DeleteCourseResponse{}, err
	}

	logger.Debug(ctx, "Courses.DeleteCourse succeed")
	return NewDeleteCourseResponse(resp), nil
}

func (s *CoursesServiceClient) EnrollUser(ctx context.Context, req EnrollUserRequest) (EnrollUserResponse, error) {
	logger.Debug(ctx, "Enrolling user", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := s.Client.EnrollUser(ctx, NewEnrollUserRequest(req))
	if err != nil {
		return EnrollUserResponse{}, err
	}

	logger.Debug(ctx, "Courses.EnrollUser succeed")
	return NewEnrollUserResponse(resp), nil
}

func (s *CoursesServiceClient) ExpelUser(ctx context.Context, req ExpelUserRequest) (ExpelUserResponse, error) {
	logger.Debug(ctx, "Expelling user", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := s.Client.ExpelUser(ctx, NewExpelUserRequest(req))
	if err != nil {
		return ExpelUserResponse{}, err
	}

	logger.Debug(ctx, "Courses.ExpelUser succeed")
	return NewExpelUserResponse(resp), nil
}

func (s *CoursesServiceClient) IsTeacher(ctx context.Context, rc *redis.Client, req IsTeacherRequest) (IsTeacherResponse, error) {
	logger.Debug(ctx, "Checking if user is teacher", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := rds.Get[IsTeacherResponse](rc, ctx, "Courses.IsTeacher", fmt.Sprintf("%s:%s", req.UserID, req.CourseID))
	if err == nil {
		logger.Debug(ctx, "Response was cached")
		return resp, nil
	}
	logger.Debug(ctx, "Response was not cached", slog.Any("error", err))

	pbresp, err := s.Client.IsTeacher(ctx, NewIsTeacherRequest(req))
	if err != nil {
		return IsTeacherResponse{}, err
	}

	resp = NewIsTeacherResponse(pbresp)
	rds.Put(rc, ctx, "Courses.IsTeacher", fmt.Sprintf("%s:%s", req.UserID, req.CourseID), resp, 24 * time.Hour)

	logger.Debug(ctx, "Courses.IsTeacher succeed")
	return resp, nil
}

func (s *CoursesServiceClient) IsMember(ctx context.Context, rc *redis.Client, req IsMemberRequest) (IsMemberResponse, error) {
	logger.Debug(ctx, "Checking if user is member", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := rds.Get[IsMemberResponse](rc, ctx, "Courses.IsMember", fmt.Sprintf("%s:%s", req.UserID, req.CourseID))
	if err == nil {
		logger.Debug(ctx, "Response was cached")
		return resp, nil
	}
	logger.Debug(ctx, "Response was not cached", slog.Any("error", err))

	pbresp, err := s.Client.IsMember(ctx, NewIsMemberRequest(req))
	if err != nil {
		return IsMemberResponse{}, err
	}

	resp = NewIsMemberResponse(pbresp)
	rds.Put(rc, ctx, "Courses.IsMember", fmt.Sprintf("%s:%s", req.UserID, req.CourseID), resp, 24 * time.Hour)

	logger.Debug(ctx, "Courses.IsMember succeed")
	return resp, nil
}

func (s *CoursesServiceClient) GetCourseStudents(ctx context.Context, req GetCourseStudentsRequest) (GetCourseStudentsResponse, error) {
	logger.Debug(ctx, "Getting course students", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := s.Client.GetCourseStudents(ctx, NewGetCourseStudentsRequest(req))
	if err != nil {
		return GetCourseStudentsResponse{}, err
	}

	logger.Debug(ctx, "Courses.GetCourseStudents succeed")
	return NewGetCourseStudentsResponse(resp), nil
}
