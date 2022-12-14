package redis

import (
	"crypto/tls"
	"github.com/cubny/cart/internal/config"
	"github.com/go-redis/redis/v9"
)

// Client is a wrapper around a default Redis client with some helper methods.
type Client struct {
	redis.UniversalClient
}

// NewRedis returns a new client for Redis with proper configuration.
func NewRedis(cfg *config.Redis) *Client {
	return &Client{UniversalClient: redisClient(cfg)}
}

func redisClient(cfg *config.Redis) redis.UniversalClient {
	var client redis.UniversalClient

	var tlsCfg *tls.Config
	if cfg.IsTLSEnabled.Get() {
		tlsCfg = &tls.Config{}
	}

	if cfg.IsRedisCluster.Get() {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        []string{cfg.URL.Get()},
			Password:     cfg.Password.Get(),
			MaxRetries:   int(cfg.MaxRetries.Get()),
			ReadTimeout:  cfg.ReadTimeout.Get(),
			WriteTimeout: cfg.WriteTimeout.Get(),
			DialTimeout:  cfg.DialTimeout.Get(),
			TLSConfig:    tlsCfg,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:         cfg.URL.Get(),
			Password:     cfg.Password.Get(),
			MaxRetries:   int(cfg.MaxRetries.Get()),
			ReadTimeout:  cfg.ReadTimeout.Get(),
			WriteTimeout: cfg.WriteTimeout.Get(),
			DialTimeout:  cfg.DialTimeout.Get(),
			TLSConfig:    tlsCfg,
		})
	}
	return client
}
