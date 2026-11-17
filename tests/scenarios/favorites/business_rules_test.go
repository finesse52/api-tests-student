package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	favoriteManager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/response"
	"api-tests-template/tests/fixtures"
)

func TestCannotAddSameAdvertisementTwice(t *testing.T) {
	viewer := fixtures.Login(t, fixtures.PrimaryLoginEnv, fixtures.PrimaryPasswordEnv)
	owner := fixtures.Login(t, fixtures.SecondLoginEnv, fixtures.SecondPasswordEnv)
	request := fixtures.NewCreateRequest(t, 1)
	created := fixtures.CreateOwnedAdvertisement(t, owner.Token, request)
	fixtures.AddFavorite(t, viewer.Token, request.Title, created.ID)

	t.Log("Шаг: повторно добавляем то же объявление в избранное")
	result := fixtures.RetryServiceUnavailable(t, "повторное POST /favorites", func() response.Result {
		return favoriteManager.Add(t, viewer.Token, created.ID)
	})
	apiError := result.RequireError(t, http.StatusConflict, "")
	require.Equal(t, "Advertisement is already in favorites", apiError.Message, "неверное сообщение для дубликата")
	fixtures.RequireFavorite(t, viewer.Token, request.Title, created.ID, true)
	t.Logf("Результат: дубликат отклонён, HTTP 409, message=%q", apiError.Message)
}
