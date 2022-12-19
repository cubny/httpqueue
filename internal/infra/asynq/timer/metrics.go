package timer

import (
	"errors"
	"github.com/hibiken/asynq"
	"github.com/prometheus/client_golang/prometheus"

	timerRepo "github.com/cubny/httpqueue/internal/infra/redis/timer"
)

var (
	relayErrorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "httpqueue",
			Subsystem: "relay",
			Name:      "error_counter",
			Help:      "Counter of relay errors",
		}, []string{"type", "reason"})
)

func init() {
	prometheus.MustRegister(relayErrorCount)
}

type dequeueError string

const (
	dequeueErrorURL          dequeueError = "url"
	dequeueErrorDeserialized dequeueError = "deserialize"
	dequeueErrorOthers       dequeueError = "others"
)

func (e dequeueError) String() string {
	return string(e)
}

func dequeueErrorReason(err error) dequeueError {
	switch err {
	case timerRepo.ErrInvalidURL:
		return dequeueErrorURL
	case timerRepo.ErrDeserialization:
		return dequeueErrorDeserialized
	default:
		return dequeueErrorOthers
	}
}

func relayDequeueErrorTypeInc(err error) {
	reason := dequeueErrorReason(err)
	relayErrorCount.With(prometheus.Labels{"type": "dequeue", "reason": reason.String()}).Inc()
}

type producerError string

const (
	producerErrorDuplicateTask  producerError = "duplicateTask"
	producerErrorTaskIDConflict producerError = "idConflict"
	producerErrorOthers         producerError = "others"
)

func (e producerError) String() string {
	return string(e)
}

func producerErrorReason(err error) producerError {
	switch {
	case errors.Is(err, asynq.ErrDuplicateTask):
		return producerErrorDuplicateTask
	case errors.Is(err, asynq.ErrTaskIDConflict):
		return producerErrorTaskIDConflict
	default:
		return producerErrorOthers
	}
}

func relayProducerErrorTypeInc(err error) {
	reason := producerErrorReason(err)
	relayErrorCount.With(prometheus.Labels{"type": "producer", "reason": reason.String()}).Inc()
}
