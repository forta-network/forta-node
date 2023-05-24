package agentgrpc

import (
	"github.com/forta-network/forta-node/config"
)

// BotDialer dials a bot.
type BotDialer interface {
	DialBot(ac config.AgentConfig) (Client, error)
}

type botDialer struct{}

// NewBotDialer creates a new bot dialer.
func NewBotDialer() BotDialer {
	return &botDialer{}
}

func (bd *botDialer) DialBot(ac config.AgentConfig) (Client, error) {
	client := NewClient()
	err := client.DialWithRetry(ac)
	if err != nil {
		return nil, err
	}
	return client, nil
}
