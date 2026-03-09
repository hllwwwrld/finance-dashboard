package models

import "time"

type CreatePaymentRequest struct {
	Name     string `json:"name"`
	Amount   int    `json:"amount"`
	DueDay   int    `json:"dueDay"`
	Category string `json:"category"`
	Color    string `json:"color"`
}

type CreatePaymentResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	DueDay    int       `json:"dueDate"`
	Category  string    `json:"category"`
	Color     string    `json:"color"`
	DaysUntil int       `json:"daysUntil"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
