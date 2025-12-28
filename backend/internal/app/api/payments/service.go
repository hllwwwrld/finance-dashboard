package payments

import "github.com/finance-dashboard/backend/internal/pkg/tables"

type Implementation struct {
	paymentsTable tables.Payments
}

func New(paymentsTable tables.Payments) *Implementation {
	return &Implementation{paymentsTable: paymentsTable}
}
