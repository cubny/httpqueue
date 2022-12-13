package config

import "time"

// Redis is the configuration for Redis.
type Redis struct {
	URL      string `seed:"" env:"REDIS_URL"`
	Password string `seed:"" env:"REDIS_PASSWORD"`
	DB       int64  `seed:"0" env:"REDIS_DB"`

	IsRedisCluster bool `seed:"true" env:"REDIS_IS_CLUSTER"`
	IsTLSEnabled   bool `seed:"true" env:"REDIS_IS_TLS"`

	MaxRetries   int64         `seed:"2" env:"REDIS_MAX_RETRIES"`
	ReadTimeout  time.Duration `seed:"1s" env:"REDIS_READ_TIMEOUT"`
	WriteTimeout time.Duration `seed:"1s" env:"REDIS_WRITE_TIMEOUT"`
	DialTimeout  time.Duration `seed:"2s" env:"REDIS_DIAL_TIMEOUT"`
}
