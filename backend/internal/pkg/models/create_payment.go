package models

type CreatePaymentRequest struct {
	UserID   string `json:"userId"`
	Name     string `json:"name"`
	Amount   int    `json:"amount"`
	DueDate  string `json:"dueDate"`
	Category string `json:"category"`
	Color    string `json:"color"`
}
