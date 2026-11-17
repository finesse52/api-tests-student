package favorites

import (
	"testing"

	"api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/managers/response"
)

// List retrieves a filtered page of the authorized user's favorites.
func List(t *testing.T, token, title string) response.Result {
	t.Helper()
	return response.Read(t, favorites.HttpList(t, token, title))
}

// Add submits an advertisement ID to the authorized user's favorites.
func Add(t *testing.T, token, id string) response.Result {
	t.Helper()
	return response.Read(t, favorites.HttpAdd(t, token, id))
}

// BulkAdd merges multiple advertisement IDs into the authorized user's favorites.
func BulkAdd(t *testing.T, token string, ids []string) response.Result {
	t.Helper()
	return response.Read(t, favorites.HttpBulkAdd(t, token, ids))
}

// Delete removes an advertisement ID from the authorized user's favorites.
func Delete(t *testing.T, token, id string) response.Result {
	t.Helper()
	return response.Read(t, favorites.HttpDelete(t, token, id))
}
