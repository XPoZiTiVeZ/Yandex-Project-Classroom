package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	UserID      string `json:"user_id"`
	IsSuperUser bool   `json:"is_superuser"`
	jwt.RegisteredClaims
}

func VerifyJWT(tokenString string, secret []byte) (*JwtClaims, error) {
	claims := &JwtClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}

func SignJWT(userID string, isSuperUser bool, secretKey []byte, ttl time.Duration) (string, error) {
	claims := JwtClaims{
		UserID:      userID,
		IsSuperUser: isSuperUser,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return token, nil
}
