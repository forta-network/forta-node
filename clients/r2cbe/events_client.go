package r2cbe

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/forta-network/forta-core-go/protocol"
	"google.golang.org/protobuf/proto"
)

type combinedBlockEventsClient struct {
	dispatcherClient *http.Client
	r2Client         *http.Client

	dispatcherURL  *url.URL
	dispatcherPath string
}

func NewCombinedBlockEventsClient(dispatcherURL string) *combinedBlockEventsClient {
	u, _ := url.Parse(dispatcherURL)

	dispatcherClient := http.DefaultClient
	dispatcherClient.Timeout = 10 * time.Second

	return &combinedBlockEventsClient{
		dispatcherClient: dispatcherClient,
		r2Client:         http.DefaultClient,
		dispatcherURL:    u,
		dispatcherPath:   u.Path,
	}
}

type PresignedURLItem struct {
	Bucket       int64  `json:"bucket"`
	PresignedURL string `json:"presignedURL"`
	ExpiresAt    int64  `json:"expiresAt"`
}

func (c *combinedBlockEventsClient) GetCombinedBlockEvents(bucket int64) (_ *protocol.CombinedBlockEvents, err error) {
	c.dispatcherURL.Path, err = url.JoinPath(c.dispatcherPath, fmt.Sprintf("%d", bucket))
	if err != nil {
		return nil, err
	}

	resp, err := c.dispatcherClient.Get(c.dispatcherURL.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var item PresignedURLItem
	err = json.NewDecoder(resp.Body).Decode(&item)
	if err != nil {
		return nil, err
	}

	resp, err = c.r2Client.Get(item.PresignedURL)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(brotli.NewReader(resp.Body))
	if err != nil {
		return nil, err
	}

	var events protocol.CombinedBlockEvents

	err = proto.Unmarshal(b, &events)
	if err != nil {
		return nil, err
	}

	return &events, nil
}
