package users

import (
	"net/http"

	"github.com/finance-dashboard/backend/internal/pkg/middlewares"
	"github.com/finance-dashboard/backend/internal/pkg/tables"
)

type Implementation struct {
	usersTable tables.Users
}

func New(paymentsTable tables.Users) *Implementation {
	return &Implementation{usersTable: paymentsTable}
}

// GetUserFromContext - Получение пользователя из контекста
func GetUserFromContext(r *http.Request) *middlewares.Claims {
	user, ok := r.Context().Value(middlewares.UserContextKey).(*middlewares.Claims)
	if !ok {
		return nil
	}

	return user
}
