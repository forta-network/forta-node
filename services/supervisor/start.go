package supervisor

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/forta-protocol/forta-node/clients"
	"github.com/forta-protocol/forta-node/clients/health"
	"github.com/forta-protocol/forta-node/clients/messaging"
	"github.com/forta-protocol/forta-node/config"
)

// SupervisorService manages the scanner node's service and agent containers.
type SupervisorService struct {
	ctx          context.Context
	client       clients.DockerClient
	globalClient clients.DockerClient

	msgClient   clients.MessageClient
	config      SupervisorServiceConfig
	maxLogSize  string
	maxLogFiles int

	scannerContainer *clients.DockerContainer
	jsonRpcContainer *clients.DockerContainer
	containers       []*clients.DockerContainer
	mu               sync.RWMutex

	lastRun  health.TimeTracker
	lastStop health.TimeTracker
}

type SupervisorServiceConfig struct {
	Config     config.Config
	Passphrase string
}

func (sup *SupervisorService) Start() error {
	if err := sup.start(); err != nil {
		return err
	}

	go sup.healthCheck()

	return nil
}

func (sup *SupervisorService) start() error {
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
	releaseInfo := config.ReleaseInfoFromString(os.Getenv(config.EnvReleaseInfo))

	sup.maxLogSize = sup.config.Config.Log.MaxLogSize
	sup.maxLogFiles = sup.config.Config.Log.MaxLogFiles

	if err := sup.client.Prune(sup.ctx); err != nil {
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
			hostFortaDir: config.DefaultContainerFortaDirPath,
		},
		Ports: map[string]string{
			"": config.DefaultHealthPort, // random host port
		},
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

func (sup *SupervisorService) Stop() error {
	sup.mu.RLock()
	defer sup.mu.RUnlock()

	ctx := context.Background()
	for _, cnt := range sup.containers {
		if err := sup.client.StopContainer(ctx, cnt.ID); err != nil {
			log.Error(fmt.Sprintf("error stopping %s container", cnt.ID), err)
		} else {
			log.Infof("Container %s is stopped", cnt.ID)
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
			Name:    "events.run-agent",
			Status:  health.StatusInfo,
			Details: sup.lastRun.String(),
		},
		&health.Report{
			Name:    "events.stop-agent",
			Status:  health.StatusInfo,
			Details: sup.lastStop.String(),
		},
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
	return &SupervisorService{
		ctx:          ctx,
		client:       dockerClient,
		globalClient: globalClient,
		config:       cfg,
	}, nil
}
