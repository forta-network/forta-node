package blocksdata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/andybalholm/brotli"
	backoff "github.com/cenkalti/backoff/v4"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-core-go/utils/httpclient"
	"google.golang.org/protobuf/proto"
)

var (
	minBackoff = 1 * time.Second
	maxBackoff = 1 * time.Minute
)

type blocksDataClient struct {
	dispatcherURL  *url.URL
	dispatcherPath string
}

func NewCombinedBlockEventsClient(dispatcherURL string) *blocksDataClient {
	u, _ := url.Parse(dispatcherURL)

	return &blocksDataClient{
		dispatcherURL:  u,
		dispatcherPath: u.Path,
	}
}

type PresignedURLItem struct {
	Bucket       int64  `json:"bucket"`
	PresignedURL string `json:"presignedURL"`
	ExpiresAt    int64  `json:"expiresAt"`
}

func (c *blocksDataClient) GetBlocksData(bucket int64) (_ *protocol.BlocksData, err error) {
	c.dispatcherURL.Path, err = url.JoinPath(c.dispatcherPath, fmt.Sprintf("%d", bucket))
	if err != nil {
		return nil, err
	}

	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = minBackoff
	bo.MaxInterval = maxBackoff

	var item PresignedURLItem

	err = backoff.Retry(func() error {
		resp, err := httpclient.Default.Get(c.dispatcherURL.String())
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&item)
		if err != nil {
			return err
		}

		if item.ExpiresAt < time.Now().Unix() {
			return backoff.Permanent(fmt.Errorf("presigned URL expired"))
		}
		return nil
	}, bo)

	if err != nil {
		return nil, err
	}

	var blocks protocol.BlocksData

	err = backoff.Retry(func() error {
		resp, err := httpclient.Default.Get(item.PresignedURL)
		if err != nil {
			return err
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		b, err := io.ReadAll(brotli.NewReader(resp.Body))
		if err != nil {
			return err
		}

		err = proto.Unmarshal(b, &blocks)
		if err != nil {
			return backoff.Permanent(err)
		}

		return nil
	}, bo)

	if err != nil {
		return nil, err
	}

	return &blocks, nil
}
