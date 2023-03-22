package public_api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/ethereum"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-core-go/security"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/forta-network/forta-node/clients/botauth"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/clients/ratelimiter"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/metrics"
	jwt_provider "github.com/forta-network/forta-node/services/jwt-provider"
	"github.com/rs/cors"
)

// PublicAPIProxy proxies requests from agents to json-rpc endpoint
type PublicAPIProxy struct {
	ctx context.Context
	cfg config.JsonRpcConfig
	Key *keystore.Key

	server *http.Server

	rateLimiter *ratelimiter.RateLimiter

	lastErr          health.ErrorTracker
	botAuthenticator *botauth.BotAuthenticator
}

func (p *PublicAPIProxy) Start() error {
	p.botAuthenticator.RegisterMessageHandlers()

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

		bot, found := p.botAuthenticator.FindAgentFromRemoteAddr(r.Host)
		if found {
			return
		}

		claims := map[string]interface{}{"owner": bot.Owner}

		jwtToken, err := jwt_provider.CreateBotJWT(p.Key, bot.ID, claims)
		if err != nil {
			return
		}

		bearerToken := fmt.Sprintf("Bearer %s", jwtToken)

		r.Header.Set("Authorization", bearerToken)
	}

	c := cors.New(
		cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
		},
	)

	p.server = &http.Server{
		Addr:    ":8545",
		Handler: p.metricHandler(c.Handler(rp)),
	}
	utils.GoListenAndServe(p.server)
	return nil
}

func (p *PublicAPIProxy) metricHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			t := time.Now()
			agentConfig, foundAgent := p.botAuthenticator.FindAgentFromRemoteAddr(req.RemoteAddr)
			if foundAgent && p.rateLimiter.ExceedsLimit(agentConfig.ID) {
				writeTooManyReqsErr(w, req)
				p.botAuthenticator.MsgClient().PublishProto(
					messaging.SubjectMetricAgent, &protocol.AgentMetricList{
						Metrics: metrics.GetJSONRPCMetrics(*agentConfig, t, 0, 1, 0),
					},
				)
				return
			}

			h.ServeHTTP(w, req)

			if foundAgent {
				duration := time.Since(t)
				p.botAuthenticator.MsgClient().PublishProto(
					messaging.SubjectMetricAgent, &protocol.AgentMetricList{
						Metrics: metrics.GetJSONRPCMetrics(*agentConfig, t, 1, 0, duration),
					},
				)
			}
		},
	)
}

func (p *PublicAPIProxy) Stop() error {
	if p.server != nil {
		return p.server.Close()
	}
	return nil
}

func (p *PublicAPIProxy) Name() string {
	return "json-rpc-proxy"
}

// Health implements health.Reporter interface.
func (p *PublicAPIProxy) Health() health.Reports {
	return health.Reports{
		p.lastErr.GetReport("api"),
	}
}

func (p *PublicAPIProxy) apiHealthChecker() {
	p.testAPI()
	ticker := time.NewTicker(time.Minute * 5)
	for range ticker.C {
		p.testAPI()
	}
}

func (p *PublicAPIProxy) testAPI() {
	err := ethereum.TestAPI(p.ctx, "http://localhost:8545")
	p.lastErr.Set(err)
}

func NewPublicAPIProxy(ctx context.Context, cfg config.Config) (*PublicAPIProxy, error) {
	key, err := security.LoadKey(config.DefaultContainerKeyDirPath)
	if err != nil {
		return nil, err
	}

	botAuthenticator, err := botauth.NewBotAuthenticator(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &PublicAPIProxy{
		ctx:              ctx,
		botAuthenticator: botAuthenticator,
		Key:              key,
		rateLimiter: ratelimiter.NewRateLimiter(
			1000,
			1,
		),
	}, nil
}
