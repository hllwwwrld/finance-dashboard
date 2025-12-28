package postgres

// Пример использования PostgreSQL с Squirrel
//
// Этот файл содержит примеры использования Service для работы с БД.
// В продакшене этот файл можно удалить.

import (
	"context"
	"log"
)

// ExampleUsage демонстрирует базовое использование Service
func ExampleUsage() {
	// Создаем подключение к БД
	service, err := New(context.Background())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer service.Close()

	// Пример 1: SELECT запрос
	users, err := service.GetUsers(context.Background())
	if err != nil {
		log.Printf("Error getting users: %v", err)
	}
	log.Printf("Found %d users", len(users))

	// Пример 2: INSERT запрос
	userID, err := service.CreateUser(context.Background(), "John Doe", "john@example.com")
	if err != nil {
		log.Printf("Error creating user: %v", err)
	}
	log.Printf("Created user with ID: %d", userID)
}

// GetUsers пример SELECT запроса с Squirrel
func (s *Service) GetUsers(ctx context.Context) ([]User, error) {
	// Строим SQL запрос с помощью Squirrel
	query, args, err := s.Builder.
		Select("id", "name", "email", "created_at").
		From("users").
		Where("active = ?", true).
		OrderBy("created_at DESC").
		Limit(10).
		ToSql()

	if err != nil {
		return nil, err
	}

	// Выполняем запрос
	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, rows.Err()
}

// CreateUser пример INSERT запроса с Squirrel
func (s *Service) CreateUser(ctx context.Context, name, email string) (int64, error) {
	// Строим INSERT запрос
	query, args, err := s.Builder.
		Insert("users").
		Columns("name", "email").
		Values(name, email).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return 0, err
	}

	// Выполняем запрос и получаем ID
	var id int64
	err = s.DB.QueryRowContext(ctx, query, args...).Scan(&id)
	return id, err
}

// User пример модели данных
type User struct {
	ID        int64
	Name      string
	Email     string
	CreatedAt string
}
