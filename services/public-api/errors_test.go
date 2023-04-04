package public_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type errorWriter struct{}

func (w *errorWriter) Header() http.Header {
	return http.Header{}
}

func (w *errorWriter) Write(b []byte) (int, error) {
	return 0, fmt.Errorf("write failed")
}

func (w *errorWriter) WriteHeader(statusCode int) {}

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

	// bad response writer
	r = require.New(t)

	req, err = http.NewRequest("POST", "http://asdf.asdf", nil)
	r.NoError(err)
	recorderBad := errorWriter{}

	writeTooManyReqsErr(&recorderBad, req)

	resp = recorder.Result()
	r.Equal(http.StatusTooManyRequests, resp.StatusCode)
	r.Error(json.NewDecoder(resp.Body).Decode(&errResp))
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

	// bad response writer
	r = require.New(t)

	req, err = http.NewRequest("POST", "http://asdf.asdf", nil)
	r.NoError(err)
	recorderBad := errorWriter{}

	writeAuthError(&recorderBad, req)

	resp = recorder.Result()
	r.Equal(http.StatusUnauthorized, resp.StatusCode)
	r.Error(json.NewDecoder(resp.Body).Decode(&errResp))
}
