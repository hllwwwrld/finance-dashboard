package payments

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/finance-dashboard/backend/internal/pkg/models"
)

func (i *Implementation) DeletePayment(resp http.ResponseWriter, req *http.Request) {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(resp, fmt.Sprintf("PaymentsList.io.ReadAll err: %v", err), http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	reqBody := &models.DeletePaymentRequest{}
	err = json.Unmarshal(bodyBytes, reqBody)
	if err != nil {
		http.Error(resp, fmt.Sprintf("PaymentsList.json.Unmarshal err: %v", err), http.StatusInternalServerError)
		return
	}

	err = i.paymentsTable.DeleteByID(req.Context(), reqBody.ID)
	if err != nil {
		http.Error(resp, fmt.Sprintf("usersTable.Create err: %v", err), http.StatusInternalServerError)
		return
	}

	respBytes, err := json.Marshal(
		models.DeletePaymentResponse{
			Success: true,
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
