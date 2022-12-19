package config

import (
	"context"
	"time"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	// AppMode could be one of: all, api, relay, workers.
	//- all: all components are active in one process. (default)
	//- api: only activate the API server
	//- relay: only activate the relay component (it reads from the outbox table and publish it in the queue)
	//- workers: only activate the workers component. consumers of the queue.
	AppMode string `env:"APP_MODE,default=all"`
	// ConsumerConcurrency number of concurrent consumers.
	// Note: the words consumers and workers are used interchangeably
	ConsumerConcurrency int `env:"CONSUMER_CONCURRENCY,default=10"`

	HTTP     HTTP
	DB       DB
	Producer Producer
	Redis    Redis
	Relay    Relay
}

type HTTP struct {
	Port        int `env:"HTTP_PORT"`
	MetricsPort int `env:"HTTP_METRICS_PORT"`
}

type DB struct {
	// TimerMaxTTLDays indicates the TTL of the timers record in the DB
	TimerMaxTTLDays int `env:"DB_TIMER_MAX_TTL_DAYS,default=180"`
}

type Producer struct {
	// MaxRetry indicates how many times the timer.Timer webhook is allowed to be called at maximum in case of failure.
	MaxRetry int `env:"PRODUCER_MAX_RETRY,default=10"`
}

// Redis is the configuration for Redis.
type Redis struct {
	URL      string `env:"REDIS_URL"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int64  `env:"REDIS_DB,default=0"`

	IsRedisCluster bool `env:"REDIS_IS_CLUSTER,default=true"`
	IsTLSEnabled   bool `env:"REDIS_IS_TLS,default=true"`

	MaxRetries   int64         `env:"REDIS_MAX_RETRIES,default=2"`
	ReadTimeout  time.Duration `env:"REDIS_READ_TIMEOUT,default=1s"`
	WriteTimeout time.Duration `env:"REDIS_WRITE_TIMEOUT,default=1s"`
	DialTimeout  time.Duration `env:"REDIS_DIAL_TIMEOUT,default=2s"`
}

type Relay struct {
	// BatchSize indicates the number of items the relay dequeues from the Outbox table in each iteration.
	BatchSize int `env:"RELAY_BATCH_SIZE,default=1"`
	// FrequencyMilliSeconds is the intervals between outbox dequeue.
	FrequencyMilliSeconds int `env:"RELAY_FREQUENCY_MILLI_SECONDS,default=500"`
}

// New constructs the config.
// variables are populated using the envars and default values.
func New(ctx context.Context) (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process(ctx, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
