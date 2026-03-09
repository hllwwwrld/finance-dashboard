package test_helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/finance-dashboard/backend/internal/app/api/users"
	"github.com/finance-dashboard/backend/internal/config"
	"github.com/finance-dashboard/backend/internal/pkg/models"
	"github.com/stretchr/testify/require"
)

func RegisterUser(t *testing.T, ctx context.Context, usersService *users.Implementation) (login, password string) {
	login = time.Now().String()
	password = time.Now().String()
	reqBody := &models.RegisterUserRequest{
		Login:    login,
		Password: password,
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)
	require.NotEmpty(t, jsonBody)

	registerUserReq, err := http.NewRequestWithContext(
		ctx, "POST", config.UserRegisterEndpoint, bytes.NewBuffer(jsonBody),
	)
	require.NoError(t, err)

	resp := httptest.NewRecorder()
	usersService.Register(resp, registerUserReq)
	require.Equal(t, resp.Code, http.StatusOK)

	return login, password
}

func LoginUser(t *testing.T, ctx context.Context, login, password string, usersService *users.Implementation) *http.Cookie {
	reqBody := &models.LoginUserRequest{
		Login:    login,
		Password: password,
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)
	require.NotEmpty(t, jsonBody)

	loginReq, err := http.NewRequestWithContext(
		ctx, "POST", config.UserLoginEndpoint, bytes.NewBuffer(jsonBody),
	)
	require.NoError(t, err)

	resp := httptest.NewRecorder()
	usersService.Login(resp, loginReq)
	require.Equal(t, resp.Code, http.StatusOK)

	respBody := &models.LoginUserResponse{}
	err = json.Unmarshal(resp.Body.Bytes(), respBody)
	require.NoError(t, err)
	require.Equal(t, respBody.Success, true)

	cookies := resp.Result().Cookies()
	require.NotEmpty(t, cookies)

	authCookie := FindCookie(cookies, "auth_token")
	require.NotNil(t, authCookie)
	require.NotEmpty(t, authCookie.Value)

	return authCookie
}

func RegisterAndLoginUser(t *testing.T, ctx context.Context, usersService *users.Implementation) (string, string, *http.Cookie) {
	login, password := RegisterUser(t, ctx, usersService)
	return login, password, LoginUser(t, ctx, login, password, usersService)
}

// FindCookie поиск куки в слайсе
func FindCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}
