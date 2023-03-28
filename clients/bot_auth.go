package clients

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/forta-network/forta-core-go/protocol/settings"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/config"
)

// BotAuthenticator makes sure ip is an assigned bot
type botAuthenticator struct {
	ctx          context.Context
	server       *http.Server
	dockerClient DockerClient
	msgClient    MessageClient

	agentConfigs  []config.AgentConfig
	agentConfigMu sync.RWMutex
}

func (p *botAuthenticator) FindAgentFromRemoteAddr(hostPort string) (*config.AgentConfig, error) {
	agentContainer, err := p.dockerClient.GetContainerFromRemoteAddr(p.ctx, hostPort)
	if err != nil {
		return nil, err
	}

	containerName := agentContainer.Names[0][1:]

	p.agentConfigMu.RLock()
	defer p.agentConfigMu.RUnlock()

	for _, agentConfig := range p.agentConfigs {
		if agentConfig.ContainerName() == containerName {
			return &agentConfig, nil
		}
	}

	return nil, err
}

func (p *botAuthenticator) handleAgentVersionsUpdate(payload messaging.AgentPayload) error {
	p.agentConfigMu.Lock()
	p.agentConfigs = payload
	p.agentConfigMu.Unlock()
	return nil
}

func NewBotAuthenticator(ctx context.Context, cfg config.Config) (BotAuthenticator, error) {
	globalClient, err := NewDockerClient("")
	if err != nil {
		return nil, fmt.Errorf("failed to create the global docker client: %v", err)
	}
	msgClient := messaging.NewClient("bot-auth", fmt.Sprintf("%s:%s", config.DockerNatsContainerName, config.DefaultNatsPort))


	rateLimiting := cfg.JsonRpcProxy.RateLimitConfig
	if rateLimiting == nil {
		rateLimiting = (*config.RateLimitConfig)(settings.GetChainSettings(cfg.ChainID).JsonRpcRateLimiting)
	}

	b := &botAuthenticator{
		ctx:          ctx,
		dockerClient: globalClient,
		msgClient:    msgClient,
	}

	msgClient.Subscribe(messaging.SubjectAgentsVersionsLatest, messaging.AgentsHandler(b.handleAgentVersionsUpdate))

	return b, nil
}
