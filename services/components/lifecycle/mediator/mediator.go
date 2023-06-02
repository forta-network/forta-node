package mediator

import (
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/services/components/lifecycle"
	"github.com/forta-network/forta-node/services/components/metrics"
)

type lifecycleMediator struct {
	msgClient        clients.MessageClient
	lifecycleMetrics metrics.Lifecycle
}

// Mediator helps in connecting the bot manager with bot pool.
type Mediator interface {
	ConnectBotPool(botPool lifecycle.BotPoolUpdater)
	ConnectBotMonitor(botMonitor lifecycle.BotMonitorUpdater)
	lifecycle.BotPoolUpdater
}

// New creates a new bot lifecycle mediator for given bot client pool.
// It lets the bot client pool subscribe to the messaging coming from the bot manager so
// the bot manager and the bot client pool are connected.
// This helps in defining the manager-pool communication concretely.
func New(msgClient clients.MessageClient, lifecycleMetrics metrics.Lifecycle) Mediator {
	return &lifecycleMediator{
		msgClient:        msgClient,
		lifecycleMetrics: lifecycleMetrics,
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
		messaging.SubjectAgentsStatusRestarted, messaging.AgentsHandler(botPool.ReconnectToBotsWithConfigs),
	)
}

// ConnectBotMonitor connects given bot monitor by subscribing to lifecycle management messages.
func (lm *lifecycleMediator) ConnectBotMonitor(botMonitor lifecycle.BotMonitorUpdater) {
	lm.msgClient.Subscribe(
		messaging.SubjectMetricAgent, messaging.AgentMetricHandler(botMonitor.UpdateWithMetrics),
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

func (lm *lifecycleMediator) ReconnectToBotsWithConfigs(payload messaging.AgentPayload) error {
	lm.msgClient.Publish(messaging.SubjectAgentsStatusRestarted, payload)
	return nil
}
