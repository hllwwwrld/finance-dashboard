package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/finance-dashboard/backend/internal/pkg/models"
)

func Ping(w http.ResponseWriter, _ *http.Request) {
	resp := models.PingResponse{Ok: true}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		errText := fmt.Sprintf("Ping.json.Marshal err: %v", err)
		slog.Error(errText)

		http.Error(w, errText, http.StatusInternalServerError)
	}

	if _, err = w.Write(respBytes); err != nil {
		slog.Error(fmt.Sprintf("Ping.w.Write err: %v", err))
	}
}
