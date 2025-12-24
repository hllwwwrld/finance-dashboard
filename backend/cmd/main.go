package main

import (
	"log"

	"github.com/finance-dashboard/backend/internal/app/api"
	"github.com/finance-dashboard/backend/internal/config"
)

func main() {
	server, err := api.New(
		api.HandlersMap{
			config.PingEndpoint:         api.Ping,
			config.PaymentsListEndpoint: api.PaymentsList,
		},
	)
	if err != nil {
		log.Fatalf("small_boss.New err: %v", err)
	}

	server.Run()
}
