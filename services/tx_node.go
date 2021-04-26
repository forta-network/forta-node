package services

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"OpenZeppelin/zephyr-node/clients"
	"OpenZeppelin/zephyr-node/config"
)

const EnvJsonRpcUrl = "JSON_RPC_URL"
const EnvStartBlock = "START_BLOCK"
const EnvAgents = "AGENTS"

const ZephyrPrefix = "zephyr"

// TxNodeService manages the safe-node docker container as a service
type TxNodeService struct {
	ctx             context.Context
	client          clients.DockerClient
	nodeContainer   *clients.DockerContainer
	agentContainers []*clients.DockerContainer
	config          TxNodeServiceConfig
}

type TxNodeServiceConfig struct {
	JsonRpcUrl     string
	LogLevel       string
	ContainerImage string
	StartBlock     int
	AgentConfigs   []config.AgentConfig
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

	if err := t.client.Prune(t.ctx); err != nil {
		return err
	}

	var networkIDs []string
	var agentNames []string
	for _, agent := range t.config.AgentConfigs {
		nwID, err := t.client.CreateInternalNetwork(t.ctx, agent.Name)
		if err != nil {
			return err
		}
		networkIDs = append(networkIDs, nwID)
		name := fmt.Sprintf("%s-agent-%s", ZephyrPrefix, agent.Name)
		agentNames = append(agentNames, name)
		agentContainer, err := t.client.StartContainer(t.ctx, clients.DockerContainerConfig{
			Name:      name,
			Image:     agent.Image,
			NetworkID: nwID,
		})
		if err != nil {
			return err
		}
		t.agentContainers = append(t.agentContainers, agentContainer)
	}

	nwID, err := t.client.CreatePublicNetwork(t.ctx, fmt.Sprintf("%s-node", ZephyrPrefix))
	if err != nil {
		return err
	}
	container, err := t.client.StartContainer(t.ctx, clients.DockerContainerConfig{
		Name:  "zephyr-node",
		Image: t.config.ContainerImage,
		Env: map[string]string{
			EnvJsonRpcUrl:      t.config.JsonRpcUrl,
			EnvStartBlock:      startBlock,
			EnvAgents:          strings.Join(agentNames, ","),
			config.EnvLogLevel: t.config.LogLevel,
		},
		NetworkID:      nwID,
		LinkNetworkIDs: networkIDs,
	})
	if err != nil {
		return err
	}

	t.nodeContainer = container
	return nil
}

func (t *TxNodeService) Stop() error {
	log.Infof("Stopping %s", t.Name())
	ctx := context.Background()
	if t.nodeContainer != nil {
		if err := t.client.StopContainer(ctx, t.nodeContainer.ID); err != nil {
			log.Error("error stopping node container", err)
		} else {
			log.Infof("Container %s is stopped", t.nodeContainer.ID)
		}
	}
	for _, agt := range t.agentContainers {
		if err := t.client.StopContainer(ctx, agt.ID); err != nil {
			log.Error(fmt.Sprintf("error stopping %s container", agt.ID), err)
		} else {
			log.Infof("Container %s is stopped", agt.ID)
		}
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
