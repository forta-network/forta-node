package scanner

import (
	"bytes"
	"sort"

	"github.com/bits-and-blooms/bloom"
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


const (
	maxAddressesLength       = 50
	addressBloomFilterFPRate = 1e-3
)

func truncateFinding(finding *protocol.Finding) (bloomFilter *protocol.BloomFilter, truncated bool) {
	sort.Strings(finding.Addresses)

	// create bloom filter from addresses
	bf := bloom.NewWithEstimates(uint(len(finding.Addresses)), addressBloomFilterFPRate)
	for _, address := range finding.Addresses {
		bf.Add([]byte(address))
	}

	// extract bitset from bloom filter
	var b bytes.Buffer

	_, err := bf.WriteTo(&b)
	if err != nil {
		return nil, false
	}

	if len(finding.Addresses) > maxAddressesLength {
		finding.Addresses = finding.Addresses[:maxAddressesLength]
		truncated = true
	}

	return &protocol.BloomFilter{
		K:      uint64(bf.K()),
		M:      uint64(bf.Cap()),
		Bitset: b.Bytes(),
	}, truncated
}