package database

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(cfg *config.PostgresConfig, ctx context.Context) (*pgxpool.Pool, error) {
	sslMode := "require"
	if !cfg.SSLMode {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB, sslMode)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("бд не отвечает: %w", err)
	}

	return pool, nil
}
