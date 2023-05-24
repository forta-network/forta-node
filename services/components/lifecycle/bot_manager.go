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
	ExitInactiveBots(ctx context.Context) error
	RestartExitedBots(ctx context.Context) error
	TearDownRunningBots(ctx context.Context)
}

type botLifecycleManager struct {
	botRegistry      registry.BotRegistry
	botClient        containers.BotClient
	botPool          BotPoolUpdater
	lifecycleMetrics metrics.Lifecycle
	botMonitor       BotMonitor

	runningBots []config.AgentConfig
}

var _ BotLifecycleManager = &botLifecycleManager{}

// NewManager creates new.
func NewManager(
	botRegistry registry.BotRegistry, botClient containers.BotClient,
	botPool BotPoolUpdater, lifecycleMetrics metrics.Lifecycle,
	botMonitor BotMonitor,
) *botLifecycleManager {
	return &botLifecycleManager{
		botRegistry:      botRegistry,
		botClient:        botClient,
		botPool:          botPool,
		lifecycleMetrics: lifecycleMetrics,
		botMonitor:       botMonitor,
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
	// this is just for avoiding bot client error noise
	time.Sleep(botRemoveTimeout)

	// then stop the containers
	for _, removedBotConfig := range removedBotConfigs {
		err := blm.botClient.TearDownBot(ctx, removedBotConfig)
		if err != nil {
			log.WithError(err).WithField("container", removedBotConfig.ContainerName()).
				Warn("failed to tear down unassigned bot container")
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
			log.WithError(err).WithField("container", addedBotConfig.ContainerName()).
				Warn("failed to launch bot")
			// drop the bot from the list so it can be picked again next time
			assignedBots = Drop(addedBotConfig, assignedBots)
			blm.lifecycleMetrics.FailureLaunch(addedBotConfig)
			continue
		}
	}

	// then update the pool with latest bots
	blm.botPool.UpdateBotsWithLatestConfigs(assignedBots)
	blm.lifecycleMetrics.StatusRunning(assignedBots...)
	blm.botMonitor.MonitorBots(GetBotIDs(assignedBots))

	blm.runningBots = assignedBots
	return nil
}

// ExitInactiveBots exits inactive bots so the restart can pick them up later.
func (blm *botLifecycleManager) ExitInactiveBots(ctx context.Context) error {
	inactiveBotIDs := blm.botMonitor.GetInactiveBots()
	if len(inactiveBotIDs) == 0 {
		return nil
	}
	for _, inactiveBotID := range inactiveBotIDs {
		botConfig, found := blm.findBotConfigByID(inactiveBotID)
		logger := log.WithField("bot", inactiveBotID)
		if !found {
			logger.Warn("could not find the config for inactive bot - skipping stop")
			continue
		}
		if err := blm.botClient.StopBot(ctx, botConfig); err != nil {
			logger.WithError(err).Error("failed to stop the inactive bot")
		}
	}
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
	if len(restartedBotConfigs) > 0 {
		blm.botPool.ReinitBotsWithConfigs(restartedBotConfigs)
	}
	return nil
}

// TearDownRunningBots tears down all running bots.
func (blm *botLifecycleManager) TearDownRunningBots(ctx context.Context) {
	if len(blm.runningBots) == 0 {
		return
	}
	log.WithField("count", len(blm.runningBots)).Info("tearing down running bots")

	// remove all bots from the pool
	blm.botPool.RemoveBotsWithConfigs(blm.runningBots)

	// then wait a little to let the bot pool process this
	// this is just for avoiding bot client error noise
	time.Sleep(botRemoveTimeout)

	// then stop the containers
	for _, runningBotConfig := range blm.runningBots {
		err := blm.botClient.TearDownBot(ctx, runningBotConfig)
		if err != nil {
			log.WithError(err).WithField("container", runningBotConfig.ContainerName()).
				Warn("failed to tear down running bot container")
		}
	}
}

func (blm *botLifecycleManager) findBotConfig(containerName string) (config.AgentConfig, bool) {
	for _, bot := range blm.runningBots {
		if bot.ContainerName() == containerName {
			return bot, true
		}
	}
	return config.AgentConfig{}, false
}

func (blm *botLifecycleManager) findBotConfigByID(botID string) (config.AgentConfig, bool) {
	for _, bot := range blm.runningBots {
		if bot.ID == botID {
			return bot, true
		}
	}
	return config.AgentConfig{}, false
}
