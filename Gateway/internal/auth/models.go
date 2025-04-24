package auth

import (
	pb "Classroom/Gateway/pkg/api/auth"
)

// RegisterRequest - запрос на регистрацию пользователя
// @Description Содержит данные, необходимые для создания нового пользователя
type RegisterRequest struct {
    // Email пользователя
    Email string `json:"email" example:"user@example.com" extensions:"x-order=0"`
    // Имя пользователя
    FirstName string `json:"first_name" example:"Иван" extensions:"x-order=1"`
    // Фамилия пользователя
    LastName string `json:"last_name" example:"Иванов" extensions:"x-order=2"`
    // Пароль
    Password string `json:"password" example:"Qwerty123!" extensions:"x-order=3"`
} // @name AuthRegisterRequest

func NewRegisterRequest(req RegisterRequest) *pb.RegisterRequest {
    return &pb.RegisterRequest{
        Email:          req.Email,
        Password:       req.Password,
        FirstName:      req.FirstName,
        LastName:       req.LastName,
    }
}

// RegisterResponse - ответ на успешную регистрацию
// @Description Возвращает ID созданного пользователя
type RegisterResponse struct {
    // Уникальный идентификатор пользователя
    UserID string `json:"user_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
} // @name AuthRegisterResponse

func NewRegisterResponse(resp *pb.RegisterResponse) RegisterResponse {
    return RegisterResponse{
        resp.GetUserId(),
    }
}

// LoginRequest - запрос на аутентификацию
// @Description Содержит учетные данные для входа в систему
type LoginRequest struct {
    // Email пользователя
    Email string `json:"email" example:"user@example.com" extensions:"x-order=0"`
    // Пароль пользователя
    Password string `json:"password" example:"Qwerty123!" extensions:"x-order=1"`
} // @name AuthLoginRequest

func NewLoginRequest(req LoginRequest) *pb.LoginRequest {
    return &pb.LoginRequest{
        Email:    req.Email,
        Password: req.Password,
    }
}

// LoginResponse - ответ с токенами доступа
// @Description Возвращает пару access/refresh токенов для аутентификации
type LoginResponse struct {
    // Токен для обновления access токена
    RefreshToken string `json:"refresh_token" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
    // Токен для доступа к API
    AccessToken string `json:"access_token" example:"eyJhbG..." extensions:"x-order=1"`
} // @name AuthLoginResponse

func NewLoginResponse(resp *pb.LoginResponse) LoginResponse {
    return LoginResponse{
        RefreshToken: resp.GetRefreshToken(),
        AccessToken: resp.GetAccessToken(),
    }
}

// RefreshRequest - запрос на обновление токена
// @Description Содержит refresh токен для получения нового access токена
type RefreshRequest struct {
    // Действующий refresh токен
    RefreshToken string `json:"refresh_token" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
} // @name AuthRefreshRequest

func NewRefreshRequest(req RefreshRequest) *pb.RefreshRequest {
    return &pb.RefreshRequest{
        RefreshToken: req.RefreshToken,
    }
}

// RefreshResponse - ответ с новым токеном доступа
// @Description Возвращает новый access токен
type RefreshResponse struct {
    // Новый access токен
    AccessToken string `json:"access_token" example:"eyJhbG..." extensions:"x-order=0"`
} // @name AuthRefreshResponse

func NewRefreshResponse(resp *pb.RefreshResponse) RefreshResponse {
    return RefreshResponse{
        AccessToken: resp.GetAccessToken(),
    }
}

// LogoutRequest - запрос на выход из системы
// @Description Содержит refresh токен для инвалидации сессии
type LogoutRequest struct {
    // Refresh токен для удаления
    RefreshToken string `json:"refresh_token" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
} // @name AuthLogoutRequest

func NewLogoutRequest(req LogoutRequest) *pb.LogoutRequest {
    return &pb.LogoutRequest{
        RefreshToken: req.RefreshToken,
    }
}

// LogoutResponse - подтверждение выхода
// @Description Пустой ответ, указывающий на успешный выход
type LogoutResponse struct{

} // @name AuthLogoutResponse

func NewLogoutResponse(resp *pb.LogoutResponse) LogoutResponse {
    return LogoutResponse{}
}

type GetUserInfoRequest struct {
    UserID string `schema:"user_id"`
}

func NewGetUserInfoRequest(req GetUserInfoRequest) *pb.GetUserInfoRequest {
    return &pb.GetUserInfoRequest{
        UserId: req.UserID,
    }
}

// UserInfoResponse - полная информация о пользователе
// @Description Возвращает все доступные данные пользователя
type GetUserInfoResponse struct {
    // Уникальный идентификатор
    UserID string `json:"user_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
    // Email адрес
    Email string `json:"email" example:"user@example.com" extensions:"x-order=1"`
    // Имя
    FirstName string `json:"first_name" example:"Иван" extensions:"x-order=2"`
    // Фамилия
    LastName string `json:"last_name" example:"Иванов" extensions:"x-order=3"`
    // Признак администратора
    IsSuperUser bool `json:"is_superuser" example:"false" extensions:"x-order=4"`
} // @name AuthUserInfoResponse

func NewGetUserInfoResponse(resp *pb.GetUserInfoResponse) GetUserInfoResponse {
    return GetUserInfoResponse{
        UserID:      resp.GetUserId(),
        Email:       resp.GetEmail(),
        FirstName:   resp.GetFirstName(),
        LastName:    resp.GetLastName(),
        IsSuperUser: resp.GetIsSuperuser(),
    }
}