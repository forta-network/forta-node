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
	SubjectImagesLatest         = "images.latest"
)

// AgentPayload is the message payload.
type AgentPayload []config.AgentConfig

// AgentMetricPayload is the message payload for metrics.
type AgentMetricPayload *protocol.AgentMetricList

// ImagesPayload is message payload for Forta node image references.
type ImagesPayload config.FortaImages
