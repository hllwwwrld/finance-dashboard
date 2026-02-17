package payments

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/finance-dashboard/backend/internal/app/api/users"
)

func (i *Implementation) PaymentsList(resp http.ResponseWriter, req *http.Request) {
	authCookie := users.GetUserFromContext(req)
	if authCookie == nil {
		http.Error(resp, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	userPayments, err := i.paymentsTable.GetByUserID(req.Context(), authCookie.UserID)
	if err != nil {
		http.Error(resp, fmt.Sprintf("PaymentsList.i.paymentsTable.GetByUserID err: %v", err), http.StatusInternalServerError)
		return
	}

	for _, payment := range userPayments {
		payment.DaysUntil = calculateDaysUntil(payment)
	}

	respBytes, err := json.Marshal(userPayments)
	if err != nil {
		http.Error(resp, fmt.Sprintf("PaymentsList.json.Marshal err: %v", err), http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	if _, err = resp.Write(respBytes); err != nil {
		slog.Error(fmt.Sprintf("PaymentsList.resp.Write err: %v", err))
	}
}
