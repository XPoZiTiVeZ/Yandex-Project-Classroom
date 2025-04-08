package controller_test

import (
	"Classroom/Auth/internal/controller"
	"Classroom/Auth/internal/controller/mocks"
	"Classroom/Auth/internal/dto"
	"Classroom/Auth/internal/service"
	pb "Classroom/Auth/pkg/api/auth"
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuthController_Register(t *testing.T) {
	type MockBehavior func(svc *mocks.AuthService, req *pb.RegisterRequest)

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		req          *pb.RegisterRequest
		want         *pb.RegisterResponse
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(svc *mocks.AuthService, req *pb.RegisterRequest) {
				svc.EXPECT().Register(mock.Anything, dto.RegisterDTO{
					Email:     req.Email,
					Password:  req.Password,
					FirstName: req.FirstName,
					LastName:  req.LastName,
				}).Return("123", nil)
			},
			req: &pb.RegisterRequest{
				Email:          "email@email.com",
				Password:       "password",
				PasswordRepeat: "password",
				FirstName:      "name",
				LastName:       "last",
			},
			want: &pb.RegisterResponse{UserId: "123"},
		},
		{
			name:         "passwords do not match",
			mockBehavior: func(svc *mocks.AuthService, req *pb.RegisterRequest) {},
			req: &pb.RegisterRequest{
				Email:          "email@email.com",
				Password:       "password",
				PasswordRepeat: "passwordwew",
				FirstName:      "name",
				LastName:       "last",
			},
			wantErr: status.Error(codes.InvalidArgument, "passwords do not match"),
		},
		{
			name:         "invalid email",
			mockBehavior: func(svc *mocks.AuthService, req *pb.RegisterRequest) {},
			req: &pb.RegisterRequest{
				Email:          "email@e",
				Password:       "password",
				PasswordRepeat: "password",
				FirstName:      "name",
				LastName:       "last",
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid request: Key: 'RegisterDTO.Email' Error:Field validation for 'Email' failed on the 'email' tag"),
		},
		{
			name: "user already exists",
			mockBehavior: func(svc *mocks.AuthService, req *pb.RegisterRequest) {
				svc.EXPECT().Register(mock.Anything, dto.RegisterDTO{
					Email:     req.Email,
					Password:  req.Password,
					FirstName: req.FirstName,
					LastName:  req.LastName,
				}).Return("", service.ErrUserAlreadyExists)
			},
			req: &pb.RegisterRequest{
				Email:          "email@email.com",
				Password:       "password",
				PasswordRepeat: "password",
				FirstName:      "name",
				LastName:       "last",
			},
			wantErr: status.Error(codes.AlreadyExists, "user already exists"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocks.NewAuthService(t)
			tc.mockBehavior(svc, tc.req)
			c := controller.NewAuthController(slog.Default(), svc)
			got, err := c.Register(context.Background(), tc.req)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAuthController_Login(t *testing.T) {
	type MockBehavior func(svc *mocks.AuthService, req *pb.LoginRequest)

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		req          *pb.LoginRequest
		want         *pb.LoginResponse
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(svc *mocks.AuthService, req *pb.LoginRequest) {
				svc.EXPECT().Login(mock.Anything, dto.LoginDTO{
					Email:    req.Email,
					Password: req.Password,
				}).Return(dto.TokensDTO{AccessToken: "token", RefreshToken: "token2"}, nil)
			},
			req: &pb.LoginRequest{
				Email:    "email@email.com",
				Password: "password",
			},
			want: &pb.LoginResponse{
				AccessToken:  "token",
				RefreshToken: "token2",
			},
		},
		{
			name:         "invalid email",
			mockBehavior: func(svc *mocks.AuthService, req *pb.LoginRequest) {},
			req: &pb.LoginRequest{
				Email:    "email@e",
				Password: "password",
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid request: Key: 'LoginDTO.Email' Error:Field validation for 'Email' failed on the 'email' tag"),
		},
		{
			name: "invalid credentials",
			mockBehavior: func(svc *mocks.AuthService, req *pb.LoginRequest) {
				svc.EXPECT().Login(mock.Anything, dto.LoginDTO{
					Email:    req.Email,
					Password: req.Password,
				}).Return(dto.TokensDTO{}, service.ErrInvalidCredentials)
			},
			req: &pb.LoginRequest{
				Email:    "email@email.com",
				Password: "password",
			},
			wantErr: status.Error(codes.Unauthenticated, "invalid credentials"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocks.NewAuthService(t)
			tc.mockBehavior(svc, tc.req)
			c := controller.NewAuthController(slog.Default(), svc)
			got, err := c.Login(context.Background(), tc.req)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAuthController_Refresh(t *testing.T) {
	type MockBehavior func(svc *mocks.AuthService, req *pb.RefreshRequest)

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		req          *pb.RefreshRequest
		want         *pb.RefreshResponse
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(svc *mocks.AuthService, req *pb.RefreshRequest) {
				svc.EXPECT().Refresh(mock.Anything, req.RefreshToken).Return("token", nil)
			},
			req: &pb.RefreshRequest{
				RefreshToken: "token",
			},
			want: &pb.RefreshResponse{
				AccessToken: "token",
			},
		},
		{
			name: "invalid token",
			mockBehavior: func(svc *mocks.AuthService, req *pb.RefreshRequest) {
				svc.EXPECT().Refresh(mock.Anything, req.RefreshToken).Return("", service.ErrInvalidToken)
			},
			req: &pb.RefreshRequest{
				RefreshToken: "token",
			},
			wantErr: status.Error(codes.Unauthenticated, "invalid token"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocks.NewAuthService(t)
			tc.mockBehavior(svc, tc.req)
			c := controller.NewAuthController(slog.Default(), svc)
			got, err := c.Refresh(context.Background(), tc.req)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAuthController_Logout(t *testing.T) {
	type MockBehavior func(svc *mocks.AuthService, req *pb.LogoutRequest)

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		req          *pb.LogoutRequest
		want         *pb.LogoutResponse
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(svc *mocks.AuthService, req *pb.LogoutRequest) {
				svc.EXPECT().Logout(mock.Anything, req.RefreshToken).Return(nil)
			},
			req: &pb.LogoutRequest{
				RefreshToken: "token",
			},
			want: &pb.LogoutResponse{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocks.NewAuthService(t)
			tc.mockBehavior(svc, tc.req)
			c := controller.NewAuthController(slog.Default(), svc)
			got, err := c.Logout(context.Background(), tc.req)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
