package service_test

import (
	"Classroom/Lessons/internal/domain"
	"Classroom/Lessons/internal/dto"
	"Classroom/Lessons/internal/service"
	mocks "Classroom/Lessons/internal/service/mocks"
	"context"
	"log/slog"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLessonService_Create(t *testing.T) {
	type MockBehavior func(repo *mocks.MockLessonRepo, payload dto.CreateLessonDTO)

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		payload      dto.CreateLessonDTO
		want         domain.Lesson
		wantErr      error
	}{
		{
			name: "success",
			payload: dto.CreateLessonDTO{
				Title:    "Math",
				CourseID: "course-id",
				Content:  "Math content",
			},
			mockBehavior: func(repo *mocks.MockLessonRepo, payload dto.CreateLessonDTO) {
				lesson := domain.Lesson{
					ID:       "lesson-id",
					Title:    payload.Title,
					CourseID: payload.CourseID,
					Content:  payload.Content,
				}
				repo.EXPECT().Create(mock.Anything, payload).Return(lesson, nil)
			},
			want: domain.Lesson{
				ID:       "lesson-id",
				Title:    "Math",
				CourseID: "course-id",
				Content:  "Math content",
			},
		},
		{
			name: "failed to create",
			payload: dto.CreateLessonDTO{
				Title:    "Math",
				CourseID: "course-id",
				Content:  "Math content",
			},
			mockBehavior: func(repo *mocks.MockLessonRepo, payload dto.CreateLessonDTO) {
				repo.EXPECT().Create(mock.Anything, payload).Return(domain.Lesson{}, assert.AnError)
			},
			wantErr: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewMockLessonRepo(t)
			tc.mockBehavior(repo, tc.payload)
			svc := service.NewLessonService(slog.Default(), repo)
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
