package models

type PaymentsListRequest struct {
	UserID string `json:"userId"`
}

type PaymentsListResponse struct {
	Payments []Payment `json:"payments"`
}
