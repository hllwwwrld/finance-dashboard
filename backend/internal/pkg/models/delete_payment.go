package models

type DeletePaymentRequest struct {
	ID string `json:"id"`
}

type DeletePaymentResponse struct {
	Success bool `json:"success"`
}
