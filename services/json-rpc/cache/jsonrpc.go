package json_rpc_cache

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type jsonRpcReq struct {
	ID     json.RawMessage `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type jsonRpcResp struct {
	ID      json.RawMessage `json:"id"`
	JsonRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *jsonRpcError   `json:"error,omitempty"`
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

func writeJsonResponse(w http.ResponseWriter, req *jsonRpcReq, result any) error {
	b, err := json.Marshal(result)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)

	return json.NewEncoder(w).Encode(&jsonRpcResp{
		ID:      req.ID,
		JsonRPC: "2.0",
		Result:  b,
	})
}

func writeBadRequest(w http.ResponseWriter, req *jsonRpcReq, err error) {
	if req == nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusBadRequest)

	if err := json.NewEncoder(w).Encode(&errorResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Error: jsonRpcError{
			Code:    -32600,
			Message: err.Error(),
		},
	}); err != nil {
		log.WithError(err).Error("failed to write jsonrpc error response body")
	}
}

func writeUnauthorized(w http.ResponseWriter, req *jsonRpcReq) {
	w.WriteHeader(http.StatusUnauthorized)

	if err := json.NewEncoder(w).Encode(&errorResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Error: jsonRpcError{
			Code:    -32000,
			Message: "unauthorized",
		},
	}); err != nil {
		log.WithError(err).Error("failed to write jsonrpc error response body")
	}
}

func writeNotFound(w http.ResponseWriter, req *jsonRpcReq) {
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&jsonRpcResp{
		ID:      req.ID,
		JsonRPC: "2.0",
		Result:  nil,
		Error: &jsonRpcError{
			Code:    -32603,
			Message: "result not found in cache",
		},
	}); err != nil {
		log.WithError(err).Error("failed to write jsonrpc error response body")
	}
}

func writeInternalError(w http.ResponseWriter, req *jsonRpcReq, err error) {
	w.WriteHeader(http.StatusInternalServerError)

	if err := json.NewEncoder(w).Encode(&errorResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Error: jsonRpcError{
			Code:    -32603,
			Message: err.Error(),
		},
	}); err != nil {
		log.WithError(err).Error("failed to write jsonrpc error response body")
	}
}
