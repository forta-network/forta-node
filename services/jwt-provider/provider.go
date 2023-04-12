package jwt_provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/security"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/config"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// JWTProvider provides jwt tokens to bots, signed with node's private key..
type JWTProvider struct {
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

// CreateBotJWT returns a bot JWT token. Basically security.ScannerJWT with bot&request info.
func CreateBotJWT(key *keystore.Key, agentID string, claims map[string]interface{}) (string, error) {
	if key == nil {
		return "", fmt.Errorf("provider has no private key")
	}
	if claims == nil {
		claims = make(map[string]interface{})
	}

	claims["bot-id"] = agentID

	return security.CreateScannerJWT(key, claims)
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
	addr := fmt.Sprintf(":%s", config.DefaultJWTProviderPort)

	// setup routes
	r := mux.NewRouter()
	r.HandleFunc("/create", j.createJWTHandler).Methods(http.MethodPost)

	j.srv = &http.Server{
		Addr:    addr,
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