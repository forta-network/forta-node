package scanner

import (
	"OpenZeppelin/fotify-node/config"
	"OpenZeppelin/fotify-node/protocol"
)

type AnalyzerAgent struct {
	config config.AgentConfig
	client protocol.AgentClient
}
