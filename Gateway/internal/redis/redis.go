package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func MustNew(url string) *redis.Client {
	opts, err := redis.ParseURL(url)
	if err != nil {
		log.Fatalf("failed to parse redis url: %v", err)
	}

	client := redis.NewClient(opts)

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	return client
}

func Get[T any](rc *redis.Client, ctx context.Context, method string, id string) (T, error) {
	var payload T

	data, err := rc.Get(ctx, encodeKey(method, id)).Bytes()
	if errors.Is(err, redis.Nil) {
		return payload, redis.Nil
	}
	if err != nil {
		return payload, err
	}

	if err := json.Unmarshal(data, &payload); err != nil {
		return payload, err
	}

	return payload, nil
}

func Put(rc *redis.Client, ctx context.Context, method string, id string, body any, ttl time.Duration) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	if err := rc.Set(ctx, encodeKey(method, id), data, ttl).Err(); err != nil {
		return err
	}

	return nil
}

func Delete(rc *redis.Client, ctx context.Context, method string, id string) error {
    key := encodeKey(method, id)
    _, err := rc.Del(ctx, key).Result()
    if err != nil {
        return fmt.Errorf("failed to delete key %s: %w", key, err)
    }
    return nil
}

func encodeKey(method, id string) string {
	return fmt.Sprintf("%s:%s", method, id)
}