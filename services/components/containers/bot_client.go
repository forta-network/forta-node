package containers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/errdefs"
	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/docker"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/services/components/metrics"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// Timeouts
const (
	BotPullTimeout     = time.Minute * 10
	BotStartTimeout    = time.Minute * 5
	BotShutdownTimeout = time.Minute

	ImagePullCooldownThreshold = 5
	ImagePullCooldownDuration  = time.Hour * 4

	DockerResourcesPollingInterval = time.Second * 10
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
	logConfig        config.LogConfig
	resourcesConfig  config.ResourcesConfig
	tokenExchangeURL string
	client           clients.DockerClient
	botImageClient   clients.DockerClient
	msgClient        clients.MessageClient
}

// NewBotClient creates a new bot client to manage bot containers.
func NewBotClient(
	logConfig config.LogConfig, resourcesConfig config.ResourcesConfig,
	tokenExchangeURL string,
	client clients.DockerClient, botImageClient clients.DockerClient,
	msgClient clients.MessageClient,
) *botClient {
	botImageClient.SetImagePullCooldown(ImagePullCooldownThreshold, ImagePullCooldownDuration)
	return &botClient{
		logConfig:        logConfig,
		resourcesConfig:  resourcesConfig,
		tokenExchangeURL: tokenExchangeURL,
		client:           client,
		botImageClient:   botImageClient,
		msgClient:        msgClient,
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
		botContainerCfg := NewBotContainerConfig(
			botNetworkID, botConfig, bc.logConfig, bc.resourcesConfig, bc.tokenExchangeURL,
		)
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
	err = bc.attachServiceContainers(ctx, botNetworkID)
	if err != nil {
		return fmt.Errorf("failed to attach service containers to the bot network: %v", err)
	}

	go bc.pollDockerResources(botConfig.ContainerName(), botConfig)

	return nil
}

// pollDockerResources polls docker resources for bot container and sends them to the publisher.
func (bc *botClient) pollDockerResources(containerID string, agentConfig config.AgentConfig) {
	ctx := context.Background()
	ticker := initTicker(DockerResourcesPollingInterval)
	defer ticker.Stop()

	var previousResources docker.ContainerResources

	for t := range ticker.C {
		logrus.WithField("containerID", containerID).Debug("polling docker resources")
		// request docker stats
		resources, err := bc.client.ContainerStats(ctx, containerID)
		if errdefs.IsNotFound(err) {
			logrus.WithError(err).
				WithField("containerID", containerID).
				WithField("agentID", agentConfig.ID).
				Warn("bot container can't be found, stopping docker resources poller")
			return
		} else if err != nil {
			logrus.WithError(err).Error("error while getting container stats", containerID)
			continue
		}

		var (
			bytesSent uint64
			bytesRecv uint64
		)

		for _, network := range resources.NetworkStats {
			bytesSent += network.TxBytes
			bytesRecv += network.RxBytes
		}

		logrus.WithField("containerID", containerID).
			WithField("resources", resources).
			Debug("sending docker resources metrics")

		cpuPercent := calculateCPUPercentUnix(&previousResources, resources)

		metrics.SendAgentMetrics(bc.msgClient, []*protocol.AgentMetric{
			metrics.CreateAgentResourcesMetric(
				agentConfig, t, domain.MetricDockerResourcesCPU, cpuPercent),
			metrics.CreateAgentResourcesMetric(
				agentConfig, t, domain.MetricDockerResourcesMemory, float64(resources.MemoryStats.Usage)),
			metrics.CreateAgentResourcesMetric(
				agentConfig, t, domain.MetricDockerResourcesNetworkSent, float64(bytesSent)),
			metrics.CreateAgentResourcesMetric(
				agentConfig, t, domain.MetricDockerResourcesNetworkReceive, float64(bytesRecv)),
		})

		previousResources = *resources
	}
}

func calculateCPUPercentUnix(previousResources, resources *docker.ContainerResources) float64 {
	var (
		cpuPercent = 0.0
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(resources.CPUStats.CPUUsage.TotalUsage) - float64(previousResources.CPUStats.CPUUsage.TotalUsage)
		// calculate the change for the entire system between readings
		systemDelta = float64(resources.CPUStats.SystemUsage) - float64(previousResources.CPUStats.SystemUsage)
	)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(resources.CPUStats.OnlineCPUs) * 100.0
	}
	return cpuPercent
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

func initTicker(interval time.Duration) *time.Ticker {
	nextTick := time.Now().Truncate(interval).Add(interval)
	initialSleepDuration := time.Until(nextTick)

	// Sleep until the next interval
	time.Sleep(initialSleepDuration)

	// Start a ticker that ticks every interval
	ticker := time.NewTicker(interval)

	return ticker
}
