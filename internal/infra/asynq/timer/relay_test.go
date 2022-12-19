package timer

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cubny/httpqueue/internal/app/timer"
	"github.com/cubny/httpqueue/internal/config"
	mocks "github.com/cubny/httpqueue/internal/mocks/app/timer"
)

func TestNewRelay(t *testing.T) {
	ctrl := gomock.NewController(t)

	type args struct {
		cfg      *config.Relay
		outbox   timer.Outbox
		producer timer.Producer
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "valid",
			args: args{
				cfg:      &config.Relay{FrequencyMilliSeconds: 500},
				outbox:   mocks.NewOutbox(ctrl),
				producer: mocks.NewProducer(ctrl),
			},
			wantErr: assert.NoError,
		},
		{
			name: "nil outbox is unaccepted",
			args: args{
				cfg:      &config.Relay{FrequencyMilliSeconds: 500},
				outbox:   nil,
				producer: mocks.NewProducer(ctrl),
			},
			wantErr: assert.Error,
		},
		{
			name: "nil producer is unaccepted",
			args: args{
				cfg:      &config.Relay{FrequencyMilliSeconds: 500},
				outbox:   mocks.NewOutbox(ctrl),
				producer: nil,
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewRelay(tt.args.cfg, tt.args.outbox, tt.args.producer)
			if !tt.wantErr(t, err, fmt.Sprintf("NewRelay(%v, %v, %v)", tt.args.cfg, tt.args.outbox, tt.args.producer)) {
				return
			}
		})
	}
}

func TestRelay_Start(t *testing.T) {
	t.Run("with frequency of 100ms call dispatch at least 3 times in 350ms", func(t *testing.T) {

		cfg := &config.Relay{
			BatchSize:             1,
			FrequencyMilliSeconds: 100,
		}

		ctrl := gomock.NewController(t)

		now := time.Now()
		tms := []*timer.Timer{
			{ID: "1", URL: url.URL{Scheme: "http://", Host: "valid1.url"}, FireAt: now},
			{ID: "2", URL: url.URL{Scheme: "http://", Host: "valid2.url"}, FireAt: now},
			{ID: "3", URL: url.URL{Scheme: "http://", Host: "valid3.url"}, FireAt: now},
		}
		outbox := mocks.NewOutbox(ctrl)
		outbox.EXPECT().DequeueOutbox(gomock.Any(), gomock.Any()).Return(tms, nil).Times(3)

		producer := mocks.NewProducer(ctrl)
		producer.EXPECT().Send(gomock.Any(), gomock.Any()).Return(nil).MinTimes(3 * 3) // producer is called per timer

		r, err := NewRelay(cfg, outbox, producer)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())

		go r.Start(ctx)
		time.Sleep(350 * time.Millisecond)
		cancel()
	})
}

func TestRelay_dispatch(t *testing.T) {
	tests := []struct {
		name   string
		mockFn func(producer *mocks.Producer, outbox *mocks.Outbox)
	}{
		{
			name: "producer is called when some timers are dequeued",
			mockFn: func(producer *mocks.Producer, outbox *mocks.Outbox) {
				now := time.Now()
				tms := []*timer.Timer{
					{ID: "1", URL: url.URL{Scheme: "http://", Host: "valid1.url"}, FireAt: now},
					{ID: "2", URL: url.URL{Scheme: "http://", Host: "valid2.url"}, FireAt: now},
					{ID: "3", URL: url.URL{Scheme: "http://", Host: "valid3.url"}, FireAt: now},
				}
				outbox.EXPECT().DequeueOutbox(gomock.Any(), gomock.Any()).Return(tms, nil)
				producer.EXPECT().Send(gomock.Any(), gomock.Any()).Return(nil).Times(3)
			},
		},
		{
			name: "producer is not called when queue is empty",
			mockFn: func(producer *mocks.Producer, outbox *mocks.Outbox) {
				tms := make([]*timer.Timer, 0)
				outbox.EXPECT().DequeueOutbox(gomock.Any(), gomock.Any()).Return(tms, nil)
				producer.EXPECT().Send(gomock.Any(), gomock.Any()).Return(nil).Times(0)
			},
		},
		{
			name: "producer is not called when dequeue returns error",
			mockFn: func(producer *mocks.Producer, outbox *mocks.Outbox) {
				outbox.EXPECT().DequeueOutbox(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
				producer.EXPECT().Send(gomock.Any(), gomock.Any()).Return(nil).Times(0)
			},
		},
		{
			name: "producer errors are instrumented",
			mockFn: func(producer *mocks.Producer, outbox *mocks.Outbox) {
				now := time.Now()
				tms := []*timer.Timer{
					{ID: "1", URL: url.URL{Scheme: "http://", Host: "valid1.url"}, FireAt: now},
					{ID: "2", URL: url.URL{Scheme: "http://", Host: "valid2.url"}, FireAt: now},
					{ID: "3", URL: url.URL{Scheme: "http://", Host: "valid3.url"}, FireAt: now},
				}
				outbox.EXPECT().DequeueOutbox(gomock.Any(), gomock.Any()).Return(tms, nil)
				producer.EXPECT().Send(gomock.Any(), gomock.Any()).Return(assert.AnError).Times(3)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			producer := mocks.NewProducer(ctrl)
			outbox := mocks.NewOutbox(ctrl)
			cfg := &config.Relay{
				BatchSize:             1,
				FrequencyMilliSeconds: 100,
			}

			r, err := NewRelay(cfg, outbox, producer)
			require.NoError(t, err)

			// the expectations are asserted in mockFn
			tt.mockFn(producer, outbox)
			r.dispatch(context.Background())
		})
	}
}
