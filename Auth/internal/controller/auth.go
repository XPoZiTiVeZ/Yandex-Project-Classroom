package controller

import (
	"context"
	"errors"
	"log/slog"

	"Classroom/Auth/internal/dto"
	"Classroom/Auth/internal/entities"
	"Classroom/Auth/internal/service"
	pb "Classroom/Auth/pkg/api/auth"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	Register(ctx context.Context, dto dto.RegisterDTO) (string, error)
	Login(ctx context.Context, dto dto.LoginDTO) (dto.TokensDTO, error)
	Refresh(ctx context.Context, refreshToken string) (string, error)
	Logout(ctx context.Context, refreshToken string) error
	GetUserInfo(ctx context.Context, userID string) (entities.User, error)
}

type authController struct {
	svc      AuthService
	logger   *slog.Logger // Для логирования ошибок и дебага запросов
	validate *validator.Validate
	pb.UnimplementedAuthServiceServer
}

func NewAuthController(logger *slog.Logger, svc AuthService) *authController {
	validate := validator.New()
	return &authController{
		svc:      svc,
		logger:   logger,
		validate: validate,
	}
}

func (c *authController) Init(srv *grpc.Server) {
	pb.RegisterAuthServiceServer(srv, c)
}

func (c *authController) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	const op = "controller.Register"
	logger := c.logger.With(slog.String("op", op))

	// преобразование в dto для передачи между слоями
	dto := dto.RegisterDTO{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	// валидация данных
	if err := c.validate.Struct(dto); err != nil {
		logger.Debug("invalid request", "err", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	userID, err := c.svc.Register(ctx, dto)

	if errors.Is(err, service.ErrUserAlreadyExists) {
		logger.Debug("user already exists")
		return nil, status.Errorf(codes.AlreadyExists, "user already exists")
	}
	if err != nil {
		logger.Error("failed to register user", "err", err)
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &pb.RegisterResponse{UserId: userID}, nil
}

func (c *authController) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	const op = "controller.Login"
	logger := c.logger.With(slog.String("op", op))

	// преобразование в dto для передачи между слоями
	dto := dto.LoginDTO{
		Email:    req.Email,
		Password: req.Password,
	}

	// валидация данных
	if err := c.validate.Struct(dto); err != nil {
		logger.Debug("invalid request", "err", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	tokens, err := c.svc.Login(ctx, dto)

	if errors.Is(err, service.ErrInvalidCredentials) {
		logger.Debug("invalid credentials")
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}
	if err != nil {
		logger.Error("failed to login user", "err", err)
		return nil, status.Error(codes.Internal, "failed to login user")
	}

	return &pb.LoginResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func (c *authController) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	const op = "controller.Refresh"
	logger := c.logger.With(slog.String("op", op))

	accessToken, err := c.svc.Refresh(ctx, req.RefreshToken)

	if errors.Is(err, service.ErrInvalidToken) {
		logger.Debug("invalid token")
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	if err != nil {
		logger.Error("failed to refresh token", "err", err)
		return nil, status.Error(codes.Internal, "failed to refresh token")
	}

	return &pb.RefreshResponse{AccessToken: accessToken}, nil
}

func (c *authController) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	const op = "controller.Logout"
	logger := c.logger.With(slog.String("op", op))

	if err := c.svc.Logout(ctx, req.RefreshToken); err != nil {
		logger.Error("failed to logout user", "err", err)
		return nil, status.Error(codes.Internal, "failed to logout user")
	}

	return &pb.LogoutResponse{}, nil
}

func (c *authController) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	if err := c.validate.Var(req.UserId, "required,uuid"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}
	user, err := c.svc.GetUserInfo(ctx, req.UserId)
	if errors.Is(err, service.ErrUserNotFound) {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	if err != nil {
		c.logger.Error("failed to get user info", "err", err, "id", req.UserId)
		return nil, status.Error(codes.Internal, "failed to get user info")
	}

	return &pb.GetUserInfoResponse{
		UserId:      user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		IsSuperuser: user.IsSuperUser,
	}, nil
}
