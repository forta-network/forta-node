package jwt_provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/security"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/config"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// JWTProvider provides jwt tokens to bots, signed with node's private key..
type JWTProvider struct {
	botConfigs      []config.AgentConfig
	botConfigsMutex sync.RWMutex

	// to match request ip <-> bot id
	dockerClient clients.DockerClient

	cfg *JWTProviderConfig

	lastErr health.ErrorTracker

	srv *http.Server
}

type JWTProviderConfig struct {
	Key    *keystore.Key
	Config config.Config
}

func NewJWTProvider(
	cfg config.Config,
) (*JWTProvider, error) {
	key, err := security.LoadKey(config.DefaultContainerKeyDirPath)
	if err != nil {
		return nil, err
	}

	return initProvider(
		&JWTProviderConfig{
			Key:    key,
			Config: cfg,
		},
	)
}

func initProvider(cfg *JWTProviderConfig) (*JWTProvider, error) {
	globalClient, err := clients.NewDockerClient("")
	if err != nil {
		return nil, fmt.Errorf("failed to create the global docker client: %v", err)
	}

	return &JWTProvider{dockerClient: globalClient, cfg: cfg}, nil
}

// Start spawns a jwt provider routine and returns.
func (j *JWTProvider) Start() error {
	return j.StartWithContext(context.Background())
}

func (j *JWTProvider) Stop() error {
	return j.srv.Close()
}

// StartWithContext subscribe to bot updates and spawn a Bot JWT Provider http server.
func (j *JWTProvider) StartWithContext(ctx context.Context) error {
	if j.cfg.Config.JWTProvider.Addr == "" {
		j.cfg.Config.JWTProvider.Addr = fmt.Sprintf(":%s", config.DefaultJWTProviderPort)
	}

	// setup routes
	r := mux.NewRouter()
	r.HandleFunc("/create", j.createJWTHandler).Methods(http.MethodPost)

	j.srv = &http.Server{
		Addr:    j.cfg.Config.JWTProvider.Addr,
		Handler: r,
	}

	go func() {
		err := j.listenAndServeWithContext(ctx)
		if err != nil {
			logrus.WithError(err).Panic("server error")
		}
	}()

	return nil
}

func (j *JWTProvider) listenAndServeWithContext(ctx context.Context) error {
	errChan := make(chan error)

	go func() {
		logrus.Infof("Starting Bot JWT Provider Service on: %s", j.srv.Addr)
		err := j.srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	// gracefully handle stopping server
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		_ = j.srv.Close()
		return nil
	}
}

// agentIDReverseLookup reverse lookup from ip to agent id.
func (j *JWTProvider) agentIDReverseLookup(ctx context.Context, ipAddr string) (string, error) {
	container, err := j.findContainerByIP(ctx, ipAddr)
	if err != nil {
		return "", err
	}

	botID, err := j.extractBotIDFromContainer(ctx, container)
	if err != nil {
		return "", err
	}

	return botID, nil
}

const envPrefix = config.EnvFortaBotID + "="

func (j *JWTProvider) extractBotIDFromContainer(ctx context.Context, container types.Container) (string, error) {
	// container struct doesn't have the "env" information, inspection required.
	c, err := j.dockerClient.InspectContainer(ctx, container.ID)
	if err != nil {
		return "", err
	}

	// find the env variable with bot id
	for _, s := range c.Config.Env {
		if env := strings.SplitAfter(s, envPrefix); len(env) == 2 {
			return env[1], nil
		}
	}

	return "", fmt.Errorf("can't extract bot id from container")
}

func (j *JWTProvider) findContainerByIP(ctx context.Context, ipAddr string) (types.Container, error) {
	containers, err := j.dockerClient.GetContainers(ctx)
	if err != nil {
		return types.Container{}, err
	}

	// find the container that has the same ip
	for _, container := range containers {
		for _, network := range container.NetworkSettings.Networks {
			if network.IPAddress == ipAddr {
				return container, nil
			}
		}
	}
	return types.Container{}, fmt.Errorf("can't find container %s", ipAddr)
}

func (j *JWTProvider) testAPI(_ context.Context) {
	j.lastErr.Set(nil)
}

func (j *JWTProvider) apiHealthChecker(ctx context.Context) {
	j.testAPI(ctx)
	ticker := time.NewTicker(time.Minute * 5)
	for range ticker.C {
		j.testAPI(ctx)
	}
}

func (j *JWTProvider) Name() string {
	return "jwt-provider"
}

func (j *JWTProvider) Health() health.Reports {
	return health.Reports{
		j.lastErr.GetReport("api"),
	}
}

// requestHash used for "hash" claim in JWT token
func requestHash(uri string, payload []byte) common.Hash {
	requestStr := fmt.Sprintf("%s%s", uri, payload)

	return crypto.Keccak256Hash([]byte(requestStr))
}

// CreateBotJWT returns a bot JWT token. Basically security.ScannerJWT with bot&request info.
func CreateBotJWT(key *keystore.Key, agentID string, claims map[string]interface{}) (string, error) {
	if claims == nil {
		claims = make(map[string]interface{})
	}

	claims["bot-id"] = agentID

	return security.CreateScannerJWT(key, claims)
}
