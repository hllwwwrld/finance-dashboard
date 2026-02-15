package models

// UpdateMonthlyIncomeRequest ...
type UpdateMonthlyIncomeRequest struct {
	Income int `json:"income"`
}

// UpdateMonthlyIncomeResponse ...
type UpdateMonthlyIncomeResponse struct {
	bool `json:"success"`
}
