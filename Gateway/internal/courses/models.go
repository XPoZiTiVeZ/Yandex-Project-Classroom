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
	CourseID    string     `json:"course_id"`
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

type Member struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type GetCourseMembersRequest struct {
	CourseID string `json:"course_id"`
	Index    int32  `json:"index"`
}

func NewGetCourseMembersRequest(req GetCourseMembersRequest) *pb.GetCourseMembersRequest {
	return &pb.GetCourseMembersRequest{
		CourseId: req.CourseID,
		Index:    req.Index,
	}
}

type GetCourseMembersResponse struct {
	Index   int32    `json:"index"`
	Total   int32    `json:"total"`
	Members []Member `json:"members"`
}

func NewGetCourseMembersResponse(resp *pb.GetCourseMembersResponse) GetCourseMembersResponse {
	return GetCourseMembersResponse{
		Index: resp.GetIndex(),
		Total: resp.GetTotal(),
		Members: func() []Member {
			members := make([]Member, len(resp.GetMembers()))
			for i, m := range resp.GetMembers() {
				members[i] = Member{
					UserID:    m.GetUserId(),
					Email:     m.GetEmail(),
					FirstName: m.GetFirstName(),
					LastName:  m.GetLastName(),
				}
			}
			
			return members
		}(),
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

type ExpelUserRequest struct {
	CourseID string `json:"course_id"`
	UserID   string `json:"user_id"`
}

func NewExpelUserRequest(req ExpelUserRequest) *pb.ExpelUserRequest {
	return &pb.ExpelUserRequest{
		CourseId: req.CourseID,
		UserId:   req.UserID,
	}
}

type ExpelUserResponse struct{}

func NewExpelUserResponse(resp *pb.ExpelUserResponse) ExpelUserResponse {
	return ExpelUserResponse{}
}

type IsTeacherRequest struct {
	UserID   string `json:"user_id"`
	CourseID string `json:"course_id"`
}

func NewIsTeacherRequest(req IsTeacherRequest) *pb.IsTeacherRequest {
	return &pb.IsTeacherRequest{
		UserId:   req.UserID,
		CourseId: req.CourseID,
	}
}

type IsTeacherResponse struct {
	IsTeacher bool `json:"is_teacher"`
}

func NewIsTeacherResponse(resp *pb.IsTeacherResponse) IsTeacherResponse {
	return IsTeacherResponse{
		IsTeacher: resp.GetIsTeacher(),
	}
}

type IsMemberRequest struct {
	UserID   string `json:"user_id"`
	CourseID string `json:"course_id"`
}

func NewIsMemberRequest(req IsMemberRequest) *pb.IsMemberRequest {
	return &pb.IsMemberRequest{
		UserId:   req.UserID,
		CourseId: req.CourseID,
	}
}

type IsMemberResponse struct {
	IsMember bool `json:"is_member"`
}

func NewIsMemberResponse(resp *pb.IsMemberResponse) IsMemberResponse {
	return IsMemberResponse{
		IsMember: resp.GetIsMember(),
	}
}
