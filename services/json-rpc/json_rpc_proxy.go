package json_rpc

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/ethereum"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-core-go/protocol/settings"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/metrics"
	"github.com/octo/retry"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

// JsonRpcProxy proxies requests from agents to json-rpc endpoint
type JsonRpcProxy struct {
	ctx          context.Context
	cfg          config.JsonRpcConfig
	server       *http.Server
	dockerClient clients.DockerClient
	msgClient    clients.MessageClient

	agentConfigs  []config.AgentConfig
	agentConfigMu sync.RWMutex

	rateLimiter *RateLimiter

	lastErr health.ErrorTracker
}

const maxRetries = 4
const retryDuration = time.Second * 60

func (p *JsonRpcProxy) Start() error {
	p.registerMessageHandlers()

	rpcUrl, err := url.Parse(p.cfg.Url)
	if err != nil {
		return err
	}
	rp := httputil.NewSingleHostReverseProxy(rpcUrl)

	rp.Transport = retry.NewTransport(
		http.DefaultTransport, retry.Attempts(maxRetries), retry.Timeout(retryDuration),
		// retry for a minute, every 15 seconds, 4 times
		retry.ExpBackoff{Max: retryDuration, Base: time.Second * 10, Factor: 2.0},
	)

	d := rp.Director
	rp.Director = func(r *http.Request) {
		d(r)
		r.Host = rpcUrl.Host
		r.URL = rpcUrl
		for h, v := range p.cfg.Headers {
			r.Header.Set(h, v)
		}
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	p.server = &http.Server{
		Addr:    ":8545",
		Handler: p.metricHandler(c.Handler(rp)),
	}
	utils.GoListenAndServe(p.server)
	return nil
}

func (p *JsonRpcProxy) metricHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		t := time.Now()
		agentConfig, foundAgent := p.findAgentFromRemoteAddr(req.RemoteAddr)
		if foundAgent && p.rateLimiter.ExceedsLimit(agentConfig.ID) {
			writeTooManyReqsErr(w, req)
			p.msgClient.PublishProto(messaging.SubjectMetricAgent, &protocol.AgentMetricList{
				Metrics: metrics.GetJSONRPCMetrics(*agentConfig, t, 0, 1, 0),
			})
			return
		}

			// handle request with 5xx code retry
		ctx, cancel := context.WithTimeout(req.Context(), time.Minute*3)
			defer cancel()
			h.ServeHTTP(w, req.WithContext(ctx))

		if foundAgent {
			duration := time.Since(t)
			p.msgClient.PublishProto(messaging.SubjectMetricAgent, &protocol.AgentMetricList{
				Metrics: metrics.GetJSONRPCMetrics(*agentConfig, t, 1, 0, duration),
			})
		}
	})
}

func (p *JsonRpcProxy) findAgentFromRemoteAddr(hostPort string) (*config.AgentConfig, bool) {
	containers, err := p.dockerClient.GetContainers(p.ctx)
	if err != nil {
		log.WithError(err).Error("failed to get the container list")
		return nil, false
	}
	ipAddr := strings.Split(hostPort, ":")[0]

	var agentContainer *types.Container
	for _, container := range containers {
		for _, network := range container.NetworkSettings.Networks {
			if network.IPAddress == ipAddr {
				agentContainer = &container
				break
			}
		}
		if agentContainer != nil {
			break
		}
	}
	if agentContainer == nil {
		log.WithField("agentIpAddr", ipAddr).Warn("could not found agent container from ip address")
		return nil, false
	}

	p.agentConfigMu.RLock()
	defer p.agentConfigMu.RUnlock()

	containerName := agentContainer.Names[0][1:]
	for _, agentConfig := range p.agentConfigs {
		if agentConfig.ContainerName() == containerName {
			return &agentConfig, true
		}
	}

	log.WithFields(log.Fields{
		"agentIpAddr":   ipAddr,
		"containerName": containerName,
	}).Warn("could not find agent config for container")
	return nil, false
}

func (p *JsonRpcProxy) handleAgentVersionsUpdate(payload messaging.AgentPayload) error {
	p.agentConfigMu.Lock()
	p.agentConfigs = payload
	p.agentConfigMu.Unlock()
	return nil
}

func (p *JsonRpcProxy) Stop() error {
	if p.server != nil {
		return p.server.Close()
	}
	return nil
}

func (p *JsonRpcProxy) Name() string {
	return "json-rpc-proxy"
}

// Health implements health.Reporter interface.
func (p *JsonRpcProxy) Health() health.Reports {
	return health.Reports{
		p.lastErr.GetReport("api"),
	}
}

func (p *JsonRpcProxy) apiHealthChecker() {
	p.testAPI()
	ticker := time.NewTicker(time.Minute * 5)
	for range ticker.C {
		p.testAPI()
	}
}

func (p *JsonRpcProxy) testAPI() {
	err := ethereum.TestAPI(p.ctx, "http://localhost:8545")
	p.lastErr.Set(err)
}

func (p *JsonRpcProxy) registerMessageHandlers() {
	p.msgClient.Subscribe(messaging.SubjectAgentsVersionsLatest, messaging.AgentsHandler(p.handleAgentVersionsUpdate))
}

func NewJsonRpcProxy(ctx context.Context, cfg config.Config) (*JsonRpcProxy, error) {
	jCfg := cfg.Scan.JsonRpc
	if len(cfg.JsonRpcProxy.JsonRpc.Url) > 0 {
		jCfg = cfg.JsonRpcProxy.JsonRpc
	}
	globalClient, err := clients.NewDockerClient("")
	if err != nil {
		return nil, fmt.Errorf("failed to create the global docker client: %v", err)
	}
	msgClient := messaging.NewClient("json-rpc-proxy", fmt.Sprintf("%s:%s", config.DockerNatsContainerName, config.DefaultNatsPort))

	rateLimiting := cfg.JsonRpcProxy.RateLimitConfig
	if rateLimiting == nil {
		rateLimiting = (*config.RateLimitConfig)(settings.GetChainSettings(cfg.ChainID).JsonRpcRateLimiting)
	}

	return &JsonRpcProxy{
		ctx:          ctx,
		cfg:          jCfg,
		dockerClient: globalClient,
		msgClient:    msgClient,
		rateLimiter: NewRateLimiter(
			rateLimiting.Rate,
			rateLimiting.Burst,
		),
	}, nil
}
