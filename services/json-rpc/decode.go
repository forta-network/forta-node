package json_rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type JsonRpcReq struct {
	ID     json.RawMessage `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type JsonRpcResp struct {
	ID     json.RawMessage `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  *JsonRpcError   `json:"error"`
}

func DecodeBody(req *http.Request) (*JsonRpcReq, error) {
	var decodedBody JsonRpcReq
	if err := json.NewDecoder(req.Body).Decode(&decodedBody); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to decode json-rpc request body")
	}
	return &decodedBody, nil
}

func decodeAndReplaceBody(req *http.Request) (*JsonRpcReq, error) {
	b, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body")
	}
	req.Body.Close()
	req.Body = io.NopCloser(bytes.NewBuffer(b))

	decodedBody, err := DecodeBody(req)
	if err != nil {
		return nil, fmt.Errorf("failed to decode json-rpc request body")
	}

	req.Body.Close()

	req.Body = io.NopCloser(bytes.NewBuffer(b))
	return decodedBody, nil
}
