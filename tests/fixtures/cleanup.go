package fixtures

import (
	"net/http"
	"testing"
	"time"

	"api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/response"
)

const (
	retryTimeout  = 10 * time.Second
	retryInterval = 250 * time.Millisecond
)

// RequireAdvertisementDeleted polls until a deleted advertisement is no longer readable.
func RequireAdvertisementDeleted(t *testing.T, id string) {
	t.Helper()
	deadline := time.Now().Add(retryTimeout)
	for {
		result := advertisement.Get(t, id)
		if result.StatusCode == http.StatusNotFound {
			return
		}
		if result.StatusCode != http.StatusOK && result.StatusCode != http.StatusServiceUnavailable {
			result.RequireStatus(t, http.StatusNotFound)
		}
		if time.Now().After(deadline) {
			t.Fatalf("объявление %s осталось доступным после удаления: HTTP %d, %s", id, result.StatusCode, result.Body)
		}
		time.Sleep(retryInterval)
	}
}

// DeleteUntilCompleted retries an idempotent DELETE after a transient gateway error.
func DeleteUntilCompleted(t *testing.T, entity string, remove func() response.Result) {
	t.Helper()
	deadline := time.Now().Add(retryTimeout)
	sawGatewayError := false
	for {
		result := remove()
		if result.StatusCode == http.StatusNoContent || (sawGatewayError && result.StatusCode == http.StatusNotFound) {
			t.Logf("Результат удаления %s: HTTP %d", entity, result.StatusCode)
			return
		}
		if result.StatusCode != http.StatusServiceUnavailable {
			result.RequireStatus(t, http.StatusNoContent)
		}
		if time.Now().After(deadline) {
			t.Fatalf("удаление %s не завершилось за %s: HTTP %d, %s", entity, retryTimeout, result.StatusCode, result.Body)
		}
		sawGatewayError = true
		t.Logf("Временный HTTP 503 при удалении %s; повторяем запрос", entity)
		time.Sleep(retryInterval)
	}
}

// cleanupDeletedEntity retries transient gateway errors while deleting test-owned data.
func cleanupDeletedEntity(t *testing.T, entity string, remove func() response.Result) {
	t.Helper()
	deadline := time.Now().Add(retryTimeout)
	for {
		result := remove()
		if result.StatusCode == http.StatusNoContent || result.StatusCode == http.StatusNotFound {
			return
		}
		if result.StatusCode != http.StatusServiceUnavailable || time.Now().After(deadline) {
			t.Errorf("не удалось очистить %s: HTTP %d, %s", entity, result.StatusCode, result.Body)
			return
		}
		time.Sleep(retryInterval)
	}
}
