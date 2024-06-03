package registry

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/forta-network/forta-node/store"

	"github.com/ethereum/go-ethereum/common"
	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-node/config"
	log "github.com/sirupsen/logrus"
)

// BotRegistry loads the latest bots from the registry store.
type BotRegistry interface {
	LoadAssignedBots() ([]config.AgentConfig, error)
	LoadHeartbeatBot() (*config.AgentConfig, error)
	GetConfigByID(agentID string) (*config.AgentConfig, error)
	health.Reporter
}

// botRegistry retrieves the bot list changes so the node can stay in sync.
type botRegistry struct {
	cfg            config.Config
	scannerAddress common.Address

	registryStore store.RegistryStore

	mu         *sync.RWMutex
	botConfigs []config.AgentConfig

	lastChecked        health.TimeTracker
	lastChangeDetected health.TimeTracker
	lastErr            health.ErrorTracker
}

// New creates a new service.
func New(cfg config.Config, scannerAddress common.Address) (BotRegistry, error) {
	service := &botRegistry{
		cfg:            cfg,
		scannerAddress: scannerAddress,
		mu:             &sync.RWMutex{},
	}
	var (
		regStr store.RegistryStore
		err    error
	)
	if cfg.LocalModeConfig.Enable {
		regStr, err = store.NewPrivateRegistryStore(context.Background(), cfg)
	} else {
		regStr, err = store.NewRegistryStore(context.Background(), cfg)
	}
	if err != nil {
		return nil, err
	}
	service.registryStore = regStr
	return service, nil
}

func (br *botRegistry) LoadHeartbeatBot() (*config.AgentConfig, error) {
	ac, err := br.registryStore.FindAgentGlobally(config.HeartbeatBotID)
	if err != nil {
		return nil, err
	}
	if ac == nil {
		return nil, errors.New("cannot not find heartbeat bot")
	}
	return ac, nil
}

// LoadAssignedBots returns the latest bot list for the running scanner.
func (br *botRegistry) LoadAssignedBots() ([]config.AgentConfig, error) {

	br.lastChecked.Set()
	agts, changed, err := br.registryStore.GetAgentsIfChanged(br.scannerAddress.Hex())
	if err != nil {
		br.lastErr.Set(err)
		return nil, fmt.Errorf("failed to get the latest bot list: %v", err)
	}

	logger := log.WithField("component", "bot-loader")
	if changed {
		br.lastChangeDetected.Set()
		
		br.mu.Lock()
		br.botConfigs = agts
		br.mu.Unlock()
		
		logger.WithField("count", len(agts)).Info("updated bot list")
	} else {
		logger.Debug("no bot list changes detected")
	}

	return br.botConfigs, nil
}

func (br *botRegistry) GetConfigByID(agentID string) (*config.AgentConfig, error) {
	br.mu.RLock()
	defer br.mu.RUnlock()

	for _, ac := range br.botConfigs {
		if ac.ID == agentID {
			return &ac, nil
		}
	}
	return nil, fmt.Errorf("cannot find bot with ID %s", agentID)
}

// Name implements health.Reporter interface.
func (br *botRegistry) Name() string {
	return "bot-registry"
}

// Health implements the health.Reporter interface.
func (br *botRegistry) Health() health.Reports {
	return health.Reports{
		br.lastErr.GetReport("event.checked.error"),
		&health.Report{
			Name:    "event.checked.time",
			Status:  health.StatusInfo,
			Details: br.lastChecked.String(),
		},
		&health.Report{
			Name:    "event.change-detected.time",
			Status:  health.StatusInfo,
			Details: br.lastChangeDetected.String(),
		},
	}
}
