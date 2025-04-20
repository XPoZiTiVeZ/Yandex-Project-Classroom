package auth

import (
	pb "Classroom/Gateway/pkg/api/auth"
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthServiceClient struct {
	logger *slog.Logger
	Conn   *grpc.ClientConn
	// Это интерфейс, тут не нужен указатель
	Client         pb.AuthServiceClient
	DefaultTimeout time.Duration
}

func NewAuthServiceClient(logger *slog.Logger, address string, port int, DefaultTimeout *time.Duration) (*AuthServiceClient, error) {
	var opts []grpc.DialOption
	opts = append(
		opts, grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", address, port), opts...)
	if err != nil {
		logger.Error("fail to dial: %v", slog.Any("error", err))
		return nil, err
	}

	state := conn.GetState()
	// if state != connectivity.Ready {
	// 	return nil, fmt.Errorf("connection is not ready, state: %v", state)
	// }

	logger.Info("Connected to grpc Auth", slog.String("address", address), slog.Int("port", port), slog.String("state", state.String()))

	client := pb.NewAuthServiceClient(conn)

	timeout := 10 * time.Second
	if DefaultTimeout != nil {
		timeout = *DefaultTimeout
	}

	return &AuthServiceClient{
		logger:         logger,
		Conn:           conn,
		Client:         client,
		DefaultTimeout: timeout,
	}, nil
}

func (s *AuthServiceClient) Register(ctx context.Context, req RegisterRequest) (RegisterResponse, error) {
	s.logger.Debug("registering user", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := s.Client.Register(ctx, NewRegisterRequest(req))
	if err != nil {
		return RegisterResponse{}, err
	}

	s.logger.Debug("auth.Register succeed")
	return NewRegisterResponse(resp), nil
}

func (s *AuthServiceClient) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	s.logger.Debug("logging in user", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := s.Client.Login(ctx, NewLoginRequest(req))
	if err != nil {
		return LoginResponse{}, err
	}

	s.logger.Debug("auth.Login succeed")
	return NewLoginResponse(resp), nil
}

func (s *AuthServiceClient) Refresh(ctx context.Context, req RefreshRequest) (RefreshResponse, error) {
	s.logger.Debug("refreshing user token", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := s.Client.Refresh(ctx, NewRefreshRequest(req))
	if err != nil {
		return RefreshResponse{}, err
	}

	s.logger.Debug("auth.Refresh succeed")
	return NewRefreshResponse(resp), nil
}

func (s *AuthServiceClient) Logout(ctx context.Context, req LogoutRequest) (LogoutResponse, error) {
	s.logger.Debug("logging out user", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := s.Client.Logout(ctx, NewLogoutRequest(req))
	if err != nil {
		return LogoutResponse{}, err
	}

	s.logger.Debug("auth.Logout succeed")
	return NewLogoutResponse(resp), nil
}

func (s *AuthServiceClient) GetUserInfo(ctx context.Context, req GetUserInfoRequest) (GetUserInfoResponse, error) {
	s.logger.Debug("getting user info", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := s.Client.GetUserInfo(ctx, NewGetUserInfoRequest(req))
	if err != nil {
		return GetUserInfoResponse{}, err
	}

	s.logger.Debug("auth.GetUserInfo succeed")
	return NewGetUserInfoResponse(resp), nil
}
