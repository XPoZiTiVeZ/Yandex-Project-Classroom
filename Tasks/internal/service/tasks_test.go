package service_test

import (
	"Classroom/Tasks/internal/domain"
	"Classroom/Tasks/internal/dto"
	"Classroom/Tasks/internal/service"
	mocks "Classroom/Tasks/internal/service/mocks"
	"Classroom/Tasks/pkg/events"
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTaskService_Create(t *testing.T) {
	type MockBehavior func(repo *mocks.MockTaskRepo, pr *mocks.MockProducer, payload dto.CreateTaskDTO)

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		payload      dto.CreateTaskDTO
		want         string
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(repo *mocks.MockTaskRepo, pr *mocks.MockProducer, payload dto.CreateTaskDTO) {
				repo.EXPECT().CourseExists(context.Background(), payload.CourseID).Return(true, nil)
				repo.EXPECT().Create(context.Background(), payload).Return(domain.Task{
					ID:       "task-id",
					Title:    payload.Title,
					Content:  payload.Content,
					CourseID: payload.CourseID,
				}, nil)

				pr.EXPECT().PublishTaskCreated(events.TaskCreated{CourseID: payload.CourseID, TaskID: "task-id"}).Return(nil)
			},
			payload: dto.CreateTaskDTO{
				CourseID: "course-id",
				Title:    "title",
				Content:  "content",
			},
			want: "task-id",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewMockTaskRepo(t)
			pr := mocks.NewMockProducer(t)
			tc.mockBehavior(repo, pr, tc.payload)
			svc := service.NewTaskService(slog.Default(), repo, nil, pr)
			got, err := svc.Create(context.Background(), tc.payload)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestTaskService_Update(t *testing.T) {
	type MockBehavior func(repo *mocks.MockTaskRepo, payload dto.UpdateTaskDTO)
	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		payload      dto.UpdateTaskDTO
		want         domain.Task
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(repo *mocks.MockTaskRepo, payload dto.UpdateTaskDTO) {
				task := domain.Task{
					ID:      payload.TaskID,
					Title:   "old",
					Content: "old",
				}
				repo.EXPECT().GetByID(mock.Anything, payload.TaskID).Return(task, nil)
				repo.EXPECT().Update(mock.Anything, domain.Task{
					ID:      payload.TaskID,
					Title:   *payload.Title,
					Content: *payload.Content,
				}).Return(nil)
			},
			payload: dto.UpdateTaskDTO{
				TaskID:  "task-id",
				Title:   strPtr("title"),
				Content: strPtr("content"),
			},
			want: domain.Task{
				ID:      "task-id",
				Title:   "title",
				Content: "content",
			},
		},
		{
			name: "task not found",
			mockBehavior: func(repo *mocks.MockTaskRepo, payload dto.UpdateTaskDTO) {
				repo.EXPECT().GetByID(mock.Anything, payload.TaskID).Return(domain.Task{}, domain.ErrNotFound)
			},
			payload: dto.UpdateTaskDTO{
				TaskID:  "task-id",
				Title:   strPtr("title"),
				Content: strPtr("content"),
			},
			wantErr: domain.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewMockTaskRepo(t)
			tc.mockBehavior(repo, tc.payload)
			svc := service.NewTaskService(slog.Default(), repo, nil, nil)
			got, err := svc.Update(context.Background(), tc.payload)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestTaskService_ToggleTaskStatus(t *testing.T) {
	type args struct {
		TaskID, UserID string
	}

	type MockBehavior func(tasks *mocks.MockTaskRepo, statuses *mocks.MockStatusRepo, args args)
	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		args         args
		want         domain.TaskStatus
		wantErr      error
	}{
		{
			name: "need to create status",
			mockBehavior: func(tasks *mocks.MockTaskRepo, statuses *mocks.MockStatusRepo, args args) {
				tasks.EXPECT().GetByID(mock.Anything, args.TaskID).Return(domain.Task{ID: args.TaskID}, nil)
				statuses.EXPECT().Get(mock.Anything, args.TaskID, args.UserID).Return(domain.TaskStatus{}, domain.ErrNotFound)
				statuses.EXPECT().Create(mock.Anything, domain.TaskStatus{
					TaskID:    args.TaskID,
					UserID:    args.UserID,
					Completed: true,
				}).Return(nil)
			},
			args: args{
				TaskID: "task-id",
				UserID: "user-id",
			},
			want: domain.TaskStatus{
				TaskID:    "task-id",
				UserID:    "user-id",
				Completed: true,
			},
		},
		{
			name: "need to update status",
			mockBehavior: func(tasks *mocks.MockTaskRepo, statuses *mocks.MockStatusRepo, args args) {
				tasks.EXPECT().GetByID(mock.Anything, args.TaskID).Return(domain.Task{ID: args.TaskID}, nil)
				statuses.EXPECT().Get(mock.Anything, args.TaskID, args.UserID).Return(domain.TaskStatus{
					TaskID:    args.TaskID,
					UserID:    args.UserID,
					Completed: false,
				}, nil)
				statuses.EXPECT().Update(mock.Anything, domain.TaskStatus{
					TaskID:    args.TaskID,
					UserID:    args.UserID,
					Completed: true,
				}).Return(nil)
			},
			args: args{
				TaskID: "task-id",
				UserID: "user-id",
			},
			want: domain.TaskStatus{
				TaskID:    "task-id",
				UserID:    "user-id",
				Completed: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tasks := mocks.NewMockTaskRepo(t)
			statuses := mocks.NewMockStatusRepo(t)
			tc.mockBehavior(tasks, statuses, tc.args)
			svc := service.NewTaskService(slog.Default(), tasks, statuses, nil)
			got, err := svc.ToggleTaskStatus(context.Background(), tc.args.TaskID, tc.args.UserID)
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
