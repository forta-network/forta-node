package public_api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTooManyReqsError(t *testing.T) {
	r := require.New(t)

	req, err := http.NewRequest("POST", "http://asdf.asdf", nil)
	r.NoError(err)
	recorder := httptest.NewRecorder()

	writeTooManyReqsErr(recorder, req)

	resp := recorder.Result()
	r.Equal(http.StatusTooManyRequests, resp.StatusCode)
	var errResp errorResponse
	r.NoError(json.NewDecoder(resp.Body).Decode(&errResp))
	r.Contains(errResp.Error.Message, "exceeds")
}

func TestAuthError(t *testing.T) {
	r := require.New(t)

	req, err := http.NewRequest("POST", "http://asdf.asdf", nil)
	r.NoError(err)
	recorder := httptest.NewRecorder()

	writeAuthError(recorder, req)

	resp := recorder.Result()
	r.Equal(http.StatusUnauthorized, resp.StatusCode)
	var errResp errorResponse
	r.NoError(json.NewDecoder(resp.Body).Decode(&errResp))
	r.Contains(errResp.Error.Message, "request source is not a deployed agent")
}
