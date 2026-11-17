package advertisement

import (
	"testing"

	"api-tests-template/internal/client/http/advertisement"
	"api-tests-template/internal/managers/advertisement/models"
	"api-tests-template/internal/managers/response"
)

// Create sends a new advertisement and returns its status and body for positive or negative checks.
func Create(t *testing.T, token string, request models.CreateRequest) response.Result {
	t.Helper()
	return response.Read(t, advertisement.HttpCreate(t, token, request))
}

// Get reads an advertisement by ID.
func Get(t *testing.T, id string) response.Result {
	t.Helper()
	return response.Read(t, advertisement.HttpGet(t, id))
}

// Delete removes an advertisement by ID using its owner's token.
func Delete(t *testing.T, token, id string) response.Result {
	t.Helper()
	return response.Read(t, advertisement.HttpDelete(t, token, id))
}

// GetPhotos reads the list of stored photos for an advertisement.
func GetPhotos(t *testing.T, token, id string) response.Result {
	t.Helper()
	return response.Read(t, advertisement.HttpGetPhotos(t, token, id))
}
