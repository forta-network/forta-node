package botauth

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/forta-network/forta-core-go/protocol/settings"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/config"
	log "github.com/sirupsen/logrus"
)

// BotAuthenticator makes sure ip is an assigned bot
type BotAuthenticator struct {
	ctx          context.Context
	server       *http.Server
	dockerClient clients.DockerClient
	msgClient    clients.MessageClient

	agentConfigs  []config.AgentConfig
	agentConfigMu sync.RWMutex
}

func (p *BotAuthenticator) MsgClient() clients.MessageClient {
	return p.msgClient
}


func (p *BotAuthenticator) FindAgentFromRemoteAddr(hostPort string) (*config.AgentConfig, bool) {
	agentContainer, err := p.dockerClient.FindContainerNameFromRemoteAddr(p.ctx, hostPort)
	if err != nil {
		return nil, false
	}

	containerName := agentContainer.Names[0][1:]

	p.agentConfigMu.RLock()
	defer p.agentConfigMu.RUnlock()

	for _, agentConfig := range p.agentConfigs {
		if agentConfig.ContainerName() == containerName {
			return &agentConfig, true
		}
	}

	log.WithFields(
		log.Fields{
			"sourceAddr":    hostPort,
			"containerName": containerName,
		},
	).Warn("could not find agent config for container")

	return nil, false
}

func (p *BotAuthenticator) handleAgentVersionsUpdate(payload messaging.AgentPayload) error {
	p.agentConfigMu.Lock()
	p.agentConfigs = payload
	p.agentConfigMu.Unlock()
	return nil
}

func (p *BotAuthenticator) RegisterMessageHandlers() {
	p.msgClient.Subscribe(messaging.SubjectAgentsVersionsLatest, messaging.AgentsHandler(p.handleAgentVersionsUpdate))
}

func NewBotAuthenticator(ctx context.Context, cfg config.Config) (*BotAuthenticator, error) {
	globalClient, err := clients.NewDockerClient("")
	if err != nil {
		return nil, fmt.Errorf("failed to create the global docker client: %v", err)
	}
	msgClient := messaging.NewClient("bot-auth", fmt.Sprintf("%s:%s", config.DockerNatsContainerName, config.DefaultNatsPort))

	rateLimiting := cfg.JsonRpcProxy.RateLimitConfig
	if rateLimiting == nil {
		rateLimiting = (*config.RateLimitConfig)(settings.GetChainSettings(cfg.ChainID).JsonRpcRateLimiting)
	}

	return &BotAuthenticator{
		ctx:          ctx,
		dockerClient: globalClient,
		msgClient:    msgClient,
	}, nil
}
