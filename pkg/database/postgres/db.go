package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"

	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool инициализирует пул подключений к PostgreSQL
// - основной трафик — публичные GET-запросы (объявления, ЖК);
func NewPool(cfg *config.PostgresConfig, ctx context.Context) (*pgxpool.Pool, error) {
	sslMode := "require"
	if !cfg.SSLMode {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB, sslMode)

	pgxCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга DSN: %w", err)
	}

	// Конфигурация пула соединений

	// 1. Максимальное количество активных соединений (одновременных клиентов):
	//    - большинство запросов — SELECT
	//    - сервер поддерживает многопоточность
	pgxCfg.MaxConns = 50

	// 2. Минимальное количество соединений в пуле:
	//    - держим 5 соединений, чтобы избежать холодного старта после простоя
	pgxCfg.MinConns = 5

	// 3. Максимальное время простоя соединения:
	//    - если соединение не используется > 30 секунд — закрываем
	//    - снижает потребление памяти
	pgxCfg.MaxConnIdleTime = 30 * time.Second

	// 4. Максимальное время жизни соединения:
	//    - перезапуск соединений каждые 5 минут защищает от "зависших" транзакций
	pgxCfg.MaxConnLifetime = 5 * time.Minute

	// 5. Инициализация параметров каждого подключения:
	//    - предотвращает долгие запросы, которые могут тормозить пул
	pgxCfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		// Ограничение выполнения запроса — не более 5 сек
		if _, err := conn.Exec(ctx, "SET statement_timeout = '5s'"); err != nil {
			return fmt.Errorf("не удалось установить statement_timeout: %w", err)
		}
		// Ограничение ожидания блокировок — не более 1 сек
		if _, err := conn.Exec(ctx, "SET lock_timeout = '1s'"); err != nil {
			return fmt.Errorf("не удалось установить lock_timeout: %w", err)
		}
		return nil
	}

	// Создание пула
	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// Проверка соединения
	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("бд не отвечает: %w", err)
	}

	return pool, nil
}
