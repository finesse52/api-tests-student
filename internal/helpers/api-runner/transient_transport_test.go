package api_runner

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

// RoundTrip invokes the test's transport behavior for one request attempt.
func (call roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return call(request)
}

func TestTransientTransportReplaysMultipartBodyAfterUnexpectedEOF(t *testing.T) {
	const payload = "multipart request body"
	attempts := 0
	transport := transientTransport{base: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		attempts++
		body, err := io.ReadAll(request.Body)
		require.NoError(t, err)
		require.NoError(t, request.Body.Close())
		require.Equal(t, payload, string(body), "multipart body changed during retry")
		if attempts == 1 {
			return nil, io.ErrUnexpectedEOF
		}
		return &http.Response{StatusCode: http.StatusBadRequest, Body: io.NopCloser(bytes.NewReader(nil))}, nil
	})}
	request, err := http.NewRequest(http.MethodPost, "https://example.com/advertisement", bytes.NewBufferString(payload))
	require.NoError(t, err)

	response, err := transport.RoundTrip(request)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, response.StatusCode)
	require.Equal(t, 2, attempts, "transport should retry once after unexpected EOF")
	require.NoError(t, response.Body.Close())
}
