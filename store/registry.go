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
)

var (
	errInvalidBot = errors.New("invalid bot")
)

const (
	keyDefaultChainSetting = "default"
	minShardCount          = 1
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

	shouldUpdate := rs.lastCompletedVersion != hash.Hash || time.Since(rs.lastUpdate) > 5*time.Minute
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

// returns shard id for an index, distributed evenly in an increased order.
// Example:
// Target: 6, Shards: 3
// should be [0,0,0,0,0,0,1,1,1,1,1,1,2,2,2,2,2,2]
func calculateShardID(target, idx uint) uint {
	if target == 0 {
		return 0
	}

	return idx / target
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
		ID:       agentID,
		Image:    image,
		Manifest: ref,
		ChainID:  cfg.ChainID,
		Owner:    owner,
	}, signedManifest, nil
}

func (rs *registryStore) loadAssignment(assignment *registry.Assignment) (*config.AgentConfig, error) {
	botCfg, agentData, err := loadBot(rs.ctx, rs.cfg, rs.bms, assignment.AgentID, assignment.AgentManifest, assignment.AgentOwner)
	if err != nil {
		return nil, err
	}

	botCfg.Owner = assignment.AgentOwner
	botCfg.ShardConfig = populateShardConfig(assignment, agentData, rs.cfg.ChainID)

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

func populateShardConfig(assignment *registry.Assignment, agentManifest *manifest.SignedAgentManifest, chainID int) *config.ShardConfig {
	var (
		target, shards uint
	)

	// check if there is a default chain setting
	chainSetting, ok := agentManifest.Manifest.ChainSettings[keyDefaultChainSetting]
	// if not a sharded bot, shard is always 0
	if ok {
		target = chainSetting.Target
		shards = chainSetting.Shards
	}

	// check if there is a chain setting for the scanner's chain
	chainIDStr := strconv.FormatInt(int64(chainID), 10)
	chainSetting, ok = agentManifest.Manifest.ChainSettings[chainIDStr]
	// if not a sharded bot, shard is always 0
	if ok {
		target = chainSetting.Target
		shards = chainSetting.Shards
	}

	// if no sharding specified, shard count is 1 and target is total assigns
	if shards == 0 {
		target = uint(assignment.AssignedScanners)

		return createShardConfig(0, minShardCount, target)
	}

	// fallback for target, calculate it from shard to assign ratio.
	// target defaults to total assigns / shards
	if target == 0 && shards != 0 {
		target = uint(uint64(assignment.AssignedScanners) / uint64(shards))
	}

	shardID := calculateShardID(target, uint(assignment.ScannerIndex))

	return createShardConfig(shardID, shards, target)
}

func createShardConfig(shardID, shards, target uint) *config.ShardConfig {
	return &config.ShardConfig{
		ShardID: shardID,
		Target:  target,
		Shards:  shards,
	}
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

	// load by image references
	for i, agentImage := range rs.cfg.LocalModeConfig.BotImages {
		if len(agentImage) == 0 {
			continue
		}
		// forta-agent-1, forta-agent-2, forta-agent-3, ...
		agentID := strconv.Itoa(i + 1)
		agentConfigs = append(agentConfigs, *rs.makePrivateModeAgentConfig(agentID, agentImage, nil))
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
	for i, shardedBot := range rs.cfg.LocalModeConfig.ShardedBots {
		// load bot by image
		if shardedBot.BotImage != nil {
			instances := shardedBot.Shards * shardedBot.Target
			for botIdx := uint(0); botIdx < instances; botIdx++ {
				shardConfig := &config.ShardConfig{
					Shards:  shardedBot.Shards,
					Target:  shardedBot.Target,
					ShardID: calculateShardID(shardedBot.Target, botIdx),
				}

				agentID := strconv.Itoa(len(agentConfigs) + i + 1)
				agentConfigs = append(
					agentConfigs, *rs.makePrivateModeAgentConfig(agentID, *shardedBot.BotImage, shardConfig),
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
	return nil, errors.New("feature not available (private/local registry)")
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
