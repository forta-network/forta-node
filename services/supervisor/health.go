package supervisor

import (
	"errors"

	"github.com/forta-network/forta-core-go/utils"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/services"

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
	sup.mu.RLock()
	defer sup.mu.RUnlock()
	for _, knownContainer := range sup.containers {
		var foundContainer *types.Container

		// this has a threshold so that the healthcheck doesn't fail while a container is starting
		err := utils.TryTimes(func(attempt int) error {
			var err error
			foundContainer, err = sup.client.GetContainerByID(sup.ctx, knownContainer.ID)
			currAttempt := attempt + 1
			if err != nil && errors.Is(err, clients.ErrContainerNotFound) {
				log.Warnf("healthcheck: container '%s' with id '%s' was not found (attempt=%d/%d)", knownContainer.Name, knownContainer.ID, currAttempt, maxAttempts)
				return err
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
		if foundContainer == nil {
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
		logger := log.WithField("name", knownContainer.Name)

		containerDetails, err := sup.client.InspectContainer(sup.ctx, foundContainer.ID)
		if err != nil {
			return err
		}
		if containerDetails.State.ExitCode == services.ExitCodeTriggered {
			logger.Info("detected internal exit trigger - exiting")
			services.TriggerExit()
			return nil
		}

		logger.Warn("starting exited container")
		_, err = sup.client.StartContainer(sup.ctx, knownContainer.Config)
		if err != nil {
			return fmt.Errorf("failed to start container '%s': %v", knownContainer.Name, err)
		}
		return nil
	default:
		log.WithField("name", knownContainer.Name).Panicf("unhandled container state: %s", foundContainer.State)
	}
	return nil
}
