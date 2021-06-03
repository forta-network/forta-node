package scanner

import (
	"github.com/OpenZeppelin/fortify-node/config"
	"github.com/OpenZeppelin/fortify-node/protocol"
)

type AnalyzerAgent struct {
	config config.AgentConfig
	client protocol.AgentClient
}
