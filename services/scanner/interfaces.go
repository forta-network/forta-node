package scanner

import (
	"bytes"
	"encoding/base64"
	"math/big"
	"sort"

	"github.com/bits-and-blooms/bloom"
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

const (
	maxAddressesLength       = 50
	addressBloomFilterFPRate = 1e-3
)

func truncateFinding(finding *protocol.Finding) (truncated bool) {
	sort.Strings(finding.Addresses)

	// truncate finding addresses
	lenFindingAddrs := len(finding.Addresses)
	if lenFindingAddrs > maxAddressesLength {
		finding.Addresses = finding.Addresses[:maxAddressesLength]
		truncated = true
	}

	return truncated
}

func createBloomFilter(allAddresses []string) (*protocol.BloomFilter, error) {

	// create bloom filter from all addresses
	bf := bloom.NewWithEstimates(uint(len(allAddresses)), addressBloomFilterFPRate)
	for _, address := range allAddresses {
		bf.Add([]byte(address))
	}

	// extract bitset from bloom filter
	var b bytes.Buffer

	_, err := bf.WriteTo(&b)
	if err != nil {
		return nil, err
	}

	// create bloom filter
	bitset := base64.StdEncoding.EncodeToString(b.Bytes())

	kBigInt := new(big.Int).SetUint64(uint64(bf.K()))
	mBigInt := new(big.Int).SetUint64(uint64(bf.Cap()))

	kHexStr := utils.BigIntToHex(kBigInt)
	mHexStr := utils.BigIntToHex(mBigInt)

	return &protocol.BloomFilter{
		K:         kHexStr,
		M:         mHexStr,
		Bitset:    bitset,
		ItemCount: uint32(len(allAddresses)),
	}, nil
}

func reduceMapToArr(m map[string]bool) (result []string) {
	for s := range m {
		result = append(result, s)
	}

	return
}
