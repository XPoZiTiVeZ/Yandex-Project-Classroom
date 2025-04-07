package service_test

import (
	"Classroom/Auth/internal/config"
	"Classroom/Auth/internal/dto"
	"Classroom/Auth/internal/entities"
	"Classroom/Auth/internal/service"
	"Classroom/Auth/internal/service/mocks"
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register(t *testing.T) {
	type MockBehavior func(users *mocks.UserRepo, payload dto.RegisterDTO)

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		payload      dto.RegisterDTO
		want         string
		wantErr      error
	}{
		{
			name: "success",
			payload: dto.RegisterDTO{
				Email:     "test@example.com",
				Password:  "securePassword123",
				FirstName: "John",
				LastName:  "Doe",
			},
			mockBehavior: func(users *mocks.UserRepo, payload dto.RegisterDTO) {
				users.EXPECT().
					Create(mock.Anything, mock.MatchedBy(func(dto dto.CreateUserDTO) bool {
						return dto.Email == payload.Email &&
							dto.FirstName == payload.FirstName &&
							dto.LastName == payload.LastName &&
							dto.IsSuperUser == false &&
							len(dto.PasswordHash) > 0
					})).
					Return(entities.User{ID: "user-id"}, nil)
			},
			want:    "user-id",
			wantErr: nil,
		},
		{
			name: "user already exists",
			payload: dto.RegisterDTO{
				Email:     "duplicate@example.com",
				Password:  "securePassword123",
				FirstName: "Jane",
				LastName:  "Smith",
			},
			mockBehavior: func(users *mocks.UserRepo, payload dto.RegisterDTO) {
				users.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(entities.User{}, service.ErrUserAlreadyExists)
			},
			want:    "",
			wantErr: service.ErrUserAlreadyExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokenRepo := mocks.NewTokenRepo(t)
			userRepo := mocks.NewUserRepo(t)
			tc.mockBehavior(userRepo, tc.payload)
			conf := config.Auth{JwtSecret: "secret", AccessTTL: time.Minute, RefreshTTL: time.Minute}
			svc := service.NewAuthService(slog.Default(), userRepo, tokenRepo, conf)
			got, err := svc.Register(context.Background(), tc.payload)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	type MockBehavior func(users *mocks.UserRepo, tokens *mocks.TokenRepo, payload dto.LoginDTO)

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		payload      dto.LoginDTO
		want         dto.TokensDTO
		wantErr      error
	}{
		{
			name: "success",
			payload: dto.LoginDTO{
				Email:    "user@example.com",
				Password: "correct-password",
			},
			mockBehavior: func(users *mocks.UserRepo, tokens *mocks.TokenRepo, payload dto.LoginDTO) {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
				require.NoError(t, err)
				users.EXPECT().
					GetByEmail(mock.Anything, payload.Email).
					Return(entities.User{
						ID:           "user-id",
						PasswordHash: hashedPassword,
					}, nil)

				tokens.EXPECT().
					Create(mock.Anything, "user-id", mock.Anything).
					Return("refresh-token", nil)
			},
			want:    dto.TokensDTO{RefreshToken: "refresh-token"},
			wantErr: nil,
		},
		{
			name: "user not found",
			payload: dto.LoginDTO{
				Email:    "notfound@example.com",
				Password: "password",
			},
			mockBehavior: func(users *mocks.UserRepo, tokens *mocks.TokenRepo, payload dto.LoginDTO) {
				users.EXPECT().
					GetByEmail(mock.Anything, payload.Email).
					Return(entities.User{}, service.ErrUserNotFound)
			},
			want:    dto.TokensDTO{},
			wantErr: service.ErrInvalidCredentials,
		},
		{
			name: "wrong password",
			payload: dto.LoginDTO{
				Email:    "user@example.com",
				Password: "wrong-password",
			},
			mockBehavior: func(users *mocks.UserRepo, tokens *mocks.TokenRepo, payload dto.LoginDTO) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)

				users.EXPECT().
					GetByEmail(mock.Anything, payload.Email).
					Return(entities.User{
						ID:           "user-id",
						PasswordHash: hashedPassword,
					}, nil)
			},
			want:    dto.TokensDTO{},
			wantErr: service.ErrInvalidCredentials,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokenRepo := mocks.NewTokenRepo(t)
			userRepo := mocks.NewUserRepo(t)
			tc.mockBehavior(userRepo, tokenRepo, tc.payload)
			conf := config.Auth{JwtSecret: "secret", AccessTTL: time.Minute, RefreshTTL: time.Minute}
			svc := service.NewAuthService(slog.Default(), userRepo, tokenRepo, conf)
			got, err := svc.Login(context.Background(), tc.payload)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want.RefreshToken, got.RefreshToken)
			assert.NotEmpty(t, got.AccessToken)
		})
	}
}

func TestAuthService_Refresh(t *testing.T) {
	type MockBehavior func(users *mocks.UserRepo, tokens *mocks.TokenRepo, refreshToken string)

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		refreshToken string
		wantErr      error
	}{
		{
			name:         "success",
			refreshToken: "valid-refresh-token",
			mockBehavior: func(users *mocks.UserRepo, tokens *mocks.TokenRepo, refreshToken string) {
				tokens.EXPECT().
					GetInfoByToken(mock.Anything, refreshToken).
					Return(entities.RefreshToken{
						UserID:    "user-id",
						ExpiresAt: time.Now().Add(time.Minute),
					}, nil)

				users.EXPECT().
					GetByID(mock.Anything, "user-id").
					Return(entities.User{
						ID:          "user-id",
						IsSuperUser: false,
					}, nil)
			},
			wantErr: nil,
		},
		{
			name:         "token expired",
			refreshToken: "expired-refresh-token",
			mockBehavior: func(users *mocks.UserRepo, tokens *mocks.TokenRepo, refreshToken string) {
				tokens.EXPECT().
					GetInfoByToken(mock.Anything, refreshToken).
					Return(entities.RefreshToken{
						UserID:    "user-id",
						ExpiresAt: time.Now().Add(-time.Minute),
					}, nil)
			},
			wantErr: service.ErrInvalidToken,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokenRepo := mocks.NewTokenRepo(t)
			userRepo := mocks.NewUserRepo(t)
			tc.mockBehavior(userRepo, tokenRepo, tc.refreshToken)
			conf := config.Auth{JwtSecret: "secret", AccessTTL: time.Minute, RefreshTTL: time.Minute}
			svc := service.NewAuthService(slog.Default(), userRepo, tokenRepo, conf)
			got, err := svc.Refresh(context.Background(), tc.refreshToken)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, got)
		})
	}
}

func TestAuthService_Logout(t *testing.T) {
	type MockBehavior func(tokens *mocks.TokenRepo, refreshToken string)

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		refreshToken string
		wantErr      error
	}{
		{
			name:         "success",
			refreshToken: "valid-refresh-token",
			mockBehavior: func(tokens *mocks.TokenRepo, refreshToken string) {
				tokens.EXPECT().
					Revoke(mock.Anything, refreshToken).
					Return(true, nil)
			},
			wantErr: nil,
		},
		{
			name:         "not deleted",
			refreshToken: "valid-refresh-token",
			mockBehavior: func(tokens *mocks.TokenRepo, refreshToken string) {
				tokens.EXPECT().
					Revoke(mock.Anything, refreshToken).
					Return(false, nil)
			},
			wantErr: nil,
		},
		{
			name:         "unknown err",
			refreshToken: "valid-refresh-token",
			mockBehavior: func(tokens *mocks.TokenRepo, refreshToken string) {
				tokens.EXPECT().
					Revoke(mock.Anything, refreshToken).
					Return(false, assert.AnError)
			},
			wantErr: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokenRepo := mocks.NewTokenRepo(t)
			userRepo := mocks.NewUserRepo(t)
			tc.mockBehavior(tokenRepo, tc.refreshToken)
			conf := config.Auth{JwtSecret: "secret", AccessTTL: time.Minute, RefreshTTL: time.Minute}
			svc := service.NewAuthService(slog.Default(), userRepo, tokenRepo, conf)
			err := svc.Logout(context.Background(), tc.refreshToken)

			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}
