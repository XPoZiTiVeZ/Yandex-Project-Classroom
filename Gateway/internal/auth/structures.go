package auth

import (
	pb "Classroom/Gateway/pkg/api/auth"
)

type RegisterRequest struct {
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Password       string `json:"password"`
	PasswordRepeat string `json:"password_repeat"`
}

func NewRegisterRequest(req RegisterRequest) *pb.RegisterRequest {
	return &pb.RegisterRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}
}

type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `Json:"message"`
	UserID  string `json:"user_id"`
}

func NewRegisterResponse(resp *pb.RegisterResponse) RegisterResponse {
	return RegisterResponse{
		resp.GetSuccess(),
		resp.GetMessage(),
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
	Success      bool   `json:"success"`
	Message      string `Json:"message"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

func NewLoginResponse(resp *pb.LoginResponse) LoginResponse {
	return LoginResponse{
		resp.GetSuccess(),
		resp.GetMessage(),
		resp.GetAccessToken(),
		resp.GetRefreshToken(),
	}
}
