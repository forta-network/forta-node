package messaging

import (
	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/protocol"
)

// Message types
const (
	SubjectAgentsVersionsLatest = "agents.versions.latest"
	SubjectAgentsActionRun      = "agents.action.run"
	SubjectAgentsActionStop     = "agents.action.stop"
	SubjectAgentsStatusRunning  = "agents.status.running"
	SubjectAgentsStatusStopped  = "agents.status.stopped"

	SubjectAlertsStatusPending   = "alerts.status.pending"
	SubjectAlertsStatusPublished = "alerts.status.published"
)

// AgentPayload is the message payload.
type AgentPayload []config.AgentConfig

// AlertsPayload is the message payload.
type AlertsPayload []*protocol.SignedAlert
