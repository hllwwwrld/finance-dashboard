package models

type PaymentsListRequest struct{}

type PaymentsListResponse struct {
	Payments      []*Payment `json:"payments"`
	TotalExpenses int        `json:"totalExpenses"`
}
