package json_rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

const testRequestID = "123"

func TestTooManyReqsError(t *testing.T) {
	r := require.New(t)

	buf := bytes.NewBuffer([]byte(fmt.Sprintf(`{"id":"%s"}`, testRequestID)))
	req, err := http.NewRequest("POST", "http://asdf.asdf", buf)
	r.NoError(err)
	recorder := httptest.NewRecorder()

	writeTooManyReqsErr(recorder, req)

	resp := recorder.Result()
	r.Equal(http.StatusTooManyRequests, resp.StatusCode)
	var errResp errorResponse
	r.NoError(json.NewDecoder(resp.Body).Decode(&errResp))
	r.Equal("2.0", errResp.JSONRPC)
	r.Equal(testRequestID, errResp.ID)
	r.Equal(-32000, errResp.Error.Code)
	r.Contains(errResp.Error.Message, "exceeds")
}
