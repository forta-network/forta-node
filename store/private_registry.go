package store

import (
	"context"
	"errors"
	"strconv"
	"sync"

	"github.com/forta-network/forta-core-go/manifest"
	"github.com/forta-network/forta-core-go/registry"
	"github.com/forta-network/forta-node/config"
	log "github.com/sirupsen/logrus"
)

type privateRegistryStore struct {
	ctx context.Context
	cfg config.Config
	rc  registry.Client
	mc  manifest.Client
	mu  sync.Mutex
}

func (rs *privateRegistryStore) FindScannerShardIDForBot(agentID, scannerAddress string) (uint, uint, uint, error) {
	return 0, 0, 0, nil
}

func (rs *privateRegistryStore) GetAgentsIfChanged(scanner string) ([]*config.AgentConfig, bool, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	var agentConfigs []*config.AgentConfig

	// load by image references
	for i, agentImage := range rs.cfg.LocalModeConfig.BotImages {
		if len(agentImage) == 0 {
			continue
		}
		// forta-agent-1, forta-agent-2, forta-agent-3, ...
		agentID := strconv.Itoa(i + 1)
		agentConfigs = append(agentConfigs, rs.makePrivateModeAgentConfig(agentID, agentImage))
	}

	// load by bot IDs
	for _, agentID := range rs.cfg.LocalModeConfig.BotIDs {
		agt, err := rs.rc.GetAgent(agentID)
		logger := log.WithFields(
			log.Fields{
				"botID": agentID,
			},
		)
		if err != nil {
			logger.WithError(err).Error("failed to get bot from registry")
			continue
		}
		agtCfg, err := loadBot(rs.ctx, rs.cfg, rs.mc, agentID, agt.Manifest)
		if err != nil {
			logger.WithError(err).Error("failed to load bot")
			continue
		}
		agentConfigs = append(agentConfigs, agtCfg)
	}

	return agentConfigs, true, nil
}

func (rs *privateRegistryStore) FindAgentGlobally(agentID string) (*config.AgentConfig, error) {
	return nil, errors.New("feature not available (private/local registry)")
}

func (rs *privateRegistryStore) makePrivateModeAgentConfig(id string, image string) *config.AgentConfig {
	return &config.AgentConfig{
		ID:      id,
		Image:   image,
		IsLocal: true,
	}
}

func NewPrivateRegistryStore(ctx context.Context, cfg config.Config) (*privateRegistryStore, error) {
	mc, err := manifest.NewClient(cfg.Registry.IPFS.GatewayURL)
	if err != nil {
		return nil, err
	}

	rc, err := GetRegistryClient(
		ctx, cfg, registry.ClientConfig{
			JsonRpcUrl: cfg.Registry.JsonRpc.Url,
			ENSAddress: cfg.ENSConfig.ContractAddress,
			Name:       "registry-store",
		},
	)
	if err != nil {
		return nil, err
	}
	return &privateRegistryStore{
		ctx: ctx,
		cfg: cfg,
		mc:  mc,
		rc:  rc,
	}, nil
}
