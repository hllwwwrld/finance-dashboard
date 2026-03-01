package models

import "time"

type CreatePaymentRequest struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
	// todo dueDate должен быть int, но в api он string, заменить тут на int
	DueDate  string `json:"dueDate"`
	Category string `json:"category"`
	Color    string `json:"color"`
}

type CreatePaymentResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	DueDate   string    `json:"dueDate"`
	Category  string    `json:"category"`
	Color     string    `json:"color"`
	DaysUntil int       `json:"daysUntil"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
