package components

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/forta-network/forta-core-go/clients/agentlogs"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/agentgrpc"
	"github.com/forta-network/forta-node/clients/docker"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/services/components/botio"
	"github.com/forta-network/forta-node/services/components/botio/botreq"
	"github.com/forta-network/forta-node/services/components/containers"
	"github.com/forta-network/forta-node/services/components/lifecycle"
	"github.com/forta-network/forta-node/services/components/lifecycle/mediator"
	"github.com/forta-network/forta-node/services/components/metrics"
	"github.com/forta-network/forta-node/services/components/registry"
)

// BotProcessingConfig contains bot processing component configuration and dependencies.
type BotProcessingConfig struct {
	Config        config.Config
	MessageClient clients.MessageClient
}

// BotProcessing contains the bot processing components.
type BotProcessing struct {
	RequestSender botio.Sender
	Results       botreq.ReceiveOnlyChannels
}

// GetBotProcessingComponents returns the bot processing components after doing dependency injection.
func GetBotProcessingComponents(ctx context.Context, botProcCfg BotProcessingConfig) (BotProcessing, error) {
	resultChannels := botreq.MakeResultChannels()
	lifecycleMetrics := metrics.NewLifecycleClient(botProcCfg.MessageClient)
	botClientFactory := botio.NewBotClientFactory(
		resultChannels.SendOnly(), botProcCfg.MessageClient,
		lifecycleMetrics, agentgrpc.NewBotDialer(),
	)
	botPool := lifecycle.NewBotPool(
		ctx, lifecycleMetrics, botClientFactory, botProcCfg.Config.BotsToWait(),
	)
	mediator.New(botProcCfg.MessageClient, lifecycleMetrics).ConnectBotPool(botPool)

	// update the bot pool directly if we are in standalone mode
	if botProcCfg.Config.LocalModeConfig.IsStandalone() {
		botRegistry, err := registry.New(botProcCfg.Config, common.HexToAddress(utils.ZeroAddress))
		if err != nil {
			return BotProcessing{}, fmt.Errorf("failed to create the standalone mode registry: %v", err)
		}
		bots, err := botRegistry.LoadAssignedBots()
		if err != nil {
			return BotProcessing{}, fmt.Errorf("failed to load the standalone mode bots: %v", err)
		}
		if err := botPool.UpdateBotsWithLatestConfigs(bots); err != nil {
			return BotProcessing{}, fmt.Errorf("failed to update the standalone mode bot pool: %v", err)
		}
	}

	sender := botio.NewSender(ctx, botProcCfg.MessageClient, botPool)
	return BotProcessing{
		RequestSender: sender,
		Results:       resultChannels.ReceiveOnly(),
	}, nil
}

// BotLifecycleConfig contains bot lifecycle component configuration and dependencies.
type BotLifecycleConfig struct {
	Config         config.Config
	ScannerAddress common.Address
	MessageClient  clients.MessageClient
	BotRegistry    registry.BotRegistry
	Key            *keystore.Key
}

// BotLifecycle contains the bot lifecycle components.
type BotLifecycle struct {
	BotManager   lifecycle.BotLifecycleManager
	BotClient    containers.BotClient
	ImageCleanup containers.ImageCleanup
	BotLogger    lifecycle.BotLogger
}

// GetBotLifecycleComponents returns the bot lifecycle management components.
func GetBotLifecycleComponents(
	ctx context.Context,
	botLifeConfig BotLifecycleConfig,
) (BotLifecycle, error) {
	cfg := botLifeConfig.Config
	// bot image client is helpful for loading local mode agents from a restricted container registry
	var (
		botImageClient clients.DockerClient
		err            error
	)
	if cfg.LocalModeConfig.Enable && cfg.LocalModeConfig.ContainerRegistry != nil {
		botImageClient, err = docker.NewAuthDockerClient(
			"",
			cfg.LocalModeConfig.ContainerRegistry.Username,
			cfg.LocalModeConfig.ContainerRegistry.Password,
		)
	} else {
		botImageClient, err = docker.NewDockerClient("")
	}
	if err != nil {
		return BotLifecycle{}, fmt.Errorf("failed to create the bot image docker client: %v", err)
	}

	dockerClient, err := docker.NewDockerClient(containers.LabelFortaSupervisor)
	if err != nil {
		return BotLifecycle{}, fmt.Errorf("failed to create the bot docker client: %v", err)
	}

	botClient := containers.NewBotClient(
		botLifeConfig.Config.Log, botLifeConfig.Config.ResourcesConfig,
		dockerClient, botImageClient,
	)
	lifecycleMetrics := metrics.NewLifecycleClient(botLifeConfig.MessageClient)
	lifecycleMediator := mediator.New(botLifeConfig.MessageClient, lifecycleMetrics)
	botMonitor := lifecycle.NewBotMonitor(lifecycleMetrics)
	lifecycleMediator.ConnectBotMonitor(botMonitor)
	botManager := lifecycle.NewManager(
		botLifeConfig.BotRegistry, botClient, lifecycleMediator,
		lifecycleMetrics, botMonitor,
	)
	imageCleanup := containers.NewImageCleanup(dockerClient, botLifeConfig.BotRegistry)
	botLogger := lifecycle.NewBotLogger(
		botClient,
		dockerClient,
		botLifeConfig.Key,
		agentlogs.NewClient("").SendLogs,
	)

	return BotLifecycle{
		BotManager:   botManager,
		BotClient:    botClient,
		ImageCleanup: imageCleanup,
		BotLogger:    botLogger,
	}, nil
}
