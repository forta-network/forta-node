package scanner

import (
	"sort"

	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/forta-network/forta-node/config"
)

// TxResult contains request and response data.
type TxResult struct {
	AgentConfig config.AgentConfig
	Request     *protocol.EvaluateTxRequest
	Response    *protocol.EvaluateTxResponse
	Timestamps  *domain.TrackingTimestamps
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

// AgentPool contains all the agents which we can forward the alert, block and tx requests
// to and receive the results from.
type AgentPool interface {
	SendEvaluateTxRequest(req *protocol.EvaluateTxRequest)
	TxResults() <-chan *TxResult
	SendEvaluateBlockRequest(req *protocol.EvaluateBlockRequest)
	BlockResults() <-chan *BlockResult
	SendEvaluateAlertRequest(req *protocol.EvaluateAlertRequest)
	CombinationAlertResults() <-chan *CombinationAlertResult
}

func truncateFinding(finding *protocol.Finding) (truncated bool) {
	sort.Strings(finding.Addresses)

	// truncate finding addresses
	lenFindingAddrs := len(finding.Addresses)

	if lenFindingAddrs > utils.NumMaxAddressesPerAlert {
		finding.Addresses = finding.Addresses[:utils.NumMaxAddressesPerAlert]
		truncated = true
	}

	return truncated
}

func reduceMapToArr(m map[string]bool) (result []string) {
	for s := range m {
		result = append(result, s)
	}

	return
}
