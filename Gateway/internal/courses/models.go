package courses

import (
	"time"

	pb "Classroom/Gateway/pkg/api/courses"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// Course - основная информация о курсе
// @Description Полная информация о курсе включая временные метки
type Course struct {
    // Уникальный идентификатор курса
    CourseID string `json:"course_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
    // ID преподавателя (владельца курса)
    TeacherID string `json:"teacher_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=1"`
    // Название курса
    Title string `json:"title" example:"Основы программирования" extensions:"x-order=2"`
    // Подробное описание курса
    Description string `json:"description" example:"Базовый курс по алгоритмам и структурам данных" extensions:"x-order=3"`
    // Видимость курса для студентов
    Visibility bool `json:"visibility" example:"true" extensions:"x-order=4"`
    // Дата начала курса (опционально)
    StartTime *time.Time `json:"start_time,omitempty" example:"2023-09-01T00:00:00Z" extensions:"x-order=5"`
    // Дата окончания курса (опционально)
    EndTime *time.Time `json:"end_time,omitempty" example:"2023-12-31T23:59:59Z" extensions:"x-order=6"`
} // @name Course

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

// Student - информация о студенте
// @Description Основные данные студента для отображения в списках курса
type Student struct {
    // Уникальный идентификатор пользователя
    UserID string `json:"user_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
    // Email студента
    Email string `json:"email" example:"student@example.com" extensions:"x-order=1"`
    // Имя студента
    FirstName string `json:"first_name" example:"Алексей" extensions:"x-order=2"`
    // Фамилия студента
    LastName string `json:"last_name" example:"Петров" extensions:"x-order=3"`
} // @name CourseStudent

// Enrollment - запись о зачислении студента
// @Description Информация о подписке студента на курс
type Enrollment struct {
    // ID курса
    CourseID string `json:"course_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
    // ID студента
    StudentID string `json:"student_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=1"`
    // Дата и время зачисления
    EnrolledAt time.Time `json:"enrolled_at" example:"2023-09-01T12:00:00Z" extensions:"x-order=2"`
} // @name CourseEnrollment

// CreateCourseRequest - запрос на создание курса
// @Description Параметры для создания нового курса
type CreateCourseRequest struct {
    // ID преподавателя
    UserID string `json:"user_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0" swaggerignore:"true"`
    // Название курса
    Title string `json:"title" example:"Основы программирования" extensions:"x-order=1"`
    // Описание курса
    Description string `json:"description" example:"Базовый курс по алгоритмам" extensions:"x-order=2"`
    // Видимость курса
    Visibility bool `json:"visibility" example:"true" extensions:"x-order=3"`
    // Дата начала (опционально)
    StartTime *time.Time `json:"start_time,omitempty" example:"2023-09-01T00:00:00Z" extensions:"x-order=4"`
    // Дата окончания (опционально)
    EndTime *time.Time `json:"end_time,omitempty" example:"2023-12-31T23:59:59Z" extensions:"x-order=5"`
} // @name CreateCourseRequest

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

// CreateCourseResponse - ответ после создания курса
// @Description Возвращает созданный курс
type CreateCourseResponse struct {
    Course Course `json:"course" extensions:"x-order=0"`
} // @name CreateCourseResponse

func NewCreateCourseResponse(resp *pb.CreateCourseResponse) CreateCourseResponse {
	course := NewCourse(resp.GetCourse())

	return CreateCourseResponse{
		Course: course,
	}
}

// GetCourseRequest - запрос информации о курсе
// @Description Требует ID курса для получения данных
type GetCourseRequest struct {
    // ID запрашиваемого курса
    CourseID string `schema:"course_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
	UserID string
} // @name GetCourseRequest

func NewGetCourseRequest(req GetCourseRequest) *pb.GetCourseRequest {
	return &pb.GetCourseRequest{
		CourseId: req.CourseID,
		UserId: req.UserID,
	}
}

// GetCourseResponse - информация о курсе
// @Description Возвращает полные данные курса
type GetCourseResponse struct {
    Course Course `schema:"course" extensions:"x-order=0"`
} // @name GetCourseResponse

func NewGetCourseResponse(resp *pb.GetCourseResponse) GetCourseResponse {
	course := NewCourse(resp.GetCourse())

	return GetCourseResponse{
		Course: course,
	}
}

// GetCoursesRequest - запрос списка курсов
// @Description Может фильтровать по ID пользователя
type GetCoursesRequest struct {
    // ID пользователя (опционально)
    UserID string `schema:"user_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
} // @name GetCoursesRequest

func NewGetCoursesRequest(req GetCoursesRequest) *pb.GetCoursesRequest {
	return &pb.GetCoursesRequest{
		UserId: req.UserID,
	}
}

// GetCoursesByStudentRequest - запрос курсов студента
// @Description Возвращает курсы, на которые подписан студент
type GetCoursesByStudentRequest struct {
    // ID студента
    StudentId string `schema:"student_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
} // @name GetCoursesByStudentRequest

func NewGetCoursesByStudentRequest(req GetCoursesByStudentRequest) *pb.GetCoursesByStudentRequest {
	return &pb.GetCoursesByStudentRequest{
		StudentId: req.StudentId,
	}
}

// GetCoursesByTeacherRequest - запрос курсов преподавателя
// @Description Возвращает курсы, созданные преподавателем
type GetCoursesByTeacherRequest struct {
    // ID преподавателя
    TeacherID string `schema:"teacher_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
} // @name GetCoursesByTeacherRequest

func NewGetCoursesByTeacherRequest(req GetCoursesByTeacherRequest) *pb.GetCoursesByTeacherRequest {
	return &pb.GetCoursesByTeacherRequest{
		TeacherId: req.TeacherID,
	}
}

// GetCoursesResponse - список курсов
// @Description Содержит массив курсов
type GetCoursesResponse struct {
    // Массив курсов
    Courses []Course `json:"courses" extensions:"x-order=0"`
} // @name GetCoursesResponse

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

// UpdateCourseRequest - запрос обновления курса
// @Description Позволяет частично обновить данные курса
type UpdateCourseRequest struct {
    // ID обновляемого курса
    CourseID string `json:"course_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
    // Новое название (опционально)
    Title *string `json:"title,omitempty" example:"Продвинутое программирование" extensions:"x-order=1"`
    // Новое описание (опционально)
    Description *string `json:"description,omitempty" example:"Расширенный курс" extensions:"x-order=2"`
    // Новая видимость (опционально)
    Visibility *bool `json:"visibility,omitempty" example:"false" extensions:"x-order=3"`
    // Новое время начала (опционально)
    StartTime *time.Time `json:"start_time,omitempty" example:"2024-01-01T00:00:00Z" extensions:"x-order=4"`
    // Новое время окончания (опционально)
    EndTime *time.Time `json:"end_time,omitempty" example:"2024-06-30T23:59:59Z" extensions:"x-order=5"`
} // @name UpdateCourseRequest

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

// UpdateCourseResponse - результат обновления
// @Description Возвращает обновленный курс
type UpdateCourseResponse struct {
    Course Course `json:"course" extensions:"x-order=0"`
} // @name UpdateCourseResponse

func NewUpdateCourseResponse(resp *pb.UpdateCourseResponse) UpdateCourseResponse {
	course := NewCourse(resp.GetCourse())

	return UpdateCourseResponse{
		Course: course,
	}
}

// DeleteCourseRequest - запрос удаления курса
// @Description Требует ID курса для удаления
type DeleteCourseRequest struct {
    // ID удаляемого курса
    CourseID string `json:"course_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
} // @name DeleteCourseRequest

func NewDeleteCourseRequest(req DeleteCourseRequest) *pb.DeleteCourseRequest {
	return &pb.DeleteCourseRequest{
		CourseId: req.CourseID,
	}
}

// DeleteCourseResponse - результат удаления
// @Description Возвращает удаленный курс
type DeleteCourseResponse struct {
    Course Course `json:"course" extensions:"x-order=0"`
} // @name DeleteCourseResponse

func NewDeleteCourseResponse(resp *pb.DeleteCourseResponse) DeleteCourseResponse {
	course := NewCourse(resp.GetCourse())

	return DeleteCourseResponse{
		Course: course,
	}
}

// EnrollUserRequest - запрос зачисления студента
// @Description Добавляет студента на курс
type EnrollUserRequest struct {
    // ID курса
    CourseID string `json:"course_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
    // ID студента
    UserID string `json:"user_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=1"`
} // @name EnrollUserRequest

func NewEnrollUserRequest(req EnrollUserRequest) *pb.EnrollUserRequest {
	return &pb.EnrollUserRequest{
		CourseId: req.CourseID,
		UserId:   req.UserID,
	}
}

// EnrollUserResponse - результат зачисления
// @Description Подтверждение успешного зачисления
type EnrollUserResponse struct {
    Enrollment Enrollment `json:"enrollment" extensions:"x-order=0"`
} // @name EnrollUserResponse

func NewEnrollUserResponse(resp *pb.EnrollUserResponse) EnrollUserResponse {
	return EnrollUserResponse{
		Enrollment: Enrollment{
			CourseID:   resp.GetEnrollment().GetCourseId(),
			StudentID:  resp.GetEnrollment().GetStudentId(),
			EnrolledAt: resp.GetEnrollment().GetEnrolledAt().AsTime(),
		},
	}
}

// ExpelUserRequest - запрос отчисления студента
// @Description Удаляет студента с курса
type ExpelUserRequest struct {
    // ID курса
    CourseID string `json:"course_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
    // ID студента
    UserID string `json:"user_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=1"`
} // @name ExpelUserRequest

func NewExpelUserRequest(req ExpelUserRequest) *pb.ExpelUserRequest {
	return &pb.ExpelUserRequest{
		CourseId: req.CourseID,
		UserId:   req.UserID,
	}
}

// ExpelUserResponse - результат отчисления
// @Description Подтверждение успешного отчисления
type ExpelUserResponse struct {
    Enrollment Enrollment `json:"enrollment" extensions:"x-order=0"`
} // @name ExpelUserResponse

func NewExpelUserResponse(resp *pb.ExpelUserResponse) ExpelUserResponse {
	return ExpelUserResponse{
		Enrollment: Enrollment{
			CourseID:   resp.GetEnrollment().GetCourseId(),
			StudentID:  resp.GetEnrollment().GetStudentId(),
			EnrolledAt: resp.GetEnrollment().GetEnrolledAt().AsTime(),
		},
	}
}

// IsTeacherRequest - проверка роли преподавателя
// @Description Проверяет, является ли пользователь преподавателем курса
type IsTeacherRequest struct {
    // ID пользователя
    UserID string `json:"user_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
    // ID курса
    CourseID string `json:"course_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=1"`
} // @name IsTeacherRequest

func NewIsTeacherRequest(req IsTeacherRequest) *pb.IsTeacherRequest {
	return &pb.IsTeacherRequest{
		UserId:   req.UserID,
		CourseId: req.CourseID,
	}
}

// IsTeacherResponse - результат проверки
// @Description Возвращает true если пользователь - преподаватель
type IsTeacherResponse struct {
    // Результат проверки
    IsTeacher bool `json:"is_teacher" example:"true" extensions:"x-order=0"`
} // @name IsTeacherResponse

func NewIsTeacherResponse(resp *pb.IsTeacherResponse) IsTeacherResponse {
	return IsTeacherResponse{
		IsTeacher: resp.GetIsTeacher(),
	}
}

// IsMemberRequest - проверка участия в курсе
// @Description Проверяет, является ли пользователь участником курса
type IsMemberRequest struct {
    // ID пользователя
    UserID string `json:"user_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
    // ID курса
    CourseID string `json:"course_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=1"`
} // @name IsMemberRequest

func NewIsMemberRequest(req IsMemberRequest) *pb.IsMemberRequest {
	return &pb.IsMemberRequest{
		UserId:   req.UserID,
		CourseId: req.CourseID,
	}
}

// IsMemberResponse - результат проверки
// @Description Возвращает true если пользователь - участник
type IsMemberResponse struct {
    // Результат проверки
    IsMember bool `json:"is_member" example:"true" extensions:"x-order=0"`
} // @name IsMemberResponse

func NewIsMemberResponse(resp *pb.IsMemberResponse) IsMemberResponse {
	return IsMemberResponse{
		IsMember: resp.GetIsMember(),
	}
}

// GetCourseStudentsRequest - запрос списка студентов
// @Description Возвращает студентов курса с пагинацией
type GetCourseStudentsRequest struct {
    // ID курса
    CourseID string `schema:"course_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
    // Начальный индекс (для пагинации)
    Index int32 `schema:"index" example:"0" extensions:"x-order=1"`
    // Лимит записей (для пагинации)
    Limit int32 `schema:"limit" example:"10" extensions:"x-order=2"`
} // @name GetCourseStudentsRequest

func NewGetCourseStudentsRequest(req GetCourseStudentsRequest) *pb.GetCourseStudentsRequest {
	return &pb.GetCourseStudentsRequest{
		CourseId: req.CourseID,
		Index:    req.Index,
		Limit:    req.Index,
	}
}

// GetCourseStudentsResponse - список студентов
// @Description Содержит массив студентов с информацией о пагинации
type GetCourseStudentsResponse struct {
    // Текущий индекс
    Index int32 `json:"index" example:"0" extensions:"x-order=0"`
    // Общее количество студентов
    Total int32 `json:"total" example:"100" extensions:"x-order=1"`
    // Массив студентов
    Students []Student `json:"students" extensions:"x-order=2"`
} // @name GetCourseStudentsResponse

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
