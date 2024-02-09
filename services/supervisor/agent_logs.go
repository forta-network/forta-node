package supervisor

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func (sup *SupervisorService) syncAgentLogs() {
	interval := time.Duration(sup.botLifecycleConfig.Config.AgentLogsConfig.SendIntervalSeconds) * time.Second
	ticker := time.NewTicker(interval)
	for range ticker.C {
		err := sup.botLifecycle.BotLogger.SendBotLogs(sup.ctx, interval)
		sup.lastAgentLogsRequest.Set()
		sup.lastAgentLogsRequestError.Set(err)
		if err != nil {
			log.WithError(err).Warn("failed to sync agent logs")
		}
	}
}
