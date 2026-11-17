package api_runner

import (
	"net/http"
	"time"

	"github.com/steinfletcher/apitest"
)

const requestTimeout = 20 * time.Second

// CreateWithTimeout prepares a request runner with a bounded network call duration.
func (at *ApiTest) CreateWithTimeout() *apitest.APITest {
	return at.Create().EnableNetworking(&http.Client{Timeout: requestTimeout})
}
