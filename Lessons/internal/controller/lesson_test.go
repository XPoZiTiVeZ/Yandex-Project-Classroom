package controller_test

import (
	"Classroom/Lessons/internal/controller"
	mocks "Classroom/Lessons/internal/controller/mocks"
	"Classroom/Lessons/internal/domain"
	"Classroom/Lessons/internal/dto"
	pb "Classroom/Lessons/pkg/api/lessons"
	"context"
	"log/slog"

	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLessonController_CreateLesson(t *testing.T) {
	type MockBehavior func(svc *mocks.MockLessonService, req *pb.CreateLessonRequest)

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		req          *pb.CreateLessonRequest
		want         *pb.CreateLessonResponse
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(svc *mocks.MockLessonService, req *pb.CreateLessonRequest) {
				svc.EXPECT().Create(mock.Anything, dto.CreateLessonDTO{
					Title:    req.Title,
					CourseID: req.CourseId,
					Content:  req.Content,
				}).Return(domain.Lesson{
					ID:       "lesson-id",
					Title:    req.Title,
					CourseID: req.CourseId,
					Content:  req.Content,
				}, nil)
			},
			req: &pb.CreateLessonRequest{
				Title:    "title",
				Content:  "content",
				CourseId: uuid.NewString(),
			},
			want: &pb.CreateLessonResponse{
				LessonId: "lesson-id",
			},
		},
		{
			name:         "invalid request",
			mockBehavior: func(svc *mocks.MockLessonService, req *pb.CreateLessonRequest) {},
			req: &pb.CreateLessonRequest{
				Title:    "title",
				Content:  "content",
				CourseId: "not uuid",
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid request: Key: 'CreateLessonDTO.CourseID' Error:Field validation for 'CourseID' failed on the 'uuid' tag"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocks.NewMockLessonService(t)
			tc.mockBehavior(svc, tc.req)
			c := controller.NewLessonController(slog.Default(), svc)
			got, err := c.CreateLesson(context.Background(), tc.req)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestLessonController_UpdateLesson(t *testing.T) {
	type MockBehavior func(svc *mocks.MockLessonService, req *pb.UpdateLessonRequest)

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		req          *pb.UpdateLessonRequest
		want         *pb.UpdateLessonResponse
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(svc *mocks.MockLessonService, req *pb.UpdateLessonRequest) {
				svc.EXPECT().Update(mock.Anything, dto.UpdateLessonDTO{
					LessonID: req.LessonId,
					Title:    req.Title,
					Content:  req.Content,
				}).Return(domain.Lesson{
					ID:      req.LessonId,
					Title:   *req.Title,
					Content: *req.Content,
				}, nil)
			},
			req: &pb.UpdateLessonRequest{
				LessonId: uuid.NewString(),
				Title:    strPtr("title"),
				Content:  strPtr("content"),
			},
			want: &pb.UpdateLessonResponse{
				Lesson: &pb.Lesson{
					Title:   "title",
					Content: "content",
				},
			},
		},
		{
			name:         "invalid request",
			mockBehavior: func(svc *mocks.MockLessonService, req *pb.UpdateLessonRequest) {},
			req: &pb.UpdateLessonRequest{
				LessonId: "not uuid",
				Title:    strPtr("title"),
				Content:  strPtr("content"),
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid request: Key: 'UpdateLessonDTO.LessonID' Error:Field validation for 'LessonID' failed on the 'uuid' tag"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocks.NewMockLessonService(t)
			tc.mockBehavior(svc, tc.req)
			c := controller.NewLessonController(slog.Default(), svc)
			got, err := c.UpdateLesson(context.Background(), tc.req)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want.Lesson.Content, got.Lesson.Content)
			assert.Equal(t, tc.want.Lesson.Title, got.Lesson.Title)
		})
	}
}

func TestLessonController_DeleteLesson(t *testing.T) {
	type MockBehavior func(svc *mocks.MockLessonService, req *pb.DeleteLessonRequest)

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		req          *pb.DeleteLessonRequest
		want         *pb.DeleteLessonResponse
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(svc *mocks.MockLessonService, req *pb.DeleteLessonRequest) {
				svc.EXPECT().Delete(mock.Anything, req.LessonId).Return(nil)
			},
			req: &pb.DeleteLessonRequest{
				LessonId: uuid.NewString(),
			},
			want: &pb.DeleteLessonResponse{Success: true},
		},
		{
			name:         "invalid request",
			mockBehavior: func(svc *mocks.MockLessonService, req *pb.DeleteLessonRequest) {},
			req: &pb.DeleteLessonRequest{
				LessonId: "not uuid",
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid lesson id"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocks.NewMockLessonService(t)
			tc.mockBehavior(svc, tc.req)
			c := controller.NewLessonController(slog.Default(), svc)
			got, err := c.DeleteLesson(context.Background(), tc.req)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func strPtr(s string) *string {
	return &s
}
