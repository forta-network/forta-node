package bot_jwt

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type createJWTMessage struct {
	Hash string `json:"hash"`
	Exp  uint64 `json:"exp,omitempty"`
}

// createJWTHandler returns a scanner jwt token with claims [hash] = hash(uri,payload) and [bot] = "bot id"
func (j *JWTProvider) createJWTHandler(w http.ResponseWriter, req *http.Request) {
	var msg createJWTMessage
	err := json.NewDecoder(req.Body).Decode(&msg)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	agentID, err := j.agentIDReverseLookup(req.Context(), req.RemoteAddr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	jwt, err := CreateBotJWT(j.key, agentID, msg.Hash, msg.Exp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, jwt)
}
