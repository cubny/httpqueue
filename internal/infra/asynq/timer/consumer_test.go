package timer

import (
	mocks2 "github.com/cubny/httpqueue/internal/mocks/external/asynq"
	mocks "github.com/cubny/httpqueue/internal/mocks/external/redis"
	"github.com/golang/mock/gomock"
	"github.com/hibiken/asynq"
	"testing"
)

func TestNewConsumer(t *testing.T) {
	ctrl := gomock.NewController(t)

	createAsynqServerFn := func() *asynq.Server {
		redisClient := mocks.NewRedisClient(ctrl)
		redisConnOpt := mocks2.NewRedisConnOpt(ctrl)
		redisConnOpt.EXPECT().MakeRedisClient().Return(redisClient)
		return asynq.NewServer(redisConnOpt, asynq.Config{})
	}

	tests := []struct {
		server    *asynq.Server
		processor asynq.Handler
		name      string
		wantErr   bool
	}{
		{
			name:      "valid",
			server:    createAsynqServerFn(),
			processor: mocks2.NewHandler(ctrl),
			wantErr:   false,
		},
		{
			name:      "no server",
			server:    nil,
			processor: mocks2.NewHandler(ctrl),
			wantErr:   true,
		},
		{
			name:      "no processor",
			server:    createAsynqServerFn(),
			processor: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewConsumer(tt.server, tt.processor)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConsumer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
