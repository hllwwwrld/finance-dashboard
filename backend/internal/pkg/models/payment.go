package models

import "time"

type Payment struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"userId" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	Amount    int       `json:"amount" db:"amount"`
	DueDate   time.Time `json:"dueData" db:"due_date"`
	Category  string    `json:"category" db:"category"`
	Color     string    `json:"color" db:"color"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
	DaysUntil int       `json:"daysUntil" db:"-"`
}
