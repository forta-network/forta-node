package bothttp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/forta-network/forta-core-go/utils/httpclient"
	"github.com/hashicorp/go-multierror"
)

type HealthResponse struct {
	Errors []string `json:"errors"`
}

// Client is the bot HTTP client interface.
type Client interface {
	Health(ctx context.Context) error
}

type botClient struct {
	baseUrl    string
	httpClient *http.Client
}

// NewClient creates anew client.
func NewClient(host string, port int) Client {
	return &botClient{
		baseUrl:    fmt.Sprintf("http://%s:%d", host, port),
		httpClient: httpclient.Default,
	}
}

// Health does a health check on the bot.
func (bc *botClient) Health(ctx context.Context) error {
	healthUrl := fmt.Sprintf("%s/health", bc.baseUrl)
	req, err := http.NewRequestWithContext(ctx, "GET", healthUrl, nil)
	if err != nil {
		return err
	}
	resp, err := bc.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var healthResp HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&healthResp); err != nil {
		return err
	}
	if len(healthResp.Errors) == 0 {
		return nil
	}
	for _, errMsg := range healthResp.Errors {
		err = multierror.Append(err, errors.New(errMsg))
	}
	return err
}
