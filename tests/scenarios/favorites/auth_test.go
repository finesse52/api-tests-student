package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	favoriteManager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/response"
	"api-tests-template/internal/utils"
)

func TestFavoritesRequireAuthorization(t *testing.T) {
	utils.LoadEnv()
	for _, auth := range []struct {
		name    string
		token   string
		message string
	}{
		{"без токена", "", "Authorization header required"},
		{"неверный токен", "incorrect_token", "Invalid or expired token"},
	} {
		t.Run(auth.name, func(t *testing.T) {
			for _, operation := range []struct {
				name string
				call func(*testing.T, string) response.Result
			}{
				{"GET /favorites", func(t *testing.T, token string) response.Result {
					return favoriteManager.List(t, token, "")
				}},
				{"POST /favorites", func(t *testing.T, token string) response.Result {
					return favoriteManager.Add(t, token, missingAdvertisementID)
				}},
				{"DELETE /favorites/{ad_id}", func(t *testing.T, token string) response.Result {
					return favoriteManager.Delete(t, token, missingAdvertisementID)
				}},
			} {
				t.Run(operation.name, func(t *testing.T) {
					t.Logf("Подготовка: %s", auth.name)
					t.Logf("Шаг: вызываем %s", operation.name)
					result := operation.call(t, auth.token)
					apiError := result.RequireError(t, http.StatusUnauthorized, "unauthorized")
					require.Equal(t, auth.message, apiError.Message, "неверное сообщение об ошибке авторизации")
					t.Logf("Результат: HTTP %d, message=%q", result.StatusCode, apiError.Message)
				})
			}
		})
	}
}
