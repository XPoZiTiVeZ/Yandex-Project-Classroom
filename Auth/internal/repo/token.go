package repo

import (
	"Classroom/Auth/internal/entities"
	"Classroom/Auth/internal/service"
	"Classroom/Auth/pkg/e"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type tokenRepo struct {
	storage *redis.Client
}

func NewTokenRepo(storage *redis.Client) *tokenRepo {
	return &tokenRepo{
		storage: storage,
	}
}

func (r *tokenRepo) Create(ctx context.Context, userID string, ttl time.Duration) (string, error) {
	payload := RefreshToken{
		UserID:    userID,
		ExpiresAt: time.Now().Add(ttl).Unix(),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", e.Wrap(err, "failed to marshal refresh token")
	}
	token := uuid.NewString()
	if err := r.storage.Set(ctx, tokenKey(token), data, ttl).Err(); err != nil {
		return "", err
	}
	return token, nil
}

func (r *tokenRepo) GetInfoByToken(ctx context.Context, token string) (entities.RefreshToken, error) {
	data, err := r.storage.Get(ctx, tokenKey(token)).Bytes()
	if errors.Is(err, redis.Nil) {
		return entities.RefreshToken{}, service.ErrInvalidToken
	}
	if err != nil {
		return entities.RefreshToken{}, err
	}

	var payload RefreshToken
	if err := json.Unmarshal(data, &payload); err != nil {
		return entities.RefreshToken{}, e.Wrap(err, "failed to unmarshal refresh token")
	}

	return payload.ToEntity(), nil
}

func (r *tokenRepo) Revoke(ctx context.Context, token string) (bool, error) {
	res, err := r.storage.Del(ctx, tokenKey(token)).Result()
	if err != nil {
		return false, err
	}
	if res == 0 {
		return false, nil
	}
	return true, nil
}

func tokenKey(token string) string {
	return fmt.Sprintf("refreshToken:%s", token)
}
