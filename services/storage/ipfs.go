package storage

import (
	"context"
	"io"

	ipfsapi "github.com/ipfs/go-ipfs-api"
)

// IPFSClient implements the required IPFS node interface.
type IPFSClient interface {
	AddToFiles(r io.Reader, path string, options ...ipfsapi.AddOpts) (string, error)
	Cat(path string) (io.ReadCloser, error)
	Unpin(path string) error

	FilesRead(ctx context.Context, path string, options ...ipfsapi.FilesOpt) (io.ReadCloser, error)
	FilesWrite(ctx context.Context, path string, data io.Reader, options ...ipfsapi.FilesOpt) error
	FilesRm(ctx context.Context, path string, force bool) error
	FilesCp(ctx context.Context, src string, dest string) error
	FilesStat(ctx context.Context, path string, options ...ipfsapi.FilesOpt) (*ipfsapi.FilesStatObject, error)
	FilesMkdir(ctx context.Context, path string, options ...ipfsapi.FilesOpt) error
	FilesLs(ctx context.Context, path string, options ...ipfsapi.FilesOpt) ([]*ipfsapi.MfsLsEntry, error)
	FilesMv(ctx context.Context, src string, dest string) error

	RepoGC(context.Context) error

	ID(peer ...string) (*ipfsapi.IdOutput, error)
}

// IPFSRouter implements the IPFS router server methods.
type IPFSRouter interface {
	Provide(ctx context.Context, scanner, peerID, bloomFilter string) error
}
