package users_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/finance-dashboard/backend/internal/config"
	"github.com/finance-dashboard/backend/internal/pkg/models"
	"github.com/stretchr/testify/require"
)

func Test_Login(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Создание юзера
	registerReqBody := &models.RegisterUserRequest{
		Login:    "testuser_login",
		Password: "testpassword",
	}
	registerJsonBody, err := json.Marshal(registerReqBody)
	require.NoError(t, err)

	registerReq, err := http.NewRequestWithContext(
		ctx, "POST", config.UserRegisterEndpoint, bytes.NewBuffer(registerJsonBody),
	)
	require.NoError(t, err)

	registerResp := httptest.NewRecorder()
	usersService.Register(registerResp, registerReq)
	require.Equal(t, registerResp.Code, http.StatusOK)

	// Успешный логин
	t.Run("successful login", func(t *testing.T) {
		t.Parallel()

		reqBody := &models.LoginUserRequest{
			Login:    "testuser_login",
			Password: "testpassword",
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

		respBody := &models.RegisterUserResponse{}
		err = json.Unmarshal(resp.Body.Bytes(), respBody)
		require.NoError(t, err)
		require.Equal(t, respBody.Success, true)

		cookies := resp.Result().Cookies()
		require.NotEmpty(t, cookies)

		authCookie := findCookie(cookies, "auth_token")
		require.NotNil(t, authCookie)
		require.NotEmpty(t, authCookie.Value)

		require.Equal(t, authCookie.Path, "/")
		require.True(t, authCookie.HttpOnly)
		require.Equal(t, authCookie.SameSite, http.SameSiteLaxMode)
	})

	// Авторизация не авторизация, если пароль не совпадает
	t.Run("failed login with wrong password", func(t *testing.T) {
		t.Parallel()

		reqBody := &models.LoginUserRequest{
			Login:    "testuser_login",
			Password: "wrongpassword",
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

		respBody := &models.RegisterUserResponse{}
		err = json.Unmarshal(resp.Body.Bytes(), respBody)
		require.NoError(t, err)
		require.Equal(t, respBody.Success, false)

		// Проверяем, что cookie не установлен при неуспешном логине
		cookies := resp.Result().Cookies()
		authCookie := findCookie(cookies, "auth_token")
		require.Nil(t, authCookie)
	})

	// Юзер не найден -- не авторизует
	// todo тут мне не нравится internal, когда клиент не найден,
	// todo надо бы обработчик от методов базы, чтобы маппить "Not found" и ошибки вызова
	t.Run("failed login with non-existent user", func(t *testing.T) {
		t.Parallel()

		reqBody := &models.LoginUserRequest{
			Login:    "nonexistentuser",
			Password: "anypassword",
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
		// При несуществующем пользователе хендлер возвращает InternalServerError
		require.Equal(t, resp.Code, http.StatusInternalServerError)
	})
}

// findCookie поиск куки в слайсе
func findCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}
