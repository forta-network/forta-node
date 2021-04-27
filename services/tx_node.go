package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"OpenZeppelin/zephyr-node/clients"
	"OpenZeppelin/zephyr-node/config"
)

const ZephyrPrefix = "zephyr"

var nodeName = fmt.Sprintf("%s-node", ZephyrPrefix)
var proxyName = fmt.Sprintf("%s-proxy", ZephyrPrefix)

// TxNodeService manages the safe-node docker container as a service
type TxNodeService struct {
	ctx             context.Context
	client          clients.DockerClient
	nodeContainer   *clients.DockerContainer
	proxyContainer  *clients.DockerContainer
	agentContainers []*clients.DockerContainer
	config          TxNodeServiceConfig
}

type TxNodeServiceConfig struct {
	Config config.Config
}

func (t *TxNodeService) Start() error {
	log.Infof("Starting %s", t.Name())
	_, err := log.ParseLevel(t.config.Config.Log.Level)
	if err != nil {
		log.Error("invalid log level", err)
		return err
	}

	cfgBytes, err := json.Marshal(t.config.Config)
	if err != nil {
		log.Error("cannot marshal config to json", err)
		return err
	}
	cfgJson := string(cfgBytes)

	if err := t.client.Prune(t.ctx); err != nil {
		return err
	}

	var proxyKeys []string
	for range t.config.Config.Agents {
		proxyKeys = append(proxyKeys, uuid.Must(uuid.NewRandom()).String())
	}

	nodeNetwork, err := t.client.CreatePublicNetwork(t.ctx, nodeName)
	if err != nil {
		return err
	}

	// removes ide warning about possible nil slice
	if proxyKeys == nil {
		return errors.New("proxy keys is nil")
	}

	var networkIDs []string
	for i, agent := range t.config.Config.Agents {
		nwID, err := t.client.CreateInternalNetwork(t.ctx, agent.Name)
		if err != nil {
			return err
		}
		networkIDs = append(networkIDs, nwID)
		agentContainer, err := t.client.StartContainer(t.ctx, clients.DockerContainerConfig{
			Name:           agent.ContainerName(),
			Image:          agent.Image,
			NetworkID:      nwID,
			LinkNetworkIDs: []string{},
			Env: map[string]string{
				config.EnvAgentKey:   proxyKeys[i],
				config.EnvAgentProxy: proxyName,
			},
		})
		if err != nil {
			return err
		}
		t.agentContainers = append(t.agentContainers, agentContainer)
	}

	proxyContainer, err := t.client.StartContainer(t.ctx, clients.DockerContainerConfig{
		Name:      proxyName,
		Image:     t.config.Config.Zephyr.ProxyImage,
		NetworkID: nodeNetwork,
		Env: map[string]string{
			config.EnvZephyrConfig: cfgJson,
			config.EnvAgentKeys:    strings.Join(proxyKeys, ","),
		},
		LinkNetworkIDs: networkIDs,
	})
	if err != nil {
		return err
	}
	t.proxyContainer = proxyContainer

	nodeContainer, err := t.client.StartContainer(t.ctx, clients.DockerContainerConfig{
		Name:  nodeName,
		Image: t.config.Config.Zephyr.NodeImage,
		Env: map[string]string{
			config.EnvZephyrConfig: cfgJson,
		},
		NetworkID:      nodeNetwork,
		LinkNetworkIDs: networkIDs,
	})
	if err != nil {
		return err
	}

	t.nodeContainer = nodeContainer
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
	if t.proxyContainer != nil {
		if err := t.client.StopContainer(ctx, t.proxyContainer.ID); err != nil {
			log.Error("error stopping node container", err)
		} else {
			log.Infof("Container %s is stopped", t.proxyContainer.ID)
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
