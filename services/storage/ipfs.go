package storage

import (
	"context"
	"io"
	"net/http"

	ipfsapi "github.com/ipfs/go-ipfs-api"
)

// IPFSClient implements the required IPFS node interface.
type IPFSClient interface {
	Add(r io.Reader, options ...ipfsapi.AddOpts) (string, error)
	Cat(path string) (io.ReadCloser, error)

	FilesRead(ctx context.Context, path string, options ...ipfsapi.FilesOpt) (io.ReadCloser, error)
	FilesWrite(ctx context.Context, path string, data io.Reader, options ...ipfsapi.FilesOpt) error
	FilesRm(ctx context.Context, path string, force bool) error
	FilesCp(ctx context.Context, src string, dest string) error
	FilesStat(ctx context.Context, path string, options ...ipfsapi.FilesOpt) (*ipfsapi.FilesStatObject, error)
	FilesMkdir(ctx context.Context, path string, options ...ipfsapi.FilesOpt) error
	FilesLs(ctx context.Context, path string, options ...ipfsapi.FilesOpt) ([]*ipfsapi.MfsLsEntry, error)
	FilesMv(ctx context.Context, src string, dest string) error
}

// NewIPFSClient creates a new IPFS client.
func NewIPFSClient(apiURL string) IPFSClient {
	return ipfsapi.NewShellWithClient(apiURL, http.DefaultClient)
}

// toFiles allows adding the file to MFS directly.
func toFiles(path string) ipfsapi.AddOpts {
	return func(rb *ipfsapi.RequestBuilder) error {
		rb.Option("to-files", path)
		return nil
	}
}
