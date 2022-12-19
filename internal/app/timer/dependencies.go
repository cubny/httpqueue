package timer

import (
	"context"
)

type Outbox interface {
	DequeueOutbox(ctx context.Context, batchSize int) ([]*Timer, error)
}

type Repo interface {
	Find(ctx context.Context, timerID string) (*Timer, error)
	AddTimer(ctx context.Context, timer *Timer) error
	Archive(ctx context.Context, timerID string) error
	IsArchived(ctx context.Context, timerID string) (bool, error)
}

type Producer interface {
	Send(ctx context.Context, timer *Timer) error
}

type HttpClient interface {
	Shoot(ctx context.Context, timer *Timer) error
}

// Service holds all the business logic
type Service interface {
	CreateTimer(ctx context.Context, cmd SetTimerCommand) (*Timer, error)
	GetTimer(ctx context.Context, timerID string) (*Timer, error)
	ArchiveTimer(ctx context.Context, timerID string) error
}
