package timer

import (
	"context"
	"fmt"
	"github.com/cubny/cart/internal/config"
	"time"

	"github.com/bsm/redislock"
	"github.com/cubny/cart/internal/app"
	infraredis "github.com/cubny/cart/internal/infra/redis"
	"github.com/go-redis/redis/v9"
)

const (
	timerSchedulerQueueName = "timerSortedQueue"
	timerTaskQueueName      = "timerTaskQueue"
)

var ErrTimerKeyNotFound = fmt.Errorf("timer key not found")

type Outbox struct {
	redisClient       *infraredis.Client
	distributedLocker *redislock.Client
	maxTTL            time.Duration
	ticker            *time.Ticker
	batchSize         int
}

func NewOutbox(ctx context.Context, client *infraredis.Client, distributedLocker *redislock.Client, cfg *config.Queue, errChan chan<- error) *Outbox {
	r := &Outbox{
		redisClient:       client,
		distributedLocker: distributedLocker,
		maxTTL:            time.Duration(cfg.MaxTTLDays) * time.Hour * 24,
		batchSize:         cfg.BatchSize,
		ticker:            time.NewTicker(time.Duration(cfg.SchedulerFrequencyMilliSeconds) * time.Millisecond),
	}

	go r.enqueueDuesOnIntervals(ctx, errChan)

	return r
}

//func (r *Queue) Dequeue(ctx context.Context) ([]*Task, error) {
//	tasks := make([]*Task, r.batchSize)
//	for i := 0; i < r.batchSize; {
//		timerID, err := r.redisClient.RPop(ctx, timerTaskQueueName).Result()
//		if err != nil {
//			return nil, err
//		}
//
//		val, err := r.redisClient.Get(ctx, timerID).Result()
//		switch {
//		case err == redis.Nil:
//			continue
//		case err != nil:
//			return nil, err
//		}
//
//		dsTimer, err := deserializeValue(val)
//		if err != nil {
//			return nil, err
//		}
//
//		timer, err := toInternal(dsTimer)
//		if err != nil {
//			return nil, err
//		}
//
//		task := &Task{
//			Timer: timer,
//			OnFailure: func(ctx context.Context) error {
//				return r.redisClient.RPush(ctx, timerTaskQueueName, timerID).Err()
//			},
//			OnFinish: func(ctx context.Context) error {
//				pipe := r.redisClient.TxPipeline()
//				pipe.Del(ctx, timer.ID)
//				// add to bloom filter
//				_, err := pipe.Exec(ctx)
//				return err
//			},
//		}
//
//		tasks = append(tasks, task)
//		i++
//	}
//
//	return tasks, nil
//
//}

//func (r *Queue) enqueueDuesOnIntervals(ctx context.Context, errChan chan<- error) {
//	for {
//		if _, ok := <-r.ticker.C; ok {
//			return
//		}
//
//		if err := r.enqueueDues(ctx); err != nil {
//			errChan <- err
//		}
//	}
//}
//
//func (r *Queue) Close() error {
//	r.ticker.Stop()
//	return r.redisClient.Close()
//}
//
//// acquireLock gets the lock for the item with key
//func (r *Queue) acquireLock(ctx context.Context, key string) *redislock.Lock {
//	lock, err := r.distributedLocker.Obtain(ctx, key, 1000*time.Millisecond, nil)
//	if err != nil {
//		return nil
//	}
//	return lock
//}

func (r *Outbox) Store(ctx context.Context, timer *app.Timer) error {
	internalTimer := fromInternal(timer)
	value := serializeValue(internalTimer)
	key := serializeKey(internalTimer)

	pipe := r.redisClient.TxPipeline()

	pipe.Set(ctx, key, value, r.maxTTL)
	pipe.LPush(ctx, timerTaskQueueName, value)
	//pipe.ZAdd(ctx, timerSchedulerQueueName, redis.Z{Member: key, Score: timer.DelayFromNowSeconds()})

	_, err := pipe.Exec(ctx)
	return err
}

func (r *Outbox) Find(ctx context.Context, timerID string) (*app.Timer, error) {
	val, err := r.redisClient.Get(ctx, timerID).Result()
	switch {
	case err == redis.Nil:
		return nil, ErrTimerKeyNotFound
	case err != nil:
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

	return timer, nil
}

func (r *Outbox) Dispatch(ctx context.Context) error {
	for i := 0; i < r.batchSize; {
		timerID, err := r.redisClient.RPop(ctx, timerTaskQueueName).Result()
		if err != nil {
			return err
		}

		timer, err := r.Find(ctx, timerID)
		switch {
		case err == ErrTimerKeyNotFound:
			continue
		case err != nil:
			return err
		}

		r.queue.Enqueue(timer)
		return nil
	}
}
