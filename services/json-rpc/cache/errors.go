package json_rpc_cache

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

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
