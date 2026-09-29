package postgres

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool создает и проверяет подключение к PostgreSQL через пул соединений
func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
	// Получаем строку подключения из env
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		// Резервный формат из отдельных переменных, если нет единой URL
		connStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)
	}

	// Парсим конфигурацию пула
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга конфига БД: %w", err)
	}

	// Создаем сам пул
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания пула БД: %w", err)
	}
	//Создаем конекст с таймером и выклбчателем
	pingCtx, pingCancel := context.WithTimeout(ctx, 5*time.Second)
	defer pingCancel()

	// Обязательно делаем Ping, чтобы убедиться, что база реально доступна
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("не удалось подключиться к БД (ping failed): %w", err)
	}

	return pool, nil
}
