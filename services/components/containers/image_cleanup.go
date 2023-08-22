package containers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/docker"
	"github.com/forta-network/forta-node/services/components/registry"
	"github.com/sirupsen/logrus"
)

const imageCleanupInterval = time.Hour * 1

// ImageCleanup deals with image cleanup.
type ImageCleanup interface {
	Do(context.Context) error
}

type imageCleanup struct {
	client        clients.DockerClient
	botRegistry   registry.BotRegistry
	lastCleanup   time.Time
	exclusionList []string
}

// NewImageCleanup creates new.
func NewImageCleanup(client clients.DockerClient, botRegistry registry.BotRegistry, excludeImages ...string) *imageCleanup {
	return &imageCleanup{
		client:        client,
		botRegistry:   botRegistry,
		exclusionList: excludeImages,
	}
}

// Do does the image cleanup by finding all unused Disco images and removing them.
// The logic executes only after an interval.
func (ic *imageCleanup) Do(ctx context.Context) error {
	if time.Since(ic.lastCleanup) < imageCleanupInterval {
		return nil
	}

	containers, err := ic.client.GetContainers(ctx)
	if err != nil {
		return fmt.Errorf("failed to get containers during image cleanup: %v", err)
	}

	// we list the digest references as the main references for all images
	// because we pull by digest references
	images, err := ic.client.ListDigestReferences(ctx)
	if err != nil {
		return fmt.Errorf("failed to list images during image cleanup: %v", err)
	}

	heartbeatBot, err := ic.botRegistry.LoadHeartbeatBot()
	if err != nil {
		return fmt.Errorf("failed to load the heartbeat bot during cleanup: %v", err)
	}

	logrus.WithField("image", heartbeatBot.Image).Debug("cleanup: loaded heartbeat bot image reference")

	for _, image := range images {
		logger := logrus.WithField("image", image)

		if ic.isExcluded(image) || image == heartbeatBot.Image {
			logger.Debug("image is excluded - skipping cleanup")
			continue
		}

		if ic.isImageInUse(containers, image) {
			logger.Debug("image is in use - skipping cleanup")
			continue
		}

		if err := ic.client.RemoveImage(ctx, image); err != nil {
			logger.WithError(err).Warn("failed to cleanup unused disco image")
		} else {
			logger.Info("successfully cleaned up unused image")
		}
	}

	ic.lastCleanup = time.Now()
	return nil
}

func (ic *imageCleanup) isExcluded(ref string) bool {
	// needs to be a Disco image
	if !strings.Contains(ref, "bafybei") {
		return true
	}

	for _, excluded := range ic.exclusionList {
		// expecting the ref to include the excluded ref because
		// we specify it and it can be a subset of the full reference
		if strings.Contains(ref, excluded) {
			return true
		}
	}
	return false
}

func (lc *imageCleanup) isImageInUse(containers docker.ContainerList, image string) bool {
	for _, container := range containers {
		if container.Image == image {
			return true
		}
	}
	return false
}
