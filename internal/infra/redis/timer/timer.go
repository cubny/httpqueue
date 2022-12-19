package timer

import (
	"context"
	"fmt"
	"time"

	extRedis "github.com/go-redis/redis/v8"

	"github.com/cubny/httpqueue/internal/app/timer"
	"github.com/cubny/httpqueue/internal/config"
)

const (
	timerTaskQueueName   = "timerTaskQueue"
	timerBloomFilterName = "timerBloomFilter"
)

var (
	// ErrDeserialization indicates an error in the deserialization of the timer value.
	ErrDeserialization = fmt.Errorf("failed to deserialize timer value")
	// ErrInvalidURL is
	ErrInvalidURL = fmt.Errorf("timer has invalid URL")
)

// DB holds repo functionalities for timers.
type DB struct {
	redisClient extRedis.UniversalClient
	maxTTL      time.Duration
}

// NewDB constructs a DB
func NewDB(client extRedis.UniversalClient, cfg *config.DB) *DB {
	return &DB{
		redisClient: client,
		maxTTL:      time.Duration(cfg.TimerMaxTTLDays) * time.Hour * 24,
	}
}

// AddTimer follows the outbox pattern:
// 1. adds the timer object to the repo
// 2. adds the timer key to the outbox table (for the message relay to pick it up)
func (d *DB) AddTimer(ctx context.Context, timer *timer.Timer) error {
	internalTimer := fromInternal(timer)
	value := serializeValue(internalTimer)
	key := serializeKey(internalTimer.ID)

	pipe := d.redisClient.TxPipeline()

	pipe.Set(ctx, key, value, d.maxTTL)
	pipe.LPush(ctx, timerTaskQueueName, internalTimer.ID)

	_, err := pipe.Exec(ctx)
	return err
}

// Find looks up a key in the k/v DB.
// returns nil, nil when nothing found.
func (d *DB) Find(ctx context.Context, timerID string) (*timer.Timer, error) {
	key := serializeKey(timerID)
	val, err := d.redisClient.Get(ctx, key).Result()
	switch {
	case err == extRedis.Nil:
		return nil, nil
	case err != nil:
		return nil, err
	}

	dsTimer, err := deserializeValue(val)
	if err != nil {
		return nil, ErrDeserialization
	}

	t, err := toInternal(dsTimer)
	if err != nil {
		return nil, ErrInvalidURL
	}

	return t, nil
}

// IsArchived checks whether a timer is archived.
func (d *DB) IsArchived(ctx context.Context, timerID string) (bool, error) {
	return d.redisClient.Do(ctx, "BF.EXISTS", timerBloomFilterName, timerID).Bool()
}

// Archive a timer. for space efficiency it uses a Bloom Filter.
func (d *DB) Archive(ctx context.Context, timerID string) error {
	key := serializeKey(timerID)
	pipe := d.redisClient.TxPipeline()

	pipe.Del(ctx, key)
	pipe.Do(ctx, "BF.ADD", timerBloomFilterName, timerID)

	_, err := pipe.Exec(ctx)
	return err
}

// DequeueOutbox serves a message relay. It pops timers' keys out of the outbox queue.
func (d *DB) DequeueOutbox(ctx context.Context, batchSize int) ([]*timer.Timer, error) {
	timers := make([]*timer.Timer, 0, batchSize)
	for len(timers) < batchSize {
		timerID, err := d.redisClient.RPop(ctx, timerTaskQueueName).Result()
		switch {
		case err == extRedis.Nil:
			return timers, nil
		case err != nil:
			return nil, err
		}

		t, err := d.Find(ctx, timerID)
		switch {
		case err != nil:
			return nil, err
		case t == nil:
			continue
		}

		timers = append(timers, t)
	}

	return timers, nil
}
