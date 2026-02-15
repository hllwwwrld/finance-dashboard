package users

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/finance-dashboard/backend/internal/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

func (i *Implementation) Register(resp http.ResponseWriter, req *http.Request) {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(resp, fmt.Sprintf("PaymentsList.io.ReadAll err: %v", err), http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	reqBody := &models.RegisterUserRequest{}
	err = json.Unmarshal(bodyBytes, reqBody)
	if err != nil {
		http.Error(resp, fmt.Sprintf("PaymentsList.json.Unmarshal err: %v", err), http.StatusInternalServerError)
		return
	}

	password, err := bcrypt.GenerateFromPassword(
		[]byte(reqBody.Password),
		bcrypt.DefaultCost, // сложность хеширования (можно увеличить)
	)
	if err != nil {
		http.Error(resp, fmt.Sprintf("bcrypt.GenerateFromPassword err: %v", err), http.StatusInternalServerError)
	}

	user := models.User{
		Login:    reqBody.Login,
		Password: string(password),
	}
	_, err = i.usersTable.Create(req.Context(), user)
	if err != nil {
		http.Error(resp, fmt.Sprintf("usersTable.Create err: %v", err), http.StatusInternalServerError)
	}

	respBytes, err := json.Marshal(models.RegisterUserResponse{Success: true})
	if err != nil {
		http.Error(resp, fmt.Sprintf("PaymentsList.json.Marshal err: %v", err), http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	if _, err = resp.Write(respBytes); err != nil {
		slog.Error(fmt.Sprintf("PaymentsList.resp.Write err: %v", err))
	}
}
