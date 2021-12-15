package containers

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

// TxNodeService manages the safe-node docker container as a service
type TxNodeService struct {
	ctx        context.Context
	client     clients.DockerClient
	authClient clients.DockerClient

	msgClient   clients.MessageClient
	config      TxNodeServiceConfig
	maxLogSize  string
	maxLogFiles int

	scannerContainer *clients.DockerContainer
	jsonRpcContainer *clients.DockerContainer
	containers       []*clients.DockerContainer
	mu               sync.RWMutex
}

type TxNodeServiceConfig struct {
	Config     config.Config
	Passphrase string
}

func (t *TxNodeService) Start() error {
	if err := t.start(); err != nil {
		return err
	}

	go t.healthCheck()

	return nil
}

func (t *TxNodeService) start() error {
	log.Infof("Starting %s", t.Name())
	_, err := log.ParseLevel(t.config.Config.Log.Level)
	if err != nil {
		log.Error("invalid log level", err)
		return err
	}

	t.maxLogSize = t.config.Config.Log.MaxLogSize
	t.maxLogFiles = t.config.Config.Log.MaxLogFiles

	cfgBytes, err := json.Marshal(t.config.Config)
	if err != nil {
		log.Error("cannot marshal config to json", err)
		return err
	}
	cfgJson := string(cfgBytes)

	if err := t.client.Prune(t.ctx); err != nil {
		return err
	}

	if config.UseDockerImages == "remote" {
		if err := t.ensureNodeImages(); err != nil {
			return err
		}
	}

	supervisorContainer, err := t.client.GetContainerByName(t.ctx, config.DockerSupervisorContainerName)
	if err != nil {
		return fmt.Errorf("failed to get the supervisor container: %v", err)
	}
	commonNodeImage := supervisorContainer.Image

	nodeNetworkID, err := t.client.CreatePublicNetwork(t.ctx, config.DockerNetworkName)
	if err != nil {
		return err
	}
	if err := t.attachSupervisor(nodeNetworkID); err != nil {
		return err
	}

	// start nats, wait for it and connect from the supervisor
	natsContainer, err := t.client.StartContainer(t.ctx, clients.DockerContainerConfig{
		Name:  config.DockerNatsContainerName,
		Image: "nats:2.3.2",
		Ports: map[string]string{
			"4222": "4222",
		},
		NetworkID:   nodeNetworkID,
		MaxLogFiles: t.maxLogFiles,
		MaxLogSize:  t.maxLogSize,
	})
	if err != nil {
		return err
	}

	if err := t.client.WaitContainerStart(t.ctx, natsContainer.ID); err != nil {
		return fmt.Errorf("failed while waiting for nats to start: %v", err)
	}
	// in tests, this is already set to a mock client
	if t.msgClient == nil {
		t.msgClient = messaging.NewClient("supervisor", fmt.Sprintf("%s:%s", config.DockerNatsContainerName, config.DefaultNatsPort))
	}
	t.registerMessageHandlers()

	queryContainer, err := t.client.StartContainer(t.ctx, clients.DockerContainerConfig{
		Name:  config.DockerPublisherContainerName,
		Image: commonNodeImage,
		Cmd:   []string{config.DefaultFortaNodeBinaryPath, "publisher"},
		Env: map[string]string{
			config.EnvConfig:   cfgJson,
			config.EnvFortaDir: config.DefaultContainerFortaDirPath,
			config.EnvNatsHost: config.DockerNatsContainerName,
		},
		Ports: map[string]string{
			fmt.Sprintf("%d", t.config.Config.Query.Port): "80",
		},
		Volumes: map[string]string{
			t.config.Config.FortaDir: config.DefaultContainerFortaDirPath,
		},
		Files: map[string][]byte{
			"passphrase": []byte(t.config.Passphrase),
		},
		NetworkID:   nodeNetworkID,
		MaxLogFiles: t.maxLogFiles,
		MaxLogSize:  t.maxLogSize,
	})
	if err != nil {
		return err
	}

	t.jsonRpcContainer, err = t.client.StartContainer(t.ctx, clients.DockerContainerConfig{
		Name:  config.DockerJSONRPCProxyContainerName,
		Image: commonNodeImage,
		Cmd:   []string{config.DefaultFortaNodeBinaryPath, "json-rpc"},
		Env: map[string]string{
			config.EnvConfig: cfgJson,
		},
		NetworkID:   nodeNetworkID,
		MaxLogFiles: t.maxLogFiles,
		MaxLogSize:  t.maxLogSize,
	})
	if err != nil {
		return err
	}

	t.scannerContainer, err = t.client.StartContainer(t.ctx, clients.DockerContainerConfig{
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
			t.config.Config.FortaDir: config.DefaultContainerFortaDirPath,
		},
		Files: map[string][]byte{
			"passphrase": []byte(t.config.Passphrase),
		},
		NetworkID:   nodeNetworkID,
		MaxLogFiles: t.maxLogFiles,
		MaxLogSize:  t.maxLogSize,
	})
	if err != nil {
		return err
	}

	t.addContainerUnsafe(natsContainer, queryContainer, t.jsonRpcContainer, t.scannerContainer)

	return nil
}

func (t *TxNodeService) attachSupervisor(nodeNetworkID string) error {
	container, err := t.client.GetContainerByName(t.ctx, config.DockerSupervisorContainerName)
	if err != nil {
		return fmt.Errorf("failed to get supervisor container while attaching to node network: %v", err)
	}
	if err := t.client.AttachNetwork(t.ctx, container.ID, nodeNetworkID); err != nil {
		return fmt.Errorf("failed to attach supervisor to node network: %v", err)
	}
	return nil
}

func (t *TxNodeService) ensureNodeImages() error {
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
		if err := t.ensureLocalImage(image.Name, image.Ref, image.RequireAuth); err != nil {
			return err
		}
	}
	return nil
}

func (t *TxNodeService) ensureLocalImage(name, ref string, requireAuth bool) error {
	client := t.client
	if requireAuth {
		client = t.authClient
	}
	return client.EnsureLocalImage(t.ctx, name, ref)
}

func (t *TxNodeService) Stop() error {
	t.mu.RLock()
	defer t.mu.RUnlock()

	ctx := context.Background()
	for _, cnt := range t.containers {
		if err := t.client.StopContainer(ctx, cnt.ID); err != nil {
			log.Error(fmt.Sprintf("error stopping %s container", cnt.ID), err)
		} else {
			log.Infof("Container %s is stopped", cnt.ID)
		}
	}
	return nil
}

func (t *TxNodeService) Name() string {
	return "TxNode"
}

func NewTxNodeService(ctx context.Context, cfg TxNodeServiceConfig) (*TxNodeService, error) {
	dockerAuthClient, err := clients.NewAuthDockerClient(cfg.Config.Registry.Username, cfg.Config.Registry.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to create the agent docker client: %v", err)
	}
	dockerClient, err := clients.NewDockerClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create the docker client: %v", err)
	}
	return &TxNodeService{
		ctx:        ctx,
		client:     dockerClient,
		authClient: dockerAuthClient,
		config:     cfg,
	}, nil
}
