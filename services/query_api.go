package services

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"

	"OpenZeppelin/fortify-node/protocol"
	"OpenZeppelin/fortify-node/store"
)

// AlertApi allows retrieval of alerts from the database
type AlertApi struct {
	ctx   context.Context
	store store.AlertStore
	cfg   AlertApiConfig
}

const paramStartDate = "startDate"
const paramEndDate = "endDate"

const defaultSinceDate = -24 * 7 * time.Hour // last 7 days
const defaultPageLimit = 100
const maxPageLimit = 1000

type AlertApiConfig struct {
	Port int
}

func getDateParam(r *http.Request, name string, defaultTime time.Time) (time.Time, error) {
	dtStr := r.URL.Query().Get(name)
	if dtStr == "" {
		return defaultTime, nil
	}
	return time.Parse(time.RFC3339, dtStr)
}

func parseQueryRequest(r *http.Request) (*store.AlertQueryRequest, error) {
	now := time.Now()
	startDate, err := getDateParam(r, paramStartDate, now.Add(defaultSinceDate))
	if err != nil {
		return nil, fmt.Errorf("startDate must be in RFC3339 format")
	}
	endDate, err := getDateParam(r, paramEndDate, now)
	if err != nil {
		return nil, fmt.Errorf("endDate must be in RFC3339 format")
	}

	request := &store.AlertQueryRequest{
		FromTime:  startDate,
		ToTime:    endDate,
		PageToken: r.URL.Query().Get("pageToken"),
		Limit:     defaultPageLimit,
	}

	limit := r.URL.Query().Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return nil, fmt.Errorf("limit must be an integer")
		}
		if l > maxPageLimit {
			return nil, fmt.Errorf("limit must be less than 1000")
		}
		request.Limit = l
	}
	return request, nil
}

func (t *AlertApi) getAlerts(w http.ResponseWriter, r *http.Request) {
	queryReq, err := parseQueryRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	alerts, err := t.store.QueryAlerts(queryReq)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	// convert to protocol
	resp := &protocol.AlertResponse{
		Alerts:        alerts.Alerts,
		NextPageToken: alerts.NextPageToken,
	}
	m := jsonpb.Marshaler{EmitDefaults: true}

	if err := m.Marshal(w, resp); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}

func (t *AlertApi) Start() error {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/alerts", t.getAlerts)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})
	return http.ListenAndServe(fmt.Sprintf(":%d", t.cfg.Port), c.Handler(router))
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
