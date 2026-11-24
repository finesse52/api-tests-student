package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	favoriteManager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/response"
	"api-tests-template/tests/fixtures"
)

func TestCannotAddOwnAdvertisement(t *testing.T) {
	owner := fixtures.Login(t, fixtures.PrimaryLoginEnv, fixtures.PrimaryPasswordEnv)
	request := fixtures.NewCreateRequest(t, validPhoto)
	created := fixtures.CreateOwnedAdvertisement(t, owner.Token, request)

	t.Log("Шаг: владелец пытается добавить своё объявление в избранное")
	result := fixtures.RetryServiceUnavailable(t, "POST /favorites для своего объявления", func() response.Result {
		return favoriteManager.Add(t, owner.Token, created.ID)
	})
	apiError := result.RequireError(t, http.StatusForbidden, "")
	require.Equal(t, "Cannot add own advertisement to favorites", apiError.Message, "неверное сообщение для собственного объявления")
	fixtures.RequireFavorite(t, owner.Token, request.Title, created.ID, false)
	t.Logf("HTTP 403, message=%q", apiError.Message)
}
