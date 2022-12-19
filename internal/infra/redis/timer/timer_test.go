package timer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cubny/httpqueue/internal/app/timer"
	"github.com/cubny/httpqueue/internal/config"
	mocks "github.com/cubny/httpqueue/internal/mocks/external/redis"
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
)

func TestDB_AddTimer(t *testing.T) {

	ctrl := gomock.NewController(t)
	cfg := &config.DB{TimerMaxTTLDays: 10}

	redisClient := mocks.NewRedisClient(ctrl)
	d := NewDB(redisClient, cfg)

	tm, err := timer.NewTimer("http://valid.url", 0, 0, 0)
	require.NoError(t, err)

	pipeliner := mocks.NewRedisPipeliner(ctrl)
	redisClient.EXPECT().TxPipeline().Return(pipeliner)

	pipeliner.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
	pipeliner.EXPECT().LPush(gomock.Any(), gomock.Any(), gomock.Any())
	pipeliner.EXPECT().Exec(gomock.Any()).Return(nil, nil)

	err = d.AddTimer(context.Background(), tm)
	require.NoError(t, err)
}

func TestDB_Find(t *testing.T) {
	aTimer, err := timer.NewTimer("http://valid.url", 0, 0, 0)
	require.NoError(t, err)

	aTimerInRedisTimer := redisTimer{
		ID:           aTimer.ID,
		FireAtSecond: aTimer.FireAt.Unix(),
		URL:          aTimer.URL.String(),
	}

	aTimerInRedisTimerJSON, err := json.Marshal(aTimerInRedisTimer)
	require.NoError(t, err)

	aTimerInRedisTimerJSONString := string(aTimerInRedisTimerJSON)
	fmt.Println(aTimerInRedisTimerJSONString)

	tests := []struct {
		name           string
		redisGetResult *redis.StringCmd
		want           *timer.Timer
		wantErr        bool
	}{
		{
			name:           "finds the timer in DB",
			redisGetResult: redis.NewStringResult(aTimerInRedisTimerJSONString, nil),
			want:           aTimer,
			wantErr:        false,
		},
		{
			name:           "does not exist",
			redisGetResult: redis.NewStringResult("", redis.Nil),
			want:           nil,
			wantErr:        false,
		},
		{
			name:           "Get returns error",
			redisGetResult: redis.NewStringResult("", assert.AnError),
			want:           nil,
			wantErr:        true,
		},
		{
			name:           "malformed value",
			redisGetResult: redis.NewStringResult("{malformed, json}", nil),
			want:           nil,
			wantErr:        true,
		},
		{
			name:           "data modified manually and invalid url is used",
			redisGetResult: redis.NewStringResult("{\"id\":\"1\",\"fire_at\":1631285661,\"url\":\"invalid.url/711aa1a5-9333-44da-9a7c-c5eae6edac44\"}", nil),
			want:           nil,
			wantErr:        true,
		},
	}

	ctrl := gomock.NewController(t)
	cfg := &config.DB{TimerMaxTTLDays: 10}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redisClient := mocks.NewRedisClient(ctrl)
			d := NewDB(redisClient, cfg)

			redisClient.EXPECT().Get(gomock.Any(), gomock.Any()).Return(tt.redisGetResult)

			got, err := d.Find(context.Background(), "1")
			if (err != nil) != tt.wantErr {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assertTimer(t, tt.want, got)
		})
	}
}

func assertTimer(t *testing.T, expected, actual *timer.Timer) {
	t.Helper()

	if expected == nil {
		assert.Nil(t, actual)
		return
	}

	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.URL.String(), actual.URL.String())
	assert.Equal(t, expected.FireAt.Second(), actual.FireAt.Second())
}

func TestDB_DequeueOutbox(t *testing.T) {
	ctrl := gomock.NewController(t)
	cfg := &config.DB{TimerMaxTTLDays: 10}

	tenTimers := make([]*timer.Timer, 0, 10)
	for i := 0; i < 10; i++ {
		aTimer, err := timer.NewTimer("http://valid.url", 0, 0, 0)
		require.NoError(t, err)
		tenTimers = append(tenTimers, aTimer)
	}

	tests := []struct {
		name      string
		mockFn    func(client *mocks.RedisClient)
		batchSize int
		want      []*timer.Timer
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "dequeue 10 timers",
			mockFn: func(client *mocks.RedisClient) {
				for i := 0; i < 10; i++ {
					rPopResult := redis.NewStringResult(tenTimers[i].ID, nil)
					client.EXPECT().RPop(gomock.Any(), gomock.Any()).Return(rPopResult)

					getResult := redis.NewStringResult(serializeValue(fromInternal(tenTimers[i])), nil)
					client.EXPECT().Get(gomock.Any(), serializeKey(tenTimers[i].ID)).Return(getResult)
				}

			},
			batchSize: 10,
			want:      tenTimers,
			wantErr:   assert.NoError,
		},
		{
			name: "returns early when queue is empty",
			mockFn: func(client *mocks.RedisClient) {
				rPopResult := redis.NewStringResult("", redis.Nil)
				client.EXPECT().RPop(gomock.Any(), gomock.Any()).Return(rPopResult)
			},
			batchSize: 10,
			want:      []*timer.Timer{},
			wantErr:   assert.NoError,
		},
		{
			name: "keeps popping when keys are not found",
			mockFn: func(client *mocks.RedisClient) {
				for i := 0; i < 10; i++ {
					rPopResult := redis.NewStringResult(tenTimers[i].ID, nil)
					client.EXPECT().RPop(gomock.Any(), gomock.Any()).Return(rPopResult)

					getResult := redis.NewStringResult("", redis.Nil)
					client.EXPECT().Get(gomock.Any(), serializeKey(tenTimers[i].ID)).Return(getResult)
				}

				for i := 0; i < 10; i++ {
					rPopResult := redis.NewStringResult(tenTimers[i].ID, nil)
					client.EXPECT().RPop(gomock.Any(), gomock.Any()).Return(rPopResult)

					getResult := redis.NewStringResult(serializeValue(fromInternal(tenTimers[i])), nil)
					client.EXPECT().Get(gomock.Any(), serializeKey(tenTimers[i].ID)).Return(getResult)
				}

			},
			batchSize: 10,
			want:      tenTimers,
			wantErr:   assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redisClient := mocks.NewRedisClient(ctrl)

			d := NewDB(redisClient, cfg)
			tt.mockFn(redisClient)

			got, err := d.DequeueOutbox(context.Background(), tt.batchSize)
			if !tt.wantErr(t, err, fmt.Sprintf("DequeueOutbox(ctx, %d)", tt.batchSize)) {
				return
			}

			assert.Equal(t, len(got), len(tt.want))

			if len(tt.want) == 0 {
				return
			}

			sort.Slice(got, func(i, j int) bool {
				return got[i].ID < got[j].ID
			})
			sort.Slice(tt.want, func(i, j int) bool {
				return tt.want[i].ID < tt.want[j].ID
			})

			for i := 0; i < tt.batchSize; i++ {
				assertTimer(t, got[i], tt.want[i])
			}
		})
	}
}
