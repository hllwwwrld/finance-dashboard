package main

import (
	"context"
	"log"

	"github.com/finance-dashboard/backend/internal/app/api"
	"github.com/finance-dashboard/backend/internal/config"
	"github.com/finance-dashboard/backend/internal/pkg/postgres"
	"github.com/finance-dashboard/backend/internal/pkg/tables"
)

func main() {
	ctx := context.Background()

	server, err := api.New(
		api.HandlersMap{
			config.PingEndpoint:         api.Ping,
			config.PaymentsListEndpoint: api.PaymentsList,
		},
	)
	if err != nil {
		log.Fatalf("api.New err: %v", err)
	}

	postgresConnection, err := postgres.New(ctx)
	if err != nil {
		log.Fatalf("postgres.New err: %v", err)
	}
	tables.NewPayments(postgresConnection)

	server.Run()
}
