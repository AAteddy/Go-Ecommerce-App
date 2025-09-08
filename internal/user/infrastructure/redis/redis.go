package redis

import (
	"context"
	"time"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
	"github.com/redis/go-redis/v9"
)

// RedisTokenStore implements the TokenStore interface using Redis.

type RedisTokenStore struct {
	client *redis.Client
}

func NewRedisTokenStore(addr string) *RedisTokenStore {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisTokenStore{client: client}
}

func (r *RedisTokenStore) SaveToken(ctx context.Context, token string, userID string, expiry time.Duration) error {
	return r.client.Set(ctx, "token:"+token, userID, expiry).Err()
}

func (r *RedisTokenStore) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	_, err := r.client.Get(ctx, "blacklist:"+token).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, errors.Wrap(err, "failed to check token blacklist")
	}
	return true, nil
}

func (r *RedisTokenStore) BlacklistToken(ctx context.Context, token string, expiry time.Duration) error {
	return r.client.Set(ctx, "blacklist:"+token, "blacklisted", expiry).Err()
}