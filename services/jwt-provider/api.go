package jwt_provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/forta-network/forta-node/services/jwt-provider/provider"
	"net"
	"net/http"
	"time"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-node/config"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	errBadCreateMessage = "bad create jwt message body"
)

// JWTAPI provides jwt tokens to bots, signed with node's private key..
type JWTAPI struct {
	provider provider.JWTProvider
	lastErr  health.ErrorTracker

	srv *http.Server
}

func NewJWTAPI(
	cfg config.Config,
) (*JWTAPI, error) {
	p, err := provider.NewJWTProvider(cfg)
	if err != nil {
		return nil, err
	}

	return &JWTAPI{
		provider: p,
	}, nil
}

// Start spawns a jwt provider routine and returns.
func (j *JWTAPI) Start() error {
	return j.StartWithContext(context.Background())
}

func (j *JWTAPI) Stop() error {
	return j.srv.Close()
}

// StartWithContext subscribe to bot updates and spawn a Bot JWT Provider http server.
func (j *JWTAPI) StartWithContext(ctx context.Context) error {
	addr := fmt.Sprintf(":%s", config.DefaultJWTProviderPort)

	// setup routes
	r := mux.NewRouter()
	r.HandleFunc("/create", j.handleJwtRequest).Methods(http.MethodPost)

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

func (j *JWTAPI) listenAndServeWithContext(ctx context.Context) error {
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

func (j *JWTAPI) testAPI(_ context.Context) {
	j.lastErr.Set(nil)
}

func (j *JWTAPI) apiHealthChecker(ctx context.Context) {
	j.testAPI(ctx)
	ticker := time.NewTicker(time.Minute * 5)
	for range ticker.C {
		j.testAPI(ctx)
	}
}

func (j *JWTAPI) Name() string {
	return "jwt-provider"
}

func (j *JWTAPI) Health() health.Reports {
	return health.Reports{
		j.lastErr.GetReport("api"),
	}
}

func (j *JWTAPI) handleJwtRequest(w http.ResponseWriter, req *http.Request) {
	var msg CreateJWTMessage
	if req.Body != http.NoBody {
		err := json.NewDecoder(req.Body).Decode(&msg)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(w, errBadCreateMessage)
			return
		}
	}

	ipAddr, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		j.lastErr.Set(err)
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintf(w, "can't extract ip from request: %s", req.RemoteAddr)
		return
	}

	jwt, err := j.provider.CreateJWTFromIP(req.Context(), ipAddr, msg.Claims)
	if err == provider.ErrCannotFindBotForIP {
		j.lastErr.Set(err)
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprintf(w, "can't find bot id from request source %s, err: %v", ipAddr, err)
		return
	}

	if err != nil {
		j.lastErr.Set(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, "cannot create jwt")
		return
	}

	resp, err := json.Marshal(CreateJWTResponse{Token: jwt})

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "%s", resp)
}
