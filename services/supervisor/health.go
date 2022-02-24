package supervisor

import (
	"github.com/forta-protocol/forta-node/utils"

	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	log "github.com/sirupsen/logrus"
)

const defaultHealthCheckInterval = time.Second * 5
const maxAttempts = 10

func (sup *SupervisorService) healthCheck() {
	ticker := time.NewTicker(defaultHealthCheckInterval)
	for {
		select {
		case <-sup.ctx.Done():
			ticker.Stop()
			return

		case <-ticker.C:
			if err := sup.doHealthCheck(); err != nil {
				log.Errorf("failed to do health check: %v", err)
			}
		}
	}
}

func (sup *SupervisorService) doHealthCheck() error {
	containersList, err := sup.client.GetContainers(sup.ctx)
	if err != nil {
		return fmt.Errorf("failed to get containers list: %v", err)
	}
	sup.mu.RLock()
	defer sup.mu.RUnlock()
	for _, knownContainer := range sup.containers {
		var foundContainer *types.Container
		var ok bool

		// this has a threshold so that the healthcheck doesn't fail while a container is starting
		err := utils.TryTimes(func(attempt int) error {
			foundContainer, ok = containersList.FindByID(knownContainer.ID)
			currAttempt := attempt + 1
			if !ok {
				notFoundErr := fmt.Errorf("healthcheck: container '%s' with id '%s' was not found (attempt=%d/%d)", knownContainer.Name, knownContainer.ID, currAttempt, maxAttempts)
				log.Warnf(notFoundErr.Error())
				// get containers again, so that we can get updated info
				containersList, err = sup.client.GetContainers(sup.ctx)
				if err != nil {
					return fmt.Errorf("failed to get containers list: %v", err)
				}
				return notFoundErr
			}
			// If the container is found alive at later attempts, make it obvious.
			if currAttempt > 1 {
				log.Infof("healthcheck: container '%s' with id '%s' was found alive (attempt=%d/%d)", knownContainer.Name, knownContainer.ID, currAttempt, maxAttempts)
			}
			return nil
		}, maxAttempts, 1*time.Second)
		if err != nil {
			// If this ever happens, then we have a critical gap in our logic.
			log.Error(err.Error())
			continue
		}
		if err := sup.ensureUp(knownContainer, foundContainer); err != nil {
			return err
		}
	}
	return nil
}

func (sup *SupervisorService) ensureUp(knownContainer *Container, foundContainer *types.Container) error {
	switch foundContainer.State {
	case "created", "running", "restarting", "paused", "dead":
		return nil
	case "exited":
		log.Warnf("starting exited container '%s'", knownContainer.Name)
		_, err := sup.client.StartContainer(sup.ctx, knownContainer.Config)
		if err != nil {
			return fmt.Errorf("failed to start container '%s': %v", knownContainer.Name, err)
		}
		return nil
	default:
		log.Panicf("unhandled container state: %s", foundContainer.State)
	}
	return nil
}
