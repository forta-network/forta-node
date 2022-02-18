package store

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	log "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/contracts"
	"github.com/forta-protocol/forta-node/ethereum"
	"github.com/forta-protocol/forta-node/services/registry/regtypes"
	"github.com/forta-protocol/forta-node/utils"
)

type RegistryStore interface {
	FindAgentGlobally(agentID string) (*config.AgentConfig, error)
	GetAgentsIfChanged(scanner string) ([]*config.AgentConfig, bool, error)
}

type registryStore struct {
	ctx        context.Context
	eth        ethereum.Client
	dispatch   dispatch
	agents     agentRegistry
	ipfsClient IPFSClient
	cfg        config.Config
	version    string
}

func (rs *registryStore) GetAgentsIfChanged(scanner string) ([]*config.AgentConfig, bool, error) {
	_, versionHash, err := rs.getAgentListHash(nil, scanner)
	if err != nil {
		return nil, false, err
	}
	if rs.version != versionHash {
		var agts []*config.AgentConfig
		opts, err := rs.optsWithLatestBlock()
		if err != nil {
			return nil, false, err
		}
		// get it again, so that the block is fixed throughout iteration
		// this avoids getting the opts (getBlock) in the normal nothing-changed case
		length, hash, err := rs.getAgentListHash(opts, scanner)
		if err != nil {
			return nil, false, err
		}

		scannerID := common.HexToHash(scanner).Big()

		var i int64
		var failedLoadingAny bool
		for i = 0; i < length.Int64(); i++ {
			pos := big.NewInt(i)
			res, err := rs.dispatch.AgentRefAt(opts, scannerID, pos)
			if err != nil {
				return nil, false, err
			}
			agtCfg, err := rs.makeAgentConfig(utils.BytesToHex(res.AgentId.Bytes()), res.Metadata)
			if err != nil {
				failedLoadingAny = true
				log.WithError(err).Warn("could not parse config for agent")
			} else {
				agts = append(agts, agtCfg)
			}
		}

		// failed to load all: not doing this can cause getting stuck with the latest hash and zero agents
		if len(agts) == 0 && failedLoadingAny {
			return nil, false, errors.New("loaded zero agents")
		}

		rs.version = hash
		return agts, true, nil
	}
	return nil, false, nil
}

func (rs *registryStore) getAgentListHash(opts *bind.CallOpts, scanner string) (*big.Int, string, error) {
	res, err := rs.dispatch.ScannerHash(opts, common.HexToHash(scanner).Big())
	if err != nil {
		return nil, "", err
	}

	return res.Length, utils.Bytes32ToHex(res.Manifest), nil
}

func (rs *registryStore) FindAgentGlobally(agentID string) (*config.AgentConfig, error) {
	opts, err := rs.optsWithLatestBlock()
	if err != nil {
		return nil, err
	}
	agt, err := rs.agents.GetAgent(opts, common.HexToHash(agentID).Big())
	if err != nil {
		return nil, fmt.Errorf("failed to get the latest ref: %v", err)
	}
	return rs.makeAgentConfig(agentID, agt.Metadata)
}

func (rs *registryStore) makeAgentConfig(agentID string, ref string) (*config.AgentConfig, error) {
	if len(ref) == 0 {
		return nil, nil
	}
	var agentData *regtypes.AgentFile

	var err error
	for i := 0; i < 10; i++ {
		agentData, err = rs.ipfsClient.GetAgentFile(ref)
		if err == nil {
			break
		}
	}
	if err != nil {
		err = fmt.Errorf("failed to load the agent file using ipfs ref: %v", err)
		return nil, err
	}

	image, err := utils.ValidateDiscoImageRef(rs.cfg.Registry.ContainerRegistry, agentData.Manifest.ImageReference)
	if err != nil {
		return nil, fmt.Errorf("invalid agent image reference '%s': %v", agentData.Manifest.ImageReference, err)
	}

	return &config.AgentConfig{
		ID:       agentID,
		Image:    image,
		Manifest: ref,
	}, nil
}

func (rs *registryStore) optsWithLatestBlock() (*bind.CallOpts, error) {
	blk, err := rs.eth.BlockByNumber(rs.ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get the block for agents: %v", err)
	}
	num, err := utils.HexToBigInt(blk.Number)
	if err != nil {
		return nil, err
	}
	return &bind.CallOpts{
		BlockNumber: num,
	}, nil
}

func NewRegistryStore(ctx context.Context, cfg config.Config, ethClient ethereum.Client) (*registryStore, error) {
	agentRegAddress := cfg.AgentRegistryContractAddress

	rpc, err := ethereum.NewRpcClient(cfg.Registry.JsonRpc.Url)
	if err != nil {
		return nil, err
	}
	client := ethclient.NewClient(rpc)

	ar, err := contracts.NewAgentRegistryCaller(common.HexToAddress(agentRegAddress), client)
	if err != nil {
		return nil, err
	}

	d, err := contracts.NewDispatchCaller(common.HexToAddress(cfg.Registry.ContractAddress), client)
	if err != nil {
		return nil, err
	}

	return &registryStore{
		ctx:        ctx,
		eth:        ethClient,
		cfg:        cfg,
		dispatch:   d,
		agents:     ar,
		ipfsClient: &ipfsClient{cfg.Registry.IPFS.GatewayURL},
	}, nil
}
