package advertisementCreate

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisement"
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
		{"без токена", "", "Authorization header required"},
		{"неверный токен", "incorrect_token", "Invalid or expired token"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request := fixtures.NewCreateRequest(t, 1)
			t.Logf("Подготовка: запрос на создание объявления %s", tc.name)
			result := advertisement.Create(t, tc.token, request)
			t.Log("Шаг: проверяем отказ в создании объявления")
			apiError := result.RequireError(t, http.StatusUnauthorized, "unauthorized")
			require.Equal(t, tc.message, apiError.Message, "неверное сообщение об ошибке авторизации")
			t.Logf("Результат: HTTP 401, message=%q", apiError.Message)
		})
	}
}
