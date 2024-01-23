package lifecycle

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/forta-network/forta-node/clients/docker"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/services/components/containers"
	"github.com/forta-network/forta-node/services/components/metrics"
	"github.com/forta-network/forta-node/services/components/registry"
	"github.com/forta-network/forta-node/store"
	log "github.com/sirupsen/logrus"
)

// Timeouts
var (
	botRemoveTimeout = time.Second * 5

	heartbeatBotDuration     = time.Minute * 5
	heartbeatBotLoadInterval = time.Hour
)

// BotLifecycleManager manages lifecycles of running bots.
type BotLifecycleManager interface {
	ManageBots(ctx context.Context) error
	CleanupUnusedBots(ctx context.Context) error
	ExitInactiveBots(ctx context.Context) error
	RestartExitedBots(ctx context.Context) error
	TearDownRunningBots(ctx context.Context)
}

type botLifecycleManager struct {
	botRegistry       registry.BotRegistry
	botClient         containers.BotClient
	botPool           BotPoolUpdater
	lifecycleMetrics  metrics.Lifecycle
	botMonitor        BotMonitor
	lastHeartbeatLoad time.Time

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

// addHeartbeatBotIfDue appends a heartbeat bot to the list of bots
// Every hour, it includes the bot for a 5 minute period
// It doesn't run all the time to prevent large batch volume increases
func (blm *botLifecycleManager) addHeartbeatBotIfDue(cfgs []config.AgentConfig) []config.AgentConfig {
	timeSinceLast := time.Since(blm.lastHeartbeatLoad)

	// if more than an hour ago, or if within 5 minutes, continue to include the heartbeat bot
	if timeSinceLast > heartbeatBotLoadInterval || timeSinceLast < heartbeatBotDuration {
		hb, err := blm.botRegistry.LoadHeartbeatBot()
		if err != nil && err != store.ErrLocalMode {
			blm.lifecycleMetrics.SystemError("load.heartbeat.bot", err)
		}
		if hb != nil {
			cfgs = append(cfgs, *hb)
			if timeSinceLast > heartbeatBotLoadInterval {
				blm.lastHeartbeatLoad = time.Now().UTC()
			}
		}
	}
	return cfgs
}

// ManageBots starts containers for assigned bots and stops the containers for unassigned
// bots and lets other services know.
func (blm *botLifecycleManager) ManageBots(ctx context.Context) error {
	botsToRun, err := blm.botRegistry.LoadAssignedBots()
	if err != nil {
		blm.lifecycleMetrics.SystemError("load.assigned.bots", err)
		return fmt.Errorf("failed to load assigned bots: %v", err)
	}
	blm.lifecycleMetrics.SystemStatus("load.assigned.bots", strconv.Itoa(len(botsToRun)))

	// append the heartbeat bot if due to execute
	botsToRun = blm.addHeartbeatBotIfDue(botsToRun)

	// find the removed bots and remove them from the pool
	removedBotConfigs := FindMissingBots(blm.runningBots, botsToRun)
	if len(removedBotConfigs) > 0 {
		if err := blm.botPool.RemoveBotsWithConfigs(removedBotConfigs); err != nil {
			log.WithError(err).Error("error removing bots")
			blm.lifecycleMetrics.SystemError("remove.bots.with.configs", err)
		}
		blm.lifecycleMetrics.StatusStopping(removedBotConfigs...)
	}

	// then wait a little to let the bot pool process this
	// this is just for avoiding bot client error noise
	time.Sleep(botRemoveTimeout)

	// then stop the containers
	for _, removedBotConfig := range removedBotConfigs {
		if err := blm.botClient.TearDownBot(ctx, removedBotConfig.ContainerName()); err != nil {
			log.WithError(err).WithField("container", removedBotConfig.ContainerName()).
				Warn("failed to tear down unassigned bot container")
			blm.lifecycleMetrics.BotError("unassigned.teardown", err, removedBotConfig)
		}
	}

	// find the bot containers to start
	addedBotConfigs := FindExtraBots(blm.runningBots, botsToRun)

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
			botsToRun = Drop(addedBotConfig, botsToRun)
			blm.lifecycleMetrics.FailurePull(downloadErrs[i], addedBotConfig)
			continue
		}

		// skip if the bot could not start
		err := blm.botClient.LaunchBot(ctx, addedBotConfig)
		if err != nil {
			log.WithError(err).WithField("container", addedBotConfig.ContainerName()).
				Warn("failed to launch bot")
			// drop the bot from the list so it can be picked again next time
			botsToRun = Drop(addedBotConfig, botsToRun)
			blm.lifecycleMetrics.FailureLaunch(err, addedBotConfig)
			continue
		}
	}

	// then update the pool with latest bots
	if err := blm.botPool.UpdateBotsWithLatestConfigs(botsToRun); err != nil {
		blm.lifecycleMetrics.SystemError("update.bots.with.latest.configs", err)
	}
	blm.lifecycleMetrics.StatusRunning(botsToRun...)
	blm.botMonitor.MonitorBots(GetBotIDs(botsToRun))

	blm.runningBots = botsToRun
	return nil
}

// CleanupUnusedBots cleans up unused bots.
func (blm *botLifecycleManager) CleanupUnusedBots(ctx context.Context) error {
	if len(blm.runningBots) == 0 {
		return nil
	}

	botContainers, err := blm.botClient.LoadBotContainers(ctx)
	if err != nil {
		return fmt.Errorf("failed to load bot containers during bot cleanup: %v", err)
	}

	if len(botContainers) == 0 {
		return nil
	}

	for _, botContainer := range botContainers {
		botContainerName := botContainer.Names[0][1:]
		_, ok := blm.findBotConfig(botContainerName)
		if ok {
			continue
		}

		if err := blm.botClient.TearDownBot(ctx, botContainerName); err != nil {
			log.WithField("botContainer", botContainerName).WithError(err).
				Error("error while tearing down the unused bot")
		}
	}

	return nil
}

// ExitInactiveBots exits inactive bots so the restart can pick them up later.
func (blm *botLifecycleManager) ExitInactiveBots(ctx context.Context) error {
	inactiveBotIDs := blm.botMonitor.GetInactiveBots()
	if len(inactiveBotIDs) == 0 {
		return nil
	}
	inactiveCfgs := make([]config.AgentConfig, 0, len(inactiveBotIDs))
	for _, inactiveBotID := range inactiveBotIDs {
		botConfig, found := blm.findBotConfigByID(inactiveBotID)
		logger := log.WithField("bot", inactiveBotID)
		if !found {
			logger.Warn("could not find the config for inactive bot - skipping stop", inactiveBotID)
			continue
		}

		// TODO: Do not support inactive bot restarts until v2 health check support.
		if botConfig.ProtocolVersion >= 2 {
			continue
		}

		inactiveCfgs = append(inactiveCfgs, botConfig)
		logger.Info("killing inactive bot for reinitialization")
		if err := blm.botClient.StopBot(ctx, botConfig); err != nil {
			logger.WithError(err).Error("failed to stop the inactive bot")
			blm.lifecycleMetrics.FailureStop(fmt.Errorf("failed to stop the inactive bot: %v", err.Error()), botConfig)
		}
	}
	// sends metric with shard IDs
	blm.lifecycleMetrics.StatusInactive(inactiveCfgs...)
	return nil
}

// RestartExitedBots restarts bot containers when they are down and lets other services know.
func (blm *botLifecycleManager) RestartExitedBots(ctx context.Context) error {
	botContainers, err := blm.botClient.LoadBotContainers(ctx)
	if err != nil {
		blm.lifecycleMetrics.SystemError("load.bot.containers", fmt.Errorf("failed to load bot containers: %v", err.Error()))
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
		logger = log.WithField("botId", restartedBotConfig.ID)
		logger.Warn("restarting bot container")
		blm.lifecycleMetrics.ActionRestart(restartedBotConfig)
		if err := blm.botClient.StartWaitBotContainer(ctx, botContainer.ID); err != nil {
			logger.WithError(err).Error("failed to start exited bot container")
			blm.lifecycleMetrics.BotError("start.exited.bot.container", fmt.Errorf("failed to start exited bot container: %v", err.Error()), restartedBotConfig)
			continue
		}
		restartedBotConfigs = append(restartedBotConfigs, restartedBotConfig)
	}

	// let the bot pool reconnect to the restarted bots
	if len(restartedBotConfigs) > 0 {
		if err := blm.botPool.ReconnectToBotsWithConfigs(restartedBotConfigs); err != nil {
			blm.lifecycleMetrics.SystemError("reinit.bots.with.configs", fmt.Errorf("failed to reinit bots with configs: %v", err.Error()))
		}
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
	if err := blm.botPool.RemoveBotsWithConfigs(blm.runningBots); err != nil {
		blm.lifecycleMetrics.SystemError("teardown.remove.bots.with.configs", err)
		log.WithError(err).Error("error removing bots with configs")
	}

	// then wait a little to let the bot pool process this
	// this is just for avoiding bot client error noise
	time.Sleep(botRemoveTimeout)

	// then stop the containers
	for _, runningBotConfig := range blm.runningBots {
		err := blm.botClient.TearDownBot(ctx, runningBotConfig.ContainerName())
		if err != nil {
			blm.lifecycleMetrics.BotError("teardown.bot", err, runningBotConfig)
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
