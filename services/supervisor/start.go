package supervisor

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/forta-protocol/forta-core-go/manifest"
	"github.com/forta-protocol/forta-core-go/release"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ipfs/go-cid"
	log "github.com/sirupsen/logrus"

	"github.com/forta-protocol/forta-core-go/clients/agentlogs"
	"github.com/forta-protocol/forta-core-go/clients/health"
	"github.com/forta-protocol/forta-core-go/security"
	"github.com/forta-protocol/forta-node/clients"
	"github.com/forta-protocol/forta-node/clients/messaging"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/services"
)

const (
	// SupervisorStrategyVersion is for versioning the critical changes in supervisor's management strategy.
	// It's effective in deciding if an agent container should be restarted or not.
	SupervisorStrategyVersion = "1"
)

// SupervisorService manages the scanner node's service and agent containers.
type SupervisorService struct {
	ctx            context.Context
	client         clients.DockerClient
	globalClient   clients.DockerClient
	manifestClient manifest.Client
	releaseClient  release.Client

	msgClient   clients.MessageClient
	config      SupervisorServiceConfig
	maxLogSize  string
	maxLogFiles int

	scannerContainer *clients.DockerContainer
	jsonRpcContainer *clients.DockerContainer
	containers       []*Container
	mu               sync.RWMutex

	lastRun                   health.TimeTracker
	lastStop                  health.TimeTracker
	lastTelemetryRequest      health.TimeTracker
	lastTelemetryRequestError health.ErrorTracker
	lastAgentLogsRequest      health.TimeTracker
	lastAgentLogsRequestError health.ErrorTracker

	healthClient health.HealthClient

	agentLogsClient agentlogs.Client
	prevAgentLogs   agentlogs.Agents
}

type SupervisorServiceConfig struct {
	Config     config.Config
	Passphrase string
	Key        *keystore.Key
}

// Container extends the default container data.
type Container struct {
	clients.DockerContainer
	IsAgent     bool
	AgentConfig *config.AgentConfig
}

func (sup *SupervisorService) Start() error {
	if err := sup.start(); err != nil {
		return err
	}

	go sup.healthCheck()

	return nil
}

func (sup *SupervisorService) start() error {
	if !sup.config.Config.TelemetryConfig.Disable {
		go sup.syncTelemetryData()
	}
	go sup.syncAgentLogs()

	sup.mu.Lock()
	defer sup.mu.Unlock()

	log.Infof("Starting %s", sup.Name())
	_, err := log.ParseLevel(sup.config.Config.Log.Level)
	if err != nil {
		log.Error("invalid log level", err)
		return err
	}

	hostFortaDir := os.Getenv(config.EnvHostFortaDir)
	if len(hostFortaDir) == 0 {
		return fmt.Errorf("supervisor needs to know $%s to mount to the other containers it runs", config.EnvHostFortaDir)
	}
	releaseInfo := release.ReleaseInfoFromString(os.Getenv(config.EnvReleaseInfo))
	releaseInfo, err = sup.getFullReleaseInfo(releaseInfo)
	if err != nil {
		return fmt.Errorf("failed to get full release info: %v", err)
	}
	if releaseInfo != nil {
		release.LogReleaseInfo(releaseInfo)
	}

	sup.maxLogSize = sup.config.Config.Log.MaxLogSize
	sup.maxLogFiles = sup.config.Config.Log.MaxLogFiles

	if err := sup.removeOldContainers(); err != nil {
		return err
	}

	if err := sup.ensureNodeImages(); err != nil {
		return err
	}

	supervisorContainer, err := sup.globalClient.GetContainerByName(sup.ctx, config.DockerSupervisorContainerName)
	if err != nil {
		return fmt.Errorf("failed to get the supervisor container: %v", err)
	}
	commonNodeImage := supervisorContainer.Image

	nodeNetworkID, err := sup.client.CreatePublicNetwork(sup.ctx, config.DockerNetworkName)
	if err != nil {
		return err
	}
	if err := sup.client.AttachNetwork(sup.ctx, supervisorContainer.ID, nodeNetworkID); err != nil {
		return fmt.Errorf("failed to attach supervisor container to node network: %v", err)
	}

	var natsNetworkID string
	if sup.config.Config.ExposeNats {
		natsNetworkID = nodeNetworkID
	} else {
		natsNetworkID, err = sup.client.CreateInternalNetwork(sup.ctx, config.DockerNatsContainerName)
		if err != nil {
			return err
		}
		if err := sup.client.AttachNetwork(sup.ctx, supervisorContainer.ID, natsNetworkID); err != nil {
			return fmt.Errorf("failed to attach supervisor container to nats network: %v", err)
		}
	}

	// start nats, wait for it and connect from the supervisor
	natsContainer, err := sup.client.StartContainer(sup.ctx, clients.DockerContainerConfig{
		Name:  config.DockerNatsContainerName,
		Image: "nats:2.3.2",
		Ports: map[string]string{
			"4222": "4222",
			"6222": "6222",
			"8222": "8222",
		},
		NetworkID:   natsNetworkID,
		MaxLogFiles: sup.maxLogFiles,
		MaxLogSize:  sup.maxLogSize,
	})
	if err != nil {
		return err
	}
	sup.addContainerUnsafe(natsContainer)

	if err := sup.client.WaitContainerStart(sup.ctx, natsContainer.ID); err != nil {
		return fmt.Errorf("failed while waiting for nats to start: %v", err)
	}
	// in tests, this is already set to a mock client
	if sup.msgClient == nil {
		sup.msgClient = messaging.NewClient("supervisor", fmt.Sprintf("%s:%s", config.DockerNatsContainerName, config.DefaultNatsPort))
	}
	sup.registerMessageHandlers()

	publisherContainer, err := sup.client.StartContainer(sup.ctx, clients.DockerContainerConfig{
		Name:  config.DockerPublisherContainerName,
		Image: commonNodeImage,
		Cmd:   []string{config.DefaultFortaNodeBinaryPath, "publisher"},
		Env: map[string]string{
			config.EnvReleaseInfo: releaseInfo.String(),
		},
		Volumes: map[string]string{
			hostFortaDir: config.DefaultContainerFortaDirPath,
		},
		Ports: map[string]string{
			"": config.DefaultHealthPort, // random host port
		},
		Files: map[string][]byte{
			"passphrase": []byte(sup.config.Passphrase),
		},
		NetworkID:   nodeNetworkID,
		MaxLogFiles: sup.maxLogFiles,
		MaxLogSize:  sup.maxLogSize,
	})
	if err != nil {
		return err
	}
	sup.addContainerUnsafe(publisherContainer)

	sup.jsonRpcContainer, err = sup.client.StartContainer(sup.ctx, clients.DockerContainerConfig{
		Name:  config.DockerJSONRPCProxyContainerName,
		Image: commonNodeImage,
		Cmd:   []string{config.DefaultFortaNodeBinaryPath, "json-rpc"},
		Volumes: map[string]string{
			// give access to host docker
			"/var/run/docker.sock": "/var/run/docker.sock",
			hostFortaDir:           config.DefaultContainerFortaDirPath,
		},
		Ports: map[string]string{
			"": config.DefaultHealthPort, // random host port
		},
		DialHost:    true,
		NetworkID:   nodeNetworkID,
		MaxLogFiles: sup.maxLogFiles,
		MaxLogSize:  sup.maxLogSize,
	})
	if err != nil {
		return err
	}
	sup.addContainerUnsafe(sup.jsonRpcContainer)

	sup.scannerContainer, err = sup.client.StartContainer(sup.ctx, clients.DockerContainerConfig{
		Name:  config.DockerScannerContainerName,
		Image: commonNodeImage,
		Cmd:   []string{config.DefaultFortaNodeBinaryPath, "scanner"},
		Volumes: map[string]string{
			hostFortaDir: config.DefaultContainerFortaDirPath,
		},
		Ports: map[string]string{
			"": config.DefaultHealthPort, // random host port
		},
		Files: map[string][]byte{
			"passphrase": []byte(sup.config.Passphrase),
		},
		DialHost:    true,
		NetworkID:   nodeNetworkID,
		MaxLogFiles: sup.maxLogFiles,
		MaxLogSize:  sup.maxLogSize,
	})
	if err != nil {
		return err
	}
	sup.addContainerUnsafe(sup.scannerContainer)

	if !sup.config.Config.ExposeNats {
		if err := sup.attachToNetwork(config.DockerPublisherContainerName, natsNetworkID); err != nil {
			return err
		}
		if err := sup.attachToNetwork(config.DockerScannerContainerName, natsNetworkID); err != nil {
			return err
		}
		if err := sup.attachToNetwork(config.DockerJSONRPCProxyContainerName, natsNetworkID); err != nil {
			return err
		}
	}

	return nil
}

func (sup *SupervisorService) attachToNetwork(containerName, nodeNetworkID string) error {
	container, err := sup.client.GetContainerByName(sup.ctx, containerName)
	if err != nil {
		return fmt.Errorf("failed to get '%s' container while attaching to node network: %v", containerName, err)
	}
	if err := sup.client.AttachNetwork(sup.ctx, container.ID, nodeNetworkID); err != nil {
		return fmt.Errorf("failed to attach '%s' container to node network: %v", containerName, err)
	}
	return nil
}

func (sup *SupervisorService) ensureNodeImages() error {
	for _, image := range []struct {
		Name string
		Ref  string
	}{
		{
			Name: "nats",
			Ref:  "nats:2.3.2",
		},
	} {
		if err := sup.client.EnsureLocalImage(sup.ctx, image.Name, image.Ref); err != nil {
			return err
		}
	}
	return nil
}

// removes old service containers and agents started with an old supervisor
func (sup *SupervisorService) removeOldContainers() error {
	type containerDefinition struct {
		ID   string
		Name string
	}
	var containersToRemove []*containerDefinition

	// gather old service containers
	for _, containerName := range []string{
		config.DockerScannerContainerName,
		config.DockerPublisherContainerName,
		config.DockerJSONRPCProxyContainerName,
		config.DockerNatsContainerName,
	} {
		container, err := sup.client.GetContainerByName(sup.ctx, containerName)
		if err != nil {
			log.WithError(err).WithField("containerName", containerName).Info("did not find old service container - ignoring")
			continue
		}
		containersToRemove = append(containersToRemove, &containerDefinition{
			ID:   container.ID,
			Name: containerName,
		})
	}

	// gather old agents
	containers, err := sup.client.GetContainers(sup.ctx)
	if err != nil {
		return fmt.Errorf("failed to get containers list: %v", err)
	}
	for _, container := range containers {
		containerName := container.Names[0][1:]
		logger := log.WithFields(log.Fields{
			"containerName": containerName,
			"containerId":   container.ID,
		})
		if !strings.Contains(containerName, "forta-agent-") {
			continue
		}
		if container.Labels[clients.DockerLabelFortaSupervisorStrategyVersion] != SupervisorStrategyVersion {
			logger.Info("agent container is old - need to remove")
			containersToRemove = append(containersToRemove, &containerDefinition{
				ID:   container.ID,
				Name: containerName,
			})
		}
	}

	// remove all of the gathered containers
	for _, container := range containersToRemove {
		logger := log.WithFields(log.Fields{
			"containerName": container.Name,
			"containerId":   container.ID,
		})
		if err := sup.client.RemoveContainer(sup.ctx, container.ID); err != nil {
			const msg = "failed to remove old container"
			logger.WithError(err).Error(msg)
			return fmt.Errorf("%s: %v", msg, err)
		}
		if err := sup.client.WaitContainerPrune(sup.ctx, container.ID); err != nil {
			const msg = "failed while waiting removal of old container"
			logger.WithError(err).Error(msg)
			return fmt.Errorf("%s: %v", msg, err)
		}
	}
	// after all gathered containers are removed, remove their networks
	for _, container := range containersToRemove {
		logger := log.WithFields(log.Fields{
			"containerName": container.Name,
			"containerId":   container.ID,
		})
		if err := sup.client.RemoveNetworkByName(sup.ctx, container.Name); err != nil {
			const msg = "failed to remove old network"
			logger.WithError(err).Warn(msg)
			// ignore network removal errs
		}
	}

	return nil
}

func (sup *SupervisorService) syncTelemetryData() {
	time.After(time.Second * 15)          // rate limit crash loops
	ticker := time.NewTicker(time.Minute) // slow down with auto-upgrade later
	for {
		err := sup.doSyncTelemetryData()
		sup.lastTelemetryRequest.Set()
		sup.lastTelemetryRequestError.Set(err)
		if err != nil {
			log.WithError(err).Warn("telemetry sync failed")
		}
		<-ticker.C
	}
}

func (sup *SupervisorService) doSyncTelemetryData() error {
	scannerJwt, err := security.CreateScannerJWT(sup.config.Key, map[string]interface{}{
		"access": "telemetry",
	})
	if err != nil {
		return err
	}
	return sup.healthClient.SendReports(
		fmt.Sprintf("http://host.docker.internal:%s/health", config.DefaultHealthPort),
		sup.config.Config.TelemetryConfig.URL,
		scannerJwt,
	)
}

// complete release info in case runner is old and starts supervisor by providing missing release properties
func (sup *SupervisorService) getFullReleaseInfo(releaseInfo *release.ReleaseInfo) (*release.ReleaseInfo, error) {
	if releaseInfo == nil {
		return nil, nil
	}
	if len(releaseInfo.IPFS) == 0 {
		return releaseInfo, nil
	}
	if _, err := cid.Parse(releaseInfo.IPFS); err != nil {
		return releaseInfo, nil
	}
	fullReleaseManifest, err := sup.releaseClient.GetReleaseManifest(sup.ctx, releaseInfo.IPFS)
	if err != nil {
		return nil, err
	}
	return &release.ReleaseInfo{
		FromBuild: false,
		IPFS:      releaseInfo.IPFS,
		Manifest:  *fullReleaseManifest,
	}, nil
}

func (sup *SupervisorService) Stop() error {
	sup.mu.RLock()
	defer sup.mu.RUnlock()

	ctx := context.Background()
	for _, cnt := range sup.containers {
		if services.IsGracefulShutdown() && cnt.IsAgent {
			continue // keep container agents alive
		}
		if err := sup.client.StopContainer(ctx, cnt.DockerContainer.ID); err != nil {
			log.Error(fmt.Sprintf("error stopping %s container", cnt.DockerContainer.ID), err)
		} else {
			log.Infof("Container %s is stopped", cnt.DockerContainer.ID)
		}
	}
	return nil
}

func (sup *SupervisorService) Name() string {
	return "supervisor"
}

// Health implements the health.Reporter interface.
func (sup *SupervisorService) Health() health.Reports {
	sup.mu.RLock()
	defer sup.mu.RUnlock()

	containersStatus := health.StatusOK
	if len(sup.containers) < 4 {
		containersStatus = health.StatusFailing
	}

	return health.Reports{
		&health.Report{
			Name:    "containers.managed",
			Status:  containersStatus,
			Details: strconv.Itoa(len(sup.containers)),
		},
		&health.Report{
			Name:    "event.run-agent.time",
			Status:  health.StatusInfo,
			Details: sup.lastRun.String(),
		},
		&health.Report{
			Name:    "event.stop-agent.time",
			Status:  health.StatusInfo,
			Details: sup.lastStop.String(),
		},
		sup.lastTelemetryRequest.GetReport("event.telemetry-sync.time"),
		sup.lastTelemetryRequestError.GetReport("event.telemetry-sync.error"),
		sup.lastAgentLogsRequest.GetReport("event.agent-logs-sync.time"),
		sup.lastAgentLogsRequestError.GetReport("event.agent-logs-sync.error"),
	}
}

func NewSupervisorService(ctx context.Context, cfg SupervisorServiceConfig) (*SupervisorService, error) {
	dockerClient, err := clients.NewDockerClient("supervisor")
	if err != nil {
		return nil, fmt.Errorf("failed to create the docker client: %v", err)
	}
	globalClient, err := clients.NewDockerClient("")
	if err != nil {
		return nil, fmt.Errorf("failed to create the global docker client: %v", err)
	}

	releaseClient, err := release.NewClient(cfg.Config.Registry.IPFS.GatewayURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create the release client: %v", err)
	}

	return &SupervisorService{
		ctx:             ctx,
		client:          dockerClient,
		globalClient:    globalClient,
		releaseClient:   releaseClient,
		config:          cfg,
		healthClient:    health.NewClient(),
		agentLogsClient: agentlogs.NewClient(cfg.Config.AgentLogsConfig.URL),
	}, nil
}
