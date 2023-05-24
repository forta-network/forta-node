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

const (
	botIDKey contextKey = iota
	botOwnerKey
	isScannerKey
)

const claimKeyBotOwner = "bot-owner"

// PublicAPIProxy proxies requests from agents to json-rpc endpoint
type PublicAPIProxy struct {
	ctx       context.Context
	cfg       config.PublicAPIProxyConfig
	Key       *keystore.Key
	msgClient clients.MessageClient

	server *http.Server

	rateLimiter ratelimiter.RateLimiter

	lastErr       health.ErrorTracker
	authenticator clients.IPAuthenticator
}

func (p *PublicAPIProxy) newReverseProxy() http.Handler {
	apiURL, err := url.Parse(p.cfg.Url)
	if err != nil {
		logrus.WithError(err).Panic("bad public api proxy configuration")
	}

	rp := httputil.NewSingleHostReverseProxy(apiURL)

	d := rp.Director
	rp.Director = func(r *http.Request) {
		d(r)
		r.Host = apiURL.Host
		r.URL.Host = apiURL.Host
		for h, v := range p.cfg.Headers {
			r.Header.Set(h, v)
		}
		r.Header.Set("User-Agent", "forta-scan-node")
	}

	return rp
}

func (p *PublicAPIProxy) createPublicAPIProxyHandler() http.Handler {
	c := cors.New(
		cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
		},
	)
	return p.authMiddleware(p.metricMiddleware(c.Handler(p.newReverseProxy())))
}

func (p *PublicAPIProxy) metricMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			t := time.Now()
			botID, _, isScanner, foundAgent := getBotFromContext(req.Context())
			if foundAgent && !isScanner && p.rateLimiter.ExceedsLimit(botID) {
				writeTooManyReqsErr(w, req)
				p.msgClient.PublishProto(
					messaging.SubjectMetricAgent, &protocol.AgentMetricList{
						Metrics: metrics.GetPublicAPIMetrics(botID, t, 0, 1, 0),
					},
				)
				return
			}

			h.ServeHTTP(w, req)

			if foundAgent {
				duration := time.Since(t)
				p.msgClient.PublishProto(
					messaging.SubjectMetricAgent, &protocol.AgentMetricList{
						Metrics: metrics.GetPublicAPIMetrics(botID, t, 1, 0, duration),
					},
				)
			}
		},
	)
}

func (p *PublicAPIProxy) authMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			botReq, err := p.authenticateRequest(req)
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

func (p *PublicAPIProxy) authenticateRequest(req *http.Request) (*http.Request, error) {
	containerName, err := p.authenticator.FindContainerNameFromRemoteAddr(req.Context(), req.RemoteAddr)
	if err != nil {
		return req, err
	}

	var botID, botOwner string

	isScanner := false
	// combiner feed authorization
	if containerName == config.DockerScannerContainerName {
		isScanner = true
		botID = req.Header.Get("bot-id")
		botOwner = req.Header.Get("bot-owner")
	} else {
		// bot authorization
		agentConfig, err := p.authenticator.FindAgentByContainerName(containerName)
		// request source is not a bot
		if err != nil {
			return req, err
		}

		botID = agentConfig.ID
		botOwner = agentConfig.Owner
	}

	// set authorization values as context to use in next middlewares
	ctxWithBot := context.WithValue(req.Context(), botIDKey, botID)
	ctxWithBot = context.WithValue(ctxWithBot, botOwnerKey, botOwner)
	ctxWithBot = context.WithValue(ctxWithBot, isScannerKey, isScanner)

	botReq := req.WithContext(ctxWithBot)

	return botReq, nil
}

func (p *PublicAPIProxy) setAuthBearer(r *http.Request) {
	log := logrus.WithField("addr", r.RemoteAddr)
	botID, botOwner, _, ok := getBotFromContext(r.Context())
	if !ok {
		return
	}

	claims := map[string]interface{}{claimKeyBotOwner: botOwner}

	jwtToken, err := jwt_provider.CreateBotJWT(p.Key, botID, claims)
	if err != nil {
		log.WithError(err).Warn("can't create bot jwt")
		return
	}

	bearerToken := fmt.Sprintf("Bearer %s", jwtToken)

	r.Header.Set("Authorization", bearerToken)
}

func (p *PublicAPIProxy) Start() error {
	p.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", config.DefaultPublicAPIProxyPort),
		Handler: p.createPublicAPIProxyHandler(),
	}

	utils.GoListenAndServe(p.server)

	return nil
}

func (p *PublicAPIProxy) Stop() error {
	if p.server != nil {
		return p.server.Close()
	}
	return nil
}

func (p *PublicAPIProxy) Name() string {
	return "public-api-proxy"
}

func getBotFromContext(ctx context.Context) (string, string, bool, bool) {
	botIdVal := ctx.Value(botIDKey)
	if botIdVal == nil {
		return "", "", false, false
	}

	botID, ok := botIdVal.(string)
	if !ok {
		return "", "", false, false
	}

	botOwnerVal := ctx.Value(botOwnerKey)
	if botOwnerVal == nil {
		return "", "", false, false
	}

	botOwner, ok := botOwnerVal.(string)
	if !ok {
		return "", "", false, false
	}

	isScannerVal := ctx.Value(isScannerKey)
	if botOwnerVal == nil {
		return "", "", false, false
	}

	isScanner, ok := isScannerVal.(bool)
	if !ok {
		return "", "", false, false
	}

	return botID, botOwner, isScanner, ok
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

	botAuthenticator, err := clients.NewBotAuthenticator(ctx)
	if err != nil {
		return nil, err
	}

	msgClient := messaging.NewClient("public-api", fmt.Sprintf("%s:%s", config.DockerNatsContainerName, config.DefaultNatsPort))

	rateLimiting := cfg.PublicAPIProxy.RateLimitConfig
	if rateLimiting == nil {
		rateLimiting = &config.RateLimitConfig{Rate: 1000, Burst: 1}
	}

	return newPublicAPIProxy(ctx, cfg.PublicAPIProxy, botAuthenticator, ratelimiter.NewRateLimiter(rateLimiting.Rate, rateLimiting.Burst), key, msgClient)
}

func newPublicAPIProxy(
	ctx context.Context, cfg config.PublicAPIProxyConfig, botAuthenticator clients.IPAuthenticator, rateLimiter ratelimiter.RateLimiter, key *keystore.Key, msgClient clients.MessageClient,
) (
	*PublicAPIProxy, error,
) {
	return &PublicAPIProxy{
		ctx:           ctx,
		cfg:           cfg,
		authenticator: botAuthenticator,
		msgClient:     msgClient,
		Key:           key,
		rateLimiter: rateLimiter,
	}, nil
}
