package entities

import "time"

type RefreshToken struct {
	UserID    string
	ExpiresAt time.Time
}
