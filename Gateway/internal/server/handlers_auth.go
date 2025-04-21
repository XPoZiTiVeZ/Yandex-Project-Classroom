package server

import (
	"Classroom/Gateway/internal/auth"
	"Classroom/Gateway/pkg/logger"
	"log/slog"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RegisterHandler регистрирует нового пользователя
// @Summary Регистрация пользователя
// @Description Создает новую учетную запись пользователя
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body auth.RegisterRequest true "Данные регистрации"
// @Success 201 {object} auth.RegisterResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 409 {object} ErrorResponse "Пользователь уже существует"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /auth/register [post]
func (s *Server) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[auth.RegisterRequest](r.Context())

	resp, err := s.Auth.Register(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler auth.Register error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				BadRequest(w, e.Message())
			case codes.AlreadyExists:
				AlreadyExists(w, "email already exists")
			case codes.Unavailable:
				ServiceUnavailable(w)
			}
		} else {
			InternalError(w)
		}
		return
	}

	WriteJSON(w, resp, http.StatusCreated)
}

// LoginHandler аутентифицирует пользователя
// @Summary Вход в систему
// @Description Возвращает токены доступа и обновления
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "Учетные данные"
// @Success 200 {object} auth.LoginResponse
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 401 {object} ErrorResponse "Неверные учетные данные"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /auth/login [post]
func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[auth.LoginRequest](r.Context())

	resp, err := s.Auth.Login(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler auth.Login error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				BadRequest(w, e.Message())
			case codes.Unauthenticated:
				Unauthorized(w, "invalid credentials")
			case codes.Unavailable:
				ServiceUnavailable(w)
			}
		} else {
			InternalError(w)
		}
		return
	}

	WriteJSON(w, resp, http.StatusOK)
}

// RefreshHandler обновляет токен доступа
// @Summary Обновление токена
// @Description Генерирует новую пару токенов по refresh-токену
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body auth.RefreshRequest true "Refresh токен"
// @Success 200 {object} auth.RefreshResponse
// @Failure 401 {object} ErrorResponse "Неверный токен"
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /auth/refresh [post]
func (s *Server) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[auth.RefreshRequest](r.Context())

	resp, err := s.Auth.Refresh(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler auth.Refresh error", slog.Any("error", err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Unauthenticated:
				Unauthorized(w, "invalid refresh token")
			case codes.Unavailable:
				ServiceUnavailable(w)
			}
		} else {
			InternalError(w)
		}
		return
	}

	WriteJSON(w, resp, http.StatusOK)
}

// LogoutHandler выходит из системы
// @Summary Выход из системы
// @Description Инвалидирует токены пользователя
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body auth.LogoutRequest true "Данные для выхода"
// @Success 200 {object} auth.LogoutResponse
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /auth/logout [post]
func (s *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[auth.LogoutRequest](r.Context())

	resp, err := s.Auth.Logout(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler auth.Logout error", slog.Any("error", err))

		ServiceUnavailable(w)
	}
	WriteJSON(w, resp, http.StatusOK)
}

// GetUserInfoHandler возвращает информацию о пользователе
// @Summary Информация о пользователе
// @Description Возвращает данные пользователя
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body auth.GetUserInfoRequest true "Идентификатор пользователя"
// @Success 200 {object} auth.GetUserInfoResponse
// @Failure 503 {object} ErrorResponse "Сервис недоступен"
// @Router /auth/user-info [post]
func (s *Server) GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	body := GetBody[auth.GetUserInfoRequest](r.Context())

	resp, err := s.Auth.GetUserInfo(r.Context(), body)
	if err != nil {
		logger.Error(r.Context(), "Handler auth.GetUserInfo error", slog.Any("error", err))

		ServiceUnavailable(w)
	}

	WriteJSON(w, resp, http.StatusOK)
}
