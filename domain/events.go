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
	Logs      []LogEntry
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
	Logs        []*LogEntry
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
	var logs []*protocol.TransactionEvent_Log

	for _, l := range t.Logs {
		var evtLog protocol.TransactionEvent_Log
		logJson, err := json.Marshal(l)
		if err != nil {
			return nil, err
		}
		err = um.Unmarshal(bytes.NewReader(logJson), &evtLog)

		if err != nil {
			log.Errorf("cannot unmarshal logJson: %s", err.Error())
			log.Errorf("JSON: %s", string(logJson))
			return nil, err
		}
		safeAddStrToMap(addresses, l.Address)
		logs = append(logs, &evtLog)
	}

	nw := &protocol.TransactionEvent_Network{}
	if t.BlockEvt.ChainID != nil {
		nw.ChainId = utils.BigIntToHex(t.BlockEvt.ChainID)
	}

	fakeReceipt := &protocol.TransactionEvent_EthReceipt{
		Root:              "",
		Status:            "0x1",
		CumulativeGasUsed: "0x0",
		LogsBloom:         "0x0",
		Logs:              logs,
		TransactionHash:   t.Transaction.Hash,
		ContractAddress:   "",
		GasUsed:           "0x0",
		BlockHash:         t.BlockEvt.Block.Hash,
		BlockNumber:       t.BlockEvt.Block.Number,
		TransactionIndex:  t.Transaction.TransactionIndex,
	}

	return &protocol.TransactionEvent{
		Type:        evtType,
		Transaction: &tx,
		Logs:        logs,
		Network:     nw,
		Traces:      traces,
		Addresses:   addresses,
		Receipt:     fakeReceipt,
		Block: &protocol.TransactionEvent_EthBlock{
			BlockHash:      t.BlockEvt.Block.Hash,
			BlockNumber:    t.BlockEvt.Block.Number,
			BlockTimestamp: t.BlockEvt.Block.Timestamp,
		},
	}, nil
}
