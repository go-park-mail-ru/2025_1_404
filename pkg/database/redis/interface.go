package redis

import (
	"context"
	"time"
)

//go:generate mockgen -source interface.go -destination=mocks/mock_redis.go -package=mocks

type RedisRepo interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	IsNotFound(err error) bool
}
