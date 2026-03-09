package users_test

import (
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
)

// todo дописать параметризацию, чтобы делать баланс пользака рандомным и потом проверять, что в fetch пришел нужный
func Test_FetchProfile(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	fetchProfileResp := httptest.NewRecorder()
	fetchProfileReq, err := http.NewRequestWithContext(
		ctx, "POST", config.UserFetchProfileEndpoint, nil,
	)
	require.NoError(t, err)
	require.NotNil(t, fetchProfileReq)

	login, _, authCookie := test_helpers.RegisterAndLoginUser(t, ctx, usersService)
	fetchProfileReq.AddCookie(authCookie)

	fetchProfileWithMiddleware := middlewares.Auth(usersService.FetchProfile)
	fetchProfileWithMiddleware(fetchProfileResp, fetchProfileReq)
	require.Equal(t, http.StatusOK, fetchProfileResp.Code)

	fetchProfileRespBody := &models.FetchProfileResponse{}
	err = json.Unmarshal(fetchProfileResp.Body.Bytes(), fetchProfileRespBody)
	require.NoError(t, err)

	require.Equal(t, 0, fetchProfileRespBody.MonthlyIncome)

	user, err := usersTable.GetByLogin(ctx, login)
	require.NoError(t, err)
	require.NotNil(t, user)
}
