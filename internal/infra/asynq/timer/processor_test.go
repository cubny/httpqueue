package timer

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cubny/httpqueue/internal/app/timer"
	timer2 "github.com/cubny/httpqueue/internal/infra/http/client/timer"
	mocks "github.com/cubny/httpqueue/internal/mocks/app/timer"
)

func TestNewProcessor(t *testing.T) {
	ctrl := gomock.NewController(t)

	tests := []struct {
		name       string
		service    timer.Service
		httpClient timer.HttpClient
		wantErr    bool
	}{
		{
			name:       "valid",
			service:    mocks.NewService(ctrl),
			httpClient: mocks.NewHttpClient(ctrl),
			wantErr:    false,
		},
		{
			name:       "no service",
			service:    nil,
			httpClient: mocks.NewHttpClient(ctrl),
			wantErr:    true,
		},
		{
			name:       "no http client",
			service:    mocks.NewService(ctrl),
			httpClient: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProcessor(tt.service, tt.httpClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProcessor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			assert.NotNil(t, got)
		})
	}
}

func TestProcessor_ProcessTask(t *testing.T) {
	type spec struct {
		mockFn             func(service *mocks.Service, httpClient *mocks.HttpClient)
		wantError          bool
		wantRetryableError bool
	}

	ctrl := gomock.NewController(t)

	testFn := func(s spec) func(t *testing.T) {
		return func(t *testing.T) {
			service := mocks.NewService(ctrl)
			httpClient := mocks.NewHttpClient(ctrl)

			p, err := NewProcessor(service, httpClient)
			require.NoError(t, err)

			payload := &Payload{TimerID: "1"}
			payloadBytes, err := json.Marshal(payload)
			require.NoError(t, err)

			task := asynq.NewTask(TypeName, payloadBytes)
			s.mockFn(service, httpClient)

			perr := p.ProcessTask(context.Background(), task)
			if (perr != nil) != s.wantError {
				t.Errorf("ProcessTask error = %v, wantErr %v", perr, s.wantError)
				return
			}
			if perr != nil {
				assert.Equal(t, s.wantRetryableError, !errors.Is(perr, asynq.SkipRetry))
			}
		}
	}

	t.Run("finds the timer, shoots webhook and archive the timer", testFn(spec{
		mockFn: func(service *mocks.Service, httpClient *mocks.HttpClient) {
			u, err := url.Parse("http://valid.url")
			require.NoError(t, err)

			foundTimer := &timer.Timer{ID: "1", URL: *u, FireAt: time.Now()}
			service.EXPECT().GetTimer(gomock.Any(), "1").Return(foundTimer, nil)
			service.EXPECT().ArchiveTimer(gomock.Any(), "1").Return(nil)

			httpClient.EXPECT().Shoot(gomock.Any(), foundTimer).Return(nil)
		},
		wantError: false,
	}))

	t.Run("finds the timer, shoots the webhook, get responded with non-retryable HTTP error", testFn(spec{
		mockFn: func(service *mocks.Service, httpClient *mocks.HttpClient) {
			u, err := url.Parse("http://valid.url")
			require.NoError(t, err)

			foundTimer := &timer.Timer{ID: "1", URL: *u, FireAt: time.Now()}
			service.EXPECT().GetTimer(gomock.Any(), "1").Return(foundTimer, nil)

			httpClient.EXPECT().Shoot(gomock.Any(), foundTimer).Return(assert.AnError)
		},
		wantError:          true,
		wantRetryableError: false,
	}))

	t.Run("finds the timer, shoots the webhook, get responded with retryable HTTP error", testFn(spec{
		mockFn: func(service *mocks.Service, httpClient *mocks.HttpClient) {
			u, err := url.Parse("http://valid.url")
			require.NoError(t, err)

			foundTimer := &timer.Timer{ID: "1", URL: *u, FireAt: time.Now()}
			service.EXPECT().GetTimer(gomock.Any(), "1").Return(foundTimer, nil)

			httpClient.EXPECT().Shoot(gomock.Any(), foundTimer).Return(timer2.ErrRetryableRequestFailure)
		},
		wantError:          true,
		wantRetryableError: true,
	}))

	t.Run("timer does not exist", testFn(spec{
		mockFn: func(service *mocks.Service, httpClient *mocks.HttpClient) {
			service.EXPECT().GetTimer(gomock.Any(), "1").Return(nil, timer.ErrTimerNotFound)
		},
		wantError:          true,
		wantRetryableError: false,
	}))

	t.Run("unknown issue with repo", testFn(spec{
		mockFn: func(service *mocks.Service, httpClient *mocks.HttpClient) {
			service.EXPECT().GetTimer(gomock.Any(), "1").Return(nil, assert.AnError)
		},
		wantError:          true,
		wantRetryableError: true,
	}))

	t.Run("timer is archived, does not call the webhook", testFn(spec{
		mockFn: func(service *mocks.Service, httpClient *mocks.HttpClient) {
			service.EXPECT().GetTimer(gomock.Any(), "1").Return(nil, timer.ErrTimerArchived)
		},
		wantError: false,
	}))
}
