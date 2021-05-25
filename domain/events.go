package domain

import (
	"bytes"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/protobuf/jsonpb"

	"OpenZeppelin/fortify-node/protocol"
	"OpenZeppelin/fortify-node/utils"
)

type EventType string

const (
	EventTypeReorg EventType = "reorg"
	EventTypeBlock EventType = "block"
)

type BlockEvent struct {
	EventType EventType
	ChainID   *big.Int
	Block     *types.Block
	Traces    []Trace
}

func (t *BlockEvent) ToMessage() (*protocol.BlockEvent, error) {
	evtType := protocol.BlockEvent_BLOCK
	if t.EventType == "reorg" {
		evtType = protocol.BlockEvent_REORG
	}
	return &protocol.BlockEvent{
		Type:        evtType,
		BlockHash:   t.Block.Hash().Hex(),
		BlockNumber: utils.BigIntToHex(t.Block.Number()),
		Network: &protocol.BlockEvent_Network{
			ChainId: utils.BigIntToHex(t.ChainID),
		},
	}, nil
}

type TransactionEvent struct {
	BlockEvt    *BlockEvent
	Transaction *types.Transaction
	Receipt     *types.Receipt
}

// ToMessage converts the TransactionEvent to the protocol.TransactionEvent message
func (t *TransactionEvent) ToMessage() (*protocol.TransactionEvent, error) {
	evtType := protocol.TransactionEvent_BLOCK
	if t.BlockEvt.EventType == "reorg" {
		evtType = protocol.TransactionEvent_REORG
	}

	um := jsonpb.Unmarshaler{
		AllowUnknownFields: true,
	}

	// convert trace domain model to proto (filter traces)
	var traces []*protocol.TransactionEvent_Trace
	for _, trace := range t.BlockEvt.Traces {
		if trace.TransactionHash != nil && *trace.TransactionHash == t.Transaction.Hash().Hex() {
			var pTrace protocol.TransactionEvent_Trace
			traceJson, err := json.Marshal(trace)
			if err != nil {
				return nil, err
			}
			if err := um.Unmarshal(bytes.NewReader(traceJson), &pTrace); err != nil {
				return nil, err
			}
			traces = append(traces, &pTrace)
		}
	}

	// convert tx domain model to proto
	var tx protocol.TransactionEvent_EthTransaction
	if t.Transaction != nil {
		txJson, err := t.Transaction.MarshalJSON()
		if err != nil {
			return nil, err
		}
		if err := um.Unmarshal(bytes.NewReader(txJson), &tx); err != nil {
			return nil, err
		}
	}

	// convert receipt domain model to proto
	var receipt protocol.TransactionEvent_EthReceipt
	if t.Receipt != nil {
		receiptJson, err := t.Receipt.MarshalJSON()
		if err != nil {
			return nil, err
		}
		if err := um.Unmarshal(bytes.NewReader(receiptJson), &receipt); err != nil {
			return nil, err
		}
	}

	nw := &protocol.TransactionEvent_Network{
		ChainId: utils.BigIntToHex(t.BlockEvt.ChainID),
	}

	return &protocol.TransactionEvent{
		Type:        evtType,
		Transaction: &tx,
		Receipt:     &receipt,
		Network:     nw,
		Traces:      traces,
	}, nil
}
