package timer

import (
	"encoding/json"
	"fmt"
	"github.com/cubny/httpqueue/internal/app/timer"
	"net/url"
	"time"
)

const (
	timerKeyFmt = "timer-%s" // timer-<timer_id>
)

type redisTimer struct {
	ID           string `json:"id"`
	FireAtSecond int64  `json:"fire_at"`
	URL          string `json:"url"`
}

func serializeKey(timerID string) string {
	return fmt.Sprintf(timerKeyFmt, timerID)
}

func serializeValue(t redisTimer) string {
	// ignore the error because we know the model is valid (doesn't contain channels, cyclic data structures, etc.)
	bytes, _ := json.Marshal(t)
	return string(bytes)
}

func deserializeValue(str string) (redisTimer, error) {
	t := redisTimer{}
	err := json.Unmarshal([]byte(str), &t)
	return t, err
}

func fromInternal(t *timer.Timer) redisTimer {
	return redisTimer{
		ID:           t.ID,
		FireAtSecond: t.FireAt.Unix(),
		URL:          t.URL.String(),
	}
}

func toInternal(r redisTimer) (*timer.Timer, error) {
	URL, err := url.ParseRequestURI(r.URL)
	if err != nil {
		return nil, err
	}

	return &timer.Timer{
		ID:     r.ID,
		URL:    *URL,
		FireAt: time.Unix(r.FireAtSecond, 0),
	}, nil
}
