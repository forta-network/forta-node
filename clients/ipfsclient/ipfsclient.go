package ipfsclient

import (
	"context"
	"io"
	"net/http"

	ipfsapi "github.com/ipfs/go-ipfs-api"
)

// Client wraps the IPFS client by implementing some extra methods.
type Client struct {
	*ipfsapi.Shell
}

// New creates a new client.
func New(ipfsURL string) *Client {
	return &Client{Shell: ipfsapi.NewShellWithClient(ipfsURL, http.DefaultClient)}
}

// RepoGC triggers garbage collection on the repo.
func (client *Client) RepoGC(ctx context.Context) error {
	return client.Request("repo/gc").Exec(ctx, nil)
}

// AddToFiles uses the `to-files` option to add to the given path. If no path is provided,
// then it works like the default add.
func (client *Client) AddToFiles(r io.Reader, path string, options ...ipfsapi.AddOpts) (string, error) {
	if len(path) > 0 {
		return client.Add(r, append(options, toFiles(path))...)
	}
	return client.Add(r, options...)
}

// toFiles allows adding the file to MFS directly.
func toFiles(path string) ipfsapi.AddOpts {
	return func(rb *ipfsapi.RequestBuilder) error {
		rb.Option("to-files", path)
		return nil
	}
}
