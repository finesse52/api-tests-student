package fixtures

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	favoriteManager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/internal/managers/response"
)

// FillFavoritesToLimit adds available foreign advertisements until the user's favorites reach the API limit.
func FillFavoritesToLimit(t *testing.T, token, userID, excludedID string, initial models.List) models.List {
	t.Helper()
	if initial.Total == models.FavoriteLimit {
		return initial
	}

	existing := make(map[string]struct{}, len(initial.Items))
	for _, item := range initial.Items {
		existing[item.Advertisement.ID] = struct{}{}
	}
	catalog := advertisements.GetAdvertisementsListAll(t, token, "")
	candidates := make([]string, 0, len(catalog.Items))
	seen := make(map[string]struct{}, len(catalog.Items))
	for index := len(catalog.Items) - 1; index >= 0; index-- {
		item := catalog.Items[index]
		if item.ID == "" || item.ID == excludedID || item.UserID == "" || item.UserID == userID ||
			strings.HasPrefix(item.Title, testAdvertisementPrefix) {
			continue
		}
		if _, present := existing[item.ID]; present {
			continue
		}
		if _, present := seen[item.ID]; present {
			continue
		}
		seen[item.ID] = struct{}{}
		candidates = append(candidates, item.ID)
	}
	require.GreaterOrEqualf(t, len(candidates), models.FavoriteLimit-initial.Total,
		"в каталоге недостаточно чужих объявлений для достижения лимита избранного")

	current := initial
	for next := 0; current.Total < models.FavoriteLimit; {
		if next == len(candidates) {
			t.Fatalf("не удалось заполнить избранное: %d/%d, кандидаты закончились", current.Total, models.FavoriteLimit)
		}
		batchEnd := min(next+models.FavoriteLimit-current.Total, len(candidates))
		batch := candidates[next:batchEnd]
		next = batchEnd
		result := RetryServiceUnavailable(t, "POST /favorites/bulk", func() response.Result {
			return favoriteManager.BulkAdd(t, token, batch)
		})
		result.RequireStatus(t, http.StatusOK)
		var merged models.BulkAddResult
		result.Decode(t, &merged)
		t.Logf("Подготовка лимита: добавлено %d, не найдено %d, пропущено %d", len(merged.Added), len(merged.NotFound), len(merged.Skipped))
		current = waitFavoritesTotal(t, token, current.Total+len(merged.Added))
	}
	return current
}

// waitFavoritesTotal polls until the page reflects at least the expected number of favorites.
func waitFavoritesTotal(t *testing.T, token string, expected int) models.List {
	t.Helper()
	deadline := time.Now().Add(retryTimeout)
	for {
		current := CurrentFavorites(t, token)
		if current.Total >= expected {
			return current
		}
		if time.Now().After(deadline) {
			t.Fatalf("избранное не достигло ожидаемого количества %d за %s: получено %d", expected, retryTimeout, current.Total)
		}
		time.Sleep(retryInterval)
	}
}
