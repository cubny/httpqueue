package timer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"

	"github.com/cubny/httpqueue/internal/app/timer"
	internalHttpClient "github.com/cubny/httpqueue/internal/infra/http/client/timer"
)

// Processor is the logical processor of the Timer tasks.
// it receives async.Task with the type TypeName and makes the HTTP call to the given Timer URL, using an internal
// HTTP Client that is aware of retryability of failed requests. In that way if the failed request is retryable
// Processor let the worker retry the same task later, otherwise it fails the task permanently.
type Processor struct {
	service    timer.Service
	httpClient timer.HttpClient
}

// NewProcessor constructs a Processor
func NewProcessor(service timer.Service, httpClient timer.HttpClient) (*Processor, error) {
	if service == nil {
		return nil, errors.New("service is not set up")
	}

	if httpClient == nil {
		return nil, errors.New("httpClient is not set up")
	}

	return &Processor{
		service:    service,
		httpClient: httpClient,
	}, nil
}

// ProcessTask processes task
func (p *Processor) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var payload Payload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		// TODO: remove from DB
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	t, err := p.service.GetTimer(ctx, payload.TimerID)
	switch {
	case err == timer.ErrTimerNotFound:
		return fmt.Errorf("timer does not exist: %v: %w", err, asynq.SkipRetry)
	case err == timer.ErrTimerArchived:
		return nil
	case err != nil:
		// this could be due to unavailability of the underlying DB
		return fmt.Errorf("failed to find the timer in storage: %v", err)
	}

	logrus.WithFields(logrus.Fields{"timer": t}).Debug("making HTTP call")
	err = p.httpClient.Shoot(ctx, t)
	switch {
	case errors.Is(err, internalHttpClient.ErrRetryableRequestFailure):
		return fmt.Errorf("temporarliy failed to call the timer URL: %w", err)
	case err != nil:
		// TODO: publish to DLQ for troubleshooting and remove from the main queue
		return fmt.Errorf("permenantly failed to call the timer URL: %v: %w", err, asynq.SkipRetry)
	default:
		return p.service.ArchiveTimer(ctx, payload.TimerID)
	}
}
