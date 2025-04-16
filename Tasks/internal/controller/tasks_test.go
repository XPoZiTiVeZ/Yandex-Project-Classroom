package controller_test

import (
	"Classroom/Tasks/internal/controller"
	mocks "Classroom/Tasks/internal/controller/mocks"
	"Classroom/Tasks/internal/domain"
	"Classroom/Tasks/internal/dto"
	pb "Classroom/Tasks/pkg/api/tasks"
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

func TestTaskController_CreateTask(t *testing.T) {
	type MockBehavior func(svc *mocks.MockTaskService, req *pb.CreateTaskRequest)

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		req          *pb.CreateTaskRequest
		want         *pb.CreateTaskResponse
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(svc *mocks.MockTaskService, req *pb.CreateTaskRequest) {
				svc.EXPECT().Create(context.Background(), dto.CreateTaskDTO{
					Title:    req.Title,
					Content:  req.Description,
					CourseID: req.CourseId,
				}).Return("task-id", nil)
			},
			req: &pb.CreateTaskRequest{
				CourseId:    uuid.NewString(),
				Title:       "title",
				Description: "description",
			},
			want: &pb.CreateTaskResponse{
				TaskId: "task-id",
			},
		},
		{
			name:         "invalid request",
			mockBehavior: func(svc *mocks.MockTaskService, req *pb.CreateTaskRequest) {},
			req: &pb.CreateTaskRequest{
				CourseId:    "not uuid",
				Title:       "title",
				Description: "description",
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid request: Key: 'CreateTaskDTO.CourseID' Error:Field validation for 'CourseID' failed on the 'uuid' tag"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocks.NewMockTaskService(t)
			tc.mockBehavior(svc, tc.req)
			c := controller.NewTaskController(slog.Default(), svc)
			got, err := c.CreateTask(context.Background(), tc.req)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestTaskController_GetTask(t *testing.T) {
	type MockBehavior func(svc *mocks.MockTaskService, req *pb.GetTaskRequest)

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		req          *pb.GetTaskRequest
		want         *pb.GetTaskResponse
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(svc *mocks.MockTaskService, req *pb.GetTaskRequest) {
				svc.EXPECT().GetTaskByID(mock.Anything, req.TaskId).Return(domain.Task{
					ID:       req.TaskId,
					Title:    "title",
					Content:  "content",
					CourseID: "course-id",
				}, nil)
			},
			req: &pb.GetTaskRequest{
				TaskId: uuid.NewString(),
			},
			want: &pb.GetTaskResponse{
				Task: &pb.Task{
					Title:    "title",
					Content:  "content",
					CourseId: "course-id",
				},
			},
		},
		{
			name:         "invalid request",
			mockBehavior: func(svc *mocks.MockTaskService, req *pb.GetTaskRequest) {},
			req: &pb.GetTaskRequest{
				TaskId: "not uuid",
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid task id"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocks.NewMockTaskService(t)
			tc.mockBehavior(svc, tc.req)
			c := controller.NewTaskController(slog.Default(), svc)
			got, err := c.GetTask(context.Background(), tc.req)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want.Task.Title, got.Task.Title)
			assert.Equal(t, tc.want.Task.Content, got.Task.Content)
			assert.Equal(t, tc.want.Task.CourseId, got.Task.CourseId)
		})
	}
}

func TestTaskController_UpdateTask(t *testing.T) {
	type MockBehavior func(svc *mocks.MockTaskService, req *pb.UpdateTaskRequest)

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		req          *pb.UpdateTaskRequest
		want         *pb.UpdateTaskResponse
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(svc *mocks.MockTaskService, req *pb.UpdateTaskRequest) {
				svc.EXPECT().Update(context.Background(), dto.UpdateTaskDTO{
					TaskID:  req.TaskId,
					Title:   req.Title,
					Content: req.Content,
				}).Return(domain.Task{
					Title:   *req.Title,
					Content: *req.Content,
				}, nil)
			},
			req: &pb.UpdateTaskRequest{
				Title:   strPtr("test"),
				Content: strPtr("test"),
				TaskId:  uuid.NewString(),
			},
			want: &pb.UpdateTaskResponse{
				Task: &pb.Task{
					Title:   "test",
					Content: "test",
				},
			},
		},
		{
			name:         "invalid request",
			mockBehavior: func(svc *mocks.MockTaskService, req *pb.UpdateTaskRequest) {},
			req: &pb.UpdateTaskRequest{
				Title:   strPtr("test"),
				Content: strPtr("test"),
				TaskId:  "not uuid",
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid request: Key: 'UpdateTaskDTO.TaskID' Error:Field validation for 'TaskID' failed on the 'uuid' tag"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocks.NewMockTaskService(t)
			tc.mockBehavior(svc, tc.req)
			c := controller.NewTaskController(slog.Default(), svc)
			got, err := c.UpdateTask(context.Background(), tc.req)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want.Task.Title, got.Task.Title)
			assert.Equal(t, tc.want.Task.Content, got.Task.Content)
		})
	}
}

func TestTaskController_ChangeStatusTask(t *testing.T) {
	type MockBehavior func(svc *mocks.MockTaskService, req *pb.ChangeStatusTaskRequest)

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		req          *pb.ChangeStatusTaskRequest
		want         *pb.ChangeStatusTaskResponse
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(svc *mocks.MockTaskService, req *pb.ChangeStatusTaskRequest) {
				svc.EXPECT().ToggleTaskStatus(context.Background(), req.TaskId, req.StudentId).Return(domain.TaskStatus{
					Completed: true,
					TaskID:    req.TaskId,
					UserID:    req.StudentId,
				}, nil)
			},
			req: &pb.ChangeStatusTaskRequest{
				TaskId:    uuid.NewString(),
				StudentId: uuid.NewString(),
			},
			want: &pb.ChangeStatusTaskResponse{
				TaskStatus: true,
			},
		},
		{
			name:         "invalid request",
			mockBehavior: func(svc *mocks.MockTaskService, req *pb.ChangeStatusTaskRequest) {},
			req: &pb.ChangeStatusTaskRequest{
				TaskId:    "not uuid",
				StudentId: uuid.NewString(),
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid task id"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocks.NewMockTaskService(t)
			tc.mockBehavior(svc, tc.req)
			c := controller.NewTaskController(slog.Default(), svc)
			got, err := c.ChangeStatusTask(context.Background(), tc.req)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want.TaskStatus, got.TaskStatus)
		})
	}
}

func strPtr(s string) *string {
	return &s
}
