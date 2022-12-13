package timer

import (
	"context"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	"github.com/cubny/cart/internal/app"
	infraredis "github.com/cubny/cart/internal/infra/redis"
	"github.com/go-redis/redis/v9"
)

const (
	timerSchedulerQueueName = "timerSortedQueue"
	timerTaskQueueName      = "timerTaskQueue"
	enqueueFrequency        = 500 * time.Millisecond
)

type Queue struct {
	redisClient       *infraredis.Client
	distributedLocker *redislock.Client
	maxTTL            time.Duration
	ticker            *time.Ticker
	batchSize         int
}

type Task struct {
	Timer     *app.Timer
	OnFailure func(ctx context.Context) error
	OnFinish  func(ctx context.Context) error
}

func NewQueue(ctx context.Context, client *infraredis.Client, distributedLocker *redislock.Client, batchSize int, maxTTL time.Duration, errChan chan<- error) *Queue {
	r := &Queue{
		redisClient:       client,
		distributedLocker: distributedLocker,
		maxTTL:            maxTTL,
		ticker:            time.NewTicker(enqueueFrequency),
		batchSize:         batchSize,
	}

	go r.enqueueDuesOnIntervals(ctx, errChan)

	return r
}

func (r *Queue) Dequeue(ctx context.Context) ([]Task, error) {
	tasks := make([]Task, r.batchSize)
	for i := 0; i < r.batchSize; {
		timerID, err := r.redisClient.RPop(ctx, timerTaskQueueName).Result()
		if err != nil {
			return nil, err
		}

		lock := r.acquireLock(ctx, timerQueueLockKey(timerID))
		if lock == nil {
			continue
		}

		val, err := r.redisClient.Get(ctx, timerID).Result()
		if err != nil {
			return nil, err
		}

		dsTimer, err := deserializeValue(val)
		if err != nil {
			return nil, err
		}

		timer, err := toInternal(dsTimer)
		if err != nil {
			return nil, err
		}

		task := &Task{
			Timer: timer,
			OnFailure: func(ctx context.Context) error {
				r.redisClient.RPush(ctx)
				return nil
			},
			OnFinish: func(ctx context.Context) error {
				pipe := r.redisClient.TxPipeline()
				pipe.Del(ctx, timer.ID)
				// add to bloom filter
				_, err := pipe.Exec(ctx)
				if err != nil {
					return err
				}
				return lock.Release(ctx)
			},
		}

		i++
	}

}

func (r *Queue) enqueueDuesOnIntervals(ctx context.Context, errChan chan<- error) {
	for {
		if _, ok := <-r.ticker.C; ok {
			return
		}

		if err := r.enqueueDues(ctx); err != nil {
			errChan <- err
		}
	}
}

func (r *Queue) Close() error {
	r.ticker.Stop()
	return r.redisClient.Close()
}

// acquireLock gets the lock for the item with key
func (r *Queue) acquireLock(ctx context.Context, key string) *redislock.Lock {
	lock, err := r.distributedLocker.Obtain(ctx, key, 1000*time.Millisecond, nil)
	if err != nil {
		return nil
	}
	return lock
}

func (r *Queue) EnqueueLater(ctx context.Context, timer *app.Timer) error {
	internalTimer := fromInternal(timer)
	value := serializeValue(internalTimer)
	key := serializeKey(internalTimer)

	pipe := r.redisClient.TxPipeline()

	pipe.Set(ctx, key, value, r.maxTTL)
	pipe.ZAdd(ctx, timerSchedulerQueueName, &redis.Z{Member: key, Score: timer.DelayFromNowSeconds()})

	_, err := pipe.Exec(ctx)
	return err
}

func (r *Queue) enqueueDues(ctx context.Context) error {
	start := int64(0)
	for i := r.batchSize; i >= 0; {
		values, err := r.redisClient.ZRangeWithScores(ctx, timerSchedulerQueueName, start, start).Result()
		if err != nil {
			return fmt.Errorf("failed to get range from zset: %w", err)
		}
		if len(values) == 0 || values[0].Score > float64(time.Now().Unix()) {
			break
		}

		key := values[0].Member.(string)
		lock := r.acquireLock(ctx, timerSchedulerLockKey(key))
		if lock == nil {
			start++
			continue
		}

		pipe := r.redisClient.TxPipeline()
		pipe.RPush(ctx, timerTaskQueueName, key)
		pipe.ZRem(ctx, timerSchedulerQueueName, key)
		if _, err = pipe.Exec(ctx); err != nil {
			return fmt.Errorf("failed to enqueue scheduled timers")
		}

		if err = lock.Release(ctx); err != nil {
			return fmt.Errorf("failed to release the scheulder lock: %w", err)
		}

		start++
		i--
	}
	return nil
}
