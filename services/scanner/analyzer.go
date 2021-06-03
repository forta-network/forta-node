package scanner

import (
	"fortify-node/config"
	"fortify-node/protocol"
)

type AnalyzerAgent struct {
	config config.AgentConfig
	client protocol.AgentClient
}
