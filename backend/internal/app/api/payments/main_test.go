package payments_test

import (
	"context"
	"os"
	"testing"

	"github.com/finance-dashboard/backend/internal/app/api/payments"
	"github.com/finance-dashboard/backend/internal/app/api/users"
	"github.com/finance-dashboard/backend/internal/pkg/postgres"
	"github.com/finance-dashboard/backend/internal/pkg/tables"
)

var paymentsService *payments.Implementation
var paymentsTable tables.Payments

// usersService нужен, чтобы в тестах авторизовывать пользователя
var usersService *users.Implementation
var usersTable tables.Users

func TestMain(m *testing.M) {
	ctx := context.Background()

	postgresConnection, _ := postgres.New(ctx)

	usersTable = tables.NewUsers(postgresConnection)
	usersService = users.New(usersTable)

	paymentsTable = tables.NewPayments(postgresConnection)
	paymentsService = payments.New(paymentsTable)

	os.Exit(m.Run())
}
