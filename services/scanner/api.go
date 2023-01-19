package scanner

import (
	"context"
	"net/http"
	"strconv"

	"github.com/forta-network/forta-core-go/feeds"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/goccy/go-json"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

// API allows triggering things on scanner
type API struct {
	ctx     context.Context
	started bool
	feed    feeds.BlockFeed
	server  *http.Server
}

type Message struct {
	Message string `json:"message"`
}

func message(str string) []byte {
	b, _ := json.Marshal(Message{Message: str})
	return b
}

func writeError(w http.ResponseWriter, code int, str string) {
	w.WriteHeader(code)
	if _, err := w.Write(message(str)); err != nil {
		log.WithError(err).Errorf("error writing: %s", str)
	}
}

func writeMessage(w http.ResponseWriter, str string) {
	w.WriteHeader(200)
	if _, err := w.Write(message(str)); err != nil {
		log.WithError(err).Errorf("error writing: %s", str)
	}
}

func (a *API) startBlocks(w http.ResponseWriter, r *http.Request) {
	if a.feed.IsStarted() {
		writeMessage(w, "already started")
	} else {
		start := r.URL.Query().Get("start")
		start64, err := strconv.ParseInt(start, 10, 64)
		if err != nil {
			writeError(w, 400, "?start is required and must be integer")
			return
		}

		end := r.URL.Query().Get("end")
		end64, err := strconv.ParseInt(end, 10, 64)
		if err != nil {
			writeError(w, 400, "?end is required and must be integer")
			return
		}

		rate := r.URL.Query().Get("rate")
		rate64, err := strconv.ParseInt(rate, 10, 64)
		if err != nil {
			writeError(w, 400, "?end is required and must be integer")
			return
		}
		a.feed.StartRange(start64, end64, rate64)
		writeMessage(w, "ok")
	}
}

func (t *API) Start() error {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/start", t.startBlocks)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	t.server = &http.Server{
		Addr:    ":80",
		Handler: c.Handler(router),
	}
	utils.GoListenAndServe(t.server)
	return nil
}

func (t *API) Stop() error {
	if t.server != nil {
		return t.server.Close()
	}
	return nil
}

func (t *API) Name() string {
	return "scanner-api"
}

func NewScannerAPI(ctx context.Context, feed feeds.BlockFeed) *API {
	return &API{
		ctx:  ctx,
		feed: feed,
	}
}
