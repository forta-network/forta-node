package supervisor

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/forta-protocol/forta-node/security"

	log "github.com/sirupsen/logrus"

	"github.com/forta-protocol/forta-node/clients"
	"github.com/forta-protocol/forta-node/clients/messaging"
	"github.com/forta-protocol/forta-node/config"
)

// SupervisorService manages the scanner node's service and agent containers.
type SupervisorService struct {
	ctx        context.Context
	client     clients.DockerClient
	authClient clients.DockerClient

	msgClient   clients.MessageClient
	config      SupervisorServiceConfig
	maxLogSize  string
	maxLogFiles int

	scannerContainer *clients.DockerContainer
	jsonRpcContainer *clients.DockerContainer
	containers       []*clients.DockerContainer
	mu               sync.RWMutex
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

	sup.maxLogSize = sup.config.Config.Log.MaxLogSize
	sup.maxLogFiles = sup.config.Config.Log.MaxLogFiles

	passphrase, err := security.ReadPassphrase()
	if err != nil {
		return err
	}
	sup.config.Passphrase = passphrase

	if err := sup.client.Prune(sup.ctx); err != nil {
		return err
	}

	if err := sup.ensureNodeImages(); err != nil {
		return err
	}

	supervisorContainer, err := sup.client.GetContainerByName(sup.ctx, config.DockerSupervisorContainerName)
	if err != nil {
		return fmt.Errorf("failed to get the supervisor container: %v", err)
	}
	commonNodeImage := supervisorContainer.Image

	nodeNetworkID, err := sup.client.CreatePublicNetwork(sup.ctx, config.DockerNetworkName)
	if err != nil {
		return err
	}
	if err := sup.attachToNetwork(config.DockerSupervisorContainerName, nodeNetworkID); err != nil {
		return err
	}

	var natsNetworkID string
	if sup.config.Config.ExposeNats {
		natsNetworkID = nodeNetworkID
	} else {
		natsNetworkID, err = sup.client.CreateInternalNetwork(sup.ctx, config.DockerNatsContainerName)
		if err != nil {
			return err
		}
		if err := sup.attachToNetwork(config.DockerSupervisorContainerName, natsNetworkID); err != nil {
			return err
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
		Volumes: map[string]string{
			hostFortaDir: config.DefaultContainerFortaDirPath,
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
		Name        string
		Ref         string
		RequireAuth bool
	}{
		{
			Name: "nats",
			Ref:  "nats:2.3.2",
		},
	} {
		if err := sup.ensureLocalImage(image.Name, image.Ref, image.RequireAuth); err != nil {
			return err
		}
	}
	return nil
}

func (sup *SupervisorService) ensureLocalImage(name, ref string, requireAuth bool) error {
	client := sup.client
	if requireAuth {
		client = sup.authClient
	}
	return client.EnsureLocalImage(sup.ctx, name, ref)
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
	return "Supervisor"
}

func NewSupervisorService(ctx context.Context, cfg SupervisorServiceConfig) (*SupervisorService, error) {
	dockerAuthClient, err := clients.NewAuthDockerClient("supervisor", cfg.Config.Registry.Username, cfg.Config.Registry.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to create the agent docker client: %v", err)
	}
	dockerClient, err := clients.NewDockerClient("supervisor")
	if err != nil {
		return nil, fmt.Errorf("failed to create the docker client: %v", err)
	}
	return &SupervisorService{
		ctx:        ctx,
		client:     dockerClient,
		authClient: dockerAuthClient,
		config:     cfg,
	}, nil
}
