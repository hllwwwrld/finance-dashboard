package payments

import (
	"time"

	"github.com/finance-dashboard/backend/internal/pkg/models"
	"github.com/finance-dashboard/backend/internal/pkg/tables"
)

type Implementation struct {
	paymentsTable tables.Payments
}

func New(paymentsTable tables.Payments) *Implementation {
	return &Implementation{paymentsTable: paymentsTable}
}

func calculateDaysUntil(payment *models.Payment) int {
	now := time.Now()
	needMonth := int(now.Month())

	// если в этом месяца уже наступила дата платежа, то считаем, что след платеж в будущем месяце
	if payment.DueDate.Day() < time.Now().Day() {
		needMonth++
	}

	nextDueDate := time.Date(
		now.Year(), time.Month(needMonth), payment.DueDate.Day(), 0, 0, 0, 0, time.Local,
	)

	return nextDueDate.AddDate(now.Year(), int(now.Month()), now.Day()).Day()
}
