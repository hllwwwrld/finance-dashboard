package users

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/finance-dashboard/backend/internal/pkg/middlewares"
	"github.com/finance-dashboard/backend/internal/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

func (i *Implementation) Login(resp http.ResponseWriter, req *http.Request) {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(resp, fmt.Sprintf("PaymentsList.io.ReadAll err: %v", err), http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	reqBody := &models.LoginUserRequest{}
	err = json.Unmarshal(bodyBytes, reqBody)
	if err != nil {
		http.Error(resp, fmt.Sprintf("PaymentsList.json.Unmarshal err: %v", err), http.StatusInternalServerError)
		return
	}

	user, err := i.usersTable.GetByLogin(req.Context(), reqBody.Login)
	// todo научиться ловить ошибку, что пользователь не найден и возвращать 404
	if err != nil {
		http.Error(resp, fmt.Sprintf("usersTable.Create err: %v", err), http.StatusInternalServerError)
		return
	}

	successLogin := true
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password))
	if err != nil {
		successLogin = false
	}

	respBytes, err := json.Marshal(models.RegisterUserResponse{Success: successLogin})
	if err != nil {
		http.Error(resp, fmt.Sprintf("PaymentsList.json.Marshal err: %v", err), http.StatusUnauthorized)
		return
	}

	if !successLogin {
		resp.Header().Set("Content-Type", "application/json")
		if _, err = resp.Write(respBytes); err != nil {
			slog.Error(fmt.Sprintf("PaymentsList.resp.Write err: %v", err))
		}
		return
	}

	token, err := middlewares.GenerateJWT(user.ID, user.Login)
	if err != nil {
		slog.Error(fmt.Sprintf("Login.middlewares.GenerateJWT err: %v", err))
	}

	http.SetCookie(resp, middlewares.DefaultAuthCookie(token))

	resp.Header().Set("Content-Type", "application/json")
	if _, err = resp.Write(respBytes); err != nil {
		slog.Error(fmt.Sprintf("PaymentsList.resp.Write err: %v", err))
	}
}
