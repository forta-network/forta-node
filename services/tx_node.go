package services

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"OpenZeppelin/safe-node/clients"
	"OpenZeppelin/safe-node/config"
)

const EnvJsonRpcUrl = "JSON_RPC_URL"

type TxNodeService struct {
	ctx       context.Context
	client    clients.DockerClient
	container *clients.DockerContainer
	config    TxNodeServiceConfig
}

type TxNodeServiceConfig struct {
	JsonRpcUrl string
	LogLevel   string
}

func (t *TxNodeService) Start() error {
	log.Infof("Starting %s", t.Name())
	_, err := log.ParseLevel(t.config.LogLevel)
	if err != nil {
		log.Error("invalid log level", err)
		return err
	}
	container, err := t.client.StartContainer(t.ctx, clients.DockerContainerConfig{
		Name:  fmt.Sprintf("safe-node-%s", ExecID(t.ctx)),
		Image: "openzeppelin/safe-node",
		Env: map[string]string{
			EnvJsonRpcUrl:      t.config.JsonRpcUrl,
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
