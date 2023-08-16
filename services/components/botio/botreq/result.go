package botreq

import (
	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/config"
)

// TxResult contains request and response data.
type TxResult struct {
	AgentConfig config.AgentConfig
	Request     *protocol.EvaluateTxRequest
	Response    *protocol.EvaluateTxResponse
	Timestamps  *domain.TrackingTimestamps
}

// HealthCheckResult contains request and response data.
type HealthCheckResult struct {
	AgentConfig config.AgentConfig
	Request     *protocol.HealthCheckRequest
	Response    *protocol.HealthCheckResponse
	Timestamps  *domain.TrackingTimestamps
	InvokeError error
}

// BlockResult contains request and response data.
type BlockResult struct {
	AgentConfig config.AgentConfig
	Request     *protocol.EvaluateBlockRequest
	Response    *protocol.EvaluateBlockResponse
	Timestamps  *domain.TrackingTimestamps
}

// CombinationAlertResult contains request and response data.
type CombinationAlertResult struct {
	AgentConfig config.AgentConfig
	Request     *protocol.EvaluateAlertRequest
	Response    *protocol.EvaluateAlertResponse
	Timestamps  *domain.TrackingTimestamps
}

// SendReceiveChannels has the bot result channels.
type SendReceiveChannels struct {
	Tx               chan *TxResult
	Block            chan *BlockResult
	CombinationAlert chan *CombinationAlertResult
	HealthCheck      chan *HealthCheckResult
}

// MakeResultChannels makes the result channels and returns.
func MakeResultChannels() SendReceiveChannels {
	return SendReceiveChannels{
		Tx:               make(chan *TxResult),
		Block:            make(chan *BlockResult),
		CombinationAlert: make(chan *CombinationAlertResult),
		HealthCheck:      make(chan *HealthCheckResult),
	}
}

// ReceiveOnly returns the receive-only channels so that we cannot send.
func (src SendReceiveChannels) ReceiveOnly() ReceiveOnlyChannels {
	return ReceiveOnlyChannels{
		Tx:               src.Tx,
		Block:            src.Block,
		CombinationAlert: src.CombinationAlert,
		HealthCheck:      src.HealthCheck,
	}
}

// SendOnly returns the send-only channels so that we cannot receive.
func (src SendReceiveChannels) SendOnly() SendOnlyChannels {
	return SendOnlyChannels{
		Tx:               src.Tx,
		Block:            src.Block,
		CombinationAlert: src.CombinationAlert,
		HealthCheck:      src.HealthCheck,
	}
}

// ReceiveOnlyChannels has the bot result channels.
type ReceiveOnlyChannels struct {
	Tx               <-chan *TxResult
	Block            <-chan *BlockResult
	CombinationAlert <-chan *CombinationAlertResult
	HealthCheck      <-chan *HealthCheckResult
}

// SendOnlyChannels has the bot result channels.
type SendOnlyChannels struct {
	Tx               chan<- *TxResult
	Block            chan<- *BlockResult
	CombinationAlert chan<- *CombinationAlertResult
	HealthCheck      chan<- *HealthCheckResult
}
