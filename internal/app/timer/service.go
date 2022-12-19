package timer

import (
	"context"
	"errors"
)

var (
	ErrTimerNotFound = errors.New("timer not found")
	ErrTimerArchived = errors.New("timer is archived")
)

type ServiceImp struct {
	repo Repo
}

// NewService creates a new Service
func NewService(db Repo) (*ServiceImp, error) {
	return &ServiceImp{repo: db}, nil
}

// CreateTimer creates timer
func (s *ServiceImp) CreateTimer(ctx context.Context, cmd SetTimerCommand) (*Timer, error) {
	timer, err := NewTimerFromCommand(cmd)
	if err != nil {
		return nil, err
	}

	if err = s.repo.AddTimer(ctx, timer); err != nil {
		return timer, err
	}

	return timer, nil
}

// GetTimer fetches a timer by ID from the repo
func (s *ServiceImp) GetTimer(ctx context.Context, timerID string) (*Timer, error) {
	timer, err := s.repo.Find(ctx, timerID)
	switch {
	case err != nil:
		return nil, err
	case timer != nil:
		return timer, nil
	}

	// check whether the timer is archived
	archived, err := s.repo.IsArchived(ctx, timerID)
	switch {
	case err != nil:
		return nil, err
	case !archived:
		return nil, ErrTimerNotFound
	default:
		return nil, ErrTimerArchived
	}
}

func (s *ServiceImp) ArchiveTimer(ctx context.Context, timerID string) error {
	return s.repo.Archive(ctx, timerID)
}
