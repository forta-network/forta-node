package alertapi

import (
	"bytes"
	"fmt"
	"github.com/forta-protocol/forta-node/domain"
	"github.com/goccy/go-json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

type client struct {
	apiUrl string
}

func (c *client) post(path string, body interface{}, headers map[string]string) error {
	jsonVal, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.apiUrl, path), bytes.NewBuffer(jsonVal))
	if err != nil {
		return err
	}
	for n, v := range headers {
		req.Header[n] = []string{v}
	}
	hClient := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := hClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		log.WithFields(log.Fields{
			"apiUrl":   c.apiUrl,
			"path":     path,
			"body":     string(jsonVal),
			"response": string(b),
			"status":   resp.StatusCode,
		}).Error("alert api error")
	}
	return nil
}

func (c *client) PostBatch(batch *domain.AlertBatch) error {
	path := fmt.Sprintf("/batch/%s", batch.Ref)
	headers := map[string]string{
		"x-forta-scanner": batch.Scanner,
		"content-type":    "application/json",
	}
	return c.post(path, batch, headers)
}

func NewClient(apiUrl string) *client {
	return &client{apiUrl: apiUrl}
}
