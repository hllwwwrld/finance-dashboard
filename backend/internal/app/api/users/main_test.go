package users_test

import (
	"context"
	"os"
	"testing"

	"github.com/finance-dashboard/backend/internal/app/api/users"
	"github.com/finance-dashboard/backend/internal/pkg/postgres"
	"github.com/finance-dashboard/backend/internal/pkg/tables"
)

var usersService *users.Implementation
var usersTable tables.Users

func TestMain(m *testing.M) {
	ctx := context.Background()

	postgresConnection, _ := postgres.New(ctx)
	usersTable = tables.NewUsers(postgresConnection)

	usersService = users.New(usersTable)

	os.Exit(m.Run())
}
