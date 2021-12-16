package supervisor

import (
	"context"
	"fmt"
	"sync"

	"github.com/goccy/go-json"

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

	sup.maxLogSize = sup.config.Config.Log.MaxLogSize
	sup.maxLogFiles = sup.config.Config.Log.MaxLogFiles

	cfgBytes, err := json.Marshal(sup.config.Config)
	if err != nil {
		log.Error("cannot marshal config to json", err)
		return err
	}
	cfgJson := string(cfgBytes)

	if err := sup.client.Prune(sup.ctx); err != nil {
		return err
	}

	if config.UseDockerImages == "remote" {
		if err := sup.ensureNodeImages(); err != nil {
			return err
		}
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
	if err := sup.attachSupervisor(nodeNetworkID); err != nil {
		return err
	}

	// start nats, wait for it and connect from the supervisor
	natsContainer, err := sup.client.StartContainer(sup.ctx, clients.DockerContainerConfig{
		Name:  config.DockerNatsContainerName,
		Image: "nats:2.3.2",
		Ports: map[string]string{
			"4222": "4222",
		},
		NetworkID:   nodeNetworkID,
		MaxLogFiles: sup.maxLogFiles,
		MaxLogSize:  sup.maxLogSize,
	})
	if err != nil {
		return err
	}

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
			config.EnvConfig:   cfgJson,
			config.EnvFortaDir: config.DefaultContainerFortaDirPath,
			config.EnvNatsHost: config.DockerNatsContainerName,
		},
		Ports: map[string]string{
			fmt.Sprintf("%d", sup.config.Config.Query.Port): "80",
		},
		Volumes: map[string]string{
			sup.config.Config.FortaDir: config.DefaultContainerFortaDirPath,
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

	sup.jsonRpcContainer, err = sup.client.StartContainer(sup.ctx, clients.DockerContainerConfig{
		Name:  config.DockerJSONRPCProxyContainerName,
		Image: commonNodeImage,
		Cmd:   []string{config.DefaultFortaNodeBinaryPath, "json-rpc"},
		Env: map[string]string{
			config.EnvConfig: cfgJson,
		},
		NetworkID:   nodeNetworkID,
		MaxLogFiles: sup.maxLogFiles,
		MaxLogSize:  sup.maxLogSize,
	})
	if err != nil {
		return err
	}

	sup.scannerContainer, err = sup.client.StartContainer(sup.ctx, clients.DockerContainerConfig{
		Name:  config.DockerScannerContainerName,
		Image: commonNodeImage,
		Cmd:   []string{config.DefaultFortaNodeBinaryPath, "scanner"},
		Env: map[string]string{
			config.EnvConfig:        cfgJson,
			config.EnvFortaDir:      config.DefaultContainerFortaDirPath,
			config.EnvPublisherHost: config.DockerPublisherContainerName,
			config.EnvNatsHost:      config.DockerNatsContainerName,
		},
		Ports: map[string]string{
			"8989": "80",
		},
		Volumes: map[string]string{
			sup.config.Config.FortaDir: config.DefaultContainerFortaDirPath,
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

	sup.addContainerUnsafe(natsContainer, publisherContainer, sup.jsonRpcContainer, sup.scannerContainer)

	return nil
}

func (sup *SupervisorService) attachSupervisor(nodeNetworkID string) error {
	container, err := sup.client.GetContainerByName(sup.ctx, config.DockerSupervisorContainerName)
	if err != nil {
		return fmt.Errorf("failed to get supervisor container while attaching to node network: %v", err)
	}
	if err := sup.client.AttachNetwork(sup.ctx, container.ID, nodeNetworkID); err != nil {
		return fmt.Errorf("failed to attach supervisor to node network: %v", err)
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
		{
			Name:        "node",
			Ref:         config.DockerScannerNodeImage,
			RequireAuth: true,
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
	dockerAuthClient, err := clients.NewAuthDockerClient(cfg.Config.Registry.Username, cfg.Config.Registry.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to create the agent docker client: %v", err)
	}
	dockerClient, err := clients.NewDockerClient()
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
