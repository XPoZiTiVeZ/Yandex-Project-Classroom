package auth

import (
	rds "Classroom/Gateway/internal/redis"
	pb "Classroom/Gateway/pkg/api/auth"
	"Classroom/Gateway/pkg/config"
	"Classroom/Gateway/pkg/logger"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthServiceClient struct {
	Conn           *grpc.ClientConn
	Client         pb.AuthServiceClient
	DefaultTimeout time.Duration
}

func NewAuthServiceClient(ctx context.Context, config *config.Config) (*AuthServiceClient, error) {
	address, port := config.Auth.Address, config.Auth.Port
	timeout := config.Common.Timeout

	var opts []grpc.DialOption
	opts = append(
		opts, grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", address, port), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	state := conn.GetState()
	// if !conn.WaitForStateChange(ctx, state) {
	// 	return nil, fmt.Errorf("failed to wait for state change")
	// }
	// state = conn.GetState()

	logger.Info(ctx, "Connected to gRPC auth", slog.String("address", address), slog.Int("port", port), slog.String("state", state.String()))

	client := pb.NewAuthServiceClient(conn)

	return &AuthServiceClient{
		Conn:           conn,
		Client:         client,
		DefaultTimeout: timeout,
	}, nil
}

func (s *AuthServiceClient) Register(ctx context.Context, req RegisterRequest) (RegisterResponse, error) {
	logger.Debug(ctx, "Registering user", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := s.Client.Register(ctx, NewRegisterRequest(req))
	if err != nil {
		return RegisterResponse{}, err
	}

	logger.Debug(ctx, "Auth.Register succeed")
	return NewRegisterResponse(resp), nil
}

func (s *AuthServiceClient) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	logger.Debug(ctx, "Logging in user", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, s.DefaultTimeout)
	defer cancel()

	resp, err := s.Client.Login(ctx, NewLoginRequest(req))
	if err != nil {
		return LoginResponse{}, err
	}

	logger.Debug(ctx, "Auth.Login succeed")
	return NewLoginResponse(resp), nil
}

func (s *AuthServiceClient) Refresh(ctx context.Context, req RefreshRequest) (RefreshResponse, error) {
	logger.Debug(ctx, "Refreshing user token", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := s.Client.Refresh(ctx, NewRefreshRequest(req))
	if err != nil {
		return RefreshResponse{}, err
	}

	logger.Debug(ctx, "Auth.Refresh succeed")
	return NewRefreshResponse(resp), nil
}

func (s *AuthServiceClient) Logout(ctx context.Context, req LogoutRequest) (LogoutResponse, error) {
	logger.Debug(ctx, "Logging out user", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := s.Client.Logout(ctx, NewLogoutRequest(req))
	if err != nil {
		return LogoutResponse{}, err
	}

	logger.Debug(ctx, "Auth.Logout succeed")
	return NewLogoutResponse(resp), nil
}

func (s *AuthServiceClient) GetUserInfo(ctx context.Context, rc *redis.Client, req GetUserInfoRequest) (GetUserInfoResponse, error) {
	logger.Debug(ctx, "Getting user info", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := rds.Get[GetUserInfoResponse](rc, ctx, "Auth.GetUserInfo", req.UserID)
	if err == nil {
		return resp, nil
	}
	logger.Debug(ctx, "Response was not cached", slog.Any("error", err))

	pbresp, err := s.Client.GetUserInfo(ctx, NewGetUserInfoRequest(req))
	if err != nil {
		return GetUserInfoResponse{}, err
	}

	resp = NewGetUserInfoResponse(pbresp)
	rds.Put(rc, ctx, "Auth.GetUserInfo", req.UserID, resp, 24 * time.Hour)

	logger.Debug(ctx, "Auth.GetUserInfo succeed")
	return resp, nil
}
