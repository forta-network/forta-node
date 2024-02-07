package lifecycle

import (
	"context"
	"sync"
	"time"

	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/services/components/botio"
	"github.com/forta-network/forta-node/services/components/metrics"
	log "github.com/sirupsen/logrus"
)

// BotPool contains all the bot clients which we can forward the alert, block and tx requests
// to and receive the results from, and manages the lifecyle of these bots.
type BotPool interface {
	BotPoolUpdater
	botio.BotPool
}

// BotPoolUpdater updates bots.
type BotPoolUpdater interface {
	UpdateBotsWithLatestConfigs(messaging.AgentPayload) error
	RemoveBotsWithConfigs(messaging.AgentPayload) error
	ReconnectToBotsWithConfigs(messaging.AgentPayload) error
}

type botPool struct {
	ctx context.Context

	botClients []botio.BotClient
	mu         sync.RWMutex

	waitBots int
	botWg    *sync.WaitGroup

	waitInit bool // hack to make testing synchronous

	lifecycleMetrics metrics.Lifecycle
	botClientFactory botio.BotClientFactory
}

var _ BotPool = &botPool{}

// NewBotPool creates a new bot pool.
func NewBotPool(
	ctx context.Context, lifecycleMetrics metrics.Lifecycle,
	botClientFactory botio.BotClientFactory, waitBots int,
) *botPool {
	botPool := &botPool{
		ctx:              ctx,
		waitBots:         waitBots,
		lifecycleMetrics: lifecycleMetrics,
		botClientFactory: botClientFactory,
	}
	if waitBots > 0 {
		botPool.botWg = &sync.WaitGroup{}
		botPool.botWg.Add(waitBots)
		go botPool.logBotWait()
		go botPool.logBotStatuses()
	}
	return botPool
}

func (bp *botPool) logBotWait() {
	if bp.botWg != nil {
		bp.botWg.Wait()
		log.Info("started all bots")
	}
}

func (bp *botPool) logBotStatuses() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		bp.doLogBotStatuses()
	}
}

func (bp *botPool) doLogBotStatuses() {
	for _, agent := range bp.GetCurrentBotClients() {
		agent.LogStatus()
	}
}

// UpdateBotsWithLatestConfigs starts and adds new bots and updates the config of updated bots.
func (bp *botPool) UpdateBotsWithLatestConfigs(latestConfigs messaging.AgentPayload) error {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	// add new bots
	var latestBotClients []botio.BotClient
	for _, botConfig := range latestConfigs {
		logger := botLog(botConfig)
		botClient, ok := bp.getBotClient(botConfig.ContainerName())
		if ok && !botClient.IsClosed() {
			logger.Debug("bot client already exists - skipping update")
			latestBotClients = append(latestBotClients, botClient)
			continue
		}

		if ok && botClient.IsClosed() {
			logger.Info("replacing closed bot client")
		} else {
			logger.Info("adding new bot client")
		}

		latestBotClients = append(latestBotClients, bp.startBotClient(botConfig))
	}
	bp.botClients = latestBotClients

	// updated the config of the bots that have different config
	updatedBotConfigs := FindUpdatedBots(bp.getConfigsUnsafe(), latestConfigs)
	for _, updatedBotConfig := range updatedBotConfigs {
		logger := botLog(updatedBotConfig)
		botClient, ok := bp.getBotClient(updatedBotConfig.ContainerName())
		if !ok {
			logger.Info("could not find the updated bot! skipping")
			continue
		}
		botClient.SetConfig(updatedBotConfig)
	}
	if len(updatedBotConfigs) > 0 {
		bp.lifecycleMetrics.ActionUpdate(updatedBotConfigs...)
	}

	// if the pool needs to wait for the first time, detect that and mark done
	if bp.waitBots > 0 && bp.botWg != nil {
		bp.botWg.Add(-len(latestConfigs))
		bp.botWg = nil
	}

	return nil
}

func (bp *botPool) startBotClient(botConfig config.AgentConfig) botio.BotClient {
	botClient := bp.botClientFactory.NewBotClient(bp.ctx, botConfig)

	if bp.waitInit {
		botClient.Initialize()
	} else {
		go botClient.Initialize()
	}
	botClient.StartProcessing()

	return botClient
}

// RemoveBotsWithConfigs closes and discards the bots to be removed.
func (bp *botPool) RemoveBotsWithConfigs(removedBotConfigs messaging.AgentPayload) error {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	// close and discard the removed bots
	for _, removedBotConfig := range removedBotConfigs {
		logger := botLog(removedBotConfig)
		botClient, ok := bp.getBotClient(removedBotConfig.ContainerName())
		if !ok {
			logger.Info("could not find the removed bot! skipping")
			continue
		}
		_ = botClient.Close()
	}

	// find the bots we are not supposed to remove and keep them
	var preservedBots []botio.BotClient
	for _, preservedBotConfig := range FindExtraBots(removedBotConfigs, bp.getConfigsUnsafe()) {
		botClient, ok := bp.getBotClient(preservedBotConfig.ContainerName())
		if ok {
			preservedBots = append(preservedBots, botClient)
		}
	}

	bp.botClients = preservedBots

	return nil
}

// ReconnectToBotsWithConfigs reinitializes bots.
func (bp *botPool) ReconnectToBotsWithConfigs(reconnectedBots messaging.AgentPayload) error {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	var latestBotClients []botio.BotClient
	for _, botClient := range bp.botClients {
		botConfig, found := FindBot(botClient.Config().ContainerName(), reconnectedBots)
		// if found, close old and replace with new
		if found {
			_ = botClient.Close()
			botClient = bp.startBotClient(botConfig)
		}
		// append previous or new one one, depending on the previous step
		latestBotClients = append(latestBotClients, botClient)
	}
	bp.botClients = latestBotClients

	return nil
}

// GetCurrentBotClients returns the current bot clients safely.
func (bp *botPool) GetCurrentBotClients() []botio.BotClient {
	bp.mu.RLock()
	defer bp.mu.RUnlock()

	return bp.botClients
}

// WaitForAll waits for bot clients to start if the count was provided during initialization.
func (bp *botPool) WaitForAll() {
	if bp.botWg == nil {
		return
	}
	bp.botWg.Wait()
}

func botLog(botConfig config.AgentConfig) *log.Entry {
	return log.WithField("bot", botConfig.ID).WithField("container", botConfig.ContainerName())
}

func (bp *botPool) getBotClient(containerName string) (botio.BotClient, bool) {
	for _, bot := range bp.botClients {
		if bot.Config().ContainerName() == containerName {
			return bot, true
		}
	}
	return nil, false
}

func (bp *botPool) getConfigsUnsafe() (allConfigs []config.AgentConfig) {
	for _, botClient := range bp.botClients {
		allConfigs = append(allConfigs, botClient.Config())
	}
	return
}
