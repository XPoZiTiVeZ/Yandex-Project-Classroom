package courses

import (
	"time"

	pb "Classroom/Gateway/pkg/api/courses"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Course struct {
	CourseID    string     `json:"course_id"`
	TeacherID   string     `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Visibility  bool       `json:"visibility"`
	StartTime   *time.Time `json:"start_time,omitempty"`
	EndTime     *time.Time `json:"end_time,omitempty"`
}

func NewCourse(pbCourse *pb.Course) Course {
	course := Course{
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

	return course
}

type Student struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Enrollment struct {
	CourseID string `json:"course_id"`
	StudentID string `json:"student_id"`
	EnrolledAt time.Time `json:"enrolled_at"`
}

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
		t := timestamppb.New(*req.StartTime)
		pbReq.StartTime = t
	}

	if req.EndTime != nil {
		t := timestamppb.New(*req.EndTime)
		pbReq.EndTime = t
	}

	return pbReq
}

type CreateCourseResponse struct {
	Course
}

func NewCreateCourseResponse(resp *pb.CreateCourseResponse) CreateCourseResponse {
	course := NewCourse(resp.GetCourse())

	return CreateCourseResponse{
		Course: course,
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

type GetCourseResponse struct {
	Course
}

func NewGetCourseResponse(resp *pb.GetCourseResponse) GetCourseResponse {
	course := NewCourse(resp.GetCourse())

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

type GetCoursesByStudentRequest struct {
	StudentId string `json:"student_id"`
}

func NewGetCoursesByStudentRequest(req GetCoursesByStudentRequest) *pb.GetCoursesByStudentRequest {
	return &pb.GetCoursesByStudentRequest{
		StudentId: req.StudentId,
	}
}

type GetCoursesByTeacherRequest struct {
	TeacherID string `json:"teacher_id"`
}

func NewGetCoursesByTeacherRequest(req GetCoursesByTeacherRequest) *pb.GetCoursesByTeacherRequest {
	return &pb.GetCoursesByTeacherRequest{
		TeacherId: req.TeacherID,
	}
}

type GetCoursesResponse struct {
	Courses []Course `json:"courses"`
}

func NewGetCoursesResponse(resp *pb.GetCoursesResponse) GetCoursesResponse {
	var courses []Course
	for _, pbCourse := range resp.GetCourses() {
		course := NewCourse(pbCourse)
		courses = append(courses, course)
	}

	return GetCoursesResponse{
		Courses: courses,
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
		t := timestamppb.New(*req.StartTime)
		pbReq.StartTime = t
	}
	if req.EndTime != nil {
		t := timestamppb.New(*req.EndTime)
		pbReq.EndTime = t
	}

	return pbReq
}

type UpdateCourseResponse struct {
	Course
}

func NewUpdateCourseResponse(resp *pb.UpdateCourseResponse) UpdateCourseResponse {
	course := NewCourse(resp.GetCourse())

	return UpdateCourseResponse{
		Course: course,
	}
}

type DeleteCourseRequest struct {
	CourseID string `json:"course_id"`
}

func NewDeleteCourseRequest(req DeleteCourseRequest) *pb.DeleteCourseRequest {
	return &pb.DeleteCourseRequest{
		CourseId: req.CourseID,
	}
}

type DeleteCourseResponse struct {
	Course
}

func NewDeleteCourseResponse(resp *pb.DeleteCourseResponse) DeleteCourseResponse {
	course := NewCourse(resp.GetCourse())

	return DeleteCourseResponse{
		Course: course,
	}
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

type EnrollUserResponse struct {
	Enrollment
}

func NewEnrollUserResponse(resp *pb.EnrollUserResponse) EnrollUserResponse {
	return EnrollUserResponse{
		Enrollment: Enrollment{
			CourseID:   resp.GetEnrollment().GetCourseId(),
			StudentID:  resp.GetEnrollment().GetStudentId(),
			EnrolledAt: resp.GetEnrollment().GetEnrolledAt().AsTime(),
		},
	}
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

type ExpelUserResponse struct {
	Enrollment
}

func NewExpelUserResponse(resp *pb.ExpelUserResponse) ExpelUserResponse {
	return ExpelUserResponse{
		Enrollment: Enrollment{
			CourseID:   resp.GetEnrollment().GetCourseId(),
			StudentID:  resp.GetEnrollment().GetStudentId(),
			EnrolledAt: resp.GetEnrollment().GetEnrolledAt().AsTime(),
		},
	}
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

type GetCourseStudentsRequest struct {
	CourseID string `json:"course_id"`
	Index    int32  `json:"index"`
	Limit    int32  `json:"limit"`
}

func NewGetCourseStudentsRequest(req GetCourseStudentsRequest) *pb.GetCourseStudentsRequest {
	return &pb.GetCourseStudentsRequest{
		CourseId: req.CourseID,
		Index:    req.Index,
		Limit:    req.Index,
	}
}

type GetCourseStudentsResponse struct {
	Index   int32    `json:"index"`
	Total   int32    `json:"total"`
	Students []Student `json:"students"`
}

func NewGetCourseStudentsResponse(resp *pb.GetCourseStudentsResponse) GetCourseStudentsResponse {
	return GetCourseStudentsResponse{
		Index: resp.GetIndex(),
		Total: resp.GetTotal(),
		Students: func() []Student {
			var students []Student
			for _, m := range resp.GetStudents() {
				students = append(students, Student{
					UserID:    m.GetUserId(),
					Email:     m.GetEmail(),
					FirstName: m.GetFirstName(),
					LastName:  m.GetLastName(),
				})
			}
			
			return students
		}(),
	}
}
