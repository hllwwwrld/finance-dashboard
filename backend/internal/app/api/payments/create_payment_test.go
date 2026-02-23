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

func Test_CreatePayment(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Создание юзера
	registerReqBody := &models.RegisterUserRequest{
		Login:    fmt.Sprintf("Test_CreatePayment_%d", time.Now().UnixNano()),
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

	dueDate := time.Now().Add(time.Hour * 48).Format("2006-01-02")
	// Создание платежа
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

	fmt.Println(createReq.Context().Value(middlewares.UserContextKey))

	createRespBody := &models.CreatePaymentResponse{}
	err = json.Unmarshal(createResp.Body.Bytes(), createRespBody)
	require.NoError(t, err)

	require.NotEmpty(t, createRespBody.ID)
	require.NotEmpty(t, createRespBody.UserID)
	require.Equal(t, createBody.Name, createRespBody.Name)
	require.Equal(t, createBody.Amount, createRespBody.Amount)
	require.Equal(t, createBody.Category, createRespBody.Category)
	require.Equal(t, createBody.Color, createRespBody.Color)
	require.Equal(t, 2, createRespBody.DaysUntil)

	// Проверка, что платеж есть в таблице
	user, err := usersTable.GetByLogin(ctx, registerReqBody.Login)
	require.NoError(t, err)
	require.NotNil(t, user)

	userPayments, err := paymentsTable.GetByUserID(ctx, user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, userPayments)

	require.Len(t, userPayments, 1)

	payment := userPayments[0]
	require.NotEmpty(t, payment.ID)
	require.Equal(t, user.ID, payment.UserID)
	require.Equal(t, payment.Name, createBody.Name)
	require.Equal(t, createBody.Amount, payment.Amount)
	require.Equal(t, createBody.DueDate, payment.DueDate.Format("2006-01-02"))
	require.Equal(t, createBody.Category, payment.Category)
	require.Equal(t, createBody.Color, payment.Color)
	require.False(t, payment.CreatedAt.IsZero())
	require.False(t, payment.UpdatedAt.IsZero())
	require.Equal(t, payment.CreatedAt, payment.CreatedAt)
}

// Test_CreatePayment_Unauthorized проверяет, что без авторизации
// создание платежа возвращает 401
func Test_CreatePayment_Unauthorized(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	createResp := httptest.NewRecorder()
	createBody := &models.CreatePaymentRequest{
		Name:     "Unauthorized payment",
		Amount:   100,
		DueDate:  "2026-01-10",
		Category: "test",
		Color:    "#000000",
	}
	createBytes, err := json.Marshal(createBody)
	require.NoError(t, err)
	require.NotEmpty(t, createBytes)

	createReq, err := http.NewRequestWithContext(
		ctx, http.MethodPost, config.PaymentsCreate, bytes.NewBuffer(createBytes),
	)
	require.NoError(t, err)
	require.NotNil(t, createReq)

	// без куки авторизации
	createWithMiddleware := middlewares.Auth(paymentsService.CreatePayment)
	createWithMiddleware(createResp, createReq)

	require.Equal(t, http.StatusUnauthorized, createResp.Code)
}
