package supervisor

import (
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

func (sup *SupervisorService) syncAgentLogs() {
	interval := time.Duration(sup.botLifecycleConfig.Config.AgentLogsConfig.SendIntervalSeconds) * time.Second
	ticker := time.NewTicker(interval)
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

	if err := sup.botLifecycle.botLogger.SendBotLogs(sup.ctx); err != nil {
		log.WithError(err).Error("error while sending bot logs")
		return err
	}

	return nil
}
