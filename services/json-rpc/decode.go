package json_rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type jsonRpcReq struct {
	Method string `json:"method"`
}

func decodeAndReplaceBody(req *http.Request) (*jsonRpcReq, error) {
	b, _ := io.ReadAll(req.Body)
	req.Body.Close()
	req.Body = io.NopCloser(bytes.NewBuffer(b))
	var decodedBody jsonRpcReq
	if err := json.Unmarshal(b, &decodedBody); err != nil {
		return nil, fmt.Errorf("failed to decode json-rpc request body")
	}
	return &decodedBody, nil
}
