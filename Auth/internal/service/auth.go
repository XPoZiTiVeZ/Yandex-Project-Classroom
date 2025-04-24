package service

import (
	"Classroom/Auth/internal/config"
	"Classroom/Auth/internal/dto"
	"Classroom/Auth/internal/entities"
	"Classroom/Auth/pkg/e"
	"context"
	"errors"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	Create(ctx context.Context, dto dto.CreateUserDTO) (entities.User, error)
	GetByEmail(ctx context.Context, email string) (entities.User, error)
	GetByID(ctx context.Context, id string) (entities.User, error)
}

type TokenRepo interface {
	Create(ctx context.Context, userID string, ttl time.Duration) (string, error)
	GetInfoByToken(ctx context.Context, token string) (entities.RefreshToken, error)
	Revoke(ctx context.Context, token string) (bool, error)
}

type authService struct {
	logger *slog.Logger // Для дебага и информации, ошибки логируются выше
	users  UserRepo
	tokens TokenRepo
	conf   config.Auth
}

func NewAuthService(logger *slog.Logger, users UserRepo, tokens TokenRepo, conf config.Auth) *authService {
	return &authService{
		logger: logger.With(slog.String("service", "auth")),
		users:  users,
		tokens: tokens,
		conf:   conf,
	}
}

func (a *authService) Register(ctx context.Context, payload dto.RegisterDTO) (string, error) {
	// Хешируем пароль
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", e.Wrap(err, "failed to hash password")
	}

	// Метод проверяет unique constraint на email, и возвращает ErrUserAlreadyExists
	user, err := a.users.Create(ctx, dto.CreateUserDTO{
		Email:        payload.Email,
		PasswordHash: passwordHash,
		FirstName:    payload.FirstName,
		LastName:     payload.LastName,
	})
	if err != nil {
		return "", e.Wrap(err, "failed to create user")
	}

	a.logger.Info("user registered", "id", user.ID)
	return user.ID, nil
}

func (a *authService) Login(ctx context.Context, payload dto.LoginDTO) (dto.TokensDTO, error) {
	// Проверяем существует ли пользователь
	user, err := a.users.GetByEmail(ctx, payload.Email)
	if errors.Is(err, ErrUserNotFound) {
		return dto.TokensDTO{}, ErrInvalidCredentials
	}
	if err != nil {
		return dto.TokensDTO{}, e.Wrap(err, "failed to get user by email")
	}

	// Проверяем пароль
	if bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(payload.Password)) != nil {
		return dto.TokensDTO{}, ErrInvalidCredentials
	}

	// Создаем refreshToken
	refreshToken, err := a.tokens.Create(ctx, user.ID, a.conf.RefreshTTL)
	if err != nil {
		return dto.TokensDTO{}, e.Wrap(err, "failed to create refresh token")
	}

	// Создаем accessToken
	accessToken, err := SignJWT(user.ID, user.IsSuperUser, []byte(a.conf.JwtSecret), a.conf.AccessTTL)
	if err != nil {
		return dto.TokensDTO{}, e.Wrap(err, "failed to sign access token")
	}

	a.logger.Info("user logged in", "id", user.ID)
	return dto.TokensDTO{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (a *authService) Refresh(ctx context.Context, refreshToken string) (string, error) {
	// Получаем информацию по токену
	info, err := a.tokens.GetInfoByToken(ctx, refreshToken)
	if errors.Is(err, ErrInvalidToken) {
		return "", ErrInvalidToken
	}
	if err != nil {
		return "", e.Wrap(err, "failed to get refresh token info")
	}

	if info.ExpiresAt.Before(time.Now()) {
		return "", ErrInvalidToken
	}

	// Получаем актуальную информацию о пользователе
	user, err := a.users.GetByID(ctx, info.UserID)
	if err != nil {
		return "", e.Wrap(err, "failed to get user by id")
	}

	// Создаем accessToken
	accessToken, err := SignJWT(user.ID, user.IsSuperUser, []byte(a.conf.JwtSecret), a.conf.AccessTTL)
	if err != nil {
		return "", e.Wrap(err, "failed to sign access token")
	}
	return accessToken, nil
}

func (a *authService) Logout(ctx context.Context, refreshToken string) error {
	// Удаляем refreshToken
	revoked, err := a.tokens.Revoke(ctx, refreshToken)
	if revoked {
		a.logger.Info("user logged out")
	}
	return e.WrapIfErr(err, "failed to revoke refresh token")
}

func (a *authService) GetUserInfo(ctx context.Context, userID string) (entities.User, error) {
	return a.users.GetByID(ctx, userID)
}
