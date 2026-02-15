package models

// UpdateMonthlyIncomeRequest ...
type UpdateMonthlyIncomeRequest struct {
	Income int `json:"income"`
}

// UpdateMonthlyIncomeResponse ...
type UpdateMonthlyIncomeResponse struct {
	Success bool `json:"success"`
}
