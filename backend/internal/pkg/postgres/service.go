package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

const (
	dbRetryDuration = 5 * time.Second
	dbRetryTimeout  = 10 * time.Second
	dbRetryCount    = dbRetryTimeout / dbRetryDuration
)

// Service представляет сервис для работы с PostgreSQL
type Service struct {
	DB      *sql.DB
	Builder squirrel.StatementBuilderType
}

// New создает новое подключение к PostgreSQL и инициализирует Squirrel builder
func New(ctx context.Context) (*Service, error) {
	// Получаем параметры подключения из переменных окружения
	dsn := buildDSN()

	// Открываем подключение к БД
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Проверяем подключение
	if err := pingWithRetry(ctx, db); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	//db.SetMaxOpenConns(25)        // Максимум открытых соединений
	//db.SetMaxIdleConns(5)          // Максимум неактивных соединений
	//db.SetConnMaxLifetime(5 * time.Minute)  // Время жизни соединения
	//db.SetConnMaxIdleTime(10 * time.Minute) // Время простоя соединения

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &Service{
		DB:      db,
		Builder: psql,
	}, nil
}

// buildDSN строит строку подключения (Data Source Name) для PostgreSQL
func buildDSN() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "dashboard")
	sslmode := getEnv("DB_SSLMODE", "disable") // disable для локальной разработки

	// Формат DSN для lib/pq: postgres://user:password@host:port/dbname?sslmode=disable
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
}

// getEnv получает переменную окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func pingWithRetry(ctx context.Context, db *sql.DB) error {
	err := db.Ping()

	retriesLeft := dbRetryCount

	for err != nil && retriesLeft > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-time.NewTimer(dbRetryDuration).C:
			retriesLeft--
		}

		err = db.Ping()
	}

	return err
}

// Close закрывает подключение к БД
func (s *Service) Close() error {
	if s.DB != nil {
		return s.DB.Close()
	}
	return nil
}
