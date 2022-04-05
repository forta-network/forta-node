package json_rpc

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type requestPayload struct {
	ID string `json:"id"`
}

type errorResponse struct {
	JSONRPC string       `json:"jsonrpc"`
	ID      string       `json:"id"`
	Error   jsonRpcError `json:"error"`
}

type jsonRpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func writeTooManyReqsErr(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusTooManyRequests)

	var reqPayload requestPayload
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
