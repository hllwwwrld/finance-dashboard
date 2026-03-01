package users

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func (i *Implementation) Logout(resp http.ResponseWriter, _ *http.Request) {
	authCookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1, // Удаляем куку
	}
	http.SetCookie(resp, authCookie)

	resp.Header().Set("Content-Type", "application/json")
	if _, err := resp.Write([]byte("{success: true")); err != nil {
		slog.Error(fmt.Sprintf("Logout.Write err: %v", err))
	}
}
