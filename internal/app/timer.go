package app

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/url"
	"time"
)

var ErrFireAtInPast = errors.New("time in the past")

type Timer struct {
	ID            string
	URL           *url.URL
	FireAt        time.Time
	TotalAttempts int
}

func NewTimerFromCommand(cmd *SetTimerCommand) (*Timer, error) {
	return NewTimer(cmd.URLRaw, time.Duration(cmd.Hours), time.Duration(cmd.Minutes), time.Duration(cmd.Seconds))
}

func NewTimer(rawURL string, hours, minutes, seconds time.Duration) (*Timer, error) {
	validURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid rawURL provided: %w", err)
	}

	now := time.Now()
	duration := hours*time.Hour + minutes*time.Minute + seconds*time.Second

	firedAt := now.Add(duration)
	if firedAt.Before(now) {
		return nil, ErrFireAtInPast
	}

	return &Timer{
		ID:     uuid.NewString(),
		URL:    validURL,
		FireAt: firedAt,
	}, nil
}

func (t *Timer) DelayFromNowSeconds() float64 {
	return t.FireAt.Sub(time.Now()).Seconds()
}
