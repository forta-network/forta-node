package supervisor

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/forta-network/forta-core-go/manifest"
	"github.com/forta-network/forta-core-go/release"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ipfs/go-cid"
	log "github.com/sirupsen/logrus"

	"github.com/forta-network/forta-core-go/clients/agentlogs"
	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/security"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/services"
	netmgmt "github.com/forta-network/forta-node/services/network"
)

const (
	// SupervisorStrategyVersion is for versioning the critical changes in supervisor's management strategy.
	// It's effective in deciding if an agent container should be restarted or not.
	SupervisorStrategyVersion = "1"
)

// SupervisorService manages the scanner node's service and agent containers.
type SupervisorService struct {
	ctx context.Context

	client           clients.DockerClient
	globalClient     clients.DockerClient
	agentImageClient clients.DockerClient
	botManager       netmgmt.BotManager

	manifestClient manifest.Client
	releaseClient  release.Client

	msgClient   clients.MessageClient
	config      SupervisorServiceConfig
	maxLogSize  string
	maxLogFiles int

	commonNodeImage string

	botNetworkID     string
	serviceNetworkID string

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
	IsAgent      bool
	IsAgentAdmin bool
	AgentConfig  *config.AgentConfig
}

func (sup *SupervisorService) Start() error {
	if err := sup.start(); err != nil {
		return err
	}

	go sup.healthCheck()

	return nil
}

func (sup *SupervisorService) start() error {
	// in addition to the feature disable flags, check private mode flags to disable agent logging and telemetry

	shouldDisableTelemetry := sup.config.Config.TelemetryConfig.Disable || sup.config.Config.PrivateModeConfig.Enable
	if !shouldDisableTelemetry {
		go sup.syncTelemetryData()
	}

	shouldDisableAgentLogs := sup.config.Config.AgentLogsConfig.Disable || sup.config.Config.PrivateModeConfig.Enable
	if !shouldDisableAgentLogs {
		go sup.syncAgentLogs()
	}

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
	sup.commonNodeImage = supervisorContainer.Image

	// we run an ephemeral container to detect host networking information
	hostNetContainer, err := sup.client.StartContainer(sup.ctx, clients.DockerContainerConfig{
		Name:      config.DockerHostNetContainerName,
		Image:     sup.commonNodeImage,
		Cmd:       []string{config.DefaultFortaNodeBinaryPath, "detect-host-networking"},
		NetworkID: "host", // attach to host networking so we can detect it
	})
	if err != nil {
		return err
	}
	defaultInterface, defaultSubnet, defaultGateway, err := sup.getHostNetworkingInfo(hostNetContainer)
	if err != nil {
		return fmt.Errorf("failed to get host networking info: %v", err)
	}

	// select x.x.x.128-x.x.x.255 as the ip range
	subnetParts := strings.Split(defaultSubnet, "/")
	ipAddr := subnetParts[0]
	ipAddrParts := strings.Split(ipAddr, ".")
	ipAddrParts[3] = "128"
	ipRange := strings.Join([]string{strings.Join(ipAddrParts, "."), subnetParts[1]}, "/")

	// create an ipvlan network that uses host interface as parent so we can provide internet access
	// to bots with easy isolation
	botNetworkID, err := sup.client.CreateNetwork(sup.ctx, config.DockerBotNetworkName, &types.NetworkCreate{
		Driver: "ipvlan",
		IPAM: &network.IPAM{
			Config: []network.IPAMConfig{
				{
					Subnet:  defaultSubnet,
					IPRange: ipRange,
					Gateway: defaultGateway,
				},
			},
		},
		Options: map[string]string{
			"parent": defaultInterface,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create the bot network: %v", err)
	}
	botNetwork, err := sup.client.GetNetworkByID(sup.ctx, botNetworkID)
	if err != nil {
		return fmt.Errorf("failed to get bot network id: %v", err)
	}
	botNetworkConfig := botNetwork.IPAM.Config[0]

	serviceNetworkID, err := sup.client.CreatePublicNetwork(sup.ctx, config.DockerServiceNetworkName)
	if err != nil {
		return fmt.Errorf("failed to create the service network: %v", err)
	}
	if err := sup.client.AttachNetwork(sup.ctx, supervisorContainer.ID, serviceNetworkID); err != nil {
		return fmt.Errorf("failed to attach supervisor container to node network: %v", err)
	}
	serviceNetwork, err := sup.client.GetNetworkByID(sup.ctx, serviceNetworkID)
	if err != nil {
		return fmt.Errorf("failed to get service network id: %v", err)
	}
	serviceNetworkConfig := serviceNetwork.IPAM.Config[0]

	// create the network manager for the bots
	defaultGwIPAddr := net.ParseIP(defaultGateway)
	_, botNetworkSubnet, _ := net.ParseCIDR(botNetworkConfig.Subnet)
	_, serviceNetworkSubnet, _ := net.ParseCIDR(serviceNetworkConfig.Subnet)
	sup.botManager = netmgmt.NewBotManager(sup.ctx, sup.client, &defaultGwIPAddr, []*net.IPNet{
		botNetworkSubnet, serviceNetworkSubnet,
	})

	sup.botNetworkID = botNetworkID
	sup.serviceNetworkID = serviceNetworkID

	ipfsContainer, err := sup.client.StartContainer(sup.ctx, clients.DockerContainerConfig{
		Name:  config.DockerIpfsContainerName,
		Image: "ipfs/go-ipfs:v0.12.2",
		Ports: map[string]string{
			"5001": "5001",
		},
		NetworkID:   serviceNetworkID,
		MaxLogFiles: sup.maxLogFiles,
		MaxLogSize:  sup.maxLogSize,
	})
	if err != nil {
		return err
	}
	sup.addContainerUnsafe(ipfsContainer, false)

	// start nats, wait for it and connect from the supervisor
	// TODO: Check if NATS is exposed after these changes. Try not to expose it if that's the case.
	natsContainer, err := sup.client.StartContainer(sup.ctx, clients.DockerContainerConfig{
		Name:  config.DockerNatsContainerName,
		Image: "nats:2.3.2",
		Ports: map[string]string{
			"4222": "4222",
			"6222": "6222",
			"8222": "8222",
		},
		NetworkID:   serviceNetworkID,
		MaxLogFiles: sup.maxLogFiles,
		MaxLogSize:  sup.maxLogSize,
	})
	if err != nil {
		return err
	}
	sup.addContainerUnsafe(natsContainer, false)

	if err := sup.client.WaitContainerStart(sup.ctx, natsContainer.ID); err != nil {
		return fmt.Errorf("failed while waiting for nats to start: %v", err)
	}
	// in tests, this is already set to a mock client
	if sup.msgClient == nil {
		sup.msgClient = messaging.NewClient("supervisor", fmt.Sprintf("%s:%s", config.DockerNatsContainerName, config.DefaultNatsPort))
	}
	sup.registerMessageHandlers()

	sup.jsonRpcContainer, err = sup.client.StartContainer(sup.ctx, clients.DockerContainerConfig{
		Name:  config.DockerJSONRPCProxyContainerName,
		Image: sup.commonNodeImage,
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
		NetworkID:   serviceNetworkID,
		MaxLogFiles: sup.maxLogFiles,
		MaxLogSize:  sup.maxLogSize,
	})
	if err != nil {
		return err
	}
	sup.addContainerUnsafe(sup.jsonRpcContainer, false)

	sup.scannerContainer, err = sup.client.StartContainer(sup.ctx, clients.DockerContainerConfig{
		Name:  config.DockerScannerContainerName,
		Image: sup.commonNodeImage,
		Cmd:   []string{config.DefaultFortaNodeBinaryPath, "scanner"},
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
		DialHost:    true,
		NetworkID:   serviceNetworkID,
		MaxLogFiles: sup.maxLogFiles,
		MaxLogSize:  sup.maxLogSize,
	})
	if err != nil {
		return err
	}
	sup.addContainerUnsafe(sup.scannerContainer, false)

	return nil
}

func (sup *SupervisorService) getHostNetworkingInfo(container *clients.DockerContainer) (string, string, string, error) {
	ctx, cancel := context.WithTimeout(sup.ctx, time.Second*30)
	defer cancel()
	ticker := time.NewTicker(time.Second * 2)
	logger := log.WithField("container", container.Name)
	for {
		select {
		case <-ctx.Done():
			return "", "", "", fmt.Errorf("failed to get host networking info: %v", ctx.Err())
		case <-ticker.C:
			container, err := sup.client.GetContainerByID(ctx, container.ID)
			if err != nil {
				logger.WithError(err).Warn("failed to get host networking info - retrying")
			}
			if container.State != "exited" {
				logger.WithField("state", container.State).Info("waiting for 'exited' state")
				continue
			}
			output, err := sup.client.GetContainerLogs(sup.ctx, container.ID, "", -1)
			if err != nil {
				logger.WithError(err).Error("failed to get container output")
				return "", "", "", err
			}
			parts := strings.Split(output, " ")
			return parts[0], parts[1], parts[2], nil
		}
	}
}

func (sup *SupervisorService) attachToNetwork(containerName, networkID string) error {
	container, err := sup.client.GetContainerByName(sup.ctx, containerName)
	if err != nil {
		return fmt.Errorf("failed to get '%s' container while attaching to network: %v", containerName, err)
	}
	if err := sup.client.AttachNetwork(sup.ctx, container.ID, networkID); err != nil {
		return fmt.Errorf("failed to attach '%s' container to network: %v", containerName, err)
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
		{
			Name: "ipfs/go-ipfs",
			Ref:  "ipfs/go-ipfs:v0.12.2",
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
		config.DockerJSONRPCProxyContainerName,
		config.DockerNatsContainerName,
		config.DockerIpfsContainerName,
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
	// WARNING: removing legacy networks only
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
	// TODO: Regard supervisor strategy on service and bot networks, remove them here if old so we replace.

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
		var err error
		if cnt.IsAgent {
			err = sup.client.StopContainer(ctx, cnt.DockerContainer.ID)
		} else {
			err = sup.client.InterruptContainer(ctx, cnt.DockerContainer.ID)
		}
		logger := log.WithFields(log.Fields{
			"id":      cnt.ID,
			"isAgent": cnt.IsAgent,
		})
		if err != nil {
			logger.WithError(err).Error("error stopping container")
		} else {
			logger.Info("requested to stop container")
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
	if len(sup.containers) < config.DockerSupervisorManagedContainers {
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

	// agent image client is helpful for loading private mode agents from a restricted container registry
	var agentImageClient clients.DockerClient
	if cfg.Config.PrivateModeConfig.Enable && cfg.Config.PrivateModeConfig.ContainerRegistry != nil {
		agentImageClient, err = clients.NewAuthDockerClient(
			"",
			cfg.Config.PrivateModeConfig.ContainerRegistry.Username,
			cfg.Config.PrivateModeConfig.ContainerRegistry.Password,
		)
	} else {
		agentImageClient, err = clients.NewDockerClient("")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create the private docker client: %v", err)
	}

	return &SupervisorService{
		ctx:              ctx,
		client:           dockerClient,
		globalClient:     globalClient,
		agentImageClient: agentImageClient,
		releaseClient:    releaseClient,
		config:           cfg,
		healthClient:     health.NewClient(),
		agentLogsClient:  agentlogs.NewClient(cfg.Config.AgentLogsConfig.URL),
	}, nil
}
