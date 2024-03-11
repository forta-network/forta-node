package json_rpc

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Error   JsonRpcError    `json:"error"`
}

type JsonRpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func writeTooManyReqsErr(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusTooManyRequests)

	var reqPayload JsonRpcReq
	if err := json.NewDecoder(req.Body).Decode(&reqPayload); err != nil {
		log.WithError(err).Error("failed to decode jsonrpc request body")
		return
	}

	if err := json.NewEncoder(w).Encode(&ErrorResponse{
		JSONRPC: "2.0",
		ID:      reqPayload.ID,
		Error: JsonRpcError{
			Code:    -32000,
			Message: "agent exceeds scan node request limit",
		},
	}); err != nil {
		log.WithError(err).Error("failed to write jsonrpc error response body")
	}
}
