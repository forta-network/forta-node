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
	"github.com/forta-network/forta-node/clients/r2cbe"
	"github.com/forta-network/forta-node/config"
	jrp "github.com/forta-network/forta-node/services/json-rpc"
	log "github.com/sirupsen/logrus"
)

const (
	ChainIDHeader = "X-Forta-Chain-ID"
)

type JsonRpcCache struct {
	ctx              context.Context
	cfg              config.JsonRpcCacheConfig
	botAuthenticator clients.IPAuthenticator

	server *http.Server

	cache *cache

	cbeClient clients.CombinedBlockEventsClient
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

	c.cbeClient = r2cbe.NewCombinedBlockEventsClient(c.cfg.DispatcherURL)

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
		req, err := jrp.DecodeBody(r)
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
			resp := &jrp.JsonRpcResp{
				ID:     req.ID,
				Result: nil,
				Error: &jrp.JsonRpcError{
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

		resp := &jrp.JsonRpcResp{
			ID:     req.ID,
			Result: json.RawMessage(b),
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.WithError(err).Error("failed to write jsonrpc response body")
		}
	})
}

func (c *JsonRpcCache) pollEvents() {
	for {
		time.Sleep(1 * time.Second)
		log.Info("Polling for combined block events")

		bucket := time.Now().Truncate(time.Second * 10).Unix()

		events, err := c.cbeClient.GetCombinedBlockEvents(bucket)
		if err != nil {
			log.WithError(err).Error("Failed to get combined block events")
			continue
		}

		log.Info("Added combined block events to local cache")
		c.cache.Append(events)
	}
}
