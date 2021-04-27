package services

import (
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

	"OpenZeppelin/zephyr-node/proxy"
)

// TxProxyService is a proxy service for a given agent
type TxProxyService struct {
	ctx     context.Context
	servers []proxy.HttpProxyServer
}

type TxProxyConfig struct {
	AgentAddresses []string
	Keys           []string
}

func (t *TxProxyService) Get(idx int) proxy.HttpProxyServer {
	return t.servers[idx]
}

func (t *TxProxyService) Start() error {
	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		if t.ctx.Err() != nil {
			return t.ctx.Err()
		}
		log.Info("tx-logger tick")
	}
	return nil
}

func (t *TxProxyService) Stop() error {
	log.Infof("Stopping %s", t.Name())
	return nil
}

func (t *TxProxyService) Name() string {
	return "TxProxyService"
}

func NewTxProxyService(ctx context.Context, cfg TxProxyConfig) (*TxProxyService, error) {
	var servers []proxy.HttpProxyServer
	port := 35000
	if len(cfg.AgentAddresses) != len(cfg.Keys) {
		return nil, errors.New("number of agents must == number of keys")
	}
	for i, agt := range cfg.AgentAddresses {
		key := cfg.Keys[i]
		s, err := proxy.NewHttpProxy(ctx, proxy.HttpProxyConfig{
			Name:              agt,
			Key:               key,
			WhitelistPatterns: []string{".*"},
			Port:              port,
		})
		if err != nil {
			log.Errorf("error while creating proxy for %s", agt)
			return nil, err
		}
		servers = append(servers, s)
		port++
	}
	return &TxProxyService{
		ctx:     ctx,
		servers: servers,
	}, nil
}
