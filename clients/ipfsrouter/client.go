package ipfsrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type routerClient struct {
	apiURL string
}

// NewClient creates a new client instance.
func NewClient(apiURL string) *routerClient {
	return &routerClient{apiURL: apiURL}
}

func (client *routerClient) Provide(ctx context.Context, scanner, peerID, bloomFilter string) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&struct {
		Scanner     string `json:"scanner"`
		PeerID      string `json:"peerId"`
		BloomFilter string `json:"bloomFilter"`
	}{
		Scanner:     scanner,
		PeerID:      peerID,
		BloomFilter: bloomFilter,
	}); err != nil {
		return fmt.Errorf("failed to encode provide request payload: %v", err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", client.apiURL, &buf)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("provide request failed: %v", err)
	}
	resp.Body.Close()
	return nil
}
