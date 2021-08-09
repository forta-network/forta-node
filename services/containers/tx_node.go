package containers

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"OpenZeppelin/fortify-node/clients"
	"OpenZeppelin/fortify-node/clients/messaging"
	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/store"
)

// TxNodeService manages the safe-node docker container as a service
type TxNodeService struct {
	ctx         context.Context
	client      clients.DockerClient
	agentClient clients.DockerClient

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

	t.msgClient = messaging.NewClient("cli", ":"+config.DefaultNatsPort) // accessible from localhost
	t.registerMessageHandlers()

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

	nodeNetwork, err := t.client.CreatePublicNetwork(t.ctx, config.DockerNetworkName)
	if err != nil {
		return err
	}

	natsContainer, err := t.client.StartContainer(t.ctx, clients.DockerContainerConfig{
		Name:  config.DockerNatsContainerName,
		Image: "nats:latest",
		Ports: map[string]string{
			"4222": "4222",
		},
		NetworkID:   nodeNetwork,
		MaxLogFiles: t.maxLogFiles,
		MaxLogSize:  t.maxLogSize,
	})
	if err != nil {
		return err
	}

	queryContainer, err := t.client.StartContainer(t.ctx, clients.DockerContainerConfig{
		Name:  config.DockerQueryContainerName,
		Image: t.config.Config.Query.QueryImage,
		Env: map[string]string{
			config.EnvFortifyConfig: cfgJson,
			config.EnvNatsHost:      config.DockerNatsContainerName,
		},
		Ports: map[string]string{
			fmt.Sprintf("%d", t.config.Config.Query.Port): "80",
		},
		Volumes: map[string]string{
			t.config.Config.Query.DB.Path: store.DBPath,
		},
		NetworkID:   nodeNetwork,
		MaxLogFiles: t.maxLogFiles,
		MaxLogSize:  t.maxLogSize,
	})
	if err != nil {
		return err
	}

	t.jsonRpcContainer, err = t.client.StartContainer(t.ctx, clients.DockerContainerConfig{
		Name:  config.DockerJSONRPCProxyContainerName,
		Image: t.config.Config.JsonRpcProxy.JsonRpcImage,
		Env: map[string]string{
			config.EnvFortifyConfig: cfgJson,
		},
		NetworkID:   nodeNetwork,
		MaxLogFiles: t.maxLogFiles,
		MaxLogSize:  t.maxLogSize,
	})
	if err != nil {
		return err
	}

	keyPath, err := config.GetKeyStorePath()
	if err != nil {
		return err
	}

	t.scannerContainer, err = t.client.StartContainer(t.ctx, clients.DockerContainerConfig{
		Name:  config.DockerScannerContainerName,
		Image: t.config.Config.Scanner.ScannerImage,
		Env: map[string]string{
			config.EnvFortifyConfig: cfgJson,
			config.EnvQueryNode:     config.DockerQueryContainerName,
			config.EnvNatsHost:      config.DockerNatsContainerName,
		},
		Volumes: map[string]string{
			keyPath: "/.keys",
		},
		Files: map[string][]byte{
			"passphrase": []byte(t.config.Passphrase),
		},
		NetworkID:   nodeNetwork,
		MaxLogFiles: t.maxLogFiles,
		MaxLogSize:  t.maxLogSize,
	})
	if err != nil {
		return err
	}

	t.addContainerUnsafe(natsContainer, queryContainer, t.jsonRpcContainer, t.scannerContainer)

	return nil
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
	agentDockerClient, err := clients.NewAuthDockerClient(cfg.Config.Registry.Username, cfg.Config.Registry.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to create the agent docker client: %v", err)
	}
	dockerClient, err := clients.NewDockerClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create the docker client: %v", err)
	}
	return &TxNodeService{
		ctx:         ctx,
		client:      dockerClient,
		agentClient: agentDockerClient,
		config:      cfg,
	}, nil
}
