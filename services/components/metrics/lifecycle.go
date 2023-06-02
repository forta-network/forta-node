package metrics

import (
	"fmt"
	"time"

	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/config"
)

// Bot lifecycle metrics
const (
	MetricClientDial  = "agent.client.dial"
	MetricClientClose = "agent.client.close"

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
	MetricFailureStop               = "agent.failure.stop"
	MetricFailureDial               = "agent.failure.dial"
	MetricFailureInitialize         = "agent.failure.initialize"
	MetricFailureInitializeResponse = "agent.failure.initialize.response"
	MetricFailureInitializeValidate = "agent.failure.initialize.validate"
	MetricFailureTooManyErrs        = "agent.failure.too-many-errs"
)

// Lifecycle creates lifecycle metrics. It is useful in
// understanding what is going on during lifecycle management.
type Lifecycle interface {
	ClientDial(...config.AgentConfig)
	ClientClose(...config.AgentConfig)

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

	FailurePull(error, ...config.AgentConfig)
	FailureLaunch(error, ...config.AgentConfig)
	FailureStop(error, ...config.AgentConfig)
	FailureDial(error, ...config.AgentConfig)
	FailureInitialize(error, ...config.AgentConfig)
	FailureInitializeResponse(error, ...config.AgentConfig)
	FailureInitializeValidate(error, ...config.AgentConfig)
	FailureTooManyErrs(error, ...config.AgentConfig)

	BotError(metricName string, err error, botID ...string)
	SystemError(metricName string, err error)
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

func (lc *lifecycle) ClientDial(botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricClientDial, "", botConfigs))
}

func (lc *lifecycle) ClientClose(botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricClientClose, "", botConfigs))
}

func (lc *lifecycle) StatusRunning(botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricStatusRunning, "", botConfigs))
}

func (lc *lifecycle) StatusAttached(botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricStatusAttached, "", botConfigs))
}

func (lc *lifecycle) StatusInitialized(botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricStatusInitialized, "", botConfigs))
}

func (lc *lifecycle) StatusStopping(botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricStatusStopping, "", botConfigs))
}

func (lc *lifecycle) StatusActive(botIDs []string) {
	SendAgentMetrics(lc.msgClient, fromBotIDs(MetricStatusActive, "", botIDs))
}

func (lc *lifecycle) StatusInactive(botIDs []string) {
	SendAgentMetrics(lc.msgClient, fromBotIDs(MetricStatusInactive, "", botIDs))
}

func (lc *lifecycle) ActionUpdate(botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricActionUpdate, "", botConfigs))
}

func (lc *lifecycle) ActionRestart(botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricActionRestart, "", botConfigs))
}

func (lc *lifecycle) ActionSubscribe(subscriptions []domain.CombinerBotSubscription) {
	SendAgentMetrics(lc.msgClient, fromBotSubscriptions(MetricActionSubscribe, subscriptions))
}

func (lc *lifecycle) ActionUnsubscribe(subscriptions []domain.CombinerBotSubscription) {
	SendAgentMetrics(lc.msgClient, fromBotSubscriptions(MetricActionUnsubscribe, subscriptions))
}

func (lc *lifecycle) FailurePull(err error, botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricFailurePull, err.Error(), botConfigs))
}

func (lc *lifecycle) FailureLaunch(err error, botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricFailureLaunch, err.Error(), botConfigs))
}

func (lc *lifecycle) FailureStop(err error, botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricFailureStop, err.Error(), botConfigs))
}

func (lc *lifecycle) FailureDial(err error, botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricFailureDial, err.Error(), botConfigs))
}

func (lc *lifecycle) FailureInitialize(err error, botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricFailureInitialize, err.Error(), botConfigs))
}

func (lc *lifecycle) FailureInitializeResponse(err error, botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricFailureInitializeResponse, err.Error(), botConfigs))
}

func (lc *lifecycle) FailureInitializeValidate(err error, botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricFailureInitializeValidate, err.Error(), botConfigs))
}

func (lc *lifecycle) FailureTooManyErrs(err error, botConfigs ...config.AgentConfig) {
	SendAgentMetrics(lc.msgClient, fromBotConfigs(MetricFailureTooManyErrs, err.Error(), botConfigs))
}

func (lc *lifecycle) BotError(metricName string, err error, botIDs ...string) {
	SendAgentMetrics(lc.msgClient, fromBotIDs(fmt.Sprintf("agent.error.%s", metricName), err.Error(), botIDs))
}

func (lc *lifecycle) SystemError(metricName string, err error) {
	SendAgentMetrics(lc.msgClient, fromBotIDs(fmt.Sprintf("system.error.%s", metricName), err.Error(), []string{"system"}))
}

func fromBotSubscriptions(action string, subscriptions []domain.CombinerBotSubscription) (metrics []*protocol.AgentMetric) {
	for _, botSub := range subscriptions {
		metrics = append(metrics, CreateAgentMetric(botSub.Subscriber.BotID, action, 1))
	}
	return
}

func fromBotConfigs(metricName string, details string, botConfigs []config.AgentConfig) (metrics []*protocol.AgentMetric) {
	for _, botConfig := range botConfigs {
		if details == "" && botConfig.IsSharded() {
			details = fmt.Sprintf("shard=%d", botConfig.ShardConfig.ShardID)
		}
		metrics = append(metrics, &protocol.AgentMetric{
			AgentId:   botConfig.ID,
			Timestamp: time.Now().Format(time.RFC3339),
			Name:      metricName,
			Details:   details,
			Value:     1,
		})
	}
	return
}

func fromBotIDs(metricName string, details string, botIDs []string) (metrics []*protocol.AgentMetric) {
	for _, botID := range botIDs {
		metrics = append(metrics, &protocol.AgentMetric{
			AgentId:   botID,
			Timestamp: time.Now().Format(time.RFC3339),
			Name:      metricName,
			Details:   details,
			Value:     1,
		})
	}
	return
}
