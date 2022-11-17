package store

import (
	"errors"
	rds "github.com/forta-network/forta-node/clients/redis"
	"github.com/forta-network/forta-node/config"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

const retryMax = 3

type DeduplicationStore interface {
	IsFirst(id string) (bool, error)
}

type dedupeStore struct {
	cfg config.Config
	mux sync.Mutex
	ttl time.Duration
	r   rds.Client
}

func (ds *dedupeStore) reconnect() error {
	c, err := newClient(ds.cfg)
	if err != nil {
		return err
	}
	ds.r = c
	return nil
}

func (ds *dedupeStore) IsFirst(id string) (bool, error) {
	ds.mux.Lock()
	defer ds.mux.Unlock()
	var err error
	for i := 0; i < retryMax; i++ {
		res := ds.r.SetNX(id, true, ds.ttl)
		if res.Err() == nil {
			return res.Result()
		}
		log.WithError(res.Err()).Error("error checking for duplicate on redis (reconnecting)")
		err = res.Err()
		if err := ds.reconnect(); err != nil {
			return false, err
		}
	}
	return false, err
}

func newClient(cfg config.Config) (rds.Client, error) {
	// regular redis
	if cfg.LocalModeConfig.DeduplicationConfig.Redis != nil {
		return rds.NewClient(*cfg.LocalModeConfig.DeduplicationConfig.Redis)
	}
	// redis cluster
	if cfg.LocalModeConfig.DeduplicationConfig.RedisCluster != nil {
		return rds.NewClusterClient(*cfg.LocalModeConfig.DeduplicationConfig.RedisCluster)
	}
	return nil, errors.New("redis or redisCluster is required in deduplicationConfig section")
}

func NewDeduplicationStore(cfg config.Config) (DeduplicationStore, error) {
	// no config, just return
	if cfg.LocalModeConfig.DeduplicationConfig == nil {
		return nil, nil
	}
	r, err := newClient(cfg)
	if err != nil {
		return nil, err
	}
	return &dedupeStore{
		cfg: cfg,
		r:   r,
		ttl: time.Duration(cfg.LocalModeConfig.DeduplicationConfig.TTLSeconds) * time.Second,
		mux: sync.Mutex{},
	}, nil
}
