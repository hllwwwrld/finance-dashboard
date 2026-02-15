package users

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/finance-dashboard/backend/internal/pkg/models"
)

func (i *Implementation) UpdateMonthlyIncome(resp http.ResponseWriter, req *http.Request) {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(resp, fmt.Sprintf("PaymentsList.io.ReadAll err: %v", err), http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	reqBody := &models.UpdateMonthlyIncomeRequest{}
	err = json.Unmarshal(bodyBytes, reqBody)
	if err != nil {
		http.Error(resp, fmt.Sprintf("PaymentsList.json.Unmarshal err: %v", err), http.StatusInternalServerError)
		return
	}

	authCookie := GetUserFromContext(req)
	if authCookie == nil {
		http.Error(resp, "getUserFromContext err", http.StatusUnauthorized)
		return
	}

	_, err = i.usersTable.UpdateMonthlyIncome(req.Context(), authCookie.Login, reqBody.Income)
	if err != nil {
		http.Error(resp, fmt.Sprintf("usersTable.UpdateMonthlyIncome err: %v", err), http.StatusInternalServerError)
	}

	respBytes, err := json.Marshal(models.RegisterUserResponse{Success: true})
	if err != nil {
		http.Error(resp, fmt.Sprintf("UpdateMonthlyIncome.json.Marshal err: %v", err), http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	if _, err = resp.Write(respBytes); err != nil {
		slog.Error(fmt.Sprintf("PaymentsList.resp.Write err: %v", err))
	}
}
