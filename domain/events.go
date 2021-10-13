package domain

import (
	"bytes"
	"encoding/json"
	"math/big"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	log "github.com/sirupsen/logrus"

	"github.com/forta-protocol/forta-node/protocol"
	"github.com/forta-protocol/forta-node/utils"
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

func str(val *string) string {
	if val == nil {
		return ""
	}
	return *val
}

func strArr(vals []*string) []string {
	result := make([]string, 0, len(vals))
	for _, v := range vals {
		result = append(result, str(v))
	}
	return result
}

func strPtr(val string) *string {
	return &val
}

func (t *BlockEvent) ToMessage() (*protocol.BlockEvent, error) {
	evtType := protocol.BlockEvent_BLOCK
	if t.EventType == "reorg" {
		evtType = protocol.BlockEvent_REORG
	}

	txs := make([]string, 0, len(t.Block.Transactions))
	for _, tx := range t.Block.Transactions {
		txs = append(txs, tx.Hash)
	}
	return &protocol.BlockEvent{
		Type:        evtType,
		BlockHash:   t.Block.Hash,
		BlockNumber: t.Block.Number,
		Network: &protocol.BlockEvent_Network{
			ChainId: utils.BigIntToHex(t.ChainID),
		},
		Block: &protocol.BlockEvent_EthBlock{
			Difficulty:       str(t.Block.Difficulty),
			Hash:             t.Block.Hash,
			Number:           t.Block.Number,
			ParentHash:       t.Block.ParentHash,
			Timestamp:        t.Block.Timestamp,
			Nonce:            str(t.Block.Nonce),
			ExtraData:        str(t.Block.ExtraData),
			GasLimit:         str(t.Block.GasLimit),
			GasUsed:          str(t.Block.GasUsed),
			LogsBloom:        str(t.Block.LogsBloom),
			Miner:            str(t.Block.Miner),
			MixHash:          str(t.Block.MixHash),
			Size:             str(t.Block.Size),
			StateRoot:        str(t.Block.StateRoot),
			ReceiptsRoot:     str(t.Block.ReceiptsRoot),
			TotalDifficulty:  str(t.Block.TotalDifficulty),
			Sha3Uncles:       str(t.Block.Sha3Uncles),
			Uncles:           strArr(t.Block.Uncles),
			TransactionsRoot: str(t.Block.TransactionsRoot),
			Transactions:     txs,
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
		addresses[strings.ToLower(addr)] = true
	}
}

func safeAddStrToMap(addresses map[string]bool, addr *string) {
	if addr != nil {
		safeAddStrValueToMap(addresses, *addr)
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
			// lowercase addresses
			if pTrace.Action != nil {
				pTrace.Action.To = strings.ToLower(pTrace.Action.To)
				pTrace.Action.From = strings.ToLower(pTrace.Action.From)
				pTrace.Action.RefundAddress = strings.ToLower(pTrace.Action.RefundAddress)
				pTrace.Action.Address = strings.ToLower(pTrace.Action.Address)
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

		// lowercase to/from
		tx.To = strings.ToLower(tx.To)
		tx.From = strings.ToLower(tx.From)
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

		receipt.ContractAddress = strings.ToLower(receipt.ContractAddress)
		safeAddStrValueToMap(addresses, receipt.ContractAddress)
		for _, l := range receipt.Logs {
			l.Address = strings.ToLower(l.Address)
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
