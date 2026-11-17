package fixtures

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/internal/managers/response"
)

// AddFavorite adds an advertisement and registers cleanup if the test exits early.
func AddFavorite(t *testing.T, token, id string) {
	t.Helper()
	t.Logf("Подготовка: добавляем объявление %s в избранное", id)
	result := favorites.Add(t, token, id)
	result.RequireStatus(t, http.StatusCreated)
	var added models.Item
	result.Decode(t, &added)
	require.Equal(t, id, added.Advertisement.ID, "в ответе на добавление указан другой ad_id")
	require.NotEmpty(t, added.FavoriteID, "сервер не вернул favorite_id")
	t.Cleanup(func() {
		cleanupDeletedEntity(t, "избранное для "+id, func() response.Result {
			return favorites.Delete(t, token, id)
		})
	})
	t.Log("Результат подготовки: объявление добавлено в избранное")
}

// RequireFavorite polls a filtered favorites page until the expected state appears.
func RequireFavorite(t *testing.T, token, title, id string, want bool) {
	t.Helper()
	deadline := time.Now().Add(retryTimeout)
	for {
		result := favorites.List(t, token, title)
		if result.StatusCode == http.StatusOK {
			var list models.List
			result.Decode(t, &list)
			found := false
			for _, item := range list.Items {
				if item.Advertisement.ID == id {
					require.Equal(t, title, item.Advertisement.Title, "название избранного объявления")
					found = true
				}
			}
			if found == want {
				return
			}
		} else if result.StatusCode != http.StatusServiceUnavailable {
			result.RequireStatus(t, http.StatusOK)
		}
		if time.Now().After(deadline) {
			t.Fatalf("объявление %s: ожидалось наличие в избранном=%t; последний ответ: HTTP %d, %s", id, want, result.StatusCode, result.Body)
		}
		time.Sleep(retryInterval)
	}
}
