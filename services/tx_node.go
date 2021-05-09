package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"OpenZeppelin/fortify-node/clients"
	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/store"
)

const ContainerPrefix = "fortify"

var scannerName = fmt.Sprintf("%s-scanner", ContainerPrefix)
var queryName = fmt.Sprintf("%s-query", ContainerPrefix)

// TxNodeService manages the safe-node docker container as a service
type TxNodeService struct {
	ctx              context.Context
	client           clients.DockerClient
	scannerContainer *clients.DockerContainer
	queryContainer   *clients.DockerContainer
	agentContainers  []*clients.DockerContainer
	config           TxNodeServiceConfig
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

	nodeNetwork, err := t.client.CreatePublicNetwork(t.ctx, scannerName)
	if err != nil {
		return err
	}

	// removes ide warning about possible nil slice
	if proxyKeys == nil {
		return errors.New("proxy keys is nil")
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
		})
		if err != nil {
			return err
		}
		t.agentContainers = append(t.agentContainers, agentContainer)
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
	t.queryContainer = queryContainer

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

	t.scannerContainer = scannerContainer
	return nil
}

func (t *TxNodeService) Stop() error {
	ctx := context.Background()
	if t.scannerContainer != nil {
		if err := t.client.StopContainer(ctx, t.scannerContainer.ID); err != nil {
			log.Error("error stopping node container", err)
		} else {
			log.Infof("Container %s is stopped", t.scannerContainer.ID)
		}
	}
	if t.queryContainer != nil {
		if err := t.client.StopContainer(ctx, t.queryContainer.ID); err != nil {
			log.Error("error stopping query container", err)
		} else {
			log.Infof("Container %s is stopped", t.queryContainer.ID)
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
