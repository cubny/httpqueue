package timer

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"

	"github.com/cubny/httpqueue/internal/app/timer"
	"github.com/cubny/httpqueue/internal/config"
)

type Producer struct {
	broker   Broker
	maxRetry int
}

type Broker interface {
	EnqueueContext(ctx context.Context, task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)
}

func NewProducer(broker Broker, cfg *config.Producer) *Producer {
	return &Producer{
		broker:   broker,
		maxRetry: cfg.MaxRetry,
	}
}

func (p *Producer) Send(ctx context.Context, timer *timer.Timer) error {
	if err := timer.Validate(); err != nil {
		return fmt.Errorf("timer is invalid: %w", err)
	}

	task, err := NewTask(timer.ID)
	if err != nil {
		return err
	}

	_, err = p.broker.EnqueueContext(ctx, task, asynq.MaxRetry(p.maxRetry), asynq.ProcessAt(timer.FireAt))
	return err
}
