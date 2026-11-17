package api_runner

import (
	"bytes"
	"errors"
	"io"
	"net"
	"net/http"
	"syscall"
	"time"
)

const (
	transportAttempts  = 3
	transportRetryWait = 250 * time.Millisecond
)

type transientTransport struct {
	base http.RoundTripper
}

// RoundTrip retries transient transport failures and rebuilds the request body for each attempt.
func (transport *transientTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	var body []byte
	if request.Body != nil {
		var err error
		body, err = io.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
		if err = request.Body.Close(); err != nil {
			return nil, err
		}
	}

	for attempt := range transportAttempts {
		current := request.Clone(request.Context())
		if request.Body != nil {
			current.Body = io.NopCloser(bytes.NewReader(body))
		}
		response, err := transport.base.RoundTrip(current)
		if err == nil || !isTransientTransportError(err) || attempt == transportAttempts-1 {
			return response, err
		}
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
		wait := time.NewTimer(transportRetryWait * time.Duration(attempt+1))
		select {
		case <-request.Context().Done():
			wait.Stop()
			return nil, request.Context().Err()
		case <-wait.C:
		}
	}
	return nil, errors.New("transport retry exhausted")
}

// isTransientTransportError identifies network failures that can clear on a new connection.
func isTransientTransportError(err error) bool {
	var networkError net.Error
	return errors.As(err, &networkError) && networkError.Timeout() ||
		errors.Is(err, io.EOF) ||
		errors.Is(err, io.ErrUnexpectedEOF) ||
		errors.Is(err, syscall.ECONNRESET) ||
		errors.Is(err, syscall.EPIPE)
}
