package fixtures

import (
	"net/http"
	"testing"
	"time"

	"api-tests-template/internal/managers/response"
)

// RetryServiceUnavailable repeats a safe request until the gateway returns an API response or the deadline passes.
func RetryServiceUnavailable(t *testing.T, operation string, call func() response.Result) response.Result {
	t.Helper()
	deadline := time.Now().Add(retryTimeout)
	delay := retryInterval
	for {
		result := call()
		if result.StatusCode != http.StatusServiceUnavailable || time.Now().After(deadline) {
			return result
		}
		t.Logf("Временный HTTP 503 для %s; повторяем запрос через %s", operation, delay)
		time.Sleep(delay)
		delay = min(delay*2, time.Second)
	}
}
