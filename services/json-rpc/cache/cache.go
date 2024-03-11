package json_rpc_cache

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/protocol"
)

type cache struct {
	chains      map[uint64]*chainCache
	cacheExpire time.Duration
}

func (c *cache) collectGarbage(garbage map[uint64][]string) {
	for chainID, keys := range garbage {
		c.chains[chainID].collectGarbage(keys)
	}
}

func (c *cache) Append(events *protocol.CombinedBlockEvents) {
	garbage := make(map[uint64][]string, 0)
	defer func() {
		go func() {
			time.Sleep(c.cacheExpire)
			c.collectGarbage(garbage)
		}()
	}()

	for _, event := range events.Events {
		chainID := event.ChainID
		cc, ok := c.chains[chainID]
		if !ok {
			cc = &chainCache{
				mu:   &sync.RWMutex{},
				data: make(map[string]interface{}),
			}
			c.chains[chainID] = cc
		}

		keys := make([]string, 0)

		// eth_blockNumber
		method := "eth_blockNumber"
		params := "[]"
		val, ok := cc.get(method, params)
		if ok {
			blockNumber := val.(string)
			actualBlockNumber, err := strconv.ParseInt(strings.Replace(blockNumber, "0x", "", -1), 16, 64)
			if err != nil {
				continue
			}

			newBlockNumber, err := strconv.ParseInt(strings.Replace(event.Block.Number, "0x", "", -1), 16, 64)
			if err != nil {
				continue
			}

			if newBlockNumber > actualBlockNumber {
				cc.put(method, params, event.Block.Number)
			}
		} else {
			cc.put(method, params, event.Block.Number)
		}

		keys = append(keys, cacheKey(method, params))

		// eth_getBlockByNumber
		method = "eth_getBlockByNumber"
		params = fmt.Sprintf(`["%s", "true"]`, event.Block.Number)

		blockByNumber := &domain.Block{
			Difficulty:       &event.Block.Difficulty,
			ExtraData:        &event.Block.ExtraData,
			GasLimit:         &event.Block.GasLimit,
			GasUsed:          &event.Block.GasUsed,
			Hash:             event.Block.Hash,
			LogsBloom:        &event.Block.LogsBloom,
			Miner:            &event.Block.Miner,
			MixHash:          &event.Block.MixHash,
			Nonce:            &event.Block.Nonce,
			Number:           event.Block.Number,
			ParentHash:       event.Block.ParentHash,
			ReceiptsRoot:     &event.Block.ReceiptsRoot,
			Sha3Uncles:       &event.Block.Sha3Uncles,
			Size:             &event.Block.Size,
			StateRoot:        &event.Block.StateRoot,
			Timestamp:        event.Block.Timestamp,
			TotalDifficulty:  &event.Block.TotalDifficulty,
			TransactionsRoot: &event.Block.TransactionsRoot,
		}

		for _, tx := range event.Block.Transactions {
			blockByNumber.Transactions = append(blockByNumber.Transactions, domain.Transaction{
				BlockHash:            event.Block.Hash,
				BlockNumber:          event.Block.Number,
				From:                 tx.From,
				Gas:                  tx.Gas,
				GasPrice:             tx.GasPrice,
				Hash:                 tx.Hash,
				Input:                &tx.Input,
				Nonce:                tx.Nonce,
				To:                   &tx.To,
				TransactionIndex:     tx.TransactionIndex,
				Value:                &tx.Value,
				V:                    tx.V,
				R:                    tx.R,
				S:                    tx.S,
				MaxFeePerGas:         &tx.MaxFeePerGas,
				MaxPriorityFeePerGas: &tx.MaxPriorityFeePerGas,
			})

			for _, uncle := range event.Block.Uncles {
				blockByNumber.Uncles = append(blockByNumber.Uncles, &uncle)
			}
		}

		cc.put(method, params, blockByNumber)
		keys = append(keys, cacheKey(method, params))

		// eth_getLogs
		method = "eth_getLogs"
		params = fmt.Sprintf(`[{"fromBlock":"%s","toBlock":"%s"}]`, event.Block.Number, event.Block.Number)

		logsByBlock := make([]domain.LogEntry, len(event.Logs))
		for i, log := range event.Logs {
			logsByBlock[i] = domain.LogEntry{
				Address:          &log.Address,
				BlockHash:        &log.BlockHash,
				BlockNumber:      &log.BlockNumber,
				Data:             &log.Data,
				LogIndex:         &log.LogIndex,
				Removed:          &log.Removed,
				Topics:           make([]*string, len(log.Topics)),
				TransactionHash:  &log.TransactionHash,
				TransactionIndex: &log.TransactionIndex,
			}

			for i, topic := range log.Topics {
				logsByBlock[i].Topics = append(logsByBlock[i].Topics, &topic)
			}
		}

		cc.put(method, params, logsByBlock)
		keys = append(keys, cacheKey(method, params))

		// trace_block
		method = "trace_block"
		params = fmt.Sprintf(`["%s"]`, event.Block.Number)

		traceBlock := make([]domain.Trace, len(event.Traces))
		blockNumber, err := strconv.ParseInt(strings.Replace(event.Block.Number, "0x", "", -1), 16, 64)
		if err != nil {
			continue
		}

		intBlockNumber := int(blockNumber)
		for i, trace := range event.Traces {
			traceBlock[i] = domain.Trace{
				BlockNumber: &intBlockNumber,
				Subtraces:   int(trace.Subtraces),
				TraceAddress: func() []int {
					traceAddress := make([]int, len(trace.TraceAddress))
					for i, address := range trace.TraceAddress {
						traceAddress[i] = int(address)
					}
					return traceAddress
				}(),
				TransactionHash: &trace.TransactionHash,
				TransactionPosition: func(i int64) *int {
					res := int(i)
					return &res
				}(trace.TransactionPosition),
				Type:  trace.Type,
				Error: &trace.Error,
			}

			if event.Block != nil {
				traceBlock[i].BlockHash = &event.Block.Hash
			}

			if trace.Action != nil {
				traceBlock[i].Action = domain.TraceAction{
					CallType:      &trace.Action.CallType,
					To:            &trace.Action.To,
					Input:         &trace.Action.Input,
					From:          &trace.Action.From,
					Gas:           &trace.Action.Gas,
					Value:         &trace.Action.Value,
					Init:          &trace.Action.Init,
					Address:       &trace.Action.Address,
					Balance:       &trace.Action.Balance,
					RefundAddress: &trace.Action.RefundAddress,
				}
			}

			if trace.Result != nil {
				traceBlock[i].Result = &domain.TraceResult{
					Output:  &trace.Result.Output,
					GasUsed: &trace.Result.GasUsed,
					Address: &trace.Result.Address,
					Code:    &trace.Result.Code,
				}
			}
		}

		cc.put(method, params, traceBlock)
		keys = append(keys, cacheKey(method, params))

		garbage[chainID] = keys
	}
}

func (c *cache) Get(chainId uint64, method string, params string) (interface{}, bool) {
	cc, ok := c.chains[chainId]
	if ok {
		return cc.get(method, params)
	}

	return nil, false
}

type chainCache struct {
	mu *sync.RWMutex
	// key is method + params
	data map[string]interface{}
}

func (c *chainCache) put(method string, params string, result interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[cacheKey(method, params)] = result
}

func (c *chainCache) get(method string, params string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result, ok := c.data[cacheKey(method, params)]

	return result, ok
}

func (c *chainCache) collectGarbage(keys []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, key := range keys {
		delete(c.data, key)
	}
}

func cacheKey(method, params string) string {
	return method + params
}
