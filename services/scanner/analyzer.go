package scanner

import (
	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/protocol"
)

type AnalyzerAgent struct {
	config config.AgentConfig
	client protocol.AgentClient
}
