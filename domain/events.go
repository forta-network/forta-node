package domain

import (
	"bytes"
	"encoding/json"
	"math/big"

	"github.com/golang/protobuf/jsonpb"
	log "github.com/sirupsen/logrus"

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
	Block     *Block
	Traces    []Trace
}

func (t *BlockEvent) ToMessage() (*protocol.BlockEvent, error) {
	evtType := protocol.BlockEvent_BLOCK
	if t.EventType == "reorg" {
		evtType = protocol.BlockEvent_REORG
	}
	return &protocol.BlockEvent{
		Type:        evtType,
		BlockHash:   t.Block.Hash,
		BlockNumber: t.Block.Number,
		Network: &protocol.BlockEvent_Network{
			ChainId: utils.BigIntToHex(t.ChainID),
		},
	}, nil
}

type TransactionEvent struct {
	BlockEvt    *BlockEvent
	Transaction *Transaction
	Receipt     *TransactionReceipt
}

func safeAddStrValueToMap(addresses map[string]bool, addr string) {
	if addr != "" {
		addresses[addr] = true
	}
}

func safeAddStrToMap(addresses map[string]bool, addr *string) {
	if addr != nil {
		addresses[*addr] = true
	}
}

// ToMessage converts the TransactionEvent to the protocol.TransactionEvent message
func (t *TransactionEvent) ToMessage() (*protocol.TransactionEvent, error) {
	evtType := protocol.TransactionEvent_BLOCK
	if t.BlockEvt.EventType == "reorg" {
		evtType = protocol.TransactionEvent_REORG
	}

	addresses := make(map[string]bool)

	um := jsonpb.Unmarshaler{
		AllowUnknownFields: true,
	}

	// convert trace domain model to proto (filter traces)
	var traces []*protocol.TransactionEvent_Trace
	for _, trace := range t.BlockEvt.Traces {
		if trace.TransactionHash != nil && *trace.TransactionHash == t.Transaction.Hash {
			safeAddStrToMap(addresses, trace.Action.Address)
			safeAddStrToMap(addresses, trace.Action.RefundAddress)
			safeAddStrToMap(addresses, trace.Action.To)
			safeAddStrToMap(addresses, trace.Action.From)

			var pTrace protocol.TransactionEvent_Trace
			traceJson, err := json.Marshal(trace)
			if err != nil {
				return nil, err
			}
			if err := um.Unmarshal(bytes.NewReader(traceJson), &pTrace); err != nil {
				log.Errorf("cannot unmarshal traceJson: %s", err.Error())
				log.Errorf("JSON: %s", string(traceJson))
				return nil, err
			}
			traces = append(traces, &pTrace)
		}
	}

	// convert tx domain model to proto
	var tx protocol.TransactionEvent_EthTransaction
	if t.Transaction != nil {
		safeAddStrToMap(addresses, t.Transaction.To)
		safeAddStrToMap(addresses, &t.Transaction.From)

		txJson, err := json.Marshal(t.Transaction)
		if err != nil {
			return nil, err
		}
		if err := um.Unmarshal(bytes.NewReader(txJson), &tx); err != nil {
			log.Errorf("cannot unmarshal txJson: %s", err.Error())
			log.Errorf("JSON: %s", string(txJson))
			return nil, err
		}
	}

	// convert receipt domain model to proto
	var receipt protocol.TransactionEvent_EthReceipt
	if t.Receipt != nil {
		receiptJson, err := json.Marshal(t.Receipt)
		if err != nil {
			return nil, err
		}
		err = um.Unmarshal(bytes.NewReader(receiptJson), &receipt)

		if err != nil {
			log.Errorf("cannot unmarshal receiptJson: %s", err.Error())
			log.Errorf("JSON: %s", string(receiptJson))
			return nil, err
		}

		safeAddStrValueToMap(addresses, receipt.ContractAddress)
		for _, l := range receipt.Logs {
			safeAddStrValueToMap(addresses, l.Address)
		}
	}

	nw := &protocol.TransactionEvent_Network{}
	if t.BlockEvt.ChainID != nil {
		nw.ChainId = utils.BigIntToHex(t.BlockEvt.ChainID)
	}

	return &protocol.TransactionEvent{
		Type:        evtType,
		Transaction: &tx,
		Network:     nw,
		Traces:      traces,
		Addresses:   addresses,
		Receipt:     &receipt,
		Block: &protocol.TransactionEvent_EthBlock{
			BlockHash:      t.BlockEvt.Block.Hash,
			BlockNumber:    t.BlockEvt.Block.Number,
			BlockTimestamp: t.BlockEvt.Block.Timestamp,
		},
	}, nil
}
