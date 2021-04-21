package services

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"OpenZeppelin/zephyr-node/clients"
	"OpenZeppelin/zephyr-node/config"
)

const EnvJsonRpcUrl = "JSON_RPC_URL"
const EnvStartBlock = "START_BLOCK"

// TxNodeService manages the safe-node docker container as a service
type TxNodeService struct {
	ctx       context.Context
	client    clients.DockerClient
	container *clients.DockerContainer
	config    TxNodeServiceConfig
}

type TxNodeServiceConfig struct {
	JsonRpcUrl     string
	LogLevel       string
	ContainerImage string
	StartBlock     int
}

func (t *TxNodeService) Start() error {
	log.Infof("Starting %s", t.Name())
	_, err := log.ParseLevel(t.config.LogLevel)
	if err != nil {
		log.Error("invalid log level", err)
		return err
	}
	var startBlock string
	if t.config.StartBlock != 0 {
		startBlock = fmt.Sprintf("%d", t.config.StartBlock)
	}
	container, err := t.client.StartContainer(t.ctx, clients.DockerContainerConfig{
		Name:  fmt.Sprintf("zephyr-node-%s", ExecID(t.ctx)),
		Image: t.config.ContainerImage,
		Env: map[string]string{
			EnvJsonRpcUrl:      t.config.JsonRpcUrl,
			EnvStartBlock:      startBlock,
			config.EnvLogLevel: t.config.LogLevel,
		},
	})
	if err != nil {
		return err
	}
	t.container = &container
	return nil
}

func (t *TxNodeService) Stop() error {
	log.Infof("Stopping %s", t.Name())
	if t.container != nil {
		return t.client.StopContainer(t.container.ID)
	}
	return nil
}

func (t *TxNodeService) Name() string {
	return "TxNode"
}

func NewTxNodeService(ctx context.Context, cfg TxNodeServiceConfig) (*TxNodeService, error) {
	dockerClient := clients.NewDockerClient()
	return &TxNodeService{
		ctx:    ctx,
		client: dockerClient,
		config: cfg,
	}, nil
}
