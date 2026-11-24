package fixtures

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	favoriteManager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/internal/managers/response"
)

// FillFavoritesToLimit creates test-owned advertisements and fills the user's favorites to the API limit.
func FillFavoritesToLimit(t *testing.T, viewerToken, ownerToken, photoFile string, initial models.List) models.List {
	t.Helper()
	missing := models.FavoriteLimit - initial.Total
	if missing == 0 {
		return initial
	}
	t.Logf("Подготовка лимита: создаём %d объявлений второго пользователя", missing)
	ids := make([]string, 0, missing)
	for range missing {
		request := NewCreateRequest(t, photoFile)
		created := CreateOwnedAdvertisement(t, ownerToken, request)
		ids = append(ids, created.ID)
	}
	result := RetryServiceUnavailable(t, "POST /favorites/bulk", func() response.Result {
		return favoriteManager.BulkAdd(t, viewerToken, ids)
	})
	result.RequireStatus(t, http.StatusOK)
	var merged models.BulkAddResult
	result.Decode(t, &merged)
	require.Empty(t, merged.NotFound, "пакетное добавление не нашло созданные тестом объявления")
	t.Logf("Подготовка лимита: пакетное добавление вернуло %d новых и %d пропущенных объявлений", len(merged.Added), len(merged.Skipped))
	return waitFavoritesTotal(t, viewerToken, models.FavoriteLimit)
}

// waitFavoritesTotal polls until the page reflects the expected number of favorites.
func waitFavoritesTotal(t *testing.T, token string, expected int) models.List {
	t.Helper()
	deadline := time.Now().Add(retryTimeout)
	for {
		current := CurrentFavorites(t, token)
		if current.Total == expected {
			return current
		}
		if time.Now().After(deadline) {
			t.Fatalf("избранное не достигло ожидаемого количества %d за %s: получено %d", expected, retryTimeout, current.Total)
		}
		time.Sleep(retryInterval)
	}
}
