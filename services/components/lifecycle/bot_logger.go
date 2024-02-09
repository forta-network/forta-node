package lifecycle

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/forta-network/forta-core-go/clients/agentlogs"
	"github.com/forta-network/forta-core-go/security"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/docker"
	"github.com/forta-network/forta-node/services/components/containers"
	"github.com/forta-network/forta-node/services/components/registry"
	log "github.com/sirupsen/logrus"
)

// BotLogger manages bots logging.
type BotLogger interface {
	SendBotLogs(ctx context.Context, snapshotInterval time.Duration) error
}

type botLogger struct {
	botClient     containers.BotClient
	dockerClient  clients.DockerClient
	agentRegistry registry.BotRegistry
	key           *keystore.Key
	prevAgentLogs agentlogs.Agents

	sendAgentLogs func(agents agentlogs.Agents, authToken string) error
}

var _ BotLogger = &botLogger{}

func NewBotLogger(
	botClient containers.BotClient,
	dockerClient clients.DockerClient,
	agentRegistry registry.BotRegistry,
	key *keystore.Key,
	sendAgentLogs func(agents agentlogs.Agents, authToken string) error,
) *botLogger {
	return &botLogger{
		botClient:     botClient,
		dockerClient:  dockerClient,
		agentRegistry: agentRegistry,
		key:           key,
		sendAgentLogs: sendAgentLogs,
	}
}

// adjust these better with auto-upgrade later
const (
	defaultAgentLogTailLines = 300
)

func (bl *botLogger) SendBotLogs(ctx context.Context, snapshotInterval time.Duration) error {
	var (
		sendLogs agentlogs.Agents
		keepLogs agentlogs.Agents
	)

	botContainers, err := bl.botClient.LoadBotContainers(ctx)
	if err != nil {
		return fmt.Errorf("failed to load the bot containers: %v", err)
	}

	for _, container := range botContainers {
		if container.Labels[docker.LabelFortaSettingsAgentLogsEnable] != "true" {
			continue
		}
		logs, err := bl.dockerClient.GetContainerLogs(
			ctx, container.ID,
			fmt.Sprintf("%ds", int64(snapshotInterval.Seconds())),
			defaultAgentLogTailLines,
		)
		if err != nil {
			log.WithError(err).Warn("failed to get agent container logs")
			continue
		}

		if len(logs) == 0 {
			log.WithField("agent", container.Labels[docker.LabelFortaBotID]).Debug("no logs found for agent")
			continue
		}

		agentID := container.Labels[docker.LabelFortaBotID]
		agent := &agentlogs.Agent{
			ID:   agentID,
			Logs: logs,
		}

		agentConfig, err := bl.agentRegistry.GetConfigByID(agentID)
		if err == nil {
			agent.ChainID = int64(agentConfig.ChainID)
		}

		// don't send if it's the same with previous logs but keep it for next time
		// so we can check
		keepLogs = append(keepLogs, agent)
		if !bl.prevAgentLogs.Has(agent.ID, logs) {
			log.WithField("agent", agent.ID).Debug("new agent logs found")
			sendLogs = append(sendLogs, agent)
		} else {
			log.WithField("agent", agent.ID).Debug("no new agent logs")
		}
	}

	if len(sendLogs) > 0 {
		scannerJwt, err := security.CreateScannerJWT(bl.key, map[string]interface{}{
			"access": "agent_logs",
		})
		if err != nil {
			return fmt.Errorf("failed to create scanner token: %v", err)
		}
		if err := bl.sendAgentLogs(sendLogs, scannerJwt); err != nil {
			return fmt.Errorf("failed to send agent logs: %v", err)
		}
		log.WithField("count", len(sendLogs)).Debug("successfully sent new agent logs")
	} else {
		log.Debug("no new agent logs were found - not sending")
	}

	bl.prevAgentLogs = keepLogs
	return nil
}
