package services

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"OpenZeppelin/fortify-node/clients"
	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/store"
)

const ContainerPrefix = "fortify"

var scannerName = fmt.Sprintf("%s-scanner", ContainerPrefix)
var jsonRpcProxyName = fmt.Sprintf("%s-json-rpc", ContainerPrefix)
var queryName = fmt.Sprintf("%s-query", ContainerPrefix)

// TxNodeService manages the safe-node docker container as a service
type TxNodeService struct {
	ctx        context.Context
	client     clients.DockerClient
	containers []*clients.DockerContainer
	config     TxNodeServiceConfig
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

	nodeNetwork, err := t.client.CreatePublicNetwork(t.ctx, scannerName)
	if err != nil {
		return err
	}

	var networkIDs []string
	for _, agent := range t.config.Config.Agents {
		nwID, err := t.client.CreatePublicNetwork(t.ctx, agent.Name)
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
				config.EnvJsonRpcHost: jsonRpcProxyName,
				config.EnvJsonRpcPort: "8545",
			},
		})
		if err != nil {
			return err
		}
		t.containers = append(t.containers, agentContainer)
	}

	queryContainer, err := t.client.StartContainer(t.ctx, clients.DockerContainerConfig{
		Name:  queryName,
		Image: t.config.Config.Query.QueryImage,
		Env: map[string]string{
			config.EnvFortifyConfig: cfgJson,
		},
		Ports: map[string]string{
			fmt.Sprintf("%d", t.config.Config.Query.Port): "80",
		},
		Volumes: map[string]string{
			t.config.Config.Query.DB.Path: store.DBPath,
		},
		NetworkID: nodeNetwork,
	})
	if err != nil {
		return err
	}

	jsonRpcContainer, err := t.client.StartContainer(t.ctx, clients.DockerContainerConfig{
		Name:  jsonRpcProxyName,
		Image: t.config.Config.JsonRpcProxy.JsonRpcImage,
		Env: map[string]string{
			config.EnvFortifyConfig: cfgJson,
		},
		NetworkID:      nodeNetwork,
		LinkNetworkIDs: networkIDs,
	})
	if err != nil {
		return err
	}

	scannerContainer, err := t.client.StartContainer(t.ctx, clients.DockerContainerConfig{
		Name:  scannerName,
		Image: t.config.Config.Scanner.ScannerImage,
		Env: map[string]string{
			config.EnvFortifyConfig: cfgJson,
			config.EnvQueryNode:     queryName,
		},
		NetworkID:      nodeNetwork,
		LinkNetworkIDs: networkIDs,
	})
	if err != nil {
		return err
	}

	t.containers = append(t.containers, jsonRpcContainer, scannerContainer, queryContainer)

	return nil
}

func (t *TxNodeService) Stop() error {
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
	dockerClient := clients.NewDockerClient()
	return &TxNodeService{
		ctx:    ctx,
		client: dockerClient,
		config: cfg,
	}, nil
}
