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
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"

	"github.com/forta-protocol/forta-node/clients"
	"github.com/forta-protocol/forta-node/clients/health"
	"github.com/forta-protocol/forta-node/clients/messaging"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/ethereum"
	"github.com/forta-protocol/forta-node/metrics"
	"github.com/forta-protocol/forta-node/protocol"
	"github.com/forta-protocol/forta-node/utils"
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

	lastErr health.ErrorTracker
}

func (p *JsonRpcProxy) Start() error {
	log.Infof("Starting %s", p.Name())
	rpcUrl, err := url.Parse(p.cfg.Url)
	if err != nil {
		return err
	}
	rp := httputil.NewSingleHostReverseProxy(rpcUrl)

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
		h.ServeHTTP(w, req)
		duration := time.Since(t)
		agentConfig, ok := p.findAgentFromRemoteAddr(req.RemoteAddr)
		if ok {
			p.msgClient.PublishProto(messaging.SubjectMetricAgent, &protocol.AgentMetricList{
				Metrics: metrics.GetJSONRPCMetrics(*agentConfig, t, duration),
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

	containerName := agentContainer.Names[0]
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
	log.Infof("Stopping %s", p.Name())
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
	if err != nil {
		p.lastErr.Set(err)
	}
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
	return &JsonRpcProxy{
		ctx:          ctx,
		cfg:          jCfg,
		dockerClient: globalClient,
		msgClient:    msgClient,
	}, nil
}
