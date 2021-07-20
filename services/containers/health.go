package containers

import (
	"OpenZeppelin/fortify-node/clients"
	"OpenZeppelin/fortify-node/utils"

	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	log "github.com/sirupsen/logrus"
)

const defaultHealthCheckInterval = time.Second * 5

func (t *TxNodeService) healthCheck() {
	ticker := time.NewTicker(defaultHealthCheckInterval)
	for {
		select {
		case <-t.ctx.Done():
			ticker.Stop()
			return

		case <-ticker.C:
			if err := t.doHealthCheck(); err != nil {
				log.Errorf("failed to do health check: %v", err)
			}
		}
	}
}

func (t *TxNodeService) doHealthCheck() error {
	containersList, err := t.client.GetContainers(t.ctx)
	if err != nil {
		return fmt.Errorf("failed to get containers list: %v", err)
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	for _, knownContainer := range t.containers {
		var foundContainer *types.Container
		var ok bool
		// give it some times to avoid a race condition with a newly-started container
		err := utils.TryTimes(func() error {
			foundContainer, ok = containersList.FindByID(knownContainer.ID)
			if !ok {
				return fmt.Errorf("container '%s' with id '%s' was not found", knownContainer.Name, knownContainer.ID)
			}
			return nil
		}, 5, 1*time.Second)
		if err != nil {
			// If this ever happens, then we have a critical gap in our logic.
			log.Error(err.Error())
			continue
		}
		if err := t.ensureUp(knownContainer, foundContainer); err != nil {
			return err
		}
	}
	return nil
}

func (t *TxNodeService) ensureUp(knownContainer *clients.DockerContainer, foundContainer *types.Container) error {
	switch foundContainer.State {
	case "created", "running", "restarting", "paused", "dead":
		return nil
	case "exited":
		log.Warnf("starting exited container '%s'", knownContainer.Name)
		_, err := t.client.StartContainer(t.ctx, knownContainer.Config)
		if err != nil {
			return fmt.Errorf("failed to start container '%s': %v", knownContainer.Name, err)
		}
		return nil
	default:
		log.Panicf("unhandled container state: %s", foundContainer.State)
	}
	return nil
}
