package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	favoriteManager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/internal/managers/response"
	"api-tests-template/tests/fixtures"
)

func TestCannotAddFavoriteBeyondLimit(t *testing.T) {
	viewer := fixtures.Login(t, fixtures.PrimaryLoginEnv, fixtures.PrimaryPasswordEnv)
	owner := fixtures.Login(t, fixtures.SecondLoginEnv, fixtures.SecondPasswordEnv)
	request := fixtures.NewCreateRequest(t, validPhoto)
	created := fixtures.CreateOwnedAdvertisement(t, owner.Token, request)
	initial := fixtures.PreserveFavorites(t, viewer.Token)
	t.Logf("Подготовка: исходное избранное содержит %d объявлений", initial.Total)
	full := fixtures.FillFavoritesToLimit(t, viewer.Token, owner.Token, validPhoto, initial)
	require.Equal(t, models.FavoriteLimit, full.Total, "не удалось подготовить избранное на границе лимита")
	t.Logf("Результат подготовки: достигнут лимит %d избранных", full.Total)

	t.Log("Шаг: добавляем ещё одно чужое объявление")
	result := fixtures.RetryServiceUnavailable(t, "POST /favorites сверх лимита", func() response.Result {
		return favoriteManager.Add(t, viewer.Token, created.ID)
	})
	apiError := result.RequireError(t, http.StatusBadRequest, "limit_reached")
	require.Equal(t, "Favorites limit of 50 reached", apiError.Message, "неверное сообщение об ограничении избранного")
	fixtures.RequireFavorite(t, viewer.Token, request.Title, created.ID, false)
	require.Equal(t, models.FavoriteLimit, fixtures.CurrentFavorites(t, viewer.Token).Total,
		"после отклонённого запроса количество избранных изменилось")
	t.Logf("Результат: лимит соблюдён, HTTP 400, message=%q", apiError.Message)
}
