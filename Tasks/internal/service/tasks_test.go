package service_test

import (
	"Classroom/Tasks/internal/domain"
	"Classroom/Tasks/internal/dto"
	"Classroom/Tasks/internal/service"
	mocks "Classroom/Tasks/internal/service/mocks"
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskService_Create(t *testing.T) {
	type MockBehavior func(repo *mocks.MockTaskRepo, payload dto.CreateTaskDTO)

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		payload      dto.CreateTaskDTO
		want         string
		wantErr      error
	}{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewMockTaskRepo(t)
			tc.mockBehavior(repo, tc.payload)
			svc := service.NewTaskService(slog.Default(), repo, nil)
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
	}{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewMockTaskRepo(t)
			tc.mockBehavior(repo, tc.payload)
			svc := service.NewTaskService(slog.Default(), repo, nil)
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
	}{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tasks := mocks.NewMockTaskRepo(t)
			statuses := mocks.NewMockStatusRepo(t)
			tc.mockBehavior(tasks, statuses, tc.args)
			svc := service.NewTaskService(slog.Default(), tasks, statuses)
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
