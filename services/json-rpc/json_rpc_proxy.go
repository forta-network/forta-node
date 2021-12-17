package json_rpc

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"

	"github.com/forta-protocol/forta-node/config"
)

// JsonRpcProxy proxies requests from agents to json-rpc endpoint
type JsonRpcProxy struct {
	ctx context.Context
	cfg config.JsonRpcConfig
}

func (p *JsonRpcProxy) Start() error {
	log.Infof("Starting %s", p.Name())
	rpcUrl, err := url.Parse(p.cfg.Url)
	if err != nil {
		return err
	}
	rp := httputil.NewSingleHostReverseProxy(rpcUrl)

	d := rp.Director
	rp.Director = func(r *http.Request) {
		d(r)
		r.Host = rpcUrl.Host
		r.URL = rpcUrl
		for h, v := range p.cfg.Headers {
			r.Header.Set(h, v)
		}
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	return http.ListenAndServe(":8545", c.Handler(rp))
}

func (p *JsonRpcProxy) Stop() error {
	log.Infof("Stopping %s", p.Name())
	return nil
}

func (p *JsonRpcProxy) Name() string {
	return "JsonRpcProxy"
}

func NewJsonRpcProxy(ctx context.Context, cfg config.Config) (*JsonRpcProxy, error) {
	jCfg := cfg.Scan.JsonRpc
	if cfg.JsonRpcProxy != nil {
		jCfg = cfg.JsonRpcProxy.JsonRpc
	}
	return &JsonRpcProxy{
		ctx: ctx,
		cfg: jCfg,
	}, nil
}
