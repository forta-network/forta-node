package lifecycle

import (
	"context"
	"fmt"
	"time"

	"github.com/forta-network/forta-node/clients/docker"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/services/components/containers"
	"github.com/forta-network/forta-node/services/components/metrics"
	"github.com/forta-network/forta-node/services/components/registry"
	log "github.com/sirupsen/logrus"
)

// Timeouts
var (
	botRemoveTimeout = time.Second * 5
)

// BotLifecycleManager manages lifecycles of running bots.
type BotLifecycleManager interface {
	ManageBots(ctx context.Context) error
	RestartExitedBots(ctx context.Context) error
}

type botLifecycleManager struct {
	botRegistry      registry.BotRegistry
	botClient        containers.BotClient
	botPool          BotPoolUpdater
	lifecycleMetrics metrics.Lifecycle

	runningBots []config.AgentConfig
}

var _ BotLifecycleManager = &botLifecycleManager{}

// NewManager creates new.
func NewManager(
	botRegistry registry.BotRegistry, botClient containers.BotClient,
	botPool BotPoolUpdater, lifecycleMetrics metrics.Lifecycle,
) *botLifecycleManager {
	return &botLifecycleManager{
		botRegistry:      botRegistry,
		botClient:        botClient,
		botPool:          botPool,
		lifecycleMetrics: lifecycleMetrics,
	}
}

// ManageBots starts containers for assigned bots and stops the containers for unassigned
// bots and lets other services know.
func (blm *botLifecycleManager) ManageBots(ctx context.Context) error {
	assignedBots, err := blm.botRegistry.LoadAssignedBots()
	if err != nil {
		return fmt.Errorf("failed to load assigned bots: %v", err)
	}

	// find the removed bots and remove them from the pool
	removedBotConfigs := FindMissingBots(blm.runningBots, assignedBots)
	if len(removedBotConfigs) > 0 {
		blm.botPool.RemoveBotsWithConfigs(removedBotConfigs)
		blm.lifecycleMetrics.StatusStopping(removedBotConfigs...)
	}

	// then wait a little to let the bot pool process this
	time.Sleep(botRemoveTimeout)

	// then stop the containers
	for _, removedBotConfig := range removedBotConfigs {
		err := blm.botClient.TearDownBot(ctx, removedBotConfig)
		if err != nil {
			log.WithField("container", removedBotConfig.ContainerName()).
				Warn("failed to stop unassigned bot container")
		}
	}

	// find the bot containers to start
	addedBotConfigs := FindExtraBots(blm.runningBots, assignedBots)

	// then download all images concurrently
	var downloadErrs []error
	if len(addedBotConfigs) > 0 {
		downloadErrs = blm.botClient.EnsureBotImages(ctx, addedBotConfigs)
	}

	// and start them
	for i, addedBotConfig := range addedBotConfigs {
		// skip start if we could not download
		if downloadErrs[i] != nil {
			log.WithFields(log.Fields{
				"bot":   addedBotConfig.ID,
				"image": addedBotConfig.Image,
				"error": downloadErrs[i],
			}).Error("bot image download failed - skipping launch")
			// drop the bot from the list so it can be picked again next time
			assignedBots = Drop(addedBotConfig, assignedBots)
			blm.lifecycleMetrics.FailurePull(addedBotConfig)
			continue
		}

		// skip if the bot could not start
		err := blm.botClient.LaunchBot(ctx, addedBotConfig)
		if err != nil {
			log.WithField("container", addedBotConfig.ContainerName()).
				Warn("failed to start assigned bot container")
			// drop the bot from the list so it can be picked again next time
			assignedBots = Drop(addedBotConfig, assignedBots)
			blm.lifecycleMetrics.FailureLaunch(addedBotConfig)
			continue
		}
	}

	// then update the pool with latest bots
	blm.botPool.UpdateBotsWithLatestConfigs(assignedBots)
	blm.lifecycleMetrics.StatusRunning(assignedBots...)

	blm.runningBots = assignedBots
	return nil
}

// RestartExitedBots restarts bot containers when they are down and lets other services know.
func (blm *botLifecycleManager) RestartExitedBots(ctx context.Context) error {
	botContainers, err := blm.botClient.LoadBotContainers(ctx)
	if err != nil {
		return fmt.Errorf("failed to load bot containers: %v", err)
	}

	// find exited bot containers and restart them
	var restartedBotConfigs []config.AgentConfig
	for _, botContainer := range botContainers {
		if botContainer.State != "exited" {
			continue
		}

		containerName := docker.GetContainerName(botContainer)
		logger := log.WithField("container", containerName)
		restartedBotConfig, found := blm.findBotConfig(containerName)
		if !found {
			logger.Warn("could not find config for exited bot container")
			continue
		}

		logger.Warn("restarting bot container")
		blm.lifecycleMetrics.ActionRestart(restartedBotConfig)
		if err := blm.botClient.StartWaitBotContainer(ctx, botContainer.ID); err != nil {
			logger.WithError(err).Error("failed to start exited bot container")
			continue
		}
		restartedBotConfigs = append(restartedBotConfigs, restartedBotConfig)
	}

	// let bot pool reinitialize the restart bots
	blm.botPool.ReinitBotsWithConfigs(restartedBotConfigs)
	return nil
}

func (blm *botLifecycleManager) findBotConfig(containerName string) (config.AgentConfig, bool) {
	for _, bot := range blm.runningBots {
		if bot.ContainerName() == containerName {
			return bot, true
		}
	}
	return config.AgentConfig{}, false
}