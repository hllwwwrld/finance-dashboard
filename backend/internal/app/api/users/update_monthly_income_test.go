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
	"github.com/finance-dashboard/backend/internal/pkg/test_helpers"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// todo вынести создание и логин пользака в отдельную функцию и использовать как прекондишен степ
func Test_UpdateMonthlyIncome(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	updateMonthlyIncomeResp := httptest.NewRecorder()
	updateMonthlyIncomeBody := &models.UpdateMonthlyIncomeRequest{
		Income: 261000,
	}
	updateMonthlyIncomeBytes, err := json.Marshal(updateMonthlyIncomeBody)
	require.NoError(t, err)
	require.NotEmpty(t, updateMonthlyIncomeBody)

	updateMonthlyIncomeReq, err := http.NewRequestWithContext(
		ctx, "POST", config.UserRegisterEndpoint, bytes.NewBuffer(updateMonthlyIncomeBytes),
	)
	require.NoError(t, err)
	require.NotNil(t, updateMonthlyIncomeReq)

	login, password, authCookie := test_helpers.RegisterAndLoginUser(t, ctx, usersService)
	updateMonthlyIncomeReq.AddCookie(authCookie)

	updateMonthlyIncomeWithMiddleware := middlewares.Auth(usersService.UpdateMonthlyIncome)
	updateMonthlyIncomeWithMiddleware(updateMonthlyIncomeResp, updateMonthlyIncomeReq)
	require.Equal(t, http.StatusOK, updateMonthlyIncomeResp.Code)

	user, err := usersTable.GetByLogin(ctx, login)
	require.NoError(t, err)
	require.NotNil(t, user)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	require.NoError(t, err)

	require.NotEmpty(t, user.ID)
	require.Equal(t, login, user.Login)
	require.Equal(t, updateMonthlyIncomeBody.Income, user.MonthlyIncome)
	require.False(t, user.CreatedAt.IsZero())
	require.False(t, user.UpdatedAt.IsZero())
	require.NotEqual(t, user.CreatedAt, user.UpdatedAt)
}
