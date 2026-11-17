package advertisements

import (
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// HttpGetAdvertisements ищет объявления по части названия
func HttpGetAdvertisements(t *testing.T, token, textForSearch string) *http.Response {
	t.Helper()
	return HttpGetAdvertisementsPage(t, token, textForSearch, "")
}

// HttpGetAdvertisementsPage запрашивает страницу поиска по курсору.
func HttpGetAdvertisementsPage(t *testing.T, token, textForSearch, cursor string) *http.Response {
	t.Helper()
	request := apiRunner.GetRunner().Auth(token).CreateWithTimeout().
		Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		Query("limit", "50")
	if cursor != "" {
		request = request.Query("cursor", cursor)
	}
	return request.Expect(t).End().Response
}
