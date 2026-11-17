package response

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type Result struct {
	StatusCode int
	Body       []byte
}

type Error struct {
	Error   string          `json:"error"`
	Message string          `json:"message"`
	Details json.RawMessage `json:"details"`
}

// Read consumes and closes an HTTP response so callers can inspect status and body together.
func Read(t *testing.T, raw *http.Response) Result {
	t.Helper()
	require.NotNil(t, raw, "HTTP response is nil")
	defer raw.Body.Close()
	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "cannot read HTTP response")
	return Result{StatusCode: raw.StatusCode, Body: body}
}

// RequireStatus checks the status and includes the response body in failures.
func (r Result) RequireStatus(t *testing.T, expected int) {
	t.Helper()
	require.Equalf(t, expected, r.StatusCode, "unexpected HTTP status; body: %s", r.Body)
}

// Decode parses a successful JSON response into the supplied model.
func (r Result) Decode(t *testing.T, target any) {
	t.Helper()
	require.NoErrorf(t, json.Unmarshal(r.Body, target), "invalid JSON response: %s", r.Body)
}

// RequireError checks the HTTP status and the API's error code and message.
func (r Result) RequireError(t *testing.T, status int, errorCode string) Error {
	t.Helper()
	r.RequireStatus(t, status)
	var apiError Error
	r.Decode(t, &apiError)
	if errorCode == "" {
		require.NotEmpty(t, apiError.Error, "API error code is empty")
	} else {
		require.Equal(t, errorCode, apiError.Error, "unexpected error code")
	}
	require.NotEmpty(t, apiError.Message, "API error message is empty")
	return apiError
}
