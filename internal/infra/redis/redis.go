package redis

import (
	"crypto/tls"
	"github.com/cubny/httpqueue/internal/config"
	"github.com/go-redis/redis/v8"
)

// Client is a wrapper around a default Redis client with some helper methods.
type Client struct {
	redis.UniversalClient
}

// NewRedis returns a new client for Redis with proper configuration.
func NewRedis(cfg *config.Redis) *Client {
	return &Client{UniversalClient: redisClient(cfg)}
}

// MakeRedisClient returns the initiated UniversalClient
func (c Client) MakeRedisClient() interface{} {
	return c.UniversalClient
}

func redisClient(cfg *config.Redis) redis.UniversalClient {
	var client redis.UniversalClient

	var tlsCfg *tls.Config
	if cfg.IsTLSEnabled {
		tlsCfg = &tls.Config{}
	}

	if cfg.IsRedisCluster {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        []string{cfg.URL},
			Password:     cfg.Password,
			MaxRetries:   int(cfg.MaxRetries),
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			DialTimeout:  cfg.DialTimeout,
			TLSConfig:    tlsCfg,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:         cfg.URL,
			Password:     cfg.Password,
			MaxRetries:   int(cfg.MaxRetries),
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			DialTimeout:  cfg.DialTimeout,
			TLSConfig:    tlsCfg,
		})
	}
	return client
}
