package messaging

import (
	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/config"
)

// Message types
const (
	SubjectAgentsActionRun        = "agents.action.run"
	SubjectAgentsActionStop       = "agents.action.stop"
	SubjectAgentsAlertSubscribe   = "agents.alert.subscribe"
	SubjectAgentsAlertUnsubscribe = "agents.alert.unsubscribe"
	SubjectAgentsStatusRunning    = "agents.status.running"
	SubjectAgentsStatusAttached   = "agents.status.attached"
	SubjectAgentsStatusStopping   = "agents.status.stopping"
	SubjectAgentsStatusStopped    = "agents.status.stopped"
	SubjectAgentsStatusRestarted  = "agents.status.restarted"
	SubjectMetricAgent            = "metric.agent"
	SubjectScannerBlock           = "scanner.block"
	SubjectScannerAlert           = "scanner.alert"
	SubjectInspectionDone         = "inspection.done"
)

// AgentPayload is the message payload.
type AgentPayload []config.AgentConfig

// AgentMetricPayload is the message payload for metrics.
type AgentMetricPayload *protocol.AgentMetricList

// SubscriptionPayload is the message payload for combiner bot subscriptions.
type SubscriptionPayload []*domain.CombinerBotSubscription

// ScannerPayload is the message payload for general scanner info.
type ScannerPayload struct {
	LatestBlockInput uint64 `json:"latestBlockInput"`
}
