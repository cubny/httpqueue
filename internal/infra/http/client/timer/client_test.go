package timer

import (
	"context"
	"github.com/cubny/httpqueue/internal/app/timer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Shoot(t *testing.T) {
	tests := []struct {
		name               string
		serverStatusCode   int
		wantRetryableError bool
	}{
		{
			name:               "returns retryable error for 500",
			serverStatusCode:   http.StatusInternalServerError,
			wantRetryableError: true,
		},
		{
			name:               "returns retryable error for 502",
			serverStatusCode:   http.StatusBadGateway,
			wantRetryableError: true,
		},
		{
			name:               "returns retryable error for 429",
			serverStatusCode:   http.StatusTooManyRequests,
			wantRetryableError: true,
		},
		{
			name:               "returns non-retryable error for 404",
			serverStatusCode:   http.StatusNotFound,
			wantRetryableError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.serverStatusCode)
			}))
			defer ts.Close()

			tm, err := timer.NewTimer(ts.URL, 0, 0, 0)
			require.NoError(t, err)

			client := NewClient()
			err = client.Shoot(context.Background(), tm)

			if tt.wantRetryableError {
				assert.ErrorIs(t, err, ErrRetryableRequestFailure)
			} else {
				assert.NotErrorIs(t, err, ErrRetryableRequestFailure)
			}
		})
	}
}
