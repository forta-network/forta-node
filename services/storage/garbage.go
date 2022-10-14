package storage

import (
	"context"
	"fmt"
	"path"
	"sort"
	"time"

	ipfsapi "github.com/ipfs/go-ipfs-api"
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

type userInfo struct {
	User         string
	ContentKinds []string
}

func (storage *Storage) doCollectGarbage(ctx context.Context) error {
	// just attempt creating the base dir to avoid unnecessary errors
	storage.ipfs.FilesMkdir(ctx, DefaultBasePath, ipfsapi.FilesMkdir.Parents(true))

	list, err := storage.ipfs.FilesLs(ctx, DefaultBasePath)
	if err != nil {
		return fmt.Errorf("failed to list the base storage path: %v", err)
	}

	var users []*userInfo
	for _, stat := range list {
		userName := stat.Name
		contentList, err := storage.ipfs.FilesLs(ctx, path.Join(DefaultBasePath, userName))
		if err != nil {
			return fmt.Errorf("failed to get the content kinds for user '%s': %v", userName, err)
		}
		var contentKinds []string
		for _, kind := range contentList {
			contentKinds = append(contentKinds, kind.Name)
		}
		users = append(users, &userInfo{
			User:         userName,
			ContentKinds: contentKinds,
		})
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
	list, err := storage.ipfs.FilesLs(ctx, contentDir, ipfsapi.FilesLs.Stat(true))
	if err != nil {
		return fmt.Errorf("error while listing '%s': %v", contentDir, err)
	}
	// ensure it's sorted in alphabetical order (ascending)
	sort.Slice(list, func(i, j int) bool {
		return sort.StringsAreSorted([]string{list[i].Name, list[j].Name})
	})
	limit := ContentLimit(kind)

	var oldEntries []*ipfsapi.MfsLsEntry
	oldCount := len(list) - limit
	if oldCount > 0 {
		logger.WithField("old", oldCount).Info("detected old entries - will remove")
		oldEntries = list[:len(list)-limit]
	}
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
