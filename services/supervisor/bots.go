package supervisor

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// refreshBotContainers refreshes bot containers every 15 seconds.
// This allows us to blast the latest assignment list very often
// and keep bot containers and clients in order.
func (sup *SupervisorService) refreshBotContainers() {
	for {
		select {
		case <-sup.ctx.Done():
			return

		case <-time.After(time.Second * 15):
			if err := sup.botLifecycle.BotManager.ManageBots(sup.ctx); err != nil {
				log.WithError(err).Error("error while managing bots")
			}
			if err := sup.botLifecycle.BotManager.RestartExitedBots(sup.ctx); err != nil {
				log.WithError(err).Error("error while restarting exited bots")
			}
		}
	}
}
