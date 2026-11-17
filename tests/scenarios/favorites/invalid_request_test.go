package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	favoriteManager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/response"
	"api-tests-template/tests/fixtures"
)

func TestCannotAddFavoriteWithoutAdvertisementID(t *testing.T) {
	viewer := fixtures.Login(t, fixtures.PrimaryLoginEnv, fixtures.PrimaryPasswordEnv)
	t.Log("Подготовка: тело POST /favorites содержит пустой advertisement_id")
	t.Log("Шаг: отправляем запрос на добавление в избранное")
	result := fixtures.RetryServiceUnavailable(t, "POST /favorites без advertisement_id", func() response.Result {
		return favoriteManager.Add(t, viewer.Token, "")
	})
	apiError := result.RequireError(t, http.StatusBadRequest, "invalid_request")
	require.Equal(t, "invalid UUID length: 0", apiError.Message, "неверное сообщение для пустого advertisement_id")
	t.Logf("Результат: пустой advertisement_id отклонён, HTTP 400, message=%q", apiError.Message)
}
