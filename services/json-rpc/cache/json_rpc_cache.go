package json_rpc_cache

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/blocksdata"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/services/components/metrics"
	"github.com/forta-network/forta-node/services/components/registry"
	"github.com/gorilla/mux"
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
)

type JsonRpcCache struct {
	ctx              context.Context
	cfg              config.JsonRpcCacheConfig
	botAuthenticator clients.IPAuthenticator
	botRegistry      registry.BotRegistry
	msgClient        clients.MessageClient

	server *http.Server

	cache *inMemory

	blocksDataClient clients.BlocksDataClient
}

func NewJsonRpcCache(ctx context.Context, cfg config.JsonRpcCacheConfig, botRegistry registry.BotRegistry) (*JsonRpcCache, error) {
	botAuthenticator, err := clients.NewBotAuthenticator(ctx)
	if err != nil {
		return nil, err
	}

	return &JsonRpcCache{
		ctx:              ctx,
		cfg:              cfg,
		botAuthenticator: botAuthenticator,
		botRegistry:      botRegistry,
		msgClient:        messaging.NewClient("json-rpc-cache", fmt.Sprintf("%s:%s", config.DockerNatsContainerName, config.DefaultNatsPort)),
	}, nil
}

func (c *JsonRpcCache) Start() error {
	c.cache = NewCache(time.Duration(c.cfg.CacheExpirePeriodSeconds) * time.Second)

	r := mux.NewRouter()
	r.Handle("/", c.Handler())
	r.Handle("/health/{chainID}", c.HealthHandler())

	c.server = &http.Server{
		Addr:    ":8575",
		Handler: r,
	}

	c.blocksDataClient = blocksdata.NewBlocksDataClient(c.cfg.DispatcherURL)

	utils.GoListenAndServe(c.server)

	go c.pollBlocksData()

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
		var err error
		t := time.Now()
		req, err := decodeBody(r)
		if err != nil {
			writeBadRequest(w, req, err)
			return
		}

		agentConfig, err := c.botAuthenticator.FindAgentFromRemoteAddr(r.RemoteAddr)
		if agentConfig == nil || err != nil {
			writeUnauthorized(w, req)
			return
		}

		defer func() {
			if err != nil {
				c.msgClient.PublishProto(
					messaging.SubjectMetricAgent, &protocol.AgentMetricList{
						Metrics: []*protocol.AgentMetric{
							metrics.CreateAgentMetricV1(*agentConfig, domain.MetricJSONRPCCachePollError, 1),
						},
					},
				)
			}
		}()

		chainID, err := strconv.ParseInt(r.Header.Get(ChainIDHeader), 10, 64)
		if err != nil {
			writeBadRequest(w, req, fmt.Errorf("missing or invalid chain id header"))
			return
		}

		details := fmt.Sprintf("chainID: %d method: %s params: %s", chainID, req.Method, string(req.Params))
		log.Debugf("jsonrpc cache request. %s", details)

		result, ok := c.cache.Get(uint64(chainID), req.Method, string(req.Params))
		if !ok {
			log.Debugf("cache miss. %s", details)
			c.msgClient.PublishProto(
				messaging.SubjectMetricAgent, &protocol.AgentMetricList{
					Metrics: []*protocol.AgentMetric{
						metrics.CreateDetailedAgentMetricV2(*agentConfig, fmt.Sprintf("%s.%s", domain.MetricJSONRPCCacheMiss, req.Method), 1, details, chainID),
					},
				},
			)
			writeNotFound(w, req)
			return
		}

		err = writeJsonResponse(w, req, result)
		if err != nil {
			log.WithError(err).Error("failed to write jsonrpc response body")
			writeInternalError(w, req, err)
			return
		}

		since := float64(time.Since(t).Milliseconds())
		c.msgClient.PublishProto(
			messaging.SubjectMetricAgent, &protocol.AgentMetricList{
				Metrics: []*protocol.AgentMetric{
					metrics.CreateDetailedAgentMetricV2(*agentConfig, fmt.Sprintf("%s.%s", domain.MetricJSONRPCCacheHit, req.Method), 1, details, chainID),
					metrics.CreateDetailedAgentMetricV2(*agentConfig, domain.MetricJSONRPCCacheLatency, since, details, chainID),
				},
			},
		)
	})
}

func (c *JsonRpcCache) HealthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chainID, err := strconv.ParseInt(mux.Vars(r)["chainID"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("chain id must be an integer"))
			return
		}

		t, ok := c.cache.Get(uint64(chainID), "timestamp", "")
		if !ok {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		timestamp := t.(time.Time)
		if time.Since(timestamp) > time.Minute {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func (c *JsonRpcCache) pollBlocksData() {
	bucket := time.Now().Truncate(time.Second * 10).Unix()

	for {
		// wait for the next bucket
		<-time.After(time.Duration(bucket-time.Now().Unix()) * time.Second)

		agents, err := c.botRegistry.LoadAssignedBots()
		if err == nil && len(agents) == 0 {
			log.Warn("No agents assigned to the scanner, skipping polling for BlocksData")
			bucket += 10 // 10 seconds
			continue
		}

		log.Infof("Polling BlocksData from dispatcher. bucket: %d", bucket)

		// blocksDataClient internally retries on failure and to not block on the retry, we run it in a goroutine
		go func(b int64) {
			blocksData, err := c.blocksDataClient.GetBlocksData(b)
			if err != nil {
				c.msgClient.PublishProto(messaging.SubjectMetricAgent,
					metrics.CreateEventMetric(time.Now(), "system", domain.MetricJSONRPCCachePollError, err.Error()))
				log.WithError(err).Errorf("Failed to get BlocksData from dispatcher. bucket: %d", b)
				return
			}

			c.msgClient.PublishProto(
				messaging.SubjectMetricAgent, &protocol.AgentMetricList{
					Metrics: []*protocol.AgentMetric{
						metrics.CreateSystemMetric(domain.MetricJSONRPCCachePollSuccess, float64(len(blocksData.Blocks)), fmt.Sprintf("%d", b)),
						metrics.CreateSystemMetric(domain.MetricJSONRPCCacheSize, float64(c.cache.cache.ItemCount()), ""),
					},
				},
			)

			log.Infof("Added BlocksData to local cache. bucket: %d blocksData: %d", b, len(blocksData.Blocks))
			c.cache.Append(blocksData)
		}(bucket)

		bucket += 10 // 10 seconds
	}
}
