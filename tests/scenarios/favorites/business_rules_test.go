package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	favoriteManager "api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/fixtures"
)

func TestCannotAddSameAdvertisementTwice(t *testing.T) {
	viewer := fixtures.Login(t, fixtures.PrimaryLoginEnv, fixtures.PrimaryPasswordEnv)
	owner := fixtures.Login(t, fixtures.SecondLoginEnv, fixtures.SecondPasswordEnv)
	request := fixtures.NewCreateRequest(t, 1)
	created := fixtures.CreateOwnedAdvertisement(t, owner.Token, request)
	fixtures.AddFavorite(t, viewer.Token, created.ID)

	t.Log("Шаг: повторно добавляем то же объявление в избранное")
	result := favoriteManager.Add(t, viewer.Token, created.ID)
	apiError := result.RequireError(t, http.StatusConflict, "")
	require.Equal(t, "Advertisement is already in favorites", apiError.Message, "неверное сообщение для дубликата")
	fixtures.RequireFavorite(t, viewer.Token, request.Title, created.ID, true)
	t.Logf("Результат: дубликат отклонён, HTTP 409, message=%q", apiError.Message)
}

func TestCannotAddOwnAdvertisement(t *testing.T) {
	owner := fixtures.Login(t, fixtures.PrimaryLoginEnv, fixtures.PrimaryPasswordEnv)
	request := fixtures.NewCreateRequest(t, 1)
	created := fixtures.CreateOwnedAdvertisement(t, owner.Token, request)

	t.Log("Шаг: владелец пытается добавить своё объявление в избранное")
	result := favoriteManager.Add(t, owner.Token, created.ID)
	apiError := result.RequireError(t, http.StatusForbidden, "")
	require.Equal(t, "Cannot add own advertisement to favorites", apiError.Message, "неверное сообщение для собственного объявления")
	fixtures.RequireFavorite(t, owner.Token, request.Title, created.ID, false)
	t.Logf("Результат: своё объявление отклонено, HTTP 403, message=%q", apiError.Message)
}
