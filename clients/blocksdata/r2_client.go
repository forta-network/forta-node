package blocksdata

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	backoff "github.com/cenkalti/backoff/v4"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-core-go/utils/httpclient"
	"google.golang.org/protobuf/proto"
)

const (
	minBackoff     = 1 * time.Second
	maxBackoff     = 10 * time.Second
	maxElapsedTime = 5 * time.Minute
)

type blocksDataClient struct {
	dispatcherURL *url.URL
}

func NewBlocksDataClient(dispatcherURL string) *blocksDataClient {
	u, _ := url.Parse(dispatcherURL)

	return &blocksDataClient{
		dispatcherURL: u,
	}
}

type PresignedURLItem struct {
	Bucket       int64  `json:"bucket"`
	PresignedURL string `json:"presignedURL"`
	ExpiresAt    int64  `json:"expiresAt"`
}

func (c *blocksDataClient) GetBlocksData(bucket int64) (_ *protocol.BlocksData, err error) {
	dispatcherUrl, err := url.JoinPath(c.dispatcherURL.String(), fmt.Sprintf("%d", bucket))
	if err != nil {
		return nil, err
	}

	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = minBackoff
	bo.MaxInterval = maxBackoff
	bo.MaxElapsedTime = maxElapsedTime

	var item PresignedURLItem

	err = backoff.Retry(func() error {
		resp, err := httpclient.Default.Get(dispatcherUrl)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if resp.StatusCode == http.StatusForbidden {
			return backoff.Permanent(fmt.Errorf("forbidden"))
		}

		if resp.StatusCode == http.StatusNotFound && bytes.Contains(b, []byte("too old")) {
			return fmt.Errorf("%s", b)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, b)
		}

		err = json.Unmarshal(b, &item)
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

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return err
		}

		b, err := io.ReadAll(gzipReader)
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
