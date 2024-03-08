package json_rpc

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type errorResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Error   jsonRpcError    `json:"error"`
}

type jsonRpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func writeTooManyReqsErr(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusTooManyRequests)

	var reqPayload jsonRpcReq
	if err := json.NewDecoder(req.Body).Decode(&reqPayload); err != nil {
		log.WithError(err).Error("failed to decode jsonrpc request body")
		return
	}

	if err := json.NewEncoder(w).Encode(&errorResponse{
		JSONRPC: "2.0",
		ID:      reqPayload.ID,
		Error: jsonRpcError{
			Code:    -32000,
			Message: "agent exceeds scan node request limit",
		},
	}); err != nil {
		log.WithError(err).Error("failed to write jsonrpc error response body")
	}
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
