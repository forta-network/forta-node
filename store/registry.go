package store

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"time"

	mock_manifest2 "github.com/forta-network/forta-core-go/manifest/mocks"
	mock_registry "github.com/forta-network/forta-core-go/registry/mocks"
	"github.com/ipfs/go-cid"
	log "github.com/sirupsen/logrus"

	"github.com/forta-network/forta-core-go/ethereum"
	"github.com/forta-network/forta-core-go/feeds"
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
	GetAgentsIfChanged(scanner string) ([]*config.AgentConfig, bool, error)
	FindScannerShardIDForBot(agentID, scannerAddress string) (uint, uint, uint, error)
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
		botCfg, err := loadBot(rs.ctx, rs.cfg, rs.mc, bot.AgentID, bot.Manifest)
		switch {
		case err == nil: // yay
			// get sharding information
			shardID, shards, target, err := rs.FindScannerShardIDForBot(botCfg.ID, scanner)
			if err != nil {
				logger.WithError(err).Warn("could not find shard information for bot")
				return err
			}

			botCfg.ShardConfig = &config.ShardConfig{ShardID: shardID, Shards: shards, Target: target}
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

	if agt == nil {
		return nil, fmt.Errorf("agent does not exist: %s", agentID)
	}

	return loadBot(rs.ctx, rs.cfg, rs.mc, agentID, agt.Manifest)
}

func (rs *registryStore) FindScannerShardIDForBot(agentID, scannerAddress string) (
	shardID, shards, target uint, err error,
) {
	chainID := rs.cfg.ChainID

	// get manifest cid
	agt, err := rs.FindAgentGlobally(agentID)
	if err != nil {
		return shardID, shards, target, err
	}

	// fetch manifest
	agentManifest, err := rs.mc.GetAgentManifest(rs.ctx, agt.Manifest)
	if err != nil {
		return shardID, shards, target, err
	}

	// if bot is sharded, get total number of scanners by chain
	assigns, err := rs.rc.NumScannersForByChain(agentID, big.NewInt(int64(chainID)))
	if err != nil {
		return shardID, shards, target, fmt.Errorf("failed to get assign count: %v", err)
	}

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
		target = uint(assigns.Uint64())

		return 0, minShardCount, target, nil
	}

	// fallback for target, calculate it from shard to assign ratio.
	// target defaults to total assigns / shards
	if target == 0 && shards != 0 {
		target = uint(assigns.Uint64() / uint64(shards))
	}

	// get index of the scanner among scanners assigned to the bot for the same chain
	idx, err := rs.rc.IndexOfAssignedScannerByChain(agentID, scannerAddress, big.NewInt(int64(rs.cfg.ChainID)))
	if err != nil {
		return 0, minShardCount, target, fmt.Errorf("failed to get the index of scanner: %v, agentID: %s", err, agentID)
	}

	if idx == nil {
		return 0, minShardCount, target, fmt.Errorf("index for %s and %s not found", agentID, scannerAddress)
	}

	return calculateShardID(target, uint(idx.Uint64())), shards, target, nil
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

func loadBot(ctx context.Context, cfg config.Config, mc manifest.Client, agentID string, ref string) (*config.AgentConfig, error) {
	_, err := cid.Parse(ref)
	if len(ref) == 0 || err != nil {
		return nil, fmt.Errorf("%w: invalid bot cid '%s'", errInvalidBot, ref)
	}

	var agentData *manifest.SignedAgentManifest
	for i := 0; i < 10; i++ {
		agentData, err = mc.GetAgentManifest(ctx, ref)
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

	image, err := utils.ValidateDiscoImageRef(cfg.Registry.ContainerRegistry, *agentData.Manifest.ImageReference)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid bot image reference '%s': %v", errInvalidBot, *agentData.Manifest.ImageReference, err)
	}

	return &config.AgentConfig{
		ID:       agentID,
		Image:    image,
		Manifest: ref,
	}, nil
}

func NewRegistryStore(ctx context.Context, cfg config.Config, ethClient ethereum.Client, blockFeed feeds.BlockFeed) (*registryStore, error) {
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

	// make sure the registry client is refreshed and in sync.
	go func() {
		for err := range registry.ListenToUpgrades(ctx, rc, blockFeed) {
			if err != nil {
				log.WithError(err).Warn("error while listening to contract upgrades")
			}
		}
	}()

	return &registryStore{
		ctx: ctx,
		cfg: cfg,
		mc:  mc,
		rc:  rc,
	}, nil
}

func NewRegistryStoreFromClients(ctx context.Context, cfg config.Config, mc *mock_manifest2.MockClient, rc *mock_registry.MockClient) (*registryStore, error) {
	return &registryStore{rc: rc, mc: mc, ctx: ctx, cfg: cfg}, nil
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
