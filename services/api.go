package services

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ErrResponse struct {
	Message string `json:"message"`
}

func writeError(w http.ResponseWriter, status int, e error) {
	w.WriteHeader(status)
	er := ErrResponse{Message: e.Error()}
	b, err := json.Marshal(er)
	if err != nil {
		log.Errorf("could not marshal error json: %s", e.Error())
		return
	}
	if _, err := w.Write(b); err != nil {
		log.Errorf("could not write err resp")
		return
	}
}
