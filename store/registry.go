package store

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/contracts"
	"github.com/forta-network/forta-node/ethereum"
	"github.com/forta-network/forta-node/services/registry/regtypes"
	"github.com/forta-network/forta-node/utils"
	"math/big"
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
		length, _, err := rs.getAgentListHash(opts, scanner)
		if err != nil {
			return nil, false, err
		}

		scannerID := common.HexToHash(scanner).Big()

		var i int64
		for i = 0; i < length.Int64(); i++ {
			pos := big.NewInt(i)
			res, err := rs.dispatch.AgentRefAt(opts, scannerID, pos)
			if err != nil {
				return nil, false, err
			}
			agtCfg, err := rs.makeAgentConfig(common.Bytes2Hex(res.AgentId.Bytes()), res.Metadata)
			if err != nil {
				return nil, false, err
			}
			agts = append(agts, agtCfg)
		}

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

	image, ok := utils.ValidateImageRef(rs.cfg.Registry.ContainerRegistry, agentData.Manifest.ImageReference)
	if !ok {
		err = fmt.Errorf("invalid agent reference - skipping: %s", agentData.Manifest.ImageReference)
		return nil, err
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

func NewRegistryStore(ctx context.Context, cfg config.Config) (*registryStore, error) {
	agentRegAddress := config.GetEnvDefaults(cfg.Development).DefaultAgentRegistryContractAddress

	rpc, err := ethereum.NewRpcClient(cfg.Registry.Ethereum.JsonRpcUrl)
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

	ethClient, err := ethereum.NewStreamEthClient(context.Background(), cfg.Registry.Ethereum.JsonRpcUrl)
	if err != nil {
		return nil, err
	}

	var ipfsURL string
	if cfg.Registry.IPFSGateway != nil {
		ipfsURL = *cfg.Registry.IPFSGateway
	} else {
		ipfsURL = config.DefaultIPFSGateway
	}

	return &registryStore{
		ctx:        ctx,
		eth:        ethClient,
		cfg:        cfg,
		dispatch:   d,
		agents:     ar,
		ipfsClient: &ipfsClient{ipfsURL},
	}, nil
}
