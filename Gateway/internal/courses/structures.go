package courses

import (
	"time"

	pb "Classroom/Gateway/pkg/api/courses"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type CreateCourseRequest struct {
	UserID      string     `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Visibility  bool       `json:"visibility"`
	StartTime   *time.Time `json:"start_time,omitempty"`
	EndTime     *time.Time `json:"end_time,omitempty"`
}

func NewCreateCourseRequest(req CreateCourseRequest) *pb.CreateCourseRequest {
	pbReq := &pb.CreateCourseRequest{
		UserId:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		Visibility:  req.Visibility,
	}

	if req.StartTime != nil {
		pbReq.StartTime = timestamppb.New(*req.StartTime)
	}

	if req.EndTime != nil {
		pbReq.EndTime = timestamppb.New(*req.EndTime)
	}

	return pbReq
}

type CreateCourseResponse struct {
	CourseID string `json:"course_id"`
}

func NewCreateCourseResponse(resp *pb.CreateCourseResponse) CreateCourseResponse {
	return CreateCourseResponse{
		CourseID: resp.GetCourseId(),
	}
}

type GetCourseRequest struct {
	CourseID string `json:"course_id"`
}

func NewGetCourseRequest(req GetCourseRequest) *pb.GetCourseRequest {
	return &pb.GetCourseRequest{
		CourseId: req.CourseID,
	}
}

type Course struct {
	CourseID    string	   `json:"course_id"`
	TeacherID   string     `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Visibility  bool       `json:"visibility"`
	StartTime   *time.Time `json:"start_time,omitempty"`
	EndTime     *time.Time `json:"end_time,omitempty"`
}

type GetCourseResponse struct {
	Course *Course `json:"course,omitempty"`
}

func NewGetCourseResponse(resp *pb.GetCourseResponse) GetCourseResponse {
	pbCourse := resp.GetCourse()

	if pbCourse == nil {
		return GetCourseResponse{}
	}

	course := &Course{
		CourseID:    pbCourse.GetCourseId(),
		TeacherID:   pbCourse.GetTeacherId(),
		Title:       pbCourse.GetTitle(),
		Description: pbCourse.GetDescription(),
		Visibility:  pbCourse.GetVisibility(),
	}

	if pbCourse.GetStartTime() != nil {
		t := pbCourse.GetStartTime().AsTime()
		course.StartTime = &t
	}

	if pbCourse.GetEndTime() != nil {
		t := pbCourse.GetEndTime().AsTime()
		course.EndTime = &t
	}

	return GetCourseResponse{
		Course: course,
	}
}

type GetCoursesRequest struct {
	UserID string `json:"user_id"`
}

func NewGetCoursesRequest(req GetCoursesRequest) *pb.GetCoursesRequest {
	return &pb.GetCoursesRequest{
		UserId: req.UserID,
	}
}

type GetCoursesResponse struct {
	Courses []*pb.Course `json:"courses"`
}

func NewGetCoursesResponse(resp *pb.GetCoursesResponse) GetCoursesResponse {
	return GetCoursesResponse{
		Courses: resp.GetCourse(),
	}
}

type UpdateCourseRequest struct {
	CourseID    string     `json:"course_id"`
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Visibility  *bool      `json:"visibility,omitempty"`
	StartTime   *time.Time `json:"start_time,omitempty"`
	EndTime     *time.Time `json:"end_time,omitempty"`
}

func NewUpdateCourseRequest(req UpdateCourseRequest) *pb.UpdateCourseRequest {
	pbReq := &pb.UpdateCourseRequest{
		CourseId:    req.CourseID,
		Title:       req.Title,
		Description: req.Description,
		Visibility:  req.Visibility,
	}

	if req.StartTime != nil {
		pbReq.StartTime = timestamppb.New(*req.StartTime)
	}
	if req.EndTime != nil {
		pbReq.EndTime = timestamppb.New(*req.EndTime)
	}

	return pbReq
}

type UpdateCourseResponse struct{}

func NewUpdateCourseResponse(resp *pb.UpdateCourseResponse) UpdateCourseResponse {
	return UpdateCourseResponse{}
}

type DeleteCourseRequest struct {
	CourseID string `json:"course_id"`
}

func NewDeleteCourseRequest(req DeleteCourseRequest) *pb.DeleteCourseRequest {
	return &pb.DeleteCourseRequest{
		CourseId: req.CourseID,
	}
}

type DeleteCourseResponse struct{}

func NewDeleteCourseResponse(resp *pb.DeleteCourseResponse) DeleteCourseResponse {
	return DeleteCourseResponse{}
}

type EnrollUserRequest struct {
	CourseID string `json:"course_id"`
	UserID   string `json:"user_id"`
}

func NewEnrollUserRequest(req EnrollUserRequest) *pb.EnrollUserRequest {
	return &pb.EnrollUserRequest{
		CourseId: req.CourseID,
		UserId:   req.UserID,
	}
}

type EnrollUserResponse struct{}

func NewEnrollUserResponse(resp *pb.EnrollUserResponse) EnrollUserResponse {
	return EnrollUserResponse{}
}
