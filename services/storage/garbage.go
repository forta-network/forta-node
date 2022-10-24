package storage

import (
	"context"
	"fmt"
	"path"
	"time"

	log "github.com/sirupsen/logrus"
)

func (storage *Storage) collectGarbage(ctx context.Context) {
	ticker := time.NewTicker(time.Minute * 5)
	for {
		select {
		case <-ctx.Done():
			log.Info("exiting garbage collection")
			return
		case <-ticker.C:
			if err := storage.doCollectGarbage(context.Background()); err != nil {
				log.WithError(err).Error("error while collecting garbage")
			}
		}
	}
}

func (storage *Storage) doCollectGarbage(ctx context.Context) error {
	users, err := storage.getUsers(ctx)
	if err != nil {
		return err
	}

	for _, user := range users {
		for _, kind := range user.ContentKinds {
			if err := storage.gcContents(ctx, user.User, kind); err != nil {
				return err
			}
		}
	}

	// trigger repo GC here so it's more aggressive
	// this is more desirable for now
	if err := storage.ipfs.RepoGC(ctx); err != nil {
		return fmt.Errorf("repo gc failed: %v", err)
	}

	return nil
}

func (storage *Storage) gcContents(ctx context.Context, user, kind string) error {
	contentDir := ContentDir(user, kind)
	logger := log.WithFields(log.Fields{
		"user": user,
		"kind": kind,
		"dir":  contentDir,
	})

	_, oldEntries, err := storage.getContentInDir(ctx, user, kind)
	if err != nil {
		return err
	}
	oldCount := len(oldEntries)
	if oldCount == 0 {
		return nil
	}
	logger.WithField("old", oldCount).Info("detected old entries - will remove")

	for _, oldEntry := range oldEntries {
		contentPath := path.Join(contentDir, oldEntry.Name)
		contentLogger := log.WithFields(log.Fields{
			"path": contentPath,
			"cid":  oldEntry.Hash,
		})

		// unpin
		if err := storage.ipfs.Unpin(path.Join("/ipfs", oldEntry.Hash)); err != nil {
			contentLogger.WithError(err).Error("failed to unpin")
		} else {
			contentLogger.Info("successfully unpinned content")
		}

		// remove from MFS
		err = storage.ipfs.FilesRm(ctx, contentPath, true)
		if err != nil {
			contentLogger.WithError(err).Error("failed to remove content")
		} else {
			contentLogger.Info("successfully removed content")
		}
	}
	return nil
}
