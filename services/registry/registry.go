package registry

import (
	"context"
	"fmt"
	"github.com/forta-network/forta-node/store"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/ethereum"
	"github.com/forta-network/forta-node/services/registry/regtypes"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

// RegistryService listens to the agent scanner list changes so the node can stay in sync.
type RegistryService struct {
	cfg            config.Config
	scannerAddress common.Address
	msgClient      clients.MessageClient

	rpcClient     *rpc.Client
	registryStore store.RegistryStore

	agentsConfigs []*config.AgentConfig
	done          chan struct{}
	version       string
	sem           *semaphore.Weighted
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
func New(cfg config.Config, scannerAddress common.Address, msgClient clients.MessageClient) *RegistryService {
	return &RegistryService{
		cfg:            cfg,
		scannerAddress: scannerAddress,
		msgClient:      msgClient,
		done:           make(chan struct{}),
	}
}

// Init only initializes the service.
func (rs *RegistryService) Init() error {
	regStr, err := store.NewRegistryStore(context.TODO(), rs.cfg)
	if err != nil {
		return err
	}
	rs.registryStore = regStr
	return nil
}

// Start initializes and starts the registry service.
func (rs *RegistryService) Start() error {
	log.Infof("Starting %s", rs.Name())
	if err := rs.Init(); err != nil {
		return err
	}
	rs.sem = semaphore.NewWeighted(1)
	return rs.start()
}

func (rs *RegistryService) start() error {
	go func() {
		//TODO: possibly make this configurable, but 15s per block is normal
		ticker := time.NewTicker(15 * time.Second)
		for {
			if err := rs.publishLatestAgents(); err != nil {
				log.Errorf("failed to publish the latest agents: %v", err)
			}
			<-ticker.C
		}
	}()

	return nil
}

func (rs *RegistryService) publishLatestAgents() error {
	// only allow one executor at a time, even if slow
	if rs.sem.TryAcquire(1) {
		defer rs.sem.Release(1)
		agts, changed, err := rs.registryStore.GetAgentsIfChanged(rs.scannerAddress.Hex())
		if err != nil {
			return fmt.Errorf("failed to get the scanner list agents version: %v", err)
		}
		if changed {
			log.WithField("count", len(agts)).Infof("publishing list of agents")
			rs.agentsConfigs = agts
			rs.msgClient.Publish(messaging.SubjectAgentsVersionsLatest, agts)
		} else {
			log.Info("registry: no agent changes detected")
		}
	}
	return nil
}

// Stop stops the registry service.
func (rs *RegistryService) Stop() error {
	return nil
}

// Name returns the name of the service.
func (rs *RegistryService) Name() string {
	return "RegistryService"
}
