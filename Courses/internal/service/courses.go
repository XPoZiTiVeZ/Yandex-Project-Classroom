package service

import (
	repo "Classroom/Courses/internal/repo/postgres"
	pb "Classroom/Courses/pkg/api/courses"
	"context"
	"fmt"
	"time"

	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type CoursesService struct {
	pb.UnimplementedCoursesServiceServer
	repo *repo.CourseRepo
}

func NewCoursesService(repo *repo.CourseRepo) *CoursesService {
	return &CoursesService{repo: repo}
}

func (s *CoursesService) CreateCourse(ctx context.Context, req *pb.CreateCourseRequest) (*pb.CreateCourseResponse, error) {
	startTime := req.GetStartTime().AsTime()
	endTime := req.GetEndTime().AsTime()

	id, err := s.repo.CreateCourse(
		ctx,
		req.GetUserId(),
		req.GetTitle(),
		req.GetDescription(),
		req.GetVisibility(),
		&startTime,
		&endTime,
	)
	if err != nil {
		return nil, err
	}

	return &pb.CreateCourseResponse{CourseId: id}, nil
}

func (s *CoursesService) GetCourse(ctx context.Context, req *pb.GetCourseRequest) (*pb.GetCourseResponse, error) {
	c, err := s.repo.GetCourse(ctx, req.GetCourseId())
	if err != nil {
		return nil, err
	}

	return &pb.GetCourseResponse{
		Course: &pb.Course{
			CourseId:    c.CourseID,
			TeacherId:   c.TeacherID,
			Title:       c.Title,
			Description: c.Description,
			Visibility:  c.Visibility,
			StartTime:   timestamppb.New(c.StartTime),
			EndTime:     timestamppb.New(c.EndTime),
		},
	}, nil
}

func (s *CoursesService) UpdateCourse(ctx context.Context, req *pb.UpdateCourseRequest) (*pb.UpdateCourseResponse, error) {
	title := req.GetTitle()
	description := req.GetDescription()
	visibility := req.GetVisibility()
	pb_startTime := req.GetStartTime()
	pb_endTime := req.GetEndTime()

	var startTime *time.Time
	if pb_startTime != nil {
		time := pb_startTime.AsTime()
		startTime = &time
	}

	var endTime *time.Time
	if pb_endTime != nil {
		time := pb_endTime.AsTime()
		endTime = &time
	}

	err := s.repo.UpdateCourse(
		ctx,
		req.CourseId,
		&title,
		&description,
		&visibility,
		startTime,
		endTime,
	)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateCourseResponse{}, nil
}

func (s *CoursesService) DeleteCourse(ctx context.Context, req *pb.DeleteCourseRequest) (*pb.DeleteCourseResponse, error) {
	err := s.repo.DeleteCourse(ctx, req.CourseId)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteCourseResponse{}, nil
}

func (s *CoursesService) EnrollUser(ctx context.Context, req *pb.EnrollUserRequest) (*pb.EnrollUserResponse, error) {
	err := s.repo.EnrollUser(ctx, req.CourseId, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.EnrollUserResponse{}, nil
}

func (s *CoursesService) ExpelUser(ctx context.Context, req *pb.ExpelUserRequest) (*pb.ExpelUserResponse, error) {
	err := s.repo.EnrollUser(ctx, req.CourseId, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.ExpelUserResponse{}, nil
}

func (s *CoursesService) GetCourses(ctx context.Context, req *pb.GetCoursesRequest) (*pb.GetCoursesResponse, error) {
	courses, err := s.repo.GetUserCourses(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	var res []*pb.Course
	for _, c := range courses {
		course := pb.Course{
			CourseId:    c.CourseID,
			TeacherId:   c.TeacherID,
			Title:       c.Title,
			Description: c.Description,
			Visibility:  c.Visibility,
			StartTime:   timestamppb.New(c.StartTime),
			EndTime:     timestamppb.New(c.EndTime),
		}
		res = append(res, &course)
	}

	return &pb.GetCoursesResponse{Course: res}, nil
}

func (s *CoursesService) IsTeacher(ctx context.Context, req *pb.IsTeacherRequest) (*pb.IsTeacherResponse, error) {
	isTeacher, err := s.repo.IsTeacher(ctx, req.CourseId, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.IsTeacherResponse{IsTeacher: isTeacher}, nil
}

func (s *CoursesService) IsMember(ctx context.Context, req *pb.IsMemberRequest) (*pb.IsMemberResponse, error) {
	isMember, err := s.repo.IsMember(ctx, req.CourseId, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.IsMemberResponse{IsMember: isMember}, nil
}

func (s *CoursesService) GetCourseMembers(ctx context.Context, req *pb.GetCourseMembersRequest) (*pb.GetCourseMembersResponse, error) {
    if req.Limit <= 0 {
        req.Limit = 50
    }
    if req.Index < 0 {
        req.Index = 0
    }

    total, members, err := s.repo.GetCourseMembers(
        ctx,
        req.CourseId,
        req.Index,
        req.Limit,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to get course members: %w", err)
    }

    pbMembers := []*pb.Member{}
    for _, m := range members {
        pbMembers = append(pbMembers, &pb.Member{
            UserId:    m.UserID,
            Email:     m.Email,
            FirstName: m.FirstName,
            LastName:  m.LastName,
        })
    }

    return &pb.GetCourseMembersResponse{
        Total:   total,
        Index:   req.Index,
        Members: pbMembers,
    }, nil
}