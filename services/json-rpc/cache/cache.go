package json_rpc_cache

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

type inMemory struct {
	cache *cache.Cache
}

func NewCache(expire time.Duration) *inMemory {
	return &inMemory{
		cache: cache.New(expire, expire),
	}
}

func (c *inMemory) Append(blocksData *protocol.BlocksData) {
	for _, event := range blocksData.Blocks {
		chainID := event.ChainID

		c.cache.SetDefault(cacheKey(chainID, "timestamp", ""), time.Now())

		// eth_blockNumber
		method := "eth_blockNumber"
		params := "[]"

		if val, ok := c.cache.Get(cacheKey(chainID, method, params)); ok {
			blockNumber, ok := val.(string)

			// if the new block number is later than the cached one, update the inMemory
			if ok && isLater(blockNumber, event.Block.Number) {
				c.cache.SetDefault(cacheKey(chainID, method, params), event.Block.Number)
			}
		} else {
			c.cache.SetDefault(cacheKey(chainID, method, params), event.Block.Number)
		}

		log.Debugf("caching block number. chainID: %d method: %s params: %s", chainID, method, params)

		// eth_getBlockByNumber
		method = "eth_getBlockByNumber"
		params = fmt.Sprintf(`["%s",true]`, event.Block.Number)
		log.Debugf("caching block. chainID: %d method: %s params: %s", chainID, method, params)

		block := domain.BlockFromBlockData(event)
		c.cache.SetDefault(cacheKey(chainID, method, params), block)

		// eth_getLogs
		method = "eth_getLogs"
		params = fmt.Sprintf(`[{"fromBlock":"%s","toBlock":"%s"}]`, event.Block.Number, event.Block.Number)

		log.Debugf("caching logs. chainID: %d method: %s params: %s", chainID, method, params)

		logs := domain.LogsFromBlockData(event)
		c.cache.SetDefault(cacheKey(chainID, method, params), logs)

		// trace_block
		method = "trace_block"
		params = fmt.Sprintf(`["%s"]`, event.Block.Number)

		log.Debugf("caching traces. chainID: %d method: %s params: %s", chainID, method, params)

		traces, err := domain.TracesFromBlockData(event)
		if err == nil {
			c.cache.SetDefault(cacheKey(chainID, method, params), traces)
		}
	}
}

func (c *inMemory) Get(chainId uint64, method string, params string) (interface{}, bool) {
	return c.cache.Get(cacheKey(chainId, method, params))
}

func cacheKey(chainId uint64, method, params string) string {
	return fmt.Sprintf("%d-%s-%s", chainId, method, params)
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
