package timer

import (
	"encoding/json"
	"fmt"
	"github.com/cubny/cart/internal/app"
	"net/url"
	"time"
)

const (
	timerKeyFmt      = "timer-%s" // timer-<timer_id>
	schedulerLockFmt = "timer-schedule-%s"
	queueLockFmt     = "timer-queue-%s"
)

type redisTimer struct {
	ID            string `json:"id"`
	FireAtSecond  int64  `json:"fire_at"`
	TotalAttempts int    `json:"total_attempts"`
	URL           string `json:"url"`
}

func serializeKey(t redisTimer) string {
	return fmt.Sprintf(timerKeyFmt, t.ID)
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

func fromInternal(t *app.Timer) redisTimer {
	return redisTimer{
		ID:            t.ID,
		FireAtSecond:  t.FireAt.Unix(),
		URL:           t.URL.String(),
		TotalAttempts: t.TotalAttempts,
	}
}

func toInternal(r redisTimer) (*app.Timer, error) {
	URL, err := url.Parse(r.URL)
	if err != nil {
		return nil, err
	}

	return &app.Timer{
		ID:            r.ID,
		URL:           URL,
		FireAt:        time.Unix(r.FireAtSecond, 0),
		TotalAttempts: r.TotalAttempts,
	}, nil
}

func timerSchedulerLockKey(timerID string) string {
	return fmt.Sprintf(schedulerLockFmt, timerID)
}

func timerQueueLockKey(timerID string) string {
	return fmt.Sprintf(queueLockFmt, timerID)
}
