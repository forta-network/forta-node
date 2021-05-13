package services

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"

	"OpenZeppelin/fortify-node/config"
)

// JsonRpcProxy proxies requests from agents to json-rpc endpoint
type JsonRpcProxy struct {
	ctx context.Context
	cfg config.JsonRpcProxyConfig
}

type accessLogger struct{}

func (accessLogger) RoundTrip(r *http.Request) (*http.Response, error) {
	b, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		return nil, err
	}
	log.Info(string(b))
	return http.DefaultTransport.RoundTrip(r)
}

func (p *JsonRpcProxy) Start() error {
	log.Infof("Starting %s", p.Name())
	rpcUrl, err := url.Parse(p.cfg.Ethereum.JsonRpcUrl)
	if err != nil {
		return err
	}
	rp := httputil.NewSingleHostReverseProxy(rpcUrl)
	rp.Transport = accessLogger{}

	d := rp.Director
	rp.Director = func(r *http.Request) {
		d(r)
		r.Host = rpcUrl.Host
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		//TODO: this is where we can validate methods
		for h, v := range p.cfg.Ethereum.Headers {
			req.Header.Set(h, v)
		}
		rp.ServeHTTP(w, req)
	})
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})
	return http.ListenAndServe(":8545", c.Handler(router))
}

func (p *JsonRpcProxy) Stop() error {
	log.Infof("Stopping %s", p.Name())
	return nil
}

func (p *JsonRpcProxy) Name() string {
	return "JsonRpcProxy"
}

func NewJsonRpcProxy(ctx context.Context, cfg config.Config) (*JsonRpcProxy, error) {
	return &JsonRpcProxy{
		ctx: ctx,
		cfg: cfg.JsonRpcProxy,
	}, nil
}
