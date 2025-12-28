package models

import "time"

type PaymentsListRequest struct {
	UserID string `json:"user_id"`
}

type User struct {
	UserID        string    `db:"user_id"`
	MonthlyIncome int       `db:"monthly_income"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
