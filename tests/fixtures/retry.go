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
	attempt := 1
	for {
		result := call()
		if result.StatusCode != http.StatusServiceUnavailable || time.Now().After(deadline) {
			return result
		}
		t.Logf("Попытка %d для %s вернула HTTP 503; попытка %d через %s", attempt, operation, attempt+1, delay)
		time.Sleep(delay)
		delay = min(delay*2, time.Second)
		attempt++
	}
}
