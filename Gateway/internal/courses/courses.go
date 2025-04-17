package courses

import (
	pb "Classroom/Gateway/pkg/api/courses"
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CoursesServiceClient struct {
	Conn   *grpc.ClientConn
	Client *pb.CoursesServiceClient
	DefaultTimeout time.Duration
}

func NewCoursesServiceClient(address string, port int, DefaultTimeout *time.Duration) (*CoursesServiceClient, error) {
	var opts []grpc.DialOption
	opts = append(
		opts, grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", address, port), opts...)
	if err != nil {
		slog.Error("fail to dial: %v", slog.Any("error", err))
		return nil, err
	}

	state := conn.GetState()
	slog.Info("Connected to grpc Courses", slog.String("address", address), slog.Int("port", port), slog.String("state", state.String()))

	client := pb.NewCoursesServiceClient(conn)

	timeout := 10 * time.Second
	if DefaultTimeout != nil {
		timeout = *DefaultTimeout
	}

	return &CoursesServiceClient{
		Conn: conn,
		Client: &client,
		DefaultTimeout: timeout,
	}, nil
}

func (s *CoursesServiceClient) CreateCourse(ctx context.Context, req CreateCourseRequest) (CreateCourseResponse, error) {
	slog.Debug("creating course", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := (*s.Client).CreateCourse(ctx, NewCreateCourseRequest(req))
	if err != nil {
		return CreateCourseResponse{}, err
	}

	slog.Debug("courses.CreateCourse succeed")
	return NewCreateCourseResponse(resp), nil
}

func (s *CoursesServiceClient) GetCourse(ctx context.Context, req GetCourseRequest) (GetCourseResponse, error) {
	slog.Debug("getting course", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := (*s.Client).GetCourse(ctx, NewGetCourseRequest(req))
	if err != nil {
		return GetCourseResponse{}, err
	}

	slog.Debug("courses.GetCourse succeed")
	return NewGetCourseResponse(resp), nil
}

func (s *CoursesServiceClient) GetCourses(ctx context.Context, req GetCoursesRequest) (GetCoursesResponse, error) {
	slog.Debug("getting courses", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := (*s.Client).GetCourses(ctx, NewGetCoursesRequest(req))
	if err != nil {
		return GetCoursesResponse{}, err
	}

	slog.Debug("courses.GetCourses succeed")
	return NewGetCoursesResponse(resp), nil
}

func (s *CoursesServiceClient) GetCoursesByStudent(ctx context.Context, req GetCoursesByStudentRequest) (GetCoursesResponse, error) {
	slog.Debug("getting courses by student", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := (*s.Client).GetCoursesByStudent(ctx, NewGetCoursesByStudentRequest(req))
	if err != nil {
		return GetCoursesResponse{}, err
	}

	slog.Debug("courses.GetCoursesByStudent succeed")
	return NewGetCoursesResponse(resp), nil
}

func (s *CoursesServiceClient) GetCoursesByTeacher(ctx context.Context, req GetCoursesByTeacherRequest) (GetCoursesResponse, error) {
	slog.Debug("getting courses by teacher", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := (*s.Client).GetCoursesByTeacher(ctx, NewGetCoursesByTeacherRequest(req))
	if err != nil {
		return GetCoursesResponse{}, err
	}

	slog.Debug("courses.GetCoursesByTeacher succeed")
	return NewGetCoursesResponse(resp), nil
}

func (s *CoursesServiceClient) UpdateCourse(ctx context.Context, req UpdateCourseRequest) (UpdateCourseResponse, error) {
	slog.Debug("updating course", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := (*s.Client).UpdateCourse(ctx, NewUpdateCourseRequest(req))
	if err != nil {
		return UpdateCourseResponse{}, err
	}

	slog.Debug("courses.UpdateCourse succeed")
	return NewUpdateCourseResponse(resp), nil
}

func (s *CoursesServiceClient) DeleteCourse(ctx context.Context, req DeleteCourseRequest) (DeleteCourseResponse, error) {
	slog.Debug("deleting course", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := (*s.Client).DeleteCourse(ctx, NewDeleteCourseRequest(req))
	if err != nil {
		return DeleteCourseResponse{}, err
	}

	slog.Debug("courses.DeleteCourse succeed")
	return NewDeleteCourseResponse(resp), nil
}

func (s *CoursesServiceClient) EnrollUser(ctx context.Context, req EnrollUserRequest) (EnrollUserResponse, error) {
	slog.Debug("enrolling user", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := (*s.Client).EnrollUser(ctx, NewEnrollUserRequest(req))
	if err != nil {
		return EnrollUserResponse{}, err
	}

	slog.Debug("courses.EnrollUser succeed")
	return NewEnrollUserResponse(resp), nil
}

func (s *CoursesServiceClient) ExpelUser(ctx context.Context, req ExpelUserRequest) (ExpelUserResponse, error) {
	slog.Debug("expelling user", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := (*s.Client).ExpelUser(ctx, NewExpelUserRequest(req))
	if err != nil {
		return ExpelUserResponse{}, err
	}

	slog.Debug("courses.ExpelUser succeed")
	return NewExpelUserResponse(resp), nil
}

func (s *CoursesServiceClient) IsTeacher(ctx context.Context, req IsTeacherRequest) (IsTeacherResponse, error) {
	slog.Debug("checking if user is teacher", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := (*s.Client).IsTeacher(ctx, NewIsTeacherRequest(req))
	if err != nil {
		return IsTeacherResponse{}, err
	}

	slog.Debug("courses.IsTeacher succeed")
	return NewIsTeacherResponse(resp), nil
}

func (s *CoursesServiceClient) IsMember(ctx context.Context, req IsMemberRequest) (IsMemberResponse, error) {
	slog.Debug("checking if user is member", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()


	resp, err := (*s.Client).IsMember(ctx, NewIsMemberRequest(req))
	if err != nil {
		return IsMemberResponse{}, err
	}

	slog.Debug("courses.IsMember succeed")
	return NewIsMemberResponse(resp), nil
}

func (s *CoursesServiceClient) GetCourseStudents(ctx context.Context, req GetCourseStudentsRequest) (GetCourseStudentsResponse, error) {
	slog.Debug("getting course students", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := (*s.Client).GetCourseStudents(ctx, NewGetCourseStudentsRequest(req))
	if err != nil {
		return GetCourseStudentsResponse{}, err
	}

	slog.Debug("courses.GetCourseStudents succeed")
	return NewGetCourseStudentsResponse(resp), nil
}