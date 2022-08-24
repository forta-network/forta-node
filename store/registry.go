package store

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/ipfs/go-cid"
	log "github.com/sirupsen/logrus"

	"github.com/forta-network/forta-core-go/ethereum"
	"github.com/forta-network/forta-core-go/manifest"
	"github.com/forta-network/forta-core-go/registry"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/forta-network/forta-node/config"
)

var (
	errInvalidBot = errors.New("invalid bot")
)

type RegistryStore interface {
	FindAgentGlobally(agentID string) (*config.AgentConfig, error)
	GetAgentsIfChanged(scanner string) ([]*config.AgentConfig, bool, error)
}

type registryStore struct {
	ctx context.Context
	mc  manifest.Client
	rc  registry.Client
	cfg config.Config

	lastUpdate           time.Time
	lastCompletedVersion string
	loadedBots           []*config.AgentConfig
	invalidBots          []*registry.Agent
	mu                   sync.Mutex
}

func (rs *registryStore) GetAgentsIfChanged(scanner string) ([]*config.AgentConfig, bool, error) {
	// because we peg the latest block, it can be problematic if this is called concurrently
	rs.mu.Lock()
	defer rs.mu.Unlock()
	hash, err := rs.rc.GetAssignmentHash(scanner)
	if err != nil {
		return nil, false, err
	}

	// if the scan node is disabled, it must run no agents
	isEnabledScanner, err := rs.rc.IsEnabledScanner(scanner)
	if err != nil {
		return nil, false, fmt.Errorf("failed to check if scanner is enabled: %v", err)
	}
	if !isEnabledScanner {
		return []*config.AgentConfig{}, true, nil
	}

	shouldUpdate := rs.lastCompletedVersion != hash.Hash || time.Since(rs.lastUpdate) > 1*time.Hour
	if !shouldUpdate {
		return nil, false, nil
	}

	if err := rs.rc.PegLatestBlock(); err != nil {
		return nil, false, err
	}
	defer rs.rc.ResetOpts()

	var (
		loadedBots       []*config.AgentConfig
		invalidBots      []*registry.Agent
		failedLoadingAny bool
	)
	err = rs.rc.ForEachAssignedAgent(scanner, func(bot *registry.Agent) error {
		logger := log.WithField("botId", bot.AgentID)

		// if already invalidated, remember it for next time
		if rs.isInvalidBot(bot) {
			invalidBots = append(invalidBots, bot)
			logger.WithError(err).Warn("invalid bot - skipping")
			return nil
		}
		// if already loaded, remember it for next time
		loadedBot, ok := rs.getLoadedBot(bot)
		if ok {
			loadedBots = append(loadedBots, loadedBot)
			logger.Info("already loaded bot - skipping")
			return nil
		}

		// try loading the rest of the unrecognized bots
		botCfg, err := rs.loadBot(bot.AgentID, bot.Manifest)
		switch {
		case err == nil: // yay
			loadedBots = append(loadedBots, botCfg) // remember for next time
			logger.Info("successfully loaded bot")
			return nil

		case errors.Is(err, errInvalidBot):
			invalidBots = append(invalidBots, bot) // remember for next time
			logger.WithError(err).Warn("invalid bot - skipping")
			return nil

		default:
			failedLoadingAny = true
			logger.WithError(err).Warn("could not load bot - skipping")
			// ignore agent and move on by not returning the error
			// it will not be recognized next time and will be retried above
			return nil
		}
	})
	if err != nil {
		return nil, false, err
	}

	// failed to load all: forget that this attempt existed
	// not doing this can cause getting stuck with the latest hash and zero agents
	if len(loadedBots) == 0 && failedLoadingAny {
		return nil, false, errors.New("loaded zero bots")
	}

	// remember the bots and the update time next time
	rs.loadedBots = loadedBots
	rs.invalidBots = invalidBots
	rs.lastUpdate = time.Now()

	if failedLoadingAny {
		log.Warn("failed loading some of the bots - keeping the previous list version")
	} else {
		rs.lastCompletedVersion = hash.Hash // remember next time so we don't retry the same list
	}

	return loadedBots, true, nil
}

func (rs *registryStore) FindAgentGlobally(agentID string) (*config.AgentConfig, error) {
	agt, err := rs.rc.GetAgent(agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get the latest ref: %v, agentID: %s", err, agentID)
	}
	return rs.loadBot(agentID, agt.Manifest)
}

func (rs *registryStore) getLoadedBot(bot *registry.Agent) (*config.AgentConfig, bool) {
	for _, loadedBot := range rs.loadedBots {
		if bot.Manifest == loadedBot.Manifest {
			return loadedBot, true
		}
	}
	return nil, false
}

func (rs *registryStore) isInvalidBot(bot *registry.Agent) bool {
	for _, invalidBot := range rs.invalidBots {
		if bot.Manifest == invalidBot.Manifest {
			return true
		}
	}
	return false
}

func (rs *registryStore) loadBot(agentID string, ref string) (*config.AgentConfig, error) {
	_, err := cid.Parse(ref)
	if len(ref) == 0 || err != nil {
		return nil, fmt.Errorf("%w: invalid bot cid '%s'", errInvalidBot, ref)
	}

	var agentData *manifest.SignedAgentManifest
	for i := 0; i < 10; i++ {
		agentData, err = rs.mc.GetAgentManifest(rs.ctx, ref)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to load the bot manifest: %v", err)
	}

	if agentData.Manifest.ImageReference == nil {
		return nil, fmt.Errorf("%w: invalid bot image reference, it is nil", errInvalidBot)
	}

	image, err := utils.ValidateDiscoImageRef(rs.cfg.Registry.ContainerRegistry, *agentData.Manifest.ImageReference)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid bot image reference '%s': %v", errInvalidBot, *agentData.Manifest.ImageReference, err)
	}

	return &config.AgentConfig{
		ID:       agentID,
		Image:    image,
		Manifest: ref,
	}, nil
}

func NewRegistryStore(ctx context.Context, cfg config.Config, ethClient ethereum.Client) (*registryStore, error) {
	mc, err := manifest.NewClient(cfg.Registry.IPFS.GatewayURL)
	if err != nil {
		return nil, err
	}

	rc, err := GetRegistryClient(ctx, cfg, registry.ClientConfig{
		JsonRpcUrl: cfg.Registry.JsonRpc.Url,
		ENSAddress: cfg.ENSConfig.ContractAddress,
		Name:       "registry-store",
	})
	if err != nil {
		return nil, err
	}

	return &registryStore{
		ctx: ctx,
		cfg: cfg,
		mc:  mc,
		rc:  rc,
	}, nil
}

type privateRegistryStore struct {
	ctx context.Context
	cfg config.Config
	mu  sync.Mutex
}

func (rs *privateRegistryStore) GetAgentsIfChanged(scanner string) ([]*config.AgentConfig, bool, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	var agentConfigs []*config.AgentConfig
	for i, agentImage := range rs.cfg.LocalModeConfig.BotImages {
		if len(agentImage) == 0 {
			continue
		}
		// forta-agent-1, forta-agent-2, forta-agent-3, ...
		agentID := strconv.Itoa(i + 1)
		agentConfigs = append(agentConfigs, rs.makePrivateModeAgentConfig(agentID, agentImage))
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
	return &privateRegistryStore{
		ctx: ctx,
		cfg: cfg,
	}, nil
}

// GetRegistryClient checks the config and returns the suitaable registry.
func GetRegistryClient(ctx context.Context, cfg config.Config, registryClientCfg registry.ClientConfig) (registry.Client, error) {
	if cfg.ENSConfig.Override {
		ensStore, err := NewENSOverrideStore(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create ens override store: %v", err)
		}
		return registry.NewClientWithENSStore(ctx, registryClientCfg, ensStore)
	}
	return registry.NewClient(ctx, registryClientCfg)
}
