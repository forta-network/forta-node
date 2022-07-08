package bot_jwt

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateJWTMessage struct {
	Claims map[string]interface{} `json:"claims"`
}
type CreateJWTResponse struct {
	Token string `json:"token"`
}

// createJWTHandler returns a scanner jwt token with claims [hash] = hash(uri,payload) and [bot] = "bot id"
func (j *JWTProvider) createJWTHandler(w http.ResponseWriter, req *http.Request) {
	var msg CreateJWTMessage
	if req.Body != http.NoBody {
		err := json.NewDecoder(req.Body).Decode(&msg)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	
	agentID, err := j.agentIDReverseLookup(req.Context(), req.RemoteAddr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	jwt, err := CreateBotJWT(j.cfg.Key, agentID, msg.Claims)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	resp, err := json.Marshal(CreateJWTResponse{Token: jwt})
	
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "%s", resp)
}
