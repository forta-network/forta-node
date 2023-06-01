package metrics

import (
	"time"

	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/config"
)

// Bot lifecycle metrics
const (
	MetricStart = "agent.start"
	MetricStop  = "agent.stop"

	MetricStatusRunning     = "agent.status.running"
	MetricStatusAttached    = "agent.status.attached"
	MetricStatusInitialized = "agent.status.initialized"
	MetricStatusStopping    = "agent.status.stopping"
	MetricStatusActive      = "agent.status.active"
	MetricStatusInactive    = "agent.status.inactive"

	MetricActionUpdate      = "agent.action.update"
	MetricActionRestart     = "agent.action.restart"
	MetricActionSubscribe   = "agent.action.subscribe"
	MetricActionUnsubscribe = "agent.action.unsubscribe"

	MetricFailurePull               = "agent.failure.pull"
	MetricFailureLaunch             = "agent.failure.launch"
	MetricFailureDial               = "agent.failure.dial"
	MetricFailureInitialize         = "agent.failure.initialize"
	MetricFailureInitializeResponse = "agent.failure.initialize.response"
)

// Lifecycle creates lifecycle metrics. It is useful in
// understanding what is going on during lifecycle management.
type Lifecycle interface {
	Start(...config.AgentConfig)
	Stop(...config.AgentConfig)

	StatusRunning(...config.AgentConfig)
	StatusAttached(...config.AgentConfig)
	StatusInitialized(...config.AgentConfig)
	StatusStopping(...config.AgentConfig)
	StatusActive([]string)
	StatusInactive([]string)

	ActionUpdate(...config.AgentConfig)
	ActionRestart(...config.AgentConfig)
	ActionSubscribe([]domain.CombinerBotSubscription)
	ActionUnsubscribe([]domain.CombinerBotSubscription)

	FailurePull(...config.AgentConfig)
	FailureLaunch(...config.AgentConfig)
	FailureDial(...config.AgentConfig)
	FailureInitialize(...config.AgentConfig)
	FailureInitializeResponse(...config.AgentConfig)
}

type lifecycle struct {
	msgClient clients.MessageClient
}

// NewLifecycleClient creates a new metrics client.
func NewLifecycleClient(msgClient clients.MessageClient) Lifecycle {
	return &lifecycle{
		msgClient: msgClient,
	}
}

func (lc *lifecycle) Start(botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricStart, botConfigs))
}

func (lc *lifecycle) Stop(botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricStop, botConfigs))
}

func (lc *lifecycle) StatusRunning(botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricStatusRunning, botConfigs))
}

func (lc *lifecycle) StatusAttached(botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricStatusAttached, botConfigs))
}

func (lc *lifecycle) StatusInitialized(botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricStatusInitialized, botConfigs))
}

func (lc *lifecycle) StatusStopping(botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricStatusStopping, botConfigs))
}

func (lc *lifecycle) StatusActive(botIDs []string) {
	SendAgentMetrics(lc.msgClient, fromBotIDs(MetricStatusActive, botIDs))
}

func (lc *lifecycle) StatusInactive(botIDs []string) {
	SendAgentMetrics(lc.msgClient, fromBotIDs(MetricStatusInactive, botIDs))
}

func (lc *lifecycle) ActionUpdate(botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricActionUpdate, botConfigs))
}

func (lc *lifecycle) ActionRestart(botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricActionRestart, botConfigs))
}

func (lc *lifecycle) ActionSubscribe(subscriptions []domain.CombinerBotSubscription) {
	SendAgentMetrics(lc.msgClient, fromBotSubscriptions(MetricActionSubscribe, subscriptions))
}

func (lc *lifecycle) ActionUnsubscribe(subscriptions []domain.CombinerBotSubscription) {
	SendAgentMetrics(lc.msgClient, fromBotSubscriptions(MetricActionUnsubscribe, subscriptions))
}

func (lc *lifecycle) FailurePull(botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricFailurePull, botConfigs))
}

func (lc *lifecycle) FailureLaunch(botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricFailureLaunch, botConfigs))
}

func (lc *lifecycle) FailureDial(botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricFailureDial, botConfigs))
}

func (lc *lifecycle) FailureInitialize(botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricFailureInitialize, botConfigs))
}

func (lc *lifecycle) FailureInitializeResponse(botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricFailureInitializeResponse, botConfigs))
}

func fromBotSubscriptions(action string, subscriptions []domain.CombinerBotSubscription) (metrics []*protocol.AgentMetric) {
	for _, botSub := range subscriptions {
		metrics = append(metrics, CreateAgentMetric(botSub.Subscriber.BotID, action, 1))
	}
	return
}

func fromBotConfigs(metricName string, botConfigs []config.AgentConfig) (metrics []*protocol.AgentMetric) {
	for _, botConfig := range botConfigs {
		metrics = append(metrics, &protocol.AgentMetric{
			AgentId:   botConfig.ID,
			Timestamp: time.Now().Format(time.RFC3339),
			Name:      metricName,
			Value:     1,
		})
	}
	return
}

func fromBotIDs(metricName string, botIDs []string) (metrics []*protocol.AgentMetric) {
	for _, botID := range botIDs {
		metrics = append(metrics, &protocol.AgentMetric{
			AgentId:   botID,
			Timestamp: time.Now().Format(time.RFC3339),
			Name:      metricName,
			Value:     1,
		})
	}
	return
}
