package json_rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/config"
	"github.com/gogo/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

type JsonRpcCache struct {
	ctx              context.Context
	cfg              config.JsonRpcCacheConfig
	botAuthenticator clients.IPAuthenticator

	server *http.Server

	cache *cache
}

func NewJsonRpcCache(ctx context.Context, cfg config.JsonRpcCacheConfig) (*JsonRpcCache, error) {
	botAuthenticator, err := clients.NewBotAuthenticator(ctx)
	if err != nil {
		return nil, err
	}

	return &JsonRpcCache{
		ctx:              ctx,
		cfg:              cfg,
		botAuthenticator: botAuthenticator,
	}, nil
}

func (c *JsonRpcCache) Start() error {
	c.cache = &cache{
		chains:      make(map[uint64]*chainCache),
		cacheExpire: time.Duration(c.cfg.CacheExpirePeriodSeconds) * time.Second,
	}

	c.server = &http.Server{
		Addr:    ":8575",
		Handler: c.Handler(),
	}

	utils.GoListenAndServe(c.server)

	go c.pollEvents()

	return nil
}

func (c *JsonRpcCache) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := decodeBody(r)
		if err != nil {
			writeBadRequest(w, req, err)
			return
		}

		agentConfig, err := c.botAuthenticator.FindAgentFromRemoteAddr(r.RemoteAddr)
		if agentConfig == nil || err != nil {
			writeUnauthorized(w, nil)
			return
		}

		chainID, err := strconv.ParseInt(r.Header.Get("X-Forta-Chain-ID"), 10, 64)
		if err != nil {
			writeBadRequest(w, req, fmt.Errorf("missing or invalid chain id header"))
			return
		}

		result, ok := c.cache.Get(uint64(chainID), req.Method, req.Params)
		if !ok {
			resp := &jsonRpcResp{
				ID:     req.ID,
				Result: nil,
			}

			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.WithError(err).Error("failed to write jsonrpc response body")
			}
			return
		}

		b, err := json.Marshal(result)
		if err != nil {
			writeBadRequest(w, req, err)
			return
		}

		resp := &jsonRpcResp{
			ID:     req.ID,
			Result: json.RawMessage(b),
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.WithError(err).Error("failed to write jsonrpc response body")
		}
	})
}

type PresignedURLItem struct {
	Bucket       int64  `dynamodbav:"bucket"`
	PresignedURL string `dynamodbav:"presigned_url"`
	ExpiresAt    int64  `dynamodbav:"expires_at"`
}

func (c *JsonRpcCache) pollEvents() {
	dispatcherClient := http.DefaultClient
	dispatcherClient.Timeout = 10 * time.Second
	dispatcherURL, _ := url.Parse(c.cfg.DispatcherURL)

	r2Client := http.DefaultClient

	var err error
	for {
		time.Sleep(1 * time.Second)
		log.Info("Polling for combined block events")

		bucket := time.Now().Truncate(time.Second * 10).Unix()
		dispatcherURL.Path, err = url.JoinPath(dispatcherURL.Path, fmt.Sprintf("%d", bucket))
		if err != nil {
			continue
		}

		resp, err := dispatcherClient.Get(dispatcherURL.String())
		if err != nil {
			log.WithError(err).Error("Failed to get R2 url from dispatcher")
			continue
		}

		var item PresignedURLItem
		err = json.NewDecoder(resp.Body).Decode(&item)
		if err != nil {
			log.WithError(err).Error("Failed to decode presigned url")
			continue
		}

		err = resp.Body.Close()
		if err != nil {
			continue
		}

		resp, err = r2Client.Get(item.PresignedURL)
		if err != nil {
			log.WithError(err).Error("Failed to get combined block events from R2")
			continue
		}

		b, err := io.ReadAll(brotli.NewReader(resp.Body))
		if err != nil {
			log.WithError(err).Error("Failed to uncompress combined block events")
			continue
		}

		var events protocol.CombinedBlockEvents

		err = proto.Unmarshal(b, &events)
		if err != nil {
			log.WithError(err).Error("Failed to unmarshal combined block events")
			continue
		}

		log.Info("Added combined block events to local cache")
		c.cache.Append(&events)
	}
}

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
		val, ok := cc.get("eth_blockNumber", "")
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
				cc.put("eth_blockNumber", "", event.Block.Number)
			}
		} else {
			cc.put("eth_blockNumber", "", event.Block.Number)
		}

		keys = append(keys, "eth_blockNumber")

		// eth_getBlockByNumber
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

		cc.put("eth_getBlockByNumber", fmt.Sprintf(`["%s", "true"]`, event.Block.Number), blockByNumber)
		keys = append(keys, "eth_getBlockByNumber"+fmt.Sprintf(`["%s", "true"]`, event.Block.Number))

		// eth_getLogs
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
				logsByBlock[len(logsByBlock)-1].Topics[i] = &topic
			}
		}

		cc.put("eth_getLogs", fmt.Sprintf(`[{"fromBlock":"%s","toBlock":"%s"}]`, event.Block.Number, event.Block.Number), logsByBlock)
		keys = append(keys, "eth_getLogs"+fmt.Sprintf(`[{"fromBlock":"%s","toBlock":"%s"}]`, event.Block.Number, event.Block.Number))

		// trace_block
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

		cc.put("trace_block", fmt.Sprintf(`["%s"]`, event.Block.Number), traceBlock)
		keys = append(keys, "trace_block"+fmt.Sprintf(`["%s"]`, event.Block.Number))

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
	c.data[method+params] = result
}

func (c *chainCache) get(method string, params string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result, ok := c.data[method+params]

	return result, ok
}

func (c *chainCache) collectGarbage(keys []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, key := range keys {
		delete(c.data, key)
	}
}
