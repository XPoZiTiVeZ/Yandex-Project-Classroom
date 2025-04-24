package dto

type CreateUserDTO struct {
	Email        string
	PasswordHash []byte
	FirstName    string
	LastName     string
	IsSuperUser  bool
}
