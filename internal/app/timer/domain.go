package timer

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/url"
	"time"
)

var ErrFireAtInPast = errors.New("time in the past")

type Timer struct {
	ID     string
	URL    url.URL
	FireAt time.Time
}

func NewTimerFromCommand(cmd SetTimerCommand) (*Timer, error) {
	return NewTimer(cmd.URLRaw, time.Duration(cmd.Hours), time.Duration(cmd.Minutes), time.Duration(cmd.Seconds))
}

func NewTimer(rawURL string, hours, minutes, seconds time.Duration) (*Timer, error) {
	id := uuid.NewString()

	validURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid rawURL provided: %w", err)
	}

	validURL = validURL.JoinPath(id)

	//TODO: use UTC or local time and document the decision
	now := time.Now()
	duration := hours*time.Hour + minutes*time.Minute + seconds*time.Second

	firedAt := now.Add(duration)
	if firedAt.Before(now) {
		return nil, ErrFireAtInPast
	}

	return &Timer{
		ID:     id,
		URL:    *validURL,
		FireAt: firedAt,
	}, nil
}

func (t *Timer) DelayFromNowSeconds() float64 {
	return t.FireAt.Sub(time.Now()).Seconds()
}

func (t *Timer) Validate() error {
	if t == nil {
		return fmt.Errorf("nil timer is invalid")
	}

	if t.ID == "" {
		return fmt.Errorf("ID cannot be empty")
	}

	if t.URL.String() == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	if t.FireAt.IsZero() {
		return fmt.Errorf("FireAt cannot be zero")
	}

	return nil
}
