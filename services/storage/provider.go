package storage

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"time"

	ipfsapi "github.com/ipfs/go-ipfs-api"
	log "github.com/sirupsen/logrus"
	boom "github.com/tylertreat/BoomFilters"
)

const (
	provideContentInterval = time.Minute
)

func (storage *Storage) doProvideContent(ctx context.Context) error {
	users, err := storage.getUsers(ctx)
	if err != nil {
		return err
	}

	for _, user := range users {
		if err := storage.prepareAndSendBloom(ctx, user); err != nil {
			log.WithField("user", user.User).WithError(err).Error("failed to provide")
		}
	}

	return nil
}

func (storage *Storage) prepareAndSendBloom(ctx context.Context, user *userInfo) error {
	logger := log.WithField("user", user.User)

	var allHashes []string
	// prioritize some content over other
	for _, kind := range []string{
		KindBatchReceipt,
	} {
		if !user.HasContent(kind) {
			continue
		}

		bucketEntries, _, err := storage.getContentBuckets(ctx, user.User, kind, true)
		if err != nil {
			return err
		}

		lastIndex := len(bucketEntries) - 1
		for i, bucketEntry := range bucketEntries {
			useCache := i < lastIndex
			entries, err := storage.getContentBucketEntries(ctx, user.User, kind, bucketEntry.Name, true, useCache)
			if err != nil {
				return fmt.Errorf("failed to list bucket entries: %v", err)
			}
			for _, entry := range entries {
				allHashes = append(allHashes, entry.Hash)
			}
			if len(allHashes) > BloomLimit {
				allHashes = allHashes[:BloomLimit]
				break
			}
		}

		if len(allHashes) > BloomLimit {
			allHashes = allHashes[:BloomLimit]
			break
		}
	}
	if len(allHashes) == 0 {
		logger.Info("no entries found - skipping provide call")
		return nil
	}

	filter := boom.NewBloomFilter(BloomLimit, BloomFalsePositiveRate)
	logger.WithField("count", len(allHashes)).Debug("adding entries to bloom filter")
	for _, hash := range allHashes {
		logger.WithField("cid", hash).Trace("adding to bloom filter")
		filter.Add([]byte(hash))
	}
	var buf bytes.Buffer
	_, err := filter.WriteTo(&buf)
	if err != nil {
		return fmt.Errorf("failed to write bloom filter bytes: %v", err)
	}
	bloomEncoded := base64.StdEncoding.EncodeToString(buf.Bytes())

	r, err := storage.ipfs.FilesRead(ctx, BloomPath(user.User))
	if err == nil {
		prevBloom, _ := io.ReadAll(r)
		if string(prevBloom) == bloomEncoded {
			logger.Info("bloom filter remains the same - skipping provide call")
			return nil
		}
	}

	storage.ipfs.FilesMkdir(storage.ctx, RepoDir(user.User), ipfsapi.FilesMkdir.Parents(true))
	bloomPath := BloomPath(user.User)
	storage.ipfs.FilesRm(storage.ctx, bloomPath, true)
	if err := storage.ipfs.FilesWrite(
		ctx, bloomPath, bytes.NewBuffer([]byte(bloomEncoded)), ipfsapi.FilesWrite.Create(true),
	); err != nil {
		return fmt.Errorf("failed to write bloom: %v", err)
	}

	idResp, err := storage.ipfs.ID()
	if err != nil {
		return fmt.Errorf("failed to get peer id: %v", err)
	}

	if err := storage.router.Provide(ctx, user.User, idResp.ID, bloomEncoded); err != nil {
		return fmt.Errorf("failed to update router: %v", err)
	}

	return nil
}
