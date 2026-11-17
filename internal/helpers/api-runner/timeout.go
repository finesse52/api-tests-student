package api_runner

import (
	"net/http"
	"time"

	"github.com/steinfletcher/apitest"
)

const (
	requestTimeout      = 20 * time.Second
	retryingCallTimeout = 35 * time.Second
	headerTimeout       = 8 * time.Second
)

// CreateWithTimeout prepares a request runner with a bounded network call duration.
func (at *ApiTest) CreateWithTimeout() *apitest.APITest {
	return at.Create().EnableNetworking(&http.Client{Timeout: requestTimeout})
}

// CreateWithTransientTransportRetry retries connection failures for requests known to have no successful side effect.
func (at *ApiTest) CreateWithTransientTransportRetry() *apitest.APITest {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.ResponseHeaderTimeout = headerTimeout
	transport.DisableKeepAlives = true
	return at.Create().EnableNetworking(&http.Client{
		Timeout:   retryingCallTimeout,
		Transport: &transientTransport{base: transport},
	})
}
