package json_rpc_cache

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type jsonRpcReq struct {
	ID     json.RawMessage `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type jsonRpcResp struct {
	ID     json.RawMessage `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  *jsonRpcError   `json:"error"`
}

type errorResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Error   jsonRpcError    `json:"error"`
}

type jsonRpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func decodeBody(req *http.Request) (*jsonRpcReq, error) {
	var decodedBody jsonRpcReq
	if err := json.NewDecoder(req.Body).Decode(&decodedBody); err != nil {
		return nil, fmt.Errorf("failed to decode json-rpc request body")
	}
	return &decodedBody, nil
}
