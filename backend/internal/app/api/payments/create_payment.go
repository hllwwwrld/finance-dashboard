package payments

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/finance-dashboard/backend/internal/app/api/users"
	"github.com/finance-dashboard/backend/internal/pkg/models"
)

func (i *Implementation) CreatePayment(resp http.ResponseWriter, req *http.Request) {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Println(">>>>>1")
		http.Error(resp, fmt.Sprintf("PaymentsList.io.ReadAll err: %v", err), http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	reqBody := &models.CreatePaymentRequest{}
	err = json.Unmarshal(bodyBytes, reqBody)
	if err != nil {
		fmt.Println(">>>>>1")
		http.Error(resp, fmt.Sprintf("PaymentsList.json.Unmarshal err: %v", err), http.StatusInternalServerError)
		return
	}

	authCookie := users.GetUserFromContext(req)
	if authCookie == nil {
		fmt.Println(">>>>>1")
		http.Error(resp, "GetUserFromContext err", http.StatusUnauthorized)
		return
	}

	payment := models.Payment{
		UserID:   authCookie.UserID,
		Name:     reqBody.Name,
		Amount:   reqBody.Amount,
		DueDay:   reqBody.DueDay,
		Category: reqBody.Category,
		Color:    reqBody.Color,
	}
	createdPayment, err := i.paymentsTable.Create(req.Context(), payment)
	if err != nil {
		fmt.Printf(">>>>> %v", err)
		http.Error(resp, fmt.Sprintf("usersTable.Create err: %v", err), http.StatusInternalServerError)
		return
	}

	respBytes, err := json.Marshal(
		models.CreatePaymentResponse{
			ID:        createdPayment.ID,
			UserID:    createdPayment.UserID,
			Name:      createdPayment.Name,
			Amount:    createdPayment.Amount,
			DueDay:    createdPayment.DueDay,
			Category:  createdPayment.Category,
			Color:     createdPayment.Color,
			DaysUntil: calculateDaysUntil(createdPayment),
			CreatedAt: createdPayment.CreatedAt,
			UpdatedAt: createdPayment.UpdatedAt,
		},
	)
	if err != nil {
		fmt.Println(">>>>>1")
		http.Error(resp, fmt.Sprintf("PaymentsList.json.Marshal err: %v", err), http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	if _, err = resp.Write(respBytes); err != nil {
		fmt.Println(">>>>>2")
		slog.Error(fmt.Sprintf("PaymentsList.resp.Write err: %v", err))
	}
}
