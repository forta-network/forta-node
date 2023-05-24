package registry

import (
	"context"
	"fmt"
	"time"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/metrics"
	"github.com/forta-network/forta-node/store"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/ethereum"
	"github.com/forta-network/forta-core-go/feeds"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/services/registry/regtypes"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

// RegistryService listens to the agent scanner list changes so the node can stay in sync.
type RegistryService struct {
	cfg            config.Config
	scannerAddress common.Address
	msgClient      clients.MessageClient
	ethClient      ethereum.Client
	blockFeed      feeds.BlockFeed

	rpcClient     *rpc.Client
	registryStore store.RegistryStore

	agentsConfigs []*config.AgentConfig
	done          chan struct{}
	version       string
	sem           *semaphore.Weighted

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
func New(cfg config.Config, scannerAddress common.Address, msgClient clients.MessageClient, ethClient ethereum.Client, blockFeed feeds.BlockFeed) *RegistryService {
	return &RegistryService{
		cfg:            cfg,
		scannerAddress: scannerAddress,
		msgClient:      msgClient,
		ethClient:      ethClient,
		done:           make(chan struct{}),
		blockFeed:      blockFeed,
	}
}

// Init only initializes the service.
func (rs *RegistryService) Init() error {
	var (
		regStr store.RegistryStore
		err    error
	)
	if rs.cfg.LocalModeConfig.Enable {
		regStr, err = store.NewPrivateRegistryStore(context.Background(), rs.cfg)
	} else {
		regStr, err = store.NewRegistryStore(context.Background(), rs.cfg, rs.ethClient, rs.blockFeed)
	}
	if err != nil {
		return err
	}
	rs.registryStore = regStr
	return nil
}

// Start initializes and starts the registry service.
func (rs *RegistryService) Start() error {
	if err := rs.Init(); err != nil {
		return err
	}
	rs.sem = semaphore.NewWeighted(1)
	return rs.start()
}

func (rs *RegistryService) start() error {
	go func() {
		ticker := time.NewTicker(time.Duration(rs.cfg.Registry.CheckIntervalSeconds) * time.Second)
		for {
			err := rs.publishLatestAgents()
			rs.lastErr.Set(err)
			if err != nil {
				log.WithError(err).Error("failed to publish the latest agents")
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
		rs.lastChecked.Set()
		agts, changed, err := rs.registryStore.GetAgentsIfChanged(rs.scannerAddress.Hex())
		if err != nil {
			return fmt.Errorf("failed to get the scanner list agents version: %v", err)
		}

		if changed {
			rs.lastChangeDetected.Set()
			log.WithField("count", len(agts)).Infof("publishing list of agents")
			rs.agentsConfigs = agts

			// emit metrics for each detected bot
			var ms []*protocol.AgentMetric
			for _, agt := range agts {
				ms = append(ms, metrics.CreateAgentMetric(agt.ID, metrics.MetricAgentRegistryDetected, 1))
			}

			metrics.SendAgentMetrics(rs.msgClient, ms)

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
	return "registry"
}

// Health implements the health.Reporter interface.
func (rs *RegistryService) Health() health.Reports {
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
