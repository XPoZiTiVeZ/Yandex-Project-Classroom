package auth

import (
	pb "Classroom/Gateway/pkg/api/auth"
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
)

type AuthServiceClient struct {
	Client *pb.AuthServiceClient
}

func NewAuthServiceClient(address string, port int) (*AuthServiceClient, error) {
	var opts []grpc.DialOption

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", address, port), opts...)
	if err != nil {
		slog.Error("fail to dial: %v", slog.Any("error", err))
		return nil, err
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	return &AuthServiceClient{
		Client: &client,
	}, nil
}

func (s AuthServiceClient) Register(ctx context.Context, req RegisterRequest) (RegisterResponse, error) {
	slog.Info("registering user", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := (*s.Client).Register(ctx, NewRegisterRequest(req))

	if err != nil {
		slog.Error("auth.Register failed", slog.Any("error", err))
		return RegisterResponse{}, err
	}

	slog.Info("auth.Register succeed")
	return NewRegisterResponse(resp), nil
}

func (s AuthServiceClient) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	slog.Info("registering user", slog.Any("request", req))
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := (*s.Client).Login(ctx, NewLoginRequest(req))

	if err != nil {
		slog.Error("auth.Register failed", slog.Any("error", err))
		return LoginResponse{}, err
	}

	slog.Info("auth.Register succeed")
	return NewLoginResponse(resp), nil
}

// func (s AuthServiceClient) Refresh(ctx context.Context, req RefreshUserRequest) error {
// 	slog.Info("registering user", slog.Any("request", req))
// 	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
// 	defer cancel()

// 	resp, err := (*s.Client).Register(ctx, &pb.RegisterRequest{
// 		RefreshToken: req.RefreshToken,
// 	})

// 	if err != nil {
// 		slog.Error("auth.Register failed", slog.Any("error", err))
// 		return err
// 	}

// 	if !resp.GetSuccess() {
// 		slog.Error("auth.Register failed", slog.String("error", resp.GetMessage()))
// 		return fmt.Errorf(resp.GetMessage())
// 	}

// 	slog.Info("auth.Register succeed")
// 	return nil
// }

// func (s AuthServiceClient) Logout(ctx context.Context, req LogoutUserRequest) error {
// 	slog.Info("registering user", slog.Any("request", req))
// 	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
// 	defer cancel()

// 	resp, err := (*s.Client).Register(ctx, &pb.RegisterRequest{
// 		RefreshToken: req.RefreshToken,
// 	})

// 	if err != nil {
// 		slog.Error("auth.Register failed", slog.Any("error", err))
// 		return err
// 	}

// 	if !resp.GetSuccess() {
// 		slog.Error("auth.Register failed", slog.String("error", resp.GetMessage()))
// 		return fmt.Errorf(resp.GetMessage())
// 	}

// 	slog.Info("auth.Register succeed")
// 	return nil
// }
