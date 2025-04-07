package repo

import (
	"Classroom/Auth/internal/entities"
	"time"
)

type User struct {
	ID           string `db:"user_id"`
	Email        string `db:"email"`
	PasswordHash []byte `db:"password_hash"`
	IsSuperUser  bool   `db:"is_superuser"`
	FirstName    string `db:"first_name"`
	LastName     string `db:"last_name"`
}

func (u User) ToEntity() entities.User {
	return entities.User{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		IsSuperUser:  u.IsSuperUser,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
	}
}

type RefreshToken struct {
	UserID    string `json:"user_id"`
	ExpiresAt int64  `json:"expires_at"`
}

func (r RefreshToken) ToEntity() entities.RefreshToken {
	return entities.RefreshToken{
		UserID:    r.UserID,
		ExpiresAt: time.Unix(r.ExpiresAt, 0),
	}
}
