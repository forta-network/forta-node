package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/elazarl/goproxy"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type HttpProxyServer interface {
	Start() error
	Stop() error
}

type httpProxy struct {
	ctx       context.Context
	config    HttpProxyConfig
	whitelist []*regexp.Regexp
	cancel    context.CancelFunc
}

type HttpProxyConfig struct {
	Name              string
	Key               string
	WhitelistPatterns []string
	Port              int
}

type proxyError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (p *proxyError) json() string {
	b, err := json.Marshal(p)
	if err != nil {
		log.Fatal("could not produce json from a proxyError", err)
		panic(err)
	}
	return string(b)
}

func newErrResponse(r *http.Request, code int, message string) *http.Response {
	pErr := &proxyError{code, message}
	return goproxy.NewResponse(r,
		"application/json", code,
		pErr.json())
}

func (p *httpProxy) authorize(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	authHeader := strings.Split(r.Header.Get("Authorization"), " ")
	if len(authHeader) == 0 {
		log.Errorf("agent %s did not send a key", p.config.Name)
		return r, newErrResponse(r, http.StatusUnauthorized, "key required")
	}
	if authHeader[0] != p.config.Key {
		log.Errorf("agent %s sent invalid key", p.config.Name)
		return r, newErrResponse(r, http.StatusUnauthorized, "invalid key")
	}
	url := r.URL.String()
	for _, re := range p.whitelist {
		if re.MatchString(url) {
			return r, nil
		}
	}
	return r, newErrResponse(r, http.StatusForbidden, "url not found in whitelist")
}

func (p *httpProxy) Key() string {
	return p.config.Key
}

func (p *httpProxy) Port() int {
	return p.config.Port
}

func (p *httpProxy) Start() error {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	proxy.OnRequest().DoFunc(p.authorize)

	ctx, cancel := context.WithCancel(p.ctx)
	p.cancel = cancel
	grp, ctx := errgroup.WithContext(p.ctx)

	grp.Go(func() error {
		return http.ListenAndServe(fmt.Sprintf(":%d", p.config.Port), proxy)
	})

	grp.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})

	return grp.Wait()
}

func (p *httpProxy) Stop() error {
	p.cancel()
	return nil
}

func NewHttpProxy(ctx context.Context, config HttpProxyConfig) (*httpProxy, error) {
	var whitelist []*regexp.Regexp
	for _, wp := range config.WhitelistPatterns {
		re, err := regexp.Compile(wp)
		if err != nil {
			log.Errorf("%s is not a valid regex: %v", wp, err)
			return nil, err
		}
		whitelist = append(whitelist, re)
	}
	ctx, cancel := context.WithCancel(ctx)
	return &httpProxy{
		ctx:       ctx,
		config:    config,
		whitelist: whitelist,
		cancel:    cancel,
	}, nil
}
