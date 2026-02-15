package main

import (
	"context"
	"log"

	"github.com/finance-dashboard/backend/internal/app/api"
	"github.com/finance-dashboard/backend/internal/app/api/payments"
	"github.com/finance-dashboard/backend/internal/app/api/users"
	"github.com/finance-dashboard/backend/internal/config"
	"github.com/finance-dashboard/backend/internal/pkg/middlewares"
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
	usersTable := tables.NewUsers(postgresConnection)

	paymentsService := payments.New(paymentsTable)
	usersService := users.New(usersTable)

	server, err := api.New(
		api.HandlersMap{
			config.PingEndpoint: api.Ping,

			// хендлеры для юзеров
			config.UserRegisterEndpoint:      usersService.Register,
			config.UserLoginEndpoint:         usersService.Login,
			config.UserProfileEndpoint:       middlewares.Auth(usersService.FetchProfile),
			config.UserProfileUpdateEndpoint: middlewares.Auth(usersService.UpdateMonthlyIncome),
			// todo config.UserLogoutEndpoint: usersService.Logout

			// хендлеры для payments
			config.PaymentsCreate:       middlewares.Auth(paymentsService.CreatePayment),
			config.PaymentsListEndpoint: middlewares.Auth(paymentsService.PaymentsList),
		},
	)
	if err != nil {
		log.Fatalf("api.New err: %v", err)
	}

	server.Run()
}
