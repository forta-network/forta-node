package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/clients/ipfsclient"
	"github.com/forta-network/forta-node/clients/ipfsrouter"
	"github.com/forta-network/forta-node/config"
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
	ctx    context.Context
	ipfs   IPFSClient
	router IPFSRouter
	server *grpc.Server

	protocol.UnimplementedStorageServer
}

// New creates a new storage service.
func NewStorage(ctx context.Context, ipfsURL, routerURL string) (*Storage, error) {
	storage := &Storage{
		ctx:    ctx,
		ipfs:   ipfsclient.New(ipfsURL),
		router: ipfsrouter.NewClient(routerURL),
		server: grpc.NewServer(),
	}
	protocol.RegisterStorageServer(storage.server, storage)
	return storage, nil
}

// Start starts the service.
func (storage *Storage) Start() error {
	// just attempt creating the base dir to avoid unnecessary errors
	storage.ipfs.FilesMkdir(storage.ctx, DefaultBasePath, ipfsapi.FilesMkdir.Parents(true))

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", config.DefaultStoragePort))
	if err != nil {
		return err
	}
	go func() {
		log.Info("starting storage server...")
		err := storage.server.Serve(lis)
		log.WithError(err).Info("storage server stopped")
	}()

	go storage.loop(storage.ctx)

	return nil
}

func (storage *Storage) loop(ctx context.Context) {
	ticker := time.NewTicker(time.Minute * 1)
	var lastTick *time.Time
	for {
		select {
		case <-ctx.Done():
			log.Info("exiting content provider")
			return
		case currTick := <-ticker.C:
			if err := storage.doProvideContent(context.Background()); err != nil {
				log.WithError(err).Error("error while providing content")
			} else {
				log.Info("finished providing content refs")
			}

			if lastTick != nil && currTick.Sub(*lastTick) >= defaultGCInterval {
				if err := storage.doCollectGarbage(context.Background()); err != nil {
					log.WithError(err).Error("error while collecting garbage")
				} else {
					log.Info("finished collecting garbage")
				}
			}

			lastTick = &currTick
		}
	}
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

// Health implements the health.Reporter interface.
func (storage *Storage) Health() health.Reports {
	return nil
}

// Put puts given content to IPFS MFS.
func (storage *Storage) Put(ctx context.Context, req *protocol.PutRequest) (*protocol.PutResponse, error) {
	if len(req.User) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user")
	}
	if len(req.Kind) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "user not provided")
	}

	contentPath, bucketDir := NewContentPath(req.User, req.Kind)
	storage.ipfs.FilesMkdir(ctx, bucketDir, ipfsapi.FilesMkdir.Parents(true))
	contentID, err := storage.ipfs.AddToFiles(bytes.NewBuffer(req.Bytes), contentPath)
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
		return nil, status.Errorf(codes.NotFound, "content not found: %s", contentRef)
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

	contentDir := ContentDir(req.User, req.Kind)
	bucketEntries, err := storage.ipfs.FilesLs(ctx, contentDir, ipfsapi.FilesLs.Stat(true))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list the content directory '%s': %v", contentDir, err)
	}

	skipCount := int(req.Offset)
	limit := int(req.Limit)
	var allEntries []*ipfsapi.MfsLsEntry
	for _, bucketEntry := range bucketEntries {
		if limit == 0 {
			break
		}

		list, err := storage.getContentBucketEntries(ctx, req.User, req.Kind, bucketEntry.Name)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to list the bucket directory '%s' in '%s': %v", bucketEntry.Name, contentDir, err)
		}

		var added []*ipfsapi.MfsLsEntry

		// reduce skip count each time we skip content
		// if we completed or are about to complete the skip count, mark content for aggregation
		if skipCount >= len(list) {
			skipCount -= len(list)
		} else {
			added = list[skipCount:]
			skipCount = 0
		}

		// reduce total limit each time we aggregate content
		// if we are reaching the limit, use the remaining on the current list for the last time
		if len(added) <= limit {
			limit -= len(added)
		} else {
			added = added[:limit]
			limit = 0
		}

		// finally, append content
		allEntries = append(allEntries, added...)
	}

	if len(allEntries) >= int(req.Offset) {
		allEntries = allEntries[req.Offset:]
	}
	if len(allEntries) > int(req.Limit) {
		allEntries = allEntries[:req.Limit]
	}

	if req.Sort == protocol.SortDirection_DESC {
		sort.Slice(allEntries, func(i, j int) bool {
			return !sort.StringsAreSorted([]string{allEntries[i].Name, allEntries[j].Name})
		})
	}

	var resp protocol.ListResponse
	for _, entry := range allEntries {
		resp.Contents = append(resp.Contents, &protocol.ContentInfo{
			ContentId:   entry.Hash,
			ContentPath: path.Join(contentDir, entry.Name),
		})
	}

	return &resp, nil
}

// Provider returns the provider info.
func (storage *Storage) Provider(ctx context.Context, req *protocol.ProviderRequest) (*protocol.ProviderResponse, error) {
	idResp, err := storage.ipfs.ID()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get the id: %v", err)
	}
	return &protocol.ProviderResponse{
		Provider: &protocol.Provider{
			Id: idResp.ID,
		},
	}, nil
}

type userInfo struct {
	User         string
	ContentKinds []string
}

func (user *userInfo) HasContent(kind string) bool {
	for _, storedKind := range user.ContentKinds {
		if kind == storedKind {
			return true
		}
	}
	return false
}

func (storage *Storage) getUsers(ctx context.Context) ([]*userInfo, error) {
	list, err := storage.ipfs.FilesLs(ctx, DefaultBasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to list the base storage path: %v", err)
	}

	var users []*userInfo
	for _, stat := range list {
		userName := strings.Trim(stat.Name, "/")
		contentList, err := storage.ipfs.FilesLs(ctx, path.Join(DefaultBasePath, userName))
		if err != nil {
			return nil, fmt.Errorf("failed to get the content kinds for user '%s': %v", userName, err)
		}
		var contentKinds []string
		for _, kind := range contentList {
			contentKind := strings.Trim(kind.Name, "/")
			if contentKind == "bloom" {
				continue
			}
			contentKinds = append(contentKinds, contentKind)
		}
		log.WithField("user", userName).WithField("kinds", strings.Join(contentKinds, ",")).Debug("detected kinds")
		users = append(users, &userInfo{
			User:         userName,
			ContentKinds: contentKinds,
		})
	}

	return users, nil
}

func (storage *Storage) getContentBuckets(ctx context.Context, user, kind string) ([]*ipfsapi.MfsLsEntry, []*ipfsapi.MfsLsEntry, error) {
	contentDir := ContentDir(user, kind)
	list, err := storage.ipfs.FilesLs(ctx, contentDir, ipfsapi.FilesLs.Stat(true))
	if err != nil {
		return nil, nil, fmt.Errorf("error while listing '%s': %v", contentDir, err)
	}
	// ensure it's sorted in alphabetical order (ascending)
	sort.Slice(list, func(i, j int) bool {
		return sort.StringsAreSorted([]string{list[i].Name, list[j].Name})
	})
	limit := MaxBuckets

	var (
		newestEntries = list
		oldEntries    []*ipfsapi.MfsLsEntry
	)
	oldCount := len(list) - limit
	if oldCount > 0 {
		newestEntries = list[oldCount:]
		oldEntries = list[:oldCount]
	}
	return newestEntries, oldEntries, nil
}

func (storage *Storage) getContentBucketEntries(ctx context.Context, user, kind, bucket string) ([]*ipfsapi.MfsLsEntry, error) {
	bucketDir := BucketDir(user, kind, bucket)
	list, err := storage.ipfs.FilesLs(ctx, bucketDir, ipfsapi.FilesLs.Stat(true))
	if err != nil {
		return nil, fmt.Errorf("error while listing '%s': %v", bucketDir, err)
	}
	// ensure it's sorted in alphabetical order (ascending)
	sort.Slice(list, func(i, j int) bool {
		return sort.StringsAreSorted([]string{list[i].Name, list[j].Name})
	})
	return list, nil
}
