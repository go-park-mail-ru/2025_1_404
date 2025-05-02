package redis

import (
	"context"
	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/redis/go-redis/v9"
	"time"
)

type redisRepo struct {
	client *redis.Client
	logger logger.Logger
}

func New(cfg *config.RedisConfig, logger logger.Logger) (*redisRepo, error) {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return &redisRepo{client: client, logger: logger}, nil
}

func (repo *redisRepo) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := repo.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (repo *redisRepo) Get(ctx context.Context, key string) (string, error) {
	value, err := repo.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (repo *redisRepo) IsNotFound(err error) bool {
	return err == redis.Nil
}
