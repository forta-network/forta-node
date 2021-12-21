package store

import (
	"io"
	"net/http"

	"github.com/goccy/go-json"
)

//TODO: this is 100% throwaway code until the contract exists
type versions struct {
	Scanner string `json:"scanner"`
}

//WebUpdaterStore gets a version from an endpoint
type WebUpdaterStore struct {
	url string
}

//NewWebUpdaterStore creates a WebUpdaterStore
func NewWebUpdaterStore(url string) *WebUpdaterStore {
	return &WebUpdaterStore{url: url}
}

type autotaskResponse struct {
	Result string `json:"result"`
}

// GetLatestReference parses the response from an autotask
func (ws *WebUpdaterStore) GetLatestReference() (string, error) {
	res, err := http.Post(ws.url, "application/json", nil)
	if err != nil {
		return "", err
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var ar autotaskResponse
	if err := json.Unmarshal(b, &ar); err != nil {
		return "", err
	}
	var v versions
	if err := json.Unmarshal([]byte(ar.Result), &v); err != nil {
		return "", err
	}
	return v.Scanner, nil
}
