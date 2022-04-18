package alertapi

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/forta-network/forta-core-go/domain"
	"github.com/goccy/go-json"
	log "github.com/sirupsen/logrus"
)

type client struct {
	apiUrl string
}

func (c *client) post(path string, body interface{}, headers map[string]string, target interface{}) error {
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
	b, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.WithFields(log.Fields{
			"apiUrl":   c.apiUrl,
			"path":     path,
			"body":     string(jsonVal),
			"response": string(b),
			"status":   resp.StatusCode,
		}).Error("alert api error")
		return fmt.Errorf("%d error: %s", resp.StatusCode, string(b))
	}
	return json.Unmarshal(b, target)
}

func (c *client) PostBatch(batch *domain.AlertBatchRequest, token string) (*domain.AlertBatchResponse, error) {
	path := fmt.Sprintf("/batch/%s", batch.Ref)
	headers := map[string]string{
		"content-type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}
	var resp domain.AlertBatchResponse
	if err := c.post(path, batch, headers, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func NewClient(apiUrl string) *client {
	return &client{apiUrl: apiUrl}
}
