package redis

import (
	"errors"
	"github.com/forta-network/forta-node/config"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type Client interface {
	SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	Ping() *redis.StatusCmd
}

func tryClient(c Client) error {
	pRes := c.Ping()
	pong, err := pRes.Result()
	if err != nil {
		return err
	}
	if !strings.EqualFold(pong, "pong") {
		return errors.New("could not receive PONG from redis (connection issue?)")
	}
	return nil
}

func NewClusterClient(cfg config.RedisClusterConfig) (Client, error) {
	r := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:       cfg.Addresses,
		Password:    cfg.Password,
		MaxRetries:  5,
		DialTimeout: 10 * time.Second,
	})

	if err := tryClient(r); err != nil {
		log.WithError(err).Error("failed to connect to redis cluster")
		return nil, err
	}
	log.Info("initialized redis cluster connection")
	return r, nil
}

func NewClient(cfg config.RedisConfig) (Client, error) {
	r := redis.NewClient(&redis.Options{
		Addr:        cfg.Address,
		Password:    cfg.Password,
		DB:          cfg.DB,
		MaxRetries:  5,
		DialTimeout: 10 * time.Second,
	})

	if err := tryClient(r); err != nil {
		log.WithError(err).Error("failed to connect to redis")
		return nil, err
	}
	log.Info("initialized redis connection")
	return r, nil
}
