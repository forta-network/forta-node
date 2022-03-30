package store

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/forta-protocol/forta-core-go/manifest"
	"github.com/forta-protocol/forta-core-go/registry"

	log "github.com/sirupsen/logrus"

	"github.com/forta-protocol/forta-core-go/ethereum"
	"github.com/forta-protocol/forta-core-go/utils"
	"github.com/forta-protocol/forta-node/config"
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

	lastUpdate time.Time
	version    string
	mu         sync.Mutex
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

	if rs.version != hash.Hash || time.Since(rs.lastUpdate) > 1*time.Hour {
		if err := rs.rc.PegLatestBlock(); err != nil {
			return nil, false, err
		}
		defer rs.rc.ResetOpts()
		var agts []*config.AgentConfig

		var failedLoadingAny bool
		err := rs.rc.ForEachAssignedAgent(scanner, func(a *registry.Agent) error {
			agtCfg, err := rs.makeAgentConfig(a.AgentID, a.Manifest)
			if err != nil {
				failedLoadingAny = true
				log.WithField("agentId", a.AgentID).WithError(err).Warn("could not parse config for agent")
				// ignore agent and move on by not returning the error
				return nil
			}
			agts = append(agts, agtCfg)
			return nil
		})

		if err != nil {
			return nil, false, err
		}

		// failed to load all: not doing this can cause getting stuck with the latest hash and zero agents
		if len(agts) == 0 && failedLoadingAny {
			return nil, false, errors.New("loaded zero agents")
		}

		rs.version = hash.Hash
		rs.lastUpdate = time.Now()
		return agts, true, nil
	}
	return nil, false, nil
}

func (rs *registryStore) FindAgentGlobally(agentID string) (*config.AgentConfig, error) {
	agt, err := rs.rc.GetAgent(agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get the latest ref: %v, agentID: %s", err, agentID)
	}
	return rs.makeAgentConfig(agentID, agt.Manifest)
}

func (rs *registryStore) makeAgentConfig(agentID string, ref string) (*config.AgentConfig, error) {
	if len(ref) == 0 {
		return nil, nil
	}
	var agentData *manifest.SignedAgentManifest

	var err error
	for i := 0; i < 10; i++ {
		agentData, err = rs.mc.GetAgentManifest(rs.ctx, ref)
		if err == nil {
			break
		}
	}
	if err != nil {
		err = fmt.Errorf("failed to load the agent file using ipfs ref: %v", err)
		return nil, err
	}

	if agentData.Manifest.ImageReference == nil {
		return nil, fmt.Errorf("invalid agent image reference, it is nil")
	}

	image, err := utils.ValidateDiscoImageRef(rs.cfg.Registry.ContainerRegistry, *agentData.Manifest.ImageReference)
	if err != nil {
		return nil, fmt.Errorf("invalid agent image reference '%s': %v", *agentData.Manifest.ImageReference, err)
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
