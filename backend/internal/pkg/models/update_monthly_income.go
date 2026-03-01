package models

// UpdateMonthlyIncomeRequest ...
type UpdateMonthlyIncomeRequest struct {
	Income int `json:"income"`
}

// UpdateMonthlyIncomeResponse ...
type UpdateMonthlyIncomeResponse struct {
	MonthlyIncome int `json:"monthlyIncome"`
}
