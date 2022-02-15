package supervisor

import (
	"fmt"
	"strconv"
	"time"

	"github.com/forta-protocol/forta-node/clients/agentlogs"
	log "github.com/sirupsen/logrus"
)

// adjust these better with auto-upgrade later
const (
	defaultAgentLogSendInterval = time.Minute
	defaultAgentLogTailLines    = 20
	defaultAgentLogMaxLength    = 4096
)

func (sup *SupervisorService) syncAgentLogs() {
	time.After(time.Second * 15) // rate limit crash loops
	ticker := time.NewTicker(defaultAgentLogSendInterval)
	for range ticker.C {
		err := sup.doSyncAgentLogs()
		if err != nil {
			log.WithError(err).Warn("failed to sync agent logs")
		}
		sup.lastAgentLogsRequest.Set()
		sup.lastAgentLogsRequestError.Set(err)
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
		logs, err := sup.client.GetContainerLogs(sup.ctx, container.ID, strconv.Itoa(defaultAgentLogTailLines), defaultAgentLogMaxLength)
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
			sendLogs = append(sendLogs, agent)
		}
	}

	if err := sup.agentLogsClient.SendLogs(sendLogs); err != nil {
		return fmt.Errorf("failed to send agent logs: %v", err)
	}

	sup.prevAgentLogs = keepLogs
	return nil
}
