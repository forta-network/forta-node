package registry

import (
	"context"
	"fmt"

	"github.com/forta-network/forta-node/store"

	"github.com/ethereum/go-ethereum/common"
	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/ethereum"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/services/registry/regtypes"
	log "github.com/sirupsen/logrus"
)

// RegistryService is the registry service interface.
type RegistryService interface {
	GetLatestBots() ([]*config.AgentConfig, error)
}

// registryService listens to the agent scanner list changes so the node can stay in sync.
type registryService struct {
	cfg            config.Config
	scannerAddress common.Address
	ethClient      ethereum.Client

	registryStore store.RegistryStore

	agentsConfigs []*config.AgentConfig

	lastChecked        health.TimeTracker
	lastChangeDetected health.TimeTracker
	lastErr            health.ErrorTracker
}

// IPFSClient interacts with an IPFS Gateway.
type IPFSClient interface {
	GetAgentFile(cid string) (*regtypes.AgentFile, error)
}

// EthClient interacts with the Ethereum API.
type EthClient interface {
	ethereum.Client
}

// New creates a new service.
func New(cfg config.Config, scannerAddress common.Address, ethClient ethereum.Client) *registryService {
	return &registryService{
		cfg:            cfg,
		scannerAddress: scannerAddress,
		ethClient:      ethClient,
	}
}

// Init only initializes the service.
func (rs *registryService) Init() error {
	var (
		regStr store.RegistryStore
		err    error
	)
	if rs.cfg.LocalModeConfig.Enable {
		regStr, err = store.NewPrivateRegistryStore(context.Background(), rs.cfg)
	} else {
		regStr, err = store.NewRegistryStore(context.Background(), rs.cfg, rs.ethClient)
	}
	if err != nil {
		return err
	}
	rs.registryStore = regStr
	return nil
}

// Start initializes and starts the registry service.
func (rs *registryService) Start() error {
	return rs.Init()
}

// GetLatestBots returns the latest bot list for the running scanner.
func (rs *registryService) GetLatestBots() ([]*config.AgentConfig, error) {
	rs.lastChecked.Set()
	agts, changed, err := rs.registryStore.GetAgentsIfChanged(rs.scannerAddress.Hex())
	if err != nil {
		return nil, fmt.Errorf("failed to get the latest bot list: %v", err)
	}

	logger := log.WithField("service", rs.Name())
	if changed {
		rs.lastChangeDetected.Set()
		rs.agentsConfigs = agts
		logger.WithField("count", len(agts)).Info("updated bot list")
	} else {
		logger.Info("no bot list changes detected")
	}

	return rs.agentsConfigs, nil
}

// Stop stops the registry service.
func (rs *registryService) Stop() error {
	return nil
}

// Name returns the name of the service.
func (rs *registryService) Name() string {
	return "registry"
}

// Health implements the health.Reporter interface.
func (rs *registryService) Health() health.Reports {
	return health.Reports{
		rs.lastErr.GetReport("event.checked.error"),
		&health.Report{
			Name:    "event.checked.time",
			Status:  health.StatusInfo,
			Details: rs.lastChecked.String(),
		},
		&health.Report{
			Name:    "event.change-detected.time",
			Status:  health.StatusInfo,
			Details: rs.lastChangeDetected.String(),
		},
	}
}
