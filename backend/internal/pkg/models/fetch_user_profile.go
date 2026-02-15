package models

type FetchProfileRequest struct{}

type FetchProfileResponse struct {
	MonthlyIncome int `json:"monthly_income"`
}
