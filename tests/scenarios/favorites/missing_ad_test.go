package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	favoriteManager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/response"
	"api-tests-template/tests/fixtures"
)

const missingAdvertisementID = "00000000-0000-4000-8000-000000000000"

func TestFavoriteOperationsWithMissingAdvertisement(t *testing.T) {
	viewer := fixtures.Login(t, fixtures.PrimaryLoginEnv, fixtures.PrimaryPasswordEnv)
	for _, tc := range []struct {
		name    string
		message string
		call    func(*testing.T) response.Result
	}{
		{"добавление", "Advertisement not found", func(t *testing.T) response.Result {
			return favoriteManager.Add(t, viewer.Token, missingAdvertisementID)
		}},
		{"удаление", "Favorite entry not found", func(t *testing.T) response.Result {
			return favoriteManager.Delete(t, viewer.Token, missingAdvertisementID)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Подготовка: используем отсутствующий ad_id=%s", missingAdvertisementID)
			t.Logf("Шаг: %s отсутствующего объявления в избранном", tc.name)
			result := fixtures.RetryServiceUnavailable(t, tc.name+" отсутствующего ad_id", func() response.Result {
				return tc.call(t)
			})
			apiError := result.RequireError(t, http.StatusNotFound, "")
			require.Equal(t, tc.message, apiError.Message, "неверное сообщение для отсутствующего ad_id")
			t.Logf("Результат: HTTP %d, message=%q", result.StatusCode, apiError.Message)
		})
	}
}
