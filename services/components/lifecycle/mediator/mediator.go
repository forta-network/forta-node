package mediator

import (
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/services/components/lifecycle"
)

type lifecycleMediator struct {
	msgClient clients.MessageClient
}

// Mediator helps in connecting the bot manager with bot pool.
type Mediator interface {
	ConnectBotPool(botPool lifecycle.BotPoolUpdater)
	lifecycle.BotPoolUpdater
}

// New creates a new bot lifecycle mediator for given bot client pool.
// It lets the bot client pool subscribe to the messaging coming from the bot manager so
// the bot manager and the bot client pool are connected.
// This helps in defining the manager-pool communication concretely.
func New(msgClient clients.MessageClient) Mediator {
	return &lifecycleMediator{
		msgClient: msgClient,
	}
}

// ConnectBotPool connects given bot pool by subscribing to lifecycle management messages.
func (lm *lifecycleMediator) ConnectBotPool(botPool lifecycle.BotPoolUpdater) {
	lm.msgClient.Subscribe(
		messaging.SubjectAgentsStatusRunning, messaging.AgentsHandler(botPool.UpdateBotsWithLatestConfigs),
	)
	lm.msgClient.Subscribe(
		messaging.SubjectAgentsStatusStopping, messaging.AgentsHandler(botPool.RemoveBotsWithConfigs),
	)
	lm.msgClient.Subscribe(
		messaging.SubjectAgentsStatusRestarted, messaging.AgentsHandler(botPool.ReinitBotsWithConfigs),
	)
}

// implement the BotPoolUpdater interface by publishing the lifecycle management messages

func (lm *lifecycleMediator) UpdateBotsWithLatestConfigs(payload messaging.AgentPayload) error {
	lm.msgClient.Publish(messaging.SubjectAgentsStatusRunning, payload)
	return nil
}

func (lm *lifecycleMediator) RemoveBotsWithConfigs(payload messaging.AgentPayload) error {
	lm.msgClient.Publish(messaging.SubjectAgentsStatusStopping, payload)
	return nil
}

func (lm *lifecycleMediator) ReinitBotsWithConfigs(payload messaging.AgentPayload) error {
	lm.msgClient.Publish(messaging.SubjectAgentsStatusRestarted, payload)
	return nil
}
