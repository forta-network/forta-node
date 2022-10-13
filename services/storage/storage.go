package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"path"
	"sort"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/services/storage/content"
	"github.com/ipfs/go-cid"
	ipfsapi "github.com/ipfs/go-ipfs-api"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const defaultListLimit = 1000

// Storage persists node content.
type Storage struct {
	ipfs   IPFSClient
	server *grpc.Server

	protocol.UnimplementedStorageServer
}

// New creates a new storage service.
func NewStorage(ipfsURL string) (*Storage, error) {
	storage := &Storage{
		ipfs:   ipfsapi.NewShellWithClient(ipfsURL, http.DefaultClient),
		server: grpc.NewServer(),
	}
	protocol.RegisterStorageServer(storage.server, storage)
	return storage, nil
}

// Start starts the service.
func (storage *Storage) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", config.DefaultStoragePort))
	if err != nil {
		return err
	}
	go func() {
		log.Info("starting storage server...")
		err := storage.server.Serve(lis)
		log.WithError(err).Info("storage server stopped")
	}()

	return nil
}

// Stop stops the service.
func (storage *Storage) Stop() error {
	storage.server.GracefulStop()
	return nil
}

// Name returns the name of the service.
func (storage *Storage) Name() string {
	return "storage"
}

// Put puts given content to IPFS MFS.
func (storage *Storage) Put(ctx context.Context, req *protocol.PutRequest) (*protocol.PutResponse, error) {
	if len(req.User) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user")
	}
	if len(req.Kind) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "user not provided")
	}

	// TODO: Remove or enable after testing.
	// contentDir := content.ContentDir(req.User, req.Kind)
	// if err := storage.ipfs.FilesMkdir(ctx, contentDir, ipfsapi.FilesMkdir.Parents(true)); err != nil {
	// 	log.WithError(err).Error("failed to create parent directories")
	// 	return nil, err
	// }

	contentPath := content.NewContentPath(req.User, req.Kind)
	contentID, err := storage.ipfs.Add(bytes.NewBuffer(req.Bytes), toFiles(contentPath))
	if err != nil {
		return nil, err
	}

	return &protocol.PutResponse{
		ContentId:   contentID,
		ContentPath: contentPath,
	}, nil
}

// Get get gets requested content either by using the content ID or MFS path.
func (storage *Storage) Get(ctx context.Context, req *protocol.GetRequest) (*protocol.GetResponse, error) {
	if len(req.ContentId) == 0 && len(req.ContentPath) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "please provide one of contentId or contentPath")
	}
	if req.Download && len(req.ContentId) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "need the contentId for download")
	}

	var contentRef string
	if len(req.ContentId) > 0 {
		_, err := cid.Parse(req.ContentId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid contentId: %s", req.ContentId)
		}
		contentRef = path.Join("/ipfs", req.ContentId)
	} else {
		contentRef = req.ContentPath
	}

	// if we should not download with content id, check if content exists before requesting it
	// so we skip content resolution.
	// if we already cannot download using content path, check the content anyways so we return
	// a meaningful not found error.
	var err error
	if !req.Download {
		_, err = storage.ipfs.FilesStat(ctx, contentRef)
	}
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "content not found: %s")
	}

	var r io.ReadCloser
	if len(req.ContentId) > 0 {
		r, err = storage.ipfs.Cat(contentRef)
	} else {
		r, err = storage.ipfs.FilesRead(ctx, contentRef)
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get file bytes: %v", err)
	}
	defer r.Close()

	b, err := io.ReadAll(r)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to read bytes: %v", err)
	}

	return &protocol.GetResponse{Bytes: b}, nil
}

// List returns the list of entries.
func (storage *Storage) List(ctx context.Context, req *protocol.ListRequest) (*protocol.ListResponse, error) {
	if req.Offset < 0 {
		req.Offset = 0
	}
	if req.Limit <= 0 || req.Limit > defaultListLimit {
		req.Limit = defaultListLimit
	}

	contentDir := content.ContentDir(req.User, req.Kind)
	list, err := storage.ipfs.FilesLs(ctx, contentDir, ipfsapi.FilesLs.Stat(true))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list the directory '%s': %v", contentDir, err)
	}

	// TODO: Check these again.
	if len(list) >= int(req.Offset) {
		list = list[req.Offset:]
	}
	if len(list) > int(req.Limit) {
		list = list[:req.Limit]
	}

	if req.Sort == protocol.SortDirection_DESC {
		sort.Slice(list, func(i, j int) bool {
			return !sort.StringsAreSorted([]string{list[i].Name, list[j].Name})
		})
	}

	var resp protocol.ListResponse
	for _, entry := range list {
		resp.Contents = append(resp.Contents, &protocol.ContentInfo{
			ContentId:   entry.Hash,
			ContentPath: path.Join(contentDir, entry.Name),
		})
	}

	return &resp, nil
}
