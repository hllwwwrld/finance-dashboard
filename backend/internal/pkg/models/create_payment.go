package models

import "time"

type CreatePaymentRequest struct {
	Name     string `json:"name"`
	Amount   int    `json:"amount"`
	DueDate  string `json:"due_date"`
	Category string `json:"category"`
	Color    string `json:"color"`
}

type CreatePaymentResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	DueDate   string    `json:"due_date"`
	Category  string    `json:"category"`
	Color     string    `json:"color"`
	DaysUntil int       `json:"days_until"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
