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

		if val, ok := cc.get(method, params); ok {
			blockNumber, ok := val.(string)

			// if the new block number is later than the cached one, update the cache
			if ok && isLater(blockNumber, event.Block.Number) {
				cc.put(method, params, event.Block.Number)
			}
		} else {
			cc.put(method, params, event.Block.Number)
		}

		keys = append(keys, cacheKey(method, params))

		// eth_getBlockByNumber
		method = "eth_getBlockByNumber"
		params = fmt.Sprintf(`["%s", "true"]`, event.Block.Number)

		block := domain.BlockFromCombinedBlockEvent(event)
		cc.put(method, params, block)
		keys = append(keys, cacheKey(method, params))

		// eth_getLogs
		method = "eth_getLogs"
		params = fmt.Sprintf(`[{"fromBlock":"%s","toBlock":"%s"}]`, event.Block.Number, event.Block.Number)

		logs := domain.LogsFromCombinedBlockEvent(event)
		cc.put(method, params, logs)
		keys = append(keys, cacheKey(method, params))

		// trace_block
		method = "trace_block"
		params = fmt.Sprintf(`["%s"]`, event.Block.Number)

		traces, err := domain.TracesFromCombinedBlockEvent(event)
		if err == nil {
			cc.put(method, params, traces)
			keys = append(keys, cacheKey(method, params))
		}

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

func isLater(actual, new string) bool {
	actualInt, err := strconv.ParseInt(strings.Replace(actual, "0x", "", -1), 16, 64)
	if err != nil {
		return false
	}

	newInt, err := strconv.ParseInt(strings.Replace(new, "0x", "", -1), 16, 64)
	if err != nil {
		return false
	}

	return newInt > actualInt
}
