package json_rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type jsonRpcReq struct {
	ID     json.RawMessage `json:"id"`
	Method string          `json:"method"`
	Params string          `json:"params"`
}

type jsonRpcResp struct {
	ID     json.RawMessage `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  jsonRpcError    `json:"error"`
}

func decodeBody(req *http.Request) (*jsonRpcReq, error) {
	var decodedBody jsonRpcReq
	if err := json.NewDecoder(req.Body).Decode(&decodedBody); err != nil {
		return nil, fmt.Errorf("failed to decode json-rpc request body")
	}
	return &decodedBody, nil
}

func decodeAndReplaceBody(req *http.Request) (*jsonRpcReq, error) {
	b, _ := io.ReadAll(req.Body)
	req.Body.Close()
	req.Body = io.NopCloser(bytes.NewBuffer(b))
	return decodeBody(req)
}
