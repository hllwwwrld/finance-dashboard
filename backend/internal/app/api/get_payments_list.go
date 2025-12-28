package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/finance-dashboard/backend/internal/pkg/models"
)

func (i *Implementation) PaymentsList(resp http.ResponseWriter, req *http.Request) {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		// todo пиздец одно дублирование, надо сделать мидлварь с логами запросов/ответов,
		// todo чтобы не приходилось так писать (смотреть это и ниже)
		errText := fmt.Sprintf("PaymentsList.io.ReadAll err: %v", err)
		slog.Error(errText)

		http.Error(resp, errText, http.StatusInternalServerError)
	}
	defer req.Body.Close()

	reqBody := &models.PaymentsListRequest{}
	err = json.Unmarshal(bodyBytes, reqBody)
	if err != nil {
		errText := fmt.Sprintf("PaymentsList.json.Unmarshal err: %v", err)
		slog.Error(errText)

		http.Error(resp, errText, http.StatusInternalServerError)
	}

	userPayments, err := i.paymentsTable.GetByUserID(req.Context(), reqBody.UserID)
	if err != nil {
		errText := fmt.Sprintf("PaymentsList.i.paymentsTable.GetByUserID err: %v", err)
		slog.Error(errText)

		http.Error(resp, errText, http.StatusInternalServerError)
	}

	respBytes, err := json.Marshal(userPayments)
	if err != nil {
		errText := fmt.Sprintf("PaymentsList.json.Marshal err: %v", err)
		slog.Error(errText)

		http.Error(resp, errText, http.StatusInternalServerError)
	}

	resp.Header().Set("Content-Type", "application/json")
	if _, err = resp.Write(respBytes); err != nil {
		slog.Error(fmt.Sprintf("PaymentsList.resp.Write err: %v", err))
	}
}
