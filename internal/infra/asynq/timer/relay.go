package timer

import (
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/cubny/httpqueue/internal/app/timer"
	"github.com/cubny/httpqueue/internal/config"
)

// Relay is responsible for relaying messages from the outbox to the producer
type Relay struct {
	outbox    timer.Outbox
	producer  timer.Producer
	ticker    *time.Ticker
	batchSize int
}

// NewRelay constructs a Relay.
func NewRelay(cfg *config.Relay, outbox timer.Outbox, producer timer.Producer) (*Relay, error) {
	if outbox == nil {
		return nil, errors.New("outbox is not set up")
	}

	if producer == nil {
		return nil, errors.New("producer is not set up")
	}

	relay := &Relay{
		outbox:    outbox,
		producer:  producer,
		ticker:    time.NewTicker(time.Duration(cfg.FrequencyMilliSeconds) * time.Millisecond),
		batchSize: cfg.BatchSize,
	}

	return relay, nil

}

// Start makes periodic dispatch based on the ticker frequency.
func (r *Relay) Start(ctx context.Context) {
	for {
		select {
		case <-r.ticker.C:
			r.dispatch(ctx)
		case <-ctx.Done():
			r.Stop()
			return
		}
	}
}

// Stop the ticker.
func (r *Relay) Stop() {
	r.ticker.Stop()
}

// dispatch dequeues the outbox and send them to the workers via the producer.
func (r *Relay) dispatch(ctx context.Context) {
	timers, err := r.outbox.DequeueOutbox(ctx, r.batchSize)
	if err != nil {
		relayDequeueErrorTypeInc(err)
		log.WithContext(ctx).Errorf("unable to dequeue the outbox, %v", err)
	}

	for _, t := range timers {
		if err = r.producer.Send(ctx, t); err != nil {
			relayProducerErrorTypeInc(err)
			log.WithContext(ctx).Errorf("unable to relay messages to the producer, %v", err)
		}
	}
}
