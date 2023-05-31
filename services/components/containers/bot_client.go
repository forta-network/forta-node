package containers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/docker"
	"github.com/forta-network/forta-node/config"
	log "github.com/sirupsen/logrus"
)

// Timeouts
const (
	BotPullTimeout  = time.Minute * 5
	BotStartTimeout = time.Minute * 5
)

// BotClient launches a bot.
type BotClient interface {
	EnsureBotImages(ctx context.Context, botConfigs []config.AgentConfig) []error
	LaunchBot(ctx context.Context, botConfig config.AgentConfig) error
	TearDownBot(ctx context.Context, botConfig config.AgentConfig) error
	StopBot(ctx context.Context, botConfig config.AgentConfig) error
	LoadBotContainers(ctx context.Context) ([]types.Container, error)
	StartWaitBotContainer(ctx context.Context, containerID string) error
	CleanupUnusedBots(ctx context.Context, botConfigs []config.AgentConfig) error
}

type botClient struct {
	logConfig       config.LogConfig
	resourcesConfig config.ResourcesConfig
	client          clients.DockerClient
	botImageClient  clients.DockerClient
}

// NewBotClient creates a new bot client to manage bot containers.
func NewBotClient(
	logConfig config.LogConfig, resourcesConfig config.ResourcesConfig,
	client clients.DockerClient, botImageClient clients.DockerClient,
) BotClient {
	return &botClient{
		logConfig:       logConfig,
		resourcesConfig: resourcesConfig,
		client:          client,
		botImageClient:  botImageClient,
	}
}

// EnsureBotImages ensures that all of the bot images are locally available.
func (bc *botClient) EnsureBotImages(ctx context.Context, botConfigs []config.AgentConfig) []error {
	var imagePulls []docker.ImagePull
	for _, botConfig := range botConfigs {
		imagePulls = append(imagePulls, docker.ImagePull{
			Name: botConfig.ID,
			Ref:  botConfig.Image,
		})
	}
	return bc.client.EnsureLocalImages(ctx, BotPullTimeout, imagePulls)
}

// LaunchBot launches a bot by downloading docker image and starting the container.
func (bc *botClient) LaunchBot(ctx context.Context, botConfig config.AgentConfig) error {
	ctx, cancel := context.WithTimeout(ctx, BotStartTimeout)
	defer cancel()

	_, err := bc.client.GetContainerByName(ctx, botConfig.ContainerName())
	if !errors.Is(err, docker.ErrContainerNotFound) {
		log.WithField("container", botConfig.ContainerName()).Info("bot container exists - skipping launch")
		return nil
	}

	botNetworkID, err := bc.client.CreatePublicNetwork(ctx, botConfig.ContainerName())
	if err != nil {
		return fmt.Errorf("error creating public network: %v", err)
	}

	botContainerCfg := NewBotContainerConfig(botNetworkID, botConfig, bc.logConfig, bc.resourcesConfig)
	_, err = bc.client.StartContainer(ctx, botContainerCfg)
	if err != nil {
		return fmt.Errorf("failed to start bot container: %v", err)
	}
	return bc.attachServiceContainers(ctx, botNetworkID)
}

func (bc *botClient) attachServiceContainers(ctx context.Context, botNetworkID string) error {
	serviceContainerIDs, err := bc.getServiceContainerIDs(ctx)
	if err != nil {
		return err
	}
	for _, serviceContainerID := range serviceContainerIDs {
		err := bc.client.AttachNetwork(ctx, serviceContainerID, botNetworkID)
		if err != nil {
			return fmt.Errorf(
				"failed to attach service container '%s' to bot network '%s': %v",
				serviceContainerID, botNetworkID, err,
			)
		}
	}
	return nil
}

func (bc *botClient) getServiceContainerIDs(ctx context.Context) (ids []string, err error) {
	for _, containerName := range []string{
		config.DockerScannerContainerName, config.DockerJSONRPCProxyContainerName,
		config.DockerJWTProviderContainerName, config.DockerPublicAPIProxyContainerName,
	} {
		container, err := bc.client.GetContainerByName(ctx, containerName)
		if errors.Is(err, docker.ErrContainerNotFound) {
			return nil, fmt.Errorf("failed to get service container ids: %v", err)
		}
		ids = append(ids, container.ID)
	}
	return ids, nil
}

// TearDownBot tears down a bot by shutting down the docker container and removing it.
func (bc *botClient) TearDownBot(ctx context.Context, botConfig config.AgentConfig) error {
	return bc.tearDownBotWithContainerName(ctx, botConfig.ContainerName())
}

func (bc *botClient) tearDownBotWithContainerName(ctx context.Context, containerName string) error {
	container, err := bc.client.GetContainerByName(ctx, containerName)
	if err != nil {
		return fmt.Errorf("failed to get the bot container to tear down: %v", err)
	}
	serviceContainerIDs, err := bc.getServiceContainerIDs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get service container ids during bot cleanup: %v", err)
	}
	for _, serviceContainerID := range serviceContainerIDs {
		if err := bc.client.DetachNetwork(ctx, serviceContainerID, containerName); err != nil {
			log.WithFields(log.Fields{
				"network":          containerName,
				"serviceContainer": serviceContainerID,
			}).WithError(err).Warn("failed to detach the service container from the bot network")
		}
	}
	if err := bc.client.RemoveContainer(ctx, container.ID); err != nil {
		log.WithFields(log.Fields{
			"containerId":   container.ID,
			"containerName": containerName,
		}).WithError(err).Warn("failed to destroy the bot container")
	}
	if err := bc.client.RemoveNetworkByName(ctx, containerName); err != nil {
		log.WithFields(log.Fields{
			"network": containerName,
		}).WithError(err).Warn("failed to destroy the bot network")
	}
	return nil
}

// StopBot shuts down a bot container.
func (bc *botClient) StopBot(ctx context.Context, botConfig config.AgentConfig) error {
	container, err := bc.client.GetContainerByName(ctx, botConfig.ContainerName())
	if err != nil {
		return fmt.Errorf("failed to get the bot container to stop: %v", err)
	}
	if err := bc.client.StopContainer(ctx, container.ID); err != nil {
		return fmt.Errorf("failed to stop the container: %v", err)
	}
	return nil
}

// LoadBotContainers loads the latest bot list for the running scanner.
func (bc *botClient) LoadBotContainers(ctx context.Context) ([]types.Container, error) {
	return bc.client.GetContainersByLabel(ctx, docker.LabelFortaIsBot, LabelValueFortaIsBot)
}

// StartWaitBotContainer starts the bot container and waits.
func (bc *botClient) StartWaitBotContainer(ctx context.Context, containerID string) error {
	if err := bc.client.StartContainerWithID(ctx, containerID); err != nil {
		return fmt.Errorf("failed to start container with id: %v", err)
	}
	return bc.client.WaitContainerStart(ctx, containerID)
}

// CleanupUnusedBots cleans up unused bot containers and networks.
func (bc *botClient) CleanupUnusedBots(ctx context.Context, botConfigs []config.AgentConfig) error {
	if len(botConfigs) == 0 {
		return nil
	}

	botContainers, err := bc.LoadBotContainers(ctx)
	if err != nil {
		return fmt.Errorf("failed to load bot containers during bot cleanup: %v", err)
	}

	if len(botContainers) == 0 {
		return nil
	}

	for _, botContainer := range botContainers {
		botContainerName := botContainer.Names[0][1:]
		if hasBotContainer(botConfigs, botContainerName) {
			continue
		}

		if err := bc.tearDownBotWithContainerName(ctx, botContainerName); err != nil {
			log.WithField("botContainer", botContainerName).WithError(err).
				Error("error while tearing down the unused bot")
		}
	}

	return nil
}

func hasBotContainer(botList []config.AgentConfig, containerName string) bool {
	for _, currBot := range botList {
		if containerName == currBot.ContainerName() {
			return true
		}
	}
	return false
}
