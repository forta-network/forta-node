package jwt_provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

type CreateJWTMessage struct {
	Claims map[string]interface{} `json:"claims"`
}
type CreateJWTResponse struct {
	Token string `json:"token"`
}

const (
	errBadCreateMessage  = "bad create jwt message body"
	errFailedToCreateJWT = "can't find bot id from request source"
)

// createJWTHandler returns a scanner jwt token with claims [hash] = hash(uri,payload) and [bot] = "bot id"
func (j *JWTProvider) createJWTHandler(w http.ResponseWriter, req *http.Request) {
	var msg CreateJWTMessage
	if req.Body != http.NoBody {
		err := json.NewDecoder(req.Body).Decode(&msg)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(w, errBadCreateMessage)
			return
		}
	}

	jwt, err := j.doCreateJWT(req.Context(), req.RemoteAddr, msg.Claims)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "can not create jwt: %v", err)
		return
	}

	resp, err := json.Marshal(CreateJWTResponse{Token: jwt})

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "%s", resp)
}

func (j *JWTProvider) doCreateJWT(ctx context.Context, remoteAddr string, claims map[string]interface{}) (string, error) {
	ipAddr, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return "", fmt.Errorf("can't extract ip from request %s", remoteAddr)
	}

	agentID, err := j.agentIDReverseLookup(ctx, ipAddr)
	if err != nil {
		return "", fmt.Errorf("can't find bot id from request source %s, err: %v", ipAddr, err)
	}

	jwt, err := CreateBotJWT(j.cfg.Key, agentID, claims)
	if err != nil {
		return "", fmt.Errorf(errFailedToCreateJWT)
	}

	return jwt, nil
}