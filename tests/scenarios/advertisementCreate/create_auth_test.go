package advertisementCreate

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/response"
	"api-tests-template/internal/utils"
	"api-tests-template/tests/fixtures"
)

func TestCreateAdvertisementRequiresAuthorization(t *testing.T) {
	utils.LoadEnv()
	for _, tc := range []struct {
		name    string
		token   string
		message string
	}{
		{"без токена", "", fixtures.MissingAuthMessage},
		{"неверный токен", fixtures.InvalidToken, fixtures.InvalidAuthMessage},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request := fixtures.NewCreateRequest(t, validPhoto1)
			t.Logf("Подготовка: запрос на создание объявления %s", tc.name)
			result := fixtures.RetryServiceUnavailable(t, "POST /advertisement без корректной авторизации", func() response.Result {
				return advertisement.CreateExpectedError(t, tc.token, request)
			})
			t.Log("Шаг: проверяем отказ в создании объявления")
			apiError := result.RequireError(t, http.StatusUnauthorized, fixtures.UnauthorizedError)
			require.Equal(t, tc.message, apiError.Message, "неверное сообщение об ошибке авторизации")
			t.Logf("Результат: HTTP 401, message=%q", apiError.Message)
		})
	}
}
