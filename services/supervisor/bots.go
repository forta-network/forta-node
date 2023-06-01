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

		case <-time.After(time.Minute):
			sup.doRefreshBotContainers()
		}
	}
}

func (sup *SupervisorService) doRefreshBotContainers() {
	if err := sup.botLifecycle.BotManager.ManageBots(sup.ctx); err != nil {
		log.WithError(err).Error("error while managing bots")
	}
	if err := sup.botLifecycle.BotManager.CleanupUnusedBots(sup.ctx); err != nil {
		log.WithError(err).Error("error while cleaning up unused bots")
	}
	if err := sup.botLifecycle.BotManager.RestartExitedBots(sup.ctx); err != nil {
		log.WithError(err).Error("error while restarting exited bots")
	}
	// doing the exits after restarts so the exits we do here can be picked up
	// for restarts in the next round
	if err := sup.botLifecycle.BotManager.ExitInactiveBots(sup.ctx); err != nil {
		log.WithError(err).Error("error while exiting inactive bots")
	}
}
