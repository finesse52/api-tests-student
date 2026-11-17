package advertisements

import (
	"net/http"
	"testing"
	"time"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/response"
)

const (
	searchRetryTimeout  = 10 * time.Second
	searchRetryInterval = 250 * time.Millisecond
	maxSearchPages      = 100
)

// GetAdvertisementsList возвращает первую страницу объявлений по части названия.
func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	t.Helper()
	return getAdvertisementsListPage(t, token, text, "")
}

// GetAdvertisementsListAll собирает все страницы поиска, чтобы не зависеть от позиции объявления в выдаче.
func GetAdvertisementsListAll(t *testing.T, token, text string) models.AdvertisementsList {
	t.Helper()
	all := models.AdvertisementsList{Items: make([]models.Advertisement, 0)}
	cursor := ""
	seenCursors := make(map[string]struct{})
	for pageNumber := 1; pageNumber <= maxSearchPages; pageNumber++ {
		page := getAdvertisementsListPage(t, token, text, cursor)
		all.Items = append(all.Items, page.Items...)
		if !page.HasNext {
			return all
		}
		if page.NextCursor == "" {
			t.Fatalf("страница %d поиска по %q сообщает has_next=true без next_cursor", pageNumber, text)
		}
		if _, exists := seenCursors[page.NextCursor]; exists {
			t.Fatalf("поиск по %q вернул повторный курсор на странице %d", text, pageNumber)
		}
		seenCursors[page.NextCursor] = struct{}{}
		cursor = page.NextCursor
	}
	t.Fatalf("поиск по %q превысил лимит %d страниц", text, maxSearchPages)
	return all
}

// getAdvertisementsListPage читает одну страницу и повторяет только временный ответ 503.
func getAdvertisementsListPage(t *testing.T, token, text, cursor string) models.AdvertisementsList {
	t.Helper()
	deadline := time.Now().Add(searchRetryTimeout)
	for {
		result := response.Read(t, advertisements.HttpGetAdvertisementsPage(t, token, text, cursor))
		if result.StatusCode == http.StatusServiceUnavailable && time.Now().Before(deadline) {
			t.Logf("Временный HTTP 503 при поиске %q; повторяем GET", text)
			time.Sleep(searchRetryInterval)
			continue
		}
		result.RequireStatus(t, http.StatusOK)
		var page models.AdvertisementsList
		result.Decode(t, &page)
		return page
	}
}
