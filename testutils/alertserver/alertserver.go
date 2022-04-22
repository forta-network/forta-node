package alertserver

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/forta-network/forta-core-go/utils/apiutils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// AlertServer is a fake alert server.
type AlertServer struct {
	ctx    context.Context
	cancel func()
	port   int
	router *mux.Router

	knownBatches map[string][]byte
	mu           sync.RWMutex
}

// New creates a new alert server.
func New(ctx context.Context, port int) *AlertServer {
	alertServer := &AlertServer{
		port:         port,
		knownBatches: make(map[string][]byte),
	}
	alertServer.ctx, alertServer.cancel = context.WithCancel(ctx)

	r := mux.NewRouter()
	r.HandleFunc("/batch/{ref}", alertServer.AddAlert).Methods("POST")
	alertServer.router = r

	return alertServer
}

// Start starts the server.
func (as *AlertServer) Start() {
	apiutils.ListenAndServe(as.ctx, &http.Server{
		Handler:      as.router,
		Addr:         fmt.Sprintf("0.0.0.0:%d", as.port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}, "started alert server")
}

// Close closes the server.
func (as *AlertServer) Close() error {
	as.cancel()
	return nil
}

func (as *AlertServer) GetAlert(ref string) ([]byte, bool) {
	as.mu.RLock()
	defer as.mu.RUnlock()
	b, ok := as.knownBatches[ref]
	return b, ok
}

func (as *AlertServer) AddAlert(w http.ResponseWriter, r *http.Request) {
	as.mu.Lock()
	defer as.mu.Unlock()

	vars := mux.Vars(r)
	ref := vars["ref"]
	b, _ := ioutil.ReadAll(r.Body)
	logrus.WithField("ref", ref).Info("received alert: ", string(b))
	as.knownBatches[ref] = b
	return
}
