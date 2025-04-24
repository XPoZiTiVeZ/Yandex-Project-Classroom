package entities

type User struct {
	ID           string
	Email        string
	FirstName    string
	LastName     string
	IsSuperUser  bool
	PasswordHash []byte
}
