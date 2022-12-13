package service

import (
	"context"
	"errors"
	"github.com/cubny/cart/internal/app"
	"github.com/cubny/cart/internal/storage"
	"time"
)

var (
	ErrTimerNotFound = errors.New("timer not found")
)

const defaultTTL = time.Hour * 24 * 30 * 6 // 6months

// Service contains all the business logic of the shopping cart
type Service struct {
	repo Repo
}

// Repo provides the methods to CRUD persisted data.
type Repo interface {
	SaveTimer(ctx context.Context, timer *app.Timer) error
	GetTimer(ctx context.Context, id string) (*app.Timer, error)
	Close() error
}

// New creates a new Service
func New(db Repo) (*Service, error) {
	return &Service{repo: db}, nil
}

// CreateTimer creates timer
func (s *Service) CreateTimer(ctx context.Context, cmd *app.SetTimerCommand) (string, error) {
	timer, err := app.NewTimerFromCommand(cmd)
	if err != nil {
		return timer.ID, err
	}

	if err = s.repo.SaveTimer(ctx, timer); err != nil {
		return timer.ID, err
	}

	return timer.ID, nil
}

// GetTimer fetches a timer by ID from the repo
func (s *Service) GetTimer(ctx context.Context, timerID string) (*app.Timer, error) {
	// check the ownership of the cart
	timer, err := s.repo.GetTimer(ctx, timerID)
	switch {
	case err == storage.ErrRecordNotFound:
		return nil, ErrTimerNotFound
	case err != nil:
		return nil, err
	}

	return timer, nil
}
