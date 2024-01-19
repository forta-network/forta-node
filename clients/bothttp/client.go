package bothttp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/forta-network/forta-core-go/utils/httpclient"
	"github.com/hashicorp/go-multierror"
	log "github.com/sirupsen/logrus"
)

var (
	// responseSizeLimit is the maximum number of bytes to read from the response body.
	responseSizeLimit = int64(2 << 20) // 2MB
)

type HealthResponse struct {
	Errors  []string  `json:"errors"`
	Metrics []Metrics `json:"metrics"`
}

type Metrics struct {
	// ChainID is the id of the chain the metrics are for
	ChainID    int64                `json:"chainId"`
	DataPoints map[string][]float64 `json:"dataPoints"`
}

// Client is the bot HTTP client interface.
type Client interface {
	Health(ctx context.Context) ([]Metrics, error)
}

type botClient struct {
	baseUrl    string
	httpClient *http.Client
}

// NewClient creates a new client.
func NewClient(host string, port int) Client {
	return &botClient{
		baseUrl:    fmt.Sprintf("http://%s:%d", host, port),
		httpClient: httpclient.Default,
	}
}

// Health does a health check on the bot.
func (bc *botClient) Health(ctx context.Context) ([]Metrics, error) {
	healthUrl := fmt.Sprintf("%s/health", bc.baseUrl)
	req, err := http.NewRequestWithContext(ctx, "GET", healthUrl, nil)
	if err != nil {
		return nil, err
	}

	// TODO: circuit breaker for the response size
	resp, err := bc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var healthResp HealthResponse

	// Limit the response size to a certain number of bytes
	limitedReader := io.LimitReader(resp.Body, responseSizeLimit)
	if err := json.NewDecoder(limitedReader).Decode(&healthResp); err != nil {
		if strings.Contains(err.Error(), "EOF") {
			log.WithError(err).Warn("response size limit for health check is reached")
		}

		return nil, nil // ignore decoding errors
	}

	if len(healthResp.Errors) == 0 {
		return healthResp.Metrics, nil
	}

	for _, errMsg := range healthResp.Errors {
		err = multierror.Append(err, errors.New(errMsg))
	}

	return healthResp.Metrics, err
}
