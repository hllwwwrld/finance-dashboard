package payments

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/finance-dashboard/backend/internal/pkg/models"
)

func (i *Implementation) CreatePayment(resp http.ResponseWriter, req *http.Request) {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(resp, fmt.Sprintf("PaymentsList.io.ReadAll err: %v", err), http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	reqBody := &models.PaymentsListRequest{}
	err = json.Unmarshal(bodyBytes, reqBody)
	if err != nil {
		http.Error(resp, fmt.Sprintf("PaymentsList.json.Unmarshal err: %v", err), http.StatusInternalServerError)
		return
	}

	userPayments, err := i.paymentsTable.GetByUserID(req.Context(), reqBody.UserID)
	if err != nil {
		http.Error(resp, fmt.Sprintf("PaymentsList.i.paymentsTable.GetByUserID err: %v", err), http.StatusInternalServerError)
		return
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
