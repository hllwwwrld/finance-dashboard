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
	"github.com/finance-dashboard/backend/internal/pkg/test_helpers"
	"github.com/stretchr/testify/require"
)

func Test_Login(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Создание юзера
	login, password := test_helpers.RegisterUser(t, ctx, usersService)

	// Успешный логин
	t.Run("successful login", func(t *testing.T) {
		t.Parallel()

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

		authCookie := test_helpers.FindCookie(cookies, "auth_token")
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
			Login:    login,
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

		respBody := &models.LoginUserResponse{}
		err = json.Unmarshal(resp.Body.Bytes(), respBody)
		require.NoError(t, err)
		require.Equal(t, respBody.Success, false)

		// Проверяем, что cookie не установлен при неуспешном логине
		cookies := resp.Result().Cookies()
		authCookie := test_helpers.FindCookie(cookies, "auth_token")
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
