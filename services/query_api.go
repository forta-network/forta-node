package services

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"OpenZeppelin/zephyr-node/protocol"
	"OpenZeppelin/zephyr-node/store"
)

// AlertApi allows retrieval of alerts from the database
type AlertApi struct {
	ctx   context.Context
	store store.AlertStore
	cfg   AlertApiConfig
}

const defaultSinceDate = -1 * 24 * 7 * time.Hour

type AlertApiConfig struct {
	Port int
}

func (t *AlertApi) getAlerts(w http.ResponseWriter, r *http.Request) {
	//TODO: make dates a query param
	endDate := time.Now()
	startDate := endDate.Add(defaultSinceDate)
	alerts, err := t.store.GetAlerts(startDate, endDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp := &protocol.AlertResponse{
		Alerts: alerts,
	}
	m := jsonpb.Marshaler{EmitDefaults: true}

	if err := m.Marshal(w, resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (t *AlertApi) Start() error {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/alerts", t.getAlerts)
	return http.ListenAndServe(fmt.Sprintf(":%d", t.cfg.Port), router)
}

func (t *AlertApi) Stop() error {
	log.Infof("Stopping %s", t.Name())
	return nil
}

func (t *AlertApi) Name() string {
	return "AlertApi"
}

func NewAlertApi(ctx context.Context, store store.AlertStore, cfg AlertApiConfig) (*AlertApi, error) {
	return &AlertApi{
		ctx:   ctx,
		store: store,
		cfg:   cfg,
	}, nil
}
