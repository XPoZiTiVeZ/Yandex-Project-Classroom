package dto

type RegisterDTO struct {
	Email     string `validate:"required,email"`
	Password  string `validate:"required"`
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
}

type LoginDTO struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type TokensDTO struct {
	AccessToken  string
	RefreshToken string
}
