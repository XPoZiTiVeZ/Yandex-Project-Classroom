package service_test

import (
	"Classroom/Courses/internal/domain"
	"Classroom/Courses/internal/dto"
	"Classroom/Courses/internal/service"
	mocks "Classroom/Courses/internal/service/mocks"
	pb "Classroom/Courses/pkg/api/courses"
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func TestCoursesService_CreateCourse(t *testing.T) {
	type MockBehavior func(svc *mocks.MockCourseRepo, req *pb.CreateCourseRequest)

	now := time.Now()
	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		req          *pb.CreateCourseRequest
		want         *pb.CreateCourseResponse
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(svc *mocks.MockCourseRepo, req *pb.CreateCourseRequest) {
				svc.EXPECT().Create(mock.Anything, dto.CreateCourseDTO{
					TeacherID:   req.UserId,
					Title:       req.Title,
					Description: req.Description,
				}).Return(domain.Course{ID: "course-id", Title: req.Title, TeacherID: "user-id", CreatedAt: now}, nil)
			},
			req: &pb.CreateCourseRequest{
				Title:       "title",
				Description: "description",
				UserId:      uuid.NewString(),
			},
			want: &pb.CreateCourseResponse{
				Course: &pb.Course{
					CourseId:  "course-id",
					Title:     "title",
					TeacherId: "user-id",
					CreatedAt: timestamppb.New(now),
				},
			},
		},
		{
			name:         "invalid request",
			mockBehavior: func(svc *mocks.MockCourseRepo, req *pb.CreateCourseRequest) {},
			req: &pb.CreateCourseRequest{
				Title:       "title",
				Description: "description",
				UserId:      "not uuid",
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid request: Key: 'CreateCourseDTO.TeacherID' Error:Field validation for 'TeacherID' failed on the 'uuid' tag"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewMockCourseRepo(t)
			svc := service.NewCoursesService(slog.Default(), repo)
			tc.mockBehavior(repo, tc.req)
			got, err := svc.CreateCourse(context.Background(), tc.req)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestCoursesService_GetCourse(t *testing.T) {
	type MockBehavior func(svc *mocks.MockCourseRepo, req *pb.GetCourseRequest)

	courseID := uuid.NewString()
	userID := uuid.NewString()
	now := time.Now()

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		req          *pb.GetCourseRequest
		want         *pb.GetCourseResponse
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(svc *mocks.MockCourseRepo, req *pb.GetCourseRequest) {
				svc.EXPECT().GetByID(mock.Anything, req.CourseId).
					Return(domain.Course{ID: req.CourseId, TeacherID: req.UserId, Visibility: true, CreatedAt: now}, nil)
			},
			req: &pb.GetCourseRequest{
				CourseId: courseID,
				UserId:   userID,
			},
			want: &pb.GetCourseResponse{
				Course: &pb.Course{
					CourseId:   courseID,
					TeacherId:  userID,
					Visibility: true,
					CreatedAt:  timestamppb.New(now),
				},
			},
		},
		{
			name: "course not found",
			mockBehavior: func(svc *mocks.MockCourseRepo, req *pb.GetCourseRequest) {
				svc.EXPECT().GetByID(mock.Anything, req.CourseId).
					Return(domain.Course{}, domain.ErrNotFound)
			},
			req: &pb.GetCourseRequest{
				CourseId: courseID,
				UserId:   userID,
			},
			wantErr: status.Error(codes.NotFound, "course not found"),
		},
		{
			name:         "invalid course id",
			mockBehavior: func(svc *mocks.MockCourseRepo, req *pb.GetCourseRequest) {},
			req: &pb.GetCourseRequest{
				CourseId: "invalid-uuid",
				UserId:   userID,
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid course id"),
		},
		{
			name:         "invalid user id",
			mockBehavior: func(svc *mocks.MockCourseRepo, req *pb.GetCourseRequest) {},
			req: &pb.GetCourseRequest{
				CourseId: courseID,
				UserId:   "invalid-uuid",
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid user id"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewMockCourseRepo(t)
			svc := service.NewCoursesService(slog.Default(), repo)
			tc.mockBehavior(repo, tc.req)
			got, err := svc.GetCourse(context.Background(), tc.req)

			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.wantErr.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
func TestCoursesService_UpdateCourse(t *testing.T) {
	type MockBehavior func(svc *mocks.MockCourseRepo, req *pb.UpdateCourseRequest)

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		req          *pb.UpdateCourseRequest
		want         *pb.UpdateCourseResponse
		wantErr      error
	}{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewMockCourseRepo(t)
			svc := service.NewCoursesService(slog.Default(), repo)
			tc.mockBehavior(repo, tc.req)
			got, err := svc.UpdateCourse(context.Background(), tc.req)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestCoursesService_DeleteCourse(t *testing.T) {
	type MockBehavior func(svc *mocks.MockCourseRepo, req *pb.DeleteCourseRequest)

	now := time.Now()
	courseID := uuid.NewString()
	teacherID := uuid.NewString()

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		req          *pb.DeleteCourseRequest
		want         *pb.DeleteCourseResponse
		wantErr      error
	}{
		{
			name: "success - course deleted",
			mockBehavior: func(svc *mocks.MockCourseRepo, req *pb.DeleteCourseRequest) {
				svc.EXPECT().Delete(mock.Anything, req.CourseId).
					Return(domain.Course{
						ID:          courseID,
						TeacherID:   teacherID,
						Title:       "course title",
						Description: "course description",
						Visibility:  true,
						CreatedAt:   now,
					}, nil)
			},
			req: &pb.DeleteCourseRequest{
				CourseId: courseID,
			},
			want: &pb.DeleteCourseResponse{
				Course: &pb.Course{
					CourseId:    courseID,
					TeacherId:   teacherID,
					Title:       "course title",
					Description: "course description",
					Visibility:  true,
					CreatedAt:   timestamppb.New(now),
				},
			},
		},
		{
			name: "error - course not found",
			mockBehavior: func(svc *mocks.MockCourseRepo, req *pb.DeleteCourseRequest) {
				svc.EXPECT().Delete(mock.Anything, req.CourseId).
					Return(domain.Course{}, domain.ErrNotFound)
			},
			req: &pb.DeleteCourseRequest{
				CourseId: courseID,
			},
			wantErr: status.Error(codes.NotFound, "course not found"),
		},
		{
			name:         "error - invalid course id",
			mockBehavior: func(svc *mocks.MockCourseRepo, req *pb.DeleteCourseRequest) {},
			req: &pb.DeleteCourseRequest{
				CourseId: "invalid-uuid",
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid course id"),
		},
		{
			name: "error - internal server error",
			mockBehavior: func(svc *mocks.MockCourseRepo, req *pb.DeleteCourseRequest) {
				svc.EXPECT().Delete(mock.Anything, req.CourseId).
					Return(domain.Course{}, errors.New("internal error"))
			},
			req: &pb.DeleteCourseRequest{
				CourseId: courseID,
			},
			wantErr: status.Error(codes.Internal, "failed to delete course"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewMockCourseRepo(t)
			svc := service.NewCoursesService(slog.Default(), repo)
			tc.mockBehavior(repo, tc.req)
			got, err := svc.DeleteCourse(context.Background(), tc.req)

			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.wantErr.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestCoursesService_EnrollUser(t *testing.T) {
	type MockBehavior func(svc *mocks.MockCourseRepo, req *pb.EnrollUserRequest)

	now := time.Now()
	courseID := uuid.NewString()
	studentID := uuid.NewString()

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		req          *pb.EnrollUserRequest
		want         *pb.EnrollUserResponse
		wantErr      error
	}{
		{
			name: "success - user enrolled",
			mockBehavior: func(svc *mocks.MockCourseRepo, req *pb.EnrollUserRequest) {
				svc.EXPECT().EnrollUser(mock.Anything, req.CourseId, req.UserId).
					Return(domain.Enrollment{
						CourseID:   req.CourseId,
						StudentID:  req.UserId,
						EnrolledAt: now,
					}, nil)
			},
			req: &pb.EnrollUserRequest{
				CourseId: courseID,
				UserId:   studentID,
			},
			want: &pb.EnrollUserResponse{
				Enrollment: &pb.Enrollment{
					CourseId:   courseID,
					StudentId:  studentID,
					EnrolledAt: timestamppb.New(now),
				},
			},
		},
		{
			name:         "error - invalid course id",
			mockBehavior: func(svc *mocks.MockCourseRepo, req *pb.EnrollUserRequest) {},
			req: &pb.EnrollUserRequest{
				CourseId: "invalid-uuid",
				UserId:   studentID,
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid course id"),
		},
		{
			name:         "error - invalid user id",
			mockBehavior: func(svc *mocks.MockCourseRepo, req *pb.EnrollUserRequest) {},
			req: &pb.EnrollUserRequest{
				CourseId: courseID,
				UserId:   "invalid-uuid",
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid user id"),
		},
		{
			name: "error - internal server error",
			mockBehavior: func(svc *mocks.MockCourseRepo, req *pb.EnrollUserRequest) {
				svc.EXPECT().EnrollUser(mock.Anything, req.CourseId, req.UserId).
					Return(domain.Enrollment{}, errors.New("internal error"))
			},
			req: &pb.EnrollUserRequest{
				CourseId: courseID,
				UserId:   studentID,
			},
			wantErr: status.Error(codes.Internal, "failed to enroll user"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewMockCourseRepo(t)
			svc := service.NewCoursesService(slog.Default(), repo)
			tc.mockBehavior(repo, tc.req)
			got, err := svc.EnrollUser(context.Background(), tc.req)

			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.wantErr.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestCoursesService_IsTeacher(t *testing.T) {
	type MockBehavior func(svc *mocks.MockCourseRepo, req *pb.IsTeacherRequest)

	courseID := uuid.NewString()
	teacherID := uuid.NewString()
	nonTeacherID := uuid.NewString()

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		req          *pb.IsTeacherRequest
		want         *pb.IsTeacherResponse
		wantErr      error
	}{
		{
			name: "success - user is teacher",
			mockBehavior: func(svc *mocks.MockCourseRepo, req *pb.IsTeacherRequest) {
				svc.EXPECT().IsTeacher(mock.Anything, req.CourseId, req.UserId).
					Return(true, nil)
			},
			req: &pb.IsTeacherRequest{
				CourseId: courseID,
				UserId:   teacherID,
			},
			want: &pb.IsTeacherResponse{
				IsTeacher: true,
			},
		},
		{
			name: "success - user is not teacher",
			mockBehavior: func(svc *mocks.MockCourseRepo, req *pb.IsTeacherRequest) {
				svc.EXPECT().IsTeacher(mock.Anything, req.CourseId, req.UserId).
					Return(false, nil)
			},
			req: &pb.IsTeacherRequest{
				CourseId: courseID,
				UserId:   nonTeacherID,
			},
			want: &pb.IsTeacherResponse{
				IsTeacher: false,
			},
		},
		{
			name:         "error - invalid course id",
			mockBehavior: func(svc *mocks.MockCourseRepo, req *pb.IsTeacherRequest) {},
			req: &pb.IsTeacherRequest{
				CourseId: "invalid-uuid",
				UserId:   teacherID,
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid course id"),
		},
		{
			name:         "error - invalid user id",
			mockBehavior: func(svc *mocks.MockCourseRepo, req *pb.IsTeacherRequest) {},
			req: &pb.IsTeacherRequest{
				CourseId: courseID,
				UserId:   "invalid-uuid",
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid user id"),
		},
		{
			name: "error - internal server error",
			mockBehavior: func(svc *mocks.MockCourseRepo, req *pb.IsTeacherRequest) {
				svc.EXPECT().IsTeacher(mock.Anything, req.CourseId, req.UserId).
					Return(false, errors.New("internal error"))
			},
			req: &pb.IsTeacherRequest{
				CourseId: courseID,
				UserId:   teacherID,
			},
			wantErr: status.Error(codes.Internal, "failed to check teacher"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewMockCourseRepo(t)
			svc := service.NewCoursesService(slog.Default(), repo)
			tc.mockBehavior(repo, tc.req)
			got, err := svc.IsTeacher(context.Background(), tc.req)

			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.wantErr.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestCoursesService_IsMember(t *testing.T) {
	type MockBehavior func(svc *mocks.MockCourseRepo, req *pb.IsMemberRequest)

	courseID := uuid.NewString()
	memberID := uuid.NewString()
	nonMemberID := uuid.NewString()
	teacherID := uuid.NewString()

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		req          *pb.IsMemberRequest
		want         *pb.IsMemberResponse
		wantErr      error
	}{
		{
			name: "success - user is member",
			mockBehavior: func(svc *mocks.MockCourseRepo, req *pb.IsMemberRequest) {
				svc.EXPECT().IsMember(mock.Anything, req.CourseId, req.UserId).
					Return(true, nil)
			},
			req: &pb.IsMemberRequest{
				CourseId: courseID,
				UserId:   memberID,
			},
			want: &pb.IsMemberResponse{
				IsMember: true,
			},
		},
		{
			name: "success - user is teacher (considered member)",
			mockBehavior: func(svc *mocks.MockCourseRepo, req *pb.IsMemberRequest) {
				svc.EXPECT().IsMember(mock.Anything, req.CourseId, req.UserId).
					Return(false, nil)
				svc.EXPECT().IsTeacher(mock.Anything, req.CourseId, req.UserId).
					Return(true, nil)
			},
			req: &pb.IsMemberRequest{
				CourseId: courseID,
				UserId:   teacherID,
			},
			want: &pb.IsMemberResponse{
				IsMember: true,
			},
		},
		{
			name: "success - user is not member",
			mockBehavior: func(svc *mocks.MockCourseRepo, req *pb.IsMemberRequest) {
				svc.EXPECT().IsMember(mock.Anything, req.CourseId, req.UserId).
					Return(false, nil)
				svc.EXPECT().IsTeacher(mock.Anything, req.CourseId, req.UserId).
					Return(false, nil)
			},
			req: &pb.IsMemberRequest{
				CourseId: courseID,
				UserId:   nonMemberID,
			},
			want: &pb.IsMemberResponse{
				IsMember: false,
			},
		},
		{
			name:         "error - invalid course id",
			mockBehavior: func(svc *mocks.MockCourseRepo, req *pb.IsMemberRequest) {},
			req: &pb.IsMemberRequest{
				CourseId: "invalid-uuid",
				UserId:   memberID,
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid course id"),
		},
		{
			name:         "error - invalid user id",
			mockBehavior: func(svc *mocks.MockCourseRepo, req *pb.IsMemberRequest) {},
			req: &pb.IsMemberRequest{
				CourseId: courseID,
				UserId:   "invalid-uuid",
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid user id"),
		},
		{
			name: "error - internal server error (isMember check)",
			mockBehavior: func(svc *mocks.MockCourseRepo, req *pb.IsMemberRequest) {
				svc.EXPECT().IsMember(mock.Anything, req.CourseId, req.UserId).
					Return(false, errors.New("internal error"))
			},
			req: &pb.IsMemberRequest{
				CourseId: courseID,
				UserId:   memberID,
			},
			wantErr: status.Error(codes.Internal, "failed to check member"),
		},
		{
			name: "error - internal server error (isTeacher check)",
			mockBehavior: func(svc *mocks.MockCourseRepo, req *pb.IsMemberRequest) {
				svc.EXPECT().IsMember(mock.Anything, req.CourseId, req.UserId).
					Return(false, nil)
				svc.EXPECT().IsTeacher(mock.Anything, req.CourseId, req.UserId).
					Return(false, errors.New("internal error"))
			},
			req: &pb.IsMemberRequest{
				CourseId: courseID,
				UserId:   teacherID,
			},
			wantErr: status.Error(codes.Internal, "failed to check teacher"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewMockCourseRepo(t)
			svc := service.NewCoursesService(slog.Default(), repo)
			tc.mockBehavior(repo, tc.req)
			got, err := svc.IsMember(context.Background(), tc.req)

			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.wantErr.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
