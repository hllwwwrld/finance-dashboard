package payments

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/finance-dashboard/backend/internal/app/api/users"
	"github.com/finance-dashboard/backend/internal/pkg/models"
)

func (i *Implementation) CreatePayment(resp http.ResponseWriter, req *http.Request) {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(resp, fmt.Sprintf("PaymentsList.io.ReadAll err: %v", err), http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	reqBody := &models.CreatePaymentRequest{}
	err = json.Unmarshal(bodyBytes, reqBody)
	if err != nil {
		http.Error(resp, fmt.Sprintf("PaymentsList.json.Unmarshal err: %v", err), http.StatusInternalServerError)
		return
	}

	authCookie := users.GetUserFromContext(req)
	if authCookie == nil {
		http.Error(resp, "GetUserFromContext err", http.StatusUnauthorized)
		return
	}

	dueDate, err := time.Parse(time.DateOnly, reqBody.DueDate)
	if err != nil {
		http.Error(resp, fmt.Sprintf("time.Parse err: %v", err), http.StatusInternalServerError)

	}

	user := models.Payment{
		UserID:   authCookie.UserID,
		Name:     reqBody.Name,
		Amount:   reqBody.Amount,
		DueDate:  dueDate,
		Category: reqBody.Category,
		Color:    reqBody.Color,
	}
	createdPayment, err := i.paymentsTable.Create(req.Context(), user)
	if err != nil {
		http.Error(resp, fmt.Sprintf("usersTable.Create err: %v", err), http.StatusInternalServerError)
	}

	respBytes, err := json.Marshal(
		models.CreatePaymentResponse{
			ID:        createdPayment.ID,
			UserID:    createdPayment.UserID,
			Name:      createdPayment.Name,
			Amount:    createdPayment.Amount,
			DueDate:   createdPayment.DueDate.Format("02-01-2006"),
			Category:  createdPayment.Category,
			Color:     createdPayment.Color,
			DaysUntil: calculateDaysUntil(createdPayment),
			CreatedAt: createdPayment.CreatedAt,
			UpdatedAt: createdPayment.UpdatedAt,
		},
	)
	if err != nil {
		http.Error(resp, fmt.Sprintf("PaymentsList.json.Marshal err: %v", err), http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	if _, err = resp.Write(respBytes); err != nil {
		slog.Error(fmt.Sprintf("PaymentsList.resp.Write err: %v", err))
	}
}
