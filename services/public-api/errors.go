package public_api

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type errorResponse struct {
	Error publicAPIProxyError `json:"error"`
}

type publicAPIProxyError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func writeAuthError(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)


	if err := json.NewEncoder(w).Encode(&errorResponse{
		Error: publicAPIProxyError{
			Code:    -33000,
			Message: "request source is not a deployed agent",
		},
	}); err != nil {
		log.WithError(err).Error("failed to write jsonrpc error response body")
	}
}
func writeTooManyReqsErr(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusTooManyRequests)

	if err := json.NewEncoder(w).Encode(&errorResponse{
		Error: publicAPIProxyError{
			Code:    -32000,
			Message: "bot exceeds request rate limit",
		},
	}); err != nil {
		log.WithError(err).Error("failed to write jsonrpc error response body")
	}
}
