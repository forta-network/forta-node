package store

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/ipfs/go-cid"
	log "github.com/sirupsen/logrus"

	"github.com/forta-network/forta-core-go/ens"
	"github.com/forta-network/forta-core-go/manifest"
	"github.com/forta-network/forta-core-go/registry"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/store/sharding"
)

var (
	errInvalidBot = errors.New("invalid bot")
	ErrLocalMode  = errors.New("feature not available (private/local registry)")
)

const (
	// This is the force reload interval that helps us ignore the on-chain assignment
	// list hash. This helps avoid getting stuck with bad state.
	//
	// WARNING: This also affects how fast the nodes react to shard ID changes
	// because the bot assignment hash may not change for a scanner when
	// the scanner list for a bot changes (i.e. when another scanner is unassigned).
	assignmentForceReloadInterval = time.Minute * 5
)

type RegistryStore interface {
	FindAgentGlobally(agentID string) (*config.AgentConfig, error)
	GetAgentsIfChanged(scanner string) ([]config.AgentConfig, bool, error)
}

type registryStore struct {
	ctx context.Context
	bms BotManifestStore
	rc  registry.Client
	cfg config.Config

	lastUpdate           time.Time
	lastCompletedVersion string
	loadedBots           []config.AgentConfig
	invalidAssignments   []*registry.Assignment
	mu                   sync.Mutex
}

func (rs *registryStore) GetAgentsIfChanged(scanner string) ([]config.AgentConfig, bool, error) {
	// because we peg the latest block, it can be problematic if this is called concurrently
	rs.mu.Lock()
	defer rs.mu.Unlock()

	hash, err := rs.rc.GetAssignmentHash(scanner)
	if err != nil {
		return nil, false, err
	}

	shouldUpdate := rs.lastCompletedVersion != hash.Hash ||
		time.Since(rs.lastUpdate) > assignmentForceReloadInterval
	if !shouldUpdate {
		return nil, false, nil
	}

	if err := rs.rc.PegLatestBlock(); err != nil {
		return nil, false, err
	}
	defer rs.rc.ResetOpts()

	var (
		loadedBots         []config.AgentConfig
		invalidAssignments []*registry.Assignment
		failedLoadingAny   bool
	)

	chainId := big.NewInt(int64(rs.cfg.ChainID))
	assignments, err := rs.rc.GetAssignmentList(nil, chainId, scanner)
	if err != nil {
		return nil, false, err
	}

	for _, assignment := range assignments {
		logger := log.WithField("botId", assignment.AgentID)

		// if already invalidated, remember it for next time
		if rs.isInvalidBot(assignment) {
			invalidAssignments = append(invalidAssignments, assignment)
			logger.Warn("invalid bot - skipping")
			continue
		}

		// try loading the rest of the unrecognized bots
		botCfg, err := rs.loadAssignment(assignment)
		switch {
		case err == nil: // yay
			// get sharding information
			loadedBots = append(loadedBots, *botCfg) // remember for next time
			logger.Info("successfully loaded bot")

		case errors.Is(err, errInvalidBot):
			invalidAssignments = append(invalidAssignments, assignment) // remember for next time
			logger.WithError(err).Warn("invalid bot - skipping")
		default:
			failedLoadingAny = true
			logger.WithError(err).Warn("could not load bot - skipping")
			// ignore agent and move on by not returning the error
			// it will not be recognized next time and will be retried above
			continue
		}
	}

	// failed to load all: forget that this attempt existed
	// not doing this can cause getting stuck with the latest hash and zero agents
	if len(loadedBots) == 0 && failedLoadingAny {
		return nil, false, errors.New("loaded zero bots")
	}

	// remember the bots and the update time next time
	rs.loadedBots = loadedBots
	rs.invalidAssignments = invalidAssignments
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

	botCfg, _, err := loadBot(rs.ctx, rs.cfg, rs.bms, agentID, agt.Manifest, agt.Owner)
	return botCfg, err
}

func (rs *registryStore) getLoadedBot(manifest string) (config.AgentConfig, bool) {
	for _, loadedBot := range rs.loadedBots {
		if manifest == loadedBot.Manifest {
			return loadedBot, true
		}
	}
	return config.AgentConfig{}, false
}

func (rs *registryStore) isInvalidBot(bot *registry.Assignment) bool {
	for _, invalidBot := range rs.invalidAssignments {
		if bot.AgentManifest == invalidBot.AgentManifest {
			return true
		}
	}
	return false
}

func loadBot(ctx context.Context, cfg config.Config, bms BotManifestStore, agentID string, ref string, owner string) (*config.AgentConfig, *manifest.SignedAgentManifest, error) {
	_, err := cid.Parse(ref)
	if len(ref) == 0 || err != nil {
		return nil, nil, fmt.Errorf("%w: invalid bot cid '%s'", errInvalidBot, ref)
	}

	signedManifest, err := bms.GetBotManifest(ctx, ref)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load the bot manifest: %v", err)
	}

	if signedManifest.Manifest.ImageReference == nil {
		return nil, nil, fmt.Errorf("%w: invalid bot image reference, it is nil", errInvalidBot)
	}

	image, err := utils.ValidateDiscoImageRef(
		cfg.Registry.ContainerRegistry, *signedManifest.Manifest.ImageReference,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: invalid bot image reference '%s': %v", errInvalidBot, *signedManifest.Manifest.ImageReference, err)
	}

	return &config.AgentConfig{
		ID:              agentID,
		Image:           image,
		Manifest:        ref,
		ChainID:         cfg.ChainID,
		Owner:           owner,
		ProtocolVersion: signedManifest.Manifest.ProtocolVersion,
	}, signedManifest, nil
}

func (rs *registryStore) loadAssignment(assignment *registry.Assignment) (*config.AgentConfig, error) {
	botCfg, agentData, err := loadBot(rs.ctx, rs.cfg, rs.bms, assignment.AgentID, assignment.AgentManifest, assignment.AgentOwner)
	if err != nil {
		return nil, err
	}

	botCfg.Owner = assignment.AgentOwner

	if botCfg.ProtocolVersion >= 2 {
		var ok bool
		botCfg.ShardConfig, ok = sharding.CalculateShardConfigV2(assignment, agentData)
		if !ok {
			return nil, fmt.Errorf("%w: invalid sharding config", errInvalidBot)
		}
		botCfg.ChainID = int(botCfg.ShardConfig.ChainID)
	} else {
		botCfg.ShardConfig = sharding.CalculateShardConfig(assignment, agentData, rs.cfg.ChainID)
	}

	return botCfg, nil
}

func NewRegistryStore(ctx context.Context, cfg config.Config) (*registryStore, error) {
	mc, err := manifest.NewClient(cfg.Registry.IPFS.GatewayURL)
	if err != nil {
		return nil, err
	}
	bms := NewBotManifestStore(mc)

	rc, err := GetRegistryClient(
		ctx, cfg, registry.ClientConfig{
			JsonRpcUrl:       cfg.Registry.JsonRpc.Url,
			ENSAddress:       cfg.ENSConfig.ContractAddress,
			Name:             "registry-store",
			MulticallAddress: cfg.AdvancedConfig.MulticallAddress,
		},
	)
	if err != nil {
		return nil, err
	}

	// make sure the registry client is refreshed and in sync.
	go func() {
		ticker := time.NewTicker(time.Minute * 15)
		for {
			select {
			case <-ticker.C:
				err := rc.RefreshContracts()
				if err != nil {
					log.WithError(err).Warn("error while refreshing the registry contracts")
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	return &registryStore{
		ctx: ctx,
		cfg: cfg,
		bms: bms,
		rc:  rc,
	}, nil
}

type privateRegistryStore struct {
	ctx context.Context
	cfg config.Config
	rc  registry.Client
	bms BotManifestStore
	mu  sync.Mutex
}

func (rs *privateRegistryStore) GetAgentsIfChanged(scanner string) ([]config.AgentConfig, bool, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	var agentConfigs []config.AgentConfig

	var agentID int

	for _, bot := range rs.cfg.LocalModeConfig.Bots {
		if bot.BotImage == nil {
			continue
		}

		agentID++
		agentConfig := rs.makePrivateModeAgentConfig(strconv.Itoa(agentID), *bot.BotImage, nil)
		if bot.ProtocolVersion != nil {
			agentConfig.ProtocolVersion = *bot.ProtocolVersion
		}

		agentConfigs = append(agentConfigs, *agentConfig)
	}

	// load by image references
	for _, agentImage := range rs.cfg.LocalModeConfig.BotImages {
		if len(agentImage) == 0 {
			continue
		}
		// forta-agent-1, forta-agent-2, forta-agent-3, ...
		agentID++
		agentConfigs = append(agentConfigs, *rs.makePrivateModeAgentConfig(strconv.Itoa(agentID), agentImage, nil))
	}

	// load by bot IDs
	for _, agentID := range rs.cfg.LocalModeConfig.BotIDs {
		agt, err := rs.rc.GetAgent(agentID)
		logger := log.WithFields(log.Fields{
			"botID": agentID,
		})
		if err != nil {
			logger.WithError(err).Error("failed to get bot from registry")
			continue
		}
		agtCfg, _, err := loadBot(rs.ctx, rs.cfg, rs.bms, agentID, agt.Manifest, agt.Owner)
		if err != nil {
			logger.WithError(err).Error("failed to load bot")
			continue
		}

		agtCfg.Owner = agt.Owner
		agentConfigs = append(agentConfigs, *agtCfg)
	}

	// load sharded bots by image
	for _, shardedBot := range rs.cfg.LocalModeConfig.ShardedBots {
		// load bot by image
		if shardedBot.BotImage != nil {
			instances := shardedBot.Shards * shardedBot.Target
			for botIdx := uint(0); botIdx < instances; botIdx++ {
				shardConfig := &config.ShardConfig{
					Shards:  shardedBot.Shards,
					Target:  shardedBot.Target,
					ShardID: sharding.CalculateShardID(shardedBot.Target, botIdx),
				}

				agentID++
				agentConfigs = append(
					agentConfigs, *rs.makePrivateModeAgentConfig(strconv.Itoa(agentID), *shardedBot.BotImage, shardConfig),
				)
			}
		}
	}

	// load the standalone bot configs that are already running
	if rs.cfg.LocalModeConfig.IsStandalone() {
		for _, runningBot := range rs.cfg.LocalModeConfig.Standalone.BotContainers {
			agentConfigs = append(agentConfigs, config.AgentConfig{
				ID:           runningBot,
				IsStandalone: true,
				ChainID:      rs.cfg.ChainID,
			})
		}
	}

	return agentConfigs, true, nil
}

func (rs *privateRegistryStore) FindAgentGlobally(agentID string) (*config.AgentConfig, error) {
	return nil, ErrLocalMode
}

func (rs *privateRegistryStore) makePrivateModeAgentConfig(
	id string, image string,
	shardConfig *config.ShardConfig,
) *config.AgentConfig {
	return &config.AgentConfig{
		ID:          id,
		Image:       image,
		IsLocal:     true,
		ShardConfig: shardConfig,
		ChainID:     rs.cfg.ChainID,
	}
}

func NewPrivateRegistryStore(ctx context.Context, cfg config.Config) (*privateRegistryStore, error) {
	mc, err := manifest.NewClient(cfg.Registry.IPFS.GatewayURL)
	if err != nil {
		return nil, err
	}
	bms := NewBotManifestStore(mc)

	rc, err := GetRegistryClient(ctx, cfg, registry.ClientConfig{
		JsonRpcUrl: cfg.Registry.JsonRpc.Url,
		ENSAddress: cfg.ENSConfig.ContractAddress,
		Name:       "registry-store",
		NoRefresh:  cfg.LocalModeConfig.IsStandalone(),
	})
	if err != nil {
		return nil, err
	}
	return &privateRegistryStore{
		ctx: ctx,
		cfg: cfg,
		bms: bms,
		rc:  rc,
	}, nil
}

// GetRegistryClient checks the config and returns the suitaable registry.
func GetRegistryClient(ctx context.Context, cfg config.Config, registryClientCfg registry.ClientConfig) (registry.Client, error) {
	if cfg.ENSConfig.Override {
		ensResolver, err := NewENSOverrideResolver(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create ens override resolver: %v", err)
		}
		ensStore := ens.NewENStoreWithResolver(ensResolver)
		return registry.NewClientWithENSStore(ctx, registryClientCfg, ensStore)
	}
	return registry.NewClient(ctx, registryClientCfg)
}
