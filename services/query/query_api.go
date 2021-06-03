package query

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"

	"github.com/OpenZeppelin/fortify-node/protocol"
	"github.com/OpenZeppelin/fortify-node/store"
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

const prefixNot = "{not}"

var reservedParams = []string{"startDate", "endDate", "limit", "pageToken", "sort"}

type AlertApiConfig struct {
	Port int
}

func getDateParam(r *http.Request, name string, defaultTime time.Time) (time.Time, error) {
	dtStr := r.URL.Query().Get(name)
	if dtStr == "" {
		return defaultTime, nil
	}
	ms, err := strconv.Atoi(dtStr)
	if err != nil || ms < 1e12 {
		return time.Time{}, fmt.Errorf("%s must be in milliseconds format", name)
	}
	return time.Unix(0, int64(ms)*int64(time.Millisecond)), nil
}

func isReservedParam(param string) bool {
	for _, nfp := range reservedParams {
		if nfp == param {
			return true
		}
	}
	return false
}

func parseFilterCriteria(r *http.Request) ([]*store.FilterCriterion, error) {
	var filters []*store.FilterCriterion
	values := r.URL.Query()
	for k, v := range values {
		// skip startDate, endDate, limit, pageToken
		if !isReservedParam(k) {
			if len(v) > 1 {
				return nil, fmt.Errorf("%s cannot be array", k)
			}
			val := v[0]
			op := store.Equals
			if strings.HasPrefix(val, prefixNot) {
				op = store.NotEquals
				val = strings.TrimPrefix(val, prefixNot)
			}
			valList := strings.Split(val, ",")
			filters = append(filters, &store.FilterCriterion{
				Operator: op,
				Field:    k,
				Values:   valList,
			})
		}
	}
	return filters, nil
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
	filters, err := parseFilterCriteria(r)
	if err != nil {
		return nil, err
	}

	var reverse bool
	s := r.URL.Query().Get("sort")
	if s == "asc" {
		reverse = false
	} else if s == "desc" {
		reverse = true
	} else if s != "" {
		return nil, fmt.Errorf("sort must be either asc or desc (default asc)")
	}

	request := &store.AlertQueryRequest{
		StartTime: startDate,
		EndTime:   endDate,
		PageToken: r.URL.Query().Get("pageToken"),
		Limit:     defaultPageLimit,
		Criteria:  filters,
		Reverse:   reverse,
	}

	limit := r.URL.Query().Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return nil, fmt.Errorf("limit must be an integer")
		}
		if l > maxPageLimit {
			return nil, fmt.Errorf("limit cannot exceed %d", maxPageLimit)
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
	log.Infof(queryReq.Json())

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
