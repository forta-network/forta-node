package jwt_provider

import (
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

var (
	resolver = &net.Resolver{
		PreferGo:     true,
		StrictErrors: false,
		Dial:         nil,
	}
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

	ipAddr, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "can't extract ip from request %s", req.RemoteAddr)
		return
	}

	agentID, err := j.agentIDReverseLookup(req.Context(), ipAddr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "can't find bot id from request source %s, err: %v", ipAddr, err)
		return
	}

	jwt, err := CreateBotJWT(j.cfg.Key, agentID, msg.Claims)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, errFailedToCreateJWT)
		return
	}

	resp, err := json.Marshal(CreateJWTResponse{Token: jwt})

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "%s", resp)
}
