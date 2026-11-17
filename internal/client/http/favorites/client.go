package favorites

import (
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/favorites/models"
)

// HttpList retrieves a user's favorites filtered by a unique advertisement title.
func HttpList(t *testing.T, token, title string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).CreateWithTimeout().Get(path.FavoritesPath).
		Query("search", title).Query("page_size", strconv.Itoa(models.FavoriteLimit)).Expect(t).End().Response
}

// HttpAdd adds an advertisement to the authorized user's favorites.
func HttpAdd(t *testing.T, token, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).CreateWithTimeout().Post(path.FavoritesPath).
		ContentType("application/json").JSON(models.AddRequest{AdvertisementID: id}).
		Expect(t).End().Response
}

// HttpBulkAdd merges a batch of advertisement IDs into favorites using the idempotent bulk endpoint.
func HttpBulkAdd(t *testing.T, token string, ids []string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).CreateWithTransientTransportRetry().Post(path.FavoritesBulkPath).
		ContentType("application/json").JSON(models.BulkAddRequest{AdvertisementIDs: ids}).
		Expect(t).End().Response
}

// HttpDelete removes an advertisement from the authorized user's favorites.
func HttpDelete(t *testing.T, token, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).CreateWithTimeout().Delete(path.FavoritesPath + "/" + id).
		Expect(t).End().Response
}
