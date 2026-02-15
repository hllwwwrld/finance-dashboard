package users

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/finance-dashboard/backend/internal/pkg/models"
)

func (i *Implementation) FetchProfile(resp http.ResponseWriter, req *http.Request) {
	authCookie := GetUserFromContext(req)
	if authCookie == nil {
		http.Error(resp, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	userProfile, err := i.usersTable.GetByLogin(req.Context(), authCookie.Login)
	if err != nil {
		http.Error(resp, fmt.Sprintf("usersTable.UpdateMonthlyIncome err: %v", err), http.StatusInternalServerError)
	}

	respBytes, err := json.Marshal(models.FetchProfileResponse{MonthlyIncome: userProfile.MonthlyIncome})
	if err != nil {
		http.Error(resp, fmt.Sprintf("UpdateMonthlyIncome.json.Marshal err: %v", err), http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	if _, err = resp.Write(respBytes); err != nil {
		slog.Error(fmt.Sprintf("PaymentsList.resp.Write err: %v", err))
	}
}
