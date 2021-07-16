package registry

import (
	"fmt"
	"io"
	"net/http"
)

type ipfsClient struct {
	gatewayURL string
}

func (client *ipfsClient) Get(cid string) (io.ReadCloser, error) {
	resp, err := http.Get(fmt.Sprintf("%s/ipfs/%s", client.gatewayURL, cid))
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
