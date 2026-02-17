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
	"golang.org/x/crypto/bcrypt"
)

func Test_Register(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	reqBody := &models.RegisterUserRequest{
		Login:    "hllwwwrld",
		Password: "aboba",
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

	respBody := &models.RegisterUserResponse{}
	err = json.Unmarshal(resp.Body.Bytes(), respBody)
	require.NoError(t, err)
	require.Equal(t, respBody.Success, true)

	user, err := usersTable.GetByLogin(ctx, reqBody.Login)
	require.NoError(t, err)
	require.NotNil(t, user)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password))
	require.NoError(t, err)

	require.NotEmpty(t, user.ID)
	require.Equal(t, reqBody.Login, user.Login)
	require.Zero(t, user.MonthlyIncome)
	require.False(t, user.CreatedAt.IsZero())
	require.False(t, user.UpdatedAt.IsZero())
	require.Equal(t, user.CreatedAt, user.UpdatedAt)
}
