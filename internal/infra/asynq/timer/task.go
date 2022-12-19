package timer

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const TypeName = "timer:webhook"

type Payload struct {
	TimerID string
}

func NewTask(timerID string) (*asynq.Task, error) {
	payload, err := json.Marshal(Payload{TimerID: timerID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeName, payload), nil
}
