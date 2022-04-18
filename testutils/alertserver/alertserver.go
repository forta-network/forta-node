package alertserver

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/forta-network/forta-core-go/utils/apiutils"
	"github.com/gorilla/mux"
)

// AlertServer is a fake alert server.
type AlertServer struct {
	ctx    context.Context
	cancel func()
	port   int
	router *mux.Router

	knownBatches map[string][]byte
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
	r.HandleFunc("/batch/{ref}", alertServer.GetAlert).Methods("GET")
	alertServer.router = r

	return alertServer
}

// Start starts the server.
func (as *AlertServer) Start() {
	apiutils.ListenAndServe(as.ctx, &http.Server{
		Handler:      as.router,
		Addr:         fmt.Sprintf(":%d", as.port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}, "started alert server")
}

// Close closes the server.
func (as *AlertServer) Close() error {
	as.cancel()
	return nil
}

func (as *AlertServer) GetAlert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ref := vars["ref"]
	b, ok := as.knownBatches[ref]
	if !ok {
		apiutils.InternalError(w, "error getting batch")
		return
	}
	apiutils.WriteOKBody(w, b)
}

func (as *AlertServer) AddAlert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ref := vars["ref"]
	b, _ := ioutil.ReadAll(r.Body)
	as.knownBatches[ref] = b
	return
}
