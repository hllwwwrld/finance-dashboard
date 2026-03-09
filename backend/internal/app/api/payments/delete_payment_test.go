package payments_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/finance-dashboard/backend/internal/config"
	"github.com/finance-dashboard/backend/internal/pkg/middlewares"
	"github.com/finance-dashboard/backend/internal/pkg/models"
	"github.com/finance-dashboard/backend/internal/pkg/test_helpers"
	"github.com/stretchr/testify/require"
)

func Test_DeletePayment(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	dueDay := time.Now().Add(time.Hour * 48).Day()
	// Создание платежа
	createResp := httptest.NewRecorder()
	createBody := &models.CreatePaymentRequest{
		Name:     "Домашний интернет",
		Amount:   900,
		DueDay:   dueDay,
		Category: "internet",
		Color:    "#ff0000",
	}
	createBytes, err := json.Marshal(createBody)
	require.NoError(t, err)
	require.NotEmpty(t, createBytes)

	createReq, err := http.NewRequestWithContext(
		ctx, http.MethodPost, config.PaymentsCreate, bytes.NewBuffer(createBytes),
	)
	require.NoError(t, err)
	require.NotNil(t, createReq)

	login, _, authCookie := test_helpers.RegisterAndLoginUser(t, ctx, usersService)
	createReq.AddCookie(authCookie)

	createWithMiddleware := middlewares.Auth(paymentsService.CreatePayment)
	createWithMiddleware(createResp, createReq)
	require.Equal(t, http.StatusOK, createResp.Code)

	fmt.Println(createReq.Context().Value(middlewares.UserContextKey))

	createRespBody := &models.CreatePaymentResponse{}
	err = json.Unmarshal(createResp.Body.Bytes(), createRespBody)
	require.NoError(t, err)

	// Проверка, что платеж есть в таблице
	user, err := usersTable.GetByLogin(ctx, login)
	require.NoError(t, err)
	require.NotNil(t, user)

	userPayments, err := paymentsTable.GetByUserID(ctx, user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, userPayments)

	require.Len(t, userPayments, 1)

	// Удаление платежа
	deleteResp := httptest.NewRecorder()
	deleteBody := &models.DeletePaymentRequest{
		ID: userPayments[0].ID,
	}
	deleteBytes, err := json.Marshal(deleteBody)
	require.NoError(t, err)
	require.NotEmpty(t, deleteBytes)

	deleteReq, err := http.NewRequestWithContext(
		ctx, http.MethodPost, config.PaymentsDelete, bytes.NewBuffer(deleteBytes),
	)
	require.NoError(t, err)
	require.NotNil(t, deleteReq)

	deleteReq.AddCookie(authCookie)

	deleteWithMiddleware := middlewares.Auth(paymentsService.DeletePayment)
	deleteWithMiddleware(deleteResp, deleteReq)
	require.Equal(t, http.StatusOK, deleteResp.Code)

	deleteRespBody := &models.DeletePaymentResponse{}
	err = json.Unmarshal(deleteResp.Body.Bytes(), deleteRespBody)
	require.NoError(t, err)

	require.True(t, deleteRespBody.Success)

	userPayments, err = paymentsTable.GetByUserID(ctx, user.ID)
	require.NoError(t, err)
	require.Len(t, userPayments, 0)
}
