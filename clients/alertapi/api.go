package alertapi

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/utils/httpclient"
	"github.com/goccy/go-json"
	log "github.com/sirupsen/logrus"
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
	resp, err := httpclient.Default.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
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
	return nil
}

func (c *client) PostBatch(batch *domain.AlertBatchRequest, token string) error {
	path := fmt.Sprintf("/batch/%s", batch.Ref)
	headers := map[string]string{
		"content-type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}
	return c.post(path, batch, headers)
}

func NewClient(apiUrl string) *client {
	return &client{apiUrl: apiUrl}
}
