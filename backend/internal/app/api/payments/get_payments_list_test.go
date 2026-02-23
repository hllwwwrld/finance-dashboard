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
	"github.com/stretchr/testify/require"
)

func Test_GetPaymentList(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Создание юзера
	registerReqBody := &models.RegisterUserRequest{
		Login:    fmt.Sprintf("Test_GetPaymentList_%d", time.Now().UnixNano()),
		Password: "testpassword",
	}
	registerJSONBody, err := json.Marshal(registerReqBody)
	require.NoError(t, err)
	require.NotEmpty(t, registerJSONBody)

	registerReq, err := http.NewRequestWithContext(
		ctx, http.MethodPost, config.UserRegisterEndpoint, bytes.NewBuffer(registerJSONBody),
	)
	require.NoError(t, err)

	registerResp := httptest.NewRecorder()
	usersService.Register(registerResp, registerReq)
	require.Equal(t, http.StatusOK, registerResp.Code)

	// Логин
	loginResp := httptest.NewRecorder()
	loginBody := &models.LoginUserRequest{
		Login:    registerReqBody.Login,
		Password: registerReqBody.Password,
	}
	loginBytes, err := json.Marshal(loginBody)
	require.NoError(t, err)
	require.NotEmpty(t, loginBytes)

	loginReq, err := http.NewRequestWithContext(
		ctx, http.MethodPost, config.UserLoginEndpoint, bytes.NewBuffer(loginBytes),
	)
	require.NoError(t, err)
	require.NotNil(t, loginReq)

	usersService.Login(loginResp, loginReq)
	require.Equal(t, http.StatusOK, loginResp.Code)
	require.Len(t, loginResp.Result().Cookies(), 1)

	authCookie := loginResp.Result().Cookies()[0]

	// Создание платежа
	dueDate := time.Now().Add(time.Hour * 48).Format(time.DateOnly)
	createResp := httptest.NewRecorder()
	createBody := &models.CreatePaymentRequest{
		Name:     "Домашний интернет",
		Amount:   900,
		DueDate:  dueDate,
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

	createReq.AddCookie(authCookie)

	createWithMiddleware := middlewares.Auth(paymentsService.CreatePayment)
	createWithMiddleware(createResp, createReq)
	require.Equal(t, http.StatusOK, createResp.Code)

	createRespBody := &models.CreatePaymentResponse{}
	err = json.Unmarshal(createResp.Body.Bytes(), createRespBody)
	require.NoError(t, err)

	require.NotEmpty(t, createRespBody.ID)
	require.NotEmpty(t, createRespBody.UserID)
	require.Equal(t, createBody.Name, createRespBody.Name)
	require.Equal(t, createBody.Amount, createRespBody.Amount)
	require.Equal(t, createBody.Category, createRespBody.Category)
	require.Equal(t, createBody.Color, createRespBody.Color)
	require.GreaterOrEqual(t, createRespBody.DaysUntil, 0)

	// Проверка, что PaymentsList возвращает созданный платеж
	listResp := httptest.NewRecorder()
	listReq, err := http.NewRequestWithContext(
		ctx, http.MethodGet, config.PaymentsListEndpoint, nil,
	)
	require.NoError(t, err)
	require.NotNil(t, listReq)

	listReq.AddCookie(authCookie)

	listWithMiddleware := middlewares.Auth(paymentsService.PaymentsList)
	listWithMiddleware(listResp, listReq)
	require.Equal(t, http.StatusOK, listResp.Code)

	var listRespBody []*models.Payment
	err = json.Unmarshal(listResp.Body.Bytes(), &listRespBody)
	require.NoError(t, err)
	require.NotEmpty(t, listRespBody)

	var found bool
	for _, payment := range listRespBody {
		if payment.ID == createRespBody.ID {
			found = true
			require.Equal(t, createBody.Name, payment.Name)
			require.Equal(t, createBody.Amount, payment.Amount)
			require.Equal(t, createBody.Category, payment.Category)
			require.Equal(t, createBody.Color, payment.Color)
			require.Equal(t, 2, payment.DaysUntil)
		}
	}
	require.True(t, found, "created payment must be present in PaymentsList response")
}
