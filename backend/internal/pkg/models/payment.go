package models

import "time"

type Payment struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Name      string    `db:"name"`
	Amount    int       `db:"amount"`
	DueDate   time.Time `db:"due_date"`
	Category  string    `db:"category"`
	Color     string    `db:"color"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
