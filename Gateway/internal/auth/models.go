package auth

import (
	pb "Classroom/Gateway/pkg/api/auth"
)

type RegisterRequest struct {
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Password       string `json:"password"`
}

func NewRegisterRequest(req RegisterRequest) *pb.RegisterRequest {
	return &pb.RegisterRequest{
		Email:          req.Email,
		Password:       req.Password,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
	}
}

type RegisterResponse struct {
	UserID string `json:"user_id"`
}

func NewRegisterResponse(resp *pb.RegisterResponse) RegisterResponse {
	return RegisterResponse{
		resp.GetUserId(),
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewLoginRequest(req LoginRequest) *pb.LoginRequest {
	return &pb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}
}

type LoginResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

func NewLoginResponse(resp *pb.LoginResponse) LoginResponse {
	return LoginResponse{
		resp.GetAccessToken(),
		resp.GetRefreshToken(),
	}
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func NewRefreshRequest(req RefreshRequest) *pb.RefreshRequest {
	return &pb.RefreshRequest{
		RefreshToken: req.RefreshToken,
	}
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

func NewRefreshResponse(resp *pb.RefreshResponse) RefreshResponse {
	return RefreshResponse{
		resp.GetAccessToken(),
	}
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func NewLogoutRequest(req LogoutRequest) *pb.LogoutRequest {
	return &pb.LogoutRequest{
		RefreshToken: req.RefreshToken,
	}
}

type LogoutResponse struct {
}

func NewLogoutResponse(resp *pb.LogoutResponse) LogoutResponse {
	return LogoutResponse{}
}

type GetUserInfoRequest struct {
	UserID string `json:"user_id"`
}

func NewGetUserInfoRequest(req GetUserInfoRequest) *pb.GetUserInfoRequest {
	return &pb.GetUserInfoRequest{
		UserId: req.UserID,
	}
}

type GetUserInfoResponse struct {
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	IsSuperUser bool
}

func NewGetUserInfoResponse(resp *pb.GetUserInfoResponse) GetUserInfoResponse {
	return GetUserInfoResponse{
		UserID:      resp.GetUserId(),
		Email:       resp.GetEmail(),
		FirstName:   resp.GetFirstName(),
		LastName:    resp.GetLastName(),
		IsSuperUser: resp.GetIsSuperuser(),
	}
}