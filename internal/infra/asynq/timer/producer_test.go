package timer

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cubny/httpqueue/internal/app/timer"
	"github.com/cubny/httpqueue/internal/config"
	mocks2 "github.com/cubny/httpqueue/internal/mocks/external/asynq"
)

func TestProducer_Send(t *testing.T) {

	aTimer, err := timer.NewTimer("http://valid.url", 1, 1, 1)
	require.NoError(t, err)

	tests := []struct {
		name               string
		timer              *timer.Timer
		enqueueCalledTimes int
		enqueueErr         error
		wantErr            assert.ErrorAssertionFunc
	}{
		{
			name:               "valid timer, no error from broker",
			timer:              aTimer,
			enqueueErr:         nil,
			enqueueCalledTimes: 1,
			wantErr:            assert.NoError,
		},
		{
			name:               "nil timer will err",
			timer:              nil,
			enqueueErr:         nil,
			enqueueCalledTimes: 0,
			wantErr:            assert.Error,
		},
	}

	cfg := &config.Producer{}
	ctrl := gomock.NewController(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			broker := mocks2.NewBroker(ctrl)

			broker.
				EXPECT().
				EnqueueContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, tt.enqueueErr).
				Times(tt.enqueueCalledTimes)

			p := NewProducer(broker, cfg)

			tt.wantErr(t, p.Send(context.Background(), tt.timer), fmt.Sprintf("Send(ctx, %v)", tt.timer))
		})
	}
}
