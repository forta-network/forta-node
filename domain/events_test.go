package domain

import (
	"github.com/golang/protobuf/jsonpb"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func intPtr(val int) *int {
	return &val
}

func TestTransactionEvent_ToMessage(t *testing.T) {
	blockHash := "0x8d2636ff603ef946d97ad797ed13afa31234a3412dacdfecfeb3247230eb1069"
	txHash := "0x99ed5a4e541454219b444250c5c25d0306e73834b185f3aeee3f9627f0cd64c2"

	// these are checksum addresses, to confirm that logic lower-cases these
	from := "0xa7d8d9ef8D8Ce8992Df33D8b8CF4Aebabd5bD270"
	to := "0x9C025948e61aeB2EF99503c81d682045f07344c2"

	evt := &TransactionEvent{
		BlockEvt: &BlockEvent{
			EventType: "block",
			ChainID:   big.NewInt(1),
			Block: &Block{
				Difficulty:       strPtr("0x1"),
				ExtraData:        strPtr("0x1"),
				GasLimit:         strPtr("0x1"),
				GasUsed:          strPtr("0x1"),
				Hash:             blockHash,
				LogsBloom:        strPtr("0x1"),
				Miner:            strPtr("0x1"),
				MixHash:          strPtr("0x1"),
				Nonce:            strPtr("0x1"),
				Number:           "0x1",
				ParentHash:       "0xabcdef",
				ReceiptsRoot:     strPtr("0x1"),
				Sha3Uncles:       strPtr("0x1"),
				Size:             strPtr("0x1"),
				StateRoot:        strPtr("0x1"),
				Timestamp:        "0x12345",
				TotalDifficulty:  strPtr("0x1"),
				Transactions:     []Transaction{},
				TransactionsRoot: strPtr("0x1"),
				Uncles:           []*string{strPtr("0x1")},
			},
			Traces: []Trace{
				{
					Action:              TraceAction{To: &to, From: &from},
					BlockHash:           &blockHash,
					BlockNumber:         intPtr(1),
					TransactionHash:     &txHash,
					TransactionPosition: intPtr(5),
					Type:                "transaction",
				},
			},
		},
		Transaction: &Transaction{
			BlockHash:   blockHash,
			BlockNumber: "0x1",
			From:        from,
			Gas:         "0x2",
			GasPrice:    "0x3",
			Hash:        txHash,
			Nonce:       "0x5",
			To:          &to,
		},
		Receipt: &TransactionReceipt{
			BlockHash:       &blockHash,
			BlockNumber:     strPtr("0x1"),
			From:            &from,
			ContractAddress: strPtr(to),
			Logs: []LogEntry{
				{
					Address:         strPtr(to),
					BlockHash:       &blockHash,
					BlockNumber:     strPtr("0x2"),
					TransactionHash: &txHash,
				},
			},
			Status:          strPtr("0x1"),
			To:              &to,
			TransactionHash: &txHash,
		},
	}
	msg, err := evt.ToMessage()
	assert.NoError(t, err, "error returned from ToMessage")

	js := jsonpb.Marshaler{}
	str, err := js.MarshalToString(msg)
	t.Log(str)

	// I manually checked this json, so this test just ensures this behavior continues
	expected := `{"transaction":{"nonce":"0x5","gasPrice":"0x3","gas":"0x2","to":"0x9c025948e61aeb2ef99503c81d682045f07344c2","hash":"0x99ed5a4e541454219b444250c5c25d0306e73834b185f3aeee3f9627f0cd64c2","from":"0xa7d8d9ef8d8ce8992df33d8b8cf4aebabd5bd270"},"receipt":{"status":"0x1","logs":[{"address":"0x9c025948e61aeb2ef99503c81d682045f07344c2","blockNumber":"0x2","transactionHash":"0x99ed5a4e541454219b444250c5c25d0306e73834b185f3aeee3f9627f0cd64c2","blockHash":"0x8d2636ff603ef946d97ad797ed13afa31234a3412dacdfecfeb3247230eb1069"}],"transactionHash":"0x99ed5a4e541454219b444250c5c25d0306e73834b185f3aeee3f9627f0cd64c2","contractAddress":"0x9c025948e61aeb2ef99503c81d682045f07344c2","blockHash":"0x8d2636ff603ef946d97ad797ed13afa31234a3412dacdfecfeb3247230eb1069","blockNumber":"0x1"},"network":{"chainId":"0x1"},"traces":[{"action":{"to":"0x9c025948e61aeb2ef99503c81d682045f07344c2","from":"0xa7d8d9ef8d8ce8992df33d8b8cf4aebabd5bd270"},"blockHash":"0x8d2636ff603ef946d97ad797ed13afa31234a3412dacdfecfeb3247230eb1069","blockNumber":"1","transactionHash":"0x99ed5a4e541454219b444250c5c25d0306e73834b185f3aeee3f9627f0cd64c2","transactionPosition":"5","type":"transaction"}],"addresses":{"0x9c025948e61aeb2ef99503c81d682045f07344c2":true,"0xa7d8d9ef8d8ce8992df33d8b8cf4aebabd5bd270":true},"block":{"blockHash":"0x8d2636ff603ef946d97ad797ed13afa31234a3412dacdfecfeb3247230eb1069","blockNumber":"0x1","blockTimestamp":"0x12345"}}`
	assert.NoError(t, err, "error returned from json conversion")
	assert.Equal(t, str, expected)
}
