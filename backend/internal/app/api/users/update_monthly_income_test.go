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
	"golang.org/x/crypto/bcrypt"
)

// todo вынести создание и логин пользака в отдельную функцию и использовать как прекондишен степ
func Test_UpdateMonthlyIncome(t *testing.T) {
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

	loginResp := httptest.NewRecorder()

	loginBody := &models.LoginUserRequest{
		Login:    "testuser_login",
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

	updateMonthlyIncomeResp := httptest.NewRecorder()
	updateMonthlyIncomeBody := &models.UpdateMonthlyIncomeRequest{
		Income: 261000,
	}
	updateMonthlyIncomeBytes, err := json.Marshal(updateMonthlyIncomeBody)
	require.NoError(t, err)
	require.NotEmpty(t, loginBody)

	updateMonthlyIncomeReq, err := http.NewRequestWithContext(
		ctx, "POST", config.UserRegisterEndpoint, bytes.NewBuffer(updateMonthlyIncomeBytes),
	)
	require.NoError(t, err)
	require.NotNil(t, updateMonthlyIncomeReq)

	require.Len(t, loginResp.Result().Cookies(), 1)
	updateMonthlyIncomeReq.AddCookie(loginResp.Result().Cookies()[0])

	updateMonthlyIncomeWithMiddleware := middlewares.Auth(usersService.UpdateMonthlyIncome)
	updateMonthlyIncomeWithMiddleware(updateMonthlyIncomeResp, updateMonthlyIncomeReq)
	require.Equal(t, http.StatusOK, updateMonthlyIncomeResp.Code)

	user, err := usersTable.GetByLogin(ctx, loginBody.Login)
	require.NoError(t, err)
	require.NotNil(t, user)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginBody.Password))
	require.NoError(t, err)

	require.NotEmpty(t, user.ID)
	require.Equal(t, loginBody.Login, user.Login)
	require.Equal(t, updateMonthlyIncomeBody.Income, user.MonthlyIncome)
	require.False(t, user.CreatedAt.IsZero())
	require.False(t, user.UpdatedAt.IsZero())
	require.NotEqual(t, user.CreatedAt, user.UpdatedAt)
}
