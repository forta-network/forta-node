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
	BotPullTimeout     = time.Minute * 10
	BotStartTimeout    = time.Minute * 5
	BotShutdownTimeout = time.Minute

	ImagePullCooldownThreshold = 3
	ImagePullCooldownDuration  = time.Minute * 30
)

// BotClient launches a bot.
type BotClient interface {
	EnsureBotImages(ctx context.Context, botConfigs []config.AgentConfig) []error
	LaunchBot(ctx context.Context, botConfig config.AgentConfig) error
	TearDownBot(ctx context.Context, containerName string) error
	StopBot(ctx context.Context, botConfig config.AgentConfig) error
	LoadBotContainers(ctx context.Context) ([]types.Container, error)
	StartWaitBotContainer(ctx context.Context, containerID string) error
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
) *botClient {
	botImageClient.SetImagePullCooldown(ImagePullCooldownThreshold, ImagePullCooldownDuration)
	return &botClient{
		logConfig:       logConfig,
		resourcesConfig: resourcesConfig,
		client:          client,
		botImageClient:  botImageClient,
	}
}

var _ BotClient = &botClient{}

// EnsureBotImages ensures that all of the bot images are locally available.
func (bc *botClient) EnsureBotImages(ctx context.Context, botConfigs []config.AgentConfig) []error {
	var imagePulls []docker.ImagePull
	for _, botConfig := range botConfigs {
		imagePulls = append(imagePulls, docker.ImagePull{
			Name: botConfig.ID,
			Ref:  botConfig.Image,
		})
	}
	return bc.botImageClient.EnsureLocalImages(ctx, BotPullTimeout, imagePulls)
}

// LaunchBot launches a bot by downloading docker image and starting the container.
// This method can be called when the bot containers are alive and should be able to
// handle that situation.
func (bc *botClient) LaunchBot(ctx context.Context, botConfig config.AgentConfig) error {
	ctx, cancel := context.WithTimeout(ctx, BotStartTimeout)
	defer cancel()

	// first make sure that the bot's bridge network exists
	botNetworkID, err := bc.client.EnsurePublicNetwork(ctx, botConfig.ContainerName())
	if err != nil {
		return fmt.Errorf("error creating public network: %v", err)
	}

	_, err = bc.client.GetContainerByName(ctx, botConfig.ContainerName())
	switch {
	case err == nil:
		// do not create a new container - we already have it

	case errors.Is(err, docker.ErrContainerNotFound):
		// if the bot container doesn't exist, create and start the container
		botContainerCfg := NewBotContainerConfig(botNetworkID, botConfig, bc.logConfig, bc.resourcesConfig)
		_, err = bc.client.StartContainer(ctx, botContainerCfg)
		if err != nil {
			return fmt.Errorf("failed to start bot container: %v", err)
		}

	default:
		return fmt.Errorf("unexpected error while getting the bot container '%s': %v", botConfig.ContainerName(), err)
	}

	// at this point we have created a new bot container and a new bridge network for the bot
	// or found the existing container and the network: it's time to ensure that all service containers
	// are reattached to the bot's network
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
	for _, containerName := range getServiceContainerNames() {
		container, err := bc.client.GetContainerByName(ctx, containerName)
		if err != nil {
			return nil, fmt.Errorf("failed to get service container ids: %v", err)
		}
		ids = append(ids, container.ID)
	}
	return ids, nil
}

func getServiceContainerNames() []string {
	return []string{
		config.DockerScannerContainerName, config.DockerJSONRPCProxyContainerName,
		config.DockerJWTProviderContainerName, config.DockerPublicAPIProxyContainerName,
	}
}

// TearDownBot tears down a bot removing the docker container and network.
func (bc *botClient) TearDownBot(ctx context.Context, containerName string) error {
	container, err := bc.client.GetContainerByName(ctx, containerName)
	if err != nil {
		return fmt.Errorf("failed to get the bot container to tear down: %v", err)
	}
	serviceContainerIDs, err := bc.getServiceContainerIDs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get service container ids during bot cleanup: %v", err)
	}
	defer log.WithField("botContainer", containerName).Info("done tearing down the bot and the associated docker resources")
	// not returning any errors in `if`s below so we keep on by removing whatever is left
	for _, serviceContainerID := range serviceContainerIDs {
		if err := bc.client.DetachNetwork(ctx, serviceContainerID, containerName); err != nil {
			log.WithFields(
				log.Fields{
					"network":          containerName,
					"serviceContainer": serviceContainerID,
				},
			).WithError(err).Warn("failed to detach the service container from the bot network")
		}
	}

	terminateCtx, terminateCancel := context.WithTimeout(ctx, BotShutdownTimeout)
	defer terminateCancel()

	timeout := BotShutdownTimeout

	if err := bc.client.ShutdownContainer(terminateCtx, container.ID, &timeout); err != nil {
		log.WithFields(
			log.Fields{
				"containerId":   container.ID,
				"containerName": containerName,
			},
		).WithError(err).Warn("failed to terminate the bot container")
	}

	if err := bc.client.RemoveContainer(ctx, container.ID); err != nil {
		log.WithFields(
			log.Fields{
				"containerId":   container.ID,
				"containerName": containerName,
			},
		).WithError(err).Warn("failed to remove the bot container")
	}
	if err := bc.client.RemoveNetworkByName(ctx, containerName); err != nil {
		log.WithFields(
			log.Fields{
				"network": containerName,
			},
		).WithError(err).Warn("failed to destroy the bot network")
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
