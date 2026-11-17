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

// AddFavorite adds an advertisement, confirms uncertain gateway responses, and registers cleanup.
func AddFavorite(t *testing.T, token, title, id string) {
	t.Helper()
	t.Logf("Подготовка: добавляем объявление %s в избранное", id)
	t.Cleanup(func() {
		cleanupDeletedEntity(t, "избранное для "+id, func() response.Result {
			return favorites.Delete(t, token, id)
		})
	})

	deadline := time.Now().Add(retryTimeout)
	sawGatewayError := false
	for {
		result := favorites.Add(t, token, id)
		switch result.StatusCode {
		case http.StatusCreated:
			var added models.Item
			result.Decode(t, &added)
			require.Equal(t, id, added.Advertisement.ID, "в ответе на добавление указан другой ad_id")
			require.NotEmpty(t, added.FavoriteID, "сервер не вернул favorite_id")
			t.Log("Результат подготовки: объявление добавлено в избранное, HTTP 201")
			return
		case http.StatusServiceUnavailable:
			sawGatewayError = true
			t.Logf("Временный HTTP 503 при добавлении %s; проверяем состояние избранного", id)
		case http.StatusConflict:
			if !sawGatewayError {
				result.RequireStatus(t, http.StatusCreated)
			}
		default:
			result.RequireStatus(t, http.StatusCreated)
		}

		listResult := favorites.List(t, token, title)
		if listResult.StatusCode == http.StatusOK && favoritePresent(t, listResult, title, id) {
			t.Logf("Результат подготовки: объявление %s найдено в избранном после неоднозначного HTTP %d", id, result.StatusCode)
			return
		}
		if listResult.StatusCode != http.StatusOK && listResult.StatusCode != http.StatusServiceUnavailable {
			listResult.RequireStatus(t, http.StatusOK)
		}
		if time.Now().After(deadline) {
			t.Fatalf("не удалось подтвердить добавление %s в избранное за %s: POST HTTP %d, GET HTTP %d", id, retryTimeout, result.StatusCode, listResult.StatusCode)
		}
		time.Sleep(retryInterval)
	}
}

// RequireFavorite polls a filtered favorites page until the expected state appears.
func RequireFavorite(t *testing.T, token, title, id string, want bool) {
	t.Helper()
	deadline := time.Now().Add(retryTimeout)
	for {
		result := favorites.List(t, token, title)
		if result.StatusCode == http.StatusOK {
			if favoritePresent(t, result, title, id) == want {
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

// favoritePresent checks whether a successful favorites response contains the expected advertisement.
func favoritePresent(t *testing.T, result response.Result, title, id string) bool {
	t.Helper()
	var list models.List
	result.Decode(t, &list)
	for _, item := range list.Items {
		if item.Advertisement.ID == id {
			require.Equal(t, title, item.Advertisement.Title, "название избранного объявления")
			return true
		}
	}
	return false
}

// CurrentFavorites reads the complete first page; the API caps each user's favorites at 50.
func CurrentFavorites(t *testing.T, token string) models.List {
	t.Helper()
	result := RetryServiceUnavailable(t, "GET /favorites", func() response.Result {
		return favorites.List(t, token, "")
	})
	result.RequireStatus(t, http.StatusOK)
	var list models.List
	result.Decode(t, &list)
	require.LessOrEqual(t, list.Total, models.FavoriteLimit, "API вернул больше элементов избранного, чем разрешено")
	require.Len(t, list.Items, list.Total, "список избранного должен помещаться на одну страницу")
	return list
}

// PreserveFavorites removes only entries added after the test's initial snapshot.
func PreserveFavorites(t *testing.T, token string) models.List {
	t.Helper()
	initial := CurrentFavorites(t, token)
	existing := make(map[string]struct{}, len(initial.Items))
	for _, item := range initial.Items {
		existing[item.Advertisement.ID] = struct{}{}
	}
	t.Cleanup(func() {
		current := CurrentFavorites(t, token)
		for _, item := range current.Items {
			id := item.Advertisement.ID
			if _, present := existing[id]; present {
				continue
			}
			cleanupDeletedEntity(t, "избранное для "+id, func() response.Result {
				return favorites.Delete(t, token, id)
			})
		}
	})
	return initial
}
