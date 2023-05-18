package supervisor

import (
	"fmt"
	"strconv"
	"time"

	"github.com/forta-network/forta-core-go/clients/agentlogs"
	"github.com/forta-network/forta-core-go/security"
	"github.com/forta-network/forta-node/clients/docker"
	log "github.com/sirupsen/logrus"
)

// adjust these better with auto-upgrade later
const (
	defaultAgentLogSendInterval       = time.Minute
	defaultAgentLogTailLines          = 50
	defaultAgentLogAvgMaxCharsPerLine = 200
)

func (sup *SupervisorService) syncAgentLogs() {
	ticker := time.NewTicker(defaultAgentLogSendInterval)
	for range ticker.C {
		err := sup.doSyncAgentLogs()
		sup.lastAgentLogsRequest.Set()
		sup.lastAgentLogsRequestError.Set(err)
		if err != nil {
			log.WithError(err).Warn("failed to sync agent logs")
		}
	}
}

func (sup *SupervisorService) doSyncAgentLogs() error {
	sup.mu.RLock()
	defer sup.mu.RUnlock()

	var (
		sendLogs agentlogs.Agents
		keepLogs agentlogs.Agents
	)
	for _, container := range sup.containers {
		if !container.IsAgent {
			continue
		}
		dockerContainer, err := sup.client.GetContainerByID(sup.ctx, container.ID)
		if err != nil {
			log.WithError(err).Warn("failed to get agent container")
			continue
		}
		if dockerContainer.Labels[docker.LabelFortaSettingsAgentLogsEnable] != "true" {
			continue
		}
		logs, err := sup.client.GetContainerLogs(
			sup.ctx, container.ID,
			strconv.Itoa(defaultAgentLogTailLines),
			defaultAgentLogAvgMaxCharsPerLine*defaultAgentLogTailLines,
		)
		if err != nil {
			log.WithError(err).Warn("failed to get agent container logs")
			continue
		}
		agent := &agentlogs.Agent{
			ID:   container.AgentConfig.ID,
			Logs: logs,
		}
		// don't send if it's the same with previous logs but keep it for next time
		// so we can check
		keepLogs = append(keepLogs, agent)
		if !sup.prevAgentLogs.Has(container.AgentConfig.ID, logs) {
			log.WithField("agent", agent.ID).Debug("new agent logs found")
			sendLogs = append(sendLogs, agent)
		} else {
			log.WithField("agent", agent.ID).Debug("no new agent logs")
		}
	}

	if len(sendLogs) > 0 {
		scannerJwt, err := security.CreateScannerJWT(sup.config.Key, map[string]interface{}{
			"access": "agent_logs",
		})
		if err != nil {
			return fmt.Errorf("failed to create scanner token: %v", err)
		}
		if err := sup.agentLogsClient.SendLogs(sendLogs, scannerJwt); err != nil {
			return fmt.Errorf("failed to send agent logs: %v", err)
		}
		log.WithField("count", len(sendLogs)).Debug("successfully sent new agent logs")
	} else {
		log.Debug("no new agent logs were found - not sending")
	}

	sup.prevAgentLogs = keepLogs
	return nil
}
