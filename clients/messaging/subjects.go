package messaging

import (
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/protocol"
)

// Message types
const (
	SubjectAgentsVersionsLatest = "agents.versions.latest"
	SubjectAgentsActionRun      = "agents.action.run"
	SubjectAgentsActionStop     = "agents.action.stop"
	SubjectAgentsStatusRunning  = "agents.status.running"
	SubjectAgentsStatusAttached = "agents.status.attached"
	SubjectAgentsStatusStopped  = "agents.status.stopped"
	SubjectMetricAgent          = "metric.agent"
	SubjectScannerBlock         = "scanner.block"
)

// AgentPayload is the message payload.
type AgentPayload []config.AgentConfig

// AgentMetricPayload is the message payload for metrics.
type AgentMetricPayload *protocol.AgentMetricList

// ScannerPayload is the message payload for general scanner info.
type ScannerPayload struct {
	LatestBlockInput uint64 `json:"latestBlockInput"`
}
