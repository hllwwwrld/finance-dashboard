package users_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/finance-dashboard/backend/internal/config"
	"github.com/finance-dashboard/backend/internal/pkg/middlewares"
	"github.com/finance-dashboard/backend/internal/pkg/models"
	"github.com/stretchr/testify/require"
)

// todo дописать параметризацию, чтобы делать баланс пользака рандомным и потом проверять, что в fetch пришел нужный
func Test_FetchProfile(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Создание юзера
	registerReqBody := &models.RegisterUserRequest{
		Login:    "Test_FetchProfile",
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

	loginResp := httptest.NewRecorder()

	loginBody := &models.LoginUserRequest{
		Login:    "Test_FetchProfile",
		Password: "testpassword",
	}
	loginBytes, err := json.Marshal(loginBody)
	require.NoError(t, err)
	require.NotEmpty(t, loginBody)

	loginReq, err := http.NewRequestWithContext(
		ctx, "POST", config.UserRegisterEndpoint, bytes.NewBuffer(loginBytes),
	)
	require.NoError(t, err)
	require.NotNil(t, loginReq)
	usersService.Login(loginResp, loginReq)
	require.Equal(t, loginResp.Code, http.StatusOK)

	fetchProfileResp := httptest.NewRecorder()
	fetchProfileReq, err := http.NewRequestWithContext(
		ctx, "POST", config.UserRegisterEndpoint, nil,
	)
	require.NoError(t, err)
	require.NotNil(t, fetchProfileReq)

	require.Len(t, loginResp.Result().Cookies(), 1)
	fetchProfileReq.AddCookie(loginResp.Result().Cookies()[0])

	fetchProfileWithMiddleware := middlewares.Auth(usersService.FetchProfile)
	fetchProfileWithMiddleware(fetchProfileResp, fetchProfileReq)
	require.Equal(t, http.StatusOK, fetchProfileResp.Code)

	fetchProfileRespBody := &models.FetchProfileResponse{}
	err = json.Unmarshal(fetchProfileResp.Body.Bytes(), fetchProfileRespBody)
	require.NoError(t, err)

	require.Equal(t, 0, fetchProfileRespBody.MonthlyIncome)

	user, err := usersTable.GetByLogin(ctx, loginBody.Login)
	require.NoError(t, err)
	require.NotNil(t, user)
}
