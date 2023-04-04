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
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/clients/ratelimiter"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/metrics"
	jwt_provider "github.com/forta-network/forta-node/services/jwt-provider"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

type contextKey int

const authenticatedBotKey contextKey = 0
const claimKeyBotOwner = "bot-owner"
// PublicAPIProxy proxies requests from agents to json-rpc endpoint
type PublicAPIProxy struct {
	ctx       context.Context
	cfg       config.PublicAPIProxyConfig
	Key       *keystore.Key
	msgClient clients.MessageClient

	server *http.Server

	rateLimiter *ratelimiter.RateLimiter

	lastErr          health.ErrorTracker
	botAuthenticator clients.BotAuthenticator
}

func (p *PublicAPIProxy) Start() error {
	apiURL, err := url.Parse(p.cfg.Url)
	if err != nil {
		return err
	}

	rp := httputil.NewSingleHostReverseProxy(apiURL)

	d := rp.Director
	rp.Director = func(r *http.Request) {
		d(r)
		r.Host = apiURL.Host
		r.URL.Host = apiURL.Host
		r.Header.Set("User-Agent","forta-scan-node")
		for h, v := range p.cfg.Headers {
			r.Header.Set(h, v)
		}
	}

	c := cors.New(
		cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
		},
	)

	p.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", config.DefaultPublicAPIProxyPort),
		Handler: p.authMiddleware(p.metricMiddleware(c.Handler(rp))),
	}

	utils.GoListenAndServe(p.server)

	return nil
}

func (p *PublicAPIProxy) metricMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			t := time.Now()
			agentConfig, foundAgent := getBotFromContext(req.Context())
			if foundAgent && p.rateLimiter.ExceedsLimit(agentConfig.ID) {
				writeTooManyReqsErr(w, req)
				p.msgClient.PublishProto(
					messaging.SubjectMetricAgent, &protocol.AgentMetricList{
						Metrics: metrics.GetJSONRPCMetrics(*agentConfig, t, 0, 1, 0),
					},
				)
				return
			}

			h.ServeHTTP(w, req)

			if foundAgent {
				duration := time.Since(t)
				p.msgClient.PublishProto(
					messaging.SubjectMetricAgent, &protocol.AgentMetricList{
						Metrics: metrics.GetPublicAPIMetrics(*agentConfig, t, 1, 0, duration),
					},
				)
			}
		},
	)
}

func (p *PublicAPIProxy) authMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			botReq, err := p.authenticateBotRequest(req)
			if err != nil {
				logrus.WithError(err).Warn("failed to authenticate bot request")
				writeAuthError(w, req)
				return
			}

			p.setAuthBearer(botReq)

			h.ServeHTTP(w, botReq)
		},
	)
}

func (p *PublicAPIProxy) authenticateBotRequest(req *http.Request) (*http.Request, error) {
	agentConfig, err := p.botAuthenticator.FindAgentFromRemoteAddr(req.RemoteAddr)
	// request source is not a bot
	if err != nil {
		return req, err
	}

	ctxWithBoth := context.WithValue(req.Context(), authenticatedBotKey, agentConfig)
	botReq := req.WithContext(ctxWithBoth)
	return botReq, nil
}

func (p *PublicAPIProxy) setAuthBearer(r *http.Request) {
	log := logrus.WithField("addr", r.RemoteAddr)
	bot, ok := getBotFromContext(r.Context())
	if !ok {
		return
	}

	claims := map[string]interface{}{claimKeyBotOwner: bot.Owner}

	jwtToken, err := jwt_provider.CreateBotJWT(p.Key, bot.ID, claims)
	if err != nil {
		log.WithError(err).Warn("can't create bot jwt")
		return
	}

	bearerToken := fmt.Sprintf("Bearer %s", jwtToken)

	r.Header.Set("Authorization", bearerToken)
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

func getBotFromContext(ctx context.Context) (*config.AgentConfig, bool) {
	botCtxVal := ctx.Value(authenticatedBotKey)
	if botCtxVal == nil {
		return nil, false
	}

	bot, ok := botCtxVal.(*config.AgentConfig)

	return bot, ok
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

	botAuthenticator, err := clients.NewBotAuthenticator(ctx, cfg)
	if err != nil {
		return nil, err
	}

	msgClient := messaging.NewClient("public-api", fmt.Sprintf("%s:%s", config.DockerNatsContainerName, config.DefaultNatsPort))

	rateLimiting := cfg.PublicAPIProxy.RateLimitConfig
	if rateLimiting == nil {
		rateLimiting = &config.RateLimitConfig{Rate: 1000, Burst: 1}
	}

	return &PublicAPIProxy{
		ctx:              ctx,
		cfg:              cfg.PublicAPIProxy,
		botAuthenticator: botAuthenticator,
		msgClient:        msgClient,
		Key:              key,
		// TODO: adjust rate limiting
		rateLimiter: ratelimiter.NewRateLimiter(
			rateLimiting.Rate, rateLimiting.Burst,
		),
	}, nil
}
