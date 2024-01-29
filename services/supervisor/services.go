package supervisor

import (
	"context"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/forta-network/forta-core-go/clients/agentlogs"
	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/manifest"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-core-go/release"
	"github.com/forta-network/forta-core-go/security"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/docker"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/services"
	"github.com/forta-network/forta-node/services/components"
	"github.com/forta-network/forta-node/services/components/containers"
	"github.com/ipfs/go-cid"
	log "github.com/sirupsen/logrus"
)

var knownServiceContainerNames = []string{
	config.DockerScannerContainerName,
	config.DockerInspectorContainerName,
	config.DockerJSONRPCProxyContainerName,
	config.DockerJWTProviderContainerName,
	config.DockerPublicAPIProxyContainerName,
	config.DockerNatsContainerName,
	config.DockerIpfsContainerName,
	config.DockerStorageContainerName,
}

// SupervisorService manages the scanner node's service and agent containers.
type SupervisorService struct {
	ctx context.Context

	client       clients.DockerClient
	globalClient clients.DockerClient

	botLifecycleConfig components.BotLifecycleConfig
	botLifecycle       components.BotLifecycle

	manifestClient manifest.Client
	releaseClient  release.Client

	msgClient   clients.MessageClient
	config      SupervisorServiceConfig
	maxLogSize  string
	maxLogFiles int

	scannerContainer     *docker.Container
	inspectorContainer   *docker.Container
	jsonRpcContainer     *docker.Container
	publicAPIContainer   *docker.Container
	jwtProviderContainer *docker.Container
	storageContainer     *docker.Container
	containers           []*Container
	mu                   sync.RWMutex

	lastRun                         health.TimeTracker
	lastStop                        health.TimeTracker
	lastTelemetryRequest            health.TimeTracker
	lastTelemetryRequestError       health.ErrorTracker
	lastCustomTelemetryRequest      health.TimeTracker
	lastCustomTelemetryRequestError health.ErrorTracker
	lastAgentLogsRequest            health.TimeTracker
	lastAgentLogsRequestError       health.ErrorTracker
	autoUpdatesDisabled             health.MessageTracker

	healthClient health.HealthClient

	sendAgentLogs func(agents agentlogs.Agents, authToken string) error
	prevAgentLogs agentlogs.Agents
	inspectionCh  chan *protocol.InspectionResults
}

type SupervisorServiceConfig struct {
	Config             config.Config
	Passphrase         string
	Key                *keystore.Key
	BotLifecycleConfig components.BotLifecycleConfig
}

// Container extends the default container data.
type Container struct {
	docker.Container
	IsAgent     bool
	AgentConfig *config.AgentConfig
}

func (sup *SupervisorService) Start() error {
	if err := sup.start(); err != nil {
		return err
	}

	go sup.healthCheck()
	go sup.refreshBotContainers()

	return nil
}

func (sup *SupervisorService) start() error {
	// in addition to the feature disable flags, check local mode flags to disable agent logging and telemetry

	shouldDisableTelemetry := sup.config.Config.TelemetryConfig.Disable
	if !shouldDisableTelemetry {
		go sup.syncTelemetryData()
	}

	sup.mu.Lock()
	defer sup.mu.Unlock()

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

	// start of service network and container launch
	startTime := time.Now()

	nodeNetworkID, err := sup.client.EnsurePublicNetwork(sup.ctx, config.DockerNetworkName)
	if err != nil {
		return err
	}
	if err := sup.client.AttachNetwork(sup.ctx, supervisorContainer.ID, nodeNetworkID); err != nil {
		return fmt.Errorf("failed to attach supervisor container to node network: %v", err)
	}

	natsNetworkID, err := sup.client.EnsureInternalNetwork(sup.ctx, config.DockerNatsContainerName)
	if err != nil {
		return err
	}
	if err := sup.client.AttachNetwork(sup.ctx, supervisorContainer.ID, natsNetworkID); err != nil {
		return fmt.Errorf("failed to attach supervisor container to nats network: %v", err)
	}

	manageIpfsDir(sup.config.Config)
	if sup.config.Config.AdvancedConfig.IPFSExperiment {
		ipfsContainer, err := sup.client.StartContainer(sup.ctx, docker.ContainerConfig{
			Name:  config.DockerIpfsContainerName,
			Image: "ipfs/kubo:v0.16.0",
			Ports: map[string]string{
				"5001": "5001",
			},
			Files: map[string][]byte{
				"/container-init.d/001-init.sh": []byte(`
	#!/bin/sh
	
	set -xe
	
	ipfs config --bool Discovery.MDNS.Enabled 'false' && \
	ipfs config --json Routing '{"Type":"none"}' && \
	ipfs config --json Addresses.Swarm '[]' && \
	ipfs config --json Bootstrap '[]' && \
	ipfs config Datastore.StorageMax '1GB'
	`),
			},
			Volumes: map[string]string{
				path.Join(hostFortaDir, ".ipfs"): "/data/ipfs",
			},
			NetworkID:   nodeNetworkID,
			MaxLogFiles: sup.maxLogFiles,
			MaxLogSize:  sup.maxLogSize,
			Cmd: []string{
				// default CMD - taken from https://hub.docker.com/layers/ipfs/kubo/master-latest/images/sha256-65b4c19a75987bd9bb677e8d9b1b1dafb81eec2335ba65f73dfb8256f6b3d22a?context=explore
				"daemon", "--migrate=true", "--agent-version-suffix=docker",
				// extra flags
				"--offline",
			},
			CPUQuota: config.CPUsToMicroseconds(0.5),
		})
		if err != nil {
			return err
		}
		sup.addContainerUnsafe(ipfsContainer)
	}

	// start nats, wait for it and connect from the supervisor
	natsContainer, err := sup.client.StartContainer(sup.ctx, docker.ContainerConfig{
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
	sup.botLifecycleConfig.MessageClient = sup.msgClient // we are able to set this dependency only here
	sup.botLifecycle, err = components.GetBotLifecycleComponents(sup.ctx, sup.botLifecycleConfig)
	if err != nil {
		return fmt.Errorf("failed to get bot lifecycle components: %v", err)
	}

	shouldDisableAgentLogs := sup.config.Config.AgentLogsConfig.Disable || sup.config.Config.LocalModeConfig.Enable
	if !shouldDisableAgentLogs {
		go sup.syncAgentLogs()
	}

	sup.registerMessageHandlers()

	if sup.config.Config.AdvancedConfig.IPFSExperiment {
		sup.storageContainer, err = sup.client.StartContainer(
			sup.ctx, docker.ContainerConfig{
				Name:  config.DockerStorageContainerName,
				Image: commonNodeImage,
				Cmd:   []string{config.DefaultFortaNodeBinaryPath, "storage"},
				Env: map[string]string{
					config.EnvReleaseInfo: releaseInfo.String(),
				},
				Volumes: map[string]string{
					// give access to host docker
					"/var/run/docker.sock": "/var/run/docker.sock",
					hostFortaDir:           config.DefaultContainerFortaDirPath,
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
			},
		)
		if err != nil {
			return err
		}
		sup.addContainerUnsafe(sup.storageContainer)

		if err := sup.client.WaitContainerStart(sup.ctx, sup.storageContainer.ID); err != nil {
			return fmt.Errorf("failed while waiting for the storage container to start: %v", err)
		}
	}

	sup.jsonRpcContainer, err = sup.client.StartContainer(
		sup.ctx, docker.ContainerConfig{
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
			DialHost:       true,
			NetworkID:      nodeNetworkID,
			LinkNetworkIDs: []string{natsNetworkID},
			MaxLogFiles:    sup.maxLogFiles,
			MaxLogSize:     sup.maxLogSize,
		},
	)
	if err != nil {
		return err
	}
	sup.addContainerUnsafe(sup.jsonRpcContainer)

	sup.publicAPIContainer, err = sup.client.StartContainer(
		sup.ctx, docker.ContainerConfig{
			Name:  config.DockerPublicAPIProxyContainerName,
			Image: commonNodeImage,
			Cmd:   []string{config.DefaultFortaNodeBinaryPath, "public-api"},
			Volumes: map[string]string{
				// give access to host docker
				"/var/run/docker.sock": "/var/run/docker.sock",
				hostFortaDir:           config.DefaultContainerFortaDirPath,
			},
			Ports: map[string]string{
				"": config.DefaultHealthPort, // random host port
			},
			Files: map[string][]byte{
				"passphrase": []byte(sup.config.Passphrase),
			},
			DialHost:       true,
			NetworkID:      nodeNetworkID,
			LinkNetworkIDs: []string{natsNetworkID},
			MaxLogFiles:    sup.maxLogFiles,
			MaxLogSize:     sup.maxLogSize,
		},
	)
	if err != nil {
		return err
	}
	sup.addContainerUnsafe(sup.publicAPIContainer)

	shouldInspectAtStartup := *sup.config.Config.InspectionConfig.InspectAtStartup
	if sup.config.Config.LocalModeConfig.Enable {
		shouldInspectAtStartup = shouldInspectAtStartup && sup.config.Config.LocalModeConfig.ForceEnableInspection
	}

	if shouldInspectAtStartup {
		if err := sup.client.WaitContainerStart(sup.ctx, sup.jsonRpcContainer.ID); err != nil {
			return fmt.Errorf("failed while waiting for json-rpc container to start: %v", err)
		}
	}

	sup.inspectorContainer, err = sup.client.StartContainer(
		sup.ctx, docker.ContainerConfig{
			Name:  config.DockerInspectorContainerName,
			Image: commonNodeImage,
			Cmd:   []string{config.DefaultFortaNodeBinaryPath, "inspector"},
			Volumes: map[string]string{
				hostFortaDir: config.DefaultContainerFortaDirPath,
			},
			Ports: map[string]string{
				"": config.DefaultHealthPort, // random host port
			},
			Files: map[string][]byte{
				"passphrase": []byte(sup.config.Passphrase),
			},
			DialHost:       true,
			NetworkID:      nodeNetworkID,
			LinkNetworkIDs: []string{natsNetworkID},
			MaxLogFiles:    sup.maxLogFiles,
			MaxLogSize:     sup.maxLogSize,
		},
	)
	if err != nil {
		return err
	}
	sup.addContainerUnsafe(sup.inspectorContainer)

	if shouldInspectAtStartup {
		if err := sup.client.WaitContainerStart(sup.ctx, sup.inspectorContainer.ID); err != nil {
			return fmt.Errorf("failed while waiting for inspector to start: %v", err)
		}

		// this makes sure that inspector published a message. Which means publisher has also received it and
		// inspection results will be available for every batch starting first batch.
		log.Info("waiting for the first inspection to complete...")
		<-sup.inspectionCh
		log.Info("inspection to completed")
	}

	go func() {
		// wait for the publisher so it can catch the metrics
		time.Sleep(time.Minute)
		go containers.ListenToDockerEvents(sup.ctx, sup.globalClient, sup.msgClient, startTime)
	}()

	sup.scannerContainer, err = sup.client.StartContainer(
		sup.ctx, docker.ContainerConfig{
			Name:  config.DockerScannerContainerName,
			Image: commonNodeImage,
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
			DialHost:       true,
			NetworkID:      nodeNetworkID,
			LinkNetworkIDs: []string{natsNetworkID},
			MaxLogFiles:    sup.maxLogFiles,
			MaxLogSize:     sup.maxLogSize,
		},
	)
	if err != nil {
		return err
	}
	sup.addContainerUnsafe(sup.scannerContainer)

	sup.jwtProviderContainer, err = sup.client.StartContainer(
		sup.ctx, docker.ContainerConfig{
			Name:  config.DockerJWTProviderContainerName,
			Image: commonNodeImage,
			Cmd:   []string{config.DefaultFortaNodeBinaryPath, "jwt-provider"},
			Env: map[string]string{
				config.EnvReleaseInfo: releaseInfo.String(),
			},
			Volumes: map[string]string{
				// give access to host docker
				"/var/run/docker.sock": "/var/run/docker.sock",
				hostFortaDir:           config.DefaultContainerFortaDirPath,
			},
			Ports: map[string]string{
				"": config.DefaultHealthPort, // random host port
			},
			Files: map[string][]byte{
				"passphrase": []byte(sup.config.Passphrase),
			},
			DialHost:       true,
			NetworkID:      nodeNetworkID,
			LinkNetworkIDs: []string{natsNetworkID},
			MaxLogFiles:    sup.maxLogFiles,
			MaxLogSize:     sup.maxLogSize,
		},
	)
	if err != nil {
		return err
	}
	sup.addContainerUnsafe(sup.jwtProviderContainer)

	return nil
}

func (sup *SupervisorService) addContainerUnsafe(container *docker.Container, agentConfig ...*config.AgentConfig) {
	if agentConfig != nil {
		sup.containers = append(
			sup.containers, &Container{
				Container:   *container,
				IsAgent:     true,
				AgentConfig: agentConfig[0],
			},
		)
		return
	}
	sup.containers = append(sup.containers, &Container{Container: *container})
}

func (sup *SupervisorService) registerMessageHandlers() {
	if *sup.config.Config.InspectionConfig.InspectAtStartup {
		sup.msgClient.Subscribe(messaging.SubjectInspectionDone, messaging.InspectionResultsHandler(sup.handleInspectionResults))
	}
}

func manageIpfsDir(cfg config.Config) error {
	if !cfg.AdvancedConfig.IPFSExperiment {
		// purge the dir if the experiment is disabled
		os.RemoveAll(path.Join(config.DefaultContainerFortaDirPath, ".ipfs"))
		return nil
	}
	if err := os.MkdirAll(path.Join(config.DefaultContainerFortaDirPath, ".ipfs"), 0700); err != nil {
		return fmt.Errorf("failed to create ipfs dir: %v", err)
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
		{
			Name: "ipfs/kubo",
			Ref:  "ipfs/kubo:v0.16.0",
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
	for _, containerName := range knownServiceContainerNames {
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
	containerList, err := sup.client.GetContainers(sup.ctx)
	if err != nil {
		return fmt.Errorf("failed to get containers list: %v", err)
	}
	for _, container := range containerList {
		containerName := container.Names[0][1:]
		logger := log.WithFields(log.Fields{
			"containerName": containerName,
			"containerId":   container.ID,
		})
		if !strings.Contains(containerName, "forta-agent-") {
			continue
		}
		if !containers.HasSameLabelValue(
			&container,
			docker.LabelFortaSupervisorStrategyVersion, containers.LabelValueStrategyVersion,
		) {
			logger.Info("bot container is old - need to remove")
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
	slowTicker := time.NewTicker(time.Minute * 5)
	fastTicker := time.NewTicker(time.Minute)
	for {
		var err error
		select {
		case <-slowTicker.C:
			err = sup.doSyncTelemetryDataToPublicHandler()
			if err != nil {
				log.WithError(err).Warn("telemetry sync failed (public handler)")
			}
			sup.lastTelemetryRequest.Set()
			sup.lastTelemetryRequestError.Set(err)

		case <-fastTicker.C:
			err = sup.doSyncTelemetryDataToCustomHandler()
			if err != nil {
				log.WithError(err).Warn("telemetry sync failed (custom handler)")
			}
			sup.lastCustomTelemetryRequest.Set()
			sup.lastCustomTelemetryRequestError.Set(err)

		case <-sup.ctx.Done():
			return
		}
	}
}

func (sup *SupervisorService) doSyncTelemetryDataToPublicHandler() error {
	return sup.doSyncTelemetryData(sup.config.Config.TelemetryConfig.URL)
}

func (sup *SupervisorService) doSyncTelemetryDataToCustomHandler() error {
	customURL := sup.config.Config.TelemetryConfig.CustomURL
	if len(customURL) == 0 {
		return nil
	}
	return sup.doSyncTelemetryData(customURL)
}

func (sup *SupervisorService) doSyncTelemetryData(destUrl string) error {
	scannerJwt, err := security.CreateScannerJWT(sup.config.Key, map[string]interface{}{
		"access": "telemetry",
	})
	if err != nil {
		return err
	}
	dataSrc := fmt.Sprintf("http://host.docker.internal:%s/health", config.DefaultHealthPort)
	return sup.healthClient.SendReports(
		dataSrc,
		destUrl,
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
	fullReleaseManifest, err := sup.releaseClient.GetReleaseManifest(releaseInfo.IPFS)
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

	// we use the background context here because
	// we don't want tear downs to be aborted by the closed service context
	ctx := context.Background()

	if !services.IsGracefulShutdown() {
		sup.botLifecycle.BotManager.TearDownRunningBots(ctx)
	}

	for _, cnt := range sup.containers {
		err := sup.client.InterruptContainer(ctx, cnt.Container.ID)
		logger := log.WithFields(log.Fields{
			"id": cnt.ID,
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
			Name:    "local-mode",
			Status:  health.StatusInfo,
			Details: strconv.FormatBool(sup.config.Config.LocalModeConfig.Enable),
		},
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
		sup.lastCustomTelemetryRequest.GetReport("event.custom-telemetry-sync.time"),
		sup.lastCustomTelemetryRequestError.GetReport("event.custom-telemetry-sync.error"),
		sup.lastAgentLogsRequest.GetReport("event.agent-logs-sync.time"),
		sup.lastAgentLogsRequestError.GetReport("event.agent-logs-sync.error"),
		sup.autoUpdatesDisabled.GetReport("auto-updates.disabled"),
	}
}

// handleInspectionResults listen for inspections.
func (sup *SupervisorService) handleInspectionResults(payload *protocol.InspectionResults) error {
	// do a non-blocking write because messages are consumed only at startup
	select {
	case sup.inspectionCh <- payload:
		return nil
	default:
		return nil
	}
}

func NewSupervisorService(ctx context.Context, cfg SupervisorServiceConfig) (*SupervisorService, error) {
	dockerClient, err := docker.NewDockerClient(containers.LabelFortaSupervisor)
	if err != nil {
		return nil, fmt.Errorf("failed to create the docker client: %v", err)
	}
	globalClient, err := docker.NewDockerClient("")
	if err != nil {
		return nil, fmt.Errorf("failed to create the global docker client: %v", err)
	}

	releaseClient, err := release.NewClient(cfg.Config.Registry.IPFS.GatewayURL, []string{cfg.Config.Registry.ReleaseDistributionUrl})
	if err != nil {
		return nil, fmt.Errorf("failed to create the release client: %v", err)
	}

	sup := &SupervisorService{
		ctx:                ctx,
		client:             dockerClient,
		globalClient:       globalClient,
		releaseClient:      releaseClient,
		botLifecycleConfig: cfg.BotLifecycleConfig,
		config:             cfg,
		healthClient:       health.NewClient(),
		sendAgentLogs:      agentlogs.NewClient(cfg.Config.AgentLogsConfig.URL).SendLogs,
		inspectionCh:       make(chan *protocol.InspectionResults),
	}
	sup.autoUpdatesDisabled.Set(strconv.FormatBool(cfg.Config.AutoUpdate.Disable))

	return sup, nil
}
