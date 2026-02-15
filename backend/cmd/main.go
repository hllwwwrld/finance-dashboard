package main

import (
	"context"
	"log"

	"github.com/finance-dashboard/backend/internal/app/api"
	"github.com/finance-dashboard/backend/internal/app/api/payments"
	"github.com/finance-dashboard/backend/internal/config"
	"github.com/finance-dashboard/backend/internal/pkg/postgres"
	"github.com/finance-dashboard/backend/internal/pkg/tables"
)

func main() {
	ctx := context.Background()

	postgresConnection, err := postgres.New(ctx)
	if err != nil {
		log.Fatalf("postgres.New err: %v", err)
	}
	paymentsTable := tables.NewPayments(postgresConnection)
	//usersTable := tables.NewUsers(postgresConnection)

	paymentsService := payments.New(paymentsTable)
	//usersService := users.New(usersTable)

	server, err := api.New(
		api.HandlersMap{
			config.PingEndpoint: api.Ping,
			// хендлеры для payments
			config.PaymentsListEndpoint: paymentsService.PaymentsList,
			// todo тут будут еще хендлеры пейментсов

			// хендлеры для юзеров
			// todo хенделеры для users
		},
	)
	if err != nil {
		log.Fatalf("api.New err: %v", err)
	}

	server.Run()
}
