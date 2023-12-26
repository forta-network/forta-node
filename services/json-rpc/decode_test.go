package json_rpc

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeAndReplaceBody(t *testing.T) {
	r := require.New(t)

	bodyStr := `{"method":"eth_call"}`
	req, err := http.NewRequest("GET", "/", bytes.NewBufferString(bodyStr))
	r.NoError(err)

	// decodes successfully
	decodedBody, err := decodeAndReplaceBody(req)
	r.NoError(err)
	r.Equal("eth_call", decodedBody.Method)

	// still can read body because it was replaced
	b, err := io.ReadAll(req.Body)
	r.Equal(bodyStr, string(b))
}
