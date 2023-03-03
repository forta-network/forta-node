package ipfsrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/forta-network/forta-core-go/utils/httpclient"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/sirupsen/logrus"
)

type routerClient struct {
	apiURL string
}

// NewClient creates a new client instance.
func NewClient(apiURL string) *routerClient {
	return &routerClient{apiURL: apiURL}
}

func (client *routerClient) Provide(ctx context.Context, scanner, peerID, bloomFilter string) error {
	pID, err := peer.Decode(peerID)
	if err != nil {
		return fmt.Errorf("failed to decode peer ID: %v", err)
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&struct {
		Scanner     string        `json:"scanner"`
		Updated     time.Time     `json:"updated"`
		Peer        peer.AddrInfo `json:"peer"`
		BloomFilter string        `json:"bloomFilter"`
	}{
		Scanner: scanner,
		Updated: time.Now().UTC(),
		Peer: peer.AddrInfo{
			ID: pID,
		},
		BloomFilter: bloomFilter,
	}); err != nil {
		return fmt.Errorf("failed to encode provide request payload: %v", err)
	}
	logrus.WithField("payload", string(buf.Bytes())).Trace("encoded provide request")

	req, err := http.NewRequestWithContext(ctx, "POST", client.apiURL, &buf)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	resp, err := httpclient.Default.Do(req)
	if err != nil {
		return fmt.Errorf("provide request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with '%d': %s", resp.StatusCode, string(b))
	}

	return nil
}
