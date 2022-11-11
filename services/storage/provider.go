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
	provideContentInterval        = time.Minute
	provideContentLimit           = 10000
	provideBloomFalsePositiveRate = 0.0001
)

func (storage *Storage) provideContent(ctx context.Context) {
	ticker := time.NewTicker(time.Minute * 1)
	for {
		select {
		case <-ctx.Done():
			log.Info("exiting content provider")
			return
		case <-ticker.C:
			if err := storage.doProvideContent(context.Background()); err != nil {
				log.WithError(err).Error("error while providing content")
			} else {
				log.Info("finished providing content refs")
			}
		}
	}
}

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

	var allEntries []*ipfsapi.MfsLsEntry
	// prioritize some content over other
	for _, kind := range []string{
		KindBatchReceipt,
	} {
		if !user.HasContent(kind) {
			continue
		}

		entries, _, err := storage.getContentInDir(ctx, user.User, kind)
		if err != nil {
			return err
		}
		allEntries = append(allEntries, entries...)
		if len(allEntries) > provideContentLimit {
			allEntries = allEntries[:provideContentLimit]
			break
		}
	}
	if len(allEntries) == 0 {
		logger.Info("no entries found - skipping provide call")
		return nil
	}

	filter := boom.NewBloomFilter(provideContentLimit, provideBloomFalsePositiveRate)
	logger.WithField("count", len(allEntries)).Debug("adding entries to bloom filter")
	for _, entry := range allEntries {
		logger.WithField("cid", entry.Hash).Trace("adding to bloom filter")
		filter.Add([]byte(entry.Hash))
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
	if err := storage.ipfs.FilesWrite(
		ctx, BloomPath(user.User), bytes.NewBuffer([]byte(bloomEncoded)), ipfsapi.FilesWrite.Create(true),
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
