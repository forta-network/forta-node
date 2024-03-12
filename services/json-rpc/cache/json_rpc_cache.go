package json_rpc_cache

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/forta-network/forta-core-go/utils"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/blocksdata"
	"github.com/forta-network/forta-node/config"
	log "github.com/sirupsen/logrus"
)

const (
	ChainIDHeader = "X-Forta-Chain-ID"

	// These values will be injected into the agent container to configure bot cache
	// BotCacheRequestTimeout timeout until the bot must fallback to the RPC Node
	// Value in seconds and can be a float.
	BotCacheRequestTimeoutSeconds = "20"
	// BotCacheRequestInterval interval between bot requests
	// Value in seconds and can be a float.
	BotCacheRequestIntervalSeconds = "1"
	// BotCacheSupportedChains comma separated list of supported chains
	// Chains' data not filtered on the cache side.
	BotCacheSupportedChains = "1,137,56,43114,42161,10,250,8453"
)

type JsonRpcCache struct {
	ctx              context.Context
	cfg              config.JsonRpcCacheConfig
	botAuthenticator clients.IPAuthenticator

	server *http.Server

	cache *inMemory

	blocksDataClient clients.BlocksDataClient
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
	c.cache = NewCache(time.Duration(c.cfg.CacheExpirePeriodSeconds) * time.Second)

	c.server = &http.Server{
		Addr:    ":8575",
		Handler: c.Handler(),
	}

	c.blocksDataClient = blocksdata.NewCombinedBlockEventsClient(c.cfg.DispatcherURL)

	utils.GoListenAndServe(c.server)

	go c.pollEvents()

	return nil
}

func (p *JsonRpcCache) Stop() error {
	if p.server != nil {
		return p.server.Close()
	}
	return nil
}

func (p *JsonRpcCache) Name() string {
	return "json-rpc-cache"
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

		chainID, err := strconv.ParseInt(r.Header.Get(ChainIDHeader), 10, 64)
		if err != nil {
			writeBadRequest(w, req, fmt.Errorf("missing or invalid chain id header"))
			return
		}

		result, ok := c.cache.Get(uint64(chainID), req.Method, string(req.Params))
		if !ok {
			resp := &jsonRpcResp{
				ID:     req.ID,
				Result: nil,
				Error: &jsonRpcError{
					Code:    -32603,
					Message: "result not found in cache",
				},
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

func (c *JsonRpcCache) pollEvents() {
	bucket := time.Now().Truncate(time.Second * 10).Unix()

	for {
		// wait for the next bucket
		<-time.After(time.Duration(bucket-time.Now().Unix()) * time.Second)

		log.Info("Polling for combined block events", "bucket", bucket)

		// blocksDataClient internally retries on failure and to not block on the retry, we run it in a goroutine
		go func() {
			events, err := c.blocksDataClient.GetBlocksData(bucket)
			if err != nil {
				log.WithError(err).Error("Failed to get combined block events", "bucket", bucket)
				return
			}

			log.Info("Added combined block events to local cache", "bucket", bucket, "events", len(events.Blocks))
			c.cache.Append(events)
		}()

		bucket += 10 // 10 seconds
	}
}
