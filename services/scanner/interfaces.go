package scanner

import (
	"forta-network/forta-node/config"
	"forta-network/forta-node/protocol"
)

// TxResult contains request and response data.
type TxResult struct {
	AgentConfig config.AgentConfig
	Request     *protocol.EvaluateTxRequest
	Response    *protocol.EvaluateTxResponse
}

// BlockResult contains request and response data.
type BlockResult struct {
	AgentConfig config.AgentConfig
	Request     *protocol.EvaluateBlockRequest
	Response    *protocol.EvaluateBlockResponse
}

// AgentPool contains all of the agents which we can forward the block and tx requests
// to and receive the results from.
type AgentPool interface {
	SendEvaluateTxRequest(req *protocol.EvaluateTxRequest)
	TxResults() <-chan *TxResult
	SendEvaluateBlockRequest(req *protocol.EvaluateBlockRequest)
	BlockResults() <-chan *BlockResult
}
