package messaging

import (
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/config"
)

// Message types
const (
	SubjectAgentsVersionsLatest   = "agents.versions.latest"
	SubjectAgentsActionRun        = "agents.action.run"
	SubjectAgentsActionStop       = "agents.action.stop"
	SubjectAgentsAlertSubscribe   = "agents.alert.subscribe"
	SubjectAgentsAlertUnsubscribe = "agents.alert.unsubscribe"
	SubjectAgentsStatusRunning    = "agents.status.running"
	SubjectAgentsStatusAttached   = "agents.status.attached"
	SubjectAgentsStatusStopped    = "agents.status.stopped"
	SubjectMetricAgent            = "metric.agent"
	SubjectScannerBlock           = "scanner.block"
	SubjectScannerAlert           = "scanner.alert"
	SubjectInspectionDone         = "inspection.done"
)

// AgentPayload is the message payload.
type AgentPayload []config.AgentConfig

// AgentMetricPayload is the message payload for metrics.
type AgentMetricPayload *protocol.AgentMetricList

type CombinerBotSubscription struct {
	Subscription *protocol.CombinerBotSubscription
}
type SubscriptionPayload []CombinerBotSubscription

// ScannerPayload is the message payload for general scanner info.
type ScannerPayload struct {
	LatestBlockInput uint64 `json:"latestBlockInput"`
}
